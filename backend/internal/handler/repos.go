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

	var infos []model.RepoInfo
	for _, cache := range caches {
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

		infos = append(infos, model.RepoInfo{
			Path:           cache.Path,
			Name:           cache.Name,
			CurrentBranch:  cache.CurrentBranch,
			BranchCount:    branchCount,
			FileCount:      fileCount,
			LastCommitTime: cache.LastCommitTime,
			RemoteUrl:      remoteUrl,
		})
	}

	if infos == nil {
		infos = []model.RepoInfo{}
	}

	slices.SortFunc(infos, func(a, b model.RepoInfo) int {
		return cmp.Compare(a.Name, b.Name)
	})

	writeJSON(w, "ReposList", infos)
}

func GetRepoDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}

	detail, ok := store.GlobalStore.GetRepoDetail(path)
	if !ok {
		http.Error(w, "repo not found", http.StatusNotFound)
		return
	}

	// 补充未被 store 缓存的字段
	if detail.RemoteUrl == "" {
		detail.RemoteUrl = scanner.GetRemoteUrl(path)
	}
	if detail.RepoSize == 0 {
		detail.RepoSize = scanner.GetRepoSize(path)
	}

	writeJSON(w, "RepoDetail", detail)
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
