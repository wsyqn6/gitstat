<template>
  <div class="app">
    <header class="header">
      <div class="header-left">
        <div class="logo-container">
          <h1 class="logo">GITSTAT</h1>
          <div class="logo-glow"></div>
        </div>
        <div class="header-meta" v-if="version || scanPath">
          <svg class="git-icon" viewBox="0 0 78 78" width="14" height="14">
            <path fill="currentColor" transform="translate(10 10) rotate(-45 29 29)" d="M5,58c-2.76142,0 -5,-2.23858 -5,-5v-48c0,-2.76142 2.23858,-5 5,-5h33v12.54404c-2.06553,0.94801 -3.5,3.03446 -3.5,5.45596c0,0.73514 0.13221,1.43941 0.37415,2.09031l-15.28384,15.28384c-0.6509,-0.24194 -1.35517,-0.37415 -2.09031,-0.37415c-3.31371,0 -6,2.68629 -6,6c0,3.31371 2.68629,6 6,6c3.31371,0 6,-2.68629 6,-6c0,-0.73514 -0.13221,-1.43941 -0.37415,-2.09031l14.87415,-14.87415l0,11.50851c-2.06553,0.94801 -3.5,3.03446 -3.5,5.45596c0,3.31371 2.68629,6 6,6c3.31371,0 6,-2.68629 6,-6c0,-2.42149 -1.43447,-4.50795 -3.5,-5.45596l0,-12.08808c2.06553,-0.94801 3.5,-3.03446 3.5,-5.45596c0,-2.42149 -1.43447,-4.50795 -3.5,-5.45596l0,-12.54404h10c2.76142,0 5,2.23858 5,5v48c0,2.76142 -2.23858,5 -5,5z"/>
          </svg>
          <span class="meta-version">{{ version || 'dev' }}</span>
          <span class="meta-sep">│</span>
          <span class="meta-path" :title="scanPath">{{ scanPath || '—' }}</span>
        </div>
      </div>
      <nav class="header-center">
        <a @click="setView('dashboard')" :class="{ active: currentView === 'dashboard' }">
          <span class="nav-icon">
            <svg viewBox="0 0 16 16" fill="currentColor">
              <path d="M8 1.5 14.5 8 8 14.5 1.5 8Z" fill="none" stroke="currentColor" stroke-width="1.5"/>
              <path d="M8 4.5 11.5 8 8 11.5 4.5 8Z"/>
            </svg>
          </span>
          {{ t('nav.dashboard') }}
        </a>
        <a @click="setView('analytics')" :class="{ active: currentView === 'analytics' }">
          <span class="nav-icon">
            <svg viewBox="0 0 16 16" fill="currentColor">
              <circle cx="8" cy="8" r="6.5" fill="none" stroke="currentColor" stroke-width="1.5"/>
              <circle cx="8" cy="8" r="2.5"/>
            </svg>
          </span>
          {{ t('nav.analytics') }}
        </a>
        <a @click="setView('repos')" :class="{ active: currentView === 'repos' }">
          <span class="nav-icon">
            <svg viewBox="0 0 16 16" fill="currentColor">
              <rect x="1.5" y="1.5" width="13" height="13" rx="2" fill="none" stroke="currentColor" stroke-width="1.5"/>
              <rect x="3" y="4.5" width="10" height="2" rx="0.5"/>
              <rect x="3" y="9.5" width="6" height="2" rx="0.5"/>
            </svg>
          </span>
          {{ t('nav.repos') }}
        </a>
      </nav>
      <div class="header-right">
        <button @click="showSettings = true" class="header-icon-btn btn-ghost" :title="t('nav.settings')">
          <span class="nav-icon">
            <svg viewBox="0 0 16 16" fill="currentColor">
              <circle cx="8" cy="8" r="4.5"/>
              <rect x="6.5" y=".5" width="3" height="2" rx=".3"/>
              <rect x="6.5" y="13.5" width="3" height="2" rx=".3"/>
              <rect x=".5" y="6.5" width="2" height="3" rx=".3"/>
              <rect x="13.5" y="6.5" width="2" height="3" rx=".3"/>
              <circle cx="8" cy="8" r="2" fill="var(--bg-card)"/>
            </svg>
          </span>
        </button>
        <a href="https://github.com/wsyqn6/gitstat" target="_blank" rel="noopener noreferrer" class="header-icon-btn btn-ghost" :title="t('nav.github')">
          <span class="nav-icon">
            <svg viewBox="0 0 16 16" fill="currentColor">
              <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
            </svg>
          </span>
        </a>
      </div>
    </header>
    <main class="main-content">
      <KeepAlive>
        <component :is="currentComponent" />
      </KeepAlive>
    </main>
    <SettingsModal v-if="showSettings" @close="showSettings = false" />
    <ToastContainer />
  </div>
