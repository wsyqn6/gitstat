<template>
  <div class="comparison-section">
    <div class="section-header">
      <h3>{{ t('dashboard.repoComparison') }} <span class="range-hint">{{ weekRange }}</span></h3>
    </div>
    <div v-if="loading" class="comparison-table glass card">
      <div v-for="i in 3" :key="i" class="cmp-row">
        <Skeleton w="40" />
        <Skeleton w="15" />
        <Skeleton w="40" />
        <Skeleton w="15" />
        <Skeleton w="15" />
      </div>
    </div>
    <div v-else-if="repoComparison.length > 0" class="comparison-table glass card">
      <div class="cmp-header">
        <div class="cmp-col-name">{{ t('dashboard.repo') }}</div>
        <div class="cmp-col-num">{{ t('dashboard.commits') }}</div>
        <div class="cmp-col-num">{{ t('dashboard.changes') }}</div>
        <div class="cmp-col-num">{{ t('dashboard.activeDays') }}</div>
        <div class="cmp-col-num">{{ t('dashboard.dailyAvg') }}</div>
      </div>
      <div v-for="repo in repoComparison" :key="repo.repoPath" class="cmp-row">
        <div class="cmp-col-name">{{ repo.repoName }}</div>
        <div class="cmp-col-num">{{ repo.commits }}</div>
        <div class="cmp-col-num changes">
          <span class="additions">+{{ repo.additions }}</span>
          <span class="deletions">-{{ repo.deletions }}</span>
          <span class="file-count" :title="t('dashboard.filesChanged')">{{ repo.filesChanged }}</span>
        </div>
        <div class="cmp-col-num">{{ repo.activeDays }}</div>
        <div class="cmp-col-num">{{ repo.avgCommitsPerDay }}</div>
      </div>
    </div>
    <div v-else class="glass card card-empty">
      <div class="card-empty-text">{{ t('dashboard.noComparison') }}</div>
    </div>
  </div>
</template>

<script setup>
import Skeleton from './Skeleton.vue'
import { useI18n } from '../i18n'

const { t } = useI18n()

defineProps({
  repoComparison: { type: Array, required: true },
  weekRange: { type: String, required: true },
  loading: { type: Boolean, required: true }
})
</script>

<style scoped>
.comparison-section {
  margin-bottom: 3rem;
}

.comparison-table {
  overflow: hidden;
}

.cmp-header,
.cmp-row {
  display: grid;
  grid-template-columns: 2fr 1fr 2fr 1fr 1fr;
  padding: 1rem 2rem;
  align-items: center;
}

.cmp-header {
  background: var(--bg-row-hover);
  border-bottom: 1px solid var(--border-insight-card);
  font-family: var(--font-display);
  font-size: 0.8rem;
  color: var(--color-nav-link);
  letter-spacing: 1px;
  text-transform: uppercase;
}

.cmp-row {
  border-bottom: 1px solid var(--bg-btn-hover);
  transition: all 0.3s;
}

.cmp-row:last-child {
  border-bottom: none;
}

.cmp-row:hover {
  background: var(--bg-row-hover);
}

.cmp-col-name {
  font-weight: 600;
  color: var(--color-text-primary);
}

.cmp-col-num {
  font-family: var(--font-display);
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--color-primary);
  text-align: center;
}

.cmp-col-num .additions {
  color: var(--color-green);
}

.cmp-col-num .deletions {
  color: var(--color-pink);
}

.cmp-col-num.changes {
  display: flex;
  gap: 0.75rem;
  justify-content: center;
  align-items: center;
}

.cmp-col-num .file-count {
  color: var(--color-nav-link);
  font-weight: 600;
}

.cmp-col-num .file-count::before {
  content: '|';
  margin-right: 0.35rem;
  color: var(--border-insight-card);
}
</style>
