package aggregator

import (
	"cmp"
	"slices"
	"time"

	"gitstat/internal/model"
)

type RepoAgg struct {
	EarliestDate         time.Time
	EarliestCommitAuthor string
	Contributors         []model.ContributorStat
	ChartCalendar        []model.CalendarPoint
	ChartCumulative      []model.CumulativePoint
	ChartHourly          []model.HourlyPoint
}

type RepoAccumulator struct {
	repoPath       string
	earliestDate   time.Time
	earliestAuthor string
	contributors   map[string]*contributorAcc
	calendar       map[string]int
	hourly         map[int]int
}

type contributorAcc struct {
	model.ContributorStat
	lastSeen time.Time
}

func NewRepoAccumulator(repoPath string) *RepoAccumulator {
	return &RepoAccumulator{
		repoPath:     repoPath,
		contributors: make(map[string]*contributorAcc),
		calendar:     make(map[string]int),
		hourly:       make(map[int]int),
	}
}

func (ra *RepoAccumulator) Add(c *model.Commit) {
	if ra.earliestDate.IsZero() || c.Date.Before(ra.earliestDate) {
		ra.earliestDate = c.Date
		ra.earliestAuthor = c.Author
	}

	ca, ok := ra.contributors[c.Email]
	if !ok {
		ca = &contributorAcc{}
		ca.Author = c.Author
		ca.Email = c.Email
		ra.contributors[c.Email] = ca
	}
	ca.CommitCount++
	ca.Additions += c.Additions
	ca.Deletions += c.Deletions
	if c.Date.After(ca.lastSeen) {
		ca.lastSeen = c.Date
		ca.LastCommitDate = c.Date.Format("2006-01-02 15:04:05")
	}

	dateKey := c.Date.Format("2006-01-02")
	ra.calendar[dateKey]++

	hour := c.Date.Hour()
	ra.hourly[hour]++
}

func (ra *RepoAccumulator) Build() *RepoAgg {
	agg := &RepoAgg{
		EarliestDate:         ra.earliestDate,
		EarliestCommitAuthor: ra.earliestAuthor,
	}

	agg.Contributors = make([]model.ContributorStat, 0, len(ra.contributors))
	for _, ca := range ra.contributors {
		agg.Contributors = append(agg.Contributors, ca.ContributorStat)
	}
	slices.SortFunc(agg.Contributors, func(x, y model.ContributorStat) int {
		return cmp.Compare(y.CommitCount, x.CommitCount)
	})
	if len(agg.Contributors) == 0 {
		agg.Contributors = []model.ContributorStat{}
	}

	agg.ChartCalendar = buildFullCalendar(ra.calendar, ra.earliestDate, time.Now())
	agg.ChartCumulative = buildCumulative(agg.ChartCalendar)
	agg.ChartHourly = buildHourly(ra.hourly)

	return agg
}

func buildFullCalendar(dayMap map[string]int, earliest time.Time, now time.Time) []model.CalendarPoint {
	if earliest.IsZero() {
		return []model.CalendarPoint{}
	}
	loc := now.Location()
	start := time.Date(earliest.Year(), earliest.Month(), earliest.Day(), 0, 0, 0, 0, loc)
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	result := make([]model.CalendarPoint, 0)
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		result = append(result, model.CalendarPoint{
			Date: dateKey, Count: dayMap[dateKey],
		})
	}
	return result
}

func buildCumulative(calendar []model.CalendarPoint) []model.CumulativePoint {
	result := make([]model.CumulativePoint, len(calendar))
	running := 0
	for i, cp := range calendar {
		running += cp.Count
		result[i] = model.CumulativePoint{Date: cp.Date, Total: running}
	}
	return result
}

func buildHourly(hourMap map[int]int) []model.HourlyPoint {
	result := make([]model.HourlyPoint, 24)
	for h := 0; h < 24; h++ {
		result[h] = model.HourlyPoint{Hour: h, Count: hourMap[h]}
	}
	return result
}
