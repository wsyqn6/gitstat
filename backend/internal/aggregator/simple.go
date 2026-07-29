package aggregator

import (
	"cmp"
	"slices"
	"time"

	"gitstat/internal/model"
)

type SimpleAccumulator struct {
	startDate, endDate time.Time
	totalCommits       int
	totalAdditions     int
	totalDeletions     int
	repoSet            map[string]bool
	authorMap          map[string]*model.AuthorRankItem
	lastSeen           map[string]time.Time
}

func NewSimpleAccumulator(startDate, endDate time.Time) *SimpleAccumulator {
	return &SimpleAccumulator{
		startDate: startDate,
		endDate:   endDate,
		repoSet:   make(map[string]bool),
		authorMap: make(map[string]*model.AuthorRankItem),
		lastSeen:  make(map[string]time.Time),
	}
}

func (a *SimpleAccumulator) Add(c *model.Commit, repo *model.Repository) {
	if !a.startDate.IsZero() && c.Date.Before(a.startDate) {
		return
	}
	if !a.endDate.IsZero() && c.Date.After(a.endDate) {
		return
	}

	a.totalCommits++
	a.totalAdditions += c.Additions
	a.totalDeletions += c.Deletions
	a.repoSet[repo.Path] = true

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

func (a *SimpleAccumulator) Build() *AggBucket {
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
		b.AuthorRank = make([]model.AuthorRankItem, len(b.AuthorList))
		copy(b.AuthorRank, b.AuthorList)
	} else {
		b.AuthorList = []model.AuthorRankItem{}
		b.AuthorRank = []model.AuthorRankItem{}
	}

	b.DailyByRepo = []model.RepositoryDailyStats{}
	b.MonthlyByRepo = []model.RepositoryPeriodStats{}
	b.RepoComp = []model.RepoComparison{}
	b.Heatmap = []model.ActivityHeatmapPoint{}
	b.Calendar = []model.MonthlyCalendarItem{}
	b.FileRank = []model.FileRankItem{}

	return b
}
