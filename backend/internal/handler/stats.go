package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gitstat/internal/aggregator"
	"gitstat/internal/model"
	"gitstat/internal/store"
)

func GetOverviewStatsHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get("email")
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	repoPaths := r.URL.Query()["repo"] // 支持多个 repo 参数

	log.Printf("[Overview] Request: email=%s, start=%s, end=%s, repos=%v", userEmail, startDateStr, endDateStr, repoPaths)

	// 解析时间范围
	var startDate, endDate time.Time

	if startDateStr != "" && endDateStr != "" {
		startDate, _ = time.ParseInLocation("2006-01-02", startDateStr, time.Local)
		endDate, _ = time.ParseInLocation("2006-01-02", endDateStr, time.Local)
		endDate = endDate.Add(24*time.Hour - time.Second)
	} else {
		// 默认当日
		startDate, endDate = ParseTimeRange("today")
	}

	// 如果没有提供邮箱，从仓库配置中获取
	if userEmail == "" {
		repos := store.GlobalStore.GetRepositories()
		for _, repo := range repos {
			if repo.UserEmail != "" {
				userEmail = repo.UserEmail
				break
			}
		}
	}

	// 懒加载：确保缓存覆盖请求范围
	ensureDataLoaded(repoPaths, startDate)

	repos := store.GlobalStore.GetRepositories()

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

	// 收集所有提交并按时间、邮箱过滤
	var filteredCommits []model.Commit
	for _, repo := range repos {
		for _, c := range repo.Commits {
			// 邮箱过滤
			if userEmail != "" && c.Email != userEmail {
				continue
			}
			// 时间过滤
			if !startDate.IsZero() && c.Date.Before(startDate) {
				continue
			}
			if !endDate.IsZero() && c.Date.After(endDate) {
				continue
			}
			filteredCommits = append(filteredCommits, c)
		}
	}

	stats := aggregator.AggregateOverview(filteredCommits, len(repos))

	// 打印响应JSON
	if respJSON, err := json.MarshalIndent(stats, "", "  "); err == nil {
		log.Printf("[Overview] Response JSON:\n%s", string(respJSON))
	} else {
		log.Printf("[Overview] Response error: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func GetDailyStatsHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get("email")
	timeRange := r.URL.Query().Get("range")
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	repoPaths := r.URL.Query()["repo"] // 支持多个 repo 参数

	log.Printf("[Daily] Request: email=%s, range=%s, start=%s, end=%s, repos=%v",
		userEmail, timeRange, startDateStr, endDateStr, repoPaths)

	// 解析时间范围
	var startDate, endDate time.Time

	// 优先使用自定义日期范围
	if startDateStr != "" && endDateStr != "" {
		startDate, _ = time.ParseInLocation("2006-01-02", startDateStr, time.Local)
		endDate, _ = time.ParseInLocation("2006-01-02", endDateStr, time.Local)
		endDate = endDate.Add(24*time.Hour - time.Second)
	} else {
		// 使用预设时间范围
		startDate, endDate = ParseTimeRange(timeRange)
	}

	// 统一懒加载
	ensureDataLoaded(repoPaths, startDate)

	// 重新获取更新后的仓库数据
	var repos []model.Repository
	for _, cache := range store.GlobalStore.Repos {
		repos = append(repos, model.Repository{
			Path:           cache.Path,
			Name:           cache.Name,
			UserEmail:      cache.UserEmail,
			CurrentBranch:  cache.CurrentBranch,
			LastCommitTime: cache.LastCommitTime,
			Commits:        cache.Commits,
		})
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

	dailyStats := aggregator.AggregateDailyStatsWithRange(repos, userEmail, startDate, endDate)

	// 打印响应JSON
	if respJSON, err := json.MarshalIndent(dailyStats, "", "  "); err == nil {
		log.Printf("[Daily] Response JSON:\n%s", string(respJSON))
	} else {
		log.Printf("[Daily] Response: repos=%d (marshal error: %v)", len(dailyStats), err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dailyStats)
}
