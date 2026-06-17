<template>
  <div class="comparison-section">
    <div class="section-header">
      <h3>{{ t('dashboard.repoComparison') }} <span class="range-hint">{{ weekRange }}</span></h3>
    </div>
    <div v-if="loading" class="comparison-table card">
      <div v-for="i in 3" :key="i" class="cmp-row">
        <div class="skeleton-line w40"></div>
        <div class="skeleton-line w15"></div>
        <div class="skeleton-line w15"></div>
        <div class="skeleton-line w15"></div>
        <div class="skeleton-line w15"></div>
        <div class="skeleton-line w15"></div>
      </div>
    </div>
    <div v-else-if="repoComparison.length > 0" class="comparison-table card">
      <div class="cmp-header">
        <div class="cmp-col-name">{{ t('dashboard.repo') }}</div>
        <div class="cmp-col-num">{{ t('dashboard.commits') }}</div>
        <div class="cmp-col-num">{{ t('analytics.additions') }}</div>
        <div class="cmp-col-num">{{ t('analytics.deletions') }}</div>
        <div class="cmp-col-num">{{ t('dashboard.activeDays') }}</div>
        <div class="cmp-col-num">{{ t('dashboard.dailyAvg') }}</div>
      </div>
      <div v-for="repo in repoComparison" :key="repo.repoPath" class="cmp-row">
        <div class="cmp-col-name">{{ repo.repoName }}</div>
        <div class="cmp-col-num">{{ repo.commits }}</div>
        <div class="cmp-col-num additions">+{{ repo.additions }}</div>
        <div class="cmp-col-num deletions">-{{ repo.deletions }}</div>
        <div class="cmp-col-num">{{ repo.activeDays }}</div>
        <div class="cmp-col-num">{{ repo.avgCommitsPerDay }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
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
  grid-template-columns: 2fr 1fr 1fr 1fr 1fr 1fr;
  padding: 1rem 2rem;
  align-items: center;
}

.cmp-header {
  background: rgba(0, 212, 255, 0.05);
  border-bottom: 1px solid rgba(0, 212, 255, 0.2);
  font-family: 'Orbitron', sans-serif;
  font-size: 0.8rem;
  color: #a0aec0;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.cmp-row {
  border-bottom: 1px solid rgba(0, 212, 255, 0.1);
  transition: all 0.3s;
}

.cmp-row:last-child {
  border-bottom: none;
}

.cmp-row:hover {
  background: rgba(0, 212, 255, 0.05);
}

.cmp-col-name {
  font-weight: 600;
  color: #e0e6ff;
}

.cmp-col-num {
  font-family: 'Orbitron', sans-serif;
  font-size: 0.95rem;
  font-weight: 700;
  color: #00d4ff;
  text-align: center;
}

.cmp-col-num.additions {
  color: #00ff88;
}

.cmp-col-num.deletions {
  color: #ff6b9d;
}
</style>
