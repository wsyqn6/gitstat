package handler

import (
	"net/http"

	"gitstat/internal/aggregator"
)

func GetOverviewStatsHandler(w http.ResponseWriter, r *http.Request) {
	repoPaths := r.URL.Query()["repo"]
	startDate, endDate := parseTimeParams(r, "today")

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)
	userEmail := resolveUserEmail(repos, r.URL.Query().Get("email"))

	stats := aggregator.AggregateOverview(repos, userEmail, startDate, endDate)
	writeSuccess(w, "Overview", stats)
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
			writeSuccess(w, "Daily", stats)
		case "weekly":
			stats := aggregator.AggregateWeeklyStatsWithRange(repos, userEmail, startDate, endDate)
			writeSuccess(w, "Weekly", stats)
		case "monthly":
			stats := aggregator.AggregateMonthlyStatsWithRange(repos, userEmail, startDate, endDate)
			writeSuccess(w, "Monthly", stats)
		case "yearly":
			stats := aggregator.AggregateYearlyStatsWithRange(repos, userEmail, startDate, endDate)
			writeSuccess(w, "Yearly", stats)
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
	writeSuccess(w, "AuthorRank", rank)
}

func GetActivityHeatmapHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "month")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)
	userEmail := resolveUserEmail(repos, r.URL.Query().Get("email"))

	heatmap := aggregator.AggregateActivityHeatmap(repos, userEmail, startDate, endDate)
	writeSuccess(w, "Heatmap", heatmap)
}

func GetRepoComparisonHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "month")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)
	userEmail := resolveUserEmail(repos, r.URL.Query().Get("email"))

	comparison := aggregator.AggregateRepoComparison(repos, userEmail, startDate, endDate)
	writeSuccess(w, "RepoComparison", comparison)
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
	writeSuccess(w, "FileRanking", ranking)
}
