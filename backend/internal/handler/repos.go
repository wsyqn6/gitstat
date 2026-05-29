package handler

import (
	"cmp"
	"encoding/json"
	"net/http"
	"slices"

	"gitstat/internal/model"
	"gitstat/internal/scanner"
	"gitstat/internal/store"
)

func GetReposListHandler(w http.ResponseWriter, r *http.Request) {
	caches := store.GlobalStore.GetAllCaches()

	infos := make([]model.RepoInfo, 0, len(caches))
	for _, cache := range caches {
		infos = append(infos, model.RepoInfo{
			Path:           cache.Path,
			Name:           cache.Name,
			CurrentBranch:  cache.CurrentBranch,
			BranchCount:    cache.BranchCount,
			FileCount:      cache.FileCount,
			LastCommitTime: cache.LastCommitTime,
			RemoteUrl:      cache.RemoteUrl,
		})
	}

	slices.SortFunc(infos, func(a, b model.RepoInfo) int {
		return cmp.Compare(a.Name, b.Name)
	})

	writeJSON(w, "ReposList", infos)
}

func GetRepoInfoHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}

	cache := store.GlobalStore.GetRepoCache(path)
	if cache == nil {
		http.Error(w, "repo not found", http.StatusNotFound)
		return
	}

	branchCount := cache.BranchCount
	fileCount := cache.FileCount
	remoteUrl := cache.RemoteUrl

	if branchCount == 0 {
		meta, err := scanner.GetRepoMeta(cache.Path)
		if err == nil {
			branchCount = meta.BranchCount
			fileCount = meta.FileCount
			store.GlobalStore.UpdateRepoMeta(cache.Path, branchCount, fileCount)
		}
	}

	if remoteUrl == "" {
		remoteUrl = scanner.GetRemoteUrl(cache.Path)
		cache.RemoteUrl = remoteUrl
	}

	writeJSON(w, "RepoInfo", model.RepoInfo{
		Path:           cache.Path,
		Name:           cache.Name,
		CurrentBranch:  cache.CurrentBranch,
		BranchCount:    branchCount,
		FileCount:      fileCount,
		LastCommitTime: cache.LastCommitTime,
		RemoteUrl:      remoteUrl,
	})
}

func GetRepoStatsHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}

	stats, ok := store.GlobalStore.GetRepoStats(path)
	if !ok {
		http.Error(w, "repo not found", http.StatusNotFound)
		return
	}

	if stats.RepoSize == 0 {
		stats.RepoSize = scanner.GetRepoSize(path)
		store.GlobalStore.CacheRepoSize(path, stats.RepoSize)
	}

	writeJSON(w, "RepoStats", stats)
}

func GetRepoAnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	var req model.AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		http.Error(w, "Path is required", http.StatusBadRequest)
		return
	}

	if result, ok := store.GlobalStore.GetAnalyzeCache(req.Path); ok {
		writeJSON(w, "RepoAnalyze", result)
		return
	}

	result, err := scanner.AnalyzeRepoDeep(req.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	store.GlobalStore.SetAnalyzeCache(req.Path, result)
	writeJSON(w, "RepoAnalyze", result)
}
