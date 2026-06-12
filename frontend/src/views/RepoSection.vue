<template>
  <div class="repo-page">
    <div v-if="state.reposInfo.length === 0" class="empty-state card">
      <span class="empty-icon">▤</span>
      <p>{{ t('repo.empty') }}</p>
    </div>

    <template v-else>
      <div v-if="state.reposInfo.length > 1" class="tabs">
        <button
          v-for="repo in state.reposInfo"
          :key="repo.path"
          :class="['tab', { active: activePath === repo.path }]"
          @click="switchRepo(repo.path)"
        >
          {{ repo.name }}
        </button>
      </div>

      <div class="content">
          <div v-if="loading" class="loading-state">
            <span class="spinner large"></span>
            <p>{{ t('repo.loading') }}</p>
          </div>

          <div v-else-if="!detail" class="skeleton-page">
            <div class="title-row">
              <div class="skeleton-line w40" style="height:32px;border-radius:6px"></div>
            </div>
            <div class="card-row l1-row">
              <div v-for="i in 4" :key="i" class="info-card">
                <div class="skeleton-circle" style="width:36px;height:36px;margin:0 auto 0.5rem"></div>
                <div class="skeleton-line w35" style="margin:0 auto 0.4rem;height:22px;border-radius:6px"></div>
                <div class="skeleton-line w55" style="margin:0 auto"></div>
              </div>
            </div>
            <div class="card stats-cta">
              <div class="stats-cta-content">
                <div class="skeleton-circle" style="width:40px;height:40px;flex-shrink:0"></div>
                <div style="flex:1;display:flex;flex-direction:column;gap:0.3rem">
                  <div class="skeleton-line w45" style="height:18px;border-radius:6px"></div>
                  <div class="skeleton-line w70"></div>
                </div>
                <div class="skeleton-line w20" style="height:38px;border-radius:8px"></div>
              </div>
            </div>
            <div class="card-row l2-row">
              <div v-for="i in 4" :key="i" class="info-card dim">
                <div class="skeleton-line w40" style="margin:0 auto 0.3rem;height:22px;border-radius:6px"></div>
                <div class="skeleton-line w65" style="margin:0 auto"></div>
              </div>
            </div>
            <div class="section-group">
              <div class="section card">
                <div class="skeleton-line w50" style="height:20px;border-radius:6px;margin-bottom:1rem"></div>
                <div class="lang-body">
                  <div class="skeleton-chart-area" style="height:200px"></div>
                  <div style="display:flex;flex-direction:column;gap:0.5rem;padding-top:0.3rem">
                    <div v-for="i in 4" :key="i" class="skeleton-line w90"></div>
                  </div>
                </div>
              </div>
              <div class="section card">
                <div class="skeleton-line w50" style="height:20px;border-radius:6px;margin-bottom:1rem"></div>
                <div v-for="i in 5" :key="i" class="skeleton-line w85" style="margin-bottom:0.5rem;height:12px"></div>
              </div>
            </div>
          </div>

          <template v-else-if="detail">
          <div class="title-row">
            <h2 class="repo-title">{{ detail.name }}</h2>
            <span class="repo-path">{{ detail.path }}</span>
          </div>

          <!-- L1: Fast data from RepoInfo -->
          <div class="card-row l1-row">
            <div class="info-card">
              <span class="info-icon" style="color:#00d4ff">⑂</span>
              <span class="info-value">{{ detail.currentBranch }}</span>
              <span class="info-label">{{ t('repo.currentBranch') }}</span>
            </div>
            <div class="info-card clickable" :class="{ expanded: expandedBranch }" @click="toggleBranch">
              <span class="info-icon" style="color:#a78bfa">⑂</span>
              <span class="info-value">{{ localBranchCount }}</span>
              <span class="info-label">{{ t('repo.branchCount') }} <span class="expand-icon">{{ expandedBranch ? '▼' : '▶' }}</span></span>
            </div>
            <div class="info-card">
              <span class="info-icon" style="color:#fbbf24">◈</span>
              <span class="info-value">{{ detail.fileCount ?? '--' }}</span>
              <span class="info-label">{{ t('repo.fileCount') }}</span>
            </div>
            <div class="info-card clickable" :class="{ expanded: expandedTags }" @click="toggleTags">
              <span class="info-icon" style="color:#34d399">◉</span>
              <span class="info-value">{{ tagCount }}</span>
              <span class="info-label">{{ t('repo.tags') }} <span class="expand-icon">{{ expandedTags ? '▼' : '▶' }}</span></span>
            </div>
          </div>

          <div v-if="expandedBranch" class="expand-panel card">
            <h4>{{ t('repo.allBranches') }}</h4>
            <div class="branch-groups">
              <div v-if="detail.branchCount > 0" class="branch-group">
                <span class="branch-group-label">{{ t('repo.localBranch') }}</span>
                <div class="branch-tags">
                  <span
                    v-for="b in localBranches"
                    :key="b"
                    class="branch-tag"
                    :class="{ current: b === detail.currentBranch }"
                  >{{ b }}<span v-if="b === detail.currentBranch" class="current-badge">{{ t('repo.currentBranch') }}</span></span>
                </div>
              </div>
              <div v-if="remoteBranchCount > 0" class="branch-group">
                <span class="branch-group-label">{{ t('repo.remoteBranch') }}</span>
                <div class="branch-tags">
                  <span
                    v-for="b in detail.remoteBranches || []"
                    :key="b"
                    class="branch-tag remote"
                  >{{ b }}</span>
                </div>
              </div>
            </div>
          </div>

          <div v-if="expandedTags" class="expand-panel card">
            <h4>{{ t('repo.tags') }}</h4>
            <div v-if="detail.tags && detail.tags.length > 0" class="tag-cloud">
              <span v-for="t in detail.tags" :key="t" class="tag-item">{{ t }}</span>
            </div>
            <p v-else class="expand-hint">--</p>
          </div>

          <div class="stats-group">
            <div class="card-row l2-row">
            <div class="info-card dim" :class="{ 'has-data': statsLoaded && detail.repoSize > 0 }">
              <span class="info-value">
                <template v-if="statsLoaded">{{ detail.repoSize > 0 ? formatBytes(detail.repoSize) : '--' }}</template>
                <span v-else class="skeleton-line w40" style="margin:0 auto;height:22px;border-radius:6px"></span>
              </span>
              <span class="info-label">{{ t('repo.diskSize') }}</span>
            </div>
            <div class="info-card clickable dim" :class="statsLoaded ? { expanded: expandedContributor, 'has-data': hasCommits } : {}" @click="toggleContributor">
              <span class="info-value">
                <template v-if="statsLoaded">{{ hasCommits ? detail.contributors.length : '--' }}</template>
                <span v-else class="skeleton-line w30" style="margin:0 auto;height:22px;border-radius:6px"></span>
              </span>
              <span class="info-label">{{ t('repo.contributors') }} <span class="expand-icon">{{ expandedContributor ? '▼' : '▶' }}</span></span>
            </div>
            <div class="info-card dim" :class="{ 'has-data': statsLoaded && !!detail.earliestCommitAuthor }">
              <span class="info-value">
                <template v-if="statsLoaded">{{ formatDate(detail.earliestDate) }}</template>
                <span v-else class="skeleton-line w55" style="margin:0 auto;height:22px;border-radius:6px"></span>
              </span>
              <span class="info-label">{{ t('repo.createDate') }} · {{ timeSpan }}</span>
            </div>
            <div class="info-card dim" :class="{ 'has-data': statsLoaded && !!detail.lastCommitTime }">
              <span class="info-value">
                <template v-if="statsLoaded">{{ formatTimeAgo(detail.lastCommitTime) }}</template>
                <span v-else class="skeleton-line w45" style="margin:0 auto;height:22px;border-radius:6px"></span>
              </span>
              <span class="info-label">{{ t('repo.lastCommit') }}</span>
            </div>
          </div>

            <template v-if="statsLoaded">
              <div v-if="expandedContributor" class="expand-panel card">
                <h4>{{ t('repo.contributors') }}</h4>
                <div class="contrib-table">
                  <div class="contrib-header">
                    <span>{{ t('repo.author') }}</span>
                    <span>{{ t('repo.commits') }}</span>
                    <span class="add">{{ t('repo.additions') }}</span>
                    <span class="del">{{ t('repo.deletions') }}</span>
                    <span>{{ t('repo.lastCommit') }}</span>
                  </div>
                  <div v-for="ct in detail.contributors" :key="ct.email" class="contrib-row">
                    <span class="contrib-name">{{ ct.author }}</span>
                    <span>{{ ct.commitCount }}</span>
                    <span class="add">+{{ ct.additions }}</span>
                    <span class="del">-{{ ct.deletions }}</span>
                    <span class="contrib-time">{{ ct.lastCommitDate?.slice(0, 10) }}</span>
                  </div>
                </div>
              </div>
              <RepoCharts :data="chartData" :loading="chartLoading" />
            </template>

            <div v-if="!statsLoaded" class="charts-section-skel">
              <div class="section-header">
                <h3 class="section-title">{{ t('repo.commitCalendar') }}</h3>
              </div>
              <div class="skeleton-chart-area" style="height:160px;margin-bottom:2rem;border-radius:12px"></div>
              <div class="chart-row-skel">
                <div class="card" style="flex:1;padding:1rem;min-height:300px">
                  <div class="skeleton-line w40" style="height:20px;border-radius:6px;margin-bottom:1rem"></div>
                  <div class="skeleton-chart-area" style="height:240px"></div>
                </div>
                <div class="card" style="flex:1;padding:1rem;min-height:300px">
                  <div class="skeleton-line w40" style="height:20px;border-radius:6px;margin-bottom:1rem"></div>
                  <div class="skeleton-chart-area" style="height:240px"></div>
                </div>
              </div>
            </div>

            <div v-if="!statsLoaded" class="stats-group-overlay" @click="loadStats">
              <div class="stats-group-overlay-content">
                <span class="stats-cta-icon">▦</span>
                <h4>{{ t('repo.statsTitle') }}</h4>
                <p>{{ t('repo.statsDesc') }}</p>
                <button class="btn" :disabled="loadingStats" @click.stop="loadStats">
                  <span v-if="loadingStats" class="spinner"></span>
                  {{ loadingStats ? t('repo.statsLoading') : t('repo.statsBtn') }}
                </button>
              </div>
            </div>
          </div>

          <div class="section-group">
            <div class="section card">
              <h3 class="section-title">{{ t('repo.langDistribution') }}</h3>
              <div v-if="analysisLoading" class="section-loading">
                <span class="spinner"></span>
                <span>{{ t('repo.analyzing') }}</span>
              </div>
              <div v-else-if="detail.analysis" class="lang-body">
                <div ref="langChartRef" id="lang-chart-el" class="lang-chart"></div>
                <div class="lang-table">
                  <div class="lang-header">
                    <span>{{ t('repo.lang') }}</span>
                    <span>{{ t('repo.files') }}</span>
                    <span>{{ t('repo.lines') }}</span>
                    <span>{{ t('repo.percent') }}</span>
                  </div>
                  <div v-for="lang in detail.analysis.languages" :key="lang.name" class="lang-row">
                    <span class="lang-name">{{ lang.name }}</span>
                    <span>{{ lang.fileCount }}</span>
                    <span>{{ formatNumber(lang.lines) }}</span>
                    <span class="lang-pct">
                      <span class="pct-track"><span class="pct-bar" :style="{ width: lang.percentage + '%' }"></span></span>
                      <span class="pct-text">{{ lang.percentage.toFixed(1) }}%</span>
                    </span>
                  </div>
                  <div class="lang-total">
                    <span>{{ t('repo.total') }}</span>
                    <span>{{ detail.analysis.fileCount }}</span>
                    <span>{{ formatNumber(detail.analysis.totalLines) }}</span>
                    <span>100%</span>
                  </div>
                </div>
              </div>
              <div v-else class="analyze-cta">
                <span class="analyze-icon">🔍</span>
                <h4>{{ t('repo.analyzeTitle') }}</h4>
                <p class="analyze-desc">{{ t('repo.analyzeDesc') }}</p>
                <button class="btn" :disabled="state.analyzing" @click="doAnalyze">
                  <span v-if="state.analyzing" class="spinner"></span>
                  {{ state.analyzing ? t('repo.analyzing') : t('repo.analyzeBtn') }}
                </button>
              </div>
            </div>

            <div v-if="!statsLoaded" class="section card">
              <h3 class="section-title">{{ t('repo.recentCommits') }}</h3>
              <div class="commit-list">
                <div v-for="i in 5" :key="i" class="commit-main">
                  <span class="skeleton-line w15" style="height:12px"></span>
                  <span class="skeleton-line w45" style="height:12px"></span>
                  <span class="skeleton-line w12" style="height:12px"></span>
                  <span class="skeleton-line w12" style="height:12px"></span>
                  <span class="skeleton-line w18" style="height:12px"></span>
                </div>
              </div>
            </div>

            <template v-if="statsLoaded">
              <div v-if="detail.recentCommits && detail.recentCommits.length > 0" class="section card">
                <h3 class="section-title">{{ t('repo.recentCommits') }}</h3>
                <div class="commit-list">
                  <div v-for="c in detail.recentCommits" :key="c.hash" class="commit-item">
                    <div class="commit-main" @click="toggleCommit(c.hash)">
                      <span class="commit-hash">{{ c.hash.slice(0, 7) }}</span>
                      <span class="commit-msg">{{ c.message.split('\n')[0] }}</span>
                      <span class="commit-author">{{ c.author }}</span>
                      <span class="commit-time">{{ formatTimeAgo(c.date) }}</span>
                      <span class="commit-changes">
                        <span v-if="c.additions > 0" class="add">+{{ c.additions }}</span>
                        <span v-if="c.deletions > 0" class="del">-{{ c.deletions }}</span>
                      </span>
                    </div>
                    <div v-if="expandedCommit === c.hash" class="commit-body">
                      <pre>{{ c.message }}</pre>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </div>

        </template>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useI18n } from '../i18n'
