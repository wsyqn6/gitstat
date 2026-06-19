<template>
  <div class="charts-grid">
    <ChartContainer
      :title="commitTrendTitle"
      :subtitle="t('analytics.commitTrendSub')"
      :option="commitTrendOption"
      :loading="loading"
      class="chart-card primary"
    />
    <ChartContainer
      :title="codeChangeTitle"
      :subtitle="t('analytics.codeChangeSub')"
      :option="codeChangeOption"
      :loading="loading"
      class="chart-card secondary"
    />
    <ChartContainer
      :title="t('analytics.charts.devRank.title')"
      :subtitle="t('analytics.charts.devRank.subtitle')"
      :option="authorRankOption"
      :loading="loading"
      class="chart-card tertiary"
    />
    <ChartContainer
      :title="t('analytics.charts.heatmap.title')"
      :subtitle="t('analytics.charts.heatmap.subtitle')"
      :option="heatmapOption"
      :loading="loading"
      class="chart-card quaternary"
    />
    <ChartContainer
      :title="t('analytics.charts.hourly.title')"
      :subtitle="t('analytics.charts.hourly.subtitle')"
      :option="hourlyOption"
      :loading="loading"
      class="chart-card tertiary"
    />
    <ChartContainer
      :title="t('analytics.charts.repoCompare.title')"
      :subtitle="t('analytics.charts.repoCompare.subtitle')"
      :option="repoComparisonOption"
      :loading="loading"
      class="chart-card quinary"
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
  periodStats: Array,
  currentGranularity: String,
  selectedRepos: Array,
  authorRank: Array,
  activityHeatmap: Array,
  repoComparison: Array
})

function formatPeriodLabel(period, granularity) {
  if (granularity === 'week') {
    const match = period.match(/-W(\d+)/)
    return match ? t('analytics.charts.weekLabel').replace('{0}', match[1]) : period
  }
  if (granularity === 'month')
    return t('analytics.charts.monthLabel').replace('{0}', period.slice(5))
  if (granularity === 'year')
    return t('analytics.charts.yearLabel').replace('{0}', period)
  return period
}

function adjustColor(color, amount) {
  const num = parseInt(color.replace('#', ''), 16)
  const r = Math.min(255, Math.max(0, (num >> 16) + amount))
  const g = Math.min(255, Math.max(0, ((num >> 8) & 0x00FF) + amount))
  const b = Math.min(255, Math.max(0, (num & 0x0000FF) + amount))
  return `#${(0x1000000 + r * 0x10000 + g * 0x100 + b).toString(16).slice(1)}`
}

const { theme } = useTheme()
const chartCfg = computed(() => getChartConfig(theme.value))

const activeStats = computed(() => {
  return props.currentGranularity === 'day' ? props.dailyStats : props.periodStats
})

