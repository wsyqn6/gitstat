package store

import (
	"sync"
	"time"

	"gitstat/internal/aggregator"
	"gitstat/internal/model"
)

const (
	MaxCommitsPerRepo = 5000 // 单仓库提交数上限
)

type RepoCache struct {
	initMu sync.Mutex // 保护首次初始化扫描

	Path           string
	Name           string
	UserEmail      string
	CurrentBranch  string
	LastCommitTime string
	RemoteUrl      string
	BranchCount    int
	FileCount      int
	EarliestDate   time.Time
	LatestDate     time.Time
	Commits        []model.Commit
	Initialized    bool
	RepoSize       int64
	Analyzed       bool
	Branches       []string
	TotalLines     int
	Languages      []model.LanguageStat
}

type Store struct {
	mu        sync.RWMutex
	ScanPath  string
	Repos     map[string]*RepoCache // path -> cache
}

var GlobalStore = &Store{
	Repos: make(map[string]*RepoCache),
}

func (s *Store) SetRepositories(repos []model.Repository) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, repo := range repos {
		s.Repos[repo.Path] = &RepoCache{
			Path:           repo.Path,
			Name:           repo.Name,
			UserEmail:      repo.UserEmail,
			CurrentBranch:  repo.CurrentBranch,
			LastCommitTime: repo.LastCommitTime,
			Commits:        repo.Commits,
		}
		if len(repo.Commits) > 0 {
			s.updateDateRange(s.Repos[repo.Path], repo.Commits)
		}
	}
}

// MergeCommits 增量合并提交，并检查上限
func (s *Store) MergeCommits(path string, newCommits []model.Commit) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	cache, ok := s.Repos[path]
	if !ok {
		return false
	}

	// 去重：构建已有提交的哈希集合
	existingHashes := make(map[string]bool)
	for _, c := range cache.Commits {
		existingHashes[c.Hash] = true
	}

	// 过滤出未存在的提交
	var uniqueCommits []model.Commit
	for _, nc := range newCommits {
		if !existingHashes[nc.Hash] {
			uniqueCommits = append(uniqueCommits, nc)
			existingHashes[nc.Hash] = true
		}
	}

	if len(uniqueCommits) == 0 {
		return true
	}

	if len(cache.Commits)+len(uniqueCommits) > MaxCommitsPerRepo {
		return false
	}

	cache.Commits = append(cache.Commits, uniqueCommits...)
	s.updateDateRange(cache, uniqueCommits)
	return true
}

func (s *Store) GetRepositories() []model.Repository {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var repos []model.Repository
	for _, cache := range s.Repos {
		repos = append(repos, model.Repository{
			Path:           cache.Path,
			Name:           cache.Name,
			UserEmail:      cache.UserEmail,
			CurrentBranch:  cache.CurrentBranch,
			LastCommitTime: cache.LastCommitTime,
			Commits:        cache.Commits,
		})
	}
	return repos
}

func (s *Store) GetRepoCache(path string) *RepoCache {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Repos[path]
}

// CheckInitRange 原子检查仓库初始化状态和日期范围
func (s *Store) CheckInitRange(path string) (ok bool, initialized bool, earliest, latest time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cache, exists := s.Repos[path]
	if !exists {
		return false, false, time.Time{}, time.Time{}
	}
	return true, cache.Initialized, cache.EarliestDate, cache.LatestDate
}

// GetAllCaches 获取所有缓存副本（用于增量检查）
func (s *Store) GetAllCaches() map[string]*RepoCache {
	s.mu.RLock()
	defer s.mu.RUnlock()
	caches := make(map[string]*RepoCache)
	for path, cache := range s.Repos {
		caches[path] = cache
	}
	return caches
}

func (s *Store) GetAllCommits() []model.Commit {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allCommits []model.Commit
	for _, cache := range s.Repos {
		allCommits = append(allCommits, cache.Commits...)
	}
	return allCommits
}

func (s *Store) updateDateRange(cache *RepoCache, commits []model.Commit) {
	for _, c := range commits {
		if cache.EarliestDate.IsZero() || c.Date.Before(cache.EarliestDate) {
			cache.EarliestDate = c.Date
		}
		if cache.LatestDate.IsZero() || c.Date.After(cache.LatestDate) {
			cache.LatestDate = c.Date
		}
	}
}

// ClearAll 清空所有缓存（切换路径时调用）
func (s *Store) ClearAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ScanPath = ""
	s.Repos = make(map[string]*RepoCache)
}

// RegisterRepos 注册仓库元数据（未初始化，等待懒加载）
func (s *Store) RegisterRepos(repos []model.Repository) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, repo := range repos {
		s.Repos[repo.Path] = &RepoCache{
			Path:           repo.Path,
			Name:           repo.Name,
			UserEmail:      repo.UserEmail,
			CurrentBranch:  repo.CurrentBranch,
			LastCommitTime: repo.LastCommitTime,
			Initialized:    false,
			Commits:        []model.Commit{},
		}
	}
}

