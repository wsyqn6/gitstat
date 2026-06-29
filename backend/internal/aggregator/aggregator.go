package aggregator

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"gitstat/internal/model"
)

func AggregateOverview(repos []model.Repository) model.OverviewStats {
	var totalAdditions, totalDeletions int
	var commitCount int
	authorMap := make(map[string]*model.AuthorRankItem)
	repoSet := make(map[string]bool)

	for _, repo := range repos {
		repoSet[repo.Path] = true
		for _, c := range repo.Commits {
			commitCount++
			totalAdditions += c.Additions
			totalDeletions += c.Deletions

			key := c.Email
			if _, exists := authorMap[key]; !exists {
				authorMap[key] = &model.AuthorRankItem{
					Author: c.Author,
					Email:  c.Email,
				}
			}
			stats := authorMap[key]
			stats.Commits++
			stats.Additions += c.Additions
			stats.Deletions += c.Deletions
			stats.NetChange = stats.Additions - stats.Deletions
		}
	}

	var authorList []model.AuthorRankItem
	for _, item := range authorMap {
		authorList = append(authorList, *item)
	}
	slices.SortFunc(authorList, func(a, b model.AuthorRankItem) int {
		return cmp.Compare(b.Commits, a.Commits)
	})

	return model.OverviewStats{
		TotalCommits:    commitCount,
		TotalAdditions:  totalAdditions,
		TotalDeletions:  totalDeletions,
		ActiveAuthors:   len(authorMap),
		RepositoryCount: len(repoSet),
		Authors:         authorList,
	}
}

func AggregateDailyStatsWithRange(repos []model.Repository) []model.RepositoryDailyStats {
	var result []model.RepositoryDailyStats

	for _, repo := range repos {
		authorMap := make(map[string]*model.AuthorDailyStats)
		dailyDataMap := make(map[string]map[string]*model.DayCommitData) // email -> date -> data

		for _, commit := range repo.Commits {
			commitDate := commit.Date.Format("2006-01-02")

			key := commit.Email
			if _, exists := authorMap[key]; !exists {
				authorMap[key] = &model.AuthorDailyStats{
					Author: commit.Author,
					Email:  commit.Email,
				}
				if dailyDataMap[key] == nil {
					dailyDataMap[key] = make(map[string]*model.DayCommitData)
				}
			}

			stats := authorMap[key]
			stats.Commits++
			stats.Additions += commit.Additions
			stats.Deletions += commit.Deletions

			// 收集每日数据
			if _, exists := dailyDataMap[key][commitDate]; !exists {
				dailyDataMap[key][commitDate] = &model.DayCommitData{
					Date: commitDate,
				}
			}
			dayData := dailyDataMap[key][commitDate]
			dayData.Commits++
			dayData.Additions += commit.Additions
			dayData.Deletions += commit.Deletions
		}

		if len(authorMap) > 0 {
			var authors []model.AuthorDailyStats
			daySum := make(map[string]int) // date -> total commits
			for email, stats := range authorMap {
				// 转换每日数据为数组并排序
				var dailyData []model.DayCommitData
				for _, dayData := range dailyDataMap[email] {
					dailyData = append(dailyData, *dayData)
					daySum[dayData.Date] += dayData.Commits
				}
				slices.SortFunc(dailyData, func(a, b model.DayCommitData) int {
					return cmp.Compare(a.Date, b.Date)
				})
				stats.DailyData = dailyData

				authors = append(authors, *stats)
			}

			// 按提交数排序
			slices.SortFunc(authors, func(a, b model.AuthorDailyStats) int {
				return cmp.Compare(b.Commits, a.Commits)
			})

			dailyCommits := make([]model.DayCommitCount, 0, len(daySum))
			for date, count := range daySum {
				dailyCommits = append(dailyCommits, model.DayCommitCount{Date: date, Commits: count})
			}
			slices.SortFunc(dailyCommits, func(a, b model.DayCommitCount) int {
				return cmp.Compare(a.Date, b.Date)
			})

			result = append(result, model.RepositoryDailyStats{
				RepoName:       repo.Name,
				RepoPath:       repo.Path,
				CurrentBranch:  repo.CurrentBranch,
				LastCommitTime: repo.LastCommitTime,
				Authors:        authors,
				DailyCommits:   dailyCommits,
			})
		}
	}

	if result == nil {
		result = []model.RepositoryDailyStats{}
	}
	return result
}

func getWeekKey(t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%d-W%02d", year, week)
}

