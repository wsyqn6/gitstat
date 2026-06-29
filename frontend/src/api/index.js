import { BASE, buildParams, request } from './client'

export async function setScanPath(path) {
  return request(`${BASE}/scan/path`, {
    method: 'POST',
    body: { path },
    errorMsg: 'Set path failed'
  })
}

export async function getScanPath() {
  const data = await request(`${BASE}/scan/path`, { errorMsg: 'Failed to fetch scan path' })
  return { path: data.data?.path || '', version: data.data?.version || '' }
}

export async function getBuildVersion() {
  const data = await request(`${BASE}/version`, { errorMsg: 'Failed to fetch version' })
  return data.version || 'dev'
}

export async function getOverviewStats(startDate, endDate, repos = [], scope) {
  const params = buildParams({ startDate, endDate, repos, scope })
  const res = await request(`${BASE}/stats/overview?${params}`, { errorMsg: 'Failed to fetch stats' })
  return res.data
}

export async function exportData() {
  return request(`${BASE}/export/json`, { errorMsg: 'Export failed', blob: true })
}

export async function getDailyStats(email, timeRange = 'week', repos = [], startDate, endDate) {
  const params = buildParams({ email, timeRange, repos, startDate, endDate })
  const res = await request(`${BASE}/stats/daily?${params}`, { errorMsg: 'Failed to fetch daily stats' })
  return res.data
}

export async function getAuthorRank(repos = [], startDate, endDate, timeRange = 'week') {
  const params = buildParams({ repos, startDate, endDate, timeRange })
  const res = await request(`${BASE}/stats/authors?${params}`, { errorMsg: 'Failed to fetch author rank' })
  return res.data
}

export async function getMonthlyStats(email, timeRange = 'month', repos = [], startDate, endDate) {
  const params = buildParams({ email, timeRange, repos, startDate, endDate })
  const res = await request(`${BASE}/stats/monthly?${params}`, { errorMsg: 'Failed to fetch monthly stats' })
  return res.data
}

export async function getComparisonStats(startDate, endDate, prevStartDate, prevEndDate, repos = [], scope) {
  const params = buildParams({ startDate, endDate, prevStartDate, prevEndDate, repos, scope })
  const res = await request(`${BASE}/stats/compare?${params}`, { errorMsg: 'Failed to fetch comparison' })
  return res.data
}

export async function getActivityHeatmap(repos = [], startDate, endDate) {
  const params = buildParams({ repos, startDate, endDate })
  const res = await request(`${BASE}/stats/activity-heatmap?${params}`, { errorMsg: 'Failed to fetch activity heatmap' })
  return res.data
}

export async function getRepoComparison(repos = [], startDate, endDate, timeRange = 'week') {
  const params = buildParams({ repos, startDate, endDate, timeRange })
  const res = await request(`${BASE}/stats/repo-comparison?${params}`, { errorMsg: 'Failed to fetch repo comparison' })
  return res.data
}

export async function getFileRanking(repos = [], startDate, endDate, limit, timeRange) {
  const params = buildParams({ repos, startDate, endDate, limit, timeRange })
  const res = await request(`${BASE}/stats/file-ranking?${params}`, { errorMsg: 'Failed to fetch file ranking' })
  return res.data
}

export async function getReposList() {
  return request(`${BASE}/repos/list`, { errorMsg: 'Failed to fetch repos list' })
}

export async function getRepoInfo(path) {
  const params = buildParams({ path })
  return request(`${BASE}/repos/info?${params}`, { errorMsg: 'Failed to fetch repo info' })
}

export async function getRepoStats(path) {
  const params = buildParams({ path })
  return request(`${BASE}/repos/stats?${params}`, { errorMsg: 'Failed to fetch repo stats' })
}

export async function analyzeRepo(path) {
  return request(`${BASE}/repos/analyze`, {
    method: 'POST',
    body: { path },
    errorMsg: 'Failed to analyze repo'
  })
}

export async function getRepoChart(path) {
  const params = buildParams({ path })
  return request(`${BASE}/repos/chart?${params}`, { errorMsg: 'Failed to fetch repo chart' })
}

export async function getRepoCommits(path, offset = 0, limit = 30) {
  const params = buildParams({ path, offset, limit })
  return request(`${BASE}/repos/commits?${params}`, { errorMsg: 'Failed to fetch repo commits' })
}

export async function getRepoTagsCount(path) {
  const params = buildParams({ path })
  return request(`${BASE}/repos/tags?${params}`, { errorMsg: 'Failed to fetch repo tags' })
}

export async function getRepoTagsPage(path, offset = 0, limit = 30) {
  const params = buildParams({ path, offset, limit })
  return request(`${BASE}/repos/tags?${params}`, { errorMsg: 'Failed to fetch repo tags' })
}
