package handler

import (
	"net/http"

	"gitstat/internal/aggregator"
	"gitstat/internal/model"
)

func GetOverviewStatsHandler(w http.ResponseWriter, r *http.Request) {
	userEmail := r.URL.Query().Get("email")
	repoPaths := r.URL.Query()["repo"]

	startDate, endDate := parseTimeParams(r, "today")

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)

	if userEmail == "" {
		userEmail = resolveUserEmail(repos, "")
	}

	var filteredCommits []model.Commit
	for _, repo := range repos {
		for _, c := range repo.Commits {
			if userEmail != "" && c.Email != userEmail {
				continue
			}
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
	writeJSON(w, "Overview", stats)
}

func GetDailyStatsHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)

	stats := aggregator.AggregateDailyStatsWithRange(repos, r.URL.Query().Get("email"), startDate, endDate)
	writeJSON(w, "Daily", stats)
}

func GetWeeklyStatsHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)

	stats := aggregator.AggregateWeeklyStatsWithRange(repos, r.URL.Query().Get("email"), startDate, endDate)
	writeJSON(w, "Weekly", stats)
}

func GetMonthlyStatsHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)

	stats := aggregator.AggregateMonthlyStatsWithRange(repos, r.URL.Query().Get("email"), startDate, endDate)
	writeJSON(w, "Monthly", stats)
}

func GetYearlyStatsHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)

	stats := aggregator.AggregateYearlyStatsWithRange(repos, r.URL.Query().Get("email"), startDate, endDate)
	writeJSON(w, "Yearly", stats)
}

func GetAuthorRankHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "week")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)
	userEmail := resolveUserEmail(repos, r.URL.Query().Get("email"))

	rank := aggregator.AggregateAuthorRank(repos, userEmail, startDate, endDate)
	writeJSON(w, "AuthorRank", rank)
}

func GetActivityHeatmapHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "month")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)
	userEmail := resolveUserEmail(repos, r.URL.Query().Get("email"))

	heatmap := aggregator.AggregateActivityHeatmap(repos, userEmail, startDate, endDate)
	writeJSON(w, "Heatmap", heatmap)
}

func GetRepoComparisonHandler(w http.ResponseWriter, r *http.Request) {
	startDate, endDate := parseTimeParams(r, "month")
	repoPaths := r.URL.Query()["repo"]

	ensureDataLoaded(repoPaths, startDate)
	repos := loadRepos(repoPaths)
	userEmail := resolveUserEmail(repos, r.URL.Query().Get("email"))

	comparison := aggregator.AggregateRepoComparison(repos, userEmail, startDate, endDate)
	writeJSON(w, "RepoComparison", comparison)
}
