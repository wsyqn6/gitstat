package handler

import (
	"net/http"
	"time"

	"gitstat/internal/aggregator"
	"gitstat/internal/model"
)

func GetRepoChartHandler(w http.ResponseWriter, r *http.Request) {
	cache := getRepoCache(w, r)
	if cache == nil {
		return
	}

	ensureRepoLoaded(cache, time.Time{}, time.Now())

	ra := aggregator.NewRepoAccumulator(cache.Path)
	for i := range cache.Commits {
		ra.Add(&cache.Commits[i])
	}
	repoAgg := ra.Build()

	if repoAgg == nil {
		writeJSON(w, "RepoChart", model.RepoChartData{
			DailyAgg: []model.DailyAggPoint{},
			Hourly:   []model.HourlyPoint{},
		})
		return
	}

	writeJSON(w, "RepoChart", model.RepoChartData{
		DailyAgg: repoAgg.ChartDailyAgg,
		Hourly:   repoAgg.ChartHourly,
	})
}
