import { reactive } from 'vue'
import * as api from '../api'

export const state = reactive({
  repositories: [],
  overviewStats: null,
  dailyStats: [],
  loading: false,
  error: null
})

export async function performScan(path, timeRange) {
  state.loading = true
  state.error = null
  
  try {
    const result = await api.scanRepositories(path, timeRange)
    state.repositories = result.repositories
    await refreshStats()
  } catch (err) {
    state.error = err.message
  } finally {
    state.loading = false
  }
}

export async function refreshStats() {
  try {
    state.overviewStats = await api.getOverviewStats()
  } catch (err) {
    console.error('Failed to refresh stats:', err)
  }
}

export async function refreshDailyStats() {
  try {
    state.dailyStats = await api.getDailyStats('')
  } catch (err) {
    console.error('Failed to refresh daily stats:', err)
  }
}