const filteredStats = computed(() => {
  const source = activeStats.value
  return props.selectedRepos.length === 0
    ? source
    : source.filter(stat => props.selectedRepos.includes(stat.repoPath))
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

const granularityPrefix = computed(() => {
  switch (props.currentGranularity) {
    case 'week': return t('analytics.weekly')
    case 'month': return t('analytics.monthly')
    case 'year': return t('analytics.yearly')
    default: return t('analytics.daily')
  }
})

const commitTrendTitle = computed(() => granularityPrefix.value + t('analytics.commitTrend'))
const codeChangeTitle = computed(() => granularityPrefix.value + t('analytics.codeChange'))

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
        color: chartCfg.value.axisLabel,
        formatter: (value) => formatPeriodLabel(value, props.currentGranularity)
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

const codeChangeOption = computed(() => {
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
  const series = []

  authors.forEach((author, idx) => {
    const baseColor = colors[idx % colors.length]

    const netData = dates.map(date => {
      let additions = 0
      let deletions = 0
      filteredStats.value.forEach(repo => {
        const authorStat = repo.authors?.find(a => a.email === author.email)
        if (authorStat && authorStat.dailyData) {
          const dayData = authorStat.dailyData.find(d => d.date === date)
          if (dayData) {
            additions += dayData.additions
            deletions += dayData.deletions
          }
        }
      })
      return { additions, deletions }
    })

    series.push({
      name: author.name,
      type: 'bar',
      stack: 'total',
      barWidth: '40%',
      emphasis: { focus: 'series' },
      data: netData.map(item => item.additions),
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: baseColor },
          { offset: 1, color: adjustColor(baseColor, -30) }
        ]),
        borderRadius: [4, 4, 0, 0]
      }
    })

    series.push({
      name: `${author.name} - ${t('analytics.charts.deletions')}`,
      type: 'bar',
      stack: 'total',
      barWidth: '40%',
      emphasis: { focus: 'series' },
      data: netData.map(item => -item.deletions),
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: adjustColor(baseColor, -60) },
          { offset: 1, color: adjustColor(baseColor, -90) }
        ]),
        borderRadius: [0, 0, 4, 4]
      }
    })
  })

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: chartCfg.value.tooltipBg,
      borderColor: chartCfg.value.accent,
      textStyle: { color: chartCfg.value.tooltipText },
      axisPointer: {
        type: 'shadow',
        shadowStyle: { color: 'rgba(0, 245, 255, 0.1)' }
      },
      formatter: (params) => {
        let result = `<div style="font-weight:bold;margin-bottom:5px">${params[0].axisValue}</div>`
        params.forEach(param => {
          const value = Math.abs(param.value)
          const type = param.value >= 0 ? t('analytics.charts.additions') : t('analytics.charts.deletions')
          result += `<div>${param.seriesName}: <span style="color:${param.color}">${value}</span> (${type})</div>`
        })
        return result
      }
    },
    legend: {
      data: series.map(s => s.name),
      textStyle: { color: chartCfg.value.axisLabel },
      top: 10,
      type: 'scroll'
    },
    grid: { containLabel: true },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: { lineStyle: { color: chartCfg.value.axisLine } },
      axisLabel: {
        color: chartCfg.value.axisLabel,
        formatter: (value) => formatPeriodLabel(value, props.currentGranularity)
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
    series
  }
})

const authorRankOption = computed(() => {
  if (!props.authorRank || props.authorRank.length === 0) {
    return {
      title: {
        text: t('analytics.noData'),
        left: 'center',
        top: 'center',
        textStyle: { color: chartCfg.value.titleColor, fontSize: 16 }
      }
    }
  }

  const topAuthors = props.authorRank.slice(0, 10)
  const colors = chartCfg.value.chartColors

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: chartCfg.value.tooltipBg,
      borderColor: chartCfg.value.accent,
      textStyle: { color: chartCfg.value.tooltipText },
      axisPointer: { type: 'shadow' },
      formatter: (params) => {
        const item = params[0]
        const author = topAuthors[item.dataIndex]
        return `
          <div style="font-weight:bold;margin-bottom:5px">${author.author}</div>
          <div>${t('analytics.charts.tooltipCommits').replace('{0}', author.commits)}</div>
          <div>${t('analytics.charts.tooltipAdditions').replace('{0}', author.additions)}</div>
          <div>${t('analytics.charts.tooltipDeletions').replace('{0}', author.deletions)}</div>
          <div>${t('analytics.charts.tooltipNetChange').replace('{0}', (author.netChange > 0 ? '+' : '') + author.netChange)}</div>
        `
      }
    },
    grid: { containLabel: true },
    xAxis: {
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
    yAxis: {
      type: 'category',
      data: topAuthors.map(a => a.author),
      axisLine: { lineStyle: { color: chartCfg.value.axisLine } },
      axisLabel: { color: chartCfg.value.axisLabel },
      inverse: true
    },
    series: [{
      type: 'bar',
      data: topAuthors.map((a, idx) => ({
        value: a.commits,
        itemStyle: {
          color: new echarts.graphic.LinearGradient(1, 0, 0, 0, [
            { offset: 0, color: colors[idx % colors.length] },
            { offset: 1, color: adjustColor(colors[idx % colors.length], -40) }
          ]),
          borderRadius: [0, 4, 4, 0]
        }
      })),
      barWidth: '60%'
    }]
  }
})

