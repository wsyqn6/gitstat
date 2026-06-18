<template>
  <div class="section card">
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
      <button class="cta-btn" :disabled="commitsLoading" @click="emit('loadCommits')">
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
          class="load-more-btn"
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
  const now = Date.now()
  const diff = (now - d.getTime()) / 1000
  if (diff < 60) return Math.round(diff) + 's'
  if (diff < 3600) return Math.round(diff / 60) + 'm'
  if (diff < 86400) return Math.round(diff / 3600) + 'h'
  if (diff < 604800) return Math.round(diff / 86400) + 'd'
  return s.slice(0, 10)
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
  color: #00d4ff;
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
.cta-icon { font-size: 2.5rem; opacity: 0.4; color: #00d4ff; }
.cta-text { color: #64748b; font-size: 0.9rem; font-family: var(--font-body); }
.cta-btn {
  padding: 0.6rem 1.6rem;
  border: 1px solid rgba(0, 212, 255, 0.3);
  border-radius: 8px;
  background: rgba(0, 212, 255, 0.08);
  color: #00d4ff;
  font-family: var(--font-display);
  font-size: 0.85rem;
  letter-spacing: 1px;
  cursor: pointer;
  transition: all 0.3s;
}
.cta-btn:hover:not(:disabled) { background: rgba(0, 212, 255, 0.15); border-color: #00d4ff; }
.cta-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.empty-body { padding: 2rem; text-align: center; color: #64748b; }

.commit-list { width: 100%; }
.commit-main {
  display: grid;
  grid-template-columns: 1fr 2.5fr 1fr 1fr 1.2fr;
  padding: 0.55rem 0.5rem;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.85rem;
  color: #e0e6ff;
  border-bottom: 1px solid rgba(0, 212, 255, 0.06);
  cursor: pointer;
  transition: background 0.2s;
}
.commit-main:hover { background: rgba(0, 212, 255, 0.03); }
.commit-hash { font-family: 'Rajdhani', monospace; color: #a0aec0; }
.commit-msg { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.commit-author { font-size: 0.8rem; }
.commit-time { font-family: 'Rajdhani', monospace; color: #64748b; font-size: 0.8rem; }
.commit-changes { text-align: right; }
.commit-changes .add { margin-right: 0.3rem; }
.add { color: #00ff88; }
.del { color: #ff6b6b; }
.commit-body pre {
  margin: 0; padding: 0.6rem 0.5rem;
  background: rgba(0, 0, 0, 0.2);
  font-family: 'Rajdhani', monospace;
  font-size: 0.8rem;
  color: #a0aec0;
  border-bottom: 1px solid rgba(0, 212, 255, 0.06);
  white-space: pre-wrap;
}

.load-more-wrap {
  display: flex;
  justify-content: center;
  padding: 1rem 0;
}
.load-more-btn {
  padding: 0.5rem 2rem;
  border: 1px solid rgba(0, 212, 255, 0.25);
  border-radius: 6px;
  background: rgba(0, 212, 255, 0.05);
  color: #00d4ff;
  font-family: var(--font-display);
  font-size: 0.85rem;
  letter-spacing: 1px;
  cursor: pointer;
  transition: all 0.3s;
}
.load-more-btn:hover:not(:disabled) { background: rgba(0, 212, 255, 0.12); border-color: #00d4ff; }
.load-more-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.all-loaded { color: #475569; font-size: 0.8rem; }

.commit-skeleton { padding: 0.25rem 0; }
.skel-commit-row {
  display: grid;
  grid-template-columns: 1fr 2.5fr 1fr 1fr 1.2fr;
  padding: 0.55rem 0.5rem;
  align-items: center;
  gap: 0.3rem;
  border-bottom: 1px solid rgba(0, 212, 255, 0.06);
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