import { state, fetchReposInfo, fetchRepoInfo, fetchRepoStats, fetchRepoChart, triggerAnalyze } from '../stores/data'
import echarts from '../utils/echarts'
import RepoCharts from '../components/RepoCharts.vue'

const { t } = useI18n()

const activePath = ref(localStorage.getItem('activeRepoPath') || '')
const loading = ref(false)
const detail = ref(null)
const statsLoaded = ref(false)
const loadingStats = ref(false)
const expandedBranch = ref(false)
const expandedTags = ref(false)
const expandedContributor = ref(false)
const expandedCommit = ref(null)
const analysisLoading = ref(false)
const chartData = ref(null)
const chartLoading = ref(false)
const langChartRef = ref(null)
let langChart = null

const localBranchCount = computed(() => {
  const local = detail.value?.branchCount ?? 0
  const remote = Array.isArray(detail.value?.remoteBranches) ? detail.value.remoteBranches.length : 0
  if (local === 0 && remote === 0) return '--'
  return remote > 0 ? `${local}+${remote}` : String(local)
})

const remoteBranchCount = computed(() => {
  const v = detail.value?.remoteBranches
  return Array.isArray(v) ? v.length : '--'
})

const tagCount = computed(() => {
  const v = detail.value?.tags
  return Array.isArray(v) ? v.length : '--'
})

