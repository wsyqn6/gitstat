<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-panel">
      <div class="modal-header">
        <span class="modal-title">
          <svg viewBox="0 0 16 16" fill="currentColor" class="modal-title-icon">
            <path d="M9.796 1.343c-.527-1.79-3.065-1.79-3.592 0l-.094.319a.873.873 0 0 1-1.255.52l-.292-.16c-1.64-.892-3.433.902-2.54 2.541l.159.292a.873.873 0 0 1-.52 1.255l-.319.094c-1.79.527-1.79 3.065 0 3.592l.319.094a.873.873 0 0 1 .52 1.255l-.16.292c-.892 1.64.901 3.434 2.541 2.54l.292-.159a.873.873 0 0 1 1.255.52l.094.319c.527 1.79 3.065 1.79 3.592 0l.094-.319a.873.873 0 0 1 1.255-.52l.292.16c1.64.893 3.434-.902 2.54-2.541l-.159-.292a.873.873 0 0 1 .52-1.255l.319-.094c1.79-.527 1.79-3.065 0-3.592l-.319-.094a.873.873 0 0 1-.52-1.255l.16-.292c.893-1.64-.902-3.433-2.541-2.54l-.292.159a.873.873 0 0 1-1.255-.52zm-2.633.283c.246-.835 1.428-.835 1.674 0l.094.319a1.873 1.873 0 0 0 2.693 1.115l.291-.16c.764-.415 1.6.42 1.184 1.185l-.159.292a1.873 1.873 0 0 0 1.116 2.692l.318.094c.835.246.835 1.428 0 1.674l-.319.094a1.873 1.873 0 0 0-1.115 2.693l.16.291c.415.764-.42 1.6-1.185 1.184l-.291-.159a1.873 1.873 0 0 0-2.693 1.116l-.094.318c-.246.835-1.428.835-1.674 0l-.094-.319a1.873 1.873 0 0 0-2.692-1.115l-.292.16c-.764.415-1.6-.42-1.184-1.185l.159-.291A1.873 1.873 0 0 0 1.945 8.93l-.319-.094c-.835-.246-.835-1.428 0-1.674l.319-.094A1.873 1.873 0 0 0 3.06 4.377l-.16-.292c-.415-.764.42-1.6 1.185-1.184l.292.159a1.873 1.873 0 0 0 2.692-1.115z"/>
            <path d="M8 4.754a3.246 3.246 0 1 0 0 6.492 3.246 3.246 0 0 0 0-6.492M5.754 8a2.246 2.246 0 1 1 4.492 0 2.246 2.246 0 0 1-4.492 0"/>
          </svg>
          {{ t('settings.title') }}
        </span>
        <button class="modal-close-btn" aria-label="Close" @click="$emit('close')">✕</button>
      </div>

      <div class="modal-body">
        <div class="section toggles-section">
          <div class="toggle-row">
            <span class="toggle-label">{{ t('settings.theme') }}</span>
            <div class="segment-group" :class="{ 'slide-right': theme === 'light' }">
              <div class="segment-slider"></div>
              <button :class="{ active: theme === 'neon' }" @click="setTheme('neon')">{{ t('settings.themeDark') }}</button>
              <button :class="{ active: theme === 'light' }" @click="setTheme('light')">{{ t('settings.themeLight') }}</button>
            </div>
          </div>
          <div class="toggle-row">
            <span class="toggle-label">{{ t('settings.language') }}</span>
            <div class="segment-group" :class="{ 'slide-right': locale === 'en' }">
              <div class="segment-slider"></div>
              <button :class="{ active: locale === 'zh' }" @click="setLocale('zh')">中文</button>
              <button :class="{ active: locale === 'en' }" @click="setLocale('en')">EN</button>
            </div>
          </div>
        </div>

        <div class="section">
          <h3 class="section-title">{{ t('settings.scanConfig') }}</h3>
          <div class="form-group">
            <label class="field-label">{{ t('settings.directoryPath') }}</label>
            <div class="input-row">
              <input v-model="scanPath" :placeholder="t('settings.pathPlaceholder')" />
              <button @click="handleScan" :disabled="scanning" class="btn scan-btn" style="flex-shrink:0">
                <span v-if="!scanning">{{ t('settings.startScan') }}</span>
                <span v-else><span class="spinner"></span> {{ t('settings.scanning') }}</span>
              </button>
            </div>
          </div>
          <p v-if="scanError" class="msg error-msg">{{ scanError }}</p>
          <p v-if="scanSuccess" class="msg success-msg">{{ scanSuccess }}</p>
        </div>

        <div class="section">
          <h3 class="section-title">{{ t('settings.dataManagement') }}</h3>
          <button @click="handleExport" class="btn export-btn" style="width:100%">
            <span style="font-size:1.2rem">↓</span> {{ t('settings.exportData') }}
          </button>
        </div>

        <div class="section about-section">
          <h3 class="section-title">{{ t('settings.about') }}</h3>
          <div class="about-content">
            <div class="about-logo">GITSTAT</div>
            <p class="about-version">{{ version }} — {{ t('settings.platformName') }}</p>
            <div class="tech-stack">
              <span class="tech-item">⚡ Go 1.26</span>
              <span class="tech-item">◈ Vue 3</span>
              <span class="tech-item">◉ ECharts</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from '../i18n'
import { useTheme } from '../composables/useTheme'
import { performScan, state, loadScanPath } from '../stores/data'
import { getScanPath, exportData } from '../api'
import { useToast } from '../composables/useToast'