</template>

<script setup>
import { ref, computed, defineAsyncComponent, onMounted } from 'vue'
import { useI18n } from './i18n'
import { getScanPath } from './api'
import { loadScanPath } from './stores/data'
import { useTheme } from './composables/useTheme'
import ToastContainer from './components/ToastContainer.vue'
import SettingsModal from './components/SettingsModal.vue'

const Dashboard = defineAsyncComponent(() => import('./views/Dashboard.vue'))
const Analytics = defineAsyncComponent(() => import('./views/Analytics.vue'))
const RepoSection = defineAsyncComponent(() => import('./views/RepoSection.vue'))

const componentMap = { dashboard: Dashboard, analytics: Analytics, repos: RepoSection }

const { t } = useI18n()
const { theme } = useTheme()
const currentView = ref(localStorage.getItem('currentView') || 'dashboard')
const currentComponent = computed(() => componentMap[currentView.value])
const showSettings = ref(false)
const scanPath = ref('')
const version = ref('')

onMounted(async () => {
  try {
    const info = await getScanPath()
    scanPath.value = info.path
    version.value = info.version
  } catch (err) {
    console.error('Failed to load scan info:', err)
    try {
      await loadScanPath()
      const info = await getScanPath()
      scanPath.value = info.path
      version.value = info.version
    } catch {}
  }
})

function setView(view) {
  currentView.value = view
  localStorage.setItem('currentView', view)
}

</script>

<style scoped>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: var(--bg-nav);
  backdrop-filter: blur(var(--blur-card));
  padding: 1.5rem 2rem;
  display: flex;
  align-items: center;
  border-bottom: 1px solid var(--border-nav);
  box-shadow: var(--shadow-nav);
  position: relative;
}

.header-left {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.logo-container {
  position: relative;
}

.header-meta {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.65rem;
  font-family: var(--font-body);
  color: var(--color-technical-text);
  margin-top: 0.15rem;
  max-width: 360px;
}

.git-icon {
  color: var(--color-git-icon);
  flex-shrink: 0;
}

.meta-version {
  font-family: var(--font-body);
  color: var(--color-text-secondary);
  white-space: nowrap;
}

.meta-sep {
  color: rgba(148, 163, 184, 0.25);
}

.meta-path {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}

.logo {
  margin: 0;
  font-family: var(--font-display);
  font-weight: 900;
  font-size: 2rem;
  background: var(--gradient-logo);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 3px;
  position: relative;
  z-index: 1;
}

.logo-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 150%;
  height: 150%;
  background: var(--gradient-nav-glow);
  filter: blur(20px);
  animation: pulse 3s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; transform: translate(-50%, -50%) scale(1); }
  50% { opacity: 0.8; transform: translate(-50%, -50%) scale(1.1); }
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1.5rem;
}

.header-right {
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.header-center a {
  color: var(--color-nav-link);
  text-decoration: none;
  cursor: pointer;
  font-family: var(--font-display);
  font-weight: 500;
  font-size: 0.9rem;
  letter-spacing: 1px;
  transition: all 0.3s;
  position: relative;
  padding: 0.5rem 0;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.nav-icon {
  display: inline-flex;
  align-items: center;
  font-size: 1.2rem;
  transition: transform 0.3s;
}

.header-center a:hover {
  color: var(--color-primary);
}

.header-center a:hover .nav-icon {
  transform: rotate(180deg);
}

.header-center a.active {
  color: var(--color-primary);
  font-weight: 700;
}

.header-center a.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: var(--gradient-footer-active);
  box-shadow: 0 0 10px rgba(var(--color-primary-rgb), 0.8);
  animation: slideIn 0.3s ease-out;
}

.header-icon-btn {
  padding: 0.4rem;
  line-height: 1;
}

.nav-icon svg {
  display: block;
  width: 1.2rem;
  height: 1.2rem;
}

@keyframes slideIn {
  from { width: 0; }
  to { width: 100%; }
}

.main-content {
  flex: 1;
  padding: 2rem;
}
</style>
