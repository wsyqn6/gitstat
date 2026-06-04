<template>
  <div class="app">
    <header class="header">
      <div class="logo-container">
        <h1 class="logo">GITSTAT</h1>
        <div class="logo-glow"></div>
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
      </nav>
    </header>
    <main class="main-content">
      <KeepAlive>
        <component :is="currentComponent" />
      </KeepAlive>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, defineAsyncComponent } from 'vue'
import { useI18n } from './i18n'

const Dashboard = defineAsyncComponent(() => import('./views/Dashboard.vue'))
const Analytics = defineAsyncComponent(() => import('./views/Analytics.vue'))
const RepoSection = defineAsyncComponent(() => import('./views/RepoSection.vue'))
const Settings = defineAsyncComponent(() => import('./views/Settings.vue'))

const componentMap = { dashboard: Dashboard, analytics: Analytics, repos: RepoSection, settings: Settings }

const { t, locale, setLocale } = useI18n()
const currentView = ref(localStorage.getItem('currentView') || 'dashboard')
const currentComponent = computed(() => componentMap[currentView.value])

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

.logo-container {
  position: relative;
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

@keyframes slideIn {
  from { width: 0; }
  to { width: 100%; }
}

.main-content {
  flex: 1;
  padding: 2rem;
}
</style>
