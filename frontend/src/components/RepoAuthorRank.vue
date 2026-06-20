<template>
  <div class="section glass card author-rank">
    <h3 class="section-title">{{ t('repo.authorRank') }}</h3>
    <div v-if="loading" class="author-skeleton">
      <div v-for="i in 5" :key="i" class="skel-row">
        <Skeleton circle h="28" />
        <Skeleton w="60" h="16" radius="4" />
        <Skeleton w="30" h="16" radius="4" style="margin-left:auto" />
      </div>
    </div>
    <div v-else-if="top.length === 0" class="empty-hint">--</div>
    <template v-else>
      <div v-for="(c, i) in top" :key="c.email || c.author" class="author-row">
        <span class="rank-num">{{ i + 1 }}</span>
        <span class="author-avatar" :style="{ background: colors[i] }">{{ c.author[0] }}</span>
        <span class="author-name">{{ c.author }}</span>
        <div class="bar-track">
          <div class="bar-fill" :style="{ width: barWidth(c), background: colors[i] }"></div>
        </div>
        <span class="author-commits">{{ formatNum(c.commitCount) }}</span>
      </div>
      <div class="author-footer">
        <span>{{ t('repo.total') }} {{ formatNum(totalCommits) }}</span>
        <span class="churn">+{{ formatNum(totalAdditions) }} / −{{ formatNum(totalDeletions) }}</span>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from '../i18n'
import Skeleton from './Skeleton.vue'

const { t } = useI18n()

const props = defineProps({
  contributors: { type: Array, default: () => [] },
  loading: Boolean
})

const MAX_DISPLAY = 6

const colors = [
  'var(--color-primary)',
  'var(--color-purple-soft)',
  'var(--color-amber)',
  'var(--color-emerald)',
  'var(--color-cyan)',
  'var(--color-pink)'
]

const top = computed(() => {
  const sorted = [...props.contributors].sort((a, b) => b.commitCount - a.commitCount)
  return sorted.slice(0, MAX_DISPLAY)
})

const maxCommits = computed(() => {
  if (top.value.length === 0) return 0
  return top.value[0].commitCount
})

const totalCommits = computed(() =>
  props.contributors.reduce((s, c) => s + c.commitCount, 0)
)

const totalAdditions = computed(() =>
  props.contributors.reduce((s, c) => s + c.additions, 0)
)

const totalDeletions = computed(() =>
  props.contributors.reduce((s, c) => s + c.deletions, 0)
)

function barWidth(c) {
  if (maxCommits.value === 0) return '0%'
  return (c.commitCount / maxCommits.value * 100).toFixed(1) + '%'
}

function formatNum(n) {
  if (!n) return '0'
  if (n >= 1000) return (n / 1000).toFixed(1).replace(/\.0$/, '') + 'k'
  return String(n)
}
</script>

<style scoped>
.author-rank {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.section-title {
  font-family: var(--font-display);
  font-size: 0.95rem;
  color: var(--color-primary);
  letter-spacing: 1px;
  margin: 0 0 0.5rem 0;
}

.author-skeleton {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}
.skel-row {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.empty-hint {
  color: var(--color-text-muted);
  font-family: var(--font-body);
  font-size: 0.85rem;
  text-align: center;
  padding: 1rem 0;
}

.author-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.25rem 0;
}

.rank-num {
  width: 1rem;
  font-size: 0.7rem;
  font-weight: 700;
  color: var(--color-text-muted);
  text-align: right;
  font-family: var(--font-display);
}

.author-avatar {
  width: 1.5rem;
  height: 1.5rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.65rem;
  font-weight: 700;
  color: white;
  flex-shrink: 0;
}

.author-name {
  font-family: var(--font-body);
  font-size: 0.8rem;
  color: var(--color-text-primary);
  width: 5rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex-shrink: 0;
}

.bar-track {
  flex: 1;
  height: 6px;
  background: var(--bg-cta);
  border-radius: 3px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.5s ease-out;
}

.author-commits {
  font-family: var(--font-display);
  font-size: 0.8rem;
  font-weight: 700;
  color: var(--color-text-primary);
  width: 2.8rem;
  text-align: right;
  flex-shrink: 0;
}

.author-footer {
  margin-top: 0.5rem;
  padding-top: 0.5rem;
  border-top: 1px solid var(--border-card-subtle);
  display: flex;
  justify-content: space-between;
  font-family: var(--font-body);
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.churn {
  font-family: var(--font-mono);
  color: var(--color-nav-link);
}
</style>
