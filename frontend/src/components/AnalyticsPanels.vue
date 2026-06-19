<template>
  <div>
    <div v-if="expandedSection === 'authors' && overviewStats?.authors" class="expand-panel">
      <div class="expand-panel-header">{{ timePeriodPrefix }}{{ t('dashboard.authorRank') }} · {{ t('calendar.total') }} {{ overviewStats.authors.length }} </div>
      <table class="expand-table">
        <thead>
          <tr>
            <th>{{ t('dashboard.author') }}</th>
            <th>{{ t('dashboard.commits') }}</th>
            <th>{{ t('analytics.charts.additions') }}</th>
            <th>{{ t('analytics.charts.deletions') }}</th>
            <th>{{ t('analytics.charts.netChange') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="author in overviewStats.authors" :key="author.email">
            <td class="cell-author">{{ author.author }}</td>
            <td>{{ author.commits }}</td>
            <td class="cell-additions">+{{ author.additions }}</td>
            <td class="cell-deletions">-{{ author.deletions }}</td>
            <td :class="author.netChange >= 0 ? 'cell-additions' : 'cell-deletions'">{{ author.netChange >= 0 ? '+' : '' }}{{ author.netChange }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="insights.length > 0" class="insights-grid">
      <div
        v-for="(insight, idx) in insights"
        :key="insight.section"
        class="insight-card"
        :class="{ clickable: insight.clickable, expanded: insight.clickable && expandedSection === insight.section }"
        @click="insight.clickable && emit('toggle-section', insight.section)"
      >
        <div class="insight-icon"><span :class="'icon-' + insight.icon"></span></div>
        <div class="insight-content">
          <div class="insight-title">{{ insight.title }}</div>
          <div class="insight-value">{{ insight.value }}</div>
          <div class="insight-desc">{{ insight.description }}</div>
        </div>
      </div>
    </div>

    <div v-if="expandedSection === 'repos' && repoComparison.length > 0" class="expand-panel">
      <div class="expand-panel-header">{{ timePeriodPrefix }}{{ t('analytics.panels.activeReposTitle') }} · {{ t('calendar.total') }} {{ repoComparison.filter(r => r.commits > 0).length }} </div>
      <table class="expand-table">
        <thead>
          <tr>
            <th>{{ t('analytics.repoName') }}</th>
            <th>{{ t('dashboard.commits') }}</th>
            <th>{{ t('analytics.charts.authorCount') }}</th>
            <th>{{ t('analytics.charts.additions') }}</th>
            <th>{{ t('analytics.charts.activeDays') }}</th>
            <th>{{ t('analytics.charts.dailyAvg') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="repo in repoComparison" :key="repo.repoPath">
            <td class="cell-author">{{ repo.repoName }}</td>
            <td>{{ repo.commits }}</td>
            <td>{{ repo.authors }}</td>
            <td class="cell-additions">+{{ repo.additions }}</td>
            <td>{{ repo.activeDays }}</td>
            <td>{{ repo.avgCommitsPerDay }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from '../i18n'

const { t } = useI18n()

const props = defineProps({
  expandedSection: String,
  overviewStats: Object,
  loadedTimeRange: String,
  repoComparison: Array,
  authorRank: Array
})

defineEmits(['toggle-section'])

const timePeriodPrefix = computed(() => {
  switch (props.loadedTimeRange) {
    case 'week': return t('analytics.thisWeek')
    case 'lastWeek': return t('analytics.lastWeek')
    case 'month': return t('analytics.thisMonth')
    case 'lastMonth': return t('analytics.lastMonth')
    case 'year': return t('analytics.thisYear')
    case 'custom': return t('analytics.customPeriod')
    default: return ''
  }
})

const insights = computed(() => {
  const result = []
  const prefix = timePeriodPrefix.value

  if (props.authorRank && props.authorRank.length > 0) {
    const topAuthor = props.authorRank[0]
    result.push({
      icon: 'star',
      title: t('analytics.panels.starTitle').replace('{prefix}', prefix),
      value: topAuthor.author,
      description: t('analytics.panels.starDesc').replace('{0}', topAuthor.commits).replace('{1}', topAuthor.additions)
    })
  }

  if (props.overviewStats && props.overviewStats.totalCommits > 0) {
    const avgSize = Math.round((props.overviewStats.totalAdditions + props.overviewStats.totalDeletions) / props.overviewStats.totalCommits)
    result.push({
      icon: 'grid',
      title: t('analytics.panels.avgSizeTitle'),
      value: t('analytics.panels.avgSizeValue').replace('{0}', avgSize),
      description: t('analytics.panels.avgSizeDesc').replace('{prefix}', prefix).replace('{0}', avgSize)
    })
  }

  if (props.overviewStats) {
    const netChange = props.overviewStats.totalAdditions - props.overviewStats.totalDeletions
    result.push({
      icon: 'growth',
      title: t('analytics.panels.netGrowthTitle'),
      value: `${netChange > 0 ? '+' : ''}${netChange}`,
      description: t('analytics.panels.netGrowthDesc').replace('{prefix}', prefix).replace('{0}', netChange)
    })
  }

  if (props.repoComparison && props.repoComparison.length > 0) {
    const activeRepos = props.repoComparison.filter(r => r.commits > 0).length
    result.push({
      icon: 'folder',
      title: t('analytics.panels.activeReposTitle'),
      value: t('analytics.panels.activeReposValue').replace('{0}', activeRepos),
      description: t('analytics.panels.activeReposDesc').replace('{prefix}', prefix).replace('{0}', activeRepos),
      clickable: true,
      section: 'repos'
    })
  }

  return result
})
</script>

<style scoped>
.expand-panel {
  background: var(--bg-panel);
  backdrop-filter: blur(var(--blur-card));
  border: 1px solid var(--border-insight-card);
  border-radius: 12px;
  padding: 1.25rem;
  margin-top: 1rem;
  animation: slideDown 0.25s ease;
  overflow-x: auto;
}

.expand-panel-header {
  font-family: var(--font-display);
  font-size: 0.8rem;
  color: var(--color-accent);
  letter-spacing: 1px;
  margin-bottom: 1rem;
  opacity: 0.9;
}

.expand-table {
  width: 100%;
  border-collapse: collapse;
  font-family: var(--font-body);
  font-size: 0.9rem;
}

.expand-table th {
  text-align: left;
  padding: 0.5rem 0.75rem;
  color: var(--color-text-muted);
  font-weight: 600;
  font-size: 0.8rem;
  letter-spacing: 1px;
  text-transform: uppercase;
  border-bottom: 1px solid var(--bg-stat);
}

.expand-table td {
  padding: 0.5rem 0.75rem;
  color: var(--color-input-text);
  border-bottom: 1px solid rgba(148, 163, 184, 0.1);
}

.expand-table tbody tr:hover {
  background: var(--bg-row-hover);
}

.expand-table .cell-author {
  color: var(--color-accent);
  font-weight: 600;
}

.expand-table .cell-additions {
  color: var(--color-green);
}

.expand-table .cell-deletions {
  color: var(--color-red);
}

.insights-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.insight-card {
  background: var(--bg-tab);
  backdrop-filter: blur(var(--blur-card));
  border: 1px solid var(--border-insight-card);
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: all 0.3s;
  min-height: 100px;
}

.insight-card:hover {
  border-color: rgba(var(--color-primary-rgb), 0.5);
  transform: translateY(-2px);
  box-shadow: var(--shadow-glow-card);
}

.insight-card.clickable {
  cursor: pointer;
}

.insight-card.clickable:hover {
  border-color: rgba(var(--color-primary-rgb), 0.5);
  transform: translateY(-2px);
  box-shadow: var(--shadow-glow-card);
}

.insight-card.clickable.expanded {
  border-color: var(--color-accent);
  box-shadow: 0 0 20px var(--border-dropdown), inset 0 0 20px var(--bg-tag-section);
}

.insight-icon {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-accent);
  filter: drop-shadow(0 0 6px var(--shadow-active));
}

.icon-star::after  { content: '★'; font-size: 20px; }
.icon-grid::after  { content: '▦'; font-size: 18px; }
.icon-growth::after { content: '↗'; font-size: 20px; }
.icon-folder::after { content: '📁'; font-size: 18px; }

.insight-content {
  flex: 1;
}

.insight-title {
  font-size: 0.85rem;
  color: var(--color-text-muted);
  letter-spacing: 1px;
  margin-bottom: 0.25rem;
}

.insight-value {
  font-size: 1.5rem;
  font-weight: 700;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 0.25rem;
}

.insight-desc {
  font-size: 0.8rem;
  color: var(--color-text-secondary);
}
</style>