const heatmapOption = computed(() => {
  if (!props.activityHeatmap || props.activityHeatmap.length === 0) {
    return {
      title: {
        text: t('analytics.noData'),
        left: 'center',
        top: 'center',
        textStyle: { color: chartCfg.value.titleColor, fontSize: 16 }
      }
    }
  }

  const hours = t('analytics.charts.hours')
  const days = t('analytics.charts.dayNames')

  const data = props.activityHeatmap.map(item => [
    item.dayOfWeek,
    item.hour,
    item.commitCount
  ])

  const maxCount = Math.max(...data.map(d => d[2]))

  return {
    tooltip: {
      backgroundColor: chartCfg.value.tooltipBg,
      borderColor: chartCfg.value.accent,
      textStyle: { color: chartCfg.value.tooltipText },
      formatter: (params) => {
        return `<div>${days[params.value[0]]} ${hours[params.value[1]]}</div><div>${t('analytics.charts.tooltipCommits').replace('{0}', params.value[2])}</div>`
      }
    },
    grid: { containLabel: true },
    xAxis: {
      type: 'category',
      data: days,
      axisLine: { lineStyle: { color: chartCfg.value.axisLine } },
      axisLabel: { color: chartCfg.value.axisLabel },
      splitArea: { show: true, areaStyle: { color: 'rgba(30, 41, 59, 0.3)' } }
    },
    yAxis: {
      type: 'category',
      data: hours,
      axisLine: { lineStyle: { color: chartCfg.value.axisLine } },
      axisLabel: { color: chartCfg.value.axisLabel },
      splitArea: { show: true, areaStyle: { color: 'rgba(30, 41, 59, 0.3)' } }
    },
    visualMap: {
      min: 0,
      max: maxCount,
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      bottom: '5%',
      inRange: { color: chartCfg.value.heatmapColors },
      textStyle: { color: chartCfg.value.axisLabel }
    },
    series: [{
      type: 'heatmap',
      data: data,
      label: { show: false },
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowColor: 'rgba(0, 245, 255, 0.5)'
        }
      }
    }]
  }
})

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
            { offset: 1, color: 'rgba(0, 212, 255, 0.15)' }
          ]),
          borderRadius: [3, 3, 0, 0]
        }
      }))
    }]
  }
})

const repoComparisonOption = computed(() => {
  if (!props.repoComparison || props.repoComparison.length === 0) {
    return {
      title: {
        text: t('analytics.noData'),
        left: 'center',
        top: 'center',
        textStyle: { color: chartCfg.value.titleColor, fontSize: 16 }
      }
    }
  }

  const repos = props.repoComparison.slice(0, 5)
  const colors = chartCfg.value.chartColors

  const maxCommits = Math.max(...repos.map(r => r.commits))
  const maxAdditions = Math.max(...repos.map(r => r.additions))
  const maxAuthors = Math.max(...repos.map(r => r.authors))
  const maxActiveDays = Math.max(...repos.map(r => r.activeDays))

  return {
    tooltip: {
      backgroundColor: chartCfg.value.tooltipBg,
      borderColor: chartCfg.value.accent,
      textStyle: { color: chartCfg.value.tooltipText }
    },
    legend: {
      data: repos.map(r => r.repoName),
      textStyle: { color: chartCfg.value.axisLabel },
      top: 10
    },
    radar: {
      indicator: [
        { name: t('analytics.charts.commits'), max: 100 },
        { name: t('analytics.charts.codeVolume'), max: 100 },
        { name: t('analytics.charts.authorCount'), max: 100 },
        { name: t('analytics.charts.activeDays'), max: 100 },
        { name: t('analytics.charts.dailyAvg'), max: 100 }
      ],
      shape: 'polygon',
      splitNumber: 4,
      axisName: { color: chartCfg.value.axisLabel },
      splitLine: {
        lineStyle: { color: 'rgba(0, 245, 255, 0.2)' }
      },
      splitArea: {
        areaStyle: {
          color: ['rgba(0, 245, 255, 0.05)', 'rgba(0, 245, 255, 0.1)']
        }
      }
    },
    series: [{
      type: 'radar',
      data: repos.map((repo, idx) => {
        const avgCommitsPerDayNormalized = repo.avgCommitsPerDay > 10 ? 100 : (repo.avgCommitsPerDay / 10) * 100
        return {
          value: [
            (repo.commits / maxCommits) * 100,
            (repo.additions / maxAdditions) * 100,
            (repo.authors / maxAuthors) * 100,
            (repo.activeDays / maxActiveDays) * 100,
            avgCommitsPerDayNormalized
          ],
          name: repo.repoName,
          lineStyle: { color: colors[idx % colors.length], width: 2 },
          itemStyle: { color: colors[idx % colors.length] },
          areaStyle: { opacity: 0.2, color: colors[idx % colors.length] }
        }
      })
    }]
  }
})
</script>

<style scoped>
.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
  gap: 1.5rem;
}

.chart-card {
  min-height: 450px;
  height: 450px;
}

@media (max-width: 1200px) {
  .charts-grid {
    grid-template-columns: 1fr;
  }
}
</style>
