<template>
  <div>
    <div class="card-row l1-row">
      <div class="glass info-card">
        <span class="info-icon" style="color:var(--color-primary)">⑂</span>
        <span class="info-value">{{ detail.currentBranch }}</span>
        <span class="info-label">{{ t('repo.currentBranch') }}</span>
      </div>
      <div class="glass info-card clickable" :class="{ expanded: expandedBranch }" @click="toggleBranch">
        <span class="info-icon" style="color:var(--color-purple-soft)">⑂</span>
        <span class="info-value">{{ localBranchCount }}</span>
        <span class="info-label">{{ t('repo.branchCount') }} <span class="expand-icon">{{ expandedBranch ? '▼' : '▶' }}</span></span>
      </div>
      <div class="glass info-card">
        <span class="info-icon" style="color:var(--color-amber)">◈</span>
        <span class="info-value">{{ detail.fileCount ?? '--' }}</span>
        <span class="info-label">{{ t('repo.fileCount') }}</span>
      </div>
      <div class="glass info-card clickable" :class="{ expanded: expandedTags }" @click="toggleTags">
        <span class="info-icon" style="color:var(--color-emerald)">◉</span>
        <span class="info-value">{{ detail.tagCount ?? '--' }}</span>
        <span class="info-label">{{ t('repo.tags') }} <span class="expand-icon">{{ expandedTags ? '▼' : '▶' }}</span></span>
      </div>
    </div>

    <div v-if="expandedBranch" class="expand-panel glass card">
      <h4>{{ t('repo.allBranches') }}</h4>
      <div class="branch-groups">
        <div v-if="detail.branchCount > 0" class="branch-group">
          <span class="branch-group-label">{{ t('repo.localBranch') }}</span>
          <div class="branch-tags">
            <span
              v-for="b in localBranches"
              :key="b"
              class="branch-tag"
              :class="{ current: b === detail.currentBranch }"
            >{{ b }}<span v-if="b === detail.currentBranch" class="current-badge">{{ t('repo.currentBranch') }}</span></span>
          </div>
        </div>
        <div v-if="remoteBranchCount > 0" class="branch-group">
          <span class="branch-group-label">{{ t('repo.remoteBranch') }}</span>
          <div class="branch-tags">
            <span
              v-for="b in detail.remoteBranches || []"
              :key="b"
              class="branch-tag remote"
            >{{ b }}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="expandedTags" class="expand-panel glass card">
      <h4>{{ t('repo.tags') }}</h4>
      <div v-if="tagsLoading && tags.length === 0" class="expand-hint">
        <span class="spinner" style="display:inline-block;width:14px;height:14px;margin-right:0.5rem;vertical-align:middle"></span>{{ t('repo.statsLoading') }}
      </div>
      <div v-else-if="tags.length > 0" class="tag-cloud">
        <span v-for="t in tags" :key="t" class="tag-item">{{ t }}</span>
        <div v-if="tagsHasMore" class="load-more-tags">
          <button class="btn btn-sm" :disabled="tagsLoading" @click="emit('load-more-tags')">
            <span v-if="tagsLoading" class="spinner" style="display:inline-block;width:12px;height:12px;margin-right:0.3rem;vertical-align:middle"></span>
            {{ t('repo.loadMore') }}
          </button>
        </div>
      </div>
      <p v-else class="expand-hint">--</p>
    </div>

    <div class="card-row l2-row">
      <div class="glass info-card dim" :class="{ 'has-data': statsLoaded && detail.repoSize > 0 }">
        <span class="info-icon" style="color:var(--color-cyan)">◆</span>
        <span class="info-value">
          <template v-if="statsLoaded">{{ detail.repoSize > 0 ? formatBytes(detail.repoSize) : '--' }}</template>
          <Skeleton v-else w="40" h="22" radius="6" center />
        </span>
        <span class="info-label">{{ t('repo.diskSize') }}</span>
      </div>
      <div class="glass info-card clickable dim" :class="statsLoaded ? { expanded: expandedContributor, 'has-data': hasCommits } : {}" @click="toggleContributor">
        <span class="info-icon" style="color:var(--color-purple-soft)">⊕</span>
        <span class="info-value">
          <template v-if="statsLoaded">{{ hasCommits ? detail.contributors.length : '--' }}</template>
          <Skeleton v-else w="30" h="22" radius="6" center />
        </span>
        <span class="info-label">{{ t('repo.contributors') }} <span class="expand-icon">{{ expandedContributor ? '▼' : '▶' }}</span></span>
      </div>
      <div class="glass info-card dim" :class="{ 'has-data': statsLoaded && !!detail.earliestCommitAuthor }">
        <span class="info-icon" style="color:var(--color-amber)">◎</span>
        <span class="info-value">
          <template v-if="statsLoaded">{{ formatDate(detail.earliestDate) }}</template>
          <Skeleton v-else w="55" h="22" radius="6" center />
        </span>
        <span class="info-label">{{ t('repo.createDate') }} · {{ timeSpan }}</span>
      </div>
      <div class="glass info-card dim" :class="{ 'has-data': statsLoaded && !!detail.lastCommitTime }">
        <span class="info-icon" style="color:var(--color-emerald)">◉</span>
        <span class="info-value">
          <template v-if="statsLoaded">{{ formatTimeAgo(detail.lastCommitTime) }}</template>
          <Skeleton v-else w="45" h="22" radius="6" center />
        </span>
        <span class="info-label">{{ t('repo.lastCommit') }}</span>
      </div>
    </div>

    <template v-if="statsLoaded">
      <div v-if="expandedContributor" class="expand-panel glass card">
        <h4>{{ t('repo.contributors') }}</h4>
        <div class="contrib-table">
          <div class="contrib-header">
            <span>{{ t('repo.author') }}</span>
            <span>{{ t('repo.commits') }}</span>
            <span class="add">{{ t('repo.additions') }}</span>
            <span class="del">{{ t('repo.deletions') }}</span>
            <span>{{ t('repo.lastCommit') }}</span>
          </div>
          <div v-for="ct in detail.contributors" :key="ct.email" class="contrib-row">
            <span class="contrib-name">{{ ct.author }}</span>
            <span>{{ ct.commitCount }}</span>
            <span class="add">+{{ ct.additions }}</span>
            <span class="del">-{{ ct.deletions }}</span>
            <span class="contrib-time">{{ ct.lastCommitDate?.slice(0, 10) }}</span>
          </div>
        </div>
      </div>
    </template>

    <div class="charts-wrapper">
      <template v-if="chartLoaded">
        <RepoCharts :data="chartData" :loading="chartLoading" />
      </template>
      <template v-else>
        <div class="charts-section-skel">
          <div class="section-header">
            <h3 class="section-title">{{ t('repo.commitCalendar') }}</h3>
          </div>
          <Skeleton chart h="160" mb="2rem" radius="12" />
          <div class="chart-row-skel">
            <div class="glass card" style="flex:1;padding:1rem;min-height:300px">
              <Skeleton w="40" h="20" radius="6" mb="1rem" />
              <Skeleton chart h="240" />
            </div>
            <div class="glass card" style="flex:1;padding:1rem;min-height:300px">
              <Skeleton w="40" h="20" radius="6" mb="1rem" />
              <Skeleton chart h="240" />
            </div>
          </div>
        </div>
        <div class="chart-overlay" @click="emit('load-chart')">
          <div class="chart-overlay-content">
            <span class="chart-cta-icon">▦</span>
            <h4>{{ t('repo.chartTitle') }}</h4>
            <p>{{ t('repo.chartDesc') }}</p>
            <button class="btn" :disabled="chartLoading" @click.stop="emit('load-chart')">
              <span v-if="chartLoading" class="spinner"></span>
              {{ chartLoading ? t('repo.statsLoading') : t('repo.chartBtn') }}
            </button>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '../i18n'
