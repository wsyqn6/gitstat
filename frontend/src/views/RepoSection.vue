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
          :class="['tab btn-ghost', { active: activePath === repo.path }]"
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
              <Skeleton w="40" h="32" radius="6" />
            </div>
            <div class="card-row l1-row">
              <div v-for="i in 4" :key="i" class="info-card">
                <Skeleton circle h="36" center mb="0.5rem" />
                <Skeleton w="35" h="22" radius="6" center mb="0.4rem" />
                <Skeleton w="55" center />
              </div>
            </div>
            <div class="card stats-cta">
              <div class="stats-cta-content">
                <Skeleton circle h="40" />
                <div style="flex:1;display:flex;flex-direction:column;gap:0.3rem">
                  <Skeleton w="45" h="18" radius="6" />
                  <Skeleton w="70" />
                </div>
                <Skeleton w="20" h="38" radius="8" />
              </div>
            </div>
            <div class="card-row l2-row">
              <div v-for="i in 4" :key="i" class="info-card dim">
                <Skeleton w="40" h="22" radius="6" center mb="0.3rem" />
                <Skeleton w="65" center />
              </div>
            </div>
            <div class="section-group">
              <div class="section card">
                <Skeleton w="50" h="20" radius="6" mb="1rem" />
                <Skeleton w="100" h="10" radius="5" mb="0.75rem" />
                <div style="display:flex;flex-direction:column;gap:0.5rem">
                  <Skeleton v-for="i in 5" :key="i" w="70" />
                </div>
              </div>
              <div class="section card">
                <Skeleton w="50" h="20" radius="6" mb="1rem" />
                <Skeleton v-for="i in 5" :key="i" w="85" mb="0.5rem" />
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
            :chart-data="chartData"
            :chart-loading="chartLoading"
            :chart-loaded="chartLoaded"
            :tags="tags"
            :tags-loading="tagsLoading"
            :tags-has-more="tagsHasMore"
            :tags-loaded="tagsLoaded"
            @load-chart="loadChart"
            @load-tags="loadTags"
            @load-more-tags="loadMoreTags"
          />

          <div class="section-group">
            <RepoLangChart
              :analysis="detail.analysis"
              :analysis-loading="analysisLoading"
              @analyze="doAnalyze"
            />
            <RepoCommits
              :commits="commits"
              :commits-loading="commitsLoading"
              :commits-loaded="commitsLoaded"
              :has-more="commitsHasMore"
              @load-commits="loadCommits"
              @load-more="loadMoreCommits"
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
import { state, fetchReposInfo, fetchRepoInfo, fetchRepoStats, fetchRepoChart, fetchRepoCommits, fetchRepoTagsCount, fetchRepoTagsPage, triggerAnalyze } from '../stores/data'
import Skeleton from '../components/Skeleton.vue'
import RepoInfoCards from '../components/RepoInfoCards.vue'
import RepoLangChart from '../components/RepoLangChart.vue'
import RepoCommits from '../components/RepoCommits.vue'

const { t } = useI18n()

const activePath = ref(localStorage.getItem('activeRepoPath') || '')
const loading = ref(false)
const detail = ref(null)
const statsLoaded = ref(false)
const analysisLoading = ref(false)
const chartData = ref(null)
const chartLoading = ref(false)
const chartLoaded = ref(false)
const commits = ref([])
const commitsLoading = ref(false)
const commitsLoaded = ref(false)
const commitsHasMore = ref(false)
const commitsOffset = ref(0)
const tags = ref([])
const tagsLoading = ref(false)
const tagsTotal = ref(0)
const tagsOffset = ref(0)
const tagsHasMore = ref(false)
const tagsLoaded = ref(false)

async function loadDetail(path) {
  loading.value = true
  statsLoaded.value = false
  chartLoaded.value = false
  chartData.value = null
  tags.value = []
  tagsLoaded.value = false
  tagsOffset.value = 0
  try {
    const info = await fetchRepoInfo(path)
    detail.value = {
      contributors: [], repoSize: 0,
      earliestDate: '', earliestCommitAuthor: '',
      tagCount: 0, remoteBranches: [],
      ...info,
      remoteBranches: info.remoteBranches || [],
      branches: info.branches || []
    }
    fetchRepoTagsCount(path).then(res => {
      if (detail.value) detail.value.tagCount = res.tagCount ?? 0
    })
    fetchRepoStats(path).then(stats => {
      if (detail.value) {
        detail.value = { ...detail.value, ...stats }
        statsLoaded.value = true
      }
    }).catch(err => console.error('Failed to load stats:', err))
    if (state.analyzeCache.has(path)) {
      detail.value.analysis = state.analyzeCache.get(path)
    }
    if (state.repoChartCache.has(path)) {
      chartData.value = state.repoChartCache.get(path)
      chartLoaded.value = true
    }
  } catch (err) {
    console.error('Failed to load info:', err)
  } finally {
    loading.value = false
  }
}

