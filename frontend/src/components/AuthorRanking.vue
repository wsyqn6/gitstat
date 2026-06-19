<template>
  <div class="insight-card card">
    <div class="insight-header">
      <h3>{{ t('dashboard.authorRank') }} <span class="range-hint">{{ weekRange }}</span></h3>
    </div>
    <div v-if="loading" class="skeleton-rank">
      <div class="skeleton-rank-row" v-for="i in 5" :key="i">
        <Skeleton circle />
        <Skeleton w="60" />
        <Skeleton w="20" />
      </div>
    </div>
    <template v-else-if="authorRankWithRepos.length > 0">
      <div class="rank-list">
        <template v-for="(author, index) in authorRankWithRepos.slice(0, 5)" :key="author.email">
          <div class="rank-row rank-row-compact">
            <div class="rank-num" :class="{ gold: index === 0, silver: index === 1, bronze: index === 2 }">
              {{ index + 1 }}
            </div>
            <div class="rank-info">
              <span class="rank-name">{{ author.author }}</span>
              <span v-if="author.isMe" class="me-badge-small">{{ t('dashboard.me') }}</span>
            </div>
            <div class="rank-stats">
              <span class="rank-commits">{{ author.commits }}</span>
            </div>
          </div>
          <div class="rank-repo-dist">
            <span v-for="r in author.repos" :key="r.name" class="repo-tag" :style="{ borderColor: r.color, color: r.color }">
              {{ r.name }} {{ r.commits }}
            </span>
          </div>
        </template>
      </div>
    </template>
    <div v-else class="insight-empty">{{ t('analytics.noData') }}</div>
  </div>
</template>

<script setup>
import Skeleton from './Skeleton.vue'
import { useI18n } from '../i18n'

const { t } = useI18n()

defineProps({
  authorRankWithRepos: { type: Array, required: true },
  weekRange: { type: String, required: true },
  loading: { type: Boolean, required: true }
})
</script>

<style scoped>
.skeleton-rank {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 4px 0;
}

.skeleton-rank-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
}

.skeleton-rank-row .skeleton-line {
  height: 14px;
}

.skeleton-rank-row .skeleton-line.w60 {
  flex: 1;
}

.rank-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.rank-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  transition: all 0.3s;
  border-bottom: 1px solid rgba(0, 212, 255, 0.08);
}

.rank-row:last-child {
  border-bottom: none;
}

.rank-row:hover {
  background: var(--bg-row-hover);
}

.rank-num {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-family: var(--font-display);
  font-size: 0.8rem;
  font-weight: 700;
  background: var(--bg-btn-hover);
  color: var(--color-nav-link);
  flex-shrink: 0;
}

.rank-num.gold {
  background: rgba(var(--color-gold-rgb), 0.2);
  color: var(--color-gold);
  box-shadow: 0 0 12px rgba(255, 215, 0, 0.3);
}

.rank-num.silver {
  background: rgba(var(--color-silver-rgb), 0.15);
  color: var(--color-silver);
}

.rank-num.bronze {
  background: rgba(var(--color-bronze-rgb), 0.15);
  color: var(--color-bronze);
}

.rank-info {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  min-width: 0;
}

.rank-name {
  font-weight: 600;
  color: var(--color-text-primary);
  font-size: 0.95rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.me-badge-small {
  background: var(--bg-me-badge);
  border: 1px solid var(--border-me-badge);
  padding: 0.1rem 0.4rem;
  border-radius: 4px;
  font-size: 0.6rem;
  font-weight: 700;
  color: var(--color-primary);
  font-family: var(--font-display);
  letter-spacing: 1px;
  flex-shrink: 0;
}

.rank-stats {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-shrink: 0;
}

.rank-commits {
  font-family: var(--font-display);
  font-size: 1rem;
  font-weight: 700;
  color: var(--color-primary);
}

.rank-row-compact {
  padding-bottom: 0.25rem;
}

.rank-repo-dist {
  display: flex;
  gap: 0.4rem;
  flex-wrap: wrap;
  padding: 0 1rem 0.6rem 3rem;
  margin-top: -0.25rem;
}

.repo-tag {
  font-size: 0.7rem;
  font-family: var(--font-mono);
  padding: 0.1rem 0.5rem;
  border-radius: 8px;
  border: 1px solid;
  background: var(--bg-detail-msg);
  font-weight: 600;
}
</style>
