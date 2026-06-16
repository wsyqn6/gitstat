<template>
  <div class="cal-week-view card">
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
                <span class="author-name">{{ author.author }}</span>
                <span v-if="author.isMe" class="me-badge">{{ t('dashboard.me') }}</span>
              </div>
              <div class="author-email">{{ author.email }}</div>
            </td>
            <td v-for="day in weekDays" :key="day.date"
              class="data-cell"
              :class="{ 'has-data': !!author.days[day.date] }"
              :style="{ backgroundColor: cellBg(author.days[day.date]?.commits || 0, maxCommits) }"
            >
              <div class="cell-commits">{{ author.days[day.date]?.commits ?? '-' }}</div>
              <div v-if="author.days[day.date]" class="cell-changes">
                <span class="add">+{{ author.days[day.date].additions }}</span>
                <span class="del">-{{ author.days[day.date].deletions }}</span>
              </div>
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
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '../i18n'

const props = defineProps({
  dailyStats: { type: Array, default: () => [] },
  startDate: { type: String, default: '' }
})

const { t } = useI18n()

const expandedAuthor = ref(null)

function toggleExpandAuthor(email) {
  expandedAuthor.value = expandedAuthor.value === email ? null : email
}

function pad(n) {
  return String(n).padStart(2, '0')
}

function cellBg(commits, maxC) {
  if (!commits || !maxC) return 'transparent'
  const intensity = Math.min(commits / maxC, 1)
  const alpha = 0.05 + intensity * 0.55
  return `rgba(160, 100, 200, ${alpha.toFixed(2)})`
}

const dayNames = computed(() => t('calendar.dayNamesShort'))

const weekDays = computed(() => {
  if (!props.startDate) return []
  const start = new Date(props.startDate + 'T00:00:00')
  const dow = start.getDay() || 7
  const monday = new Date(start)
  monday.setDate(start.getDate() - dow + 1)
  const days = []
  for (let i = 0; i < 7; i++) {
    const d = new Date(monday)
    d.setDate(monday.getDate() + i)
    const y = d.getFullYear()
    const mo = String(d.getMonth() + 1).padStart(2, '0')
    const day = String(d.getDate()).padStart(2, '0')
    const dateStr = `${y}-${mo}-${day}`
    const name = dayNames.value[i]
    const dateShort = `${d.getMonth()+1}/${d.getDate()}`
    days.push({ date: dateStr, name, dateShort })
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
</script>

<style scoped>
.cal-week-view.card {
  position: relative;
}
.cal-week-view.card::after {
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
  border-bottom: 1px solid rgba(160, 100, 200, 0.15);
  vertical-align: bottom;
}
.cal-table .row-header {
  text-align: center;
  padding: 0.5rem 0.3rem;
  min-width: 64px;
  color: #64748b;
  font-family: 'Orbitron', sans-serif;
  font-size: 0.75rem;
  letter-spacing: 1px;
  text-transform: uppercase;
}
.col-header {
  min-width: 64px;
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
  min-width: 80px;
}
.total-col .day-name {
  color: #ffd700;
}

.data-cell {
  text-align: center;
  padding: 0.6rem 0.3rem;
  border-bottom: 1px solid rgba(148, 163, 184, 0.08);
  transition: all 0.2s ease;
  cursor: default;
}
.data-cell.has-data {
  cursor: pointer;
}
.data-cell:hover {
  filter: brightness(1.3);
  border-color: rgba(160, 100, 200, 0.3);
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
  border-left: 1px solid rgba(160, 100, 200, 0.25);
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
  background: rgba(160, 100, 200, 0.05);
}
.author-row.expanded {
  background: rgba(160, 100, 200, 0.08);
}
.cell-author {
  text-align: center;
  padding: 0.35rem 0.15rem !important;
  line-height: 1.2;
}
.author-label {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 3px;
}
.author-name {
  font-family: 'Rajdhani', sans-serif;
  color: #00f5ff;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: none;
  letter-spacing: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.author-email {
  font-family: 'Rajdhani', sans-serif;
  font-size: 0.6rem;
  color: #64748b;
  text-transform: none !important;
  letter-spacing: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.2;
  margin-top: 1px;
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

@media (max-width: 600px) {
  .cal-table .row-header {
    min-width: 50px;
    font-size: 0.6rem;
  }
  .th-dev {
    min-width: 65px;
  }
  .col-header {
    min-width: 44px;
  }
  .cell-author {
    padding: 0.2rem 0.2rem !important;
    max-width: 45px;
  }
  .author-name {
    font-size: 0.65rem;
    max-width: 38px;
  }
  .author-email {
    font-size: 0.45rem;
    max-width: 38px;
  }
}
</style>
