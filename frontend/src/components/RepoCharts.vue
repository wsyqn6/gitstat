<template>
  <div class="charts-section">
    <div class="section-header">
      <h3 class="section-title">{{ t('repo.commitCalendar') }}</h3>
      <div v-if="availableYears.length > 1" class="year-tabs">
        <button v-for="y in availableYears" :key="y"
                :class="['year-tab btn', { active: y === selectedYear }]"
                @click="selectedYear = y">{{ y }}</button>
      </div>
    </div>
    <div v-if="loading" class="chart-placeholder">
      <span class="spinner"></span>
    </div>
    <div v-else-if="activeYear" ref="calContainer" class="cal-container" :style="{ ...calColors, '--cell-size': cellSize + 'px' }">
        <div class="cal-grid">
          <div v-for="item in gridItems" :key="item.key"
               :class="item.cls"
               :style="item.style"
               :title="item.title">{{ item.text }}</div>
        </div>
      <div class="cal-footer">
        <span class="cal-legend-label">{{ lessLabel }}</span>
        <span v-for="l in 5" :key="l" class="cal-legend-cell"
              :style="{ background: 'var(--cal-lvl' + (l - 1) + ')' }"></span>
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
        :title="t('repo.codeGrowth')"
        :option="codeGrowthOption"
        :loading="loading"
        class="chart-card"
      />
    </div>
    <div class="chart-row chart-row-full">
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
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
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

const commitsLabel = computed(() => t('repo.charts.commitsUnit'))
const lessLabel = computed(() => t('repo.charts.less'))
const moreLabel = computed(() => t('repo.charts.more'))

const dayLabels = computed(() => t('calendar.dayNamesShort').map((d, i) => i % 2 === 1 ? d : ''))

const monthNames = computed(() => t('calendar.monthNames'))

const selectedYear = ref('')

function getLevel(count) {
  if (count === 0) return 0
  if (count <= 2) return 1
  if (count <= 5) return 2
  if (count <= 10) return 3
  return 4
}

