<template>
  <div class="app">
    <header :class="['header', { 'header--hidden': navHidden }]">
      <div class="header-left">
        <div class="logo-container" :class="{ 'logo--clicked': logoClickBoost }" @click="triggerCommitsplosion">
          <h1 class="logo">{{ eggLogoText || 'GITSTAT' }}</h1>
          <div class="logo-glow"></div>
        </div>
        <div class="header-meta" v-if="version || scanPath">
          <svg class="git-icon" viewBox="0 0 78 78" width="14" height="14">
            <path fill="currentColor" transform="translate(10 10) rotate(-45 29 29)" d="M5,58c-2.76142,0 -5,-2.23858 -5,-5v-48c0,-2.76142 2.23858,-5 5,-5h33v12.54404c-2.06553,0.94801 -3.5,3.03446 -3.5,5.45596c0,0.73514 0.13221,1.43941 0.37415,2.09031l-15.28384,15.28384c-0.6509,-0.24194 -1.35517,-0.37415 -2.09031,-0.37415c-3.31371,0 -6,2.68629 -6,6c0,3.31371 2.68629,6 6,6c3.31371,0 6,-2.68629 6,-6c0,-0.73514 -0.13221,-1.43941 -0.37415,-2.09031l14.87415,-14.87415l0,11.50851c-2.06553,0.94801 -3.5,3.03446 -3.5,5.45596c0,3.31371 2.68629,6 6,6c3.31371,0 6,-2.68629 6,-6c0,-2.42149 -1.43447,-4.50795 -3.5,-5.45596l0,-12.08808c2.06553,-0.94801 3.5,-3.03446 3.5,-5.45596c0,-2.42149 -1.43447,-4.50795 -3.5,-5.45596l0,-12.54404h10c2.76142,0 5,2.23858 5,5v48c0,2.76142 -2.23858,5 -5,5z"/>
          </svg>
          <span class="meta-version">{{ version || 'dev' }}</span>
          <span class="meta-sep">│</span>
          <span class="meta-path" :title="state.scanPath">{{ state.scanPath || '—' }}</span>
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
        <a @click="setView('trends')" :class="{ active: currentView === 'trends' }">
          <span class="nav-icon">
            <svg viewBox="0 0 16 16" fill="currentColor">
              <circle cx="8" cy="8" r="6.5" fill="none" stroke="currentColor" stroke-width="1.5"/>
              <circle cx="8" cy="8" r="2.5"/>
            </svg>
          </span>
          {{ t('nav.trends') }}
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
        <button @click="showSettings = true" class="header-icon-btn btn-ghost" :aria-label="t('nav.settings')" :title="t('nav.settings')">
          <span class="nav-icon">
            <svg viewBox="0 0 16 16" fill="currentColor">
              <path d="M9.796 1.343c-.527-1.79-3.065-1.79-3.592 0l-.094.319a.873.873 0 0 1-1.255.52l-.292-.16c-1.64-.892-3.433.902-2.54 2.541l.159.292a.873.873 0 0 1-.52 1.255l-.319.094c-1.79.527-1.79 3.065 0 3.592l.319.094a.873.873 0 0 1 .52 1.255l-.16.292c-.892 1.64.901 3.434 2.541 2.54l.292-.159a.873.873 0 0 1 1.255.52l.094.319c.527 1.79 3.065 1.79 3.592 0l.094-.319a.873.873 0 0 1 1.255-.52l.292.16c1.64.893 3.434-.902 2.54-2.541l-.159-.292a.873.873 0 0 1 .52-1.255l.319-.094c1.79-.527 1.79-3.065 0-3.592l-.319-.094a.873.873 0 0 1-.52-1.255l.16-.292c.893-1.64-.902-3.433-2.541-2.54l-.292.159a.873.873 0 0 1-1.255-.52zm-2.633.283c.246-.835 1.428-.835 1.674 0l.094.319a1.873 1.873 0 0 0 2.693 1.115l.291-.16c.764-.415 1.6.42 1.184 1.185l-.159.292a1.873 1.873 0 0 0 1.116 2.692l.318.094c.835.246.835 1.428 0 1.674l-.319.094a1.873 1.873 0 0 0-1.115 2.693l.16.291c.415.764-.42 1.6-1.185 1.184l-.291-.159a1.873 1.873 0 0 0-2.693 1.116l-.094.318c-.246.835-1.428.835-1.674 0l-.094-.319a1.873 1.873 0 0 0-2.692-1.115l-.292.16c-.764.415-1.6-.42-1.184-1.185l.159-.291A1.873 1.873 0 0 0 1.945 8.93l-.319-.094c-.835-.246-.835-1.428 0-1.674l.319-.094A1.873 1.873 0 0 0 3.06 4.377l-.16-.292c-.415-.764.42-1.6 1.185-1.184l.292.159a1.873 1.873 0 0 0 2.692-1.115z"/>
              <path d="M8 4.754a3.246 3.246 0 1 0 0 6.492 3.246 3.246 0 0 0 0-6.492M5.754 8a2.246 2.246 0 1 1 4.492 0 2.246 2.246 0 0 1-4.492 0"/>
            </svg>
          </span>
        </button>
        <a href="https://github.com/wsyqn6/gitstat" target="_blank" rel="noopener noreferrer" class="header-icon-btn btn-ghost" :aria-label="t('nav.github')" :title="t('nav.github')">
          <span class="nav-icon">
            <svg viewBox="0 0 16 16" fill="currentColor">
              <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"/>
            </svg>
          </span>
        </a>
      </div>
    </header>
    <div class="commitsplosion" v-if="particles.length">
      <div v-for="p in particles" :key="p.id" class="particle"
        :style="{
          left: p.ox+'px', top: p.oy+'px',
          width: p.size+'px', height: p.size+'px',
          '--tx': p.tx+'px', '--ty': p.ty+'px', '--rot': p.rot+'deg',
          '--delay': p.delay+'s', '--dur': p.dur+'s',
          background: p.color,
        }" />
    </div>
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
import { ref, computed, defineAsyncComponent, onMounted, onUnmounted } from 'vue'
import { useI18n } from './i18n'
import { initApp, state } from './stores/data'
import { useTheme } from './composables/useTheme'
import { useToast } from './composables/useToast'
import ToastContainer from './components/ToastContainer.vue'
import SettingsModal from './components/SettingsModal.vue'

