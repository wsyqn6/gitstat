<template>
  <div v-if="change" class="change-badge" :class="cls">
    <span class="change-prefix">{{ prefix }}</span>
    <span class="change-arrow">{{ arrow }}</span>
    <span class="change-value">{{ formatted }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  change: Object,
  prefix: String,
  compact: Boolean
})

const cls = computed(() => {
  if (!props.change) return 'flat'
  const pct = props.change.pct
  if (pct > 0) return 'up'
  if (pct < 0) return 'down'
  return 'flat'
})

const arrow = computed(() => {
  const pct = props.change?.pct
  if (pct > 0) return '▲'
  if (pct < 0) return '▼'
  return '–'
})

const formatted = computed(() => {
  if (!props.change) return ''
  const pct = props.change.pct
  const abs = props.change.abs
  const pctStr = pct > 0 ? `+${pct}` : `${pct}`
  const absStr = abs > 0 ? `+${abs}` : `${abs}`
  return `${absStr} (${pctStr}%)`
})
</script>

<style scoped>
.change-badge {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  font-family: var(--font-body);
  font-size: 0.65rem;
  font-weight: 600;
  letter-spacing: 0.2px;
  padding: 0 6px;
  border-radius: 8px;
  animation: badgeIn 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

.change-badge.compact {
  padding: 0;
  border-radius: 0;
  background: none !important;
  margin: 0;
  animation: none;
}

.change-badge.up {
  color: var(--color-red);
  background: rgba(var(--color-red-rgb), 0.1);
}

.change-badge.down {
  color: var(--color-green);
  background: rgba(var(--color-green-rgb), 0.1);
}

.change-badge.flat {
  color: var(--color-text-muted);
}

.change-prefix {
  opacity: 0.55;
}

.change-arrow {
  font-size: 0.5rem;
  line-height: 1;
}

@keyframes badgeIn {
  from {
    opacity: 0;
    transform: translateY(4px) scale(0.94);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
</style>
