<template>
  <div class="cal-month-view card">
    <div class="cal-header">
      <span class="cal-title">{{ monthYearLabel }}</span>
    </div>
    <div v-if="monthGrid.length === 0" class="cal-empty">{{ t('analytics.noData') }}</div>
    <template v-else>
      <div class="cal-grid-wrapper">
        <table class="cal-grid">
          <thead>
            <tr>
              <th v-for="d in dayNames" :key="d" class="grid-day-header"
                :class="{ 'weekend-header': d === dayNames[5] || d === dayNames[6] }">{{ d }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(week, wi) in monthGrid" :key="wi">
              <td v-for="(cell, ci) in week" :key="ci"
                class="cal-cell"
                :class="{
                  'is-today': cell?.isToday,
                  'has-data': !!cell?.data,
                  'weekend-cell': ci >= 5
                }"
                :style="cell ? { backgroundColor: cellBg(cell.data?.commits || 0, monthMaxCommits) } : {}"
                @click="cell && selectDate(cell)"
              >
                <div v-if="cell" class="cell-inner">
                  <div class="cell-day" :class="{ 'today-text': cell.isToday }">{{ cell.day }}</div>
                  <div v-if="cell.data" class="cell-count">{{ cell.data.commits }}</div>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="selectedCell" class="detail-panel expand-panel">
        <div class="detail-hero">
          <div class="hero-date">{{ formattedDate }}</div>
          <div class="hero-stats">
            <div class="hero-commits">
              <span class="hero-number">{{ selectedCell.data?.commits || 0 }}</span>
              <span class="hero-label">{{ t('calendar.commitsUnit') }}</span>
            </div>
            <div v-if="selectedCell.data" class="hero-changes">
              <span class="hero-add">+{{ selectedCell.data.additions }}</span>
              <span class="hero-del">-{{ selectedCell.data.deletions }}</span>
              <span class="hero-net" :class="netChange >= 0 ? 'pos' : 'neg'">
                {{ netChange >= 0 ? '+' : '' }}{{ netChange }}
              </span>
            </div>
          </div>
        </div>
        <div class="detail-divider"></div>
        <table v-if="dateDetail && dateDetail.length > 0" class="detail-table">
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
            <tr v-for="(item, idx) in dateDetail" :key="item.repoName + '-' + item.author" :style="{ animationDelay: idx * 0.04 + 's' }">
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
    </template>
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

const selectedCell = ref(null)

function pad(n) {
  return String(n).padStart(2, '0')
}

function cellBg(commits, maxC) {
  if (!commits || !maxC) return 'transparent'
  const intensity = Math.min(commits / maxC, 1)
  const alpha = 0.05 + intensity * 0.55
  return `rgba(0, 212, 255, ${alpha.toFixed(2)})`
}

function isToday(dateStr) {
  const d = new Date()
  const today = `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}`
  return dateStr === today
}

const dayNames = computed(() => t('calendar.dayNamesShort'))
const monthNames = computed(() => t('calendar.monthNames'))

const monthYearLabel = computed(() => {
  if (!props.startDate) return ''
  const d = new Date(props.startDate + 'T00:00:00')
  const y = d.getFullYear()
  const m = d.getMonth() + 1
  return t('calendar.dateFormat')
    .replace('{y}', y).replace('{m}', m).replace('{d}', '')
    .replace(/\s*周?\{\w\}.*$/, '').replace(/日$/, '').trim()
    .replace(/,$/, '') || `${monthNames.value[m-1]} ${y}`
})

const monthGrid = computed(() => {
  if (!props.startDate) return []
  const d = new Date(props.startDate + 'T00:00:00')
  const year = d.getFullYear()
  const month = d.getMonth()

  const dateMap = new Map()
  for (const repo of props.dailyStats) {
    for (const author of repo.authors) {
      for (const day of author.dailyData) {
        if (!dateMap.has(day.date)) {
          dateMap.set(day.date, { commits: 0, additions: 0, deletions: 0 })
        }
        const data = dateMap.get(day.date)
        data.commits += day.commits
        data.additions += day.additions
        data.deletions += day.deletions
      }
    }
  }

  const firstDay = new Date(year, month, 1)
  const daysInMonth = new Date(year, month + 1, 0).getDate()
  let startCol = firstDay.getDay() - 1
  if (startCol < 0) startCol = 6

  const grid = []
  let week = []
  for (let i = 0; i < startCol; i++) week.push(null)
  for (let day = 1; day <= daysInMonth; day++) {
    const dateStr = `${year}-${pad(month+1)}-${pad(day)}`
    const data = dateMap.get(dateStr) || null
    week.push({ day, dateStr, data, isToday: isToday(dateStr) })
    if (week.length === 7) {
      grid.push(week)
      week = []
    }
  }
  if (week.length > 0) {
    while (week.length < 7) week.push(null)
    grid.push(week)
  }
  return grid
})

