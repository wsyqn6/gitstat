package aggregator

import (
	"cmp"
	"math"
	"slices"
	"time"

	"gitstat/internal/model"
)

type AggBucket struct {
	BuiltAt        time.Time
	TotalCommits   int
	TotalAdditions int
	TotalDeletions int
	ActiveAuthors  int
	RepositoryCount int

	AuthorList   []model.AuthorRankItem
	DailyByRepo  []model.RepositoryDailyStats
	MonthlyByRepo []model.RepositoryPeriodStats
	AuthorRank   []model.AuthorRankItem
	RepoComp     []model.RepoComparison
	Heatmap      []model.ActivityHeatmapPoint
	Calendar     []model.MonthlyCalendarItem
	FileRank     []model.FileRankItem
}

type repoCompAcc struct {
	commits         int
	authorSet       map[string]bool
	additions       int
	deletions       int
	lastCommitTime  string
	fileSet         map[string]bool
}

type Accumulator struct {
	startDate, endDate time.Time

	totalCommits   int
	totalAdditions int
	totalDeletions int
	repoSet        map[string]bool
	authorMap      map[string]*model.AuthorRankItem

	dailyRepoDate    map[string]map[string]*model.DayCommitData
	dailyAuthRepo    map[string]map[string]*model.AuthorDailyStats
	monthlyRepoPeriod map[string]map[string]*model.PeriodCommitData
	monthlyAuthPeriod map[string]map[string]*model.AuthorPeriodStats

	heatmapDayHour map[int]map[int]int
	fileRank       map[string]*model.FileRankItem
	fileSeen       map[string]map[string]bool
	calendar       map[string]*model.MonthlyCalendarItem
	repoCompMap     map[string]*repoCompAcc
	repoNameOf      map[string]string
	repoBranchOf    map[string]string
	repoLastCommit  map[string]string
	lastSeen        map[string]time.Time
	authorFiles     map[string]map[string]map[string]bool
}

func NewAccumulator(startDate, endDate time.Time) *Accumulator {
	return &Accumulator{
		startDate:  startDate,
		endDate:    endDate,
		repoSet:    make(map[string]bool),
		authorMap:  make(map[string]*model.AuthorRankItem),

		dailyRepoDate:     make(map[string]map[string]*model.DayCommitData),
		dailyAuthRepo:     make(map[string]map[string]*model.AuthorDailyStats),
		monthlyRepoPeriod: make(map[string]map[string]*model.PeriodCommitData),
		monthlyAuthPeriod: make(map[string]map[string]*model.AuthorPeriodStats),

		heatmapDayHour: make(map[int]map[int]int),
		fileRank:       make(map[string]*model.FileRankItem),
		fileSeen:       make(map[string]map[string]bool),
		calendar:       make(map[string]*model.MonthlyCalendarItem),
		repoCompMap:    make(map[string]*repoCompAcc),
		repoNameOf:     make(map[string]string),
		repoBranchOf:   make(map[string]string),
		repoLastCommit: make(map[string]string),
		lastSeen:       make(map[string]time.Time),
		authorFiles:    make(map[string]map[string]map[string]bool),
	}
}

func (a *Accumulator) Add(c *model.Commit, repo *model.Repository) {
	if !a.startDate.IsZero() && c.Date.Before(a.startDate) {
		return
	}
	if !a.endDate.IsZero() && c.Date.After(a.endDate) {
		return
	}

	a.addToTotals(c)
	a.addToRepoSet(repo)
	a.addToAuthorMap(c)

	dateKey := c.Date.Format("2006-01-02")
	periodKey := a.periodKey(c.Date)

	a.addToDailyStats(c, repo, dateKey)
	a.addToMonthlyStats(c, repo, periodKey)
	a.addToCalendar(c, periodKey)
	a.addToHeatmap(c)
	a.addToRepoComp(c, repo)
	a.addToFileRank(c, repo)
}

func (a *Accumulator) addToTotals(c *model.Commit) {
	a.totalCommits++
	a.totalAdditions += c.Additions
	a.totalDeletions += c.Deletions
}

func (a *Accumulator) addToRepoSet(repo *model.Repository) {
	a.repoSet[repo.Path] = true
	a.repoNameOf[repo.Path] = repo.Name
	a.repoBranchOf[repo.Path] = repo.CurrentBranch
	if repo.LastCommitTime != "" {
		a.repoLastCommit[repo.Path] = repo.LastCommitTime
	}
}

func (a *Accumulator) addToAuthorMap(c *model.Commit) {
	entry, ok := a.authorMap[c.Email]
	if !ok {
		entry = &model.AuthorRankItem{
			Author: c.Author,
			Email:  c.Email,
		}
		a.authorMap[c.Email] = entry
	}
	entry.Commits++
	entry.Additions += c.Additions
	entry.Deletions += c.Deletions
	entry.NetChange = entry.Additions - entry.Deletions
	if c.Date.After(a.lastSeen[c.Email]) {
		a.lastSeen[c.Email] = c.Date
	}
}

