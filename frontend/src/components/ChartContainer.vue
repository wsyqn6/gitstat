<template>
  <div class="chart-container card">
    <div class="chart-header">
      <div class="title-section">
        <h3>{{ title }}</h3>
        <p v-if="subtitle" class="subtitle">{{ subtitle }}</p>
      </div>
      <div class="tech-lines">
        <span></span>
        <span></span>
        <span></span>
      </div>
    </div>
    <div v-if="loading" class="chart-loading">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>
    <div v-else ref="chartRef" class="chart"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import * as echarts from 'echarts'

const props = defineProps({
  title: String,
  subtitle: String,
  option: Object,
  loading: Boolean
})

const chartRef = ref(null)
let chartInstance = null

onMounted(() => {
  chartInstance = echarts.init(chartRef.value)
  updateChart()
  
  // 响应窗口大小变化
  window.addEventListener('resize', handleResize)
})

watch(() => props.option, updateChart, { deep: true })

function updateChart() {
  if (chartInstance && props.option) {
    chartInstance.setOption(props.option, true)
  }
}

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

.tech-lines {
  display: flex;
  gap: 8px;
  margin-top: 1rem;
}

.tech-lines span {
  height: 2px;
  background: linear-gradient(90deg, #00f5ff, transparent);
  animation: techLine 2s ease-in-out infinite;
}

.tech-lines span:nth-child(1) {
  width: 60px;
  animation-delay: 0s;
}

.tech-lines span:nth-child(2) {
  width: 40px;
  animation-delay: 0.2s;
}

.tech-lines span:nth-child(3) {
  width: 20px;
  animation-delay: 0.4s;
}

@keyframes techLine {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 1; }
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
}
</style>
