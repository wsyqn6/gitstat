<template>
  <div class="charts-section" v-if="hasData || loading">
    <h3 class="section-title">{{ t('repo.commitCalendar') }}</h3>
    <div v-if="loading" class="chart-placeholder">
      <span class="spinner"></span>
    </div>
    <div v-else-if="calYears.length > 0" class="cal-container">
      <div v-for="(yr, yi) in calYears" :key="yr.year" class="cal-year">
        <div class="cal-months">
          <span v-for="m in yr.monthLabels" :key="m.label"
                :style="{ width: m.span * (cellSize + cellGap) + 'px' }"
                class="cal-month-label">{{ m.label }}</span>
        </div>
        <div class="cal-body">
          <div class="cal-days">
            <span v-for="d in dayLabels" :key="d" class="cal-day-label">{{ d }}</span>
          </div>
          <div class="cal-grid"
               :style="{ gridTemplateColumns: `repeat(53, ${cellSize}px)`, gridTemplateRows: `repeat(7, ${cellSize}px)`, gap: cellGap + 'px' }">
            <div v-for="cell in yr.cells" :key="cell.date"
                 class="cal-cell"
                 :style="{ backgroundColor: cell.outOfYear ? 'transparent' : cell.color }"
                 :title="cell.outOfYear ? '' : cell.date + ' · ' + cell.count + ' ' + commitsLabel">
            </div>
          </div>
        </div>
        <div v-if="yi === calYears.length - 1" class="cal-footer">
          <span class="cal-legend-label">{{ lessLabel }}</span>
          <span v-for="l in 5" :key="l" class="cal-legend-cell"
                :style="{ backgroundColor: getCalColor(l - 1) }"></span>
          <span class="cal-legend-label">{{ moreLabel }}</span>
        </div>
      </div>
    </div>
    <p v-else class="empty-hint">{{ t('repo.noChartData') }}</p>

    <div class="chart-row">
      <ChartContainer
        :title="t('repo.cumulativeGrowth')"
        :option="cumulativeOption"
        :loading="loading"
        class="chart-card"
      />
      <ChartContainer
        :title="t('repo.hourlyDistribution')"
        :option="hourlyOption"
        :loading="loading"
        class="chart-card"
      />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from '../i18n'
import echarts from '../utils/echarts'
import ChartContainer from './ChartContainer.vue'

const { t, locale } = useI18n()

const props = defineProps({
  data: Object,
  loading: Boolean
})

const cellSize = 11
const cellGap = 3

const hasData = computed(() =>
  props.data?.calendar?.length > 0 || props.data?.hourly?.length > 0
)

const commitsLabel = computed(() =>
  locale.value === 'zh' ? '次提交' : 'commits'
)
const lessLabel = computed(() =>
  locale.value === 'zh' ? '少' : 'Less'
)
const moreLabel = computed(() =>
  locale.value === 'zh' ? '多' : 'More'
)

const dayLabels = computed(() => {
  if (locale.value === 'zh') return ['', '一', '', '三', '', '五', '']
  return ['', 'Mon', '', 'Wed', '', 'Fri', '']
})

const monthNames = computed(() => {
  if (locale.value === 'zh') {
    return ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']
  }
  return ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
})

const CAL_LEVELS = [
  'rgba(0, 212, 255, 0.04)',
  'rgba(0, 212, 255, 0.18)',
  'rgba(0, 212, 255, 0.38)',
  'rgba(0, 140, 255, 0.6)',
  'rgba(0, 100, 255, 0.85)'
]

function getCalColor(count) {
  if (count === 0) return CAL_LEVELS[0]
  if (count <= 2) return CAL_LEVELS[1]
  if (count <= 5) return CAL_LEVELS[2]
  if (count <= 10) return CAL_LEVELS[3]
  return CAL_LEVELS[4]
}

const calYears = computed(() => {
  if (!props.data?.calendar?.length) return []
  const dayMap = new Map()
  for (const d of props.data.calendar) {
    dayMap.set(d.date, d.count)
  }

  const years = new Set()
  for (const d of props.data.calendar) {
    years.add(d.date.slice(0, 4))
  }

  return Array.from(years).sort().map(yearStr => {
    const year = parseInt(yearStr)
    const jan1 = new Date(year, 0, 1)
    const dayOfWeek = jan1.getDay()
    const mondayOffset = dayOfWeek === 0 ? -6 : 1 - dayOfWeek
    const monday = new Date(jan1)
    monday.setDate(jan1.getDate() + mondayOffset)

    const cells = []
    for (let w = 0; w < 53; w++) {
      for (let d = 0; d < 7; d++) {
        const date = new Date(monday)
        date.setDate(monday.getDate() + w * 7 + d)
        const pad = (n) => String(n).padStart(2, '0')
        const dateStr = date.getFullYear() + '-' + pad(date.getMonth() + 1) + '-' + pad(date.getDate())
        const count = dayMap.get(dateStr) || 0
        const outOfYear = date.getFullYear() !== year
        cells.push({
          date: dateStr,
          count,
          outOfYear,
          color: getCalColor(count)
        })
      }
    }

    const monthLabels = []
    let currentMonth = -1
    for (let w = 0; w < 53; w++) {
      const d = new Date(monday)
      d.setDate(monday.getDate() + w * 7 + 3)
      const m = d.getMonth()
      if (m !== currentMonth) {
        currentMonth = m
        monthLabels.push({ label: monthNames.value[m], week: w })
      }
    }
    for (let i = 0; i < monthLabels.length; i++) {
      const nextWeek = i + 1 < monthLabels.length ? monthLabels[i + 1].week : 53
      monthLabels[i].span = nextWeek - monthLabels[i].week
    }

    return { year: yearStr, cells, monthLabels }
  })
})

