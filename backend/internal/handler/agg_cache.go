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

func getAggBucket(repoPaths []string, startDate, endDate time.Time, email string) *aggregator.AggBucket {
	key := cacheKey(repoPaths, startDate, endDate, email)

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
	bucket := computeAggBucket(repoPaths, startDate, endDate, email)

	aggCacheMu.Lock()
	aggCache[key] = &aggCacheEntry{bucket: bucket, expiresAt: time.Now().Add(3 * time.Second)}
	delete(inflight, key)
	aggCacheMu.Unlock()
	close(ch)

	return bucket
}

func computeAggBucket(repoPaths []string, startDate, endDate time.Time, email string) *aggregator.AggBucket {
	ensureDataLoaded(repoPaths, startDate)

	repos := store.GlobalStore.GetReposWithRange(repoPaths, startDate, endDate)
	repos = filterCommitsByEmail(repos, email)

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

func filterDailyByRepo(dr *model.RepositoryDailyStats, startDate, endDate time.Time) *model.RepositoryDailyStats {
	result := &model.RepositoryDailyStats{
		RepoName:       dr.RepoName,
		RepoPath:       dr.RepoPath,
		CurrentBranch:  dr.CurrentBranch,
		LastCommitTime: dr.LastCommitTime,
		DailyCommits:   make([]model.DayCommitCount, 0),
		Authors:        make([]model.AuthorDailyStats, 0),
	}

	for _, dc := range dr.DailyCommits {
		if isDateInRange(dc.Date, startDate, endDate) {
			result.DailyCommits = append(result.DailyCommits, dc)
		}
	}

	for _, author := range dr.Authors {
		filtered := &model.AuthorDailyStats{
			Author:       author.Author,
			Email:        author.Email,
			Commits:      0,
			Additions:    0,
			Deletions:    0,
			FilesChanged: author.FilesChanged,
			DailyData:    make([]model.DayCommitData, 0),
		}
		for _, dd := range author.DailyData {
			if isDateInRange(dd.Date, startDate, endDate) {
				filtered.Commits += dd.Commits
				filtered.Additions += dd.Additions
				filtered.Deletions += dd.Deletions
				filtered.DailyData = append(filtered.DailyData, dd)
			}
		}
		if filtered.Commits > 0 {
			result.Authors = append(result.Authors, *filtered)
		}
	}

	sort.Slice(result.Authors, func(i, j int) bool {
		return result.Authors[j].Commits < result.Authors[i].Commits
	})

	return result
}

func filterMonthlyByRepo(mr *model.RepositoryPeriodStats, startDate, endDate time.Time) *model.RepositoryPeriodStats {
	result := &model.RepositoryPeriodStats{
		RepoName:       mr.RepoName,
		RepoPath:       mr.RepoPath,
		CurrentBranch:  mr.CurrentBranch,
		LastCommitTime: mr.LastCommitTime,
		Authors:        make([]model.AuthorPeriodStats, 0),
	}

	for _, author := range mr.Authors {
		filtered := &model.AuthorPeriodStats{
			Author:     author.Author,
			Email:      author.Email,
			Commits:    0,
			Additions:  0,
			Deletions:  0,
			PeriodData: make([]model.PeriodCommitData, 0),
		}
		for _, pd := range author.PeriodData {
			if isMonthInRange(pd.Period, startDate, endDate) {
				filtered.Commits += pd.Commits
				filtered.Additions += pd.Additions
				filtered.Deletions += pd.Deletions
				filtered.PeriodData = append(filtered.PeriodData, pd)
			}
		}
		if filtered.Commits > 0 {
			result.Authors = append(result.Authors, *filtered)
		}
	}

	return result
}

func isDateInRange(dateStr string, startDate, endDate time.Time) bool {
	if startDate.IsZero() && endDate.IsZero() {
		return true
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return true
	}
	if !startDate.IsZero() && date.Before(startDate) {
		return false
	}
	if !endDate.IsZero() && date.After(endDate) {
		return false
	}
	return true
}

func isMonthInRange(monthStr string, startDate, endDate time.Time) bool {
	if startDate.IsZero() && endDate.IsZero() {
		return true
	}
	month, err := time.Parse("2006-01", monthStr)
	if err != nil {
		return true
	}
	if !startDate.IsZero() && month.Before(startDate) {
		return false
	}
	if !endDate.IsZero() && month.After(endDate.AddDate(0, 1, 0)) {
		return false
	}
	return true
}

func cacheKey(repoPaths []string, startDate, endDate time.Time, email string) string {
	k := email + "|"
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

func computeBucket(repos []model.Repository, startDate, endDate time.Time) *aggregator.AggBucket {
	if startDate.IsZero() && endDate.IsZero() {
		return nil
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
