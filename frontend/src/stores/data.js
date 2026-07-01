import { reactive } from 'vue'
import { setScanPath, getScanPath, getBuildVersion, getDashboardStats, getDailyTrend, getRepoComparison, getAuthorRank, getReposList, getRepoInfo, getRepoStats, getRepoChart, analyzeRepo, getRepoCommits, getRepoTagsCount, getRepoTagsPage } from '../api'
import LRUCache from '../utils/lruCache'

export const state = reactive({
  dashboardData: null,
  repoDailyTrend: [],
  repoComparison: [],
  authorRank: [],
  buildVersion: null,
  gitVersion: null,
  scanPath: '',
  loading: false,
  error: null,
  reposInfo: [],
  analyzing: false,
  analyzeCache: new LRUCache(50),
  repoInfoCache: new LRUCache(50),
  repoStatsCache: new LRUCache(50),
  repoChartCache: new LRUCache(50),
  repoTagsCache: new LRUCache(50)
})

function clearState() {
  state.dashboardData = null
  state.repoDailyTrend = []
  state.repoComparison = []
  state.authorRank = []
}

export async function performScan(path) {
  clearState()
  state.loading = true
  state.error = null
  
  try {
    await setScanPath(path)
    const dashData = await getDashboardStats(['all'])
    state.dashboardData = dashData
  } catch (err) {
    state.error = err.message
  } finally {
    state.loading = false
  }
}

export async function fetchDashboardData() {
  if (state.dashboardData) return
  try {
    state.dashboardData = await getDashboardStats(['all'])
  } catch (err) {
    console.error('Failed to fetch dashboard data:', err)
  }
}

export async function loadDashboardS1() {
  await Promise.all([
    fetchDashboardData(),
    fetchRepoDailyTrend(),
    fetchAuthorRank()
  ])
}

export async function loadDashboardS2() {
  await Promise.all([
    fetchRepoComparison()
  ])
}

export async function fetchRepoDailyTrend() {
  try {
    const daily = await getDailyTrend(['all'], null, null, 'week')
    if (!Array.isArray(daily)) {
      state.repoDailyTrend = []
      return
    }
    const result = []
    for (const repo of daily) {
      if (repo.dailyCommits && repo.dailyCommits.length > 0) {
        result.push({ repoName: repo.repoName, data: repo.dailyCommits })
      }
    }
    state.repoDailyTrend = result
  } catch (err) {
    console.error('Failed to fetch repo daily trend:', err)
  }
}

export async function fetchRepoComparison() {
  if (state.repoComparison.length > 0) return
  try {
    state.repoComparison = await getRepoComparison(['all'])
  } catch (err) {
    console.error('Failed to fetch repo comparison:', err)
  }
}

export async function fetchAuthorRank() {
  if (state.authorRank.length > 0) return
  try {
    state.authorRank = await getAuthorRank(['all'])
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

export async function fetchRepoCommits(path, offset = 0, limit = 30) {
  try {
    return await getRepoCommits(path, offset, limit)
  } catch (err) {
    console.error('Failed to fetch repo commits:', err)
    throw err
  }
}

export async function fetchRepoTagsCount(path) {
  if (state.repoTagsCache.has(path)) {
    return state.repoTagsCache.get(path)
  }
  try {
    const result = await getRepoTagsCount(path)
    state.repoTagsCache.set(path, result)
    return result
  } catch (err) {
    console.error('Failed to fetch repo tags count:', err)
    throw err
  }
}

export async function fetchRepoTagsPage(path, offset = 0, limit = 30) {
  try {
    return await getRepoTagsPage(path, offset, limit)
  } catch (err) {
    console.error('Failed to fetch repo tags page:', err)
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

export async function initApp() {
  try {
    const [scanInfo, buildInfo] = await Promise.all([
      getScanPath(),
      getBuildVersion()
    ])
    state.scanPath = scanInfo.path
    state.gitVersion = scanInfo.version
    state.buildVersion = buildInfo
  } catch (err) {
    console.error('App init failed:', err)
  }
}

export async function loadScanPath() {
  try {
    const info = await getScanPath()
    state.scanPath = info.path
    state.gitVersion = info.version
  } catch (err) {
    console.error('Failed to load scan path:', err)
  }
}