func (a *Accumulator) addToDailyStats(c *model.Commit, repo *model.Repository, dateKey string) {
	repoDaily, ok := a.dailyRepoDate[repo.Path]
	if !ok {
		repoDaily = make(map[string]*model.DayCommitData)
		a.dailyRepoDate[repo.Path] = repoDaily
	}
	dd, ok := repoDaily[dateKey]
	if !ok {
		dd = &model.DayCommitData{Date: dateKey}
		repoDaily[dateKey] = dd
	}
	dd.Commits++
	dd.Additions += c.Additions
	dd.Deletions += c.Deletions

	emailDaily, ok := a.dailyAuthRepo[c.Email]
	if !ok {
		emailDaily = make(map[string]*model.AuthorDailyStats)
		a.dailyAuthRepo[c.Email] = emailDaily
	}
	ads, ok := emailDaily[repo.Path]
	if !ok {
		ads = &model.AuthorDailyStats{
			Author: c.Author,
			Email:  c.Email,
		}
		emailDaily[repo.Path] = ads
	}
	ads.Commits++
	ads.Additions += c.Additions
	ads.Deletions += c.Deletions
	n := len(ads.DailyData)
	if n > 0 && ads.DailyData[n-1].Date == dateKey {
		ads.DailyData[n-1].Commits++
		ads.DailyData[n-1].Additions += c.Additions
		ads.DailyData[n-1].Deletions += c.Deletions
	} else {
		ads.DailyData = append(ads.DailyData, model.DayCommitData{
			Date: dateKey, Commits: 1,
			Additions: c.Additions, Deletions: c.Deletions,
		})
	}
}

func (a *Accumulator) addToMonthlyStats(c *model.Commit, repo *model.Repository, periodKey string) {
	repoMonthly, ok := a.monthlyRepoPeriod[repo.Path]
	if !ok {
		repoMonthly = make(map[string]*model.PeriodCommitData)
		a.monthlyRepoPeriod[repo.Path] = repoMonthly
	}
	mp, ok := repoMonthly[periodKey]
	if !ok {
		mp = &model.PeriodCommitData{Period: periodKey}
		repoMonthly[periodKey] = mp
	}
	mp.Commits++
	mp.Additions += c.Additions
	mp.Deletions += c.Deletions

	emailMonthly, ok := a.monthlyAuthPeriod[c.Email]
	if !ok {
		emailMonthly = make(map[string]*model.AuthorPeriodStats)
		a.monthlyAuthPeriod[c.Email] = emailMonthly
	}
	aps, ok := emailMonthly[repo.Path]
	if !ok {
		aps = &model.AuthorPeriodStats{Author: c.Author, Email: c.Email}
		emailMonthly[repo.Path] = aps
	}
	aps.Commits++
	aps.Additions += c.Additions
	aps.Deletions += c.Deletions
	np := len(aps.PeriodData)
	if np > 0 && aps.PeriodData[np-1].Period == periodKey {
		aps.PeriodData[np-1].Commits++
		aps.PeriodData[np-1].Additions += c.Additions
		aps.PeriodData[np-1].Deletions += c.Deletions
	} else {
		aps.PeriodData = append(aps.PeriodData, model.PeriodCommitData{
			Period: periodKey, Commits: 1,
			Additions: c.Additions, Deletions: c.Deletions,
		})
	}
}

func (a *Accumulator) addToCalendar(c *model.Commit, periodKey string) {
	cal, ok := a.calendar[periodKey]
	if !ok {
		cal = &model.MonthlyCalendarItem{Month: periodKey}
		a.calendar[periodKey] = cal
	}
	cal.Commits++
	cal.Additions += c.Additions
	cal.Deletions += c.Deletions
}

func (a *Accumulator) addToHeatmap(c *model.Commit) {
	loc := time.Now().Location()
	dow := int(c.Date.In(loc).Weekday())
	hour := c.Date.In(loc).Hour()
	if a.heatmapDayHour[dow] == nil {
		a.heatmapDayHour[dow] = make(map[int]int)
	}
	a.heatmapDayHour[dow][hour]++
}

func (a *Accumulator) addToRepoComp(c *model.Commit, repo *model.Repository) {
	rc, ok := a.repoCompMap[repo.Path]
	if !ok {
		rc = &repoCompAcc{
			authorSet:      make(map[string]bool),
			fileSet:        make(map[string]bool),
			lastCommitTime: repo.LastCommitTime,
		}
		a.repoCompMap[repo.Path] = rc
		a.repoNameOf[repo.Path] = repo.Name
	}
	rc.commits++
	rc.authorSet[c.Email] = true
	rc.additions += c.Additions
	rc.deletions += c.Deletions
}

