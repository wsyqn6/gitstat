<template>
  <div v-if="overviewStats" class="overview-section">
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
          <div class="stat-value">{{ overviewStats.totalCommits }}</div>
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
          <div class="stat-value">{{ overviewStats.totalAdditions }}</div>
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
          <div class="stat-value">{{ overviewStats.totalDeletions }}</div>
        </div>
      </div>
      <div class="glass stat-card clickable" :class="{ expanded: expandedSection === 'authors' }" @click="emit('toggle-section', 'authors')">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
            <circle cx="9" cy="7" r="4"></circle>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('analytics.activeAuthors') }}</div>
          <div class="stat-value">{{ overviewStats.activeAuthors }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from '../i18n'

const { t } = useI18n()

const props = defineProps({
  overviewStats: Object,
  loadedTimeRange: String,
  expandedSection: String,
  customStartDate: String,
  customEndDate: String
})

defineEmits(['toggle-section'])

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

.stat-card:hover {
  border-color: var(--border-card-hover);
  transform: translateY(-2px);
}

.stat-card.clickable {
  cursor: pointer;
}

.stat-card.clickable.expanded {
  border-color: var(--color-accent);
}

.stat-icon {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-accent);
  filter: drop-shadow(var(--shadow-accent-drop));
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
</style>
