<template>
  <div class="section card">
    <h3 class="section-title">{{ t('repo.langDistribution') }}</h3>
    <div v-if="analysisLoading" class="section-loading">
      <span class="spinner"></span>
      <span>{{ t('repo.analyzing') }}</span>
    </div>
    <div v-else-if="analysis" class="lang-body">
      <div ref="langChartRef" id="lang-chart-el" class="lang-chart"></div>
      <div class="lang-table">
        <div class="lang-header">
          <span>{{ t('repo.lang') }}</span>
          <span>{{ t('repo.files') }}</span>
          <span>{{ t('repo.lines') }}</span>
          <span>{{ t('repo.percent') }}</span>
        </div>
        <div v-for="lang in analysis.languages" :key="lang.name" class="lang-row">
          <span class="lang-name">{{ lang.name }}</span>
          <span>{{ lang.fileCount }}</span>
          <span>{{ formatNumber(lang.lines) }}</span>
          <span class="lang-pct">
            <span class="pct-track"><span class="pct-bar" :style="{ width: lang.percentage + '%' }"></span></span>
            <span class="pct-text">{{ lang.percentage.toFixed(1) }}%</span>
          </span>
        </div>
        <div class="lang-total">
          <span>{{ t('repo.total') }}</span>
          <span>{{ analysis.fileCount }}</span>
          <span>{{ formatNumber(analysis.totalLines) }}</span>
          <span>100%</span>
        </div>
      </div>
    </div>
    <div v-else class="analyze-cta">
      <span class="analyze-icon">🔍</span>
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
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useI18n } from '../i18n'
import echarts from '../utils/echarts'

const { t } = useI18n()

const props = defineProps({
  analysis: { type: Object, default: null },
  analysisLoading: Boolean
})

const emit = defineEmits(['analyze'])

const langChartRef = ref(null)
let langChart = null

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

function formatNumber(n) {
  if (n >= 1000) return (n / 1000).toFixed(1) + 'k'
  return String(n)
}

function renderChart() {
  const el = document.getElementById('lang-chart-el')
  if (!el || !props.analysis) return
  try {
    if (!langChart) langChart = echarts.init(el)
    const data = props.analysis.languages.map(l => ({
      name: l.name,
      value: l.lines || 1
    }))
    langChart.setOption({
      tooltip: {
        trigger: 'item',
        backgroundColor: 'rgba(10, 14, 39, 0.9)',
        borderColor: 'rgba(0, 212, 255, 0.3)',
        textStyle: { color: '#e0e6ff', fontSize: 12 },
        formatter: p => `${p.name}: ${p.percent}% (${formatNumber(p.value)} lines)`
      },
      series: [{
        type: 'pie',
        radius: ['30%', '65%'],
        center: ['50%', '50%'],
        avoidLabelOverlap: true,
        itemStyle: { borderRadius: 4, borderColor: 'rgba(10,14,39,0.8)', borderWidth: 2 },
        label: {
          show: true, color: '#e0e6ff', fontFamily: 'Rajdhani', fontSize: 12,
          formatter: p => `${p.name}\n${p.percent}%`
        },
        labelLine: { lineStyle: { color: 'rgba(0,212,255,0.3)' } },
        data: data.map(d => ({ ...d, itemStyle: { color: getLangColor(d.name) } }))
      }]
    })
    langChart.resize()
  } catch (e) {
    console.error('[Chart] render error:', e)
  }
}

function handleResize() {
  if (langChart) langChart.resize()
}

watch(() => props.analysis, val => {
  if (val) nextTick(renderChart)
})

watch(langChartRef, el => {
  if (el && props.analysis) nextTick(renderChart)
})

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  langChart?.dispose()
  langChart = null
})
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
.section-loading { text-align: center; padding: 2rem; color: #00d4ff; font-family: var(--font-body); }

.spinner {
  display: inline-block;
  width: 16px; height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  vertical-align: middle;
  margin-right: 0.5rem;
}

.lang-body { display: grid; grid-template-columns: 1fr 1.2fr; gap: 1.5rem; align-items: start; }
.lang-chart { height: 250px; min-height: 250px; width: 100%; }
.lang-header, .lang-row, .lang-total {
  display: grid;
  grid-template-columns: 1.5fr 0.8fr 1fr 1.5fr;
  padding: 0.5rem 0.5rem;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.85rem;
}
.lang-header {
  background: rgba(0, 212, 255, 0.05);
  border-bottom: 1px solid rgba(0, 212, 255, 0.15);
  font-family: var(--font-display);
  font-size: 0.7rem;
  color: #a0aec0;
  letter-spacing: 1px;
  text-transform: uppercase;
}
.lang-row { border-bottom: 1px solid rgba(0, 212, 255, 0.06); color: #e0e6ff; }
.lang-row:hover { background: rgba(0, 212, 255, 0.03); }
.lang-name { font-weight: 600; color: #e0e6ff; }
.lang-pct { display: flex; align-items: center; gap: 0.3rem; font-family: var(--font-display); font-size: 0.75rem; color: #00d4ff; }
.pct-track { flex: 1; height: 8px; background: rgba(0,0,0,0.15); border-radius: 4px; overflow: hidden; min-width: 4px; }
.pct-bar { display: block; height: 8px; border-radius: 4px; background: linear-gradient(90deg, #00d4ff, #7800ff); }
.pct-text { flex-shrink: 0; font-family: 'Rajdhani', 'Orbitron', monospace; }
.lang-total {
  border-top: 1px solid rgba(0, 212, 255, 0.3);
  font-family: var(--font-display);
  font-size: 0.8rem;
  color: #00d4ff;
  font-weight: 700;
  padding-top: 0.6rem;
}
.lang-total span { text-align: center; }
.lang-total span:first-child { text-align: left; }

.analyze-cta { text-align: center; padding: 2rem 1rem; }
.analyze-icon { font-size: 2.5rem; display: block; margin-bottom: 0.8rem; }
.analyze-cta h4 { font-family: var(--font-display); font-size: 1.1rem; color: #00d4ff; letter-spacing: 1px; margin-bottom: 0.4rem; }
.analyze-desc { color: #64748b; font-size: 0.9rem; margin-bottom: 1rem; font-family: var(--font-body); }
</style>