const Dashboard = defineAsyncComponent(() => import('./views/Dashboard.vue'))
const Trends = defineAsyncComponent(() => import('./views/Trends.vue'))
const RepoSection = defineAsyncComponent(() => import('./views/RepoSection.vue'))

const componentMap = { dashboard: Dashboard, trends: Trends, repos: RepoSection }

const { t } = useI18n()
const { theme } = useTheme()
const savedView = localStorage.getItem('currentView')
const currentView = ref(savedView === 'analytics' ? 'trends' : (savedView || 'dashboard'))
const currentComponent = computed(() => componentMap[currentView.value])
const showSettings = ref(false)
  const version = ref('')
  const navHidden = ref(false)
  let lastScrollY = 0
  let ticking = false

  function handleScroll() {
    if (!ticking) {
      window.requestAnimationFrame(() => {
        const currentScrollY = window.scrollY
        const delta = currentScrollY - lastScrollY
        const threshold = 60
        if (currentScrollY > threshold && delta > 15) {
          navHidden.value = true
        } else if (delta < -25) {
          navHidden.value = false
        }
        lastScrollY = currentScrollY
        ticking = false
      })
      ticking = true
    }
  }

  onMounted(async () => {
    await initApp()
    version.value = state.gitVersion
    window.addEventListener('scroll', handleScroll, { passive: true })
  })

  onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll)
  })

function setView(view) {
  currentView.value = view
  localStorage.setItem('currentView', view)
}

const { show } = useToast()
const particles = ref([])
const eggLogoText = ref('')
let pid = 0
const logoClickBoost = ref(false)
let eggClicks = 0
let eggClickTimer = 0
const EGG_CLICK_THRESHOLD = 5

const EGG_TEXTS = ['COMMIT', 'PUSH', 'MERGE', 'REBASE']

const EGG_MSGS = [
  '+42 contributions today! 🎉', '100% code coverage! 🧪',
  'No bugs found! (just kidding) 🐛', 'git push --force 🚀',
  'Merge conflict resolved! ⚔️', "You're a 10x engineer ⚡",
  'Rebase completed without drama ✨', 'All tests pass! (probably) ✅',
  'Zero dependencies added 🏆', 'Commit squashed successfully 🥫',
]