const calYears = computed(() => {
  if (!props.data?.dailyAgg?.length) return []
  const dayMap = new Map()
  for (const d of props.data.dailyAgg) {
    dayMap.set(d.date, d.commits)
  }

  const years = new Set()
  for (const d of props.data.dailyAgg) {
    years.add(d.date.slice(0, 4))
  }

  const pad = (n) => String(n).padStart(2, '0')
  const lastDate = props.data.dailyAgg[props.data.dailyAgg.length - 1]?.date || ''

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
        cells.push({
          date: dateStr,
          count,
          outOfYear: date.getFullYear() !== year || dateStr > lastDate,
          row: d + 2,
          col: w + 2
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

const calColors = computed(() => {
  const lv = chartCfg.value.calLevels
  return {
    '--cal-lvl0': lv[0],
    '--cal-lvl1': lv[1],
    '--cal-lvl2': lv[2],
    '--cal-lvl3': lv[3],
    '--cal-lvl4': lv[4],
  }
})

const gridItems = computed(() => {
  const y = activeYear.value
  if (!y) return []
  const dl = dayLabels.value
  const cl = commitsLabel.value
  const items = []
  for (const m of y.monthLabels) {
    items.push({
      key: 'm' + m.week,
      cls: 'cal-month',
      style: 'grid-row:1;grid-column:' + (m.week + 2) + '/span ' + m.span,
      text: m.label
    })
  }
  for (let i = 0; i < 7; i++) {
    if (dl[i]) {
      items.push({
        key: 'd' + i,
        cls: 'cal-day',
        style: 'grid-row:' + (i + 2) + ';grid-column:1',
        text: dl[i]
      })
    }
  }
  for (const c of y.cells) {
    const o = c.outOfYear
    items.push({
      key: c.date,
      cls: 'cal-cell' + (o ? ' cal-out' : ' cal-lvl' + getLevel(c.count)),
      style: 'grid-row:' + c.row + ';grid-column:' + c.col
    })
    if (!o) {
      items[items.length - 1].title = c.date + ' · ' + c.count + ' ' + cl
    }
  }
  return items
})

const cellSize = ref(12)
const calContainer = ref(null)

function updateCellSize() {
  const el = calContainer.value
  if (!el) return
  const w = el.clientWidth
  const fontPx = parseFloat(getComputedStyle(el).fontSize) || 14
  const avail = w - 2 * fontPx - 53 * 2
  cellSize.value = Math.max(6, Math.floor(avail / 53))
}

let ro = null
onMounted(() => {
  if (activeYear.value && calContainer.value) {
    updateCellSize()
    ro = new ResizeObserver(updateCellSize)
    ro.observe(calContainer.value)
  }
})
watch(activeYear, (y) => {
  ro?.disconnect(); ro = null
  if (!y || !calContainer.value) return
  updateCellSize()
  ro = new ResizeObserver(updateCellSize)
  ro.observe(calContainer.value)
}, { flush: 'post' })
onUnmounted(() => ro?.disconnect())

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
  const data = props.data?.dailyAgg
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
      smooth: data.length > 1,
      symbol: data.length > 1 ? 'none' : 'circle',
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

const codeGrowthOption = computed(() => {
  const data = props.data?.dailyAgg
  if (!data?.length) return emptyOption()
  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: chartCfg.value.tooltipBg,
      borderColor: chartCfg.value.accent,
      textStyle: { color: chartCfg.value.tooltipText },
      formatter: (p) => `${p[0].axisValue}<br/>${t('repo.charts.codeGrowthLabel')}: ${p[0].value}`
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
      smooth: data.length > 1,
      symbol: data.length > 1 ? 'none' : 'circle',
      lineStyle: { width: 2, color: chartCfg.value.chartColors[3] },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: chartCfg.value.chartColors[3] + '4D' },
          { offset: 1, color: chartCfg.value.chartColors[3] + '05' }
        ])
      },
      data: data.map(d => d.netLines)
    }]
  }
})

const hourlyOption = computed(() => {
  const data = props.data?.hourly
  if (!data?.some(d => d.count > 0)) return emptyOption()
  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: chartCfg.value.tooltipBg,
      borderColor: chartCfg.value.accent,
      textStyle: { color: chartCfg.value.tooltipText },
      formatter: (p) => `${p[0].name}<br/>${t('repo.commits')}: ${p[0].value}`
    },
    grid: { left: '3%', right: '3%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      data: data.map(d => String(d.hour).padStart(2, '0') + ':00'),
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
  border-radius: 4px;
  font-size: 0.7rem;
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
.cal-grid {
  display: grid;
  grid-template-columns: 2em repeat(53, var(--cell-size, 12px));
  grid-template-rows: auto repeat(7, var(--cell-size, 12px));
  gap: 2px;
}
.cal-month {
  font-size: 0.65rem;
  color: var(--color-text-muted);
  font-family: var(--font-body);
  text-align: left;
  white-space: nowrap;
}
.cal-day {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  font-size: 0.6rem;
  color: var(--color-text-muted);
  font-family: var(--font-body);
  padding-right: 2px;
}
.cal-cell {
  border-radius: 2px;
  transition: transform 0.15s;
  cursor: default;
}
.cal-cell:hover {
  transform: scale(1.3);
  outline: var(--outline-primary);
  outline-offset: 1px;
  z-index: 1;
  position: relative;
}
.cal-lvl0 { background: var(--cal-lvl0); }
.cal-lvl1 { background: var(--cal-lvl1); }
.cal-lvl2 { background: var(--cal-lvl2); }
.cal-lvl3 { background: var(--cal-lvl3); }
.cal-lvl4 { background: var(--cal-lvl4); }
.cal-out { background: transparent; }
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

.chart-row-full {
  grid-template-columns: 1fr;
}
</style>