function emptyOption() {
  return {
    title: {
      text: t('repo.noChartData'),
      left: 'center',
      top: 'center',
      textStyle: { color: '#64748b', fontSize: 14 }
    }
  }
}

const cumulativeOption = computed(() => {
  const data = props.data?.cumulative
  if (!data?.length) return emptyOption()
  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(10, 14, 39, 0.95)',
      borderColor: '#00f5ff',
      textStyle: { color: '#fff' },
      formatter: (p) => `${p[0].axisValue}<br/>${t('repo.commits')}: ${p[0].value}`
    },
    grid: { left: '3%', right: '3%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      data: data.map(d => d.date),
      boundaryGap: false,
      axisLine: { lineStyle: { color: '#334155' } },
      axisLabel: { color: '#94a3b8', fontSize: 10 },
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisLabel: { color: '#94a3b8' },
      splitLine: { lineStyle: { color: '#1e293b', type: 'dashed' } }
    },
    series: [{
      type: 'line',
      smooth: true,
      symbol: 'none',
      lineStyle: { width: 2, color: '#00d4ff' },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(0, 212, 255, 0.3)' },
          { offset: 1, color: 'rgba(0, 212, 255, 0.02)' }
        ])
      },
      data: data.map(d => d.total)
    }]
  }
})

const hourlyOption = computed(() => {
  const data = props.data?.hourly
  if (!data?.length) return emptyOption()
  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(10, 14, 39, 0.95)',
      borderColor: '#00f5ff',
      textStyle: { color: '#fff' },
      formatter: (p) => `${p[0].name}:00<br/>${t('repo.commits')}: ${p[0].value}`
    },
    grid: { left: '3%', right: '3%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      data: data.map(d => String(d.hour)),
      axisLine: { lineStyle: { color: '#334155' } },
      axisLabel: { color: '#94a3b8', fontSize: 10 },
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisLabel: { color: '#94a3b8' },
      splitLine: { lineStyle: { color: '#1e293b', type: 'dashed' } }
    },
    series: [{
      type: 'bar',
      barWidth: '60%',
      data: data.map(d => ({
        value: d.count,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#00d4ff' },
            { offset: 1, color: 'rgba(0, 212, 255, 0.15)' }
          ]),
          borderRadius: [3, 3, 0, 0]
        }
      }))
    }]
  }
})
</script>

<style scoped>
.charts-section {
  margin-top: 2rem;
}
.section-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.95rem;
  color: #00d4ff;
  letter-spacing: 1px;
  margin: 0 0 1rem 0;
}

/* Loading */
.chart-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 120px;
}

/* Calendar */
.cal-container {
  margin-bottom: 2rem;
}
.cal-year {
  margin-bottom: 0.5rem;
}
.cal-year:last-child {
  margin-bottom: 0;
}
.cal-months {
  display: flex;
  padding-left: 32px;
  margin-bottom: 2px;
}
.cal-month-label {
  font-size: 0.65rem;
  color: #64748b;
  font-family: 'Rajdhani', sans-serif;
  text-align: left;
  white-space: nowrap;
}
.cal-body {
  display: flex;
  gap: 4px;
}
.cal-days {
  display: flex;
  flex-direction: column;
  gap: 3px;
  padding-top: 0;
}
.cal-day-label {
  height: 11px;
  line-height: 11px;
  font-size: 0.6rem;
  color: #64748b;
  font-family: 'Rajdhani', sans-serif;
  text-align: right;
  padding-right: 2px;
}
.cal-grid {
  display: grid;
}
.cal-cell {
  width: 11px;
  height: 11px;
  border-radius: 2px;
  transition: all 0.15s;
  cursor: default;
}
.cal-cell:hover {
  transform: scale(1.3);
  outline: 1px solid rgba(0, 212, 255, 0.6);
  outline-offset: 1px;
}
.cal-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 3px;
  margin-top: 6px;
  padding-right: 0;
}
.cal-legend-label {
  font-size: 0.6rem;
  color: #64748b;
  font-family: 'Rajdhani', sans-serif;
}
.cal-legend-cell {
  width: 10px;
  height: 10px;
  border-radius: 2px;
}

.empty-hint {
  color: #64748b;
  font-family: 'Rajdhani', sans-serif;
  font-size: 0.85rem;
  margin-bottom: 2rem;
}

/* Chart row */
.chart-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}
.chart-card {
  min-height: 300px;
}
.chart-card :deep(.chart) {
  height: 280px;
}
</style>
