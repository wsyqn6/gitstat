package handler

import (
	"encoding/json"
	"net/http"

	"gitstat/internal/aggregator"
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

	dailyStats := aggregator.AggregateDailyStats(repos, userEmail)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dailyStats)
}
