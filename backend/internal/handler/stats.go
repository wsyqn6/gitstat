package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"gitstat/internal/aggregator"
	"gitstat/internal/model"
	"gitstat/internal/store"
)

func GetOverviewStatsHandler(w http.ResponseWriter, r *http.Request) {
	commits := store.GlobalStore.GetAllCommits()
	repos := store.GlobalStore.GetRepositories()

	stats := aggregator.AggregateOverview(commits, len(repos))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func GetDailyStatsHandler(w http.ResponseWriter, r *http.Request) {
	repos := store.GlobalStore.GetRepositories()
	userEmail := r.URL.Query().Get("email")
	timeRange := r.URL.Query().Get("range")
	repoPaths := r.URL.Query()["repo"] // 支持多个 repo 参数

	// 解析时间范围
	var startDate time.Time
	now := time.Now()

	switch timeRange {
	case "7": // 本周（从周一开始）
		weekday := int(now.Weekday())
		if weekday == 0 { // 周日
			weekday = 7
		}
		daysFromMonday := weekday - 1
		startDate = now.AddDate(0, 0, -daysFromMonday)
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	case "30":
		startDate = now.AddDate(0, 0, -30)
	case "90":
		startDate = now.AddDate(0, 0, -90)
	default: // all 或其他
		startDate = time.Time{} // 零值表示不限制
	}

	// 如果指定了仓库，过滤仓库列表
	if len(repoPaths) > 0 {
		repoMap := make(map[string]bool)
		for _, path := range repoPaths {
			repoMap[path] = true
		}
		filteredRepos := make([]model.Repository, 0)
		for _, repo := range repos {
			if repoMap[repo.Path] {
				filteredRepos = append(filteredRepos, repo)
			}
		}
		repos = filteredRepos
	}

	dailyStats := aggregator.AggregateDailyStatsWithRange(repos, userEmail, startDate)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dailyStats)
}
