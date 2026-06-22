<template>
  <div class="section glass card">
    <h3 class="section-title">{{ t('repo.recentCommits') }}</h3>

    <div v-if="commitsLoading && !commitsLoaded" class="commit-skeleton">
      <Skeleton w="25" h="14" mb="1rem" />
      <div v-for="i in 5" :key="i" class="skel-commit-row">
        <Skeleton w="70" />
        <Skeleton w="85" />
        <Skeleton w="60" />
        <Skeleton w="55" />
        <Skeleton w="50" />
      </div>
    </div>
    <div v-else-if="!commitsLoaded" class="cta-body">
      <div class="cta-icon">⎔</div>
      <p class="cta-text">{{ t('repo.commitsCta') }}</p>
      <button class="cta-btn btn" :disabled="commitsLoading" @click="emit('loadCommits')">
        <span v-if="commitsLoading" class="spinner"></span>
        <span v-else>{{ t('repo.loadCommits') }}</span>
      </button>
    </div>

    <div v-else-if="commits.length === 0" class="empty-body">
      <p>{{ t('repo.noCommits') }}</p>
    </div>

    <div v-else class="commit-list">
      <div v-for="c in commits" :key="c.hash" class="commit-item">
        <div class="commit-main" @click="toggleCommit(c.hash)">
          <span class="commit-hash">{{ c.hash.slice(0, 7) }}</span>
          <span class="commit-msg">{{ c.message.split('\n')[0] }}</span>
          <span class="commit-author">{{ c.author }}</span>
          <span class="commit-time">{{ formatTimeAgo(c.date) }}</span>
          <span class="commit-changes">
            <span v-if="c.additions > 0" class="add">+{{ c.additions }}</span>
            <span v-if="c.deletions > 0" class="del">-{{ c.deletions }}</span>
          </span>
        </div>
        <div v-if="expandedCommit === c.hash" class="commit-body">
          <pre>{{ c.message }}</pre>
        </div>
      </div>

      <div class="load-more-wrap">
        <button
          v-if="hasMore"
          class="load-more-btn btn"
          :disabled="commitsLoading"
          @click="emit('loadMore')"
        >
          <span v-if="commitsLoading" class="spinner"></span>
          <span v-else>{{ t('repo.loadMore') }}</span>
        </button>
        <span v-else class="all-loaded">{{ t('repo.allLoaded') }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Skeleton from './Skeleton.vue'
import { useI18n } from '../i18n'

const { t } = useI18n()

defineProps({
  commits: { type: Array, default: () => [] },
  commitsLoading: Boolean,
  commitsLoaded: Boolean,
  hasMore: Boolean
})

const emit = defineEmits(['loadCommits', 'loadMore'])

const expandedCommit = ref(null)

function formatTimeAgo(s) {
  if (!s) return '-'
  const d = new Date(s)
  const now = new Date()
  const pad = n => String(n).padStart(2, '0')
  const hm = pad(d.getHours()) + ':' + pad(d.getMinutes())
  const diffSec = (now - d.getTime()) / 1000

  const startOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  if (d >= startOfToday) {
    if (diffSec < 60) return Math.round(diffSec) + 's'
    if (diffSec < 3600) return Math.round(diffSec / 60) + 'm'
    return Math.round(diffSec / 3600) + 'h'
  }

  const startOfYesterday = new Date(startOfToday)
  startOfYesterday.setDate(startOfYesterday.getDate() - 1)
  if (d >= startOfYesterday) return t('repo.yesterday') + ' ' + hm

  const y = d.getFullYear()
  const mo = pad(d.getMonth() + 1)
  const dd = pad(d.getDate())
  return `${y}-${mo}-${dd} ${hm}`
}

function toggleCommit(hash) {
  expandedCommit.value = expandedCommit.value === hash ? null : hash
}
</script>

<style scoped>
.section { padding: 1.5rem; }
.section-title {
  font-family: var(--font-display);
  font-size: 0.95rem;
  color: var(--color-primary);
  letter-spacing: 1px;
  margin: 0 0 1rem 0;
}

.cta-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 2rem 1rem;
  text-align: center;
  gap: 0.8rem;
}
.cta-icon { font-size: 2.5rem; opacity: 0.4; color: var(--color-primary); }
.cta-text { color: var(--color-text-muted); font-size: 0.9rem; font-family: var(--font-body); }
.cta-btn {
  padding: 0.6rem 1.6rem;
  font-size: 0.85rem;
}

.empty-body { padding: 2rem; text-align: center; color: var(--color-text-muted); }

.commit-list { width: 100%; }
.commit-main {
  display: grid;
  grid-template-columns: 0.5fr 4fr 0.7fr 0.9fr 0.7fr;
  padding: 0.55rem 0.5rem;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.85rem;
  color: var(--color-text-primary);
  border-bottom: 1px solid rgba(var(--color-primary-rgb), 0.06);
  cursor: pointer;
  transition: background 0.2s;
}
.commit-main:hover { background: rgba(var(--color-primary-rgb), 0.03); }
.commit-hash { font-family: var(--font-mono); color: var(--color-nav-link); }
.commit-msg { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.commit-author { font-size: 0.8rem; }
.commit-time { font-family: var(--font-mono); color: var(--color-text-muted); font-size: 0.8rem; }
.commit-changes { text-align: right; }
.commit-changes .add { margin-right: 0.3rem; }
.add { color: var(--color-green); }
.del { color: var(--color-red); }
.commit-body pre {
  margin: 0; padding: 0.6rem 0.5rem;
  background: var(--bg-detail-msg);
  font-family: var(--font-mono);
  font-size: 0.8rem;
  color: var(--color-nav-link);
  border-bottom: 1px solid rgba(var(--color-primary-rgb), 0.06);
  white-space: pre-wrap;
}

.load-more-wrap {
  display: flex;
  justify-content: center;
  padding: 1rem 0;
}
.load-more-btn {
  padding: 0.5rem 2rem;
  font-size: 0.85rem;
}
.all-loaded { color: var(--color-all-loaded); font-size: 0.8rem; }

.commit-skeleton { padding: 0.25rem 0; }
.skel-commit-row {
  display: grid;
  grid-template-columns: 0.5fr 4fr 0.7fr 0.9fr 0.7fr;
  padding: 0.55rem 0.5rem;
  align-items: center;
  gap: 0.3rem;
  border-bottom: 1px solid rgba(var(--color-primary-rgb), 0.06);
}

.spinner {
  display: inline-block;
  width: 16px; height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  vertical-align: middle;
}
</style>
