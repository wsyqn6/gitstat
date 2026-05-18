<template>
  <div class="chart-container card">
    <div class="chart-header">
      <h3>{{ title }}</h3>
      <div class="tech-lines">
        <span></span>
        <span></span>
        <span></span>
      </div>
    </div>
    <div ref="chartRef" class="chart"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import * as echarts from 'echarts'

const props = defineProps({
  title: String,
  option: Object
})

const chartRef = ref(null)
let chartInstance = null

onMounted(() => {
  chartInstance = echarts.init(chartRef.value)
  updateChart()
})

watch(() => props.option, updateChart)

function updateChart() {
  if (chartInstance && props.option) {
    chartInstance.setOption(props.option)
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

.chart-header h3 {
  font-family: 'Orbitron', sans-serif;
  font-size: 1.3rem;
  color: #00d4ff;
  letter-spacing: 2px;
  text-transform: uppercase;
  margin: 0;
}

.tech-lines {
  display: flex;
  gap: 8px;
  margin-top: 1rem;
}

.tech-lines span {
  height: 2px;
  background: linear-gradient(90deg, #00d4ff, transparent);
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

.chart {
  height: 400px;
}
</style>
