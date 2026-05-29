package handler

import (
	"log"
	"sync"
	"time"

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

	var wg sync.WaitGroup
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

		wg.Add(1)
		cache := cache
		go func() {
			defer wg.Done()
			ensureRepoLoaded(cache.Path, startDate, now)
		}()
	}
	wg.Wait()
}

func ensureRepoLoaded(repoPath string, startDate, now time.Time) {
	mu := getScanMutex(repoPath)
	mu.Lock()
	defer mu.Unlock()

	cache := store.GlobalStore.GetRepoCache(repoPath)
	if cache == nil {
		return
	}

	if !cache.Initialized {
		commits, err := scanner.ScanCommitsByRange(repoPath, startDate, now)
		if err != nil {
			log.Printf("[LazyLoad] Scan failed for %s: %v", repoPath, err)
			return
		}
		store.GlobalStore.SetRepoCommits(repoPath, commits)
		log.Printf("[LazyLoad] Initialized %s with %d commits", repoPath, len(commits))
		return
	}

	if !startDate.IsZero() && !cache.EarliestDate.IsZero() && cache.EarliestDate.After(startDate) {
		log.Printf("[LazyLoad] Scanning backward: %s from %s to %s", repoPath, startDate.Format("2006-01-02"), cache.EarliestDate.Format("2006-01-02"))
		newCommits, err := scanner.ScanCommitsByRange(repoPath, startDate, cache.EarliestDate)
		if err != nil {
			log.Printf("[LazyLoad] Backward scan failed for %s: %v", repoPath, err)
		} else if len(newCommits) > 0 {
			store.GlobalStore.MergeCommits(repoPath, newCommits)
			log.Printf("[LazyLoad] Merged %d commits for %s", len(newCommits), repoPath)
		}
	}

	if !cache.LatestDate.IsZero() && now.After(cache.LatestDate) {
		log.Printf("[LazyLoad] Forward scanning: %s from %s to %s", repoPath, cache.LatestDate.Format("2006-01-02"), now.Format("2006-01-02"))
		newCommits, err := scanner.ScanCommitsByRange(repoPath, cache.LatestDate, now)
		if err != nil {
			log.Printf("[LazyLoad] Forward scan failed for %s: %v", repoPath, err)
		} else if len(newCommits) > 0 {
			store.GlobalStore.MergeCommits(repoPath, newCommits)
			log.Printf("[LazyLoad] Forward scan found %d new commits for %s", len(newCommits), repoPath)
		}
	}
}
