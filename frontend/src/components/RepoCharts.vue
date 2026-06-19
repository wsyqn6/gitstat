<template>
  <div class="charts-section">
    <div class="section-header">
      <h3 class="section-title">{{ t('repo.commitCalendar') }}</h3>
      <div v-if="availableYears.length > 1" class="year-tabs">
        <button v-for="y in availableYears" :key="y"
                :class="['year-tab btn-ghost', { active: y === selectedYear }]"
                @click="selectedYear = y">{{ y }}</button>
      </div>
    </div>
    <div v-if="loading" class="chart-placeholder">
      <span class="spinner"></span>
    </div>
    <div v-else-if="activeYear" class="cal-container">
      <div class="cal-months">
    <span v-for="m in activeYear.monthLabels" :key="m.label"
          class="cal-month-label" :style="{ gridColumn: 'span ' + m.span }">{{ m.label }}</span>
      </div>
      <div class="cal-body">
        <div class="cal-days">
          <span v-for="d in dayLabels" :key="d" class="cal-day-label">{{ d }}</span>
        </div>
        <div class="cal-grid">
          <div v-for="cell in activeYear.cells" :key="cell.date"
               class="cal-cell"
               :style="{ backgroundColor: cell.outOfYear ? 'transparent' : cell.color }"
               :title="cell.outOfYear ? '' : cell.date + ' · ' + cell.count + ' ' + commitsLabel">
          </div>
        </div>
      </div>
      <div class="cal-footer">
        <span class="cal-legend-label">{{ lessLabel }}</span>
        <span v-for="l in 5" :key="l" class="cal-legend-cell"
              :style="{ backgroundColor: getCalColor(l - 1) }"></span>
        <span class="cal-legend-label">{{ moreLabel }}</span>
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
import { ref, computed, watch } from 'vue'
import { useI18n } from '../i18n'
import { getChartConfig } from '../utils/constants'
import { useTheme } from '../composables/useTheme'
import echarts from '../utils/echarts'
import ChartContainer from './ChartContainer.vue'

const { t, locale } = useI18n()
const { theme } = useTheme()
const chartCfg = computed(() => getChartConfig(theme.value))

const props = defineProps({
  data: Object,
  loading: Boolean
})

const cellSize = 11
const cellGap = 3

const commitsLabel = computed(() => t('repo.charts.commitsUnit'))
const lessLabel = computed(() => t('repo.charts.less'))
const moreLabel = computed(() => t('repo.charts.more'))

const dayLabels = computed(() => t('calendar.dayNamesShort').map((d, i) => i % 2 === 1 ? d : ''))

const monthNames = computed(() => t('calendar.monthNames'))

const selectedYear = ref('')

