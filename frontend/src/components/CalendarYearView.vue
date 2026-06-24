<template>
  <div class="cal-year-view glass card">
    <div class="cal-header">
      <span class="cal-title">{{ yearLabel }}</span>
    </div>
    <div v-if="yearData.length === 0" class="cal-empty">{{ t('analytics.noData') }}</div>
    <div v-else class="year-grid">
      <div v-for="month in yearData" :key="month.month"
        class="year-month-card"
        :class="{ 'future-month': month.isFuture, 'selected': selectedMonth?.key === month.key }"

        @click="!month.isFuture && (selectedMonth = selectedMonth?.key === month.key ? null : month)"
      >
        <div class="ym-header">{{ month.name }}</div>
        <div class="ym-commits">{{ month.commits }}</div>
        <div class="ym-changes">
          <span class="ym-add">+{{ month.additions }}</span>
          <span class="ym-del">-{{ month.deletions }}</span>
          <span class="ym-net">{{ (month.additions || 0) - (month.deletions || 0) >= 0 ? '+' : '' }}{{ (month.additions || 0) - (month.deletions || 0) }}</span>
        </div>
      </div>
    </div>
    <div v-if="selectedMonth" class="detail-panel expand-panel">
      <div class="detail-hero">
        <div class="hero-date">{{ selectedMonth.name }} {{ startDate?.slice(0, 4) }}</div>
        <div class="hero-stats">
          <div class="hero-commits">
            <span class="hero-number">{{ selectedMonth.commits }}</span>
            <span class="hero-label">{{ t('calendar.commitsUnit') }}</span>
          </div>
          <div class="hero-changes">
            <span class="hero-add">+{{ selectedMonth.additions }}</span>
            <span class="hero-del">-{{ selectedMonth.deletions }}</span>
            <span class="hero-net" :class="(selectedMonth.additions - selectedMonth.deletions) >= 0 ? 'pos' : 'neg'">
              {{ (selectedMonth.additions - selectedMonth.deletions) >= 0 ? '+' : '' }}{{ selectedMonth.additions - selectedMonth.deletions }}
            </span>
          </div>
        </div>
      </div>
      <div class="detail-divider"></div>
      <table v-if="monthDetail.length > 0" class="detail-table">
        <thead>
          <tr>
            <th>{{ t('analytics.repoName') }}</th>
            <th>{{ t('analytics.developer') }}</th>
            <th class="num-col">{{ t('dashboard.commits') }}</th>
            <th class="num-col">{{ t('analytics.additions') }}</th>
            <th class="num-col">{{ t('analytics.deletions') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, idx) in monthDetail" :key="item.repoName + '-' + item.author" :style="{ animationDelay: idx * 0.04 + 's' }">
            <td class="cell-repo">{{ item.repoName }}</td>
            <td class="cell-author">{{ item.author }}</td>
            <td class="num-col">{{ item.commits }}</td>
            <td class="num-col cell-additions">+{{ item.additions }}</td>
            <td class="num-col cell-deletions">-{{ item.deletions }}</td>
          </tr>
        </tbody>
      </table>
      <div v-else class="detail-empty">
        <span class="empty-icon">⌧</span>
        {{ t('calendar.noDetail') }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '../i18n'

const props = defineProps({
  dailyStats: { type: Array, default: () => [] },
  startDate: { type: String, default: '' }
})

const { t, locale } = useI18n()

const selectedMonth = ref(null)

function pad(n) {
  return String(n).padStart(2, '0')
}

const monthNames = computed(() => t('calendar.monthNames'))

const yearLabel = computed(() => {
  if (!props.startDate) return ''
  return props.startDate.slice(0, 4) + (locale.value === 'zh' ? t('calendar.yearSuffix') : '')
})

const yearData = computed(() => {
  if (!props.startDate) return []
  const year = parseInt(props.startDate.slice(0, 4))
  const now = new Date()
  const currentMonth = now.getMonth() + 1
  const currentYear = now.getFullYear()

  const monthMap = new Map()
  for (let m = 1; m <= 12; m++) {
    const key = `${year}-${pad(m)}`
    monthMap.set(key, { month: m, key, commits: 0, additions: 0, deletions: 0 })
  }

  for (const repo of props.dailyStats) {
    for (const author of repo.authors) {
      for (const day of (author.dailyData || [])) {
        const month = day.date.slice(0, 7)
        if (monthMap.has(month)) {
          const data = monthMap.get(month)
          data.commits += day.commits
          data.additions += day.additions
          data.deletions += day.deletions
        }
      }
    }
  }

  let maxC = 0
  const result = []
  for (let m = 1; m <= 12; m++) {
    const key = `${year}-${pad(m)}`
    const data = monthMap.get(key)
    data.isFuture = (year > currentYear) || (year === currentYear && m > currentMonth)
    data.name = monthNames.value[m - 1]
    if (data.commits > maxC) maxC = data.commits
    result.push(data)
  }

  result.forEach(d => d.maxCommits = maxC)
  return result
})

const monthDetail = computed(() => {
  if (!selectedMonth.value) return []
  const key = selectedMonth.value.key
  const items = []
  for (const repo of props.dailyStats) {
    for (const author of repo.authors) {
      for (const day of (author.dailyData || [])) {
        const month = day.date.slice(0, 7)
        if (month === key && day.commits > 0) {
          items.push({
            repoName: repo.repoName,
            author: author.author,
            commits: day.commits,
            additions: day.additions,
            deletions: day.deletions
          })
        }
      }
    }
  }
  return items.sort((a, b) => b.commits - a.commits)
})
</script>

<style scoped>


.cal-header {
  margin-bottom: 1.25rem;
}
.cal-title {
  font-family: var(--font-display);
  font-size: 1rem;
  color: var(--color-accent);
  letter-spacing: 2px;
}

.cal-empty {
  text-align: center;
  padding: 3rem 0;
  color: var(--color-text-muted);
  font-size: 0.95rem;
  letter-spacing: 1px;
}

.year-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
}
.year-month-card {
  background: var(--bg-card);
  border: 1px solid var(--border-card);
  border-radius: var(--radius-card);
  padding: 1.5rem;
  text-align: center;
  cursor: pointer;
  content-visibility: auto;
  contain-intrinsic-size: 120px;
}
.year-month-card.selected {
  border-color: var(--border-card-hover);
}
.year-month-card.future-month {
  opacity: 0.3;
  cursor: not-allowed;
  filter: grayscale(0.5);
}
.ym-header {
  position: relative;
  z-index: 1;
  font-family: var(--font-display);
  font-size: 0.7rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  letter-spacing: 2px;
  text-transform: uppercase;
  margin-bottom: 0.35rem;
}
.ym-commits {
  position: relative;
  z-index: 1;
  font-family: var(--font-display);
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-input-text);
  line-height: 1.15;
  margin-bottom: 0.4rem;
}
.ym-changes {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  font-family: var(--font-body);
  font-size: 0.65rem;
  font-weight: 600;
}
.ym-changes .ym-add { color: var(--color-green); }
.ym-changes .ym-del { color: var(--color-red); }
.ym-changes .ym-net { color: var(--color-text-muted); }

@media (max-width: 900px) {
  .year-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
@media (max-width: 600px) {
  .year-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
