package handler

import (
	"math"
	"net/http"
	"time"

	"gitstat/internal/aggregator"
	"gitstat/internal/model"
)

func GetOverviewStatsHandler(w http.ResponseWriter, r *http.Request) {
	repoPaths := r.URL.Query()["repo"]
	startDate, endDate := parseTimeParams(r, "today")

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)

	var userEmail string
	if r.URL.Query().Get("scope") != "all" {
		userEmail = resolveUserEmail(repos, r.URL.Query().Get("email"))
	}

	stats := aggregator.AggregateOverview(repos, userEmail, startDate, endDate)
	writeJSON(w, "Overview", model.ApiResponse{Code: 200, Data: stats})
}

func GetStatsHandler(period string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate, endDate := parseTimeParams(r, "")
		repoPaths := r.URL.Query()["repo"]

		ensureDataLoaded(repoPaths, startDate)
		repos := loadRepos(repoPaths)
		userEmail := resolveUserEmail(repos, r.URL.Query().Get("email"))

		switch period {
		case "daily":
			stats := aggregator.AggregateDailyStatsWithRange(repos, userEmail, startDate, endDate)
			writeJSON(w, "Daily", model.ApiResponse{Code: 200, Data: stats})
		case "weekly":
			stats := aggregator.AggregateWeeklyStatsWithRange(repos, userEmail, startDate, endDate)
			writeJSON(w, "Weekly", model.ApiResponse{Code: 200, Data: stats})
		case "monthly":
			stats := aggregator.AggregateMonthlyStatsWithRange(repos, userEmail, startDate, endDate)
			writeJSON(w, "Monthly", model.ApiResponse{Code: 200, Data: stats})
		case "yearly":
			stats := aggregator.AggregateYearlyStatsWithRange(repos, userEmail, startDate, endDate)
			writeJSON(w, "Yearly", model.ApiResponse{Code: 200, Data: stats})
		default:
			writeError(w, ErrCodeInvalidRequest, "invalid period", http.StatusBadRequest)
		}
	}
}

func GetAuthorRankHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "week")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)
	userEmail := resolveUserEmail(repos, r.URL.Query().Get("email"))

	rank := aggregator.AggregateAuthorRank(repos, userEmail, startDate, endDate)
	writeJSON(w, "AuthorRank", model.ApiResponse{Code: 200, Data: rank})
}

func GetActivityHeatmapHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "month")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)
	userEmail := resolveUserEmail(repos, r.URL.Query().Get("email"))

	heatmap := aggregator.AggregateActivityHeatmap(repos, userEmail, startDate, endDate)
	writeJSON(w, "Heatmap", model.ApiResponse{Code: 200, Data: heatmap})
}

func GetRepoComparisonHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "month")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)
	userEmail := resolveUserEmail(repos, r.URL.Query().Get("email"))

	comparison := aggregator.AggregateRepoComparison(repos, userEmail, startDate, endDate)
	writeJSON(w, "RepoComparison", model.ApiResponse{Code: 200, Data: comparison})
}

func GetComparisonHandler(w http.ResponseWriter, r *http.Request) {
	repoPaths := r.URL.Query()["repo"]

	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")
	prevStartDateStr := r.URL.Query().Get("prevStartDate")
	prevEndDateStr := r.URL.Query().Get("prevEndDate")

	var startDate, endDate, prevStartDate, prevEndDate time.Time

	if startDateStr != "" && endDateStr != "" {
		startDate, _ = time.ParseInLocation("2006-01-02", startDateStr, time.Local)
		endDate, _ = time.ParseInLocation("2006-01-02", endDateStr, time.Local)
		endDate = endDate.Add(24*time.Hour - time.Second)
	}
	if prevStartDateStr != "" && prevEndDateStr != "" {
		prevStartDate, _ = time.ParseInLocation("2006-01-02", prevStartDateStr, time.Local)
		prevEndDate, _ = time.ParseInLocation("2006-01-02", prevEndDateStr, time.Local)
		prevEndDate = prevEndDate.Add(24*time.Hour - time.Second)
	}

	ensureDataLoaded(repoPaths, startDate)
	ensureDataLoaded(repoPaths, prevStartDate)
	repos := loadRepos(repoPaths)

	var userEmail string
	if r.URL.Query().Get("scope") != "all" {
		userEmail = resolveUserEmail(repos, r.URL.Query().Get("email"))
	}

	current := aggregator.AggregateOverview(repos, userEmail, startDate, endDate)
	previous := aggregator.AggregateOverview(repos, userEmail, prevStartDate, prevEndDate)

	result := computeComparison(current, previous)
	writeJSON(w, "Comparison", model.ApiResponse{Code: 200, Data: result})
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
			Current:  cur,
			Previous: prev,
			Abs:      abs,
			Pct:      pct,
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
			Author: ca.Author,
			Email:  ca.Email,
			IsMe:   ca.IsMe,
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

func GetFileRankingHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "month")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)
	userEmail := resolveUserEmail(repos, r.URL.Query().Get("email"))

	limit := parseIntParam(r, "limit", 5)
	if limit > 100 {
		limit = 100
	}

	ranking := aggregator.AggregateFileRanking(repos, userEmail, startDate, endDate, limit)
	writeJSON(w, "FileRanking", model.ApiResponse{Code: 200, Data: ranking})
}
