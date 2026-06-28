<template>
  <div class="analytics">
    <div class="header-section">
      <div class="title-row">
        <div>
          <h2 class="page-title">{{ t('trends.title') }}</h2>
          <p class="subtitle">{{ t('trends.subtitle') }}</p>
        </div>
        <AnalyticsControls
          :repositories="repositories"
          :selectedRepos="selectedRepos"
          :selectedTimeRange="selectedTimeRange"
          :showCustomPicker="showCustomPicker"
          :customStartDate="customStartDate"
          :customEndDate="customEndDate"
          :loading="loading"
          singleSelect
          @update:selectedRepos="selectedRepos = $event"
          @update:selectedTimeRange="selectedTimeRange = $event"
          @update:showCustomPicker="showCustomPicker = $event"
          @update:customStartDate="customStartDate = $event"
          @update:customEndDate="customEndDate = $event"
          @analyze="loadData"
        />
      </div>
    </div>

    <OverviewCards
      :overviewStats="overviewStats"
      :comparison="comparison"
      :loadedTimeRange="loadedTimeRange"
      :customStartDate="customStartDate"
      :customEndDate="customEndDate"
    />

    <div class="section">
      <div class="section-header">
        <h3 class="section-title">{{ t('trends.timeSeriesTitle') }}</h3>
        <div class="view-toggle-bar">
          <SegmentToggle
            :options="viewOptions"
            v-model="viewMode"
            v-slot="{ active }"
          >
            <button
              v-for="opt in viewOptions"
              :key="opt.value"
              :class="{ active: active === opt.value }"
            >
              <svg v-if="opt.value === 'chart'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="toggle-icon">
                <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="toggle-icon">
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                <line x1="16" y1="2" x2="16" y2="6"></line>
                <line x1="8" y1="2" x2="8" y2="6"></line>
                <line x1="3" y1="10" x2="21" y2="10"></line>
              </svg>
              <span>{{ opt.label }}</span>
            </button>
          </SegmentToggle>
        </div>
      </div>

      <TimeSeriesCharts
        v-show="viewMode === 'chart'"
        :loading="loading"
        :dailyStats="dailyStats"
        :selectedRepos="selectedRepos"
        :startDate="currentStartDate"
        :endDate="currentEndDate"
      />

      <CalendarView
        v-show="viewMode === 'calendar'"
        :viewType="calendarViewType"
        :dailyStats="dailyStats"
        :startDate="currentStartDate"
      />
    </div>

    <div class="section">
      <AnalyticsCharts
        :loading="loading"
        :dailyStats="dailyStats"
        :selectedRepos="selectedRepos"
        :activityHeatmap="activityHeatmap"
        :startDate="currentStartDate"
        :endDate="currentEndDate"
      />
    </div>

    <div class="section">
      <FileRanking
        :loading="loading"
        :data="fileRanking"
        :hasMore="fileRankingHasMore"
        @loadMore="loadFileRanking()"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { getDailyStats, getMonthlyStats, getReposList, getOverviewStats, getActivityHeatmap, getFileRanking, getComparisonStats } from '../api'
import AnalyticsControls from '../components/AnalyticsControls.vue'
import OverviewCards from '../components/OverviewCards.vue'
import TimeSeriesCharts from '../components/TimeSeriesCharts.vue'
import AnalyticsCharts from '../components/AnalyticsCharts.vue'
import CalendarView from '../components/CalendarView.vue'
import FileRanking from '../components/FileRanking.vue'
import { useFileRanking } from '../composables/useFileRanking'
import { useI18n } from '../i18n'
import SegmentToggle from '../components/SegmentToggle.vue'

