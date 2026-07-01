package handler

import (
	"math"
	"net/http"
	"time"

	"gitstat/internal/model"
	"gitstat/internal/store"
)

func GetDashboardHandler(w http.ResponseWriter, r *http.Request) {
	repoPaths := r.URL.Query()["repo"]
	startDate, endDate := parseTimeParams(r, "today")
	email := emailForHandler(r)
	bucket := getAggBucket(repoPaths, startDate, endDate, email)

	duration := endDate.Sub(startDate)
	prevStartDate := startDate.Add(-duration)
	prevEndDate := endDate.Add(-duration)
	prevBucket := getAggBucket(repoPaths, prevStartDate, prevEndDate, email)

	current := overviewFromBucket(bucket)
	previous := overviewFromBucket(prevBucket)
	comparison := computeComparison(current, previous)

	resp := model.DashboardData{Comparison: comparison}
	if bucket != nil {
		resp.Summary = model.DashboardSummary{
			TotalCommits:    bucket.TotalCommits,
			TotalAdditions:  bucket.TotalAdditions,
			TotalDeletions:  bucket.TotalDeletions,
			ActiveAuthors:   bucket.ActiveAuthors,
			RepositoryCount: bucket.RepositoryCount,
		}
		resp.DailyRepos = bucket.DailyByRepo
	}
	if resp.DailyRepos == nil {
		resp.DailyRepos = []model.RepositoryDailyStats{}
	}

	writeJSON(w, "Dashboard", model.ApiResponse{Code: 200, Data: resp})
}

func GetOverviewStatsHandler(w http.ResponseWriter, r *http.Request) {
	repoPaths := r.URL.Query()["repo"]
	startDate, endDate := parseTimeParams(r, "today")

	email := emailForHandler(r)

	bucket := getAggBucket(repoPaths, startDate, endDate, email)

	writeJSON(w, "Overview", model.ApiResponse{Code: 200, Data: overviewFromBucket(bucket)})
}

func GetStatsHandler(period string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate, endDate := parseTimeParams(r, "")
		repoPaths := r.URL.Query()["repo"]

		email := emailForHandler(r)
		bucket := getAggBucket(repoPaths, startDate, endDate, email)
		if bucket == nil || bucket.TotalCommits == 0 {
			writeJSON(w, period+" empty", model.ApiResponse{Code: 200, Data: []model.RepositoryDailyStats{}})
			return
		}

		switch period {
		case "daily":
			writeJSON(w, "Daily", model.ApiResponse{Code: 200, Data: bucket.DailyByRepo})
		case "weekly":
			writeJSON(w, "Weekly", model.ApiResponse{Code: 200, Data: bucket.DailyByRepo})
		case "monthly":
			writeJSON(w, "Monthly", model.ApiResponse{Code: 200, Data: map[string]interface{}{
				"repos":           bucket.MonthlyByRepo,
				"monthlyCalendar": bucket.Calendar,
			}})
		case "yearly":
			writeJSON(w, "Yearly", model.ApiResponse{Code: 200, Data: bucket.MonthlyByRepo})
		default:
			writeError(w, ErrCodeInvalidRequest, "invalid period", http.StatusBadRequest)
		}
	}
}

func GetAuthorRankHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "week")
	repoPaths := r.URL.Query()["repo"]

	email := emailForHandler(r)
	bucket := getAggBucket(repoPaths, startDate, endDate, email)

	if bucket == nil {
		writeJSON(w, "AuthorRank", model.ApiResponse{Code: 200, Data: []model.AuthorRankItem{}})
		return
	}

	writeJSON(w, "AuthorRank", model.ApiResponse{Code: 200, Data: bucket.AuthorRank})
}

func GetDailyTrendHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "")
	repoPaths := r.URL.Query()["repo"]
	email := emailForHandler(r)
	bucket := getAggBucket(repoPaths, startDate, endDate, email)
	if bucket == nil || bucket.TotalCommits == 0 {
		writeJSON(w, "daily-trend empty", model.ApiResponse{Code: 200, Data: []model.DailyTrendItem{}})
		return
	}
	items := make([]model.DailyTrendItem, 0, len(bucket.DailyByRepo))
	for _, repo := range bucket.DailyByRepo {
		items = append(items, model.DailyTrendItem{
			RepoName:     repo.RepoName,
			DailyCommits: repo.DailyCommits,
		})
	}
	writeJSON(w, "DailyTrend", model.ApiResponse{Code: 200, Data: items})
}

func GetActivityHeatmapHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "month")
	repoPaths := r.URL.Query()["repo"]
	bucket := getAggBucket(repoPaths, startDate, endDate, "")

	if bucket == nil {
		writeJSON(w, "Heatmap", model.ApiResponse{Code: 200, Data: []model.ActivityHeatmapPoint{}})
		return
	}
	writeJSON(w, "Heatmap", model.ApiResponse{Code: 200, Data: bucket.Heatmap})
}

func GetRepoComparisonHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "week")
	repoPaths := r.URL.Query()["repo"]
	bucket := getAggBucket(repoPaths, startDate, endDate, "")

	if bucket == nil {
		writeJSON(w, "RepoComparison", model.ApiResponse{Code: 200, Data: []model.RepoComparison{}})
		return
	}
	writeJSON(w, "RepoComparison", model.ApiResponse{Code: 200, Data: bucket.RepoComp})
}

