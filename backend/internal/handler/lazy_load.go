package handler

import (
	"log"
	"sync"
	"time"

	"gitstat/internal/model"
	"gitstat/internal/scanner"
	"gitstat/internal/store"
)

var (
	scanMutexMap   = make(map[string]*sync.Mutex)
	scanMutexGuard sync.Mutex
)

// getScanMutex 获取仓库专属锁，防止并发扫描
func getScanMutex(path string) *sync.Mutex {
	scanMutexGuard.Lock()
	defer scanMutexGuard.Unlock()

	if _, exists := scanMutexMap[path]; !exists {
		scanMutexMap[path] = &sync.Mutex{}
	}
	return scanMutexMap[path]
}

// ensureDataLoaded 统一懒加载：检查缓存 → 不足则扫描 → 合并
func ensureDataLoaded(repoPaths []string, startDate time.Time) {
	caches := store.GlobalStore.GetAllCaches()
	now := time.Now()

	if len(caches) == 0 {
		return
	}

	for _, cache := range caches {
		// 过滤未选中仓库
		if len(repoPaths) > 0 {
			found := false
			for _, p := range repoPaths {
				if p == cache.Path {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// 情况1：未初始化 → 首次扫描元数据 + 当日提交
		if !cache.Initialized {
			mu := getScanMutex(cache.Path)
			mu.Lock()
			// 双重检查
			cache = store.GlobalStore.GetRepoCache(cache.Path)
			if cache != nil && !cache.Initialized {
				initRepoCache(cache.Path, startDate, now)
			}
			mu.Unlock()
			continue
		}

		// 情况2：已初始化但缓存不足 → 增量扫描
		if cache.EarliestDate.After(startDate) {
			mu := getScanMutex(cache.Path)
			mu.Lock()
			// 双重检查
			cache = store.GlobalStore.GetRepoCache(cache.Path)
			if cache != nil && cache.EarliestDate.After(startDate) {
				newCommits, err := scanner.ScanIncremental(cache.Path, startDate, cache.EarliestDate)
				if err == nil && len(newCommits) > 0 {
					log.Printf("[LazyLoad] Merged %d commits for %s", len(newCommits), cache.Path)
					store.GlobalStore.MergeCommits(cache.Path, newCommits)
				}
			}
			mu.Unlock()
		}
	}
}

// initRepoCache 初始化仓库缓存：扫描元数据 + 指定时间范围提交
func initRepoCache(path string, startDate, endDate time.Time) {
	// 1. 扫描元数据
	meta, err := scanner.ScanMetadata(path)
	if err != nil {
		log.Printf("[LazyLoad] Failed to scan metadata for %s: %v", path, err)
		return
	}

	// 2. 扫描指定时间范围的提交
	commits, err := scanner.ScanCommitsByRange(path, startDate, endDate)
	if err != nil {
		log.Printf("[LazyLoad] Failed to scan commits for %s: %v", path, err)
		commits = []model.Commit{}
	}

	// 3. 存入缓存
	store.GlobalStore.InitRepoCache(path, meta, commits)
	log.Printf("[LazyLoad] Initialized %s with %d commits", path, len(commits))
}