const { t } = useI18n()
const loading = ref(false)
const overviewStats = ref(null)
const selectedTimeRange = ref('week')
const customStartDate = ref('')
const customEndDate = ref('')
const showCustomPicker = ref(false)
const selectedRepos = ref([])
const repositories = ref([])
const dailyStats = ref([])
const activityHeatmap = ref([])
const { data: fileRanking, hasMore: fileRankingHasMore, limit: fileRankingLimit, load, setData } = useFileRanking()
const viewMode = ref('chart')
const viewOptions = [
  { value: 'chart', label: t('calendar.chartView') },
  { value: 'calendar', label: t('calendar.calendarView') }
]
const currentStartDate = ref('')
const currentEndDate = ref('')
const loadedTimeRange = ref('')
const calendarViewType = ref('week')
const comparison = ref(null)

function toLocalDateStr(d) {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function computeCalendarViewType() {
  if (!currentStartDate.value || !currentEndDate.value) return 'week'
  if (selectedTimeRange.value === 'year') return 'year'
  if (selectedTimeRange.value === 'month' || selectedTimeRange.value === 'lastMonth') return 'month'
  if (selectedTimeRange.value === 'week' || selectedTimeRange.value === 'lastWeek' || selectedTimeRange.value === 'today') return 'week'
  const days = Math.round((new Date(currentEndDate.value) - new Date(currentStartDate.value)) / (1000 * 60 * 60 * 24))
  if (days <= 7) return 'week'
  if (days <= 31) return 'month'
  return 'year'
}

const loadFileRanking = async () => {
  if (!currentStartDate.value || !currentEndDate.value) return
  fileRankingLimit.value += 5
  await load(selectedRepos.value, currentStartDate.value, currentEndDate.value)
  await nextTick()
  document.querySelector('.file-ranking')?.scrollIntoView({ behavior: 'smooth', block: 'end' })
}

const loadData = async () => {
  if (selectedRepos.value.length === 0) return

  loading.value = true
  try {
    let startDate, endDate
    const now = new Date()

    switch (selectedTimeRange.value) {
      case 'today':
        startDate = toLocalDateStr(now)
        endDate = toLocalDateStr(now)
        break
      case 'week':
        const currentDay = now.getDay() || 7
        const monday = new Date(now)
        monday.setDate(now.getDate() - currentDay + 1)
        startDate = toLocalDateStr(monday)
        endDate = toLocalDateStr(now)
        break
      case 'lastWeek':
        const lastWeekMonday = new Date(now)
        lastWeekMonday.setDate(now.getDate() - (now.getDay() || 7) - 6)
        const lastWeekSunday = new Date(lastWeekMonday)
        lastWeekSunday.setDate(lastWeekMonday.getDate() + 6)
        startDate = toLocalDateStr(lastWeekMonday)
        endDate = toLocalDateStr(lastWeekSunday)
        break
      case 'month':
        const firstDay = new Date(now.getFullYear(), now.getMonth(), 1)
        startDate = toLocalDateStr(firstDay)
        endDate = toLocalDateStr(now)
        break
      case 'lastMonth':
        const firstDayLast = new Date(now.getFullYear(), now.getMonth() - 1, 1)
        const lastDayLast = new Date(now.getFullYear(), now.getMonth(), 0)
        startDate = toLocalDateStr(firstDayLast)
        endDate = toLocalDateStr(lastDayLast)
        break
      case 'year':
        const firstDayOfYear = new Date(now.getFullYear(), 0, 1)
        startDate = toLocalDateStr(firstDayOfYear)
        endDate = toLocalDateStr(now)
        break
      case 'custom':
        startDate = customStartDate.value
        endDate = customEndDate.value
        break
      default:
        startDate = null
        endDate = null
    }

    currentStartDate.value = startDate
    currentEndDate.value = endDate

    const days = Math.round((new Date(endDate) - new Date(startDate)) / (1000 * 60 * 60 * 24))
    const isMonth = days > 31
    const timeRange = selectedTimeRange.value === 'custom' ? '' : selectedTimeRange.value
    const statsPromise = isMonth
      ? getMonthlyStats(null, '', selectedRepos.value, startDate, endDate)
      : getDailyStats(null, timeRange, selectedRepos.value, startDate, endDate)

    let prevStartDate = '', prevEndDate = ''
    const range = selectedTimeRange.value
    if (range === 'week') {
      const ms = 7 * 86400 * 1000
      prevStartDate = toLocalDateStr(new Date(new Date(startDate).getTime() - ms))
      prevEndDate = toLocalDateStr(new Date(new Date(endDate).getTime() - ms))
    } else if (range === 'month') {
      const s = new Date(startDate)
      const e = new Date(endDate)
      prevStartDate = toLocalDateStr(new Date(s.getFullYear(), s.getMonth() - 1, 1))
      prevEndDate = toLocalDateStr(new Date(e.getFullYear(), e.getMonth(), 0))
    } else if (range === 'year') {
      const sy = new Date(startDate).getFullYear() - 1
      const ey = new Date(endDate).getFullYear() - 1
      prevStartDate = `${sy}-01-01`
      prevEndDate = `${ey}-12-31`
    }
    const comparisonPromise = prevStartDate && prevEndDate
      ? getComparisonStats(startDate, endDate, prevStartDate, prevEndDate, selectedRepos.value, 'all')
      : Promise.resolve(null)

    const [overview, rawStats, heatmap, fileRankRaw, comparisonResult] = await Promise.all([
      getOverviewStats(startDate, endDate, selectedRepos.value, 'all'),
      statsPromise,
      getActivityHeatmap(selectedRepos.value, startDate, endDate),
      getFileRanking(selectedRepos.value, startDate, endDate, fileRankingLimit.value),
      comparisonPromise
    ])
    overviewStats.value = overview
    comparison.value = comparisonResult
    setData(fileRankRaw)
    dailyStats.value = (rawStats || []).map(repo => ({
      ...repo,
      authors: (repo.authors || []).map(a => ({
        ...a,
        dailyData: a.dailyData || a.periodData || []
      }))
    }))
    activityHeatmap.value = heatmap || []

    await nextTick()
    loadedTimeRange.value = selectedTimeRange.value
    calendarViewType.value = computeCalendarViewType()
    loading.value = false
  } catch (error) {
    console.error('Failed to load analytics data:', error)
    loading.value = false
  }
}

onMounted(async () => {
  try {
    repositories.value = await getReposList()
    if (repositories.value.length === 1) {
      selectedRepos.value = [repositories.value[0].path]
    }
  } catch (error) {
    console.error('Failed to load repositories:', error)
  }
})
</script>

<style scoped>
.analytics {
  max-width: 1600px;
  margin: 0 auto;
  padding: 0 2rem 1.5rem;
}

.header-section {
  margin-bottom: 1.5rem;
  position: relative;
}

.title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
  flex-wrap: wrap;
}

.page-title {
  font-family: var(--font-display);
  font-size: 2rem;
  font-weight: 700;
  background: var(--gradient-page-title);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 3px;
  text-transform: uppercase;
  margin: 0 0 0.25rem 0;
  filter: drop-shadow(var(--shadow-glow-btn));
}

.subtitle {
  font-size: 0.85rem;
  color: var(--color-text-muted);
  letter-spacing: 4px;
  text-transform: uppercase;
  margin: 0;
}

.section {
  margin-top: 2rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  gap: 1rem;
}

.section-title {
  font-family: var(--font-display);
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-text-secondary);
  letter-spacing: 2px;
  text-transform: uppercase;
  margin: 0;
}

.view-toggle-bar {
  flex-shrink: 0;
}

.toggle-icon {
  width: 16px;
  height: 16px;
}

@media (max-width: 1200px) {
  .title-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
  }
}

@media (max-width: 768px) {
  .analytics {
    padding: 0 1rem 1rem;
  }

  .page-title {
    font-size: 1.5rem;
  }
}
</style>