func GetComparisonHandler(w http.ResponseWriter, r *http.Request) {
	repoPaths := r.URL.Query()["repo"]

	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	prevStartDateStr := r.URL.Query().Get("prevStartDate")
	prevEndDateStr := r.URL.Query().Get("prevEndDate")

	var startDate, endDate, prevStartDate, prevEndDate time.Time

	if startDateStr != "" && endDateStr != "" {
		var err error
		startDate, err = time.ParseInLocation("2006-01-02", startDateStr, time.Local)
		if err != nil {
			writeError(w, ErrCodeInvalidParams, "invalid startDate format", http.StatusBadRequest)
			return
		}
		endDate, err = time.ParseInLocation("2006-01-02", endDateStr, time.Local)
		if err != nil {
			writeError(w, ErrCodeInvalidParams, "invalid endDate format", http.StatusBadRequest)
			return
		}
		endDate = endDate.Add(24*time.Hour - time.Second)
	}
	if prevStartDateStr != "" && prevEndDateStr != "" {
		var err error
		prevStartDate, err = time.ParseInLocation("2006-01-02", prevStartDateStr, time.Local)
		if err != nil {
			writeError(w, ErrCodeInvalidParams, "invalid prevStartDate format", http.StatusBadRequest)
			return
		}
		prevEndDate, err = time.ParseInLocation("2006-01-02", prevEndDateStr, time.Local)
		if err != nil {
			writeError(w, ErrCodeInvalidParams, "invalid prevEndDate format", http.StatusBadRequest)
			return
		}
		prevEndDate = prevEndDate.Add(24*time.Hour - time.Second)
	}

	fullStart := startDate
	if !prevStartDate.IsZero() && (fullStart.IsZero() || prevStartDate.Before(fullStart)) {
		fullStart = prevStartDate
	}
	fullEnd := endDate
	if !prevEndDate.IsZero() && (fullEnd.IsZero() || prevEndDate.After(fullEnd)) {
		fullEnd = prevEndDate
	}

	ensureDataLoaded(repoPaths, fullStart)

	var userEmail string
	if r.URL.Query().Get("scope") != "all" {
		userEmail = resolveUserEmail(r.URL.Query().Get("email"))
	}

	repos := store.GlobalStore.GetReposWithRange(repoPaths, fullStart, fullEnd)
	if userEmail != "" {
		repos = filterCommitsByEmail(repos, userEmail)
	}

	currentBucket := computeBucket(repos, startDate, endDate)
	prevBucket := computeBucket(repos, prevStartDate, prevEndDate)

	markSelf(currentBucket, userEmail)
	markSelf(prevBucket, userEmail)

	current := overviewFromBucket(currentBucket)
	previous := overviewFromBucket(prevBucket)
	result := computeComparison(current, previous)
	writeJSON(w, "Comparison", model.ApiResponse{Code: 200, Data: result})
}

func GetFileRankingHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "all")
	repoPaths := r.URL.Query()["repo"]
	bucket := getAggBucket(repoPaths, startDate, endDate, "")

	if bucket == nil {
		writeJSON(w, "FileRanking", model.ApiResponse{Code: 200, Data: []model.FileRankItem{}})
		return
	}

	limit := parseIntParam(r, "limit", 5)
	if limit > 100 {
		limit = 100
	}
	rank := bucket.FileRank
	if limit > 0 && len(rank) > limit {
		rank = rank[:limit]
	}
	writeJSON(w, "FileRanking", model.ApiResponse{Code: 200, Data: rank})
}

func emailForHandler(r *http.Request) string {
	if r.URL.Query().Get("scope") == "all" {
		return ""
	}
	return resolveUserEmail(r.URL.Query().Get("email"))
}

func resolveUserEmail(email string) string {
	if email != "" {
		return email
	}
	caches := store.GlobalStore.GetAllCaches()
	for _, c := range caches {
		if c.UserEmail != "" {
			return c.UserEmail
		}
	}
	return ""
}

func computeComparison(current, previous model.OverviewStats) model.OverviewComparison {
	mk := func(cur, prev int) model.MetricChange {
		abs := cur - prev
		var pct float64
		if prev > 0 {
			pct = math.Round(float64(abs)/float64(prev)*1000) / 10
		} else if cur > 0 {
			pct = 100
		}
		return model.MetricChange{
			Current: cur, Previous: prev, Abs: abs, Pct: pct,
		}
	}

	prevMap := make(map[string]model.AuthorRankItem)
	for _, a := range previous.Authors {
		prevMap[a.Email] = a
	}

	var authorComparison []model.AuthorComparisonItem
	for _, ca := range current.Authors {
		pa, exists := prevMap[ca.Email]
		prevCommits, prevAdditions, prevDeletions := 0, 0, 0
		if exists {
			prevCommits = pa.Commits
			prevAdditions = pa.Additions
			prevDeletions = pa.Deletions
		}
		ac := model.AuthorComparisonItem{
			Author: ca.Author, Email: ca.Email, IsMe: ca.IsMe,
			Change: model.AuthorMetricChange{
				Commits:   mk(ca.Commits, prevCommits),
				Additions: mk(ca.Additions, prevAdditions),
				Deletions: mk(ca.Deletions, prevDeletions),
				NetChange: ca.NetChange - (prevAdditions - prevDeletions),
			},
		}
		authorComparison = append(authorComparison, ac)
	}
	if authorComparison == nil {
		authorComparison = []model.AuthorComparisonItem{}
	}

	return model.OverviewComparison{
		TotalCommits:   mk(current.TotalCommits, previous.TotalCommits),
		TotalAdditions: mk(current.TotalAdditions, previous.TotalAdditions),
		TotalDeletions: mk(current.TotalDeletions, previous.TotalDeletions),
		ActiveAuthors:  mk(current.ActiveAuthors, previous.ActiveAuthors),
		Authors:        authorComparison,
	}
}