const monthMaxCommits = computed(() => {
  let m = 0
  for (const week of monthGrid.value) {
    for (const cell of week) {
      if (cell?.data?.commits > m) m = cell.data.commits
    }
  }
  return m
})

function selectDate(cell) {
  selectedCell.value = selectedCell.value?.dateStr === cell.dateStr ? null : cell
}

const dateDetail = computed(() => {
  if (!selectedCell.value) return null
  const dateStr = selectedCell.value.dateStr
  const items = []
  for (const repo of props.dailyStats) {
    for (const author of repo.authors) {
      for (const day of author.dailyData) {
        if (day.date === dateStr && day.commits > 0) {
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

const formattedDate = computed(() => {
  if (!selectedCell.value) return ''
  const d = new Date(selectedCell.value.dateStr + 'T00:00:00')
  const wd = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat']
  const wdZh = ['日','一','二','三','四','五','六']
  const ms = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec']
  const fmt = t('calendar.dateFormat')
  return fmt
    .replace('{y}', d.getFullYear())
    .replace('{m}', d.getMonth() + 1)
    .replace('{d}', d.getDate())
    .replace('{w}', locale.value === 'zh' ? wdZh[d.getDay()] : wd[d.getDay()])
})

const netChange = computed(() => {
  if (!selectedCell.value?.data) return 0
  return selectedCell.value.data.additions - selectedCell.value.data.deletions
})
</script>

<style scoped>
.cal-month-view.card {
  position: relative;
}
.cal-month-view.card::after {
  content: '';
  position: absolute;
  top: 0;
  left: 12%;
  right: 12%;
  height: 2px;
  background: linear-gradient(90deg, transparent, #00f5ff, #ff00ff, transparent);
  border-radius: 2px;
  opacity: 0.45;
  pointer-events: none;
}

.cal-header {
  margin-bottom: 1.25rem;
}
.cal-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 1rem;
  color: #00f5ff;
  letter-spacing: 2px;
}

.cal-empty {
  text-align: center;
  padding: 3rem 0;
  color: #64748b;
  font-size: 0.95rem;
  letter-spacing: 1px;
}

.cal-grid-wrapper {
  overflow-x: auto;
}
.cal-grid {
  width: 100%;
  border-collapse: collapse;
  table-layout: fixed;
}
.cal-grid thead th {
  text-align: center;
  padding: 0.5rem 0.3rem;
  font-family: 'Orbitron', sans-serif;
  font-size: 0.75rem;
  color: #94a3b8;
  letter-spacing: 1px;
  text-transform: uppercase;
  border-bottom: 1px solid rgba(160, 100, 200, 0.15);
}
.weekend-header {
  color: #64748b;
  opacity: 0.6;
}

.cal-cell {
  position: relative;
  text-align: center;
  vertical-align: top;
  padding: 4px;
  border: 1px solid rgba(148, 163, 184, 0.06);
  transition: all 0.2s ease;
  cursor: default;
  height: 72px;
  width: 14.28%;
}
.cal-cell.has-data {
  cursor: pointer;
}
.cal-cell:hover {
  border-color: rgba(160, 100, 200, 0.4);
  filter: brightness(1.25);
  z-index: 1;
}
.weekend-cell {
  background: rgba(148, 163, 184, 0.03);
}

.cell-inner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 2px;
}
.cell-day {
  font-family: 'Rajdhani', sans-serif;
  font-size: 0.85rem;
  font-weight: 600;
  color: #94a3b8;
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}
.cell-day.today-text {
  background: rgba(160, 100, 200, 0.8);
  color: #0a0e27;
  font-weight: 700;
  box-shadow: 0 0 12px rgba(160, 100, 200, 0.5);
}
.cell-count {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.85rem;
  font-weight: 700;
  color: #e2e8f0;
  line-height: 1;
}

@media (max-width: 600px) {
  .cal-cell {
    height: 60px;
    padding: 2px;
  }
  .cell-day {
    width: 22px;
    height: 22px;
    font-size: 0.75rem;
  }
  .cell-count {
    font-size: 0.7rem;
  }
}
</style>
