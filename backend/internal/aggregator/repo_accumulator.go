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
	ChartDailyAgg        []model.DailyAggPoint
	ChartHourly          []model.HourlyPoint
}

type dateChangeAcc struct {
	additions int
	deletions int
}

type RepoAccumulator struct {
	repoPath       string
	earliestDate   time.Time
	earliestAuthor string
	contributors   map[string]*contributorAcc
	calendar       map[string]int
	hourly         map[int]int
	changes        map[string]*dateChangeAcc
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
		changes:      make(map[string]*dateChangeAcc),
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

	dateKey := c.Date.In(time.Local).Format("2006-01-02")
	ra.calendar[dateKey]++

	dc, ok := ra.changes[dateKey]
	if !ok {
		dc = &dateChangeAcc{}
		ra.changes[dateKey] = dc
	}
	dc.additions += c.Additions
	dc.deletions += c.Deletions

	hour := c.Date.In(time.Local).Hour()
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

	agg.ChartDailyAgg = buildDailyAgg(ra.calendar, ra.changes, ra.earliestDate, time.Now())
	agg.ChartHourly = buildHourly(ra.hourly)

	return agg
}

func buildDailyAgg(commits map[string]int, changes map[string]*dateChangeAcc, earliest time.Time, now time.Time) []model.DailyAggPoint {
	if earliest.IsZero() {
		return []model.DailyAggPoint{}
	}
	loc := now.Location()
	start := time.Date(earliest.Year(), earliest.Month(), earliest.Day(), 0, 0, 0, 0, loc)
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	result := make([]model.DailyAggPoint, 0)
	runningTotal := 0
	runningNet := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		c := commits[dateKey]
		runningTotal += c
		if dc, ok := changes[dateKey]; ok {
			runningNet += dc.additions - dc.deletions
		}
		result = append(result, model.DailyAggPoint{
			Date: dateKey, Commits: c, Total: runningTotal, NetLines: runningNet,
		})
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
