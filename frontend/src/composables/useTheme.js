import { ref, watch } from 'vue'

const STORAGE_KEY = 'gitstat-theme'
const theme = ref(localStorage.getItem(STORAGE_KEY) || 'neon')

watch(theme, (val) => {
  localStorage.setItem(STORAGE_KEY, val)
  document.documentElement.dataset.theme = val === 'ios26' ? 'ios26' : ''
}, { immediate: true })

export function useTheme() {
  function toggleTheme() {
    theme.value = theme.value === 'neon' ? 'ios26' : 'neon'
  }

  function setTheme(val) {
    theme.value = val
  }

  return { theme, toggleTheme, setTheme }
}
