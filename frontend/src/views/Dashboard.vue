<template>
  <div class="dashboard">
    <DashStatsGrid
      :loading="sectionLoading.stats"
      :summary="state.dashboardData?.summary"
      :weekly-total="weeklyTotal"
      :comparison="state.dashboardData?.comparison"
    />

    <DailyRepoStats
      :daily-stats="state.dashboardData?.dailyRepos || []"
      :loading="sectionLoading.stats"
      :comparison-authors="state.dashboardData?.comparison?.authors || []"
    />

    <div class="insight-grid">
      <WeeklyTrend
        :repo-daily-trend="state.repoDailyTrend"
        :repo-colors="repoColors"
        :week-range="weekRange"
        :loading="sectionLoading.trend"
      />
      <AuthorRanking
        :author-rank-with-repos="authorRankWithRepos"
        :week-range="weekRange"
        :loading="sectionLoading.rank"
      />
    </div>

    <div ref="section2Ref">
      <RepoComparison
        :repo-comparison="state.repoComparison"
        :week-range="weekRange"
        :loading="sectionLoading.below"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { state, fetchDashboardData, fetchRepoDailyTrend, fetchAuthorRank, loadDashboardS2 } from '../stores/data'
import WeeklyTrend from '../components/WeeklyTrend.vue'
import DailyRepoStats from '../components/DailyRepoStats.vue'
import DashStatsGrid from '../components/DashStatsGrid.vue'
import AuthorRanking from '../components/AuthorRanking.vue'
import RepoComparison from '../components/RepoComparison.vue'
import { CHART_COLORS } from '../utils/constants'

const sectionLoading = reactive({
  stats: true,
  trend: true,
  rank: true,
  below: true
})

function getWeekRange() {
  const now = new Date()
  const dow = now.getDay()
  const diff = dow === 0 ? 6 : dow - 1
  const monday = new Date(now)
  monday.setDate(now.getDate() - diff)
  const fmt = (d) => `${d.getMonth() + 1}/${d.getDate()}`
  return `${fmt(monday)} - ${fmt(now)}`
}
const weekRange = getWeekRange()

const section2Ref = ref(null)
let observer = null

const weeklyTotal = computed(() =>
  state.repoDailyTrend.reduce((sum, repo) =>
    sum + repo.data.reduce((s, d) => s + d.commits, 0), 0)
)

const repoColors = computed(() =>
  state.repoDailyTrend.map((_, i) => CHART_COLORS[i % CHART_COLORS.length])
)

const authorRankWithRepos = computed(() => {
  if (!state.authorRank || state.authorRank.length === 0) return []
  const dailyRepos = state.dashboardData?.dailyRepos || []
  const repoOf = {}
  for (const repo of dailyRepos) {
    for (const author of repo.authors || []) {
      if (!repoOf[author.email]) repoOf[author.email] = []
      const total = author.dailyData ? author.dailyData.reduce((s, d) => s + d.commits, 0) : 0
      repoOf[author.email].push({ name: repo.repoName, commits: total, color: CHART_COLORS[dailyRepos.indexOf(repo) % CHART_COLORS.length] })
    }
  }
  return state.authorRank.map(a => ({
    ...a,
    repos: (repoOf[a.email] || []).filter(r => r.commits > 0)
  }))
})

onMounted(async () => {
  await Promise.all([
    fetchDashboardData().then(() => { sectionLoading.stats = false }),
    fetchRepoDailyTrend().then(() => { sectionLoading.trend = false }),
    fetchAuthorRank().then(() => { sectionLoading.rank = false })
  ])

  observer = new IntersectionObserver(
    ([entry]) => {
      if (entry.isIntersecting && sectionLoading.below) {
        loadDashboardS2().then(() => { sectionLoading.below = false })
        observer.disconnect()
      }
    },
    { rootMargin: '200px' }
  )
  if (section2Ref.value) {
    observer.observe(section2Ref.value)
  }
})

onUnmounted(() => {
  if (observer) observer.disconnect()
})
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
  margin: 0 auto;
  animation: fadeIn 0.5s ease-out;
}

.insight-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin-bottom: 3rem;
}
</style>