const emit = defineEmits(['close'])

const { t, locale, setLocale } = useI18n()
const { theme, setTheme } = useTheme()
const { show } = useToast()

const scanPath = ref(state.scanPath || '')
const scanning = ref(false)
const scanError = ref('')
const scanSuccess = ref('')
const version = ref('dev')

onMounted(async () => {
  try {
    const info = await getScanPath()
    version.value = info.version || 'dev'
  } catch {}
  await loadScanPath()
  scanPath.value = state.scanPath || ''
})

async function handleScan() {
  scanning.value = true
  scanError.value = ''
  scanSuccess.value = ''
  try {
    state.scanPath = scanPath.value
    await performScan(scanPath.value)
    scanSuccess.value = t('settings.scanSuccess')
    setTimeout(() => { scanSuccess.value = '' }, 3000)
  } catch (err) {
    scanError.value = err.message
  } finally {
    scanning.value = false
  }
}

async function handleExport() {
  try {
    const blob = await exportData()
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'gitstat-data.json'
    a.click()
    window.URL.revokeObjectURL(url)
  } catch {
    show(t('settings.exportError'), 'error')
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  animation: fadeIn 0.2s ease;
}

.modal-panel {
  width: 520px;
  max-width: 92vw;
  max-height: 85vh;
  background: var(--bg-card);
  backdrop-filter: blur(var(--blur-card));
  border: 1px solid var(--border-card);
  border-radius: 20px;
  box-shadow: var(--shadow-dropdown), 0 0 60px rgba(0,0,0,0.2);
  display: flex;
  flex-direction: column;
  animation: slideUp 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--border-card);
  flex-shrink: 0;
}

.modal-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 1.1rem;
  color: var(--color-primary);
  letter-spacing: 1px;
}

.modal-title-icon {
  width: 1.1rem;
  height: 1.1rem;
  flex-shrink: 0;
}

.modal-close-btn {
  background: transparent;
  border: none;
  color: var(--color-text-muted);
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0.25rem;
  line-height: 1;
  transition: color 0.2s;
}
.modal-close-btn:hover { color: var(--color-text-primary); }

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.section-title {
  font-family: var(--font-display);
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--color-primary);
  letter-spacing: 1px;
  text-transform: uppercase;
  margin-bottom: 1rem;
}

.toggles-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.toggle-label {
  font-family: var(--font-body);
  font-size: 0.9rem;
  color: var(--color-text-primary);
  font-weight: 500;
}

.segment-group {
  display: flex;
  background: var(--bg-tab);
  border: 1px solid var(--border-card);
  border-radius: var(--radius-btn);
  padding: 2px;
  gap: 2px;
  position: relative;
}

.segment-slider {
  position: absolute;
  top: 2px;
  left: 2px;
  width: calc(50% - 3px);
  height: calc(100% - 4px);
  border-radius: calc(var(--radius-btn) - 2px);
  background: var(--glass-btn-bg);
  backdrop-filter: blur(var(--glass-blur));
  box-shadow: var(--glass-btn-shadow), var(--glass-btn-inner);
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  z-index: 0;
  pointer-events: none;
}

.segment-group.slide-right .segment-slider {
  transform: translateX(calc(100% + 2px));
}

.segment-group button {
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
}

.segment-group button:active {
  transform: scale(0.95);
}

.segment-group button.active {
  color: var(--glass-btn-color);
}

.segment-group button:hover:not(.active) {
  color: var(--color-text-secondary);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.field-label {
  font-family: var(--font-display);
  font-size: 0.8rem;
  color: var(--color-nav-link);
  letter-spacing: 1px;
  text-transform: uppercase;
}

.input-row {
  display: flex;
  gap: 0.75rem;
}

.input-row input {
  flex: 1;
  padding: 0.6rem 1rem;
  border: 1px solid var(--border-input);
  border-radius: var(--radius-input);
  background: var(--bg-input);
  color: var(--color-input-text);
  font-family: var(--font-body);
  font-size: 0.9rem;
  transition: all 0.3s;
}

.input-row input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 20px var(--border-input);
}

.scan-btn { min-width: 130px; }

.msg {
  margin-top: 0.5rem;
  padding: 0.6rem 0.75rem;
  border-radius: var(--radius-sm);
  font-size: 0.85rem;
}

.error-msg {
  color: var(--color-danger);
  background: var(--bg-badge-danger);
  border: 1px solid var(--border-badge-danger);
}

.success-msg {
  color: var(--color-green);
  background: var(--bg-badge-success);
  border: 1px solid var(--border-badge-success);
}

.export-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.about-section { text-align: center; }

.about-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.about-logo {
  font-family: var(--font-display);
  font-size: 2rem;
  font-weight: 900;
  background: var(--gradient-logo);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 4px;
}

.about-version {
  color: var(--color-nav-link);
  font-size: 0.9rem;
}

.tech-stack {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  justify-content: center;
}

.tech-item {
  padding: 0.5rem 1.2rem;
  background: var(--bg-btn-hover);
  border: 1px solid var(--border-input);
  border-radius: var(--radius-btn);
  font-size: 0.85rem;
  font-family: var(--font-body);
  color: var(--color-text-secondary);
  transition: all 0.3s;
}

.tech-item:hover {
  background: var(--border-insight-card);
  border-color: var(--color-primary);
  transform: translateY(-1px);
}

.spinner {
  display: inline-block;
  width: 12px;
  height: 12px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  vertical-align: middle;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px) scale(0.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
