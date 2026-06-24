<template>
  <div class="time-series-grid">
    <ChartContainer
      :title="commitTrendTitle"
      :subtitle="t('analytics.commitTrendSub')"
      :option="commitTrendOption"
      :loading="loading"
      class="chart-card"
    />
    <ChartContainer
      :title="codeChangeTitle"
      :subtitle="t('analytics.codeChangeSub')"
      :option="codeChangeOption"
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
  selectedRepos: Array
})

function adjustColor(color, amount) {
  const num = parseInt(color.replace('#', ''), 16)
  const r = Math.min(255, Math.max(0, (num >> 16) + amount))
  const g = Math.min(255, Math.max(0, ((num >> 8) & 0x00FF) + amount))
  const b = Math.min(255, Math.max(0, (num & 0x0000FF) + amount))
  return `#${(0x1000000 + r * 0x10000 + g * 0x100 + b).toString(16).slice(1)}`
}

const { theme } = useTheme()
const chartCfg = computed(() => getChartConfig(theme.value))

const filteredStats = computed(() => {
  return props.selectedRepos.length === 0
    ? props.dailyStats
    : props.dailyStats.filter(stat => props.selectedRepos.includes(stat.repoPath))
})

const processedStats = computed(() => {
  const data = filteredStats.value
  const dateSet = new Set()
  const authorMap = new Map()

  data.forEach(repo => {
    repo.authors?.forEach(author => {
      author.dailyData?.forEach(day => dateSet.add(day.date))
      if (!authorMap.has(author.email)) {
        authorMap.set(author.email, {
          name: author.author,
          email: author.email,
          isMe: author.isMe
        })
      }
    })
  })

  return {
    dates: Array.from(dateSet).sort(),
    authors: Array.from(authorMap.values())
  }
})

const commitTrendTitle = computed(() => t('analytics.daily') + t('analytics.commitTrend'))
const codeChangeTitle = computed(() => t('analytics.daily') + t('analytics.codeChange'))

const commitTrendOption = computed(() => {
  const { dates, authors } = processedStats.value

  if (dates.length === 0 || authors.length === 0) {
    return {
      title: {
        text: t('analytics.noData'),
        left: 'center',
        top: 'center',
        textStyle: { color: chartCfg.value.titleColor, fontSize: 16 }
      }
    }
  }

  const colors = chartCfg.value.chartColors
  const series = authors.map((author, idx) => {
    const data = dates.map(date => {
      let commits = 0
      filteredStats.value.forEach(repo => {
        const authorStat = repo.authors?.find(a => a.email === author.email)
        if (authorStat && authorStat.dailyData) {
          const dayData = authorStat.dailyData.find(d => d.date === date)
          if (dayData) commits += dayData.commits
        }
      })
      return commits
    })

    return {
      name: author.name,
      type: 'line',
      smooth: true,
      symbol: 'circle',
      symbolSize: 8,
      lineStyle: {
        width: 3,
        color: colors[idx % colors.length]
      },
      itemStyle: {
        color: colors[idx % colors.length],
        borderWidth: 2,
        borderColor: chartCfg.value.tooltipText
      },
      areaStyle: {
        opacity: 0.1,
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: colors[idx % colors.length] },
          { offset: 1, color: 'rgba(0,0,0,0)' }
        ])
      },
      data
    }
  })

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: chartCfg.value.tooltipBg,
      borderColor: chartCfg.value.accent,
      textStyle: { color: chartCfg.value.tooltipText },
      axisPointer: {
        type: 'cross',
        lineStyle: { color: chartCfg.value.accent, opacity: 0.5 }
      }
    },
    legend: {
      data: authors.map(a => a.name),
      textStyle: { color: chartCfg.value.axisLabel },
      top: 10
    },
    grid: { containLabel: true },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dates,
      axisLine: { lineStyle: { color: chartCfg.value.axisLine } },
      axisLabel: {
        color: chartCfg.value.axisLabel
      },
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisLabel: { color: chartCfg.value.axisLabel },
      splitLine: {
        lineStyle: {
          color: chartCfg.value.splitLine,
          type: 'dashed'
        }
      }
    },
    series
  }
})

