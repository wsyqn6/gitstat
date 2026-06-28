<template>
  <div class="segment-group">
    <div class="segment-slider" :style="sliderStyle"></div>
    <slot :active="modelValue" :activate="activate">
      <button
        v-for="opt in options"
        :key="opt.value"
        :class="{ active: modelValue === opt.value }"
        @click="activate(opt.value)"
      >
        {{ opt.label }}
      </button>
    </slot>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  options: { type: Array, required: true },
  modelValue: { required: true }
})

const emit = defineEmits(['update:modelValue'])
const activate = (val) => emit('update:modelValue', val)

const activeIndex = computed(() => {
  const i = props.options.findIndex(o => o.value === props.modelValue)
  return i >= 0 ? i : 0
})

const sliderStyle = computed(() => ({
  transform: `translateX(${activeIndex.value * 100}%)`,
  width: `${100 / props.options.length}%`
}))
</script>

<style scoped>
.segment-group {
  display: flex;
  background: var(--glass-btn-bg);
  backdrop-filter: blur(var(--glass-blur));
  border: 1px solid var(--glass-btn-border);
  border-radius: var(--radius-btn);
  padding: 2px;
  gap: 2px;
  position: relative;
}

.segment-slider {
  position: absolute;
  top: 2px;
  left: 2px;
  height: calc(100% - 4px);
  border-radius: calc(var(--radius-btn) - 2px);
  background: var(--glass-btn-hover-bg);
  backdrop-filter: blur(var(--glass-blur));
  box-shadow: var(--glass-btn-shadow), var(--glass-btn-inner);
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  z-index: 0;
  pointer-events: none;
}

.segment-group :deep(button) {
  flex: 1;
  padding: 0.4rem 1rem;
  border: none;
  border-radius: calc(var(--radius-btn) - 2px);
  background: transparent;
  color: var(--color-text-muted);
  font-family: var(--font-body);
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: color 0.25s ease;
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.segment-group :deep(button.active) {
  color: var(--glass-btn-color);
}

.segment-group :deep(button:hover:not(.active)) {
  color: var(--color-text-secondary);
}
</style>
