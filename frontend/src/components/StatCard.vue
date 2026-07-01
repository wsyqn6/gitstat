<template>
  <div class="glass stat-card" :style="cardStyle">
    <div class="stat-icon">{{ icon }}</div>
    <div class="stat-value">{{ value }}</div>
    <ChangeBadge v-if="change" :change="change" prefix="较昨日" />
    <div class="stat-label">{{ label }}</div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import ChangeBadge from './ChangeBadge.vue'

const props = defineProps({
  value: { type: [Number, String], required: true },
  label: { type: String, required: true },
  icon: { type: String, default: '◈' },
  color: { type: String, default: '#00d4ff' },
  change: { type: Object, default: null }
})

const cardStyle = computed(() => ({
  '--accent-color': props.color
}))
</script>

<style scoped>
.stat-card {
  padding: 1.5rem;
  text-align: center;
}

.stat-icon {
  font-size: 2.5rem;
  margin-bottom: 1rem;
  color: var(--accent-color);
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}

.stat-value {
  font-family: var(--font-display);
  font-size: 2.5rem;
  font-weight: 900;
  background: linear-gradient(135deg, var(--accent-color), white);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 0.5rem;
  position: relative;
  z-index: 1;
}

.stat-label {
  font-family: var(--font-display);
  color: var(--color-nav-link);
  font-size: 0.85rem;
  letter-spacing: 2px;
  text-transform: uppercase;
  position: relative;
  z-index: 1;
}
</style>
