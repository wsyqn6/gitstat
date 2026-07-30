<template>
  <div class="file-ranking glass card">
    <div class="ranking-header">
      <div class="title-section">
        <h3>{{ t('trends.fileRanking.title') }}</h3>
        <p class="subtitle">{{ t('trends.fileRanking.subtitle') }}</p>
      </div>
      <div class="sort-bar">
        <div class="sort-inner">
          <div class="sort-slider" :class="sortMode"></div>
          <button
            @click="sortMode = 'commits'"
            class="sort-btn"
            :class="{ active: sortMode === 'commits' }"
          >
            {{ t('trends.fileRanking.byCommits') }}
          </button>
          <button
            @click="sortMode = 'changes'"
            class="sort-btn"
            :class="{ active: sortMode === 'changes' }"
          >
            {{ t('trends.fileRanking.byChanges') }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="skeleton-list">
      <div class="skeleton-row" v-for="i in 6" :key="i">
        <Skeleton circle />
        <div class="skeleton-bars">
          <Skeleton w="55" />
          <Skeleton w="30" />
        </div>
      </div>
    </div>

    <template v-else-if="sortedData.length > 0">
      <div class="rank-list">
        <div
          v-for="(item, index) in sortedData"
          :key="item.filePath"
          class="rank-row"
          :style="{ animationDelay: `${index * 30}ms` }"
        >
          <div class="rank-badge" :class="medalClass(index)">
            {{ index + 1 }}
          </div>
          <div class="row-body">
            <div class="file-path-row">
              <span class="lang-icon"><LangIcon :filePath="item.filePath" /></span>
              <span class="file-path" :title="item.filePath">{{ truncatePath(item.filePath) }}</span>
            </div>
            <div class="row-stats">
              <div class="stat-item stat-commits">
                <span class="stat-value">{{ item.commits }}</span>
                <span class="stat-label">{{ t('trends.fileRanking.commits') }}</span>
              </div>
              <div class="stat-item stat-add">
                <span class="stat-value">+{{ item.additions }}</span>
              </div>
              <div class="stat-item stat-del">
                <span class="stat-value">-{{ item.deletions }}</span>
              </div>
              <div class="stat-item stat-net" :class="netClass(item.netChange)">
                <span class="stat-value">{{ netPrefix(item.netChange) }}{{ Math.abs(item.netChange) }}</span>
                <span class="stat-label">{{ t('trends.fileRanking.net') }}</span>
              </div>
            </div>
            <div class="bar-track">
              <div
                class="bar-fill"
                :style="{ width: barWidth(item) }"
              ></div>
            </div>
          </div>
        </div>
      </div>
      <button
        v-if="hasMore"
        @click="$emit('loadMore')"
        class="btn load-more-btn"
      >
        {{ t('repo.loadMore') }}
      </button>
    </template>

    <div v-else class="card-empty">
      <div class="card-empty-text">{{ t('trends.fileRanking.noData') }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '../i18n'
import Skeleton from './Skeleton.vue'
import LangIcon from './LangIcon.vue'

const { t } = useI18n()

const sortMode = ref('commits')

const props = defineProps({
  data: { type: Array, required: true },
  loading: Boolean,
  hasMore: Boolean
})

defineEmits(['loadMore'])

const sortedData = computed(() => {
  const items = props.data || []
  if (sortMode.value === 'changes') {
    return [...items].sort((a, b) => (b.additions + b.deletions) - (a.additions + a.deletions))
  }
  return items
})

const maxVal = computed(() => {
  if (sortedData.value.length === 0) return 1
  if (sortMode.value === 'changes') {
    return Math.max(...sortedData.value.map(d => d.additions + d.deletions), 1)
  }
  return Math.max(...sortedData.value.map(d => d.commits), 1)
})

function barWidth(item) {
  const val = sortMode.value === 'changes' ? item.additions + item.deletions : item.commits
  const pct = (val / maxVal.value) * 100
  return `${Math.max(pct, 2)}%`
}

function truncatePath(path) {
  if (path.length <= 55) return path
  const parts = path.split('/')
  if (parts.length <= 2) return path
  const head = parts[0]
  const tail = parts.slice(-2).join('/')
  return `${head}/.../${tail}`
}

function medalClass(index) {
  if (index === 0) return 'gold'
  if (index === 1) return 'silver'
  if (index === 2) return 'bronze'
  return ''
}

function netClass(val) {
  if (val > 0) return 'pos'
  if (val < 0) return 'neg'
  return ''
}

function netPrefix(val) {
  if (val > 0) return '+'
  if (val < 0) return '-'
  return ''
}
</script>

<style scoped>
.file-ranking {
  content-visibility: auto;
  contain-intrinsic-size: 500px;
  animation: fadeIn 0.4s ease;
}

.ranking-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.title-section h3 {
  font-family: var(--font-display);
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-text-secondary);
  letter-spacing: 2px;
  text-transform: uppercase;
  margin: 0 0 0.15rem 0;
}

.subtitle {
  font-size: 0.8rem;
  color: var(--color-text-muted);
  letter-spacing: 1px;
  margin: 0;
}

.sort-bar {
  flex-shrink: 0;
}

.sort-inner {
  display: inline-flex;
  background: var(--glass-btn-bg);
  backdrop-filter: blur(var(--glass-blur));
  border: 1px solid var(--glass-btn-border);
  border-radius: calc(var(--radius-btn) + 2px);
  padding: 3px;
  gap: 2px;
  position: relative;
}

.sort-slider {
  position: absolute;
  top: 3px;
  left: 3px;
  width: calc(50% - 4px);
  height: calc(100% - 6px);
  border-radius: calc(var(--radius-btn) - 2px);
  background: var(--glass-btn-hover-bg);
  backdrop-filter: blur(var(--glass-blur));
  box-shadow: var(--glass-btn-shadow), var(--glass-btn-inner);
  transition: transform 0.35s cubic-bezier(0.34, 1.56, 0.64, 1);
  z-index: 0;
  pointer-events: none;
}

.sort-slider.changes {
  transform: translateX(calc(100% + 2px));
}

.sort-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0.4rem 1rem;
  border: none;
  background: transparent;
  color: var(--color-text-muted);
  font-family: var(--font-body);
  font-size: 0.8rem;
  font-weight: 500;
  cursor: pointer;
  transition: color 0.25s ease;
  position: relative;
  z-index: 1;
  white-space: nowrap;
}

