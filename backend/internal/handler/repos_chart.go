package handler

import (
	"net/http"
	"time"

	"gitstat/internal/aggregator"
	"gitstat/internal/model"
	"gitstat/internal/store"
)

func GetRepoChartHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		writeError(w, ErrCodePathRequired, "path is required", http.StatusBadRequest)
		return
	}

	path, err := validatePath(path)
	if err != nil {
		writeError(w, ErrCodeInvalidRequest, err.Error(), http.StatusBadRequest)
		return
	}

	cache := store.GlobalStore.GetRepoCache(path)
	if cache == nil {
		writeError(w, ErrCodeRepoNotFound, "repo not found", http.StatusNotFound)
		return
	}

	ensureRepoLoaded(cache, time.Time{}, time.Now())

	ra := aggregator.NewRepoAccumulator(path)
	for i := range cache.Commits {
		ra.Add(&cache.Commits[i])
	}
	repoAgg := ra.Build()

	if repoAgg == nil {
		writeJSON(w, "RepoChart", model.RepoChartData{
			Calendar:   []model.CalendarPoint{},
			Cumulative: []model.CumulativePoint{},
			Hourly:     []model.HourlyPoint{},
		})
		return
	}

	writeJSON(w, "RepoChart", model.RepoChartData{
		Calendar:   repoAgg.ChartCalendar,
		Cumulative: repoAgg.ChartCumulative,
		Hourly:     repoAgg.ChartHourly,
	})
}
