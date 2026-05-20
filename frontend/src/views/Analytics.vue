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
            <div class="time-selector-cyber">
              <div class="quick-options-cyber">
                <button 
                  v-for="option in timeOptions" 
                  :key="option.value"
                  @click="selectTimeRange(option.value)"
                  class="cyber-time-btn"
                  :class="{ active: selectedTimeRange === option.value }"
                >
                  <span class="btn-glow"></span>
                  <span class="btn-text-cyber">{{ option.label }}</span>
                </button>
              </div>
              
              <div class="custom-range-cyber" v-if="selectedTimeRange === 'custom' || (customStartDate && customEndDate)">
                <div class="date-input-wrapper">
                  <input 
                    type="date" 
                    v-model="customStartDate"
                    @change="handleCustomDateChange"
                    class="cyber-date-input"
                  />
                  <div class="input-glow"></div>
                </div>
                <div class="date-connector">
                  <div class="connector-line"></div>
                  <span class="connector-icon">→</span>
                  <div class="connector-line"></div>
                </div>
                <div class="date-input-wrapper">
                  <input 
                    type="date" 
                    v-model="customEndDate"
                    @change="handleCustomDateChange"
                    class="cyber-date-input"
                  />
                  <div class="input-glow"></div>
                </div>
              </div>
            </div>
          </div>

          <div class="control-item">
            <label class="control-label">仓库筛选</label>
            <div class="repo-dropdown">
              <button 
                @click.stop="showRepoDropdown = !showRepoDropdown"
                class="repo-dropdown-btn"
                :class="{ active: showRepoDropdown, 'has-selection': selectedRepos.length > 0 }"
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
            <button 
              @click="loadData" 
              :disabled="loading || selectedRepos.length === 0" 
              class="analyze-btn"
              :class="{ 'can-analyze': selectedRepos.length > 0 }"
            >
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

    <!-- 概览统计卡片 -->
    <div v-if="overviewStats" class="overview-cards">
      <div class="stat-card">
        <div class="stat-value">{{ overviewStats.totalCommits }}</div>
        <div class="stat-label">总提交数</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ overviewStats.totalAdditions }}</div>
        <div class="stat-label">新增行数</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ overviewStats.totalDeletions }}</div>
        <div class="stat-label">删除行数</div>
      </div>
      <div class="stat-card">
        <div class="stat-value">{{ overviewStats.activeAuthors }}</div>
        <div class="stat-label">活跃作者</div>
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
import { getDailyStats, getRepositories, getOverviewStats } from '../api'

const loading = ref(false)
const overviewStats = ref(null)
const selectedTimeRange = ref('week')
const customStartDate = ref('')
const customEndDate = ref('')
const selectedRepos = ref([])
const repositories = ref([])
const dailyStats = ref([])
const showRepoDropdown = ref(false)

// 时间选项配置
const timeOptions = [
  { label: '本周', value: 'week' },
  { label: '上周', value: 'lastWeek' },
  { label: '本月', value: 'month' },
  { label: '本年', value: 'year' }
]

// 选择时间范围
const selectTimeRange = (value) => {
  selectedTimeRange.value = value
  customStartDate.value = ''
  customEndDate.value = ''
}

// 处理自定义日期变化
const handleCustomDateChange = () => {
  if (customStartDate.value && customEndDate.value) {
    selectedTimeRange.value = 'custom'
  }
}

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
  const result = selectedRepos.value.length === 0
    ? dailyStats.value
    : dailyStats.value.filter(stat => selectedRepos.value.includes(stat.repoPath))
  
  console.log('[filteredStats]', {
    selectedRepos: selectedRepos.value,
    dailyStatsCount: dailyStats.value.length,
    filteredCount: result.length
  })
  
  return result
})

// 提取所有日期
const allDates = computed(() => {
  const dateSet = new Set()
  filteredStats.value.forEach(repo => {
    repo.authors?.forEach(author => {
      author.dailyData?.forEach(day => dateSet.add(day.date))
    })
  })
  const result = Array.from(dateSet).sort()
  console.log('[allDates]', result)
  return result
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
  const result = Array.from(authorMap.values())
  console.log('[allAuthors]', result.map(a => a.name))
  return result
})

