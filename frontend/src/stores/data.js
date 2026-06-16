import { reactive } from 'vue'
import { setScanPath, getScanPath, getOverviewStats, getDailyStats, getRepoComparison, getAuthorRank, getReposList, getRepoInfo, getRepoStats, getRepoChart, analyzeRepo } from '../api'
import LRUCache from '../utils/lruCache'

export const state = reactive({
  overviewStats: null,
  dailyStats: [],
  repoDailyTrend: [],
  repoComparison: [],
  authorRank: [],
  scanPath: '',
  loading: false,
  error: null,
  reposInfo: [],
  analyzing: false,
  analyzeCache: new LRUCache(50),
  repoInfoCache: new LRUCache(50),
  repoStatsCache: new LRUCache(50),
  repoChartCache: new LRUCache(50)
})

export async function performScan(path) {
  state.loading = true
  state.error = null
  
  try {
    await setScanPath(path)
    const [overview, daily] = await Promise.all([
      getOverviewStats(null, null, ['all']),
      getDailyStats('', 'today', ['all'])
    ])
    state.overviewStats = overview
    state.dailyStats = daily
  } catch (err) {
    state.error = err.message
  } finally {
    state.loading = false
  }
}

export async function fetchOverviewStats() {
  try {
    state.overviewStats = await getOverviewStats(null, null, ['all'])
  } catch (err) {
    console.error('Failed to fetch overview stats:', err)
  }
}

export async function fetchDailyStatsToday() {
  try {
    state.dailyStats = await getDailyStats('', 'today', ['all'])
  } catch (err) {
    console.error('Failed to fetch daily stats:', err)
  }
}

export async function loadDashboardS1() {
  await Promise.all([
    fetchOverviewStats(),
    fetchRepoDailyTrend(),
    fetchAuthorRank()
  ])
}

export async function loadDashboardS2() {
  await Promise.all([
    fetchRepoComparison(),
    fetchDailyStatsToday()
  ])
}

export async function fetchRepoDailyTrend() {
  try {
    const daily = await getDailyStats('', 'week', ['all'])
    if (!Array.isArray(daily)) {
      state.repoDailyTrend = []
      return
    }
    const result = []
    for (const repo of daily) {
      const dateMap = {}
      for (const author of repo.authors) {
        for (const day of (author.dailyData || [])) {
          if (!dateMap[day.date]) dateMap[day.date] = 0
          dateMap[day.date] += day.commits
        }
      }
      const sorted = Object.entries(dateMap)
        .map(([date, commits]) => ({ date, commits }))
        .sort((a, b) => a.date.localeCompare(b.date))
      if (sorted.length > 0) {
        result.push({ repoName: repo.repoName, data: sorted, authors: repo.authors })
      }
    }
    state.repoDailyTrend = result
  } catch (err) {
    console.error('Failed to fetch repo daily trend:', err)
  }
}

export async function fetchRepoComparison() {
  try {
    state.repoComparison = await getRepoComparison(['all'], null, null, 'week')
  } catch (err) {
    console.error('Failed to fetch repo comparison:', err)
  }
}

export async function fetchAuthorRank() {
  try {
    state.authorRank = await getAuthorRank(['all'], null, null, 'week')
  } catch (err) {
    console.error('Failed to fetch author rank:', err)
  }
}

export async function fetchReposInfo() {
  try {
    state.reposInfo = await getReposList()
  } catch (err) {
    console.error('Failed to fetch repos info:', err)
  }
}

export async function fetchRepoInfo(path) {
  if (state.repoInfoCache.has(path)) {
    return state.repoInfoCache.get(path)
  }
  try {
    const result = await getRepoInfo(path)
    state.repoInfoCache.set(path, result)
    return result
  } catch (err) {
    console.error('Failed to fetch repo info:', err)
    throw err
  }
}

export async function fetchRepoStats(path) {
  if (state.repoStatsCache.has(path)) {
    return state.repoStatsCache.get(path)
  }
  try {
    const result = await getRepoStats(path)
    state.repoStatsCache.set(path, result)
    return result
  } catch (err) {
    console.error('Failed to fetch repo stats:', err)
    throw err
  }
}

export async function fetchRepoChart(path) {
  if (state.repoChartCache.has(path)) {
    return state.repoChartCache.get(path)
  }
  try {
    const result = await getRepoChart(path)
    state.repoChartCache.set(path, result)
    return result
  } catch (err) {
    console.error('Failed to fetch repo chart:', err)
    throw err
  }
}

export async function triggerAnalyze(path) {
  if (state.analyzeCache.has(path)) {
    return state.analyzeCache.get(path)
  }
  state.analyzing = true
  try {
    const result = await analyzeRepo(path)
    state.analyzeCache.set(path, result)
    return result
  } catch (err) {
    console.error('Failed to analyze repo:', err)
    throw err
  } finally {
    state.analyzing = false
  }
}

export async function loadScanPath() {
  try {
    const info = await getScanPath()
    state.scanPath = info.path
  } catch (err) {
    console.error('Failed to load scan path:', err)
  }
}
