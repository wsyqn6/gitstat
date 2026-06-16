<template>
  <template v-if="loading || (commits && commits.length > 0)">
    <div class="section card">
      <h3 class="section-title">{{ t('repo.recentCommits') }}</h3>
      <div v-if="loading" class="commit-list">
        <div v-for="i in 5" :key="i" class="commit-main">
          <span class="skeleton-line w15" style="height:12px"></span>
          <span class="skeleton-line w45" style="height:12px"></span>
          <span class="skeleton-line w12" style="height:12px"></span>
          <span class="skeleton-line w12" style="height:12px"></span>
          <span class="skeleton-line w18" style="height:12px"></span>
        </div>
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
      </div>
    </div>
  </template>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from '../i18n'

const { t } = useI18n()

defineProps({
  commits: { type: Array, default: () => [] },
  loading: Boolean
})

const expandedCommit = ref(null)

function formatTimeAgo(s) {
  if (!s) return '-'
  const d = new Date(s)
  const now = Date.now()
  const diff = (now - d.getTime()) / 1000
  if (diff < 60) return Math.round(diff) + 's'
  if (diff < 3600) return Math.round(diff / 60) + 'm'
  if (diff < 86400) return Math.round(diff / 3600) + 'h'
  if (diff < 2592000) return Math.round(diff / 86400) + 'd'
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
</style>
