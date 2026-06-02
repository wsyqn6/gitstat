<template>
  <div class="calendar-area">
    <!-- ====== WEEK VIEW ====== -->
    <div v-if="viewType === 'week'" class="cal-week-view card">
      <div class="cal-header">
        <span class="cal-title">{{ weekRangeLabel }}</span>
      </div>
      <div v-if="weekData.length === 0" class="cal-empty">{{ t('analytics.noData') }}</div>
      <div v-else class="cal-table-wrapper">
        <table class="cal-table">
          <thead>
            <tr>
              <th class="row-header th-dev">{{ t('analytics.developer') }}</th>
              <th v-for="day in weekDays" :key="day.date" class="col-header">
                <div class="day-name">{{ day.name }}</div>
                <div class="day-date">{{ day.dateShort }}</div>
              </th>
              <th class="col-header total-col">{{ t('calendar.total') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="author in weekData" :key="author.email"
              @click="toggleExpandAuthor(author.email)"
              class="author-row"
              :class="{ expanded: expandedAuthor === author.email }"
            >
              <td class="row-header cell-author">
                <div class="author-label">
                  <span v-if="author.isMe" class="me-badge">{{ t('dashboard.me') }}</span>
                  <span>{{ author.author }}</span>
                </div>
              </td>
              <td v-for="day in weekDays" :key="day.date"
                class="data-cell"
                :class="{ 'has-data': !!author.days[day.date] }"
                :style="{ backgroundColor: cellBg(author.days[day.date]?.commits || 0, maxCommits) }"
              >
                <div class="cell-commits">{{ author.days[day.date]?.commits ?? '-' }}</div>
              </td>
              <td class="data-cell total-cell">
                <div class="cell-commits">{{ author.total.commits }}</div>
                <div class="cell-changes">
                  <span class="add">+{{ author.total.additions }}</span>
                  <span class="del">-{{ author.total.deletions }}</span>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ====== MONTH VIEW ====== -->
    <div v-if="viewType === 'month'" class="cal-month-view card">
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
          <div class="detail-header">
            {{ selectedCell.dateStr }} · {{ selectedCell.data?.commits || 0 }} {{ t('calendar.commitsUnit') }}
            <span v-if="selectedCell.data"> · +{{ selectedCell.data.additions }}/-{{ selectedCell.data.deletions }}</span>
          </div>
          <table v-if="dateDetail && dateDetail.length > 0" class="expand-table">
            <thead>
              <tr>
                <th>{{ t('analytics.repoName') }}</th>
                <th>{{ t('analytics.developer') }}</th>
                <th>{{ t('dashboard.commits') }}</th>
                <th>{{ t('analytics.additions') }}</th>
                <th>{{ t('analytics.deletions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, idx) in dateDetail" :key="idx">
                <td class="cell-author">{{ item.repoName }}</td>
                <td>{{ item.author }}</td>
                <td>{{ item.commits }}</td>
                <td class="cell-additions">+{{ item.additions }}</td>
                <td class="cell-deletions">-{{ item.deletions }}</td>
              </tr>
            </tbody>
          </table>
          <div v-else class="detail-empty">{{ t('calendar.noDetail') }}</div>
        </div>
      </template>
    </div>

    <!-- ====== YEAR VIEW ====== -->
    <div v-if="viewType === 'year'" class="cal-year-view card">
      <div class="cal-header">
        <span class="cal-title">{{ yearLabel }}</span>
      </div>
      <div v-if="yearData.length === 0" class="cal-empty">{{ t('analytics.noData') }}</div>
      <div v-else class="year-grid">
        <div v-for="month in yearData" :key="month.month"
          class="year-month-card"
          :class="{ 'future-month': month.isFuture }"
          @click="!month.isFuture && $emit('switch-to-month', month.month)"
        >
          <div class="ym-header">{{ month.name }}</div>
          <div class="ym-commits">{{ month.commits }}</div>
          <div class="ym-bar-track">
            <div class="ym-bar" :style="{ width: barPct(month.commits, yearMaxCommits) }"></div>
          </div>
          <div class="ym-changes">
            <span class="add">+{{ month.additions }}</span>
            <span class="del">-{{ month.deletions }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '../i18n'

const props = defineProps({
  viewType: { type: String, required: true },
  dailyStats: { type: Array, default: () => [] },
  periodStats: { type: Array, default: () => [] },
  startDate: { type: String, default: '' },
  endDate: { type: String, default: '' }
})

defineEmits(['switch-to-month'])

const { t, locale } = useI18n()

const expandedAuthor = ref(null)
const selectedCell = ref(null)

function toggleExpandAuthor(email) {
  expandedAuthor.value = expandedAuthor.value === email ? null : email
}

// ====== Helpers ======

const dayNames = computed(() => {
  return locale.value === 'zh'
    ? ['一', '二', '三', '四', '五', '六', '日']
    : ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
})

const monthNames = computed(() => {
  if (locale.value === 'zh') {
    return ['一月','二月','三月','四月','五月','六月','七月','八月','九月','十月','十一月','十二月']
  }
  return ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec']
})

function pad(n) {
  return String(n).padStart(2, '0')
}

function cellBg(commits, maxC) {
  if (!commits || !maxC) return 'transparent'
  const intensity = Math.min(commits / maxC, 1)
  const alpha = 0.05 + intensity * 0.55
  return `rgba(0, 245, 255, ${alpha.toFixed(2)})`
}

function barPct(val, maxV) {
  if (!val || !maxV) return '0%'
  return Math.round((val / maxV) * 100) + '%'
}

function isToday(dateStr) {
  const d = new Date()
  const today = `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}`
  return dateStr === today
}

// ====== WEEK VIEW ======

const weekDays = computed(() => {
  if (!props.startDate || !props.endDate) return []
  const start = new Date(props.startDate + 'T00:00:00')
  const end = new Date(props.endDate + 'T00:00:00')
  const days = []
  const d = new Date(start)
  while (d <= end) {
    const dateStr = d.toISOString().split('T')[0]
    const dow = d.getDay() || 7
    const name = dayNames.value[dow - 1]
    const dateShort = `${d.getMonth()+1}/${d.getDate()}`
    days.push({ date: dateStr, name, dateShort })
    d.setDate(d.getDate() + 1)
  }
  return days
})

const weekRangeLabel = computed(() => {
  if (weekDays.value.length === 0) return ''
  const first = weekDays.value[0]
  const last = weekDays.value[weekDays.value.length - 1]
  return `${first.date} ~ ${last.date}`
})

const weekData = computed(() => {
  const authorMap = new Map()
  const daySet = new Set(weekDays.value.map(d => d.date))

  for (const repo of props.dailyStats) {
    for (const author of repo.authors) {
      if (!authorMap.has(author.email)) {
        authorMap.set(author.email, {
          author: author.author,
          email: author.email,
          isMe: author.isMe,
          days: {},
          total: { commits: 0, additions: 0, deletions: 0 }
        })
      }
      const entry = authorMap.get(author.email)
      for (const day of author.dailyData) {
        if (daySet.has(day.date)) {
          if (!entry.days[day.date]) {
            entry.days[day.date] = { commits: 0, additions: 0, deletions: 0 }
          }
          entry.days[day.date].commits += day.commits
          entry.days[day.date].additions += day.additions
          entry.days[day.date].deletions += day.deletions
          entry.total.commits += day.commits
          entry.total.additions += day.additions
          entry.total.deletions += day.deletions
        }
      }
    }
  }
  return Array.from(authorMap.values()).sort((a, b) => b.total.commits - a.total.commits)
})

const maxCommits = computed(() => {
  let m = 0
  for (const author of weekData.value) {
    for (const day of Object.values(author.days)) {
      if (day.commits > m) m = day.commits
    }
  }
  return m
})

// ====== MONTH VIEW ======

const monthYearLabel = computed(() => {
  if (!props.startDate) return ''
  const d = new Date(props.startDate + 'T00:00:00')
  const y = d.getFullYear()
  const m = d.getMonth() + 1
  return locale.value === 'zh' ? `${y}年${m}月` : `${monthNames.value[m-1]} ${y}`
})

const monthGrid = computed(() => {
  if (!props.startDate) return []
  const d = new Date(props.startDate + 'T00:00:00')
  const year = d.getFullYear()
  const month = d.getMonth()

  // Build date→data map
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

// ====== YEAR VIEW ======

const yearLabel = computed(() => {
  if (!props.startDate) return ''
  return props.startDate.slice(0, 4) + (locale.value === 'zh' ? '年' : '')
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

  for (const repo of props.periodStats) {
    for (const author of repo.authors) {
      for (const period of author.dailyData) {
        if (monthMap.has(period.date)) {
          const data = monthMap.get(period.date)
          data.commits += period.commits
          data.additions += period.additions
          data.deletions += period.deletions
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

const yearMaxCommits = computed(() => {
  let m = 0
  for (const month of yearData.value) {
    if (month.commits > m) m = month.commits
  }
  return m
})
</script>

<style scoped>
.calendar-area {
  animation: slideDown 0.3s ease;
}

/* ====== Shared Header ====== */
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

/* ====== Week View Table ====== */
.cal-table-wrapper {
  overflow-x: auto;
}
.cal-table {
  width: 100%;
  border-collapse: collapse;
  font-family: 'Rajdhani', sans-serif;
  font-size: 0.9rem;
}
.cal-table thead th {
  text-align: center;
  padding: 0.5rem 0.4rem;
  border-bottom: 1px solid rgba(0, 245, 255, 0.15);
  vertical-align: bottom;
}
.cal-table .row-header {
  text-align: left;
  padding: 0.5rem 0.75rem;
  min-width: 100px;
  color: #64748b;
  font-family: 'Orbitron', sans-serif;
  font-size: 0.75rem;
  letter-spacing: 1px;
  text-transform: uppercase;
  position: sticky;
  left: 0;
  background: rgba(10, 14, 39, 0.95);
  z-index: 2;
}
.th-dev {
  min-width: 120px;
}
.col-header {
  min-width: 72px;
}
.day-name {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.72rem;
  color: #00f5ff;
  letter-spacing: 1px;
  text-transform: uppercase;
}
.day-date {
  font-size: 0.75rem;
  color: #64748b;
  margin-top: 2px;
}
.total-col {
  min-width: 90px;
}
.total-col .day-name {
  color: #ffd700;
}

/* table body */
.data-cell {
  text-align: center;
  padding: 0.6rem 0.3rem;
  border-bottom: 1px solid rgba(148, 163, 184, 0.08);
  transition: background-color 0.2s ease;
  cursor: default;
}
.data-cell.has-data {
  cursor: pointer;
}
.data-cell:hover {
  filter: brightness(1.3);
}
.cell-commits {
  font-size: 1.05rem;
  font-weight: 700;
  color: #e2e8f0;
}
.cell-changes {
  font-size: 0.65rem;
  margin-top: 2px;
  display: flex;
  gap: 4px;
  justify-content: center;
}
.cell-changes .add { color: #00ff88; }
.cell-changes .del { color: #ff6b6b; }

.total-cell {
  border-left: 1px solid rgba(0, 245, 255, 0.2);
}
.total-cell .cell-commits {
  color: #ffd700;
  font-size: 1.15rem;
}

.author-row {
  cursor: pointer;
  transition: background 0.2s;
}
.author-row:hover {
  background: rgba(0, 245, 255, 0.03);
}
.author-row.expanded {
  background: rgba(0, 245, 255, 0.06);
}
.cell-author {
  color: #00f5ff !important;
  font-weight: 600;
  font-size: 0.85rem;
}
.author-label {
  display: flex;
  align-items: center;
  gap: 6px;
}
.me-badge {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.6rem;
  padding: 1px 5px;
  border: 1px solid #ffd700;
  border-radius: 3px;
  color: #ffd700;
  letter-spacing: 0.5px;
}

/* ====== Month View Grid ====== */
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
  border-bottom: 1px solid rgba(0, 245, 255, 0.15);
}
.weekend-header {
  color: #64748b;
  opacity: 0.6;
}

.cal-cell {
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
  border-color: rgba(0, 245, 255, 0.3);
  filter: brightness(1.2);
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
  background: #00f5ff;
  color: #0a0e27;
  font-weight: 700;
  box-shadow: 0 0 10px rgba(0, 245, 255, 0.5);
}
.cell-count {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.85rem;
  font-weight: 700;
  color: #e2e8f0;
  line-height: 1;
}

/* detail panel */
.detail-panel {
  margin-top: 1rem;
  animation: slideDown 0.25s ease;
}
.detail-header {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.8rem;
  color: #00f5ff;
  letter-spacing: 1px;
  margin-bottom: 0.75rem;
  opacity: 0.9;
}
.detail-empty {
  color: #64748b;
  font-size: 0.85rem;
  text-align: center;
  padding: 1rem 0;
}

/* ====== Year View Grid ====== */
.year-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1rem;
}
.year-month-card {
  background: rgba(10, 14, 39, 0.4);
  border: 1px solid rgba(0, 212, 255, 0.15);
  border-radius: 10px;
  padding: 1rem;
  text-align: center;
  transition: all 0.3s ease;
  cursor: pointer;
}
.year-month-card:hover:not(.future-month) {
  border-color: rgba(0, 212, 255, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 8px 32px rgba(0, 212, 255, 0.15);
}
.year-month-card.future-month {
  opacity: 0.35;
  cursor: not-allowed;
}
.ym-header {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.75rem;
  color: #00f5ff;
  letter-spacing: 1px;
  text-transform: uppercase;
  margin-bottom: 0.5rem;
}
.ym-commits {
  font-family: 'Orbitron', sans-serif;
  font-size: 1.3rem;
  font-weight: 700;
  background: linear-gradient(135deg, #00d4ff, #7800ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 0.5rem;
}
.ym-bar-track {
  height: 5px;
  background: rgba(148, 163, 184, 0.15);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}
.ym-bar {
  height: 100%;
  background: linear-gradient(90deg, #00f5ff, #ff00ff);
  border-radius: 3px;
  transition: width 0.5s ease;
}
.ym-changes {
  font-size: 0.65rem;
  display: flex;
  gap: 6px;
  justify-content: center;
}
.ym-changes .add { color: #00ff88; }
.ym-changes .del { color: #ff6b6b; }

/* responsive */
@media (max-width: 900px) {
  .year-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
@media (max-width: 600px) {
  .year-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .cal-table .row-header {
    min-width: 80px;
    font-size: 0.65rem;
  }
  .col-header {
    min-width: 52px;
  }
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

@keyframes slideDown {
  from { opacity: 0; transform: translateY(-8px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