func (a *Accumulator) addToFileRank(c *model.Commit, repo *model.Repository) {
	rc := a.repoCompMap[repo.Path]
	for _, f := range c.Files {
		item, exists := a.fileRank[f.Path]
		if !exists {
			item = &model.FileRankItem{
				FilePath: f.Path,
			}
			a.fileRank[f.Path] = item
			a.fileSeen[f.Path] = make(map[string]bool)
		}
		if !a.fileSeen[f.Path][c.Hash] {
			a.fileSeen[f.Path][c.Hash] = true
			item.Commits++
		}
		item.Additions += f.Additions
		item.Deletions += f.Deletions
		item.NetChange = item.Additions - item.Deletions

		if a.authorFiles[c.Email] == nil {
			a.authorFiles[c.Email] = make(map[string]map[string]bool)
		}
		if a.authorFiles[c.Email][repo.Path] == nil {
			a.authorFiles[c.Email][repo.Path] = make(map[string]bool)
		}
		a.authorFiles[c.Email][repo.Path][f.Path] = true
		rc.fileSet[f.Path] = true
	}
}

func (a *Accumulator) periodKey(t time.Time) string {
	return t.Format("2006-01")
}

func (a *Accumulator) Build() *AggBucket {
	b := &AggBucket{
		BuiltAt:         time.Now(),
		TotalCommits:    a.totalCommits,
		TotalAdditions:  a.totalAdditions,
		TotalDeletions:  a.totalDeletions,
		ActiveAuthors:   len(a.authorMap),
		RepositoryCount: len(a.repoSet),
	}

	if len(a.authorMap) > 0 {
		b.AuthorList = make([]model.AuthorRankItem, 0, len(a.authorMap))
		for _, item := range a.authorMap {
			b.AuthorList = append(b.AuthorList, *item)
		}
		slices.SortFunc(b.AuthorList, func(x, y model.AuthorRankItem) int {
			return cmp.Compare(y.Commits, x.Commits)
		})
	} else {
		b.AuthorList = []model.AuthorRankItem{}
	}

	for repoPath, dateMap := range a.dailyRepoDate {
		dailyCommits := make([]model.DayCommitCount, 0, len(dateMap))
		for _, dd := range dateMap {
			dailyCommits = append(dailyCommits, model.DayCommitCount{
				Date: dd.Date, Commits: dd.Commits,
			})
		}
		slices.SortFunc(dailyCommits, func(x, y model.DayCommitCount) int {
			return cmp.Compare(x.Date, y.Date)
		})

		repoName := a.repoNameOf[repoPath]
		if repoName == "" {
			repoName = repoPath
		}

		dds := model.RepositoryDailyStats{
			RepoName:       repoName,
			RepoPath:       repoPath,
			CurrentBranch:  a.repoBranchOf[repoPath],
			LastCommitTime: a.repoLastCommit[repoPath],
			DailyCommits:   dailyCommits,
			Authors:        make([]model.AuthorDailyStats, 0),
		}
		b.DailyByRepo = append(b.DailyByRepo, dds)
	}
	slices.SortFunc(b.DailyByRepo, func(x, y model.RepositoryDailyStats) int {
		return cmp.Compare(x.RepoName, y.RepoName)
	})

	for _, repoMap := range a.dailyAuthRepo {
		for repoPath, ads := range repoMap {
			ads.FilesChanged = len(a.authorFiles[ads.Email][repoPath])
			for i := range b.DailyByRepo {
				if b.DailyByRepo[i].RepoPath == repoPath {
					b.DailyByRepo[i].Authors = append(b.DailyByRepo[i].Authors, *ads)
					break
				}
			}
		}
	}
	for i := range b.DailyByRepo {
		slices.SortFunc(b.DailyByRepo[i].Authors, func(x, y model.AuthorDailyStats) int {
			return cmp.Compare(y.Commits, x.Commits)
		})
	}

	monthlyByRepoMap := make(map[string]*model.RepositoryPeriodStats)
	for repoPath := range a.monthlyRepoPeriod {
		repoName := a.repoNameOf[repoPath]
		if repoName == "" {
			repoName = repoPath
		}
		monthlyByRepoMap[repoPath] = &model.RepositoryPeriodStats{
			RepoName:       repoName,
			RepoPath:       repoPath,
			CurrentBranch:  a.repoBranchOf[repoPath],
			LastCommitTime: a.repoLastCommit[repoPath],
			Authors:        make([]model.AuthorPeriodStats, 0),
		}
	}

	for _, repoMap := range a.monthlyAuthPeriod {
		for repoPath, aps := range repoMap {
			rps, ok := monthlyByRepoMap[repoPath]
			if !ok {
				continue
			}
			sorted := make([]model.PeriodCommitData, len(aps.PeriodData))
			copy(sorted, aps.PeriodData)
			slices.SortFunc(sorted, func(x, y model.PeriodCommitData) int {
				return cmp.Compare(x.Period, y.Period)
			})
			rps.Authors = append(rps.Authors, model.AuthorPeriodStats{
				Author:     aps.Author,
				Email:      aps.Email,
				Commits:    aps.Commits,
				Additions:  aps.Additions,
				Deletions:  aps.Deletions,
				PeriodData: sorted,
			})
		}
	}

	b.MonthlyByRepo = make([]model.RepositoryPeriodStats, 0, len(monthlyByRepoMap))
	for _, rps := range monthlyByRepoMap {
		if len(rps.Authors) == 0 {
			rps.Authors = []model.AuthorPeriodStats{}
		}
		b.MonthlyByRepo = append(b.MonthlyByRepo, *rps)
	}
	slices.SortFunc(b.MonthlyByRepo, func(x, y model.RepositoryPeriodStats) int {
		return cmp.Compare(x.RepoName, y.RepoName)
	})

	b.AuthorRank = make([]model.AuthorRankItem, len(b.AuthorList))
	copy(b.AuthorRank, b.AuthorList)

	activeDaysMap := make(map[string]map[string]bool)
	for repoPath, dateMap := range a.dailyRepoDate {
		for date := range dateMap {
			if activeDaysMap[repoPath] == nil {
				activeDaysMap[repoPath] = make(map[string]bool)
			}
			activeDaysMap[repoPath][date] = true
		}
	}

	for repoPath, rc := range a.repoCompMap {
		activeDays := len(activeDaysMap[repoPath])
		var avg float64
		if activeDays > 0 {
			avg = float64(rc.commits) / float64(activeDays)
			avg = math.Round(avg*10) / 10
		}
		repoName := a.repoNameOf[repoPath]
		if repoName == "" {
			repoName = repoPath
		}
		b.RepoComp = append(b.RepoComp, model.RepoComparison{
			RepoName:         repoName,
			RepoPath:         repoPath,
			Commits:          rc.commits,
			Authors:          len(rc.authorSet),
			Additions:        rc.additions,
			Deletions:        rc.deletions,
			FilesChanged:     len(rc.fileSet),
			LastCommitTime:   rc.lastCommitTime,
			ActiveDays:       activeDays,
			AvgCommitsPerDay: avg,
		})
	}
	slices.SortFunc(b.RepoComp, func(x, y model.RepoComparison) int {
		return cmp.Compare(y.Commits, x.Commits)
	})

	for dow, hourMap := range a.heatmapDayHour {
		for hour, count := range hourMap {
			b.Heatmap = append(b.Heatmap, model.ActivityHeatmapPoint{
				DayOfWeek: dow, Hour: hour, CommitCount: count,
			})
		}
	}
	slices.SortFunc(b.Heatmap, func(x, y model.ActivityHeatmapPoint) int {
		if x.DayOfWeek != y.DayOfWeek {
			return cmp.Compare(x.DayOfWeek, y.DayOfWeek)
		}
		return cmp.Compare(x.Hour, y.Hour)
	})
	if b.Heatmap == nil {
		b.Heatmap = []model.ActivityHeatmapPoint{}
	}

	b.Calendar = make([]model.MonthlyCalendarItem, 0, len(a.calendar))
	for _, item := range a.calendar {
		b.Calendar = append(b.Calendar, *item)
	}
	slices.SortFunc(b.Calendar, func(x, y model.MonthlyCalendarItem) int {
		return cmp.Compare(x.Month, y.Month)
	})
	if b.Calendar == nil {
		b.Calendar = []model.MonthlyCalendarItem{}
	}

	b.FileRank = make([]model.FileRankItem, 0, len(a.fileRank))
	for _, item := range a.fileRank {
		b.FileRank = append(b.FileRank, *item)
	}
	slices.SortFunc(b.FileRank, func(x, y model.FileRankItem) int {
		return cmp.Compare(y.Commits, x.Commits)
	})
	if len(b.FileRank) == 0 {
		b.FileRank = []model.FileRankItem{}
	}

	if len(b.DailyByRepo) == 0 {
		b.DailyByRepo = []model.RepositoryDailyStats{}
	}
	if len(b.MonthlyByRepo) == 0 {
		b.MonthlyByRepo = []model.RepositoryPeriodStats{}
	}
	if len(b.AuthorRank) == 0 {
		b.AuthorRank = []model.AuthorRankItem{}
	}
	if len(b.RepoComp) == 0 {
		b.RepoComp = []model.RepoComparison{}
	}
	if len(b.AuthorList) == 0 {
		b.AuthorList = []model.AuthorRankItem{}
	}

	return b
}
