package handler

import (
	"net/http"
	"sort"
	"time"

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

	ensureRepoLoaded(path, time.Time{}, time.Now())

	now := time.Now()
	yearAgo := now.AddDate(0, 0, -364)

	dateCount := make(map[string]int)
	dayMap := make(map[string]int)
	hourCount := make(map[int]int)

	for _, c := range cache.Commits {
		dateKey := c.Date.Format("2006-01-02")
		if !c.Date.Before(yearAgo) {
			dateCount[dateKey]++
		}
		dayMap[dateKey]++
		hourCount[c.Date.Hour()]++
	}

	calendar := make([]model.CalendarPoint, 0, 365)
	for d := yearAgo; !d.After(now); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		calendar = append(calendar, model.CalendarPoint{
			Date:  dateKey,
			Count: dateCount[dateKey],
		})
	}

	days := make([]string, 0, len(dayMap))
	for d := range dayMap {
		days = append(days, d)
	}
	sort.Strings(days)

	cumulative := make([]model.CumulativePoint, 0, len(days))
	runningTotal := 0
	for _, d := range days {
		runningTotal += dayMap[d]
		cumulative = append(cumulative, model.CumulativePoint{
			Date:  d,
			Total: runningTotal,
		})
	}

	hourly := make([]model.HourlyPoint, 24)
	for h := 0; h < 24; h++ {
		hourly[h] = model.HourlyPoint{
			Hour:  h,
			Count: hourCount[h],
		}
	}

	writeJSON(w, "RepoChart", model.RepoChartData{
		Calendar:   calendar,
		Cumulative: cumulative,
		Hourly:     hourly,
	})
}