const localBranches = computed(() => {
  if (detail.value?.branches) return detail.value.branches
  return []
})

const totalCommitCount = computed(() =>
  detail.value?.contributors?.reduce((s, c) => s + c.commitCount, 0) || 0
)

const hasCommits = computed(() =>
  detail.value?.contributors && detail.value.contributors.length > 0
)

const mainLangName = computed(() => {
  if (!detail.value?.analysis?.languages?.length) return '--'
  const top = detail.value.analysis.languages[0]
  return `${top.name} ${top.percentage.toFixed(0)}%`
})

const timeSpan = computed(() => {
  if (!detail.value?.earliestDate || !detail.value?.lastCommitTime) return ''
  const d1 = new Date(detail.value.earliestDate)
  const d2 = new Date(detail.value.lastCommitTime)
  const ms = d2 - d1
  if (ms <= 0) return ''
  const days = Math.round(ms / 86400000)
  if (days < 30) return days + t('repo.dayUnit')
  if (days < 365) return Math.round(days / 30) + t('repo.monthUnit')
  const y = Math.floor(days / 365)
  const m = Math.round((days % 365) / 30)
  return y + t('repo.yearUnit') + (m ? m + t('repo.monthUnit') : '')
})

function formatDate(s) {
  if (!s) return '-'
  return s.slice(0, 10)
}

