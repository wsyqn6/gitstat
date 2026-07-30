<template>
  <div class="section glass card lang-card">
    <h3 class="section-title">{{ t('repo.langDistribution') }}</h3>
    <div v-if="analysisLoading" class="lang-skeleton">
      <Skeleton w="30" h="14" mb="0.75rem" />
      <Skeleton w="100" h="10" radius="5" mb="0.75rem" />
      <div class="lang-list">
        <div v-for="i in 5" :key="i" class="lang-row">
          <Skeleton circle h="10" />
          <Skeleton w="35" />
          <Skeleton w="15" />
          <Skeleton w="20" />
        </div>
      </div>
    </div>
    <template v-else-if="analysis">
      <div class="lang-bar">
        <span
          v-for="lang in analysis.languages"
          :key="lang.name"
          class="lang-bar-seg"
          :style="{ width: lang.percentage + '%', background: getLangColor(lang.name) }"
          :title="lang.name + ': ' + lang.percentage.toFixed(1) + '%'"
        ></span>
      </div>
      <div class="lang-list">
        <div v-for="lang in analysis.languages" :key="lang.name" class="lang-row">
          <span class="lang-dot" :style="{ background: getLangColor(lang.name) }"></span>
          <span class="lang-name">{{ lang.name }}</span>
          <span class="lang-pct">{{ lang.percentage.toFixed(1) }}%</span>
          <span class="lang-meta">{{ formatNumber(lang.lines) }} lines</span>
        </div>
      </div>
    </template>
    <div v-else class="analyze-cta">
      <span class="analyze-icon" aria-hidden="true">🔍</span>
      <h4>{{ t('repo.analyzeTitle') }}</h4>
      <p class="analyze-desc">{{ t('repo.analyzeDesc') }}</p>
      <button class="btn" :disabled="analysisLoading" @click="emit('analyze')">
        <span v-if="analysisLoading" class="spinner"></span>
        {{ analysisLoading ? t('repo.analyzing') : t('repo.analyzeBtn') }}
      </button>
    </div>
  </div>
</template>

<script setup>
import Skeleton from './Skeleton.vue'
import { useI18n } from '../i18n'
import { formatNumber } from '../utils/format'

const { t } = useI18n()

defineProps({
  analysis: { type: Object, default: null },
  analysisLoading: Boolean
})

const emit = defineEmits(['analyze'])

const LANG_COLORS = {
  Rust: '#dea584', Go: '#00add8', Python: '#3572a5',
  JavaScript: '#f1e05a', TypeScript: '#3178c6',
  'React (JSX)': '#61dafb', 'React (TSX)': '#61dafb',
  Vue: '#41b883', Java: '#b07219', Kotlin: '#a97bff',
  C: '#555555', 'C++': '#f34b7d', 'C#': '#178600',
  Ruby: '#701516', PHP: '#4f5d95', Swift: '#ffac45',
  Dart: '#00b4ab', CSS: '#563d7c', SCSS: '#c6538c',
  HTML: '#e34c26', Shell: '#89e051', Dockerfile: '#384d54',
  SQL: '#e38c00', Markdown: '#083fa1', YAML: '#cb171e',
  TOML: '#9c4221', JSON: '#292929', Lua: '#000080',
  Scala: '#c22d40', Zig: '#ec915c', Default: '#a0aec0'
}

function getLangColor(name) {
  return LANG_COLORS[name] || LANG_COLORS.Default
}

</script>

<style scoped>
.lang-card { display: flex; flex-direction: column; }
.section-title {
  font-family: var(--font-display);
  font-size: 0.95rem;
  color: var(--color-primary);
  letter-spacing: 1px;
  margin: 0 0 0.75rem 0;
  flex-shrink: 0;
}
.lang-skeleton { padding: 0.25rem 0; }
.section-loading { text-align: center; padding: 2rem; color: var(--color-primary); font-family: var(--font-body); }



.lang-bar {
  display: flex;
  height: 10px;
  border-radius: 5px;
  overflow: hidden;
  flex-shrink: 0;
  margin-bottom: 0.75rem;
}
.lang-bar-seg {
  height: 100%;
  transition: filter 0.15s;
}
.lang-bar-seg:hover { filter: brightness(1.3); }

.lang-list {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  padding-right: 0.3rem;
}
.lang-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.85rem;
  color: var(--color-text-primary);
  padding: 0.2rem 0;
  flex-shrink: 0;
}
.lang-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}
.lang-name { flex: 1; font-weight: 600; font-size: 0.85rem; }
.lang-pct {
  font-family: var(--font-display);
  font-size: 0.85rem;
  color: var(--color-primary);
  width: 4.5rem;
  text-align: right;
  flex-shrink: 0;
}
.lang-meta {
  font-family: var(--font-body);
  font-size: 0.8rem;
  color: var(--color-text-muted);
  width: 5rem;
  text-align: right;
  flex-shrink: 0;
}

.analyze-cta { text-align: center; padding: 1.5rem 1rem; margin: auto 0; }
.analyze-icon { font-size: 2rem; display: block; margin-bottom: 0.5rem; }
.analyze-cta h4 { font-family: var(--font-display); font-size: 1rem; color: var(--color-primary); letter-spacing: 1px; margin-bottom: 0.3rem; }
.analyze-desc { color: var(--color-text-muted); font-size: 0.85rem; margin-bottom: 0.8rem; font-family: var(--font-body); }
</style>
