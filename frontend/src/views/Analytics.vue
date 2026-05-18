<template>
  <div class="analytics">
    <!-- 标题和控制面板整合 -->
    <div class="header-section">
      <div class="title-row">
        <div>
          <h2 class="page-title">数据分析中心</h2>
          <p class="subtitle">实时监控 · 多维度洞察</p>
        </div>
        
        <div class="controls-inline">
          <div class="control-item">
            <label class="control-label">时间范围</label>
            <select v-model="selectedTimeRange" class="cyber-select">
              <option value="7">本周 (7天)</option>
              <option value="30">本月 (30天)</option>
              <option value="90">近三月 (90天)</option>
              <option value="all">全部时间</option>
            </select>
          </div>

          <div class="control-item">
            <label class="control-label">仓库筛选</label>
            <div class="repo-dropdown">
              <button 
                @click.stop="showRepoDropdown = !showRepoDropdown"
                class="repo-dropdown-btn"
                :class="{ active: showRepoDropdown }"
              >
                <span class="btn-text">
                  {{ selectedRepos.length === 0 ? '选择仓库' : 
                     selectedRepos.length === repositories.length ? '全部仓库' : 
                     `${selectedRepos.length} 个仓库` }}
                </span>
                <svg class="dropdown-icon" :class="{ rotated: showRepoDropdown }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"></polyline>
                </svg>
              </button>
              
              <div v-show="showRepoDropdown" class="repo-dropdown-menu">
                <div class="dropdown-header">
                  <button @click.stop="toggleAllRepos" class="select-all-btn">
                    {{ allReposSelected ? '取消全选' : '全选' }}
                  </button>
                  <span class="selected-count">{{ selectedRepos.length }}/{{ repositories.length }}</span>
                </div>
                <div class="dropdown-list">
                  <button
                    v-for="repo in repositories"
                    :key="repo.path"
                    @click.stop="toggleRepo(repo.path)"
                    class="repo-option"
                    :class="{ active: selectedRepos.includes(repo.path) }"
                  >
                    <div class="option-checkbox">
                      <svg v-if="selectedRepos.includes(repo.path)" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                        <polyline points="20 6 9 17 4 12"></polyline>
                      </svg>
                    </div>
                    <span class="option-text">{{ repo.name }}</span>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div class="control-item">
            <button @click="loadData" :disabled="loading" class="analyze-btn">
              <svg v-if="!loading" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
              </svg>
              <div v-else class="btn-spinner"></div>
              <span>{{ loading ? '分析中...' : '开始分析' }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="charts-grid">
      <!-- 提交趋势图 -->
      <ChartContainer 
        title="每日提交趋势" 
        subtitle="多用户对比分析"
        :option="commitTrendOption" 
        :loading="loading"
        class="chart-card primary"
      />
      
      <!-- 代码变更图 -->
      <ChartContainer 
        title="代码变更分布" 
        subtitle="新增 vs 删除"
        :option="codeChangeOption" 
        :loading="loading"
        class="chart-card secondary"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import ChartContainer from '../components/ChartContainer.vue'
import { getDailyStats, getRepositories } from '../api'

const loading = ref(false)
const selectedTimeRange = ref('7')
const selectedRepos = ref([])
const repositories = ref([])
const dailyStats = ref([])
const showRepoDropdown = ref(false)

// 计算所有仓库是否都选中
const allReposSelected = computed(() => {
  return repositories.value.length > 0 && 
         selectedRepos.value.length === repositories.value.length
})

// 切换单个仓库
const toggleRepo = (repoPath) => {
  const index = selectedRepos.value.indexOf(repoPath)
  if (index > -1) {
    selectedRepos.value.splice(index, 1)
  } else {
    selectedRepos.value.push(repoPath)
  }
}

// 全选/取消全选
const toggleAllRepos = () => {
  if (allReposSelected.value) {
    selectedRepos.value = []
  } else {
    selectedRepos.value = repositories.value.map(r => r.path)
  }
}

// 获取过滤后的统计数据
const filteredStats = computed(() => {
  if (selectedRepos.value.length === 0) {
    return dailyStats.value
  }
  return dailyStats.value.filter(stat => 
    selectedRepos.value.includes(stat.repoPath)
  )
})

// 提取所有日期
const allDates = computed(() => {
  const dateSet = new Set()
  filteredStats.value.forEach(repo => {
    repo.authors?.forEach(author => {
      author.dailyData?.forEach(day => dateSet.add(day.date))
    })
  })
  return Array.from(dateSet).sort()
})

// 提取所有用户
const allAuthors = computed(() => {
  const authorMap = new Map()
  filteredStats.value.forEach(repo => {
    repo.authors?.forEach(author => {
      if (!authorMap.has(author.email)) {
        authorMap.set(author.email, {
          name: author.author,
          email: author.email,
          isMe: author.isMe
        })
      }
    })
  })
  return Array.from(authorMap.values())
})

// 提交趋势图配置
const commitTrendOption = computed(() => {
  const dates = allDates.value
  const authors = allAuthors.value
  
  if (dates.length === 0 || authors.length === 0) {
    return {
      title: { 
        text: '暂无数据', 
        left: 'center', 
        top: 'center',
        textStyle: { color: '#64748b', fontSize: 16 }
      }
    }
  }

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

    const colors = ['#00f5ff', '#ff00ff', '#ffd700', '#00ff88', '#ff6b6b', '#a78bfa']
    
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
        borderColor: '#fff'
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
      backgroundColor: 'rgba(10, 14, 39, 0.95)',
      borderColor: '#00f5ff',
      textStyle: { color: '#fff' },
      axisPointer: {
        type: 'cross',
        lineStyle: { color: '#00f5ff', opacity: 0.5 }
      }
    },
    legend: {
      data: authors.map(a => a.name),
      textStyle: { color: '#94a3b8' },
      top: 10
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dates,
      axisLine: { lineStyle: { color: '#334155' } },
      axisLabel: { color: '#94a3b8' },
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisLabel: { color: '#94a3b8' },
      splitLine: { 
        lineStyle: { 
          color: '#1e293b',
          type: 'dashed'
        } 
      }
    },
    series
  }
})

