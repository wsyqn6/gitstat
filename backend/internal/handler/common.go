package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gitstat/internal/store"
)

func getRepoCache(w http.ResponseWriter, r *http.Request) *store.RepoCache {
	path := r.URL.Query().Get("path")
	if path == "" {
		writeError(w, ErrCodePathRequired, "path is required", http.StatusBadRequest)
		return nil
	}
	path, err := validatePath(path)
	if err != nil {
		writeError(w, ErrCodeInvalidRequest, err.Error(), http.StatusBadRequest)
		return nil
	}
	cache := store.GlobalStore.GetRepoCache(path)
	if cache == nil {
		writeError(w, ErrCodeRepoNotFound, "repo not found", http.StatusNotFound)
		return nil
	}
	return cache
}

func parseIntParam(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil || n < 0 {
		return defaultVal
	}
	return n
}

func validatePath(p string) (string, error) {
	clean := filepath.Clean(p)
	if strings.Contains(clean, "..") {
		return "", fmt.Errorf("invalid path: traversal not allowed")
	}
	if !filepath.IsAbs(clean) {
		abs, err := filepath.Abs(clean)
		if err != nil {
			return "", fmt.Errorf("invalid path: %w", err)
		}
		return abs, nil
	}
	return clean, nil
}

func parseTimeParams(r *http.Request, defaultRange string) (startDate, endDate time.Time, err error) {
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	if startDateStr != "" && endDateStr != "" {
		startDate, err = time.ParseInLocation("2006-01-02", startDateStr, time.Local)
		if err != nil {
			return startDate, endDate, fmt.Errorf("invalid startDate %q: %w", startDateStr, err)
		}
		endDate, err = time.ParseInLocation("2006-01-02", endDateStr, time.Local)
		if err != nil {
			return startDate, endDate, fmt.Errorf("invalid endDate %q: %w", endDateStr, err)
		}
		endDate = endDate.Add(24*time.Hour - time.Second)
	} else {
		timeRange := r.URL.Query().Get("range")
		if timeRange == "" {
			timeRange = defaultRange
		}
		startDate, endDate = ParseTimeRange(timeRange)
	}

	return startDate, endDate, nil
}

func writeJSON(w http.ResponseWriter, tag string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("[%s] Serialize error: %v", tag, err)
	}
}
