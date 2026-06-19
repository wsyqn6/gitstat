<template>
  <Teleport to="body">
    <div class="toast-container">
      <TransitionGroup name="toast">
        <div v-for="t in toasts" :key="t.id" :class="['toast', 'toast-' + t.type]">
          <span class="toast-icon">{{ icons[t.type] }}</span>
          <span class="toast-msg">{{ t.message }}</span>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>
import { useToast } from '../composables/useToast'

const { toasts } = useToast()

const icons = { success: '✓', error: '✗', info: 'ℹ' }
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 1rem;
  right: 1rem;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  pointer-events: none;
}

.toast {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.25rem;
  border-radius: 8px;
  font-family: var(--font-body);
  font-size: 0.9rem;
  font-weight: 600;
  backdrop-filter: blur(var(--blur-card));
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  pointer-events: auto;
  min-width: 200px;
}

.toast-success {
  background: var(--bg-toast-success);
  border: 1px solid var(--border-badge-success);
  color: var(--color-green);
}

.toast-error {
  background: var(--bg-toast-error);
  border: 1px solid rgba(255, 107, 107, 0.3);
  color: var(--color-red);
}

.toast-info {
  background: var(--bg-toast-info);
  border: 1px solid var(--border-dropdown);
  color: var(--color-accent);
}

.toast-icon {
  font-size: 1.1rem;
  flex-shrink: 0;
}

.toast-enter-active { animation: toastIn 0.25s ease; }
.toast-leave-active { animation: toastOut 0.2s ease; }

@keyframes toastIn {
  from { opacity: 0; transform: translateX(100%); }
  to { opacity: 1; transform: translateX(0); }
}

@keyframes toastOut {
  from { opacity: 1; transform: translateX(0); }
  to { opacity: 0; transform: translateX(100%); }
}
</style>