// 代码变更图配置
const codeChangeOption = computed(() => {
  const dates = allDates.value
  const authors = allAuthors.value
  
  if (dates.length === 0 || authors.length === 0) {
    return {
      title: { 
        text: '暂无数据', 
        left: 'center', 
        top: 'center',
        textStyle: { color: '#64748b', fontSize: 16 }
      }
    }
  }

  // 为每个用户创建新增和删除系列
  const series = []
  const colors = ['#00f5ff', '#ff00ff', '#ffd700', '#00ff88', '#ff6b6b', '#a78bfa']
  
  authors.forEach((author, idx) => {
    const baseColor = colors[idx % colors.length]
    
    // 新增系列
    const additionsData = dates.map(date => {
      let additions = 0
      filteredStats.value.forEach(repo => {
        const authorStat = repo.authors?.find(a => a.email === author.email)
        if (authorStat && authorStat.dailyData) {
          const dayData = authorStat.dailyData.find(d => d.date === date)
          if (dayData) additions += dayData.additions
        }
      })
      return additions
    })

    // 删除系列
    const deletionsData = dates.map(date => {
      let deletions = 0
      filteredStats.value.forEach(repo => {
        const authorStat = repo.authors?.find(a => a.email === author.email)
        if (authorStat && authorStat.dailyData) {
          const dayData = authorStat.dailyData.find(d => d.date === date)
          if (dayData) deletions += dayData.deletions
        }
      })
      return deletions
    })

    series.push({
      name: `${author.name} - 新增`,
      type: 'bar',
      stack: `add-${idx}`,
      emphasis: { focus: 'series' },
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: baseColor },
          { offset: 1, color: adjustColor(baseColor, -30) }
        ]),
        borderRadius: [4, 4, 0, 0]
      },
      data: additionsData
    })

    series.push({
      name: `${author.name} - 删除`,
      type: 'bar',
      stack: `del-${idx}`,
      emphasis: { focus: 'series' },
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: adjustColor(baseColor, -60) },
          { offset: 1, color: adjustColor(baseColor, -90) }
        ]),
        borderRadius: [0, 0, 4, 4]
      },
      data: deletionsData.map(v => -v) // 负值向下显示
    })
  })

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(10, 14, 39, 0.95)',
      borderColor: '#00f5ff',
      textStyle: { color: '#fff' },
      axisPointer: {
        type: 'shadow',
        shadowStyle: { color: 'rgba(0, 245, 255, 0.1)' }
      },
      formatter: (params) => {
        let result = `<strong>${params[0].axisValue}</strong><br/>`
        params.forEach(param => {
          const value = Math.abs(param.value)
          const marker = param.marker
          result += `${marker} ${param.seriesName}: ${value}<br/>`
        })
        return result
      }
    },
    legend: {
      data: series.map(s => s.name),
      textStyle: { color: '#94a3b8' },
      top: 10,
      type: 'scroll'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: { lineStyle: { color: '#334155' } },
      axisLabel: { color: '#94a3b8' },
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisLabel: { 
        color: '#94a3b8',
        formatter: (value) => Math.abs(value)
      },
      splitLine: { 
        lineStyle: { 
          color: '#1e293b',
          type: 'dashed'
        } 
      }
    },
    series
  }
})