function triggerCommitsplosion(e) {
  if (particles.value.length) return
  eggClicks++
  clearTimeout(eggClickTimer)
  logoClickBoost.value = true
  setTimeout(() => { logoClickBoost.value = false }, 250)
  if (eggClicks < EGG_CLICK_THRESHOLD) {
    eggClickTimer = setTimeout(() => { eggClicks = 0 }, 800)
    return
  }
  eggClicks = 0
  const rect = e.currentTarget.getBoundingClientRect()
  const cx = rect.left + rect.width / 2
  const cy = rect.top + rect.height / 2

  const COLORS = ['#0e4429', '#006d32', '#26a641', '#39d353', '#ff6b6b', '#ffd93d']
  const items = []
  for (let i = 0; i < 36; i++) {
    const angle = (Math.PI * 2 * i) / 36 + (Math.random() - 0.5) * 0.5
    const dist = 80 + Math.random() * 180
    const size = 6 + Math.random() * 8
    items.push({
      id: ++pid,
      ox: cx - size / 2,
      oy: cy - size / 2,
      tx: Math.cos(angle) * dist,
      ty: Math.sin(angle) * dist,
      rot: (Math.random() - 0.5) * 720,
      delay: Math.random() * 0.15,
      dur: 0.8 + Math.random() * 0.5,
      size,
      color: COLORS[Math.floor(Math.random() * COLORS.length)],
    })
  }
  particles.value = items

  let i = 0
  eggLogoText.value = EGG_TEXTS[i]
  const iv = setInterval(() => {
    i++
    if (i >= EGG_TEXTS.length) {
      clearInterval(iv)
      eggLogoText.value = ''
      return
    }
    eggLogoText.value = EGG_TEXTS[i]
  }, 600)

  setTimeout(() => { particles.value = [] }, 2500)

  show(EGG_MSGS[Math.floor(Math.random() * EGG_MSGS.length)])
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
  backdrop-filter: blur(var(--blur-nav));
  padding: 1rem 2rem;
  display: flex;
  align-items: center;
  border-bottom: 1px solid var(--border-nav);
  box-shadow: var(--shadow-nav);
  position: sticky;
  top: 0;
  z-index: 100;
  transition: transform 0.35s cubic-bezier(0.4, 0, 0.2, 1);
}

.header--hidden {
  transform: translateY(-100%);
}

.header-left {
  display: flex;
  flex-direction: row;
  align-items: flex-end;
  gap: 0.5rem;
}

.logo-container {
  position: relative;
  user-select: none;
  -webkit-user-select: none;
}

.header-meta {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  font-size: 0.6rem;
  line-height: 1;
  font-family: var(--font-body);
  color: var(--color-technical-text);
  margin-top: 0;
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
  color: var(--color-separator);
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
  font-size: 1.75rem;
  line-height: 1;
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
  width: 130%;
  height: 130%;
  background: var(--gradient-nav-glow);
  filter: blur(20px);
  animation: pulse 3s ease-in-out infinite;
  transition: opacity 0.15s, transform 0.15s;
}

.logo--clicked .logo-glow {
  opacity: 1 !important;
  transform: translate(-50%, -50%) scale(1.3);
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; transform: translate(-50%, -50%) scale(1); }
  50% { opacity: 0.8; transform: translate(-50%, -50%) scale(1.1); }
}

.header-center {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1.5rem;
}

.header-right {
  margin-left: auto;
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
  font-size: 0.8rem;
  letter-spacing: 1px;
  transition: all 0.3s;
  position: relative;
  padding: 0.3rem 0;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.nav-icon {
  display: inline-flex;
  align-items: center;
  font-size: 0.9rem;
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
  width: 0.9rem;
  height: 0.9rem;
}

@keyframes slideIn {
  from { width: 0; }
  to { width: 100%; }
}

.main-content {
  flex: 1;
  padding: 2rem;
}

.commitsplosion {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 9999;
}

.particle {
  position: fixed;
  border-radius: 2px;
  animation: particleFly var(--dur) ease-out var(--delay) forwards;
  opacity: 0;
}

@keyframes particleFly {
  0% {
    transform: translate(0, 0) scale(0);
    opacity: 1;
  }
  15% {
    opacity: 1;
    transform: translate(calc(var(--tx) * 0.3), calc(var(--ty) * 0.3)) scale(1) rotate(calc(var(--rot) * 0.3));
  }
  60% {
    opacity: 1;
    transform: translate(var(--tx), var(--ty)) rotate(var(--rot));
  }
  100% {
    opacity: 0;
    transform: translate(var(--tx), var(--ty)) rotate(var(--rot)) scale(0.3);
  }
}
</style>