func getMonthKey(t time.Time) string {
	return t.Format("2006-01")
}

func getYearKey(t time.Time) string {
	return t.Format("2006")
}

func aggregatePeriodStats(repos []model.Repository, periodKeyFn func(time.Time) string) []model.RepositoryPeriodStats {
	var result []model.RepositoryPeriodStats

	for _, repo := range repos {
		authorMap := make(map[string]*model.AuthorPeriodStats)
		periodDataMap := make(map[string]map[string]*model.PeriodCommitData)

		for _, commit := range repo.Commits {
			key := commit.Email
			if _, exists := authorMap[key]; !exists {
				authorMap[key] = &model.AuthorPeriodStats{
					Author: commit.Author,
					Email:  commit.Email,
				}
				if periodDataMap[key] == nil {
					periodDataMap[key] = make(map[string]*model.PeriodCommitData)
				}
			}

			stats := authorMap[key]
			stats.Commits++
			stats.Additions += commit.Additions
			stats.Deletions += commit.Deletions

			periodKey := periodKeyFn(commit.Date)
			if _, exists := periodDataMap[key][periodKey]; !exists {
				periodDataMap[key][periodKey] = &model.PeriodCommitData{
					Period: periodKey,
				}
			}
			pd := periodDataMap[key][periodKey]
			pd.Commits++
			pd.Additions += commit.Additions
			pd.Deletions += commit.Deletions
		}

		if len(authorMap) > 0 {
			var authors []model.AuthorPeriodStats
			for email, stats := range authorMap {
				var periodData []model.PeriodCommitData
				for _, pd := range periodDataMap[email] {
					periodData = append(periodData, *pd)
				}
				slices.SortFunc(periodData, func(a, b model.PeriodCommitData) int {
					return cmp.Compare(a.Period, b.Period)
				})
				stats.PeriodData = periodData
				authors = append(authors, *stats)
			}

			slices.SortFunc(authors, func(a, b model.AuthorPeriodStats) int {
				return cmp.Compare(b.Commits, a.Commits)
			})

			result = append(result, model.RepositoryPeriodStats{
				RepoName:       repo.Name,
				RepoPath:       repo.Path,
				CurrentBranch:  repo.CurrentBranch,
				LastCommitTime: repo.LastCommitTime,
				Authors:        authors,
			})
		}
	}

	if result == nil {
		result = []model.RepositoryPeriodStats{}
	}
	return result
}

func AggregateWeeklyStatsWithRange(repos []model.Repository) []model.RepositoryPeriodStats {
	return aggregatePeriodStats(repos, getWeekKey)
}

func AggregateMonthlyStatsWithRange(repos []model.Repository) []model.RepositoryPeriodStats {
	return aggregatePeriodStats(repos, getMonthKey)
}

func AggregateYearlyStatsWithRange(repos []model.Repository) []model.RepositoryPeriodStats {
	return aggregatePeriodStats(repos, getYearKey)
}

func AggregateMonthlyCalendar(repos []model.Repository) []model.MonthlyCalendarItem {
	monthMap := make(map[string]*model.MonthlyCalendarItem)
	monthKeys := make([]string, 0, 12)

	for _, repo := range repos {
		for _, c := range repo.Commits {
			month := c.Date.Format("2006-01")
			item, ok := monthMap[month]
			if !ok {
				item = &model.MonthlyCalendarItem{Month: month}
				monthMap[month] = item
				monthKeys = append(monthKeys, month)
			}
			item.Commits++
			item.Additions += c.Additions
			item.Deletions += c.Deletions
		}
	}

	slices.Sort(monthKeys)
	result := make([]model.MonthlyCalendarItem, len(monthKeys))
	for i, k := range monthKeys {
		result[i] = *monthMap[k]
	}
	return result
}

// 开发者排行榜聚合
func AggregateAuthorRank(repos []model.Repository) []model.AuthorRankItem {
	authorMap := make(map[string]*model.AuthorRankItem)
	lastTimes := make(map[string]time.Time)

	for _, repo := range repos {
		for _, commit := range repo.Commits {
			key := commit.Email
			if _, exists := authorMap[key]; !exists {
				authorMap[key] = &model.AuthorRankItem{
					Author: commit.Author,
					Email:  commit.Email,
				}
			}

			stats := authorMap[key]
			stats.Commits++
			stats.Additions += commit.Additions
			stats.Deletions += commit.Deletions
			stats.NetChange = stats.Additions - stats.Deletions
			if commit.Date.After(lastTimes[key]) {
				lastTimes[key] = commit.Date
				stats.LastCommitDate = commit.Date.Format("2006-01-02 15:04:05")
			}
		}
	}

	var result []model.AuthorRankItem
	for _, item := range authorMap {
		result = append(result, *item)
	}

	// 按提交数排序
	slices.SortFunc(result, func(a, b model.AuthorRankItem) int {
		return cmp.Compare(b.Commits, a.Commits)
	})

	if result == nil {
		result = []model.AuthorRankItem{}
	}
	return result
}

