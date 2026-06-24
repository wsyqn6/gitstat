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
          @update:selectedRepos="selectedRepos = $event"
          @update:selectedTimeRange="selectedTimeRange = $event"
          @update:showCustomPicker="showCustomPicker = $event"
          @update:customStartDate="customStartDate = $event"
          @update:customEndDate="customEndDate = $event"
          @analyze="loadData"
        />
      </div>
      <div class="view-toggle-bar">
        <div class="view-toggle-inner">
          <div class="toggle-slider" :class="viewMode"></div>
          <button
            @click="viewMode = 'chart'"
            class="view-toggle-btn"
            :class="{ active: viewMode === 'chart' }"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="toggle-icon">
              <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline>
            </svg>
            <span>{{ t('calendar.chartView') }}</span>
          </button>
          <button
            @click="viewMode = 'calendar'"
            class="view-toggle-btn"
            :class="{ active: viewMode === 'calendar' }"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="toggle-icon">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="16" y1="2" x2="16" y2="6"></line>
              <line x1="8" y1="2" x2="8" y2="6"></line>
              <line x1="3" y1="10" x2="21" y2="10"></line>
            </svg>
            <span>{{ t('calendar.calendarView') }}</span>
          </button>
        </div>
      </div>
    </div>

    <OverviewCards
      :overviewStats="overviewStats"
      :loadedTimeRange="loadedTimeRange"
      :customStartDate="customStartDate"
      :customEndDate="customEndDate"
    />

    <AnalyticsCharts
      v-show="viewMode === 'chart'"
      :loading="loading"
      :dailyStats="dailyStats"
      :selectedRepos="selectedRepos"
      :authorRank="authorRank"
      :activityHeatmap="activityHeatmap"
    />

    <CalendarView
      v-show="viewMode === 'calendar'"
      :viewType="calendarViewType"
      :dailyStats="dailyStats"
      :startDate="currentStartDate"
      :endDate="currentEndDate"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { getDailyStats, getRepositories, getOverviewStats, getAuthorRank, getActivityHeatmap } from '../api'
import AnalyticsControls from '../components/AnalyticsControls.vue'
import OverviewCards from '../components/OverviewCards.vue'
import AnalyticsCharts from '../components/AnalyticsCharts.vue'
import CalendarView from '../components/CalendarView.vue'
import { useI18n } from '../i18n'

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
const authorRank = ref([])
const activityHeatmap = ref([])
const viewMode = ref('chart')
const currentStartDate = ref('')
const currentEndDate = ref('')
const loadedTimeRange = ref('')
const calendarViewType = ref('week')

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

    const [overview, stats, authors, heatmap] = await Promise.all([
      getOverviewStats(startDate, endDate, selectedRepos.value),
      getDailyStats(null, selectedTimeRange.value === 'custom' ? '' : selectedTimeRange.value, selectedRepos.value, startDate, endDate),
      getAuthorRank(selectedRepos.value, startDate, endDate),
      getActivityHeatmap(selectedRepos.value, startDate, endDate)
    ])
    overviewStats.value = overview
    dailyStats.value = stats || []
    authorRank.value = authors || []
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
    repositories.value = await getRepositories()
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

.view-toggle-bar {
  text-align: center;
  margin-top: 0.75rem;
}

.view-toggle-inner {
  display: inline-flex;
  background: var(--glass-btn-bg);
  backdrop-filter: blur(var(--glass-blur));
  border: 1px solid var(--glass-btn-border);
  border-radius: calc(var(--radius-btn) + 2px);
  padding: 3px;
  gap: 2px;
  position: relative;
}

.toggle-slider {
  position: absolute;
  top: 3px;
  left: 3px;
  width: calc(50% - 4px);
  height: calc(100% - 6px);
  border-radius: calc(var(--radius-btn) - 2px);
  background: var(--glass-btn-hover-bg);
  backdrop-filter: blur(var(--glass-blur));
  box-shadow: var(--glass-btn-shadow), var(--glass-btn-inner);
  transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
  z-index: 0;
  pointer-events: none;
}

.toggle-slider.calendar {
  transform: translateX(calc(100% + 2px));
}

.view-toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0.5rem 1.2rem;
  border: none;
  background: transparent;
  color: var(--color-text-muted);
  font-family: var(--font-body);
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: color 0.25s ease;
  position: relative;
  z-index: 1;
}

.view-toggle-btn.active {
  color: var(--glass-btn-color);
}

.view-toggle-btn:hover:not(.active) {
  color: var(--color-text-secondary);
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
