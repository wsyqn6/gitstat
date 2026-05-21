package store

import (
	"log"
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
	EarliestDate   time.Time
	LatestDate     time.Time
	Commits        []model.Commit
	Initialized    bool // 是否已完成首次扫描
}

type Store struct {
	mu    sync.RWMutex
	Repos map[string]*RepoCache // path -> cache
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
