<template>
  <div class="chart-container card">
    <div class="chart-header">
      <div class="title-section">
        <h3>{{ title }}</h3>
        <p v-if="subtitle" class="subtitle">{{ subtitle }}</p>
      </div>
    </div>
    <div v-show="loading" class="chart-loading">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>
    <div v-show="!loading" ref="chartRef" class="chart"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'
import echarts from '../utils/echarts'

const props = defineProps({
  title: String,
  subtitle: String,
  option: Object,
  loading: Boolean
})

const chartRef = ref(null)
let chartInstance = null

onMounted(() => {
  // 不在 mounted 时初始化，等待 loading 结束
  window.addEventListener('resize', handleResize)
})

watch(() => props.loading, async (newVal) => {
  console.log('[ChartContainer] loading changed:', newVal, 'chartRef:', !!chartRef.value)
  if (!newVal && chartRef.value) {
    // loading 结束后初始化或更新图表
    await nextTick()
    console.log('[ChartContainer] After nextTick, initializing chart, has option:', !!props.option)
    if (!chartInstance) {
      chartInstance = echarts.init(chartRef.value)
      console.log('[ChartContainer] Chart instance created')
    }
    if (props.option) {
      console.log('[ChartContainer] Setting initial option')
      chartInstance.setOption(props.option, true)
      chartInstance.resize()
    } else {
      console.warn('[ChartContainer] No option to set')
    }
  } else if (!newVal) {
    console.warn('[ChartContainer] loading=false but chartRef is null')
  }
}, { immediate: false })

watch(() => props.option, (newVal) => {
  console.log('[ChartContainer] option changed, has chartInstance:', !!chartInstance)
  if (newVal && chartInstance) {
    console.log('[ChartContainer] Updating chart with new option')
    chartInstance.setOption(newVal, true)
    setTimeout(() => {
      if (chartInstance) {
        chartInstance.resize()
      }
    }, 50)
  } else if (!chartInstance) {
    console.warn('[ChartContainer] No chart instance to update')
  }
}, { deep: true })

function handleResize() {
  if (chartInstance) {
    chartInstance.resize()
  }
}
</script>

<style scoped>
.chart-container {
  margin-bottom: 2rem;
}

.chart-header {
  position: relative;
  margin-bottom: 1.5rem;
}

.title-section {
  display: flex;
  align-items: baseline;
  gap: 1rem;
  margin-bottom: 0.5rem;
}

.chart-header h3 {
  font-family: 'Orbitron', sans-serif;
  font-size: 1.3rem;
  color: #00f5ff;
  letter-spacing: 2px;
  text-transform: uppercase;
  margin: 0;
}

.subtitle {
  font-size: 0.85rem;
  color: #64748b;
  letter-spacing: 1px;
  margin: 0;
}

.chart-loading {
  height: 400px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  color: #64748b;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 3px solid rgba(0, 245, 255, 0.1);
  border-top-color: #00f5ff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  box-shadow: 0 0 20px rgba(0, 245, 255, 0.3);
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.chart {
  height: 400px;
  width: 100%;
}
</style>
