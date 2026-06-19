<template>
  <div class="stat-card" :style="cardStyle">
    <div class="stat-icon">{{ icon }}</div>
    <div class="stat-value">{{ value }}</div>
    <div class="stat-label">{{ label }}</div>
    <div class="stat-glow"></div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  value: { type: [Number, String], required: true },
  label: { type: String, required: true },
  icon: { type: String, default: '◈' },
  color: { type: String, default: '#00d4ff' }
})

const cardStyle = computed(() => ({
  '--accent-color': props.color
}))
</script>

<style scoped>
.stat-card {
  background: var(--bg-card);
  backdrop-filter: blur(var(--blur-card));
  padding: 1.5rem;
  border-radius: var(--radius-card);
  text-align: center;
  box-shadow: 
    var(--shadow-card),
    var(--shadow-inset);
  border: 1px solid var(--border-card);
  position: relative;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.stat-card::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(circle, var(--accent-color) 0%, transparent 70%);
  opacity: 0;
  transition: opacity 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  border-color: var(--accent-color);
  box-shadow: 
    var(--shadow-card-hover-dark),
    0 0 30px var(--accent-color);
}

.stat-card:hover::before {
  opacity: 0.1;
}

.stat-icon {
  font-size: 2.5rem;
  margin-bottom: 1rem;
  color: var(--accent-color);
  filter: drop-shadow(0 0 10px var(--accent-color));
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

.stat-glow {
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 60%;
  height: 2px;
  background: linear-gradient(90deg, transparent, var(--accent-color), transparent);
  box-shadow: 0 0 10px var(--accent-color);
}
</style>