// InitRepoCache 初始化仓库缓存（首次扫描后调用）
func (s *Store) InitRepoCache(path string, meta model.Repository, commits []model.Commit) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cache, ok := s.Repos[path]
	if !ok {
		return
	}
	cache.Name = meta.Name
	cache.UserEmail = meta.UserEmail
	cache.CurrentBranch = meta.CurrentBranch
	cache.LastCommitTime = meta.LastCommitTime
	cache.Commits = commits
	cache.Initialized = true
	if len(commits) > 0 {
		s.updateDateRange(cache, commits)
	}
}

// EnsureFirstInit 首次初始化保护（双检锁，仅一个 goroutine 扫库）
func (s *Store) EnsureFirstInit(path string, scanFn func() ([]model.Commit, error)) (bool, error) {
	s.mu.RLock()
	cache, exists := s.Repos[path]
	s.mu.RUnlock()
	if !exists {
		return false, nil
	}

	cache.initMu.Lock()
	_, initialized, _, _ := s.CheckInitRange(path)
	if initialized {
		cache.initMu.Unlock()
		return true, nil
	}

	commits, err := scanFn()
	if err != nil {
		cache.initMu.Unlock()
		return false, err
	}

	s.SetRepoCommits(path, commits)
	cache.initMu.Unlock()
	return true, nil
}

// SetRepoCommits 设置仓库提交数据（保留元数据，仅更新 commits + Initialized）
func (s *Store) SetRepoCommits(path string, commits []model.Commit) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cache, ok := s.Repos[path]
	if !ok {
		return
	}
	cache.Commits = commits
	cache.Initialized = true
	if len(commits) > 0 {
		s.updateDateRange(cache, commits)
	}
}

// UpdateRepoMeta 更新仓库快数据（分支数、文件数）
func (s *Store) UpdateRepoMeta(path string, branchCount, fileCount int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cache, ok := s.Repos[path]
	if !ok {
		return
	}
	cache.BranchCount = branchCount
	cache.FileCount = fileCount
}

// GetAnalyzeCache 获取缓存的深度分析结果
func (s *Store) GetAnalyzeCache(path string) (model.AnalyzeResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cache, ok := s.Repos[path]
	if !ok || !cache.Analyzed {
		return model.AnalyzeResult{}, false
	}
	return model.AnalyzeResult{
		Name:        cache.Name,
		Path:        cache.Path,
		BranchCount: cache.BranchCount,
		Branches:    cache.Branches,
		FileCount:   cache.FileCount,
		TotalLines:  cache.TotalLines,
		Languages:   cache.Languages,
	}, true
}

// GetRepoStats 获取仓库统计（只含计算字段）
func (s *Store) GetRepoStats(path string) (model.RepoStats, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cache, ok := s.Repos[path]
	if !ok {
		return model.RepoStats{}, false
	}

	repos := []model.Repository{{
		Path:    cache.Path,
		Name:    cache.Name,
		Commits: cache.Commits,
	}}
	rank := aggregator.AggregateAuthorRank(repos, "", time.Time{}, time.Time{})
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

	recentCommits := cache.Commits
	if len(recentCommits) > 20 {
		recentCommits = recentCommits[len(recentCommits)-20:]
	}
	for i, j := 0, len(recentCommits)-1; i < j; i, j = i+1, j-1 {
		recentCommits[i], recentCommits[j] = recentCommits[j], recentCommits[i]
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

	stats := model.RepoStats{
		Path:                 cache.Path,
		Name:                 cache.Name,
		CurrentBranch:        cache.CurrentBranch,
		LastCommitTime:       cache.LastCommitTime,
		EarliestDate:         earliestDate,
		EarliestCommitAuthor: earliestAuthor,
		RepoSize:             cache.RepoSize,
		RecentCommits:        recentCommits,
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

	return stats, true
}

// CacheRepoSize 写入仓库大小到缓存（由 handler 按需调用）
func (s *Store) CacheRepoSize(path string, size int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cache, ok := s.Repos[path]; ok {
		cache.RepoSize = size
	}
}

// SetAnalyzeCache 写入深度分析结果到缓存
func (s *Store) SetAnalyzeCache(path string, result model.AnalyzeResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cache, ok := s.Repos[path]
	if !ok {
		return
	}
	cache.BranchCount = result.BranchCount
	cache.Branches = result.Branches
	cache.FileCount = result.FileCount
	cache.TotalLines = result.TotalLines
	cache.Languages = result.Languages
	cache.Analyzed = true
}

// SetScanPath 设置扫描目录
func (s *Store) SetScanPath(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ScanPath = path
}

// GetScanPath 获取当前扫描目录
func (s *Store) GetScanPath() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ScanPath
}
