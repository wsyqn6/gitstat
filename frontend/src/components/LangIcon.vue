<template>
  <svg v-if="iconData" viewBox="0 0 24 24" width="18" height="18" class="fi" aria-hidden="true">
    <template v-if="iconData.letter">
      <rect x="0" y="0" width="24" height="24" rx="4" :fill="iconData.bg" />
      <text x="12" y="16" text-anchor="middle" :font-size="iconData.letter.length > 2 ? 7 : 9" font-weight="700" :fill="iconData.color" font-family="system-ui,sans-serif">{{ iconData.letter }}</text>
    </template>
    <template v-else-if="iconData.path">
      <path :d="iconData.path" :fill="iconData.color" />
    </template>
  </svg>
  <svg v-else viewBox="0 0 24 24" width="18" height="18" class="fi" aria-hidden="true">
    <circle cx="12" cy="12" r="10" fill="none" stroke="currentColor" opacity=".15" />
    <text x="12" y="16" text-anchor="middle" font-size="8" fill="currentColor" opacity=".35" font-family="system-ui,sans-serif">{{ fallbackLabel }}</text>
  </svg>
</template>

<script setup>
import { computed } from 'vue'
import { getLangIcon } from '../utils/langIcons'

const props = defineProps({
  filePath: { type: String, required: true }
})

const iconData = computed(() => {
  const name = props.filePath.split('/').pop() || ''
  const dot = name.lastIndexOf('.')
  const ext = dot === -1 ? '' : name.slice(dot + 1).toLowerCase()
  return getLangIcon(ext)
})

const fallbackLabel = computed(() => {
  const name = props.filePath.split('/').pop() || ''
  const dot = name.lastIndexOf('.')
  const ext = dot === -1 ? '' : name.slice(dot + 1).toLowerCase()
  return ext.slice(0, 2) || '?'
})
</script>

<style scoped>
.fi {
  display: block;
  border-radius: 3px;
  overflow: hidden;
  flex-shrink: 0;
}
</style>