import Skeleton from './Skeleton.vue'
import RepoCharts from './RepoCharts.vue'

const { t } = useI18n()

const props = defineProps({
  detail: { type: Object, required: true },
  statsLoaded: Boolean,
  chartData: { default: null },
  chartLoading: Boolean,
  chartLoaded: Boolean,
  tags: { type: Array, default: () => [] },
  tagsLoading: Boolean,
  tagsHasMore: Boolean,
  tagsLoaded: Boolean
})

const emit = defineEmits(['load-chart', 'load-tags', 'load-more-tags'])

const expandedBranch = ref(false)
const expandedTags = ref(false)
const expandedContributor = ref(false)

const localBranchCount = computed(() => {
  const local = props.detail?.branchCount ?? 0
  const remote = Array.isArray(props.detail?.remoteBranches) ? props.detail.remoteBranches.length : 0
  if (local === 0 && remote === 0) return '--'
  return remote > 0 ? `${local}+${remote}` : String(local)
})

const remoteBranchCount = computed(() => {
  const v = props.detail?.remoteBranches
  return Array.isArray(v) ? v.length : '--'
})



const localBranches = computed(() => {
  if (props.detail?.branches) return props.detail.branches
  return []
})

const totalCommitCount = computed(() =>
  props.detail?.contributors?.reduce((s, c) => s + c.commitCount, 0) || 0
)