// 辅助函数：调整颜色亮度
function adjustColor(color, amount) {
  const num = parseInt(color.replace('#', ''), 16)
  const r = Math.min(255, Math.max(0, (num >> 16) + amount))
  const g = Math.min(255, Math.max(0, ((num >> 8) & 0x00FF) + amount))
  const b = Math.min(255, Math.max(0, (num & 0x0000FF) + amount))
  return `#${(0x1000000 + r * 0x10000 + g * 0x100 + b).toString(16).slice(1)}`
}

// 加载数据
const loadData = async () => {
  if (repositories.value.length === 0) {
    try {
      repositories.value = await getRepositories()
    } catch (error) {
      console.error('加载仓库失败:', error)
      return
    }
  }
  
  loading.value = true
  try {
    const stats = await getDailyStats(null, selectedTimeRange.value, selectedRepos.value)
    console.log('后端返回数据:', JSON.stringify(stats, null, 2))
    dailyStats.value = stats
    console.log('dailyStats:', dailyStats.value)
    console.log('filteredStats:', filteredStats.value)
    console.log('allDates:', allDates.value)
    console.log('allAuthors:', allAuthors.value)
    console.log('commitTrendOption:', commitTrendOption.value)
  } catch (error) {
    console.error('加载数据失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.analytics {
  max-width: 1600px;
  margin: 0 auto;
  padding: 1.5rem 2rem;
}

/* 标题区域 */
.header-section {
  margin-bottom: 1.5rem;
  position: relative;
}

.title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 2rem;
  flex-wrap: wrap;
}

.page-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #00f5ff 0%, #ff00ff 50%, #ffd700 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 3px;
  text-transform: uppercase;
  margin: 0 0 0.25rem 0;
  filter: drop-shadow(0 0 15px rgba(0, 245, 255, 0.3));
}

.subtitle {
  font-size: 0.85rem;
  color: #64748b;
  letter-spacing: 4px;
  text-transform: uppercase;
  margin: 0;
}

/* 控制面板 */
.controls-inline {
  display: flex;
  gap: 1.5rem;
  align-items: flex-end;
}

.control-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.analyze-btn {
  background: linear-gradient(135deg, #00f5ff 0%, #ff00ff 100%);
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1.5rem;
  color: white;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  outline: none;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-family: 'Orbitron', sans-serif;
  letter-spacing: 1px;
  box-shadow: 0 4px 15px rgba(0, 245, 255, 0.4);
  height: fit-content;
  align-self: flex-end;
}

.analyze-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 25px rgba(0, 245, 255, 0.6);
}

.analyze-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.analyze-btn svg {
  width: 18px;
  height: 18px;
}

.btn-spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.control-label {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.75rem;
  color: #00f5ff;
  letter-spacing: 2px;
  text-transform: uppercase;
  font-weight: 600;
}

