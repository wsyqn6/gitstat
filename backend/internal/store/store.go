package store

import (
	"log"
	"sort"
	"sync"
	"time"

	"gitstat/internal/model"
)

const (
	MaxCommitsPerRepo = 5000 // 单仓库提交数上限
)

type RepoCache struct {
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
		log.Printf("[MergeCommits] No new commits for %s (all %d duplicates)", path, len(newCommits))
		return true
	}

	// 检查上限
	if len(cache.Commits)+len(uniqueCommits) > MaxCommitsPerRepo {
		log.Printf("[MergeCommits] Reject: would exceed limit (%d + %d > %d)", len(cache.Commits), len(uniqueCommits), MaxCommitsPerRepo)
		return false
	}

	cache.Commits = append(cache.Commits, uniqueCommits...)
	s.updateDateRange(cache, uniqueCommits)
	log.Printf("[MergeCommits] Added %d unique commits to %s (total: %d)", len(uniqueCommits), path, len(cache.Commits))
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

// GetRepoDetail 获取仓库详情（聚合数据）
func (s *Store) GetRepoDetail(path string) (model.RepoDetail, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cache, ok := s.Repos[path]
	if !ok {
		return model.RepoDetail{}, false
	}

	// 聚合贡献者
	authorMap := make(map[string]*model.ContributorStat)
	for _, c := range cache.Commits {
		key := c.Email
		if _, ok := authorMap[key]; !ok {
			authorMap[key] = &model.ContributorStat{
				Author: c.Author,
				Email:  c.Email,
			}
		}
		authorMap[key].CommitCount++
		authorMap[key].Additions += c.Additions
		authorMap[key].Deletions += c.Deletions
		if c.Date.After(parseTime(authorMap[key].LastCommitDate)) {
			authorMap[key].LastCommitDate = c.Date.Format("2006-01-02 15:04:05")
		}
	}

	var contributors []model.ContributorStat
	for _, cs := range authorMap {
		contributors = append(contributors, *cs)
	}
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].CommitCount > contributors[j].CommitCount
	})

	// 最近 20 条提交
	recentCommits := cache.Commits
	if len(recentCommits) > 20 {
		recentCommits = recentCommits[len(recentCommits)-20:]
	}
	// 反转为时间降序
	for i, j := 0, len(recentCommits)-1; i < j; i, j = i+1, j-1 {
		recentCommits[i], recentCommits[j] = recentCommits[j], recentCommits[i]
	}

	// 最早日期
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

	detail := model.RepoDetail{
		Path:                 cache.Path,
		Name:                 cache.Name,
		CurrentBranch:        cache.CurrentBranch,
		LastCommitTime:       cache.LastCommitTime,
		RemoteUrl:            cache.RemoteUrl,
		EarliestDate:         earliestDate,
		EarliestCommitAuthor: earliestAuthor,
		BranchCount:          cache.BranchCount,
		FileCount:            cache.FileCount,
		RecentCommits:        recentCommits,
		Contributors:         contributors,
	}

	// 检查分析缓存
	if cache.Analyzed {
		detail.Analysis = &model.AnalyzeResult{
			Name:        cache.Name,
			Path:        cache.Path,
			BranchCount: cache.BranchCount,
			Branches:    cache.Branches,
			FileCount:   cache.FileCount,
			TotalLines:  cache.TotalLines,
			Languages:   cache.Languages,
		}
	}

	return detail, true
}

func parseTime(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return time.Time{}
	}
	return t
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
