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
        <button class="modal-close-btn" :aria-label="t('settings.close')" @click="$emit('close')">✕</button>
      </div>

      <div class="modal-body">
        <div class="section toggles-section">
          <div class="toggle-row">
            <span class="toggle-label">{{ t('settings.theme') }}</span>
            <SegmentToggle
              :options="[{ value: 'neon', label: t('settings.themeDark') }, { value: 'light', label: t('settings.themeLight') }]"
              :modelValue="theme"
              @update:modelValue="setTheme"
            />
          </div>
          <div class="toggle-row">
            <span class="toggle-label">{{ t('settings.language') }}</span>
            <SegmentToggle
              :options="[{ value: 'zh', label: '中文' }, { value: 'en', label: 'EN' }]"
              :modelValue="locale"
              @update:modelValue="setLocale"
            />
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
                <span v-else><span class="spinner spinner-sm"></span> {{ t('settings.scanning') }}</span>
              </button>
            </div>
          </div>
          <p v-if="scanError" class="msg error-msg">{{ scanError }}</p>
          <p v-if="scanSuccess" class="msg success-msg">{{ scanSuccess }}</p>
        </div>

        <div class="section about-section">
          <h3 class="section-title">{{ t('settings.about') }}</h3>
          <div class="about-content">
            <div class="about-logo">GITSTAT</div>
            <div class="about-versions">
              <div class="app-version-row">
                <span class="app-version">{{ state.buildVersion || 'dev' }}</span>
                <span v-if="!updateUpToDate && !updateAvailable && !updateChecking && !updateErrorMsg" class="update-action update-check-link" @click="handleUpdateCheck">{{ t('settings.updateCheck') }}</span>
                <span v-if="updateChecking" class="update-action"><span class="spinner spinner-xs"></span> {{ t('settings.updateChecking') }}</span>
                <span v-if="updateUpToDate" class="update-action update-latest">✓</span>
                <template v-if="updateAvailable">
                  <span class="update-action update-new">→ <strong>{{ updateLatest }}</strong></span>
                  <a :href="updateUrl" target="_blank" rel="noopener noreferrer" class="update-download-link-sm">{{ t('settings.updateDownload') }}</a>
                </template>
                <span v-if="updateErrorMsg" class="update-action update-check-link" @click="handleUpdateCheck">{{ t('settings.updateCheck') }}</span>
              </div>
              <div class="git-version-row">{{ t('settings.gitVersion') }} {{ version }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from '../i18n'
import { useTheme } from '../composables/useTheme'
import { performScan, state } from '../stores/data'
import { checkUpdate } from '../api'
import { useToast } from '../composables/useToast'
import SegmentToggle from './SegmentToggle.vue'

const emit = defineEmits(['close'])

const { t, locale, setLocale } = useI18n()
const { theme, setTheme } = useTheme()
const { show } = useToast()

const scanPath = ref(state.scanPath || '')
const scanning = ref(false)
const scanError = ref('')
const scanSuccess = ref('')
const version = ref('dev')
const updateChecking = ref(false)
const updateLatest = ref('')
const updateUrl = ref('')
const updateUpToDate = ref(false)
const updateAvailable = ref(false)
const updateErrorMsg = ref('')
let successTimer = null

onMounted(async () => {
  scanPath.value = state.scanPath || ''
  version.value = state.gitVersion || 'dev'
})

async function handleScan() {
  scanning.value = true
  scanError.value = ''
  scanSuccess.value = ''
  try {
    state.scanPath = scanPath.value
    await performScan(scanPath.value)
    scanSuccess.value = t('settings.scanSuccess')
    successTimer = setTimeout(() => { scanSuccess.value = '' }, 3000)
  } catch (err) {
    scanError.value = err.message
  } finally {
    scanning.value = false
  }
}

async function handleUpdateCheck() {
  updateChecking.value = true
  updateUpToDate.value = false
  updateAvailable.value = false
  updateErrorMsg.value = ''
  try {
    const res = await checkUpdate()
    if (res.hasUpdate) {
      updateLatest.value = res.latestVersion
      updateUrl.value = res.downloadUrl
      updateAvailable.value = true
    } else {
      updateUpToDate.value = true
    }
  } catch (err) {
    updateErrorMsg.value = t('settings.updateError')
  } finally {
    updateChecking.value = false
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

onUnmounted(() => {
  if (successTimer) clearTimeout(successTimer)
})
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--modal-overlay);
  backdrop-filter: blur(4px);
  animation: fadeIn 0.2s ease;
}

.modal-panel {
  width: 520px;
  max-width: 92vw;
  max-height: 85vh;
  background: var(--modal-bg);
  backdrop-filter: blur(var(--blur-card));
  border: 1px solid var(--border-card);
  border-radius: 20px;
  box-shadow: var(--shadow-dropdown), 0 0 60px rgba(0,0,0,0.2);
  display: flex;
  flex-direction: column;
  animation: slideUp 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: hidden;
}

[data-theme="light"] .modal-panel {
  border-color: rgba(60, 60, 67, 0.1);
  box-shadow: var(--shadow-dropdown), 0 0 40px rgba(60, 60, 67, 0.08);
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
  margin-bottom: 0.75rem;
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

@keyframes spin { to { transform: rotate(360deg); } }
.spinner-xs {
  display: inline-block;
  width: 0.7rem;
  height: 0.7rem;
  border: 1.5px solid var(--color-text-muted);
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  vertical-align: middle;
}

.about-section { text-align: center; }

.about-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.6rem;
}

.about-logo {
  font-family: var(--font-display);
  font-size: 1.6rem;
  font-weight: 900;
  background: var(--gradient-logo);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 4px;
}

.about-versions {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.3rem;
}

.app-version-row {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.app-version {
  font-family: var(--font-display);
  font-size: 1.2rem;
  font-weight: 700;
  color: var(--color-text-primary);
  letter-spacing: 1px;
}

.git-version-row {
  font-family: var(--font-body);
  font-size: 0.7rem;
  color: var(--color-nav-link);
  opacity: 0.5;
}

.update-action {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}
.update-check-link {
  cursor: pointer;
  color: var(--color-primary);
  transition: opacity 0.2s;
}
.update-check-link:hover { opacity: 0.7; }
.update-latest { color: var(--color-green); }
.update-new { color: var(--color-primary); }

.update-download-link-sm {
  padding: 0.1rem 0.5rem;
  background: var(--color-primary);
  color: #fff;
  border-radius: var(--radius-btn);
  font-size: 0.7rem;
  font-weight: 600;
  text-decoration: none;
  line-height: 1.5;
  transition: opacity 0.2s;
}
.update-download-link-sm:hover { opacity: 0.8; }

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px) scale(0.97); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}
</style>