async function loadChart() {
  if (!activePath.value || chartLoading.value) return
  chartLoading.value = true
  try {
    const data = await fetchRepoChart(activePath.value)
    chartData.value = data
    chartLoaded.value = true
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

async function loadCommits() {
  if (!activePath.value || commitsLoading.value) return
  commitsLoading.value = true
  commitsOffset.value = 0
  try {
    const result = await fetchRepoCommits(activePath.value, 0, 30)
    commits.value = result.commits
    commitsHasMore.value = result.hasMore
    commitsLoaded.value = true
  } catch (err) {
    console.error('Failed to load commits:', err)
  } finally {
    commitsLoading.value = false
  }
}

async function loadMoreCommits() {
  if (!activePath.value || commitsLoading.value || !commitsHasMore.value) return
  commitsLoading.value = true
  const nextOffset = commitsOffset.value + 30
  try {
    const result = await fetchRepoCommits(activePath.value, nextOffset, 30)
    commits.value = commits.value.concat(result.commits)
    commitsHasMore.value = result.hasMore
    commitsOffset.value = nextOffset
  } catch (err) {
    console.error('Failed to load more commits:', err)
  } finally {
    commitsLoading.value = false
  }
}

async function loadTags() {
  if (!activePath.value || tagsLoading.value) return
  tagsLoading.value = true
  tagsOffset.value = 0
  try {
    const result = await fetchRepoTagsPage(activePath.value, 0, 30)
    tags.value = result.tags
    tagsTotal.value = result.total
    tagsHasMore.value = result.hasMore
    tagsLoaded.value = true
  } catch (err) {
    console.error('Failed to load tags:', err)
  } finally {
    tagsLoading.value = false
  }
}

async function loadMoreTags() {
  if (!activePath.value || tagsLoading.value || !tagsHasMore.value) return
  tagsLoading.value = true
  const nextOffset = tagsOffset.value + 30
  try {
    const result = await fetchRepoTagsPage(activePath.value, nextOffset, 30)
    tags.value = tags.value.concat(result.tags)
    tagsHasMore.value = result.hasMore
    tagsOffset.value = nextOffset
  } catch (err) {
    console.error('Failed to load more tags:', err)
  } finally {
    tagsLoading.value = false
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
.empty-icon { font-size: 4rem; color: var(--color-primary); margin-bottom: 1rem; opacity: 0.5; }
.empty-state p { color: var(--color-text-muted); font-size: 1.1rem; font-family: var(--font-body); letter-spacing: 1px; }

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
  border: 1px solid var(--glass-outline-border);
  border-radius: var(--radius-btn);
  font-size: 0.85rem;
}
.tab:hover { border-color: var(--color-primary); }
.tab.active {
  background: var(--glass-btn-bg);
  backdrop-filter: blur(var(--glass-blur));
  border-color: var(--color-accent);
  color: var(--glass-btn-color);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}
.loading-state p { font-family: var(--font-display); color: var(--color-primary); letter-spacing: 1px; margin-top: 1rem; }

.spinner {
  display: inline-block;
  width: 16px; height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  vertical-align: middle;
}
.spinner.large { width: 32px; height: 32px; border-width: 3px; border-color: var(--border-insight-card); border-top-color: var(--color-primary); }

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
  background: var(--gradient-repo-title);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.repo-path { font-family: var(--font-mono); color: var(--color-text-muted); font-size: 0.9rem; }

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
  background: var(--bg-card);
  backdrop-filter: blur(var(--blur-card));
  padding: 1.5rem;
  border-radius: var(--radius-card);
  border: 1px solid var(--border-card);
  text-align: center;
}
.info-card.dim { opacity: 0.5; }

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
</style>
