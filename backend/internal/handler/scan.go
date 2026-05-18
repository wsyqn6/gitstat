package handler

import (
	"encoding/json"
	"net/http"

	"gitstat/internal/model"
	"gitstat/internal/scanner"
	"gitstat/internal/store"
)

func ScanHandler(w http.ResponseWriter, r *http.Request) {
	var req model.ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	repos, err := scanner.ScanDirectory(req.Path, req.TimeRange)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	store.GlobalStore.SetRepositories(repos)

	totalCommits := 0
	for _, repo := range repos {
		totalCommits += len(repo.Commits)
	}

	resp := model.ScanResponse{
		Repositories: repos,
		TotalCommits: totalCommits,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	repos := store.GlobalStore.GetRepositories()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}
