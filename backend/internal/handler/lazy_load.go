package handler

import (
	"log"
	"sync"
	"time"

	"gitstat/internal/scanner"
	"gitstat/internal/store"
)

func ensureDataLoaded(repoPaths []string, startDate time.Time) {
	now := time.Now()
	caches := store.GlobalStore.GetAllCaches()

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
		go func(c *store.RepoCache) {
			defer wg.Done()
			ensureRepoLoaded(c, startDate, now)
		}(cache)
	}
	wg.Wait()
}

func ensureRepoLoaded(cache *store.RepoCache, startDate, now time.Time) {
	cache.ScanMu.Lock()
	defer cache.ScanMu.Unlock()

	if cache.Initialized && !now.After(cache.LatestDate) {
		return
	}

	var since time.Time
	if !cache.Initialized {
		since = startDate
		if !since.IsZero() {
			if minStart := now.AddDate(0, 0, -6); since.After(minStart) {
				since = minStart
			}
		}
	} else {
		since = cache.LatestDate
	}

	log.Printf("[LazyLoad] Scan %s range=%s～%s", cache.Path, since.Format("01-02"), now.Format("01-02"))
	commits, err := scanner.ScanCommitsByRange(cache.Path, since, now)
	if err != nil {
		log.Printf("[LazyLoad] Scan failed for %s: %v", cache.Path, err)
		return
	}
	if len(commits) == 0 {
		return
	}

	if !cache.Initialized {
		store.GlobalStore.SetRepoCommits(cache.Path, commits)
	} else {
		store.GlobalStore.MergeCommits(cache.Path, commits)
	}
}

func PreWarmData() {
	caches := store.GlobalStore.GetAllCaches()
	if len(caches) == 0 {
		return
	}
	log.Printf("[WarmUp] Start loading %d repos...", len(caches))
	var wg sync.WaitGroup
	for _, cache := range caches {
		wg.Add(1)
		go func(c *store.RepoCache) {
			defer wg.Done()
			ensureRepoLoaded(c, time.Time{}, time.Now())
		}(cache)
	}
	wg.Wait()
	log.Printf("[WarmUp] All repos loaded")
}


