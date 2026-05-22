<template>
  <div class="scan-form card">
    <div class="form-header">
      <h3>{{ t('scan.title') }}</h3>
      <div class="tech-lines">
        <span></span>
        <span></span>
        <span></span>
      </div>
    </div>
    <div class="form-group">
      <label>{{ t('scan.directoryPath') }}:</label>
      <input v-model="scanPath" placeholder="/path/to/repos" />
      <div class="input-glow"></div>
    </div>
    <button @click="handleScan" :disabled="loading" class="btn scan-btn">
      <span v-if="!loading">{{ t('scan.startScan') }}</span>
      <span v-else class="loading-text">
        <span class="spinner"></span>
        {{ t('scan.scanning') }}
      </span>
    </button>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '../i18n'
import { performScan, state, loadScanPath } from '../stores/data'

const { t } = useI18n()
const scanPath = computed({
  get: () => state.scanPath,
  set: (v) => { state.scanPath = v }
})
const loading = ref(false)
const error = ref('')

const emit = defineEmits(['scan-complete'])

async function handleScan() {
  loading.value = true
  error.value = ''
  
  try {
    await performScan(state.scanPath, '1d')
    emit('scan-complete')
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.scan-form {
  max-width: 600px;
  margin-bottom: 2rem;
}

.form-header {
  position: relative;
  margin-bottom: 2rem;
}

.form-header h3 {
  font-family: 'Orbitron', sans-serif;
  font-size: 1.5rem;
  color: #00d4ff;
  letter-spacing: 2px;
  text-transform: uppercase;
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

.form-group {
  margin-bottom: 1.5rem;
  position: relative;
}

.form-group label {
  display: block;
  margin-bottom: 0.75rem;
  font-family: 'Orbitron', sans-serif;
  font-size: 0.85rem;
  color: #a0aec0;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid rgba(0, 212, 255, 0.3);
  border-radius: 8px;
  background: rgba(10, 14, 39, 0.6);
  color: #e0e6ff;
  font-family: 'Rajdhani', sans-serif;
  font-size: 1rem;
  transition: all 0.3s;
  position: relative;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #00d4ff;
  box-shadow: 0 0 20px rgba(0, 212, 255, 0.3);
  background: rgba(10, 14, 39, 0.8);
}

.input-glow {
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, #00d4ff, transparent);
  transition: width 0.3s;
}

.form-group input:focus ~ .input-glow {
  width: 100%;
}

.scan-btn {
  width: 100%;
  padding: 1rem;
  font-size: 1rem;
  position: relative;
}

.loading-text {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error {
  color: #ff4757;
  margin-top: 1rem;
  padding: 0.75rem;
  background: rgba(255, 71, 87, 0.1);
  border: 1px solid rgba(255, 71, 87, 0.3);
  border-radius: 8px;
  font-family: 'Rajdhani', sans-serif;
}
</style>
