<template>
  <div class="analysis-grid">
    <div class="chart-card rank-card">
      <RepoAuthorRank :contributors="contributors" :loading="loading" />
    </div>
    <ChartContainer
      :title="t('analytics.charts.heatmap.title')"
      :subtitle="t('analytics.charts.heatmap.subtitle')"
      :option="heatmapOption"
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
import ChartContainer from './ChartContainer.vue'
import RepoAuthorRank from './RepoAuthorRank.vue'
import { getChartConfig } from '../utils/constants'
import { useTheme } from '../composables/useTheme'

const { t } = useI18n()

const props = defineProps({
  loading: Boolean,
  authorRank: Array,
  activityHeatmap: Array
})

const { theme } = useTheme()
const chartCfg = computed(() => getChartConfig(theme.value))

const contributors = computed(() =>
  (props.authorRank || []).map(a => ({
    ...a,
    commitCount: a.commits
  }))
)

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
          shadowColor: `rgba(${chartCfg.value.accentRgb}, 0.5)`
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

.rank-card {
  min-height: 400px;
  height: 400px;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
}

.rank-card :deep(.author-rank) {
  flex: 1;
}

@media (max-width: 1200px) {
  .analysis-grid {
    grid-template-columns: 1fr;
  }
}
</style>