function formatTimeAgo(s) {
  if (!s) return '-'
  const d = new Date(s)
  const now = Date.now()
  const diff = (now - d.getTime()) / 1000
  if (diff < 60) return Math.round(diff) + 's'
  if (diff < 3600) return Math.round(diff / 60) + 'm'
  if (diff < 86400) return Math.round(diff / 3600) + 'h'
  if (diff < 2592000) return Math.round(diff / 86400) + 'd'
  return s.slice(0, 10)
}

function formatBytes(n) {
  if (!n) return '0B'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  let v = n
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++ }
  return v.toFixed(i > 0 ? 1 : 0) + units[i]
}

function formatNumber(n) {
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k'
  return String(n)
}

function toggleBranch() {
  expandedBranch.value = !expandedBranch.value
}

function toggleTags() {
  expandedTags.value = !expandedTags.value
}

function toggleContributor() {
  expandedContributor.value = !expandedContributor.value
}

function toggleCommit(hash) {
  expandedCommit.value = expandedCommit.value === hash ? null : hash
}

async function loadDetail(path) {
  loading.value = true
  statsLoaded.value = false
  expandedBranch.value = false
  expandedTags.value = false
  expandedContributor.value = false
  expandedCommit.value = null
  langChart = null
  try {
    const info = await fetchRepoInfo(path)
    detail.value = {
      contributors: [], recentCommits: [], repoSize: 0,
      earliestDate: '', earliestCommitAuthor: '',
      tags: [], remoteBranches: [],
      ...info,
      tags: info.tags || [],
      remoteBranches: info.remoteBranches || [],
      branches: info.branches || []
    }
    if (state.repoStatsCache[path]) {
      detail.value = { ...detail.value, ...state.repoStatsCache[path] }
      statsLoaded.value = true
      await nextTick()
      renderChart()
      fetchChartData()
    }
    if (state.analyzeCache[path]) {
      detail.value.analysis = state.analyzeCache[path]
      await nextTick()
      renderChart()
    }
    if (state.repoChartCache[path]) {
      chartData.value = state.repoChartCache[path]
    }
  } catch (err) {
    console.error('Failed to load info:', err)
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  if (!activePath.value || loadingStats.value) return
  loadingStats.value = true
  try {
    const stats = await fetchRepoStats(activePath.value)
    detail.value = { ...detail.value, ...stats }
    statsLoaded.value = true
    await nextTick()
    renderChart()
    fetchChartData()
  } catch (err) {
    console.error('Failed to load stats:', err)
  } finally {
    loadingStats.value = false
  }
}

async function fetchChartData() {
  chartLoading.value = true
  try {
    const data = await fetchRepoChart(activePath.value)
    chartData.value = data
  } catch (err) {
    console.error('Failed to load chart data:', err)
  } finally {
    chartLoading.value = false
  }
}

function switchRepo(path) {
  activePath.value = path
  localStorage.setItem('activeRepoPath', path)
  loadDetail(path)
}

async function doAnalyze() {
  analysisLoading.value = true
  try {
    const result = await triggerAnalyze(activePath.value)
    detail.value.analysis = result
  } catch (err) {
    console.error('Analysis failed:', err)
  } finally {
    analysisLoading.value = false
    await nextTick()
    renderChart()
  }
}

const LANG_COLORS = {
  Rust: '#dea584', Go: '#00add8', Python: '#3572a5',
  JavaScript: '#f1e05a', TypeScript: '#3178c6',
  'React (JSX)': '#61dafb', 'React (TSX)': '#61dafb',
  Vue: '#41b883', Java: '#b07219', Kotlin: '#a97bff',
  C: '#555555', 'C++': '#f34b7d', 'C#': '#178600',
  Ruby: '#701516', PHP: '#4f5d95', Swift: '#ffac45',
  Dart: '#00b4ab', CSS: '#563d7c', SCSS: '#c6538c',
  HTML: '#e34c26', Shell: '#89e051', Dockerfile: '#384d54',
  SQL: '#e38c00', Markdown: '#083fa1', YAML: '#cb171e',
  TOML: '#9c4221', JSON: '#292929', Lua: '#000080',
  Scala: '#c22d40', Zig: '#ec915c', Default: '#a0aec0'
}

function getLangColor(name) {
  return LANG_COLORS[name] || LANG_COLORS.Default
}

function renderChart() {
  const el = document.getElementById('lang-chart-el')
  if (!el || !detail.value?.analysis) return
  try {
    if (!langChart) langChart = echarts.init(el)
    const data = detail.value.analysis.languages.map(l => ({
      name: l.name,
      value: l.lines || 1
    }))
    langChart.setOption({
      tooltip: {
        trigger: 'item',
        backgroundColor: 'rgba(10, 14, 39, 0.9)',
        borderColor: 'rgba(0, 212, 255, 0.3)',
        textStyle: { color: '#e0e6ff', fontSize: 12 },
        formatter: p => `${p.name}: ${p.percent}% (${formatNumber(p.value)} lines)`
      },
      series: [{
        type: 'pie',
        radius: ['30%', '65%'],
        center: ['50%', '50%'],
        avoidLabelOverlap: true,
        itemStyle: { borderRadius: 4, borderColor: 'rgba(10,14,39,0.8)', borderWidth: 2 },
        label: {
          show: true, color: '#e0e6ff', fontFamily: 'Rajdhani', fontSize: 12,
          formatter: p => `${p.name}\n${p.percent}%`
        },
        labelLine: { lineStyle: { color: 'rgba(0,212,255,0.3)' } },
        data: data.map(d => ({ ...d, itemStyle: { color: getLangColor(d.name) } }))
      }]
    })
    langChart.resize()
  } catch (e) {
    console.error('[Chart] render error:', e)
  }
}

function handleResize() {
  if (langChart) langChart.resize()
}

function init() {
  fetchReposInfo().then(() => {
    if (state.reposInfo.length > 0) {
      const saved = activePath.value
      const exists = state.reposInfo.some(r => r.path === saved)
      activePath.value = exists ? saved : state.reposInfo[0].path
      localStorage.setItem('activeRepoPath', activePath.value)
      loadDetail(activePath.value)
    }
  })
  window.addEventListener('resize', handleResize)
}

watch(() => detail.value?.analysis, val => {
  if (val) nextTick(renderChart)
})

onMounted(init)

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  langChart?.dispose()
  langChart = null
})
</script>

