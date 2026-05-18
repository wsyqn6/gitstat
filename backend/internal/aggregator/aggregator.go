package aggregator

import (
	"cmp"
	"slices"
	"time"

	"gitstat/internal/model"
)

func AggregateOverview(commits []model.Commit, repoCount int) model.OverviewStats {
	var totalAdditions, totalDeletions int
	authors := make(map[string]bool)

	for _, c := range commits {
		totalAdditions += c.Additions
		totalDeletions += c.Deletions
		authors[c.Author] = true
	}

	return model.OverviewStats{
		TotalCommits:    len(commits),
		TotalAdditions:  totalAdditions,
		TotalDeletions:  totalDeletions,
		ActiveAuthors:   len(authors),
		RepositoryCount: repoCount,
	}
}

func AggregateByAuthor(commits []model.Commit) map[string]map[string]interface{} {
	authorStats := make(map[string]map[string]interface{})

	for _, c := range commits {
		if _, exists := authorStats[c.Author]; !exists {
			authorStats[c.Author] = map[string]interface{}{
				"commits":   0,
				"additions": 0,
				"deletions": 0,
			}
		}

		stats := authorStats[c.Author]
		stats["commits"] = stats["commits"].(int) + 1
		stats["additions"] = stats["additions"].(int) + c.Additions
		stats["deletions"] = stats["deletions"].(int) + c.Deletions
	}

	return authorStats
}

func AggregateByTime(commits []model.Commit, granularity string) []map[string]interface{} {
	timeSeries := make(map[string]map[string]interface{})

	for _, c := range commits {
		dateKey := c.Date.Format("2006-01-02")
		if _, exists := timeSeries[dateKey]; !exists {
			timeSeries[dateKey] = map[string]interface{}{
				"date":      dateKey,
				"commits":   0,
				"additions": 0,
				"deletions": 0,
			}
		}

		stats := timeSeries[dateKey]
		stats["commits"] = stats["commits"].(int) + 1
		stats["additions"] = stats["additions"].(int) + c.Additions
		stats["deletions"] = stats["deletions"].(int) + c.Deletions
	}

	var result []map[string]interface{}
	for _, v := range timeSeries {
		result = append(result, v)
	}

	slices.SortFunc(result, func(a, b map[string]interface{}) int {
		return cmp.Compare(a["date"].(string), b["date"].(string))
	})

	return result
}

func AggregateDailyStats(repos []model.Repository, userEmail string) []model.RepositoryDailyStats {
	return AggregateDailyStatsWithRange(repos, userEmail, time.Time{})
}

func AggregateDailyStatsWithRange(repos []model.Repository, userEmail string, startDate time.Time) []model.RepositoryDailyStats {
	var result []model.RepositoryDailyStats

	// 如果没有提供邮箱，从仓库配置中获取（使用第一个有配置的）
	if userEmail == "" {
		for _, repo := range repos {
			if repo.UserEmail != "" {
				userEmail = repo.UserEmail
				break
			}
		}
	}

	for _, repo := range repos {
		authorMap := make(map[string]*model.AuthorDailyStats)
		dailyDataMap := make(map[string]map[string]*model.DayCommitData) // email -> date -> data

		for _, commit := range repo.Commits {
			commitDate := commit.Date.Format("2006-01-02")

			// 如果指定了起始日期，过滤掉之前的数据
			if !startDate.IsZero() && commit.Date.Before(startDate) {
				continue
			}

			key := commit.Email
			if _, exists := authorMap[key]; !exists {
				authorMap[key] = &model.AuthorDailyStats{
					Author: commit.Author,
					Email:  commit.Email,
					IsMe:   commit.Email == userEmail,
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
			for email, stats := range authorMap {
				// 转换每日数据为数组并排序
				var dailyData []model.DayCommitData
				for _, dayData := range dailyDataMap[email] {
					dailyData = append(dailyData, *dayData)
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

			result = append(result, model.RepositoryDailyStats{
				RepoName:       repo.Name,
				RepoPath:       repo.Path,
				CurrentBranch:  repo.CurrentBranch,
				LastCommitTime: repo.LastCommitTime,
				Authors:        authors,
			})
		}
	}

	return result
}
