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
      <nav>
        <a @click="setView('dashboard')" :class="{ active: currentView === 'dashboard' }">
          <span class="nav-icon">◈</span>
          {{ t('nav.dashboard') }}
        </a>
        <a @click="setView('analytics')" :class="{ active: currentView === 'analytics' }">
          <span class="nav-icon">◉</span>
          {{ t('nav.analytics') }}
        </a>
        <a @click="setView('repos')" :class="{ active: currentView === 'repos' }">
          <span class="nav-icon">▤</span>
          {{ t('nav.repos') }}
        </a>
        <a @click="setView('settings')" :class="{ active: currentView === 'settings' }">
          <span class="nav-icon">⚙</span>
          {{ t('nav.settings') }}
        </a>
        <button @click="toggleLanguage" class="lang-switcher" :title="locale === 'zh' ? 'Switch to English' : '切换到中文'">
          <span class="lang-icon">{{ locale === 'zh' ? '中' : 'EN' }}</span>
        </button>
        <a href="https://github.com/wsyqn6/gitstat" target="_blank" rel="noopener noreferrer" class="github-link" :title="t('nav.github')">
          <svg viewBox="0 0 16 16" width="18" height="18" fill="currentColor">
            <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
          </svg>
        </a>
      </nav>
    </header>
    <main class="main-content">
      <KeepAlive>
        <component :is="currentComponent" />
      </KeepAlive>
    </main>
    <ToastContainer />
  </div>
</template>

<script setup>
import { ref, computed, defineAsyncComponent, onMounted } from 'vue'
import { useI18n } from './i18n'
import * as api from './api'
import { loadScanPath } from './stores/data'
import ToastContainer from './components/ToastContainer.vue'

const Dashboard = defineAsyncComponent(() => import('./views/Dashboard.vue'))
const Analytics = defineAsyncComponent(() => import('./views/Analytics.vue'))
const RepoSection = defineAsyncComponent(() => import('./views/RepoSection.vue'))
const Settings = defineAsyncComponent(() => import('./views/Settings.vue'))

const componentMap = { dashboard: Dashboard, analytics: Analytics, repos: RepoSection, settings: Settings }

const { t, locale, setLocale } = useI18n()
const currentView = ref(localStorage.getItem('currentView') || 'dashboard')
const currentComponent = computed(() => componentMap[currentView.value])
const scanPath = ref('')
const version = ref('')

onMounted(async () => {
  try {
    const info = await api.getScanPath()
    scanPath.value = info.path
    version.value = info.version
  } catch (err) {
    console.error('Failed to load scan info:', err)
    try {
      await loadScanPath()
      const info = await api.getScanPath()
      scanPath.value = info.path
      version.value = info.version
    } catch {}
  }
})

function setView(view) {
  currentView.value = view
  localStorage.setItem('currentView', view)
}

function toggleLanguage() {
  setLocale(locale.value === 'zh' ? 'en' : 'zh')
}
</script>

<style scoped>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: rgba(10, 14, 39, 0.8);
  backdrop-filter: blur(20px);
  padding: 1.5rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(0, 212, 255, 0.2);
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.3);
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
  font-family: 'Rajdhani', sans-serif;
  color: #475569;
  margin-top: 0.15rem;
  max-width: 360px;
}

.git-icon {
  color: #f05033;
  flex-shrink: 0;
}

.meta-version {
  font-family: 'Rajdhani', sans-serif;
  color: #94a3b8;
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
  font-family: 'Orbitron', sans-serif;
  font-weight: 900;
  font-size: 2rem;
  background: linear-gradient(135deg, #00d4ff 0%, #7800ff 50%, #00ff88 100%);
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
  background: radial-gradient(circle, rgba(0, 212, 255, 0.3) 0%, transparent 70%);
  filter: blur(20px);
  animation: pulse 3s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; transform: translate(-50%, -50%) scale(1); }
  50% { opacity: 0.8; transform: translate(-50%, -50%) scale(1.1); }
}

.header nav a {
  color: #a0aec0;
  text-decoration: none;
  margin-left: 2rem;
  cursor: pointer;
  font-family: 'Orbitron', sans-serif;
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
  font-size: 1.2rem;
  transition: transform 0.3s;
}

.header nav a:hover {
  color: #00d4ff;
}

.header nav a:hover .nav-icon {
  transform: rotate(180deg);
}

.header nav a.active {
  color: #00d4ff;
  font-weight: 700;
}

.header nav a.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 2px;
  background: linear-gradient(90deg, #00d4ff, #7800ff);
  box-shadow: 0 0 10px rgba(0, 212, 255, 0.8);
  animation: slideIn 0.3s ease-out;
}

.lang-switcher {
  margin-left: 1.5rem;
  background: rgba(0, 212, 255, 0.1);
  border: 1px solid rgba(0, 212, 255, 0.3);
  border-radius: 6px;
  padding: 0.4rem 0.8rem;
  cursor: pointer;
  transition: all 0.3s;
  font-family: 'Orbitron', sans-serif;
  font-size: 0.8rem;
  color: #00d4ff;
  letter-spacing: 1px;
}

.lang-switcher:hover {
  background: rgba(0, 212, 255, 0.2);
  border-color: #00d4ff;
  box-shadow: 0 0 15px rgba(0, 212, 255, 0.3);
  transform: translateY(-1px);
}

.lang-icon {
  font-weight: 700;
}

.github-link {
  margin-left: 0.5rem;
  border: none;
  border-radius: 6px;
  padding: 0.4rem 0.5rem;
  cursor: pointer;
  transition: all 0.3s;
  color: #64748b;
  display: inline-flex;
  align-items: center;
  text-decoration: none;
}

.github-link:hover {
  color: #00d4ff;
  background: rgba(0, 212, 255, 0.1);
  box-shadow: 0 0 12px rgba(0, 212, 255, 0.2);
  transform: translateY(-1px);
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
