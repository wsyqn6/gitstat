package store

import (
	"sort"
	"sync"
	"time"

	"gitstat/internal/model"
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
	FullyLoaded    bool
	RepoSize       int64
	Analyzed       bool
	Branches       []string
	RemoteBranches []string
	TotalLines     int
	Languages      []model.LanguageStat
	Tags           []string
}

type Store struct {
	mu              sync.RWMutex
	ScanPath        string
	Repos           map[string]*RepoCache // path -> cache
	lastMutationAt  time.Time
}

func (s *Store) Touch() {
	s.mu.Lock()
	s.lastMutationAt = time.Now()
	s.mu.Unlock()
}

func (s *Store) LastMutationAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastMutationAt
}

var GlobalStore = &Store{
	Repos: make(map[string]*RepoCache),
}

// MergeCommits 增量合并提交，保持 Commits 按日期升序
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

	cache.Commits = append(cache.Commits, uniqueCommits...)
	s.updateDateRange(cache, uniqueCommits)
	sort.Slice(cache.Commits, func(i, j int) bool {
		return cache.Commits[i].Date.Before(cache.Commits[j].Date)
	})
	s.lastMutationAt = time.Now()
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

func (s *Store) GetReposWithRange(paths []string, startDate, endDate time.Time) []model.Repository {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pathSet := make(map[string]bool, len(paths))
	for _, p := range paths {
		pathSet[p] = true
	}

	var repos []model.Repository
	for _, cache := range s.Repos {
		if len(pathSet) > 0 && !pathSet[cache.Path] {
			continue
		}
		if !cache.Initialized {
			continue
		}

		// quick skip: no date range overlap
		if !startDate.IsZero() && !cache.LatestDate.IsZero() && startDate.After(cache.LatestDate) {
			continue
		}
		if !endDate.IsZero() && !cache.EarliestDate.IsZero() && endDate.Before(cache.EarliestDate) {
			continue
		}

		var filtered []model.Commit
		startIdx := 0
		if !startDate.IsZero() {
			startIdx = sort.Search(len(cache.Commits), func(i int) bool {
				return !cache.Commits[i].Date.Before(startDate)
			})
		}
		for i := startIdx; i < len(cache.Commits); i++ {
			c := cache.Commits[i]
			if !endDate.IsZero() && c.Date.After(endDate) {
				break
			}
			filtered = append(filtered, c)
		}
		if filtered == nil {
			filtered = []model.Commit{}
		}

		repos = append(repos, model.Repository{
			Path:           cache.Path,
			Name:           cache.Name,
			UserEmail:      cache.UserEmail,
			CurrentBranch:  cache.CurrentBranch,
			LastCommitTime: cache.LastCommitTime,
			Commits:        filtered,
		})
	}
	return repos
}

func (s *Store) GetAllInitializedRepos() []model.Repository {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var repos []model.Repository
	for _, cache := range s.Repos {
		if !cache.Initialized {
			continue
		}
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

// UpdateRepo 线程安全更新仓库缓存的指定字段
func (s *Store) UpdateRepo(path string, fn func(*RepoCache)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cache, ok := s.Repos[path]; ok {
		fn(cache)
	}
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
	s.lastMutationAt = time.Now()
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

	if len(commits) == 0 {
		cache.initMu.Unlock()
		return false, nil
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
		sort.Slice(cache.Commits, func(i, j int) bool {
			return cache.Commits[i].Date.Before(cache.Commits[j].Date)
		})
	}
	s.lastMutationAt = time.Now()
}

// SetFullyLoaded 标记仓库已加载完整历史
func (s *Store) SetFullyLoaded(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c, ok := s.Repos[path]; ok {
		c.FullyLoaded = true
	}
}

// IsFullyLoaded 检查仓库是否已加载完整历史
func (s *Store) IsFullyLoaded(path string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if c, ok := s.Repos[path]; ok {
		return c.FullyLoaded
	}
	return false
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
		Name:           cache.Name,
		Path:           cache.Path,
		BranchCount:    cache.BranchCount,
		Branches:       cache.Branches,
		RemoteBranches: cache.RemoteBranches,
		FileCount:      cache.FileCount,
		TotalLines:     cache.TotalLines,
		Languages:      cache.Languages,
		Tags:           cache.Tags,
	}, true
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
