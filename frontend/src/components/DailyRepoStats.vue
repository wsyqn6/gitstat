<template>
  <div class="daily-stats-section">
    <div class="section-header">
      <h3>{{ t('dashboard.todayDetails') }}</h3>
    </div>

    <div v-if="loading" class="skeleton-daily">
      <div v-for="i in 2" :key="i" class="repo-daily-card glass card">
        <div class="repo-daily-header" style="display:flex;align-items:center;gap:1rem;">
          <Skeleton w="50" />
          <Skeleton w="20" />
        </div>
        <div class="authors-table">
          <div v-for="j in 2" :key="j" class="table-row">
            <Skeleton w="35" />
            <Skeleton w="15" />
            <Skeleton w="35" />
          </div>
        </div>
      </div>
    </div>
    <template v-else-if="dailyStats && dailyStats.length > 0">
      <div v-for="repo in dailyStats" :key="repo.repoPath" class="repo-daily-card glass card">
        <div class="repo-daily-header">
          <div class="repo-info">
            <h4>{{ repo.repoName }}</h4>
            <span class="branch-badge">{{ repo.currentBranch }}</span>
            <span class="last-commit">{{ t('dashboard.lastCommit') }}: {{ repo.lastCommitTime }}</span>
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
            <div class="col-commits">
              <span class="value-with-change">
                <span>{{ author.commits }}</span>
                <ChangeBadge v-if="authorChangeMap[author.email]" :change="authorChangeMap[author.email].commits" compact />
              </span>
            </div>
            <div class="col-changes">
              <span class="value-with-change additions">
                <span>+{{ author.additions }}</span>
                <ChangeBadge v-if="authorChangeMap[author.email]" :change="authorChangeMap[author.email].additions" compact />
              </span>
              <span class="value-with-change deletions">
                <span>-{{ author.deletions }}</span>
                <ChangeBadge v-if="authorChangeMap[author.email]" :change="authorChangeMap[author.email].deletions" compact />
              </span>
              <span class="file-count" :title="t('dashboard.filesChanged')">{{ author.filesChanged }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>
    <div v-else class="glass card card-empty">
      <div class="card-empty-text">{{ t('dashboard.noDailyStats') }}</div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import Skeleton from './Skeleton.vue'
import ChangeBadge from './ChangeBadge.vue'
import { useI18n } from '../i18n'

const { t } = useI18n()

const props = defineProps({
  dailyStats: { type: Array, required: true },
  loading: { type: Boolean, required: true },
  comparisonAuthors: { type: Array, default: () => [] }
})

const authorChangeMap = computed(() => {
  const map = {}
  for (const a of props.comparisonAuthors) {
    map[a.email] = a.change
  }
  return map
})
</script>

<style scoped>
.daily-stats-section {
  margin-top: 2rem;
  margin-bottom: 3rem;
}

.repo-daily-card {
  margin-bottom: 2rem;
  overflow: hidden;
}

.repo-daily-header {
  margin: -1.5rem -1.5rem 1rem;
  padding: 1rem 2rem;
  background: var(--bg-row-hover);
  border-bottom: 1px solid var(--border-insight-card);
}

.repo-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-wrap: wrap;
}

.repo-daily-header h4 {
  font-family: var(--font-display);
  font-size: 1.1rem;
  color: var(--color-primary);
  margin: 0;
  letter-spacing: 1px;
}

.branch-badge {
  background: var(--bg-badge-add-intense);
  border: 1px solid var(--border-badge-add-intense);
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.8rem;
  color: var(--color-green);
  font-weight: 600;
  font-family: var(--font-mono);
}

.last-commit {
  font-size: 0.85rem;
  color: var(--color-nav-link);
  font-family: var(--font-mono);
}

.table-header,
.table-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1.5fr;
  padding: 1rem 2rem;
  align-items: center;
}

.table-header {
  background: var(--bg-row-hover);
  border-bottom: 1px solid var(--border-insight-card);
  font-family: var(--font-display);
  font-size: 0.8rem;
  color: var(--color-nav-link);
  letter-spacing: 1px;
  text-transform: uppercase;
}

.table-row {
  border-bottom: 1px solid var(--bg-btn-hover);
  transition: all 0.3s;
}

.table-row:last-child {
  border-bottom: none;
}

.table-row:hover {
  background: var(--bg-row-hover);
}

.col-author {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.author-name {
  font-weight: 600;
  color: var(--color-text-primary);
}

.me-badge {
  background: var(--bg-me-badge);
  border: 1px solid var(--border-me-badge);
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-size: 0.7rem;
  font-weight: 700;
  color: var(--color-primary);
  font-family: var(--font-display);
  letter-spacing: 1px;
  cursor: help;
}

.col-commits {
  font-family: var(--font-display);
  font-size: 1.1rem;
  font-weight: 700;
  color: var(--color-primary);
  text-align: center;
}

.value-with-change {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
}

.col-changes {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
  align-items: center;
}

.additions {
  color: var(--color-green);
  font-weight: 600;
}

.deletions {
  color: var(--color-pink);
  font-weight: 600;
}

.file-count {
  color: var(--color-nav-link);
  font-weight: 600;
}

.file-count::before {
  content: '|';
  margin-right: 0.35rem;
  color: var(--border-insight-card);
}

.skeleton-daily {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.skeleton-daily .repo-daily-header .skeleton-line {
  height: 18px;
}

.skeleton-daily .authors-table .table-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1.5fr;
  padding: 1rem 2rem;
  align-items: center;
}
</style>
