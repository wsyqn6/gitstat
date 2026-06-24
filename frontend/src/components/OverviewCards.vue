<template>
  <div class="overview-section">
    <div class="overview-period-label">{{ timePeriodLabel }}{{ t('analytics.overviewTitle') }}</div>
    <div class="overview-cards">
      <div class="glass stat-card">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('analytics.totalCommits') }}</div>
          <div class="stat-value">{{ displayValue(overviewStats?.totalCommits) }}</div>
        </div>
      </div>
      <div class="glass stat-card">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="19" x2="12" y2="5"></line>
            <polyline points="5 12 12 5 19 12"></polyline>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('analytics.totalAdditions') }}</div>
          <div class="stat-value">{{ displayValue(overviewStats?.totalAdditions) }}</div>
        </div>
      </div>
      <div class="glass stat-card">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <polyline points="19 12 12 19 5 12"></polyline>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('analytics.totalDeletions') }}</div>
          <div class="stat-value">{{ displayValue(overviewStats?.totalDeletions) }}</div>
        </div>
      </div>
      <div class="glass stat-card clickable" :class="{ expanded: expanded }" @click="toggleExpand">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
            <circle cx="9" cy="7" r="4"></circle>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('analytics.activeAuthors') }} <span class="expand-icon">{{ expanded ? '▼' : '▶' }}</span></div>
          <div class="stat-value">{{ displayValue(overviewStats?.activeAuthors) }}</div>
        </div>
      </div>
    </div>

    <div v-if="expanded && overviewStats?.authors" class="expand-panel glass card">
      <h4>{{ timePeriodLabel }}{{ t('dashboard.authorRank') }} · {{ t('calendar.total') }} {{ overviewStats.authors.length }}</h4>
      <div class="contrib-table">
        <div class="contrib-header">
          <span>{{ t('dashboard.author') }}</span>
          <span>{{ t('dashboard.commits') }}</span>
          <span class="add">{{ t('analytics.charts.additions') }}</span>
          <span class="del">{{ t('analytics.charts.deletions') }}</span>
          <span>{{ t('analytics.charts.netChange') }}</span>
        </div>
        <div v-for="author in overviewStats.authors" :key="author.email" class="contrib-row">
          <span class="contrib-name">{{ author.author }}</span>
          <span>{{ author.commits }}</span>
          <span class="add">+{{ author.additions }}</span>
          <span class="del">-{{ author.deletions }}</span>
          <span :class="author.netChange >= 0 ? 'add' : 'del'">{{ author.netChange >= 0 ? '+' : '' }}{{ author.netChange }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '../i18n'

const { t } = useI18n()

const props = defineProps({
  overviewStats: Object,
  loadedTimeRange: String,
  customStartDate: String,
  customEndDate: String
})

const expanded = ref(false)

function toggleExpand() {
  expanded.value = !expanded.value
}

function displayValue(val) {
  return val !== undefined && val !== null ? val : '--'
}

const timePeriodLabel = computed(() => {
  switch (props.loadedTimeRange) {
    case 'week': return t('analytics.thisWeek')
    case 'lastWeek': return t('analytics.lastWeek')
    case 'month': return t('analytics.thisMonth')
    case 'lastMonth': return t('analytics.lastMonth')
    case 'year': return t('analytics.thisYear')
    case 'custom': return props.customStartDate && props.customEndDate
      ? `${props.customStartDate} ~ ${props.customEndDate}`
      : t('analytics.customPeriod')
    default: return ''
  }
})
</script>

<style scoped>
.overview-section {
  margin-bottom: 2rem;
}

.overview-period-label {
  font-family: var(--font-display);
  font-size: 0.85rem;
  color: var(--color-accent);
  letter-spacing: 2px;
  text-transform: uppercase;
  margin-bottom: 0.75rem;
  opacity: 0.8;
}

.overview-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
}

.stat-card {
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.stat-card.clickable {
  cursor: pointer;
}

.stat-card.clickable:hover {
  border-color: var(--border-card-hover);
}

.stat-card.expanded {
  border-color: var(--border-card-hover);
}

.stat-icon {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-accent);
}

.stat-icon svg {
  width: 24px;
  height: 24px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-label {
  color: var(--color-text-muted);
  font-size: 0.85rem;
  letter-spacing: 1px;
  margin-bottom: 0.25rem;
}

.expand-icon {
  font-size: 0.6rem;
  margin-left: 0.3rem;
  vertical-align: middle;
}

.expand-panel {
  margin-top: 1.5rem;
  padding: 1.5rem;
  animation: slideDown 0.25s ease;
}

.expand-panel h4 {
  font-family: var(--font-display);
  font-size: 0.9rem;
  color: var(--color-primary);
  letter-spacing: 1px;
  margin: 0 0 1rem 0;
}

.contrib-table {
  width: 100%;
}

.contrib-header,
.contrib-row {
  display: grid;
  grid-template-columns: 1.5fr 1fr 1fr 1fr 1fr;
  padding: 0.6rem 0.5rem;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.85rem;
}

.contrib-header {
  background: var(--bg-row-header);
  border-bottom: 1px solid var(--border-card-subtle);
  font-family: var(--font-display);
  font-size: 0.7rem;
  color: var(--color-nav-link);
  letter-spacing: 1px;
  text-transform: uppercase;
}

.contrib-row {
  border-bottom: 1px solid var(--border-row-subtle);
  color: var(--color-text-primary);
}

.contrib-row:hover {
  background: var(--bg-row-hover);
}

.contrib-name {
  font-weight: 600;
  color: var(--color-accent);
}

.add {
  color: var(--color-green);
}

.del {
  color: var(--color-red);
}
</style>
