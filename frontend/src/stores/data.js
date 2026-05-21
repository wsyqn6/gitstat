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
    // 调用新接口设置路径，触发后端注册仓库
    await api.setScanPath(path)
    // 重新加载数据（后端会懒加载）
    await refreshStats()
  } catch (err) {
    state.error = err.message
  } finally {
    state.loading = false
  }
}

export async function refreshStats() {
  try {
    // 并行请求 overview 和 daily，触发后端懒加载（所有仓库）
    const [overview, daily] = await Promise.all([
      api.getOverviewStats(null, null, ['all']),
      api.getDailyStats('', 'today', ['all'])
    ])
    state.overviewStats = overview
    state.dailyStats = daily
  } catch (err) {
    console.error('Failed to refresh stats:', err)
  }
}

export async function refreshDailyStats() {
  try {
    state.dailyStats = await api.getDailyStats('', 'today', ['all'])
  } catch (err) {
    console.error('Failed to refresh daily stats:', err)
  }
}
