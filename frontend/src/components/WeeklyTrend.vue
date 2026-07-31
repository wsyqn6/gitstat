<template>
  <div class="insight-card glass card">
    <div class="insight-header">
      <h3>{{ t('dashboard.weeklyTrend') }} <span class="range-hint">{{ weekRange }}</span></h3>
    </div>
    <div v-if="loading" class="skeleton-chart-area trend-skeleton">
      <Skeleton w="80" />
      <Skeleton w="60" />
      <Skeleton w="90" />
      <Skeleton w="40" />
      <Skeleton w="70" />
      <Skeleton w="50" />
    </div>
    <div v-else-if="repoDailyTrend.length === 0" class="card-empty">
      <div class="card-empty-text">{{ t('analytics.noData') }}</div>
    </div>
    <template v-else>
      <div ref="trendChartRef" class="trend-chart"></div>
      <div v-if="repoDailyTrend.length > 1" class="chart-legend">
        <span v-for="(repo, i) in repoDailyTrend" :key="repo.repoName" class="legend-item">
          <span class="legend-dot" :style="{ background: repoColors[i] }"></span>
          {{ repo.repoName }}
        </span>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import Skeleton from './Skeleton.vue'
import { useI18n } from '../i18n'
import { getChartConfig } from '../utils/constants'
import { useTheme } from '../composables/useTheme'
import echarts from '../utils/echarts'

const { t } = useI18n()
const { theme } = useTheme()
const chartCfg = computed(() => getChartConfig(theme.value))

const props = defineProps({
  repoDailyTrend: { type: Array, required: true },
  repoColors: { type: Array, required: true },
  weekRange: { type: String, required: true },
  loading: { type: Boolean, required: true }
})

const trendChartRef = ref(null)
let trendChart = null

function renderTrendChart() {
  if (!trendChartRef.value || props.repoDailyTrend.length === 0) return

  if (!trendChart) {
    trendChart = echarts.init(trendChartRef.value)
  }

  const cfg = chartCfg.value
  const colors = cfg.chartColors8
  const allDates = [...new Set(props.repoDailyTrend.flatMap(r => r.data.map(d => d.date)))].sort()
  const labels = allDates.map(d => d.slice(5))
  const singlePoint = allDates.length <= 1

  const series = props.repoDailyTrend.map((repo, i) => {
    const dateMap = Object.fromEntries(repo.data.map(d => [d.date, d.commits]))
    return {
      name: repo.repoName,
      type: 'line',
      stack: 'total',
      smooth: !singlePoint,
      symbol: singlePoint ? 'circle' : 'none',
      lineStyle: { width: 1.5, color: colors[i % colors.length] },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: colors[i % colors.length] + '60' },
          { offset: 1, color: colors[i % colors.length] + '05' }
        ])
      },
      itemStyle: { color: colors[i % colors.length] },
      data: allDates.map(d => dateMap[d] || 0)
    }
  })

  trendChart.setOption({
    tooltip: {
      trigger: 'axis',
      backgroundColor: cfg.tooltipBg,
      borderColor: cfg.accent + '4D',
      textStyle: { color: cfg.tooltipText, fontSize: 12 }
    },
    legend: { show: false },
    grid: { left: 40, right: 16, top: 16, bottom: 24 },
    xAxis: {
      type: 'category',
      data: labels,
      axisLine: { lineStyle: { color: cfg.axisLine } },
      axisLabel: { color: cfg.axisLabel, fontSize: 11 }
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      splitLine: { lineStyle: { color: cfg.splitLine, type: 'dashed' } },
      axisLabel: { color: cfg.axisLabel, fontSize: 11 }
    },
    series
  })
}

watch([() => props.loading, () => props.repoDailyTrend, chartCfg], () => {
  if (props.loading) return
  nextTick(() => renderTrendChart())
})

function handleResize() {
  if (trendChart) trendChart.resize()
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  trendChart = null
})
</script>

<style scoped>
.trend-skeleton {
  height: 220px;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  padding: 20px 10px;
}

.trend-skeleton .skeleton-line {
  height: 10px;
}

.trend-chart {
  height: 220px;
  width: 100%;
}

.chart-legend {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  padding: 0.75rem 0 0 0;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.75rem;
  color: var(--color-nav-link);
  font-family: var(--font-body);
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
</style>
