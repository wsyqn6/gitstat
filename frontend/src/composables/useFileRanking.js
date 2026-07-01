import { ref } from 'vue'
import { getFileRanking } from '../api'

const INITIAL_LIMIT = 5

export function useFileRanking() {
  const data = ref([])
  const loading = ref(false)
  const loaded = ref(false)
  const hasMore = ref(false)
  const limit = ref(INITIAL_LIMIT)

  async function load(repos, startDate = '', endDate = '') {
    if (loading.value) return
    loading.value = true
    try {
      const res = await getFileRanking(repos, startDate, endDate, limit.value)
      if (res) {
        const existing = new Set(data.value.map(i => i.filePath))
        const newItems = res.filter(i => !existing.has(i.filePath))
        data.value = [...data.value, ...newItems]
        hasMore.value = res.length >= limit.value
      }
      loaded.value = true
    } catch (err) {
      console.error('Failed to load file ranking:', err)
    } finally {
      loading.value = false
    }
  }

  function setData(raw) {
    data.value = raw || []
    hasMore.value = (raw || []).length >= limit.value
    loaded.value = true
  }

  function loadMore(repos, startDate = '', endDate = '') {
    limit.value += 5
    return load(repos, startDate, endDate)
  }

  function reset() {
    data.value = []
    loading.value = false
    loaded.value = false
    hasMore.value = false
    limit.value = INITIAL_LIMIT
  }

  return { data, loading, loaded, hasMore, limit, load, setData, loadMore, reset }
}
