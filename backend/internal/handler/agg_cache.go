package handler

import (
	"log"
	"sort"
	"sync"
	"time"

	"gitstat/internal/aggregator"
	"gitstat/internal/model"
	"gitstat/internal/store"
)

type aggCacheEntry struct {
	bucket    *aggregator.AggBucket
	expiresAt time.Time
}

var (
	aggCacheMu sync.Mutex
	aggCache   = map[string]*aggCacheEntry{}
)

// inflight dedup: 相同 key 的并发请求只计算一次
var (
	inflightMu sync.Mutex
	inflight   = map[string]chan struct{}{}
)

func getAggBucket(repoPaths []string, startDate, endDate time.Time, email string, simple bool) *aggregator.AggBucket {
	key := cacheKey(repoPaths, startDate, endDate, email, simple)

	aggCacheMu.Lock()
	entry, ok := aggCache[key]
	if ok && time.Now().Before(entry.expiresAt) {
		aggCacheMu.Unlock()
		log.Printf("[AggCache] TTL hit key=%s", key[:min(len(key), 60)])
		return entry.bucket
	}
	aggCacheMu.Unlock()

	inflightMu.Lock()
	ch, ok := inflight[key]
	if ok {
		inflightMu.Unlock()
		log.Printf("[AggCache] Inflight wait key=%s", key[:min(len(key), 60)])
		<-ch
		aggCacheMu.Lock()
		entry, ok = aggCache[key]
		aggCacheMu.Unlock()
		if ok {
			log.Printf("[AggCache] Inflight served key=%s", key[:min(len(key), 60)])
			return entry.bucket
		}
		return nil
	}
	ch = make(chan struct{})
	inflight[key] = ch
	inflightMu.Unlock()

	log.Printf("[AggCache] Compute key=%s", key[:min(len(key), 60)])
	bucket := computeAggBucket(repoPaths, startDate, endDate, email, simple)

	aggCacheMu.Lock()
	aggCache[key] = &aggCacheEntry{bucket: bucket, expiresAt: time.Now().Add(3 * time.Second)}
	delete(inflight, key)
	aggCacheMu.Unlock()
	close(ch)

	return bucket
}

func computeAggBucket(repoPaths []string, startDate, endDate time.Time, email string, simple bool) *aggregator.AggBucket {
	ensureDataLoaded(repoPaths, startDate)

	repos := store.GlobalStore.GetReposWithRange(repoPaths, startDate, endDate)
	repos = filterCommitsByEmail(repos, email)

	if simple {
		acc := aggregator.NewSimpleAccumulator(startDate, endDate)
		for ri := range repos {
			repo := &repos[ri]
			for ci := range repo.Commits {
				acc.Add(&repo.Commits[ci], repo)
			}
		}
		bucket := acc.Build()
		markSelf(bucket, email)
		return bucket
	}

	acc := aggregator.NewAccumulator(startDate, endDate)
	for ri := range repos {
		repo := &repos[ri]
		for ci := range repo.Commits {
			acc.Add(&repo.Commits[ci], repo)
		}
	}
	bucket := acc.Build()
	markSelf(bucket, email)
	return bucket
}

func cacheKey(repoPaths []string, startDate, endDate time.Time, email string, simple bool) string {
	k := email + "|"
	if simple {
		k = "S|" + k
	} else {
		k = "F|" + k
	}
	if !startDate.IsZero() {
		k += startDate.Format(time.RFC3339)
	}
	k += "|"
	if !endDate.IsZero() {
		k += endDate.Format(time.RFC3339)
	}
	k += "|"
	sorted := make([]string, len(repoPaths))
	copy(sorted, repoPaths)
	sort.Strings(sorted)
	for _, p := range sorted {
		k += p + ","
	}
	return k
}

func filterCommitsByEmail(repos []model.Repository, email string) []model.Repository {
	if email == "" {
		return repos
	}
	for i := range repos {
		filtered := make([]model.Commit, 0, len(repos[i].Commits))
		for _, c := range repos[i].Commits {
			if c.Email == email {
				filtered = append(filtered, c)
			}
		}
		repos[i].Commits = filtered
	}
	return repos
}

func computeBucket(repos []model.Repository, startDate, endDate time.Time, simple bool) *aggregator.AggBucket {
	if startDate.IsZero() && endDate.IsZero() {
		return nil
	}

	if simple {
		acc := aggregator.NewSimpleAccumulator(startDate, endDate)
		for ri := range repos {
			repo := &repos[ri]
			for ci := range repo.Commits {
				acc.Add(&repo.Commits[ci], repo)
			}
		}
		return acc.Build()
	}

	acc := aggregator.NewAccumulator(startDate, endDate)
	for ri := range repos {
		repo := &repos[ri]
		for ci := range repo.Commits {
			acc.Add(&repo.Commits[ci], repo)
		}
	}
	return acc.Build()
}

func markSelf(bucket *aggregator.AggBucket, email string) {
	if email == "" || bucket == nil {
		return
	}
	for i := range bucket.AuthorList {
		if bucket.AuthorList[i].Email == email {
			bucket.AuthorList[i].IsMe = true
			break
		}
	}
	for i := range bucket.AuthorRank {
		if bucket.AuthorRank[i].Email == email {
			bucket.AuthorRank[i].IsMe = true
			break
		}
	}
}

func overviewFromBucket(bucket *aggregator.AggBucket) model.OverviewStats {
	if bucket == nil {
		return model.OverviewStats{}
	}
	return model.OverviewStats{
		TotalCommits:    bucket.TotalCommits,
		TotalAdditions:  bucket.TotalAdditions,
		TotalDeletions:  bucket.TotalDeletions,
		ActiveAuthors:   bucket.ActiveAuthors,
		RepositoryCount: bucket.RepositoryCount,
		Authors:         bucket.AuthorList,
	}
}