const hasCommits = computed(() =>
  props.detail?.contributors && props.detail.contributors.length > 0
)

const timeSpan = computed(() => {
  if (!props.detail?.earliestDate || !props.detail?.lastCommitTime) return ''
  const d1 = new Date(props.detail.earliestDate)
  const d2 = new Date(props.detail.lastCommitTime)
  const ms = d2 - d1
  if (ms <= 0) return ''
  const days = Math.round(ms / 86400000)
  if (days < 30) return days + t('repo.dayUnit')
  if (days < 365) return Math.round(days / 30) + t('repo.monthUnit')
  const y = Math.floor(days / 365)
  const m = Math.round((days % 365) / 30)
  return y + t('repo.yearUnit') + (m ? m + t('repo.monthUnit') : '')
})

function formatDate(s) {
  if (!s) return '-'
  return s.slice(0, 10)
}

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

function formatBytes(n) {
  if (!n) return '0B'
  const units = ['B', 'KB', 'MB', 'GB']
  let i = 0
  let v = n
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++ }
  return v.toFixed(i > 0 ? 1 : 0) + units[i]
}

function toggleBranch() {
  expandedBranch.value = !expandedBranch.value
}

function toggleTags() {
  expandedTags.value = !expandedTags.value
  if (expandedTags.value && !props.tagsLoaded) {
    emit('load-tags')
  }
}

function toggleContributor() {
  expandedContributor.value = !expandedContributor.value
}
</script>

<style scoped>
.card-row {
  display: grid;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}
.l1-row { grid-template-columns: repeat(4, 1fr); }
.l2-row { grid-template-columns: repeat(4, 1fr); }

.info-card {
  padding: 1.5rem;
  text-align: center;
}
.info-card.dim { opacity: 0.5; }
.info-card.dim.has-data { opacity: 1; }
.info-card.clickable { cursor: pointer; }
.info-card.expanded {
  border-color: var(--border-card-hover);
}

.info-icon { display: block; font-size: 1.8rem; margin-bottom: 0.5rem; }
.info-value {
  display: block;
  font-family: var(--font-display);
  font-size: 1.3rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 0.3rem;
}
.dim .info-value { color: var(--color-range-hint); }
.dim.has-data .info-value { color: var(--color-text-primary); }
.info-label {
  font-family: var(--font-display);
  font-size: 0.65rem;
  color: var(--color-nav-link);
  letter-spacing: 1px;
  text-transform: uppercase;
}
.expand-icon { font-size: 0.5rem; margin-left: 0.2rem; vertical-align: middle; }

.expand-panel {
  margin-bottom: 1.5rem;
  padding: 1.5rem;
}
.expand-panel h4 {
  font-family: var(--font-display);
  font-size: 0.9rem;
  color: var(--color-primary);
  letter-spacing: 1px;
  margin: 0 0 1rem 0;
}
.expand-hint { color: var(--color-text-muted); font-family: var(--font-body); font-size: 0.9rem; }

