package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gitstat/internal/model"
	"gitstat/internal/store"
)

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

func parseTimeParams(r *http.Request, defaultRange string) (startDate, endDate time.Time) {
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	if startDateStr != "" && endDateStr != "" {
		var parseErr error
		startDate, parseErr = time.ParseInLocation("2006-01-02", startDateStr, time.Local)
		if parseErr != nil {
			log.Printf("warning: invalid startDate %q: %v", startDateStr, parseErr)
		}
		endDate, parseErr = time.ParseInLocation("2006-01-02", endDateStr, time.Local)
		if parseErr != nil {
			log.Printf("warning: invalid endDate %q: %v", endDateStr, parseErr)
		}
		endDate = endDate.Add(24*time.Hour - time.Second)
	} else {
		timeRange := r.URL.Query().Get("range")
		if timeRange == "" {
			timeRange = defaultRange
		}
		startDate, endDate = ParseTimeRange(timeRange)
	}

	return startDate, endDate
}

func loadRepos(repoPaths []string, startDate, endDate time.Time) []model.Repository {
	if len(repoPaths) > 50 {
		repoPaths = repoPaths[:50]
	}
	return store.GlobalStore.GetReposWithRange(repoPaths, startDate, endDate)
}

func filterCommitsByDate(repos []model.Repository, startDate, endDate time.Time) []model.Repository {
	if startDate.IsZero() && endDate.IsZero() {
		return repos
	}
	result := make([]model.Repository, len(repos))
	for i, repo := range repos {
		var filtered []model.Commit
		for _, c := range repo.Commits {
			if !startDate.IsZero() && c.Date.Before(startDate) {
				continue
			}
			if !endDate.IsZero() && c.Date.After(endDate) {
				continue
			}
			filtered = append(filtered, c)
		}
		if filtered == nil {
			filtered = []model.Commit{}
		}
		result[i] = repo
		result[i].Commits = filtered
	}
	return result
}

func filterCommitsByEmail(repos []model.Repository, email string) []model.Repository {
	if email == "" {
		return repos
	}
	result := make([]model.Repository, 0, len(repos))
	for _, repo := range repos {
		var filtered []model.Commit
		for _, c := range repo.Commits {
			if c.Email == email {
				filtered = append(filtered, c)
			}
		}
		if len(filtered) > 0 {
			result = append(result, repo)
			result[len(result)-1].Commits = filtered
		}
	}
	return result
}

func resolveUserEmail(repos []model.Repository, email string) string {
	if email != "" {
		return email
	}
	for _, repo := range repos {
		if repo.UserEmail != "" {
			return repo.UserEmail
		}
	}
	return ""
}

func writeJSON(w http.ResponseWriter, tag string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		log.Printf("[%s] Serialize error: %v", tag, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "serialization failed", Code: ErrCodeInternalError})
		return
	}

	w.Write(buf.Bytes())
}