.sort-btn.active {
  color: var(--glass-btn-color);
}

.sort-btn:hover:not(.active) {
  color: var(--color-text-secondary);
}

/* === Skeleton === */
.skeleton-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.skeleton-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 0;
}

.skeleton-bars {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

/* === Rank list === */
.rank-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.rank-row {
  display: flex;
  align-items: stretch;
  gap: 12px;
  padding: 10px 14px;
  border-radius: 10px;
  transition: background 0.2s ease;
  animation: rowFadeIn 0.3s ease both;
  cursor: default;
}

.rank-row:hover {
  background: var(--bg-row-hover);
}

.rank-badge {
  width: 28px;
  min-height: 52px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  font-family: var(--font-display);
  font-size: 0.75rem;
  font-weight: 700;
  background: var(--bg-btn-hover);
  color: var(--color-nav-link);
  flex-shrink: 0;
  align-self: center;
}

.rank-badge.gold {
  background: rgba(var(--color-gold-rgb), 0.2);
  color: var(--color-gold);
  box-shadow: 0 0 12px rgba(255, 215, 0, 0.25);
}

.rank-badge.silver {
  background: rgba(var(--color-silver-rgb), 0.15);
  color: var(--color-silver);
}

.rank-badge.bronze {
  background: rgba(var(--color-bronze-rgb), 0.15);
  color: var(--color-bronze);
}

.row-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.file-path-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.lang-icon {
  display: flex;
  align-items: center;
  flex-shrink: 0;
  line-height: 0;
}

.file-path {
  font-family: var(--font-mono, 'Cascadia Code', 'Fira Code', monospace);
  font-size: 0.85rem;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  letter-spacing: 0.3px;
}

.row-stats {
  display: flex;
  align-items: center;
  gap: 14px;
}

.stat-item {
  display: flex;
  align-items: baseline;
  gap: 3px;
}

.stat-value {
  font-family: var(--font-display);
  font-size: 0.85rem;
  font-weight: 600;
}

.stat-label {
  font-family: var(--font-body);
  font-size: 0.6rem;
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.stat-commits .stat-value {
  color: var(--color-primary);
}

.stat-add .stat-value {
  color: var(--color-green);
  font-size: 0.8rem;
}

.stat-del .stat-value {
  color: var(--color-red);
  font-size: 0.8rem;
}

.stat-net .stat-value {
  font-size: 0.8rem;
}

.stat-net.pos .stat-value {
  color: var(--color-green);
}

.stat-net.neg .stat-value {
  color: var(--color-red);
}

.bar-track {
  height: 3px;
  background: var(--bg-btn-hover);
  border-radius: 2px;
  overflow: hidden;
  margin-top: 2px;
}

.bar-fill {
  height: 100%;
  border-radius: 2px;
  background: linear-gradient(90deg, var(--color-primary), rgba(var(--color-primary-rgb), 0.3));
  transition: width 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.load-more-btn {
  margin-top: 1rem;
  width: 100%;
  text-align: center;
}
.load-more-btn:focus { outline: none; }

@media (max-width: 768px) {
  .ranking-header {
    flex-direction: column;
  }

  .row-stats {
    flex-wrap: wrap;
    gap: 10px;
  }

  .file-path {
    font-size: 0.8rem;
  }
}
</style>
