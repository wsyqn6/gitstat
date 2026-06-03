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
              
              <div class="custom-range-cyber">
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
                    {{ selectedRepos.length === 0 ? t('analytics.selectRepo') : 
                       selectedRepos.length === repositories.length ? t('analytics.allRepos') : 
                       `${selectedRepos.length} ${t('analytics.reposSelected')}` }}
                  </span>
                <svg class="dropdown-icon" :class="{ rotated: showRepoDropdown }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="6 9 12 15 18 9"></polyline>
                </svg>
              </button>
              
              <div v-show="showRepoDropdown" class="repo-dropdown-menu">
                <div class="dropdown-header">
                  <button @click.stop="toggleAllRepos" class="select-all-btn">
                    {{ allReposSelected ? t('analytics.cancelSelectAll') : t('analytics.selectAll') }}
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
              <span>{{ loading ? t('analytics.analyzing') : t('analytics.startAnalyze') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 视图切换 -->
    <div class="view-toggle-bar">
      <div class="view-toggle-inner">
        <button
          @click="viewMode = 'chart'"
          class="view-toggle-btn"
          :class="{ active: viewMode === 'chart' }"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="toggle-icon">
            <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline>
          </svg>
          <span>{{ t('calendar.chartView') }}</span>
        </button>
        <button
          @click="viewMode = 'calendar'"
          class="view-toggle-btn"
          :class="{ active: viewMode === 'calendar' }"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="toggle-icon">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
            <line x1="16" y1="2" x2="16" y2="6"></line>
            <line x1="8" y1="2" x2="8" y2="6"></line>
            <line x1="3" y1="10" x2="21" y2="10"></line>
          </svg>
          <span>{{ t('calendar.calendarView') }}</span>
        </button>
      </div>
    </div>

    <!-- 概览统计卡片 -->
    <div v-if="overviewStats" class="overview-section">
      <div class="overview-period-label">{{ timePeriodLabel }}{{ t('analytics.overviewTitle') }}</div>
      <div class="overview-cards">
      <div class="stat-card">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('analytics.totalCommits') }}</div>
          <div class="stat-value">{{ overviewStats.totalCommits }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="19" x2="12" y2="5"></line>
            <polyline points="5 12 12 5 19 12"></polyline>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('analytics.totalAdditions') }}</div>
          <div class="stat-value">{{ overviewStats.totalAdditions }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <polyline points="19 12 12 19 5 12"></polyline>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('analytics.totalDeletions') }}</div>
          <div class="stat-value">{{ overviewStats.totalDeletions }}</div>
        </div>
      </div>
      <div class="stat-card clickable" :class="{ expanded: expandedSection === 'authors' }" @click="toggleSection('authors')">
        <div class="stat-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
            <circle cx="9" cy="7" r="4"></circle>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ t('analytics.activeAuthors') }}</div>
          <div class="stat-value">{{ overviewStats.activeAuthors }}</div>
        </div>
      </div>
    </div>

      <!-- 活跃作者展开面板 -->
      <div v-if="expandedSection === 'authors' && overviewStats.authors" class="expand-panel">
        <div class="expand-panel-header">{{ timePeriodPrefix }}活跃作者 · 共 {{ overviewStats.authors.length }} 人</div>
        <table class="expand-table">
          <thead>
            <tr>
              <th>作者</th>
              <th>提交数</th>
              <th>新增</th>
              <th>删除</th>
              <th>净变更</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="author in overviewStats.authors" :key="author.email">
              <td class="cell-author">{{ author.author }}</td>
              <td>{{ author.commits }}</td>
              <td class="cell-additions">+{{ author.additions }}</td>
              <td class="cell-deletions">-{{ author.deletions }}</td>
              <td :class="author.netChange >= 0 ? 'cell-additions' : 'cell-deletions'">{{ author.netChange >= 0 ? '+' : '' }}{{ author.netChange }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 洞察卡片 -->
    <div v-if="insights.length > 0" class="insights-grid">
      <div v-for="(insight, idx) in insights" :key="idx" class="insight-card" :class="{ clickable: insight.clickable, expanded: insight.clickable && expandedSection === insight.section }" @click="insight.clickable && toggleSection(insight.section)">
        <div class="insight-icon" v-html="insight.iconSvg"></div>
        <div class="insight-content">
          <div class="insight-title">{{ insight.title }}</div>
          <div class="insight-value">{{ insight.value }}</div>
          <div class="insight-desc">{{ insight.description }}</div>
        </div>
      </div>
    </div>

    <!-- 活跃仓库展开面板 -->
    <div v-if="expandedSection === 'repos' && repoComparison.length > 0" class="expand-panel">
      <div class="expand-panel-header">{{ timePeriodPrefix }}活跃仓库 · 共 {{ repoComparison.filter(r => r.commits > 0).length }} 个</div>
      <table class="expand-table">
        <thead>
          <tr>
            <th>仓库名</th>
            <th>提交数</th>
            <th>作者数</th>
            <th>新增行数</th>
            <th>活跃天数</th>
            <th>日均提交</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="repo in repoComparison" :key="repo.repoPath">
            <td class="cell-author">{{ repo.repoName }}</td>
            <td>{{ repo.commits }}</td>
            <td>{{ repo.authors }}</td>
            <td class="cell-additions">+{{ repo.additions }}</td>
            <td>{{ repo.activeDays }}</td>
            <td>{{ repo.avgCommitsPerDay }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 图表区域 -->
    <div v-if="viewMode === 'chart'" class="charts-grid">
      <!-- 提交趋势图 -->
      <ChartContainer 
        :title="commitTrendTitle" 
        :subtitle="t('analytics.commitTrendSub')"
        :option="commitTrendOption" 
        :loading="loading"
        class="chart-card primary"
      />
      
      <!-- 代码变更图 -->
      <ChartContainer 
        :title="codeChangeTitle" 
        :subtitle="t('analytics.codeChangeSub')"
        :option="codeChangeOption" 
        :loading="loading"
        class="chart-card secondary"
      />

      <!-- 开发者排行榜 -->
      <ChartContainer 
        title="开发者贡献榜" 
        subtitle="Top 10 按提交数排序"
        :option="authorRankOption" 
        :loading="loading"
        class="chart-card tertiary"
      />

      <!-- 活动热力图 -->
      <ChartContainer 
        title="提交时间热力图" 
        subtitle="星期 × 小时维度"
        :option="heatmapOption" 
        :loading="loading"
        class="chart-card quaternary"
      />

      <!-- 仓库对比雷达图 -->
      <ChartContainer 
        title="仓库活跃度对比" 
        subtitle="多维度雷达分析"
        :option="repoComparisonOption" 
        :loading="loading"
        class="chart-card quinary"
      />
    </div>

    <!-- 日历表格视图 -->
    <CalendarView
      v-else
      :viewType="calendarViewType"
      :dailyStats="dailyStats"
      :periodStats="periodStats"
      :startDate="currentStartDate"
      :endDate="currentEndDate"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useI18n } from '../i18n'
import echarts from '../utils/echarts'
import ChartContainer from '../components/ChartContainer.vue'
import CalendarView from '../components/CalendarView.vue'
import { getDailyStats, getWeeklyStats, getMonthlyStats, getYearlyStats, getRepositories, getOverviewStats, getAuthorRank, getActivityHeatmap, getRepoComparison } from '../api'

const { t } = useI18n()
const loading = ref(false)
const overviewStats = ref(null)
const selectedTimeRange = ref('week')
const customStartDate = ref('')
const customEndDate = ref('')
const selectedRepos = ref([])
const repositories = ref([])
const dailyStats = ref([])
const periodStats = ref([])
const currentGranularity = ref('day')
const showRepoDropdown = ref(false)
const authorRank = ref([])
const activityHeatmap = ref([])
const repoComparison = ref([])
const expandedSection = ref(null)
const viewMode = ref('chart')
const currentStartDate = ref('')
const currentEndDate = ref('')
const loadedTimeRange = ref('')

function toLocalDateStr(d) {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

const toggleSection = (section) => {
  expandedSection.value = expandedSection.value === section ? null : section
}

const calendarViewType = computed(() => {
  if (!currentStartDate.value || !currentEndDate.value) return 'week'
  if (selectedTimeRange.value === 'year') return 'year'
  if (selectedTimeRange.value === 'month' || selectedTimeRange.value === 'lastMonth') return 'month'
  if (selectedTimeRange.value === 'week' || selectedTimeRange.value === 'lastWeek' || selectedTimeRange.value === 'today') return 'week'
  const days = Math.round((new Date(currentEndDate.value) - new Date(currentStartDate.value)) / (1000 * 60 * 60 * 24))
  if (days <= 7) return 'week'
  if (days <= 31) return 'month'
  return 'year'
})

// 切到日历且粒度为周时重取日粒度数据
watch(viewMode, (newMode) => {
  if (!currentStartDate.value || loading.value) return
  if (newMode === 'calendar' && currentGranularity.value !== 'day') {
    loadData()
  }
})

// 时间选项配置
const timeOptions = [
  { label: computed(() => t('analytics.thisWeek')), value: 'week' },
  { label: computed(() => t('analytics.thisMonth')), value: 'month' },
  { label: computed(() => t('analytics.thisYear')), value: 'year' }
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

// 周期标签格式化
const formatPeriodLabel = (period, granularity) => {
  if (granularity === 'week') {
    const match = period.match(/-W(\d+)/)
    return match ? `第${match[1]}周` : period
  }
  if (granularity === 'month') {
    return period.slice(5) + '月'
  }
  if (granularity === 'year') {
    return period + '年'
  }
  return period
}

// 根据当前粒度选择数据源
const activeStats = computed(() => {
  return currentGranularity.value === 'day' ? dailyStats.value : periodStats.value
})

// 获取过滤后的统计数据
const filteredStats = computed(() => {
  const source = activeStats.value
  const result = selectedRepos.value.length === 0
    ? source
    : source.filter(stat => selectedRepos.value.includes(stat.repoPath))
  
  console.log('[filteredStats]', {
    selectedRepos: selectedRepos.value,
    statsCount: source.length,
    filteredCount: result.length,
    granularity: currentGranularity.value
  })
  
  return result
})

// 提取所有日期 (原始key)
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

// 动态图表标题
const granularityPrefix = computed(() => {
  switch (currentGranularity.value) {
    case 'week': return t('analytics.weekly')
    case 'month': return t('analytics.monthly')
    case 'year': return t('analytics.yearly')
    default: return t('analytics.daily')
  }
})

const commitTrendTitle = computed(() => granularityPrefix.value + t('analytics.commitTrend'))
const codeChangeTitle = computed(() => granularityPrefix.value + t('analytics.codeChange'))

// 时间周期文本标签（用于概览和洞察）
const timePeriodLabel = computed(() => {
  switch (loadedTimeRange.value) {
    case 'week': return t('analytics.thisWeek')
    case 'lastWeek': return t('analytics.lastWeek')
    case 'month': return t('analytics.thisMonth')
    case 'lastMonth': return t('analytics.lastMonth')
    case 'year': return t('analytics.thisYear')
    case 'custom': return currentStartDate.value && currentEndDate.value
      ? `${currentStartDate.value} ~ ${currentEndDate.value}`
      : t('analytics.customPeriod')
    default: return ''
  }
})

const timePeriodPrefix = computed(() => {
  switch (loadedTimeRange.value) {
    case 'week': return '本周'
    case 'lastWeek': return '上周'
    case 'month': return '本月'
    case 'lastMonth': return '上月'
    case 'year': return '本年'
    case 'custom': return '统计周期'
    default: return ''
  }
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
      right: '8%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dates,
      axisLine: { lineStyle: { color: '#334155' } },
      axisLabel: { 
        color: '#94a3b8',
        formatter: (value) => formatPeriodLabel(value, currentGranularity.value)
      },
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
      right: '8%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: { lineStyle: { color: '#334155' } },
      axisLabel: { 
        color: '#94a3b8',
        formatter: (value) => formatPeriodLabel(value, currentGranularity.value)
      },
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

// 开发者排行榜图配置
const authorRankOption = computed(() => {
  if (!authorRank.value || authorRank.value.length === 0) {
    return {
      title: { 
        text: '暂无数据', 
        left: 'center', 
        top: 'center',
        textStyle: { color: '#64748b', fontSize: 16 }
      }
    }
  }

  const topAuthors = authorRank.value.slice(0, 10)
  const colors = ['#00f5ff', '#ff00ff', '#ffd700', '#00ff88', '#ff6b6b', '#a78bfa']

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(10, 14, 39, 0.95)',
      borderColor: '#00f5ff',
      textStyle: { color: '#fff' },
      axisPointer: { type: 'shadow' },
      formatter: (params) => {
        const item = params[0]
        const author = topAuthors[item.dataIndex]
        return `
          <div style="font-weight:bold;margin-bottom:5px">${author.author}</div>
          <div>提交数: ${author.commits}</div>
          <div>新增: <span style="color:#00ff88">${author.additions}</span></div>
          <div>删除: <span style="color:#ff6b6b">${author.deletions}</span></div>
          <div>净变化: ${author.netChange > 0 ? '+' : ''}${author.netChange}</div>
        `
      }
    },
    grid: {
      left: '3%',
      right: '8%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
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
    yAxis: {
      type: 'category',
      data: topAuthors.map(a => a.author),
      axisLine: { lineStyle: { color: '#334155' } },
      axisLabel: { color: '#94a3b8' },
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

// 活动热力图配置
const heatmapOption = computed(() => {
  if (!activityHeatmap.value || activityHeatmap.value.length === 0) {
    return {
      title: { 
        text: '暂无数据', 
        left: 'center', 
        top: 'center',
        textStyle: { color: '#64748b', fontSize: 16 }
      }
    }
  }

  const hours = Array.from({ length: 24 }, (_, i) => `${i}:00`)
  const days = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  
  // 构建热力图数据
  const data = activityHeatmap.value.map(item => [
    item.dayOfWeek,
    item.hour,
    item.commitCount
  ])

  const maxCount = Math.max(...data.map(d => d[2]))

  return {
    tooltip: {
      backgroundColor: 'rgba(10, 14, 39, 0.95)',
      borderColor: '#00f5ff',
      textStyle: { color: '#fff' },
      formatter: (params) => {
        return `<div>${days[params.value[0]]} ${hours[params.value[1]]}</div><div>提交数: ${params.value[2]}</div>`
      }
    },
    grid: {
      left: '3%',
      right: '8%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: days,
      axisLine: { lineStyle: { color: '#334155' } },
      axisLabel: { color: '#94a3b8' },
      splitArea: { show: true, areaStyle: { color: 'rgba(30, 41, 59, 0.3)' } }
    },
    yAxis: {
      type: 'category',
      data: hours,
      axisLine: { lineStyle: { color: '#334155' } },
      axisLabel: { color: '#94a3b8' },
      splitArea: { show: true, areaStyle: { color: 'rgba(30, 41, 59, 0.3)' } }
    },
    visualMap: {
      min: 0,
      max: maxCount,
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      bottom: '0%',
      inRange: {
        color: [
          'rgba(0, 245, 255, 0.1)',
          'rgba(0, 245, 255, 0.4)',
          'rgba(0, 245, 255, 0.7)',
          '#00f5ff'
        ]
      },
      textStyle: { color: '#94a3b8' }
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

// 仓库对比雷达图配置
const repoComparisonOption = computed(() => {
  if (!repoComparison.value || repoComparison.value.length === 0) {
    return {
      title: { 
        text: '暂无数据', 
        left: 'center', 
        top: 'center',
        textStyle: { color: '#64748b', fontSize: 16 }
      }
    }
  }

  const repos = repoComparison.value.slice(0, 5)
  const colors = ['#00f5ff', '#ff00ff', '#ffd700', '#00ff88', '#ff6b6b']

  // 计算最大值用于归一化
  const maxCommits = Math.max(...repos.map(r => r.commits))
  const maxAdditions = Math.max(...repos.map(r => r.additions))
  const maxAuthors = Math.max(...repos.map(r => r.authors))
  const maxActiveDays = Math.max(...repos.map(r => r.activeDays))

  return {
    tooltip: {
      backgroundColor: 'rgba(10, 14, 39, 0.95)',
      borderColor: '#00f5ff',
      textStyle: { color: '#fff' }
    },
    legend: {
      data: repos.map(r => r.repoName),
      textStyle: { color: '#94a3b8' },
      top: 10
    },
    radar: {
      indicator: [
        { name: '提交数', max: 100 },
        { name: '代码量', max: 100 },
        { name: '作者数', max: 100 },
        { name: '活跃天数', max: 100 },
        { name: '日均提交', max: 100 }
      ],
      shape: 'polygon',
      splitNumber: 4,
      axisName: { color: '#94a3b8' },
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

// 洞察卡片数据
const insights = computed(() => {
  const result = []
  const prefix = timePeriodPrefix.value
  
  // 最活跃开发者
  if (authorRank.value && authorRank.value.length > 0) {
    const topAuthor = authorRank.value[0]
    result.push({
      iconSvg: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>`,
      title: `${prefix}之星`,
      value: topAuthor.author,
      description: `${topAuthor.commits} 次提交 · ${topAuthor.additions} 行新增`
    })
  }
  
  // 平均提交大小
  if (overviewStats.value && overviewStats.value.totalCommits > 0) {
    const avgSize = Math.round((overviewStats.value.totalAdditions + overviewStats.value.totalDeletions) / overviewStats.value.totalCommits)
    result.push({
      iconSvg: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>`,
      title: '平均提交规模',
      value: `${avgSize} 行`,
      description: `${prefix}每次提交平均变更 ${avgSize} 行`
    })
  }
  
  // 代码净增长
  if (overviewStats.value) {
    const netChange = overviewStats.value.totalAdditions - overviewStats.value.totalDeletions
    result.push({
      iconSvg: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 7 13.5 15.5 8.5 10.5 2 17"/><polyline points="16 7 22 7 22 13"/></svg>`,
      title: '代码净增长',
      value: `${netChange > 0 ? '+' : ''}${netChange}`,
      description: `${prefix}代码净增 ${netChange} 行`
    })
  }
  
  // 活跃仓库数
  if (repoComparison.value && repoComparison.value.length > 0) {
    const activeRepos = repoComparison.value.filter(r => r.commits > 0).length
    result.push({
      iconSvg: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>`,
      title: '活跃仓库',
      value: `${activeRepos} 个`,
      description: `${prefix}${activeRepos} 个仓库有提交活动`,
      clickable: true,
      section: 'repos'
    })
  }
  
  return result
})

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
        startDate = toLocalDateStr(now)
        endDate = toLocalDateStr(now)
        break
      case 'week': // 本周（从周一开始）
        const currentDay = now.getDay() || 7
        const monday = new Date(now)
        monday.setDate(now.getDate() - currentDay + 1)
        startDate = toLocalDateStr(monday)
        endDate = toLocalDateStr(now)
        break
      case 'lastWeek': // 上周
        const lastWeekMonday = new Date(now)
        lastWeekMonday.setDate(now.getDate() - (now.getDay() || 7) - 6)
        const lastWeekSunday = new Date(lastWeekMonday)
        lastWeekSunday.setDate(lastWeekMonday.getDate() + 6)
        startDate = toLocalDateStr(lastWeekMonday)
        endDate = toLocalDateStr(lastWeekSunday)
        break
      case 'month': // 本月
        const firstDay = new Date(now.getFullYear(), now.getMonth(), 1)
        startDate = toLocalDateStr(firstDay)
        endDate = toLocalDateStr(now)
        break
      case 'lastMonth': // 上月
        const firstDayLast = new Date(now.getFullYear(), now.getMonth() - 1, 1)
        const lastDayLast = new Date(now.getFullYear(), now.getMonth(), 0)
        startDate = toLocalDateStr(firstDayLast)
        endDate = toLocalDateStr(lastDayLast)
        break
      case 'year': // 本年
        const firstDayOfYear = new Date(now.getFullYear(), 0, 1)
        startDate = toLocalDateStr(firstDayOfYear)
        endDate = toLocalDateStr(now)
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
    
    currentStartDate.value = startDate
    currentEndDate.value = endDate

    // 根据时间范围选择聚合粒度
    let granularity = 'day'
    let statsPromise
    const tr = selectedTimeRange.value

    if (tr === 'today' || tr === 'week' || tr === 'lastWeek') {
      granularity = 'day'
      statsPromise = getDailyStats(null, tr === 'custom' ? '' : tr, selectedRepos.value, startDate, endDate)
    } else if (tr === 'month' || tr === 'lastMonth') {
      if (viewMode.value === 'calendar') {
        granularity = 'day'
        statsPromise = getDailyStats(null, tr, selectedRepos.value, startDate, endDate)
      } else {
        granularity = 'week'
        statsPromise = getWeeklyStats(null, tr, selectedRepos.value, startDate, endDate)
      }
    } else if (tr === 'year') {
      granularity = 'month'
      statsPromise = getMonthlyStats(null, tr, selectedRepos.value, startDate, endDate)
    } else {
      // custom: 根据天数选择粒度
      const days = Math.round((new Date(endDate) - new Date(startDate)) / (1000 * 60 * 60 * 24))
      if (days <= 31) {
        granularity = 'day'
        statsPromise = getDailyStats(null, '', selectedRepos.value, startDate, endDate)
      } else if (days <= 180 && viewMode.value !== 'calendar') {
        granularity = 'week'
        statsPromise = getWeeklyStats(null, '', selectedRepos.value, startDate, endDate)
      } else {
        granularity = 'month'
        statsPromise = getMonthlyStats(null, '', selectedRepos.value, startDate, endDate)
      }
    }
    currentGranularity.value = granularity

    // 并行请求 overview 和 stats
    const [overview, stats, authors, heatmap, comparison] = await Promise.all([
      getOverviewStats(startDate, endDate, selectedRepos.value),
      statsPromise,
      getAuthorRank(selectedRepos.value, startDate, endDate),
      getActivityHeatmap(selectedRepos.value, startDate, endDate),
      getRepoComparison(selectedRepos.value, startDate, endDate)
    ])
    overviewStats.value = overview
    console.log(`[loadData] ${granularity} API response:`, stats)

    if (granularity === 'day') {
      dailyStats.value = stats || []
      periodStats.value = []
    } else {
      periodStats.value = stats || []
      dailyStats.value = []
    }
    authorRank.value = authors || []
    activityHeatmap.value = heatmap || []
    repoComparison.value = comparison || []
    console.log(`[loadData] ${granularity}Stats assigned, length:`, stats?.length || 0)
    
    // 等待 computed 重新计算
    await new Promise(resolve => setTimeout(resolve, 50))
    
    console.log('[loadData] After assignment:')
    console.log('  - granularity:', currentGranularity.value)
    console.log('  - filteredStats.length:', filteredStats.value.length)
    console.log('  - allDates:', allDates.value)
    console.log('  - allAuthors:', allAuthors.value.map(a => a.name))
    
    // 先关闭 loading，让图表容器显示
    loadedTimeRange.value = selectedTimeRange.value
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
.overview-section {
  margin-bottom: 2rem;
}

.overview-period-label {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.85rem;
  color: #00f5ff;
  letter-spacing: 2px;
  text-transform: uppercase;
  margin-bottom: 0.75rem;
  opacity: 0.8;
}

/* 可点击卡片 */
.stat-card.clickable,
.insight-card.clickable {
  cursor: pointer;
}

.insight-card.clickable:hover {
  border-color: rgba(0, 212, 255, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 8px 32px rgba(0, 212, 255, 0.15);
}

.stat-card.clickable.expanded,
.insight-card.clickable.expanded {
  border-color: #00f5ff;
  box-shadow: 0 0 20px rgba(0, 245, 255, 0.3), inset 0 0 20px rgba(0, 245, 255, 0.05);
}

/* 展开面板 */
.expand-panel {
  background: rgba(10, 14, 39, 0.8);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(0, 212, 255, 0.2);
  border-radius: 12px;
  padding: 1.25rem;
  margin-top: 1rem;
  animation: slideDown 0.25s ease;
  overflow-x: auto;
}

.expand-panel-header {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.8rem;
  color: #00f5ff;
  letter-spacing: 1px;
  margin-bottom: 1rem;
  opacity: 0.9;
}

.expand-table {
  width: 100%;
  border-collapse: collapse;
  font-family: 'Rajdhani', sans-serif;
  font-size: 0.9rem;
}

.expand-table th {
  text-align: left;
  padding: 0.5rem 0.75rem;
  color: #64748b;
  font-weight: 600;
  font-size: 0.8rem;
  letter-spacing: 1px;
  text-transform: uppercase;
  border-bottom: 1px solid rgba(0, 245, 255, 0.15);
}

.expand-table td {
  padding: 0.5rem 0.75rem;
  color: #e2e8f0;
  border-bottom: 1px solid rgba(148, 163, 184, 0.1);
}

.expand-table tbody tr:hover {
  background: rgba(0, 245, 255, 0.05);
}

.expand-table .cell-author {
  color: #00f5ff;
  font-weight: 600;
}

.expand-table .cell-additions {
  color: #00ff88;
}

.expand-table .cell-deletions {
  color: #ff6b6b;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.overview-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
}

.stat-card {
  background: rgba(10, 14, 39, 0.6);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(0, 212, 255, 0.2);
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: all 0.3s;
  min-height: 120px;
}

.stat-card:hover {
  border-color: rgba(0, 212, 255, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 8px 32px rgba(0, 212, 255, 0.15);
}

.stat-icon {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #00f5ff;
  filter: drop-shadow(0 0 6px rgba(0, 245, 255, 0.4));
}

.stat-icon svg {
  width: 24px;
  height: 24px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  background: linear-gradient(135deg, #00d4ff, #7800ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.stat-label {
  color: #64748b;
  font-size: 0.85rem;
  letter-spacing: 1px;
  margin-bottom: 0.25rem;
}

/* 洞察卡片 */
.insights-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.insight-card {
  background: rgba(10, 14, 39, 0.6);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(0, 212, 255, 0.2);
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: all 0.3s;
  min-height: 100px;
}

.insight-card:hover {
  border-color: rgba(0, 212, 255, 0.5);
  transform: translateY(-2px);
  box-shadow: 0 8px 32px rgba(0, 212, 255, 0.15);
}

.insight-icon {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #00f5ff;
  filter: drop-shadow(0 0 6px rgba(0, 245, 255, 0.4));
}

.insight-icon svg {
  width: 24px;
  height: 24px;
}

.insight-content {
  flex: 1;
}

.insight-title {
  font-size: 0.85rem;
  color: #64748b;
  letter-spacing: 1px;
  margin-bottom: 0.25rem;
}

.insight-value {
  font-size: 1.5rem;
  font-weight: 700;
  background: linear-gradient(135deg, #00d4ff, #7800ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 0.25rem;
}

.insight-desc {
  font-size: 0.8rem;
  color: #94a3b8;
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
  flex-direction: row;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.6rem;
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
  gap: 0.5rem;
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
}

.cyber-date-input {
  width: 130px;
  background: rgba(10, 14, 39, 0.8);
  border: 1px solid rgba(0, 245, 255, 0.3);
  border-radius: 4px;
  padding: 0.5rem 0.6rem;
  color: #e2e8f0;
  font-size: 0.78rem;
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
  grid-template-columns: repeat(auto-fit, minmax(500px, 1fr));
  gap: 1.5rem;
}

.chart-card {
  min-height: 450px;
  height: 450px;
}

/* 视图切换栏 */
.view-toggle-bar {
  margin-bottom: 1.5rem;
  display: flex;
  justify-content: center;
}
.view-toggle-inner {
  display: flex;
  background: rgba(10, 14, 39, 0.6);
  border: 1px solid rgba(0, 245, 255, 0.2);
  border-radius: 8px;
  padding: 3px;
  gap: 2px;
}
.view-toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0.5rem 1.2rem;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: #64748b;
  font-size: 0.8rem;
  font-family: 'Orbitron', sans-serif;
  letter-spacing: 1px;
  cursor: pointer;
  transition: all 0.3s ease;
  text-transform: uppercase;
}
.view-toggle-btn:hover {
  color: #94a3b8;
}
.view-toggle-btn.active {
  background: rgba(0, 245, 255, 0.15);
  color: #00f5ff;
  box-shadow: 0 0 20px rgba(0, 245, 255, 0.2);
}
.toggle-icon {
  width: 16px;
  height: 16px;
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
