<template>
  <div class="app">
    <header class="header">
      <div class="logo-container">
        <h1 class="logo">GITSTAT</h1>
        <div class="logo-glow"></div>
      </div>
      <nav>
        <a @click="currentView = 'dashboard'" :class="{ active: currentView === 'dashboard' }">
          <span class="nav-icon">◈</span>
          Dashboard
        </a>
        <a @click="currentView = 'analytics'" :class="{ active: currentView === 'analytics' }">
          <span class="nav-icon">◉</span>
          Analytics
        </a>
        <a @click="currentView = 'settings'" :class="{ active: currentView === 'settings' }">
          <span class="nav-icon">⚙</span>
          Settings
        </a>
      </nav>
    </header>
    <main class="main-content">
      <Dashboard v-if="currentView === 'dashboard'" />
      <Analytics v-if="currentView === 'analytics'" />
      <Settings v-if="currentView === 'settings'" />
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Dashboard from './views/Dashboard.vue'
import Analytics from './views/Analytics.vue'
import Settings from './views/Settings.vue'

const currentView = ref('dashboard')
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

@keyframes slideIn {
  from { width: 0; }
  to { width: 100%; }
}

.main-content {
  flex: 1;
  padding: 2rem;
}
</style>
