<template>
  <div
    :class="[baseClass, widthClass]"
    :style="computedStyle"
  />
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  w: { type: [String, Number], default: '100' },
  h: { type: [String, Number], default: null },
  circle: Boolean,
  chart: Boolean,
  radius: { type: [String, Number], default: null },
  mb: { type: String, default: '' },
  center: Boolean,
})

const baseClass = computed(() => {
  if (props.circle) return 'skeleton-circle'
  if (props.chart) return 'skeleton-chart-area'
  return 'skeleton-line'
})

const widthClass = computed(() => {
  if (props.circle) return ''
  return `w${props.w}`
})

const computedStyle = computed(() => {
  const style = {}
  if (props.h !== null) {
    style.height = `${props.h}px`
    if (props.circle) style.width = `${props.h}px`
  }
  if (props.radius !== null) {
    style.borderRadius = `${props.radius}px`
  }
  if (props.center && props.mb) {
    style.margin = `0 auto ${props.mb}`
  } else if (props.center) {
    style.margin = '0 auto'
  } else if (props.mb) {
    style.marginBottom = props.mb
  }
  return style
})
</script>
