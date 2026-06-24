<template>
  <div class="analysis-grid">
    <ChartContainer
:title="t('analytics.charts.netChangeChart.title')"
:subtitle="t('analytics.charts.netChangeChart.subtitle')"
      :option="netChangeOption"
      :loading="loading"
      class="chart-card"
    />
    <ChartContainer
      :title="t('analytics.charts.hourly.title')"
      :subtitle="t('analytics.charts.hourly.subtitle')"
      :option="hourlyOption"
      :loading="loading"
      class="chart-card"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from '../i18n'
import echarts from '../utils/echarts'
import ChartContainer from './ChartContainer.vue'
import { getChartConfig } from '../utils/constants'
import { useTheme } from '../composables/useTheme'

const { t } = useI18n()

const props = defineProps({
  loading: Boolean,
  dailyStats: Array,
  selectedRepos: Array,
  activityHeatmap: Array,
  startDate: String,
  endDate: String
})

const { theme } = useTheme()
const chartCfg = computed(() => getChartConfig(theme.value))

const filteredStats = computed(() => {
  return props.selectedRepos?.length
    ? (props.dailyStats || []).filter(s => props.selectedRepos.includes(s.repoPath))
    : (props.dailyStats || [])
})

function eachDay(start, end) {
  const days = []
  const cur = new Date(start)
  const endDate = new Date(end)
  while (cur <= endDate) {
    days.push(cur.toISOString().slice(0, 10))
    cur.setDate(cur.getDate() + 1)
  }
  return days
}

function eachMonth(start, end) {
  const months = []
  const cur = new Date(start)
  const endDate = new Date(end)
  while (cur <= endDate) {
    months.push(cur.toISOString().slice(0, 7))
    cur.setMonth(cur.getMonth() + 1)
  }
  return months
}

const granularity = computed(() => {
  if (!props.startDate || !props.endDate) return 'day'
  const days = Math.round((new Date(props.endDate) - new Date(props.startDate)) / (1000 * 60 * 60 * 24))
  return days > 31 ? 'month' : 'day'
})

const netChangeData = computed(() => {
  const dateMap = {}
  for (const repo of filteredStats.value) {
    for (const author of (repo.authors || [])) {
      for (const day of (author.dailyData || [])) {
        if (!dateMap[day.date]) {
          dateMap[day.date] = { additions: 0, deletions: 0 }
        }
        dateMap[day.date].additions += day.additions
        dateMap[day.date].deletions += day.deletions
      }
    }
  }
  const isMonth = granularity.value === 'month'
  const dates = props.startDate && props.endDate
    ? (isMonth ? eachMonth(props.startDate, props.endDate) : eachDay(props.startDate, props.endDate))
    : Object.keys(dateMap).sort()
  if (isMonth) {
    const monthMap = {}
    for (const d of dates) {
      let add = 0, del = 0
      Object.keys(dateMap).filter(k => k.startsWith(d)).forEach(k => {
        add += dateMap[k].additions
        del += dateMap[k].deletions
      })
      monthMap[d] = { additions: add, deletions: del }
    }
    return dates.map(d => ({
      date: d,
      additions: monthMap[d].additions,
      deletions: monthMap[d].deletions,
      net: monthMap[d].additions - monthMap[d].deletions
    }))
  }
  return dates.map(d => ({
    date: d,
    additions: dateMap[d]?.additions || 0,
    deletions: dateMap[d]?.deletions || 0,
    net: (dateMap[d]?.additions || 0) - (dateMap[d]?.deletions || 0)
  }))
})

const netChangeOption = computed(() => {
  const data = netChangeData.value
  if (data.length === 0) {
    return {
      title: {
        text: t('analytics.noData'),
        left: 'center',
        top: 'center',
        textStyle: { color: chartCfg.value.titleColor, fontSize: 16 }
      }
    }
  }

  const red = '#ff6b6b'
  const green = '#00ff88'

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: chartCfg.value.tooltipBg,
      borderColor: chartCfg.value.accent,
      textStyle: { color: chartCfg.value.tooltipText },
      formatter: (params) => {
        const d = data[params[0].dataIndex]
        return [
          `<div style="font-weight:bold;margin-bottom:6px">${d.date}</div>`,
          `<div>${t('analytics.charts.additions')}: <span style="color:${red}">+${d.additions}</span></div>`,
          `<div>${t('analytics.charts.deletions')}: <span style="color:${green}">-${d.deletions}</span></div>`
        ].join('')
      }
    },
    grid: { containLabel: true },
    xAxis: {
      type: 'category',
      data: data.map(d => d.date),
      axisLine: { lineStyle: { color: chartCfg.value.axisLine } },
      axisLabel: {
        color: chartCfg.value.axisLabel,
        rotate: data.length > 14 ? 45 : 0
      },
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisLabel: { color: chartCfg.value.axisLabel },
      splitLine: {
        lineStyle: { color: chartCfg.value.splitLine, type: 'dashed' }
      }
    },
    series: [{
      type: 'bar',
      barMaxWidth: 60,
      data: data.map(d => ({
        value: d.net,
        itemStyle: {
          color: d.net >= 0 ? 'transparent' : new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: green },
            { offset: 1, color: adjustColor(green, -60) }
          ]),
          borderColor: d.net >= 0 ? red : 'transparent',
          borderWidth: d.net >= 0 ? 2 : 0,
          borderRadius: d.net >= 0 ? [4, 4, 0, 0] : [0, 0, 4, 4]
        }
      }))
    }]
  }
})

function adjustColor(color, amount) {
  const num = parseInt(color.replace('#', ''), 16)
  const r = Math.min(255, Math.max(0, (num >> 16) + amount))
  const g = Math.min(255, Math.max(0, ((num >> 8) & 0x00FF) + amount))
  const b = Math.min(255, Math.max(0, (num & 0x0000FF) + amount))
  return `#${(0x1000000 + r * 0x10000 + g * 0x100 + b).toString(16).slice(1)}`
}

const hourlyDistribution = computed(() => {
  const h = Array(24).fill(0)
  if (props.activityHeatmap) {
    for (const d of props.activityHeatmap) {
      h[d.hour] += d.commitCount
    }
  }
  return h.map((count, hour) => ({ hour, count }))
})

const hourlyOption = computed(() => {
  const data = hourlyDistribution.value
  if (!data.some(d => d.count > 0)) {
    return {
      title: {
        text: t('analytics.noData'),
        left: 'center',
        top: 'center',
        textStyle: { color: chartCfg.value.titleColor, fontSize: 16 }
      }
    }
  }
  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: chartCfg.value.tooltipBg,
      borderColor: chartCfg.value.accent,
      textStyle: { color: chartCfg.value.tooltipText },
      formatter: (p) => `${p[0].name}:00<br/>${t('analytics.charts.tooltipCommits').replace('{0}', p[0].value)}`
    },
    grid: { containLabel: true },
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
            { offset: 1, color: `rgba(${chartCfg.value.primaryRgb}, 0.15)` }
          ]),
          borderRadius: [3, 3, 0, 0]
        }
      }))
    }]
  }
})
</script>

<style scoped>
.analysis-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(420px, 1fr));
  gap: 1.5rem;
}

.chart-card {
  min-height: 400px;
  height: 400px;
}

@media (max-width: 1200px) {
  .analysis-grid {
    grid-template-columns: 1fr;
  }
}
</style>
