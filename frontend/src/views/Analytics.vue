<template>
  <div class="analytics">
    <h2 class="page-title">数据分析</h2>
    
    <ChartContainer 
      title="提交趋势" 
      :option="trendChartOption" 
    />
    
    <ChartContainer 
      title="代码变更分布" 
      :option="changeChartOption" 
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { state } from '../stores/data'
import ChartContainer from '../components/ChartContainer.vue'

const trendChartOption = computed(() => {
  return {
    xAxis: { type: 'category', data: ['Mon', 'Tue', 'Wed'] },
    yAxis: { type: 'value' },
    series: [{ data: [10, 20, 15], type: 'line' }]
  }
})

const changeChartOption = computed(() => {
  return {
    xAxis: { type: 'category', data: ['Additions', 'Deletions'] },
    yAxis: { type: 'value' },
    series: [{ 
      data: [
        state.overviewStats?.totalAdditions || 0,
        state.overviewStats?.totalDeletions || 0
      ], 
      type: 'bar' 
    }]
  }
})
</script>

<style scoped>
.analytics {
  max-width: 1400px;
}

.page-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 2rem;
  color: #00d4ff;
  letter-spacing: 3px;
  text-transform: uppercase;
  margin-bottom: 2rem;
  background: linear-gradient(135deg, #00d4ff, #7800ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
</style>