<style scoped>
.repo-page {
  max-width: 1400px;
  margin: 0 auto;
  animation: fadeIn 0.5s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Empty */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  text-align: center;
}
.empty-icon { font-size: 4rem; color: #00d4ff; margin-bottom: 1rem; opacity: 0.5; }
.empty-state p { color: #64748b; font-size: 1.1rem; font-family: 'Rajdhani', sans-serif; letter-spacing: 1px; }

/* Tabs */
.tabs {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 2rem;
  overflow-x: auto;
  padding-bottom: 0.5rem;
}
.tab {
  flex-shrink: 0;
  padding: 0.6rem 1.4rem;
  border: 1px solid rgba(0, 212, 255, 0.2);
  border-radius: 8px;
  background: rgba(20, 25, 50, 0.4);
  color: #a0aec0;
  font-family: 'Orbitron', sans-serif;
  font-size: 0.85rem;
  letter-spacing: 1px;
  cursor: pointer;
  transition: all 0.3s;
}
.tab:hover { border-color: rgba(0, 212, 255, 0.5); color: #00d4ff; }
.tab.active {
  background: rgba(0, 212, 255, 0.1);
  border-color: #00d4ff;
  color: #00d4ff;
  box-shadow: 0 0 15px rgba(0, 212, 255, 0.2);
}

/* Loading */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}
.loading-state p { font-family: 'Orbitron', sans-serif; color: #00d4ff; letter-spacing: 1px; margin-top: 1rem; }

/* Title */
.title-row {
  display: flex;
  align-items: baseline;
  gap: 1rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
}
.repo-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 1.8rem;
  margin: 0;
  background: linear-gradient(135deg, #00d4ff, #00ff88);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.repo-path { font-family: 'Rajdhani', monospace; color: #64748b; font-size: 0.9rem; }

/* Spinner */
.spinner {
  display: inline-block;
  width: 16px; height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  vertical-align: middle;
}
.spinner.large { width: 32px; height: 32px; border-width: 3px; border-color: rgba(0,212,255,0.2); border-top-color: #00d4ff; }

/* Card rows */
.card-row {
  display: grid;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}
.l1-row { grid-template-columns: repeat(4, 1fr); }
.l2-row { grid-template-columns: repeat(4, 1fr); }

.info-card {
  background: rgba(20, 25, 50, 0.6);
  backdrop-filter: blur(20px);
  padding: 1.5rem 1rem;
  border-radius: 12px;
  border: 1px solid rgba(0, 212, 255, 0.15);
  text-align: center;
  transition: all 0.3s;
}
.info-card.dim { border-color: rgba(0, 212, 255, 0.08); }
.info-card.dim.has-data { border-color: rgba(0, 212, 255, 0.15); }
.info-card.clickable { cursor: pointer; }
.info-card.clickable:hover, .info-card.expanded {
  border-color: rgba(0, 212, 255, 0.4);
  box-shadow: 0 0 20px rgba(0, 212, 255, 0.1);
  transform: translateY(-2px);
}

.info-icon { display: block; font-size: 1.8rem; margin-bottom: 0.5rem; }
.info-value {
  display: block;
  font-family: 'Orbitron', sans-serif;
  font-size: 1.3rem;
  font-weight: 700;
  color: #e0e6ff;
  margin-bottom: 0.3rem;
}
.dim .info-value { color: #4a5568; }
.dim.has-data .info-value { color: #e0e6ff; }
.info-label {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.65rem;
  color: #a0aec0;
  letter-spacing: 1px;
  text-transform: uppercase;
}
.expand-icon { font-size: 0.5rem; margin-left: 0.2rem; vertical-align: middle; }

/* Expand panel */
.expand-panel {
  margin-bottom: 1.5rem;
  padding: 1.5rem;
}
.expand-panel h4 {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.9rem;
  color: #00d4ff;
  letter-spacing: 1px;
  margin: 0 0 1rem 0;
}
.expand-hint { color: #64748b; font-family: 'Rajdhani', sans-serif; font-size: 0.9rem; }

/* Branch groups */
.branch-groups {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.branch-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.branch-group-label {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.75rem;
  color: #a0aec0;
  letter-spacing: 1px;
  text-transform: uppercase;
}
.branch-tags { display: flex; flex-wrap: wrap; gap: 0.6rem; }
.branch-tag {
  padding: 0.3rem 0.8rem;
  border-radius: 12px;
  font-size: 0.85rem;
  font-family: 'Rajdhani', monospace;
  background: rgba(0, 212, 255, 0.08);
  border: 1px solid rgba(0, 212, 255, 0.2);
  color: #a0aec0;
}
.branch-tag.current {
  background: rgba(0, 255, 136, 0.12);
  border-color: rgba(0, 255, 136, 0.4);
  color: #00ff88;
  font-weight: 600;
}
.current-badge {
  font-size: 0.6rem;
  margin-left: 0.3rem;
  padding: 0.05rem 0.35rem;
  border-radius: 4px;
  background: rgba(0, 255, 136, 0.2);
  color: #00ff88;
  font-family: 'Rajdhani', sans-serif;
  vertical-align: middle;
  letter-spacing: 0.5px;
}
.branch-tag.remote {
  background: rgba(167, 139, 250, 0.08);
  border-color: rgba(167, 139, 250, 0.2);
  color: #a78bfa;
}

/* Tag cloud */
.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}
.tag-item {
  padding: 0.25rem 0.7rem;
  border-radius: 6px;
  font-size: 0.8rem;
  font-family: 'Rajdhani', monospace;
  background: rgba(52, 211, 153, 0.1);
  border: 1px solid rgba(52, 211, 153, 0.25);
  color: #34d399;
}

/* Contributor table */
.contrib-table { width: 100%; }
.contrib-header, .contrib-row {
  display: grid;
  grid-template-columns: 1.5fr 1fr 1fr 1fr 1.5fr;
  padding: 0.6rem 0.5rem;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.85rem;
}
.contrib-header {
  background: rgba(0, 212, 255, 0.05);
  border-bottom: 1px solid rgba(0, 212, 255, 0.15);
  font-family: 'Orbitron', sans-serif;
  font-size: 0.7rem;
  color: #a0aec0;
  letter-spacing: 1px;
  text-transform: uppercase;
}
.contrib-row { border-bottom: 1px solid rgba(0, 212, 255, 0.06); color: #e0e6ff; }
.contrib-name { font-weight: 600; }
.contrib-time { font-family: 'Rajdhani', monospace; color: #64748b; font-size: 0.8rem; }
.add { color: #00ff88; }
.del { color: #ff6b6b; }

/* Section group */
.section-group {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  margin-bottom: 2rem;
}
.section { padding: 1.5rem; }
.section-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.95rem;
  color: #00d4ff;
  letter-spacing: 1px;
  margin: 0 0 1rem 0;
}
.section-loading { text-align: center; padding: 2rem; color: #00d4ff; font-family: 'Rajdhani', sans-serif; }

/* Language */
.lang-body { display: grid; grid-template-columns: 1fr 1.2fr; gap: 1.5rem; align-items: start; }
.lang-chart { height: 250px; min-height: 250px; width: 100%; }
.lang-header, .lang-row, .lang-total {
  display: grid;
  grid-template-columns: 1.5fr 0.8fr 1fr 1.5fr;
  padding: 0.5rem 0.5rem;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.85rem;
}
.lang-header {
  background: rgba(0, 212, 255, 0.05);
  border-bottom: 1px solid rgba(0, 212, 255, 0.15);
  font-family: 'Orbitron', sans-serif;
  font-size: 0.7rem;
  color: #a0aec0;
  letter-spacing: 1px;
  text-transform: uppercase;
}
.lang-row { border-bottom: 1px solid rgba(0, 212, 255, 0.06); color: #e0e6ff; }
.lang-row:hover { background: rgba(0, 212, 255, 0.03); }
.lang-name { font-weight: 600; color: #e0e6ff; }
.lang-pct { display: flex; align-items: center; gap: 0.3rem; font-family: 'Orbitron', sans-serif; font-size: 0.75rem; color: #00d4ff; }
.pct-track { flex: 1; height: 8px; background: rgba(0,0,0,0.15); border-radius: 4px; overflow: hidden; min-width: 4px; }
.pct-bar { display: block; height: 8px; border-radius: 4px; background: linear-gradient(90deg, #00d4ff, #7800ff); }
.pct-text { flex-shrink: 0; font-family: 'Rajdhani', 'Orbitron', monospace; }
.lang-total {
  border-top: 1px solid rgba(0, 212, 255, 0.3);
  font-family: 'Orbitron', sans-serif;
  font-size: 0.8rem;
  color: #00d4ff;
  font-weight: 700;
  padding-top: 0.6rem;
}
.lang-total span { text-align: center; }
.lang-total span:first-child { text-align: left; }

/* Stats CTA (blank state skeleton) */
.stats-cta {
  cursor: pointer;
  margin-bottom: 1.5rem;
  padding: 1.5rem;
  transition: all 0.3s;
}
.stats-cta:hover {
  border-color: rgba(0, 212, 255, 0.4);
  box-shadow: 0 0 20px rgba(0, 212, 255, 0.1);
}
.stats-cta-content {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}
.stats-cta-content div { flex: 1; }
.stats-cta-content h4 {
  font-family: 'Orbitron', sans-serif;
  font-size: 1rem;
  color: #00d4ff;
  letter-spacing: 1px;
  margin: 0 0 0.3rem 0;
}
.stats-cta-content p {
  color: #64748b;
  font-size: 0.85rem;
  margin: 0;
  font-family: 'Rajdhani', sans-serif;
}

/* Stats group */
.stats-group {
  position: relative;
  margin-bottom: 1.5rem;
}
.stats-group-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(10, 14, 39, 0.55);
  backdrop-filter: blur(4px);
  border-radius: 12px;
  cursor: pointer;
  z-index: 2;
  transition: all 0.3s;
}
.stats-group-overlay:hover {
  background: rgba(10, 14, 39, 0.65);
  backdrop-filter: blur(6px);
}
.stats-group-overlay-content {
  text-align: center;
  padding: 1.5rem;
}
.stats-group-overlay-content h4 {
  font-family: 'Orbitron', sans-serif;
  font-size: 1rem;
  color: #00d4ff;
  letter-spacing: 1px;
  margin: 0.5rem 0 0.3rem 0;
}
.stats-group-overlay-content p {
  color: #a0aec0;
  font-size: 0.85rem;
  margin: 0 0 1rem 0;
  font-family: 'Rajdhani', sans-serif;
}

.stats-cta-icon {
  display: block;
  font-size: 2.5rem;
  color: #00d4ff;
  opacity: 0.6;
}

/* Analyze CTA */
.analyze-cta { text-align: center; padding: 2rem 1rem; }
.analyze-icon { font-size: 2.5rem; display: block; margin-bottom: 0.8rem; }
.analyze-cta h4 { font-family: 'Orbitron', sans-serif; font-size: 1.1rem; color: #00d4ff; letter-spacing: 1px; margin-bottom: 0.4rem; }
.analyze-desc { color: #64748b; font-size: 0.9rem; margin-bottom: 1rem; font-family: 'Rajdhani', sans-serif; }

/* Commits */
.commit-list { width: 100%; }
.commit-main {
  display: grid;
  grid-template-columns: 1fr 2.5fr 1fr 1fr 1.2fr;
  padding: 0.55rem 0.5rem;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.85rem;
  color: #e0e6ff;
  border-bottom: 1px solid rgba(0, 212, 255, 0.06);
  cursor: pointer;
  transition: background 0.2s;
}
.commit-main:hover { background: rgba(0, 212, 255, 0.03); }
.commit-hash { font-family: 'Rajdhani', monospace; color: #a0aec0; }
.commit-msg { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.commit-author { font-size: 0.8rem; }
.commit-time { font-family: 'Rajdhani', monospace; color: #64748b; font-size: 0.8rem; }
.commit-changes { text-align: right; }
.commit-changes .add { margin-right: 0.3rem; }
.commit-body pre {
  margin: 0; padding: 0.6rem 0.5rem;
  background: rgba(0, 0, 0, 0.2);
  font-family: 'Rajdhani', monospace;
  font-size: 0.8rem;
  color: #a0aec0;
  border-bottom: 1px solid rgba(0, 212, 255, 0.06);
  white-space: pre-wrap;
}

/* Skeleton page */
.skeleton-page {
  animation: fadeIn 0.5s ease-out;
}

/* Charts skeleton */
.charts-section-skel {
  margin-top: 2rem;
}
.charts-section-skel .section-header {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}
.chart-row-skel {
  display: flex;
  gap: 1.5rem;
}
.chart-row-skel .card {
  border-radius: 16px;
}
</style>
