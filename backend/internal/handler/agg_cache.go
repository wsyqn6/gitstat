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

func mergeBuckets(buckets []*aggregator.AggBucket, startDate, endDate time.Time) *aggregator.AggBucket {
	result := &aggregator.AggBucket{
		BuiltAt:         time.Now(),
		TotalCommits:    0,
		TotalAdditions:  0,
		TotalDeletions:  0,
		ActiveAuthors:   0,
		RepositoryCount: len(buckets),
	}

	authorMap := make(map[string]*model.AuthorRankItem)
	dailyByRepo := make([]model.RepositoryDailyStats, 0)
	monthlyByRepo := make([]model.RepositoryPeriodStats, 0)
	repoComp := make([]model.RepoComparison, 0)
	heatmap := make(map[int]map[int]int)
	calendar := make(map[string]*model.MonthlyCalendarItem)
	fileRank := make(map[string]*model.FileRankItem)

	for _, b := range buckets {
		result.TotalCommits += b.TotalCommits
		result.TotalAdditions += b.TotalAdditions
		result.TotalDeletions += b.TotalDeletions

		for _, author := range b.AuthorList {
			entry, ok := authorMap[author.Email]
			if !ok {
				entry = &model.AuthorRankItem{
					Author: author.Author,
					Email:  author.Email,
				}
				authorMap[author.Email] = entry
			}
			entry.Commits += author.Commits
			entry.Additions += author.Additions
			entry.Deletions += author.Deletions
			entry.NetChange = entry.Additions - entry.Deletions
			if author.LastCommitDate != "" {
				entry.LastCommitDate = author.LastCommitDate
			}
		}

		for _, dr := range b.DailyByRepo {
			filtered := filterDailyByRepo(&dr, startDate, endDate)
			if len(filtered.DailyCommits) > 0 || len(filtered.Authors) > 0 {
				dailyByRepo = append(dailyByRepo, *filtered)
			}
		}

		for _, mr := range b.MonthlyByRepo {
			filtered := filterMonthlyByRepo(&mr, startDate, endDate)
			if len(filtered.Authors) > 0 {
				monthlyByRepo = append(monthlyByRepo, *filtered)
			}
		}

		for _, rc := range b.RepoComp {
			repoComp = append(repoComp, rc)
		}

		for _, hm := range b.Heatmap {
			if heatmap[hm.DayOfWeek] == nil {
				heatmap[hm.DayOfWeek] = make(map[int]int)
			}
			heatmap[hm.DayOfWeek][hm.Hour] += hm.CommitCount
		}

		for _, cal := range b.Calendar {
			if !isMonthInRange(cal.Month, startDate, endDate) {
				continue
			}
			key := cal.Month
			entry, ok := calendar[key]
			if !ok {
				entry = &model.MonthlyCalendarItem{Month: key}
				calendar[key] = entry
			}
			entry.Commits += cal.Commits
			entry.Additions += cal.Additions
			entry.Deletions += cal.Deletions
		}

		for _, fr := range b.FileRank {
			entry, ok := fileRank[fr.FilePath]
			if !ok {
				entry = &model.FileRankItem{FilePath: fr.FilePath}
				fileRank[fr.FilePath] = entry
			}
			entry.Commits += fr.Commits
			entry.Additions += fr.Additions
			entry.Deletions += fr.Deletions
			entry.NetChange = entry.Additions - entry.Deletions
		}
	}

	result.AuthorList = make([]model.AuthorRankItem, 0, len(authorMap))
	for _, item := range authorMap {
		result.AuthorList = append(result.AuthorList, *item)
	}
	sort.Slice(result.AuthorList, func(i, j int) bool {
		return result.AuthorList[j].Commits < result.AuthorList[i].Commits
	})
	result.ActiveAuthors = len(result.AuthorList)

	result.AuthorRank = make([]model.AuthorRankItem, len(result.AuthorList))
	copy(result.AuthorRank, result.AuthorList)

	sort.Slice(dailyByRepo, func(i, j int) bool {
		return dailyByRepo[i].RepoName < dailyByRepo[j].RepoName
	})
	result.DailyByRepo = dailyByRepo

	sort.Slice(monthlyByRepo, func(i, j int) bool {
		return monthlyByRepo[i].RepoName < monthlyByRepo[j].RepoName
	})
	result.MonthlyByRepo = monthlyByRepo

	sort.Slice(repoComp, func(i, j int) bool {
		return repoComp[j].Commits < repoComp[i].Commits
	})
	result.RepoComp = repoComp

	result.Heatmap = make([]model.ActivityHeatmapPoint, 0)
	for dow, hourMap := range heatmap {
		for hour, count := range hourMap {
			result.Heatmap = append(result.Heatmap, model.ActivityHeatmapPoint{
				DayOfWeek: dow, Hour: hour, CommitCount: count,
			})
		}
	}
	sort.Slice(result.Heatmap, func(i, j int) bool {
		if result.Heatmap[i].DayOfWeek != result.Heatmap[j].DayOfWeek {
			return result.Heatmap[i].DayOfWeek < result.Heatmap[j].DayOfWeek
		}
		return result.Heatmap[i].Hour < result.Heatmap[j].Hour
	})

	result.Calendar = make([]model.MonthlyCalendarItem, 0, len(calendar))
	for _, item := range calendar {
		result.Calendar = append(result.Calendar, *item)
	}
	sort.Slice(result.Calendar, func(i, j int) bool {
		return result.Calendar[i].Month < result.Calendar[j].Month
	})

	result.FileRank = make([]model.FileRankItem, 0, len(fileRank))
	for _, item := range fileRank {
		result.FileRank = append(result.FileRank, *item)
	}
	sort.Slice(result.FileRank, func(i, j int) bool {
		return result.FileRank[j].Commits < result.FileRank[i].Commits
	})

	return result
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

func filterBucketByEmail(bucket *aggregator.AggBucket, email string) *aggregator.AggBucket {
	result := &aggregator.AggBucket{
		BuiltAt:         bucket.BuiltAt,
		TotalCommits:    0,
		TotalAdditions:  0,
		TotalDeletions:  0,
		ActiveAuthors:   0,
		RepositoryCount: bucket.RepositoryCount,
	}

	var authorList []model.AuthorRankItem
	for _, author := range bucket.AuthorList {
		if author.Email == email {
			authorList = append(authorList, author)
			result.TotalCommits += author.Commits
			result.TotalAdditions += author.Additions
			result.TotalDeletions += author.Deletions
		}
	}
	result.AuthorList = authorList
	result.AuthorRank = authorList
	result.ActiveAuthors = len(authorList)

	dailyByRepo := make([]model.RepositoryDailyStats, 0)
	for _, dr := range bucket.DailyByRepo {
		filtered := &model.RepositoryDailyStats{
			RepoName:       dr.RepoName,
			RepoPath:       dr.RepoPath,
			CurrentBranch:  dr.CurrentBranch,
			LastCommitTime: dr.LastCommitTime,
			DailyCommits:   make([]model.DayCommitCount, 0),
			Authors:        make([]model.AuthorDailyStats, 0),
		}
		for _, author := range dr.Authors {
			if author.Email == email {
				filtered.Authors = append(filtered.Authors, author)
				for _, dd := range author.DailyData {
					found := false
					for i := range filtered.DailyCommits {
						if filtered.DailyCommits[i].Date == dd.Date {
							filtered.DailyCommits[i].Commits += dd.Commits
							found = true
							break
						}
					}
					if !found {
						filtered.DailyCommits = append(filtered.DailyCommits, model.DayCommitCount{
							Date: dd.Date, Commits: dd.Commits,
						})
					}
				}
			}
		}
		if len(filtered.Authors) > 0 {
			dailyByRepo = append(dailyByRepo, *filtered)
		}
	}
	result.DailyByRepo = dailyByRepo

	monthlyByRepo := make([]model.RepositoryPeriodStats, 0)
	for _, mr := range bucket.MonthlyByRepo {
		filtered := &model.RepositoryPeriodStats{
			RepoName:       mr.RepoName,
			RepoPath:       mr.RepoPath,
			CurrentBranch:  mr.CurrentBranch,
			LastCommitTime: mr.LastCommitTime,
			Authors:        make([]model.AuthorPeriodStats, 0),
		}
		for _, author := range mr.Authors {
			if author.Email == email {
				filtered.Authors = append(filtered.Authors, author)
			}
		}
		if len(filtered.Authors) > 0 {
			monthlyByRepo = append(monthlyByRepo, *filtered)
		}
	}
	result.MonthlyByRepo = monthlyByRepo

	result.RepoComp = bucket.RepoComp
	result.Heatmap = bucket.Heatmap
	result.Calendar = bucket.Calendar
	result.FileRank = bucket.FileRank

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
