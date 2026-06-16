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

          <RepoInfoCards
            :detail="detail"
            :stats-loaded="statsLoaded"
            :loading-stats="loadingStats"
            :chart-data="chartData"
            :chart-loading="chartLoading"
            @load-stats="loadStats"
          />

          <div class="section-group">
            <RepoLangChart
              :analysis="detail.analysis"
              :analysis-loading="analysisLoading"
              @analyze="doAnalyze"
            />
            <RepoCommits
              :commits="detail.recentCommits || []"
              :loading="!statsLoaded"
            />
          </div>

        </template>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from '../i18n'
import { state, fetchReposInfo, fetchRepoInfo, fetchRepoStats, fetchRepoChart, triggerAnalyze } from '../stores/data'
import RepoInfoCards from '../components/RepoInfoCards.vue'
import RepoLangChart from '../components/RepoLangChart.vue'
import RepoCommits from '../components/RepoCommits.vue'

const { t } = useI18n()

const activePath = ref(localStorage.getItem('activeRepoPath') || '')
const loading = ref(false)
const detail = ref(null)
const statsLoaded = ref(false)
const loadingStats = ref(false)
const analysisLoading = ref(false)
const chartData = ref(null)
const chartLoading = ref(false)

async function loadDetail(path) {
  loading.value = true
  statsLoaded.value = false
  chartData.value = null
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
    if (state.repoStatsCache.has(path)) {
      detail.value = { ...detail.value, ...state.repoStatsCache.get(path) }
      statsLoaded.value = true
      fetchChartData()
    }
    if (state.analyzeCache.has(path)) {
      detail.value.analysis = state.analyzeCache.get(path)
    }
    if (state.repoChartCache.has(path)) {
      chartData.value = state.repoChartCache.get(path)
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
  }
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
}

onMounted(init)
</script>

<style scoped>
.repo-page {
  max-width: 1400px;
  margin: 0 auto;
  animation: fadeIn 0.5s ease-out;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  text-align: center;
}
.empty-icon { font-size: 4rem; color: #00d4ff; margin-bottom: 1rem; opacity: 0.5; }
.empty-state p { color: #64748b; font-size: 1.1rem; font-family: var(--font-body); letter-spacing: 1px; }

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
  font-family: var(--font-display);
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

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}
.loading-state p { font-family: var(--font-display); color: #00d4ff; letter-spacing: 1px; margin-top: 1rem; }

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

.title-row {
  display: flex;
  align-items: baseline;
  gap: 1rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
}
.repo-title {
  font-family: var(--font-display);
  font-size: 1.8rem;
  margin: 0;
  background: linear-gradient(135deg, #00d4ff, #00ff88);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.repo-path { font-family: 'Rajdhani', monospace; color: #64748b; font-size: 0.9rem; }

.skeleton-page {
  animation: fadeIn 0.5s ease-out;
}

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
}
.info-card.dim { border-color: rgba(0, 212, 255, 0.08); }

.stats-cta {
  margin-bottom: 1.5rem;
  padding: 1.5rem;
}
.stats-cta-content {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}
.stats-cta-content div { flex: 1; }

.section-group {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  margin-bottom: 2rem;
}
.section { padding: 1.5rem; }

.lang-body { display: grid; grid-template-columns: 1fr 1.2fr; gap: 1.5rem; align-items: start; }
</style>
