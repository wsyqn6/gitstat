<template>
  <div v-if="loading" class="stats-grid">
    <div v-for="i in 6" :key="i" class="stat-card-ph">
      <Skeleton circle class="stat-ph-icon" />
      <Skeleton w="40" class="stat-ph-value" />
      <Skeleton w="60" class="stat-ph-label" />
    </div>
  </div>
  <div v-else class="stats-grid">
    <StatCard :value="todayCommits" :label="t('dashboard.todayCommits')" icon="◈" color="#00d4ff" />
    <StatCard :value="todayAdditions" :label="t('dashboard.todayAdditions')" icon="↑" color="#00ff88" />
    <StatCard :value="todayDeletions" :label="t('dashboard.todayDeletions')" icon="↓" color="#ff6b9d" />
    <StatCard :value="activeAuthors" :label="t('dashboard.activeAuthors')" icon="◉" color="#ffd700" />
    <StatCard :value="repositoryCount" :label="t('dashboard.repositoryCount')" icon="▦" color="#a78bfa" />
    <StatCard :value="weeklyTotal" :label="t('dashboard.weeklyTotal')" icon="⟡" color="#f472b6" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from '../i18n'
import Skeleton from './Skeleton.vue'
import StatCard from './StatCard.vue'

const { t } = useI18n()

const props = defineProps({
  loading: { type: Boolean, required: true },
  overviewStats: { type: Object, default: null },
  weeklyTotal: { type: Number, default: 0 }
})

const todayCommits = computed(() => props.overviewStats?.totalCommits ?? 0)
const todayAdditions = computed(() => props.overviewStats?.totalAdditions ?? 0)
const todayDeletions = computed(() => props.overviewStats?.totalDeletions ?? 0)
const activeAuthors = computed(() => props.overviewStats?.activeAuthors ?? 0)
const repositoryCount = computed(() => props.overviewStats?.repositoryCount ?? 0)
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2rem;
  margin-bottom: 3rem;
}

.stat-card-ph {
  background: var(--bg-card);
  backdrop-filter: blur(var(--blur-card));
  padding: 2rem;
  border-radius: 16px;
  text-align: center;
  border: 1px solid var(--border-insight-card);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.stat-ph-icon {
  width: 40px;
  height: 40px;
  margin-bottom: 0.5rem;
}

.stat-ph-value {
  height: 40px;
}

.stat-ph-label {
  height: 14px;
}
</style>
