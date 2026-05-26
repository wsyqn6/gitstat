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

func getScanMutex(path string) *sync.Mutex {
	scanMutexGuard.Lock()
	defer scanMutexGuard.Unlock()

	if _, exists := scanMutexMap[path]; !exists {
		scanMutexMap[path] = &sync.Mutex{}
	}
	return scanMutexMap[path]
}

func ensureDataLoaded(repoPaths []string, startDate time.Time) {
	caches := store.GlobalStore.GetAllCaches()
	now := time.Now()

	if len(caches) == 0 {
		return
	}

	for _, cache := range caches {
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

		if !cache.Initialized {
			mu := getScanMutex(cache.Path)
			mu.Lock()
			cache = store.GlobalStore.GetRepoCache(cache.Path)
			if cache != nil && !cache.Initialized {
				initRepoCache(cache.Path, startDate, now)
				mu.Unlock()
				continue
			}
			mu.Unlock()
		}

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

func initRepoCache(path string, startDate, endDate time.Time) {
	commits, err := scanner.ScanCommitsByRange(path, startDate, endDate)
	if err != nil {
		log.Printf("[LazyLoad] Failed to scan commits for %s: %v", path, err)
		commits = []model.Commit{}
	}

	store.GlobalStore.SetRepoCommits(path, commits)
	log.Printf("[LazyLoad] Initialized %s with %d commits", path, len(commits))
}
