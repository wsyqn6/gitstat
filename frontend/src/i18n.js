import { reactive, computed } from 'vue'

const messages = {
  zh: {
    nav: {
      dashboard: '仪表盘',
      analytics: '数据分析',
      settings: '设置'
    },
    dashboard: {
      switchDirectory: '切换目录',
      todayCommits: '今日提交数',
      todayAdditions: '今日新增行数',
      todayDeletions: '今日删除行数',
      activeAuthors: '活跃作者数',
      todayDetails: '今日提交详情',
      lastCommit: '最后提交',
      author: '作者',
      commits: '提交数',
      changes: '变更',
      me: 'ME'
    },
    scan: {
      title: '扫描Git仓库',
      directoryPath: '目录路径',
      startScan: '开始扫描',
      scanning: '扫描中...'
    },
    analytics: {
      title: '数据分析中心',
      subtitle: '实时监控 · 多维度洞察',
      timeRange: '时间范围',
      thisWeek: '本周',
      lastWeek: '上周',
      thisMonth: '本月',
      thisYear: '本年',
      repoFilter: '仓库筛选',
      selectRepo: '选择仓库',
      allRepos: '全部仓库',
      reposSelected: '个仓库',
      selectAll: '全选',
      cancelSelectAll: '取消全选',
      startAnalyze: '开始分析',
      analyzing: '分析中...',
      totalCommits: '总提交数',
      totalAdditions: '新增行数',
      totalDeletions: '删除行数',
      activeAuthors: '活跃作者',
      commitTrend: '每日提交趋势',
      commitTrendSub: '多用户对比分析',
      codeChange: '代码变更分布',
      codeChangeSub: '新增 vs 删除',
      noData: '暂无数据',
      additions: '新增',
      deletions: '删除'
    },
    settings: {
      title: '设置',
      scanConfig: '扫描配置',
      directoryPath: '目录路径',
      pathPlaceholder: '/path/to/repos',
      startScan: '开始扫描',
      scanning: '扫描中...',
      scanSuccess: '扫描完成',
      dataManagement: '数据管理',
      exportData: '导出数据 (JSON)',
      about: '关于',
      platformName: 'Git提交统计平台'
    }
  },
  en: {
    nav: {
      dashboard: 'Dashboard',
      analytics: 'Analytics',
      settings: 'Settings'
    },
    dashboard: {
      switchDirectory: 'Switch Directory',
      todayCommits: "Today's Commits",
      todayAdditions: "Today's Additions",
      todayDeletions: "Today's Deletions",
      activeAuthors: 'Active Authors',
      todayDetails: "Today's Commit Details",
      lastCommit: 'Last Commit',
      author: 'Author',
      commits: 'Commits',
      changes: 'Changes',
      me: 'ME'
    },
    scan: {
      title: 'Scan Git Repositories',
      directoryPath: 'Directory Path',
      startScan: 'Start Scan',
      scanning: 'Scanning...'
    },
    analytics: {
      title: 'Data Analytics Center',
      subtitle: 'Real-time Monitoring · Multi-dimensional Insights',
      timeRange: 'Time Range',
      thisWeek: 'This Week',
      lastWeek: 'Last Week',
      thisMonth: 'This Month',
      thisYear: 'This Year',
      repoFilter: 'Repository Filter',
      selectRepo: 'Select Repositories',
      allRepos: 'All Repositories',
      reposSelected: 'repositories selected',
      selectAll: 'Select All',
      cancelSelectAll: 'Deselect All',
      startAnalyze: 'Start Analysis',
      analyzing: 'Analyzing...',
      totalCommits: 'Total Commits',
      totalAdditions: 'Additions',
      totalDeletions: 'Deletions',
      activeAuthors: 'Active Authors',
      commitTrend: 'Daily Commit Trend',
      commitTrendSub: 'Multi-user Comparison',
      codeChange: 'Code Change Distribution',
      codeChangeSub: 'Additions vs Deletions',
      noData: 'No Data',
      additions: 'Additions',
      deletions: 'Deletions'
    },
    settings: {
      title: 'Settings',
      scanConfig: 'Scan Configuration',
      directoryPath: 'Directory Path',
      pathPlaceholder: '/path/to/repos',
      startScan: 'Start Scan',
      scanning: 'Scanning...',
      scanSuccess: 'Scan completed',
      dataManagement: 'Data Management',
      exportData: 'Export Data (JSON)',
      about: 'About',
      platformName: 'Git Commit Statistics Platform'
    }
  }
}

const state = reactive({
  locale: getBrowserLocale()
})

function getBrowserLocale() {
  const saved = localStorage.getItem('locale')
  if (saved && (saved === 'zh' || saved === 'en')) return saved
  const lang = navigator.language || navigator.userLanguage
  return lang.startsWith('zh') ? 'zh' : 'en'
}

export function useI18n() {
  const t = (key) => {
    const keys = key.split('.')
    let value = messages[state.locale]
    for (const k of keys) {
      value = value?.[k]
    }
    return value || key
  }

  const setLocale = (locale) => {
    state.locale = locale
    localStorage.setItem('locale', locale)
  }

  return {
    t,
    locale: computed(() => state.locale),
    setLocale
  }
}
