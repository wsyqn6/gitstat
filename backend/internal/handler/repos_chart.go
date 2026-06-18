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
	loc := now.Location()

	dayMap := make(map[string]int)
	hourCount := make(map[int]int)

	var earliestDate time.Time
	for _, c := range cache.Commits {
		ct := c.Date.In(loc)
		dateKey := ct.Format("2006-01-02")
		dayMap[dateKey]++
		hourCount[ct.Hour()]++
		if earliestDate.IsZero() || ct.Before(earliestDate) {
			earliestDate = ct
		}
	}

	calendar := make([]model.CalendarPoint, 0)
	if !earliestDate.IsZero() {
		start := time.Date(earliestDate.Year(), earliestDate.Month(), earliestDate.Day(), 0, 0, 0, 0, loc)
		end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			dateKey := d.Format("2006-01-02")
			calendar = append(calendar, model.CalendarPoint{
				Date:  dateKey,
				Count: dayMap[dateKey],
			})
		}
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
