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
	ok, initialized, earliest, latest := store.GlobalStore.CheckInitRange(repoPath)
	if !ok {
		return
	}

	if !initialized {
		_, err := store.GlobalStore.EnsureFirstInit(repoPath, func() ([]model.Commit, error) {
			return scanner.ScanCommitsByRange(repoPath, startDate, now)
		})
		if err != nil {
			log.Printf("[LazyLoad] Scan failed for %s: %v", repoPath, err)
		}
		if err == nil && startDate.IsZero() {
			store.GlobalStore.SetFullyLoaded(repoPath)
		}
		return
	}

	// 已初始化 → 零锁并行增量
	needReverse := false
	if startDate.IsZero() {
		needReverse = !store.GlobalStore.IsFullyLoaded(repoPath)
	} else if !earliest.IsZero() && earliest.After(startDate) {
		needReverse = true
	} else if earliest.IsZero() && !store.GlobalStore.IsFullyLoaded(repoPath) {
		needReverse = true
	}

	if needReverse {
		revStart := time.Time{}
		revEnd := earliest
		if revEnd.IsZero() {
			if !startDate.IsZero() {
				revEnd = startDate
			} else {
				revEnd = now
			}
		}
		newCommits, err := scanner.ScanCommitsByRange(repoPath, revStart, revEnd)
		if err != nil {
			log.Printf("[LazyLoad] Backward scan failed for %s: %v", repoPath, err)
		} else if len(newCommits) > 0 {
			store.GlobalStore.MergeCommits(repoPath, newCommits, false)
		}
		if startDate.IsZero() {
			store.GlobalStore.SetFullyLoaded(repoPath)
		}
	}

	if !latest.IsZero() && now.After(latest) {
		newCommits, err := scanner.ScanCommitsByRange(repoPath, latest, now)
		if err != nil {
			log.Printf("[LazyLoad] Forward scan failed for %s: %v", repoPath, err)
		} else if len(newCommits) > 0 {
			store.GlobalStore.MergeCommits(repoPath, newCommits, true)
		}
	}
}
