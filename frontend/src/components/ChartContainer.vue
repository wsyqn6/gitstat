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
      <p>{{ t('repo.loading') }}</p>
    </div>
    <div v-show="!loading" ref="chartRef" class="chart"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useI18n } from '../i18n'
import echarts from '../utils/echarts'

const { t } = useI18n()

const props = defineProps({
  title: String,
  subtitle: String,
  option: Object,
  loading: Boolean
})

const chartRef = ref(null)
let chartInstance = null

onMounted(async () => {
  window.addEventListener('resize', handleResize)
  if (!props.loading && chartRef.value) {
    await nextTick()
    if (!chartInstance) {
      chartInstance = echarts.init(chartRef.value)
    }
    if (props.option) {
      chartInstance.setOption(props.option, true)
      chartInstance.resize()
    }
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance?.dispose()
  chartInstance = null
})

watch(() => props.loading, async (newVal) => {
  if (!newVal && chartRef.value) {
    await nextTick()
    if (!chartInstance) {
      chartInstance = echarts.init(chartRef.value)
    }
    if (props.option) {
      chartInstance.setOption(props.option, true)
      chartInstance.resize()
    }
  }
})

watch(() => props.option, (newVal) => {
  if (newVal && chartInstance) {
    chartInstance.setOption(newVal, true)
    nextTick(() => {
      chartInstance?.resize()
    })
  }
})

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
  font-family: var(--font-display);
  font-size: 1.3rem;
  color: var(--color-accent);
  letter-spacing: 2px;
  text-transform: uppercase;
  margin: 0;
}

.subtitle {
  font-size: 0.85rem;
  color: var(--color-text-muted);
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
  color: var(--color-text-muted);
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 3px solid var(--bg-insight-hover);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  box-shadow: 0 0 20px var(--border-dropdown);
}

.chart {
  height: 400px;
  width: 100%;
}
</style>