// 活动热力图聚合
func AggregateActivityHeatmap(repos []model.Repository) []model.ActivityHeatmapPoint {
	heatmap := make(map[string]int) // "dayOfWeek-hour" -> count
	loc := time.Now().Location()

	for _, repo := range repos {
		for _, commit := range repo.Commits {
			key := fmt.Sprintf("%d-%d", commit.Date.In(loc).Weekday(), commit.Date.In(loc).Hour())
			heatmap[key]++
		}
	}

	// 转换为数组
	var result []model.ActivityHeatmapPoint
	for key, count := range heatmap {
		parts := strings.Split(key, "-")
		if len(parts) < 2 {
			continue
		}
		dayOfWeek, _ := strconv.Atoi(parts[0])
		hour, _ := strconv.Atoi(parts[1])
		result = append(result, model.ActivityHeatmapPoint{
			DayOfWeek:   dayOfWeek,
			Hour:        hour,
			CommitCount: count,
		})
	}

	// 排序
	slices.SortFunc(result, func(a, b model.ActivityHeatmapPoint) int {
		if a.DayOfWeek != b.DayOfWeek {
			return cmp.Compare(a.DayOfWeek, b.DayOfWeek)
		}
		return cmp.Compare(a.Hour, b.Hour)
	})

	if result == nil {
		result = []model.ActivityHeatmapPoint{}
	}
	return result
}

// 仓库对比聚合
func AggregateRepoComparison(repos []model.Repository) []model.RepoComparison {
	var result []model.RepoComparison

	for _, repo := range repos {
		commitSet := make(map[string]bool) // 用于计算活跃天数
		authorSet := make(map[string]bool)
		var commits, additions, deletions int

		for _, commit := range repo.Commits {
			commits++
			additions += commit.Additions
			deletions += commit.Deletions
			authorSet[commit.Email] = true
			commitSet[commit.Date.Format("2006-01-02")] = true
		}

		if commits > 0 {
			activeDays := len(commitSet)
			avgCommitsPerDay := float64(commits) / float64(activeDays)
			// 保留1位小数
			avgCommitsPerDay = math.Round(avgCommitsPerDay*10) / 10

			result = append(result, model.RepoComparison{
				RepoName:         repo.Name,
				RepoPath:         repo.Path,
				Commits:          commits,
				Authors:          len(authorSet),
				Additions:        additions,
				Deletions:        deletions,
				LastCommitTime:   repo.LastCommitTime,
				ActiveDays:       activeDays,
				AvgCommitsPerDay: avgCommitsPerDay,
			})
		}
	}

	// 按提交数排序
	slices.SortFunc(result, func(a, b model.RepoComparison) int {
		return cmp.Compare(b.Commits, a.Commits)
	})

	if result == nil {
		result = []model.RepoComparison{}
	}
	return result
}

func AggregateFileRanking(repos []model.Repository, limit int) []model.FileRankItem {
	fileMap := make(map[string]*model.FileRankItem)
	commitSeen := make(map[string]map[string]bool) // filePath -> hash -> seen

	for _, repo := range repos {
		for _, commit := range repo.Commits {
			for _, f := range commit.Files {
				item, exists := fileMap[f.Path]
				if !exists {
					item = &model.FileRankItem{
						FilePath: f.Path,
					}
					fileMap[f.Path] = item
					commitSeen[f.Path] = make(map[string]bool)
				}
				if !commitSeen[f.Path][commit.Hash] {
					commitSeen[f.Path][commit.Hash] = true
					item.Commits++
				}
				item.Additions += f.Additions
				item.Deletions += f.Deletions
				item.NetChange = item.Additions - item.Deletions
			}
		}
	}

	var result []model.FileRankItem
	for _, item := range fileMap {
		result = append(result, *item)
	}

	slices.SortFunc(result, func(a, b model.FileRankItem) int {
		return cmp.Compare(b.Commits, a.Commits)
	})

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	if result == nil {
		result = []model.FileRankItem{}
	}
	return result
}