.cyber-select {
  background: rgba(10, 14, 39, 0.8);
  border: 2px solid rgba(0, 245, 255, 0.3);
  border-radius: 8px;
  padding: 0.6rem 1.25rem;
  color: #e2e8f0;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  outline: none;
  min-width: 180px;
  font-family: 'Rajdhani', sans-serif;
}

.cyber-select:hover {
  border-color: #00f5ff;
  box-shadow: 0 0 15px rgba(0, 245, 255, 0.3);
}

.cyber-select:focus {
  border-color: #00f5ff;
  box-shadow: 0 0 20px rgba(0, 245, 255, 0.5);
}

/* 仓库下拉框 */
.repo-dropdown {
  position: relative;
}

.repo-dropdown-btn {
  background: rgba(10, 14, 39, 0.8);
  border: 2px solid rgba(0, 245, 255, 0.3);
  border-radius: 8px;
  padding: 0.6rem 1rem;
  color: #e2e8f0;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  outline: none;
  min-width: 180px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  font-family: 'Rajdhani', sans-serif;
}

.repo-dropdown-btn:hover {
  border-color: #00f5ff;
  box-shadow: 0 0 15px rgba(0, 245, 255, 0.3);
}

.repo-dropdown-btn.active {
  border-color: #00f5ff;
  box-shadow: 0 0 20px rgba(0, 245, 255, 0.5);
}

.btn-text {
  flex: 1;
  text-align: left;
}

.dropdown-icon {
  width: 16px;
  height: 16px;
  transition: transform 0.3s ease;
  color: #00f5ff;
}

.dropdown-icon.rotated {
  transform: rotate(180deg);
}

.repo-dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  right: 0;
  background: rgba(10, 14, 39, 0.98);
  backdrop-filter: blur(20px);
  border: 2px solid rgba(0, 245, 255, 0.3);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.6), 0 0 30px rgba(0, 245, 255, 0.2);
  z-index: 1000;
  overflow: hidden;
  animation: dropdownSlide 0.2s ease;
}

@keyframes dropdownSlide {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid rgba(0, 245, 255, 0.1);
  background: rgba(0, 245, 255, 0.05);
}

.select-all-btn {
  background: transparent;
  border: 1px solid rgba(0, 245, 255, 0.4);
  border-radius: 6px;
  padding: 0.35rem 0.75rem;
  color: #00f5ff;
  font-size: 0.8rem;
  cursor: pointer;
  transition: all 0.2s ease;
  font-family: 'Rajdhani', sans-serif;
}

.select-all-btn:hover {
  background: rgba(0, 245, 255, 0.1);
  border-color: #00f5ff;
}

.selected-count {
  font-size: 0.8rem;
  color: #64748b;
  font-family: 'Orbitron', sans-serif;
}

.dropdown-list {
  max-height: 300px;
  overflow-y: auto;
  padding: 0.5rem;
}

.dropdown-list::-webkit-scrollbar {
  width: 6px;
}

.dropdown-list::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.2);
}

.dropdown-list::-webkit-scrollbar-thumb {
  background: rgba(0, 245, 255, 0.3);
  border-radius: 3px;
}

.dropdown-list::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 245, 255, 0.5);
}

.repo-option {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.6rem 0.75rem;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: #94a3b8;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
  font-family: 'Rajdhani', sans-serif;
}

.repo-option:hover {
  background: rgba(0, 245, 255, 0.1);
  color: #e2e8f0;
}

.repo-option.active {
  background: rgba(0, 245, 255, 0.15);
  color: #00f5ff;
}

.option-checkbox {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(148, 163, 184, 0.4);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.repo-option.active .option-checkbox {
  border-color: #00f5ff;
  background: rgba(0, 245, 255, 0.2);
}

.option-checkbox svg {
  width: 12px;
  height: 12px;
  color: #00f5ff;
}

.option-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 图表网格 */
.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(600px, 1fr));
  gap: 1.5rem;
}

.chart-card {
  min-height: 500px;
}

@media (max-width: 1200px) {
  .title-row {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .controls-inline {
    width: 100%;
    flex-direction: column;
  }
  
  .charts-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .analytics {
    padding: 1rem;
  }
  
  .page-title {
    font-size: 1.5rem;
  }
}
</style>
