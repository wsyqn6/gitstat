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
				// 无提交时设置 EarliestDate = now，允许后续增量扫描
				c := store.GlobalStore.GetRepoCache(cache.Path)
				if c != nil && c.Initialized && c.EarliestDate.IsZero() {
					c.EarliestDate = now
				}
				mu.Unlock()
				continue
			}
			mu.Unlock()
		}

		// 情况2：已初始化但缓存不足 → 向后扫描（加载更早的历史）
		if !cache.EarliestDate.IsZero() && cache.EarliestDate.After(startDate) {
			mu := getScanMutex(cache.Path)
			mu.Lock()
			cache = store.GlobalStore.GetRepoCache(cache.Path)
			if cache != nil && !cache.EarliestDate.IsZero() && cache.EarliestDate.After(startDate) {
				log.Printf("[LazyLoad] Scanning backward: %s from %s to %s", cache.Path, startDate.Format("2006-01-02"), cache.EarliestDate.Format("2006-01-02"))
				newCommits, err := scanner.ScanIncremental(cache.Path, startDate, cache.EarliestDate)
				if err != nil {
					log.Printf("[LazyLoad] Scan failed for %s: %v", cache.Path, err)
				} else if len(newCommits) > 0 {
					log.Printf("[LazyLoad] Merged %d commits for %s", len(newCommits), cache.Path)
					store.GlobalStore.MergeCommits(cache.Path, newCommits)
				} else {
					log.Printf("[LazyLoad] No new commits found for %s", cache.Path)
				}
			}
			mu.Unlock()
		}

		// 情况3：已初始化且可能有新数据 → 向前扫描（检查新提交）
		if !cache.LatestDate.IsZero() && now.After(cache.LatestDate) {
			mu := getScanMutex(cache.Path)
			mu.Lock()
			cache = store.GlobalStore.GetRepoCache(cache.Path)
			if cache != nil && !cache.LatestDate.IsZero() && now.After(cache.LatestDate) {
				log.Printf("[LazyLoad] Forward scanning: %s from %s to %s", cache.Path, cache.LatestDate.Format("2006-01-02"), now.Format("2006-01-02"))
				newCommits, err := scanner.ScanIncremental(cache.Path, cache.LatestDate, now)
				if err != nil {
					log.Printf("[LazyLoad] Forward scan failed for %s: %v", cache.Path, err)
				} else if len(newCommits) > 0 {
					log.Printf("[LazyLoad] Forward scan found %d new commits for %s", len(newCommits), cache.Path)
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