function getCalColor(count) {
  const levels = chartCfg.value.calLevels
  if (count === 0) return levels[0]
  if (count <= 2) return levels[1]
  if (count <= 5) return levels[2]
  if (count <= 10) return levels[3]
  return levels[4]
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

  const pad = (n) => String(n).padStart(2, '0')
  const lastDate = props.data.calendar[props.data.calendar.length - 1]?.date || ''

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
        const dateStr = date.getFullYear() + '-' + pad(date.getMonth() + 1) + '-' + pad(date.getDate())
        const count = dayMap.get(dateStr) || 0
        const outOfYear = date.getFullYear() !== year || dateStr > lastDate
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

const availableYears = computed(() =>
  calYears.value.map(y => y.year)
)

const activeYear = computed(() =>
  calYears.value.find(y => y.year === selectedYear.value) || null
)

watch(() => props.data, () => {
  const years = availableYears.value
  if (years.length > 0 && !years.includes(selectedYear.value)) {
    selectedYear.value = years[years.length - 1]
  }
}, { immediate: true })


function emptyOption() {
  return {
    title: {
      text: t('repo.noChartData'),
      left: 'center',
      top: 'center',
      textStyle: { color: chartCfg.value.titleColor, fontSize: 14 }
    }
  }
}

const cumulativeOption = computed(() => {
  const data = props.data?.cumulative
  if (!data?.length) return emptyOption()
  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: chartCfg.value.tooltipBg,
      borderColor: chartCfg.value.accent,
      textStyle: { color: chartCfg.value.tooltipText },
      formatter: (p) => `${p[0].axisValue}<br/>${t('repo.commits')}: ${p[0].value}`
    },
    grid: { left: '3%', right: '3%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      data: data.map(d => d.date),
      boundaryGap: false,
      axisLine: { lineStyle: { color: chartCfg.value.axisLine } },
      axisLabel: { color: chartCfg.value.axisLabel, fontSize: 10 },
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisLabel: { color: chartCfg.value.axisLabel },
      splitLine: { lineStyle: { color: chartCfg.value.splitLine, type: 'dashed' } }
    },
    series: [{
      type: 'line',
      smooth: true,
      symbol: 'none',
      lineStyle: { width: 2, color: chartCfg.value.primary },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: chartCfg.value.primary + '4D' },
          { offset: 1, color: chartCfg.value.primary + '05' }
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
      backgroundColor: chartCfg.value.tooltipBg,
      borderColor: chartCfg.value.accent,
      textStyle: { color: chartCfg.value.tooltipText },
      formatter: (p) => `${p[0].name}:00<br/>${t('repo.commits')}: ${p[0].value}`
    },
    grid: { left: '3%', right: '3%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      data: data.map(d => String(d.hour)),
      axisLine: { lineStyle: { color: chartCfg.value.axisLine } },
      axisLabel: { color: chartCfg.value.axisLabel, fontSize: 10 },
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisLabel: { color: chartCfg.value.axisLabel },
      splitLine: { lineStyle: { color: chartCfg.value.splitLine, type: 'dashed' } }
    },
    series: [{
      type: 'bar',
      barWidth: '60%',
      data: data.map(d => ({
        value: d.count,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: chartCfg.value.primary },
            { offset: 1, color: chartCfg.value.primary + '26' }
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
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}
.section-title {
  font-family: var(--font-display);
  font-size: 0.95rem;
  color: var(--color-primary);
  letter-spacing: 1px;
  margin: 0;
}
.year-tabs {
  display: flex;
  gap: 0.25rem;
}
.year-tab {
  padding: 0.2rem 0.6rem;
  border: 1px solid var(--glass-outline-border);
  border-radius: 4px;
  font-size: 0.7rem;
}
.year-tab.active {
  background: var(--glass-btn-bg);
  backdrop-filter: blur(var(--glass-blur));
  border-color: var(--color-accent);
  color: var(--glass-btn-color);
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
.cal-months {
  display: grid;
  grid-template-columns: repeat(53, 1fr);
  gap: 2px;
  padding-left: 32px;
  margin-bottom: 2px;
}
.cal-month-label {
  font-size: 0.65rem;
  color: var(--color-text-muted);
  font-family: var(--font-body);
  text-align: left;
  white-space: nowrap;
}
.cal-body {
  display: flex;
  gap: 4px;
}
.cal-days {
  display: grid;
  grid-template-rows: repeat(7, 1fr);
  gap: 2px;
  padding-top: 0;
}
.cal-day-label {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  font-size: 0.6rem;
  color: var(--color-text-muted);
  font-family: var(--font-body);
  padding-right: 2px;
}
.cal-grid {
  display: grid;
  grid-template-rows: repeat(7, 1fr);
  grid-auto-flow: column;
  gap: 2px;
  width: 100%;
}
.cal-cell {
  aspect-ratio: 1;
  border-radius: 2px;
  transition: all 0.15s;
  cursor: default;
  min-width: 0;
}
.cal-cell:hover {
  transform: scale(1.3);
  outline: var(--outline-primary);
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
  color: var(--color-text-muted);
  font-family: var(--font-body);
}
.cal-legend-cell {
  width: 10px;
  height: 10px;
  border-radius: 2px;
}

.empty-hint {
  color: var(--color-text-muted);
  font-family: var(--font-body);
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