const dateAggMap = computed(() => {
  const { dates } = processedStats.value
  const map = {}
  for (const date of dates) {
    let totalAdditions = 0
    let totalDeletions = 0
    const authorAddMap = {}
    const authorDelMap = {}
    filteredStats.value.forEach(repo => {
      repo.authors?.forEach(author => {
        if (!author.dailyData) return
        const day = author.dailyData.find(d => d.date === date)
        if (!day) return
        totalAdditions += day.additions
        totalDeletions += day.deletions
        const name = author.author || author.email
        authorAddMap[name] = (authorAddMap[name] || 0) + day.additions
        authorDelMap[name] = (authorDelMap[name] || 0) + day.deletions
      })
    })
    map[date] = { totalAdditions, totalDeletions, authorAddMap, authorDelMap }
  }
  return map
})

const codeChangeOption = computed(() => {
  const { dates } = processedStats.value

  if (dates.length === 0) {
    return {
      title: {
        text: t('analytics.noData'),
        left: 'center',
        top: 'center',
        textStyle: { color: chartCfg.value.titleColor, fontSize: 16 }
      }
    }
  }

  const colors = chartCfg.value.chartColors
  const addColor = colors[0]
  const delColor = colors[1]

  const additionsData = dates.map(d => dateAggMap.value[d]?.totalAdditions || 0)
  const deletionsData = dates.map(d => -(dateAggMap.value[d]?.totalDeletions || 0))

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: chartCfg.value.tooltipBg,
      borderColor: chartCfg.value.accent,
      textStyle: { color: chartCfg.value.tooltipText },
      axisPointer: {
        type: 'shadow',
        shadowStyle: { color: `rgba(${chartCfg.value.accentRgb}, 0.1)` }
      },
      formatter: (params) => {
        const date = params[0].axisValue
        const agg = dateAggMap.value[date]
        if (!agg) return ''
        let html = `<div style="font-weight:bold;margin-bottom:6px">${date}</div>`
        html += `<div style="margin-bottom:4px">${t('analytics.charts.additions')}: <span style="color:${addColor}">+${agg.totalAdditions}</span></div>`
        const sortedAdd = Object.entries(agg.authorAddMap).filter(([, v]) => v > 0).sort((a, b) => b[1] - a[1])
        for (const [name, val] of sortedAdd) {
          html += `<div style="padding-left:12px;font-size:0.85em;color:${chartCfg.value.tooltipText}">${name}: +${val}</div>`
        }
        html += `<div style="margin-top:4px;margin-bottom:4px">${t('analytics.charts.deletions')}: <span style="color:${delColor}">-${agg.totalDeletions}</span></div>`
        const sortedDel = Object.entries(agg.authorDelMap).filter(([, v]) => v > 0).sort((a, b) => b[1] - a[1])
        for (const [name, val] of sortedDel) {
          html += `<div style="padding-left:12px;font-size:0.85em;color:${chartCfg.value.tooltipText}">${name}: -${val}</div>`
        }
        return html
      }
    },
    legend: {
      data: [t('analytics.charts.additions'), t('analytics.charts.deletions')],
      textStyle: { color: chartCfg.value.axisLabel },
      top: 10
    },
    grid: { containLabel: true },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: { lineStyle: { color: chartCfg.value.axisLine } },
      axisLabel: {
        color: chartCfg.value.axisLabel
      },
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisLabel: {
        color: chartCfg.value.axisLabel,
        formatter: (value) => Math.abs(value)
      },
      splitLine: {
        lineStyle: {
          color: chartCfg.value.splitLine,
          type: 'dashed'
        }
      }
    },
    series: [
      {
        name: t('analytics.charts.additions'),
        type: 'bar',
        stack: 'total',
        barMaxWidth: 60,
        emphasis: { focus: 'series' },
        data: additionsData,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: addColor },
            { offset: 1, color: adjustColor(addColor, -30) }
          ]),
          borderRadius: [4, 4, 0, 0]
        }
      },
      {
        name: t('analytics.charts.deletions'),
        type: 'bar',
        stack: 'total',
        barMaxWidth: 60,
        emphasis: { focus: 'series' },
        data: deletionsData,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: adjustColor(delColor, -60) },
            { offset: 1, color: adjustColor(delColor, -90) }
          ]),
          borderRadius: [0, 0, 4, 4]
        }
      }
    ]
  }
})
</script>

<style scoped>
.time-series-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
  gap: 1.5rem;
}

.chart-card {
  min-height: 420px;
  height: 420px;
}

@media (max-width: 1200px) {
  .time-series-grid {
    grid-template-columns: 1fr;
  }
}
</style>
