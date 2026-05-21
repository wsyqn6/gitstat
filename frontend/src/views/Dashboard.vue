<template>
  <div class="dashboard">
    <!-- 今日统计概览 -->
      <div v-if="state.overviewStats" class="stats-grid">
        <StatCard 
          :value="todayCommits" 
          :label="t('dashboard.todayCommits')"
          icon="◈"
          color="#00d4ff"
        />
        <StatCard 
          :value="todayAdditions" 
          :label="t('dashboard.todayAdditions')"
          icon="↑"
          color="#00ff88"
        />
        <StatCard 
          :value="todayDeletions" 
          :label="t('dashboard.todayDeletions')"
          icon="↓"
          color="#ff6b9d"
        />
        <StatCard 
          :value="state.overviewStats.activeAuthors" 
          :label="t('dashboard.activeAuthors')"
          icon="◉"
          color="#ffd700"
        />
      </div>
      
      <!-- 分仓库分人统计 -->
      <div v-if="state.dailyStats && state.dailyStats.length > 0" class="daily-stats-section">
        <div class="section-header">
          <h3>{{ t('dashboard.todayDetails') }}</h3>
        </div>
        
        <div v-for="repo in state.dailyStats" :key="repo.repoPath" class="repo-daily-card card">
          <div class="repo-daily-header">
            <div class="repo-info">
              <h4>{{ repo.repoName }}</h4>
              <div class="repo-meta">
                <span class="branch-badge">{{ repo.currentBranch }}</span>
                <span class="last-commit">{{ t('dashboard.lastCommit') }}: {{ repo.lastCommitTime }}</span>
              </div>
            </div>
          </div>
          
          <div class="authors-table">
            <div class="table-header">
              <div class="col-author">{{ t('dashboard.author') }}</div>
              <div class="col-commits">{{ t('dashboard.commits') }}</div>
              <div class="col-changes">{{ t('dashboard.changes') }}</div>
            </div>
            <div 
              v-for="author in repo.authors" 
              :key="author.email" 
              class="table-row"
            >
              <div class="col-author">
                <span class="author-name">{{ author.author }}</span>
                <span v-if="author.isMe" class="me-badge" :title="t('dashboard.me')">{{ t('dashboard.me') }}</span>
              </div>
              <div class="col-commits">{{ author.commits }}</div>
              <div class="col-changes">
                <span class="additions">+{{ author.additions }}</span>
                <span class="deletions">-{{ author.deletions }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '../i18n'
import { state, refreshDailyStats, refreshStats } from '../stores/data'
import StatCard from '../components/StatCard.vue'

const { t } = useI18n()

const today = new Date().toISOString().split('T')[0]

const todayCommits = computed(() => {
  if (!state.overviewStats) return 0
  return state.overviewStats.totalCommits
})

const todayAdditions = computed(() => {
  if (!state.overviewStats) return 0
  return state.overviewStats.totalAdditions
})

const todayDeletions = computed(() => {
  if (!state.overviewStats) return 0
  return state.overviewStats.totalDeletions
})

// 首次加载直接请求数据，触发后端懒加载
onMounted(async () => {
  await refreshStats()
})
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
  margin: 0 auto;
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.header-actions {
  margin-bottom: 2rem;
  display: flex;
  justify-content: flex-end;
}

.switch-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.btn-icon {
  font-size: 1.2rem;
  transition: transform 0.3s;
}

.switch-btn:hover .btn-icon {
  transform: rotate(180deg);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  margin-bottom: 3rem;
}

.daily-stats-section {
  margin-top: 2rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.section-header h3 {
  font-family: 'Orbitron', sans-serif;
  font-size: 1.5rem;
  color: #00d4ff;
  letter-spacing: 2px;
  text-transform: uppercase;
  margin: 0;
}

.repo-daily-card {
  margin-bottom: 2rem;
  padding: 0;
  overflow: hidden;
}

.repo-daily-header {
  padding: 1.5rem 2rem;
  background: rgba(0, 212, 255, 0.05);
  border-bottom: 1px solid rgba(0, 212, 255, 0.2);
}

.repo-daily-header h4 {
  font-family: 'Orbitron', sans-serif;
  font-size: 1.2rem;
  color: #00d4ff;
  margin: 0 0 0.75rem 0;
  letter-spacing: 1px;
}

.repo-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.branch-badge {
  background: rgba(0, 255, 136, 0.15);
  border: 1px solid rgba(0, 255, 136, 0.4);
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.8rem;
  color: #00ff88;
  font-weight: 600;
  font-family: 'Rajdhani', monospace;
}

.last-commit {
  font-size: 0.85rem;
  color: #a0aec0;
  font-family: 'Rajdhani', monospace;
}

.authors-table {
  padding: 1rem 0;
}

.table-header,
.table-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1.5fr;
  padding: 1rem 2rem;
  align-items: center;
}

.table-header {
  background: rgba(0, 212, 255, 0.05);
  border-bottom: 1px solid rgba(0, 212, 255, 0.2);
  font-family: 'Orbitron', sans-serif;
  font-size: 0.8rem;
  color: #a0aec0;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.table-row {
  border-bottom: 1px solid rgba(0, 212, 255, 0.1);
  transition: all 0.3s;
}

.table-row:last-child {
  border-bottom: none;
}

.table-row:hover {
  background: rgba(0, 212, 255, 0.05);
}

.col-author {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.author-name {
  font-weight: 600;
  color: #e0e6ff;
}

.me-badge {
  background: rgba(0, 212, 255, 0.15);
  border: 1px solid rgba(0, 212, 255, 0.4);
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-size: 0.7rem;
  font-weight: 700;
  color: #00d4ff;
  font-family: 'Orbitron', sans-serif;
  letter-spacing: 1px;
  cursor: help;
}

.col-commits {
  font-family: 'Orbitron', sans-serif;
  font-size: 1.1rem;
  font-weight: 700;
  color: #00d4ff;
  text-align: center;
}

.col-changes {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.additions {
  color: #00ff88;
  font-weight: 600;
}

.deletions {
  color: #ff6b9d;
  font-weight: 600;
}
</style>