.branch-groups {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.branch-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.branch-group-label {
  font-family: var(--font-display);
  font-size: 0.75rem;
  color: var(--color-nav-link);
  letter-spacing: 1px;
  text-transform: uppercase;
}
.branch-tags { display: flex; flex-wrap: wrap; gap: 0.6rem; }
.branch-tag {
  padding: 0.3rem 0.8rem;
  border-radius: 12px;
  font-size: 0.85rem;
  font-family: var(--font-body);
  background: var(--bg-cta);
  border: 1px solid var(--border-year-tab);
  color: var(--color-nav-link);
}
.branch-tag.current {
  background: var(--bg-badge-add-intense);
  border-color: var(--border-badge-add-intense);
  color: var(--color-green);
  font-weight: 600;
}
.current-badge {
  font-size: 0.6rem;
  margin-left: 0.3rem;
  padding: 0.05rem 0.35rem;
  border-radius: 4px;
  background: var(--bg-green-mid);
  color: var(--color-green);
  font-family: var(--font-body);
  vertical-align: middle;
  letter-spacing: 0.5px;
}
.branch-tag.remote {
  background: var(--bg-purple-soft-dim);
  border-color: var(--border-purple-soft);
  color: var(--color-purple-soft);
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}
.tag-item {
  padding: 0.25rem 0.7rem;
  border-radius: 6px;
  font-size: 0.8rem;
  font-family: var(--font-body);
  background: var(--bg-emerald-dim);
  border: 1px solid var(--border-emerald);
  color: var(--color-emerald);
}

.load-more-tags {
  width: 100%;
  display: flex;
  justify-content: center;
  margin-top: 0.5rem;
}
.btn-sm {
  padding: 0.3rem 1rem;
  font-size: 0.8rem;
}

.contrib-table { width: 100%; }
.contrib-header, .contrib-row {
  display: grid;
  grid-template-columns: 1.5fr 1fr 1fr 1fr 1.5fr;
  padding: 0.6rem 0.5rem;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.85rem;
}
.contrib-header {
  background: var(--bg-row-header);
  border-bottom: 1px solid var(--border-card-subtle);
  font-family: var(--font-display);
  font-size: 0.7rem;
  color: var(--color-nav-link);
  letter-spacing: 1px;
  text-transform: uppercase;
}
.contrib-row { border-bottom: 1px solid var(--border-row-subtle); color: var(--color-text-primary); }
.contrib-name { font-weight: 600; }
.contrib-time { font-family: var(--font-body); color: var(--color-text-muted); font-size: 0.8rem; }
.add { color: var(--color-green); }
.del { color: var(--color-red); }

.spinner {
  display: inline-block;
  width: 16px; height: 16px;
  border: 2px solid var(--bg-spinner);
  border-top-color: var(--color-white);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  vertical-align: middle;
}

.charts-wrapper {
  position: relative;
  margin-bottom: 1.5rem;
}
.chart-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-overlay);
  backdrop-filter: blur(var(--blur-card));
  border-radius: 12px;
  cursor: pointer;
  z-index: 2;
  transition: all 0.3s;
}
.chart-overlay:hover {
  background: var(--bg-overlay-hover);
  backdrop-filter: blur(var(--blur-card));
}
.chart-overlay-content {
  text-align: center;
  padding: 1.5rem;
}
.chart-overlay-content h4 {
  font-family: var(--font-display);
  font-size: 1rem;
  color: var(--color-primary);
  letter-spacing: 1px;
  margin: 0.5rem 0 0.3rem 0;
}
.chart-overlay-content p {
  color: var(--color-nav-link);
  font-size: 0.85rem;
  margin: 0 0 1rem 0;
  font-family: var(--font-body);
}

.chart-cta-icon {
  display: block;
  font-size: 2.5rem;
  color: var(--color-primary);
  opacity: 0.6;
}

.charts-section-skel {
  margin-top: 2rem;
}
.charts-section-skel .section-header {
  display: flex;
  align-items: center;
  margin-bottom: 1rem;
}
.chart-row-skel {
  display: flex;
  gap: 1.5rem;
}
.chart-row-skel .card {
  border-radius: 16px;
}
.section-title {
  font-family: var(--font-display);
  font-size: 0.95rem;
  color: var(--color-primary);
  letter-spacing: 1px;
  margin: 0;
}
</style>
