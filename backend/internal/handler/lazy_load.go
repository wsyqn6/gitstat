package handler

import (
	"log"
	"sync"
	"time"

	"gitstat/internal/model"
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

	// Fast path: hash cursor matches HEAD → up to date
	if cache.Initialized && cache.LastScannedHash != "" {
		headHash, err := scanner.GetHEADHash(cache.Path)
		if err == nil && headHash == cache.LastScannedHash {
			return
		}
	}

	var commits []model.Commit
	var err error

	// Try hash-based incremental if we have a valid cursor
	if cache.Initialized && cache.LastScannedHash != "" {
		if scanner.IsAncestor(cache.Path, cache.LastScannedHash) {
			log.Printf("[LazyLoad] Hash-scan %s since %s", cache.Path, cache.LastScannedHash[:8])
			commits, err = scanner.ScanCommitsSince(cache.Path, cache.LastScannedHash)
		} else {
			log.Printf("[LazyLoad] Hash %s not in history, fallback date-scan %s", cache.Path, cache.LastScannedHash[:8])
		}
	}

	// Fallback to date-based scan
	if commits == nil && err == nil {
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
		log.Printf("[LazyLoad] Date-scan %s range=%s～%s", cache.Path, since.Format("01-02"), now.Format("01-02"))
		commits, err = scanner.ScanCommitsByRange(cache.Path, since, now)
	}

	if err != nil {
		log.Printf("[LazyLoad] Scan failed for %s: %v", cache.Path, err)
		return
	}
	if len(commits) == 0 {
		// Still update cursor if we did a hash-scan with no new commits
		if cache.Initialized && cache.LastScannedHash != "" {
			if head, e := scanner.GetHEADHash(cache.Path); e == nil {
				store.GlobalStore.UpdateRepo(cache.Path, func(c *store.RepoCache) {
					c.LastScannedHash = head
				})
			}
		}
		return
	}

	if !cache.Initialized {
		store.GlobalStore.SetRepoCommits(cache.Path, commits)
	} else {
		store.GlobalStore.MergeCommits(cache.Path, commits)
	}

	// Update cursor to current HEAD
	if head, err := scanner.GetHEADHash(cache.Path); err == nil {
		store.GlobalStore.UpdateRepo(cache.Path, func(c *store.RepoCache) {
			c.LastScannedHash = head
		})
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


