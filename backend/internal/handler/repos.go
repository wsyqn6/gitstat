package handler

import (
	"cmp"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"gitstat/internal/aggregator"
	"gitstat/internal/model"
	"gitstat/internal/scanner"
	"gitstat/internal/store"
)

func parseIntQuery(r *http.Request, key string, defaultVal int) int {
	val := r.URL.Query().Get(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

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

	branchCount := cache.BranchCount
	fileCount := cache.FileCount
	branches := cache.Branches
	remoteUrl := cache.RemoteUrl

	if branchCount == 0 || len(branches) == 0 {
		meta, err := scanner.GetRepoMeta(cache.Path)
		if err == nil {
			branchCount = meta.BranchCount
			fileCount = meta.FileCount
			branches = meta.Branches
			store.GlobalStore.UpdateRepo(cache.Path, func(c *store.RepoCache) {
				c.BranchCount = branchCount
				c.FileCount = fileCount
				c.Branches = branches
			})
		} else {
			log.Printf("warning: failed to get repo meta for %s: %v", cache.Path, err)
		}
	}

	if remoteUrl == "" {
		remoteUrl = scanner.GetRemoteUrl(cache.Path)
		store.GlobalStore.UpdateRepo(cache.Path, func(c *store.RepoCache) {
			c.RemoteUrl = remoteUrl
		})
	}

	remoteBranches, remoteBranchCount := scanner.GetRemoteBranches(cache.Path)

	store.GlobalStore.UpdateRepo(cache.Path, func(c *store.RepoCache) {
		c.RemoteBranches = remoteBranches
	})

	writeJSON(w, "RepoInfo", model.RepoInfo{
		Path:              cache.Path,
		Name:              cache.Name,
		CurrentBranch:     cache.CurrentBranch,
		BranchCount:       branchCount,
		RemoteBranchCount: remoteBranchCount,
		FileCount:         fileCount,
		LastCommitTime:    cache.LastCommitTime,
		RemoteUrl:         remoteUrl,
		Branches:          branches,
		RemoteBranches:    remoteBranches,
	})
}

func GetRepoStatsHandler(w http.ResponseWriter, r *http.Request) {
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

	repos := []model.Repository{{
		Path:    cache.Path,
		Name:    cache.Name,
		Commits: cache.Commits,
	}}
	rank := aggregator.AggregateAuthorRank(repos)
	contributors := make([]model.ContributorStat, len(rank))
	for i, item := range rank {
		contributors[i] = model.ContributorStat{
			Author:         item.Author,
			Email:          item.Email,
			CommitCount:    item.Commits,
			Additions:      item.Additions,
			Deletions:      item.Deletions,
			LastCommitDate: item.LastCommitDate,
		}
	}

	var earliestDate string
	var earliestAuthor string
	if !cache.EarliestDate.IsZero() {
		earliestDate = cache.EarliestDate.Format("2006-01-02 15:04:05")
	}
	if len(cache.Commits) > 0 {
		earliest := cache.Commits[0]
		for _, c := range cache.Commits {
			if c.Date.Before(earliest.Date) {
				earliest = c
			}
		}
		earliestAuthor = earliest.Author
	}

	repoSize := cache.RepoSize
	if repoSize == 0 {
		repoSize = scanner.GetRepoSize(path)
		store.GlobalStore.UpdateRepo(path, func(c *store.RepoCache) {
			c.RepoSize = repoSize
		})
	}

	stats := model.RepoStats{
		Path:                 cache.Path,
		Name:                 cache.Name,
		CurrentBranch:        cache.CurrentBranch,
		LastCommitTime:       cache.LastCommitTime,
		EarliestDate:         earliestDate,
		EarliestCommitAuthor: earliestAuthor,
		RepoSize:             repoSize,
		Contributors:         contributors,
	}

	if cache.Analyzed {
		stats.Analysis = &model.AnalyzeResult{
			Name:        cache.Name,
			Path:        cache.Path,
			BranchCount: cache.BranchCount,
			Branches:    cache.Branches,
			FileCount:   cache.FileCount,
			TotalLines:  cache.TotalLines,
			Languages:   cache.Languages,
		}
	}

	writeJSON(w, "RepoStats", stats)
}

func GetRepoCommitsHandler(w http.ResponseWriter, r *http.Request) {
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

	offset := parseIntQuery(r, "offset", 0)
	limit := parseIntQuery(r, "limit", 30)
	if limit < 1 {
		limit = 30
	}

	ensureRepoLoaded(path, time.Time{}, time.Now())

	all := cache.Commits
	total := len(all)

	// cache.Commits is newest-first (git log default)
	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	commits := all[start:end]
	if commits == nil {
		commits = []model.Commit{}
	}

	hasMore := end < total

	writeJSON(w, "RepoCommits", model.CommitPage{Commits: commits, HasMore: hasMore})
}

func GetRepoAnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	var req model.AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, ErrCodeInvalidRequest, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		writeError(w, ErrCodePathRequired, "Path is required", http.StatusBadRequest)
		return
	}

	var err error
	req.Path, err = validatePath(req.Path)
	if err != nil {
		writeError(w, ErrCodeInvalidRequest, err.Error(), http.StatusBadRequest)
		return
	}

	if result, ok := store.GlobalStore.GetAnalyzeCache(req.Path); ok {
		writeJSON(w, "RepoAnalyze", result)
		return
	}

	result, err := scanner.AnalyzeRepoDeep(req.Path)
	if err != nil {
		writeError(w, ErrCodeInternalError, err.Error(), http.StatusInternalServerError)
		return
	}

	store.GlobalStore.UpdateRepo(req.Path, func(c *store.RepoCache) {
		c.BranchCount = result.BranchCount
		c.Branches = result.Branches
		c.RemoteBranches = result.RemoteBranches
		c.FileCount = result.FileCount
		c.TotalLines = result.TotalLines
		c.Languages = result.Languages
		c.Tags = result.Tags
		c.Analyzed = true
	})
	writeJSON(w, "RepoAnalyze", result)
}

func GetRepoTagsHandler(w http.ResponseWriter, r *http.Request) {
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

	offset := parseIntQuery(r, "offset", -1)
	limit := parseIntQuery(r, "limit", 30)

	// no offset → return count only
	if offset < 0 {
		count := 0
		if len(cache.Tags) > 0 {
			count = len(cache.Tags)
		} else {
			count = scanner.GetTagsCount(path)
		}
		writeJSON(w, "RepoTags", map[string]int{"tagCount": count})
		return
	}

	if limit < 1 {
		limit = 30
	}

	// lazy load full tag list into cache
	if len(cache.Tags) == 0 {
		tags := scanner.GetTags(path)
		if tags == nil {
			tags = []string{}
		}
		store.GlobalStore.UpdateRepo(path, func(c *store.RepoCache) {
			c.Tags = tags
		})
	}

	all := cache.Tags
	total := len(all)

	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	tags := all[start:end]
	if tags == nil {
		tags = []string{}
	}

	writeJSON(w, "RepoTags", model.TagPage{
		Tags:    tags,
		Total:   total,
		HasMore: end < total,
	})
}