// 提交趋势图配置
const commitTrendOption = computed(() => {
  const dates = allDates.value
  const authors = allAuthors.value
  
  console.log('[commitTrendOption]', { datesCount: dates.length, authorsCount: authors.length })
  
  if (dates.length === 0 || authors.length === 0) {
    console.warn('[commitTrendOption] No data - returning empty state')
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

  // 为每个用户创建一个系列（新增为正，删除为负）
  const series = []
  const colors = ['#00f5ff', '#ff00ff', '#ffd700', '#00ff88', '#ff6b6b', '#a78bfa']
  
  authors.forEach((author, idx) => {
    const baseColor = colors[idx % colors.length]
    
    // 计算净变化：新增为正，删除为负
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
      // 返回净变化（新增-删除），但为了可视化，我们分别显示
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
      name: `${author.name} - 删除`,
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
      backgroundColor: 'rgba(10, 14, 39, 0.95)',
      borderColor: '#00f5ff',
      textStyle: { color: '#fff' },
      axisPointer: {
        type: 'shadow',
        shadowStyle: { color: 'rgba(0, 245, 255, 0.1)' }
      },
      formatter: (params) => {
        let result = `<div style="font-weight:bold;margin-bottom:5px">${params[0].axisValue}</div>`
        params.forEach(param => {
          const value = Math.abs(param.value)
          const type = param.value >= 0 ? '新增' : '删除'
          result += `<div>${param.seriesName}: <span style="color:${param.color}">${value}</span> (${type})</div>`
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
  // 验证必选条件
  if (selectedRepos.value.length === 0) {
    return
  }
  
  loading.value = true
  try {
    // 构建时间参数
    let startDate, endDate
    const now = new Date()
    
    switch (selectedTimeRange.value) {
      case 'today': // 今日
        startDate = now.toISOString().split('T')[0]
        endDate = now.toISOString().split('T')[0]
        break
      case 'week': // 本周（从周一开始）
        const currentDay = now.getDay() || 7
        const monday = new Date(now)
        monday.setDate(now.getDate() - currentDay + 1)
        startDate = monday.toISOString().split('T')[0]
        endDate = now.toISOString().split('T')[0]
        break
      case 'lastWeek': // 上周
        const lastWeekMonday = new Date(now)
        lastWeekMonday.setDate(now.getDate() - (now.getDay() || 7) - 6)
        const lastWeekSunday = new Date(lastWeekMonday)
        lastWeekSunday.setDate(lastWeekMonday.getDate() + 6)
        startDate = lastWeekMonday.toISOString().split('T')[0]
        endDate = lastWeekSunday.toISOString().split('T')[0]
        break
      case 'month': // 本月
        const firstDay = new Date(now.getFullYear(), now.getMonth(), 1)
        startDate = firstDay.toISOString().split('T')[0]
        endDate = now.toISOString().split('T')[0]
        break
      case 'year': // 本年
        const firstDayOfYear = new Date(now.getFullYear(), 0, 1)
        startDate = firstDayOfYear.toISOString().split('T')[0]
        endDate = now.toISOString().split('T')[0]
        break
      case 'custom': // 自定义
        startDate = customStartDate.value
        endDate = customEndDate.value
        break
      default:
        startDate = null
        endDate = null
    }
    
    // TODO: 在此处增加检查逻辑，如果缓存不足，先调用 performIncrementalScan
    
    // 并行请求 overview 和 daily
    const [overview, stats] = await Promise.all([
      getOverviewStats(startDate, endDate),
      getDailyStats(null, selectedTimeRange.value, selectedRepos.value, startDate, endDate)
    ])
    overviewStats.value = overview
    console.log('[loadData] daily API response:', stats)
    dailyStats.value = stats || []
    console.log('[loadData] dailyStats assigned, length:', dailyStats.value.length)
    
    // 等待 computed 重新计算
    await new Promise(resolve => setTimeout(resolve, 50))
    
    console.log('[loadData] After assignment:')
    console.log('  - filteredStats.length:', filteredStats.value.length)
    console.log('  - allDates:', allDates.value)
    console.log('  - allAuthors:', allAuthors.value.map(a => a.name))
    
    // 先关闭 loading，让图表容器显示
    loading.value = false
  } catch (error) {
    console.error('加载数据失败:', error)
    loading.value = false
  }
}

onMounted(async () => {
  // 只加载仓库列表，不自动选择，不自动加载数据
  try {
    repositories.value = await getRepositories()
  } catch (error) {
    console.error('加载仓库失败:', error)
  }
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

/* 概览统计卡片 */
.overview-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: rgba(10, 14, 39, 0.6);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(0, 212, 255, 0.2);
  border-radius: 12px;
  padding: 1.5rem;
  text-align: center;
  transition: all 0.3s;
}

.stat-card:hover {
  border-color: rgba(0, 212, 255, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 8px 32px rgba(0, 212, 255, 0.15);
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 700;
  background: linear-gradient(135deg, #00d4ff, #7800ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 0.5rem;
}

.stat-label {
  color: #a0aec0;
  font-size: 0.9rem;
  letter-spacing: 1px;
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
  opacity: 0.5;
  cursor: not-allowed;
}

.analyze-btn.can-analyze {
  opacity: 1;
  cursor: pointer;
}

.analyze-btn.can-analyze:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 25px rgba(0, 245, 255, 0.6);
}

.analyze-btn:disabled {
  opacity: 0.5;
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

/* 时间选择器 - 赛博朋克风格 */
.time-selector-cyber {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.quick-options-cyber {
  display: flex;
  gap: 0.6rem;
}

.cyber-time-btn {
  position: relative;
  background: rgba(10, 14, 39, 0.6);
  border: 1px solid rgba(0, 245, 255, 0.2);
  border-radius: 4px;
  padding: 0.6rem 1.2rem;
  color: #94a3b8;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  outline: none;
  font-family: 'Rajdhani', sans-serif;
  font-weight: 600;
  letter-spacing: 1px;
  text-transform: uppercase;
  overflow: hidden;
  min-width: 80px;
}

.cyber-time-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(0, 245, 255, 0.1), transparent);
  transition: left 0.5s ease;
}

.cyber-time-btn:hover::before {
  left: 100%;
}

.cyber-time-btn:hover {
  border-color: rgba(0, 245, 255, 0.5);
  color: #e2e8f0;
  box-shadow: 0 0 20px rgba(0, 245, 255, 0.2), inset 0 0 20px rgba(0, 245, 255, 0.05);
  transform: translateY(-2px);
}

.cyber-time-btn.active {
  background: rgba(0, 245, 255, 0.15);
  border-color: #00f5ff;
  color: #00f5ff;
  box-shadow: 
    0 0 30px rgba(0, 245, 255, 0.4),
    inset 0 0 30px rgba(0, 245, 255, 0.1),
    0 0 60px rgba(0, 245, 255, 0.2);
  text-shadow: 0 0 10px rgba(0, 245, 255, 0.8);
}

.cyber-time-btn.active .btn-glow {
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  background: linear-gradient(45deg, #00f5ff, #ff00ff, #00f5ff);
  border-radius: 4px;
  opacity: 0.3;
  filter: blur(8px);
  animation: glow-pulse 2s ease-in-out infinite;
  z-index: -1;
}

@keyframes glow-pulse {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 0.6; }
}

.btn-text-cyber {
  position: relative;
  z-index: 1;
}

/* 自定义日期范围 - 赛博朋克风格 */
.custom-range-cyber {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem;
  background: rgba(10, 14, 39, 0.4);
  border: 1px solid rgba(0, 245, 255, 0.15);
  border-radius: 8px;
  animation: slideDown 0.3s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.date-input-wrapper {
  position: relative;
  flex: 1;
}

.cyber-date-input {
  width: 100%;
  background: rgba(10, 14, 39, 0.8);
  border: 1px solid rgba(0, 245, 255, 0.3);
  border-radius: 4px;
  padding: 0.6rem 0.9rem;
  color: #e2e8f0;
  font-size: 0.9rem;
  outline: none;
  font-family: 'Rajdhani', sans-serif;
  transition: all 0.3s ease;
  position: relative;
  z-index: 1;
}

.cyber-date-input::-webkit-calendar-picker-indicator {
  filter: invert(1) hue-rotate(180deg);
  cursor: pointer;
  opacity: 0.6;
  transition: opacity 0.3s ease;
}

.cyber-date-input::-webkit-calendar-picker-indicator:hover {
  opacity: 1;
}

.cyber-date-input:focus {
  border-color: #00f5ff;
  box-shadow: 
    0 0 20px rgba(0, 245, 255, 0.3),
    inset 0 0 20px rgba(0, 245, 255, 0.05);
}

.input-glow {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 4px;
  background: radial-gradient(circle at center, rgba(0, 245, 255, 0.1), transparent 70%);
  opacity: 0;
  transition: opacity 0.3s ease;
  pointer-events: none;
}

.date-input-wrapper:focus-within .input-glow {
  opacity: 1;
}

/* 日期连接器 */
.date-connector {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: #00f5ff;
}

.connector-line {
  width: 20px;
  height: 1px;
  background: linear-gradient(90deg, transparent, #00f5ff, transparent);
  animation: line-flow 2s ease-in-out infinite;
}

@keyframes line-flow {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 1; }
}

.connector-icon {
  font-size: 1.2rem;
  color: #00f5ff;
  text-shadow: 0 0 10px rgba(0, 245, 255, 0.8);
  animation: arrow-pulse 1.5s ease-in-out infinite;
}

@keyframes arrow-pulse {
  0%, 100% { transform: translateX(0); opacity: 0.6; }
  50% { transform: translateX(3px); opacity: 1; }
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
  height: 500px;
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
