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
          <div class="stat-value">{{ displayValue(overviewStats?.totalCommits) }}</div>
          <ChangeBadge v-if="comparison" :change="comparison.totalCommits" :prefix="vsPrefix" />
          <div class="stat-label">{{ t('analytics.totalCommits') }}</div>
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
          <div class="stat-value">{{ displayValue(overviewStats?.totalAdditions) }}</div>
          <ChangeBadge v-if="comparison" :change="comparison.totalAdditions" :prefix="vsPrefix" />
          <div class="stat-label">{{ t('analytics.totalAdditions') }}</div>
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
          <div class="stat-value">{{ displayValue(overviewStats?.totalDeletions) }}</div>
          <ChangeBadge v-if="comparison" :change="comparison.totalDeletions" :prefix="vsPrefix" />
          <div class="stat-label">{{ t('analytics.totalDeletions') }}</div>
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
          <div class="stat-value">{{ displayValue(overviewStats?.activeAuthors) }}</div>
          <ChangeBadge v-if="comparison" :change="comparison.activeAuthors" :prefix="vsPrefix" />
          <div class="stat-label">{{ t('analytics.activeAuthors') }} <span class="expand-icon">{{ expanded ? '▼' : '▶' }}</span></div>
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
          <span class="cell-with-change">
            <span>{{ author.commits }}</span>
            <ChangeBadge :change="getAuthorChange(author.email)?.commits" prefix="" compact />
          </span>
          <span class="add cell-with-change">
            <span>+{{ author.additions }}</span>
            <ChangeBadge :change="getAuthorChange(author.email)?.additions" prefix="" compact />
          </span>
          <span class="del cell-with-change">
            <span>-{{ author.deletions }}</span>
            <ChangeBadge :change="getAuthorChange(author.email)?.deletions" prefix="" compact />
          </span>
          <span :class="author.netChange >= 0 ? 'add' : 'del'">{{ author.netChange >= 0 ? '+' : '' }}{{ author.netChange }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '../i18n'
import ChangeBadge from './ChangeBadge.vue'

const { t } = useI18n()

const props = defineProps({
  overviewStats: Object,
  comparison: Object,
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

const vsPrefix = computed(() => {
  switch (props.loadedTimeRange) {
    case 'week': return t('trends.vsWeek')
    case 'month': return t('trends.vsMonth')
    case 'year': return t('trends.vsYear')
    default: return ''
  }
})

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

function getAuthorChange(email) {
  return props.comparison?.authors?.find(a => a.email === email)?.change
}
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
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
}

.stat-card {
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  animation: cardIn 0.45s cubic-bezier(0.22, 1, 0.36, 1) both;
}

.stat-card:nth-child(1) { animation-delay: 0s; }
.stat-card:nth-child(2) { animation-delay: 0.06s; }
.stat-card:nth-child(3) { animation-delay: 0.12s; }
.stat-card:nth-child(4) { animation-delay: 0.18s; }

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
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  color: var(--color-accent);
  margin-bottom: 0.5rem;
}

.stat-icon svg {
  width: 18px;
  height: 18px;
  stroke-width: 1.8;
}

.stat-content {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-family: var(--font-display);
  font-size: 1.3rem;
  font-weight: 700;
  color: var(--color-text-primary);
  line-height: 1.2;
  margin-bottom: 0.35rem;
}

.stat-label {
  font-family: var(--font-display);
  font-size: 0.65rem;
  color: var(--color-nav-link);
  letter-spacing: 1px;
  text-transform: uppercase;
}

.expand-icon {
  font-size: 0.5rem;
  margin-left: 0.2rem;
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

.cell-with-change {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
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

@keyframes cardIn {
  from {
    opacity: 0;
    transform: translateY(12px) scale(0.97);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
</style>
