<template>
  <div class="settings">
    <h2 class="page-title">{{ t('settings.title') }}</h2>
    
    <div class="card settings-card">
      <div class="card-header">
        <h3>{{ t('settings.dataManagement') }}</h3>
        <div class="tech-lines">
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>
      <button @click="handleExport" class="btn export-btn">
        <span class="btn-icon">↓</span>
        {{ t('settings.exportData') }}
      </button>
    </div>
    
    <div class="card settings-card">
      <div class="card-header">
        <h3>{{ t('settings.about') }}</h3>
        <div class="tech-lines">
          <span></span>
          <span></span>
          <span></span>
        </div>
      </div>
      <div class="about-content">
        <div class="logo-text">GITSTAT</div>
        <p class="version">{{ t('settings.version') }}</p>
        <div class="tech-stack">
          <div class="tech-item">
            <span class="tech-icon">⚡</span>
            <span>Go 1.26</span>
          </div>
          <div class="tech-item">
            <span class="tech-icon">◈</span>
            <span>Vue 3</span>
          </div>
          <div class="tech-item">
            <span class="tech-icon">◉</span>
            <span>ECharts</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '../i18n'
import * as api from '../api'

const { t } = useI18n()

async function handleExport() {
  try {
    const blob = await api.exportData()
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'gitstat-data.json'
    a.click()
    window.URL.revokeObjectURL(url)
  } catch (err) {
    alert('导出失败: ' + err.message)
  }
}
</script>

<style scoped>
.settings {
  max-width: 800px;
}

.page-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 2rem;
  color: #00d4ff;
  letter-spacing: 3px;
  text-transform: uppercase;
  margin-bottom: 2rem;
  background: linear-gradient(135deg, #00d4ff, #7800ff);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.settings-card {
  margin-bottom: 2rem;
}

.card-header {
  position: relative;
  margin-bottom: 1.5rem;
}

.card-header h3 {
  font-family: 'Orbitron', sans-serif;
  font-size: 1.3rem;
  color: #00d4ff;
  letter-spacing: 2px;
  text-transform: uppercase;
  margin: 0;
}

.tech-lines {
  display: flex;
  gap: 8px;
  margin-top: 1rem;
}

.tech-lines span {
  height: 2px;
  background: linear-gradient(90deg, #00d4ff, transparent);
  animation: techLine 2s ease-in-out infinite;
}

.tech-lines span:nth-child(1) {
  width: 60px;
  animation-delay: 0s;
}

.tech-lines span:nth-child(2) {
  width: 40px;
  animation-delay: 0.2s;
}

.tech-lines span:nth-child(3) {
  width: 20px;
  animation-delay: 0.4s;
}

@keyframes techLine {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 1; }
}

.export-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  width: 100%;
}

.btn-icon {
  font-size: 1.2rem;
}

.about-content {
  text-align: center;
  padding: 2rem 0;
}

.logo-text {
  font-family: 'Orbitron', sans-serif;
  font-size: 2.5rem;
  font-weight: 900;
  background: linear-gradient(135deg, #00d4ff 0%, #7800ff 50%, #00ff88 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 4px;
  margin-bottom: 1rem;
}

.version {
  color: #a0aec0;
  font-size: 1rem;
  margin-bottom: 2rem;
}

.tech-stack {
  display: flex;
  justify-content: center;
  gap: 2rem;
  flex-wrap: wrap;
}

.tech-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: rgba(0, 212, 255, 0.1);
  border: 1px solid rgba(0, 212, 255, 0.3);
  border-radius: 8px;
  transition: all 0.3s;
}

.tech-item:hover {
  background: rgba(0, 212, 255, 0.2);
  border-color: #00d4ff;
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(0, 212, 255, 0.3);
}

.tech-icon {
  font-size: 1.2rem;
}
</style>
