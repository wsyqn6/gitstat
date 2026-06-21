<template>
  <div class="controls-section">
    <div class="controls-inline">
      <div class="control-item">
        <label class="control-label">{{ t('analytics.timeRange') }}</label>
        <div class="time-selector-cyber">
          <div class="quick-options-cyber">
            <button
              v-for="option in timeOptions"
              :key="option.value"
              @click="selectTimeRange(option.value)"
               class="cyber-time-btn btn"
              :class="{ active: selectedTimeRange === option.value }"
            >
              <span class="btn-text-cyber">{{ t(option.labelKey) }}</span>
            </button>
            <button
              @click.stop="emit('update:showCustomPicker', !showCustomPicker)"
               class="cyber-time-btn custom-btn btn"
              :class="{ active: selectedTimeRange === 'custom' }"
            >
              <span class="btn-text-cyber">{{ t('analytics.customPeriod') }}</span>
            </button>
          </div>

          <div v-if="showCustomPicker" class="custom-picker-popup" @click.stop>
            <DatePicker
              :start="customStartDate"
              :end="customEndDate"
              @update:start="onCustomStart"
              @update:end="onCustomEnd"
            />
          </div>
        </div>
      </div>

      <div class="control-item">
        <label class="control-label">{{ t('analytics.repoFilter') }}</label>
        <div v-if="repositories.length > 1" class="repo-dropdown">
          <button
            @click.stop="showDropdown = !showDropdown"
            class="repo-dropdown-btn btn"
            :class="{ active: showDropdown, 'has-selection': selectedRepos.length > 0 }"
          >
            <span class="btn-text">
              {{ selectedRepos.length === 0 ? t('analytics.selectRepo') :
                 selectedRepos.length === repositories.length ? t('analytics.allRepos') :
                 `${selectedRepos.length} ${t('analytics.reposSelected')}` }}
            </span>
            <svg class="dropdown-icon" :class="{ rotated: showDropdown }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="6 9 12 15 18 9"></polyline>
            </svg>
          </button>

          <div v-show="showDropdown" class="repo-dropdown-menu">
            <div class="dropdown-header">
              <button @click.stop="toggleAllRepos" class="select-all-btn btn-ghost">
                {{ allReposSelected ? t('analytics.cancelSelectAll') : t('analytics.selectAll') }}
              </button>
              <span class="selected-count">{{ selectedRepos.length }}/{{ repositories.length }}</span>
            </div>
            <div class="dropdown-list">
              <button
                v-for="repo in repositories"
                :key="repo.path"
                @click.stop="toggleRepo(repo.path)"
                class="repo-option btn-ghost"
                :class="{ active: selectedRepos.includes(repo.path) }"
              >
                <div class="option-checkbox">
                  <svg v-if="selectedRepos.includes(repo.path)" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                </div>
                <span class="option-text">{{ repo.name }}</span>
              </button>
            </div>
          </div>
        </div>
        <div v-else-if="repositories.length === 1" class="single-repo-label">
          {{ repositories[0].name }}
        </div>
      </div>

      <div class="control-item">
        <button
          @click="emit('analyze')"
          :disabled="loading || selectedRepos.length === 0"
          class="analyze-btn btn"
          :class="{ 'can-analyze': selectedRepos.length > 0 }"
        >
          <svg v-if="!loading" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
          </svg>
          <div v-else class="btn-spinner"></div>
          <span>{{ loading ? t('analytics.analyzing') : t('analytics.startAnalyze') }}</span>
        </button>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from '../i18n'
import DatePicker from './DatePicker.vue'

const { t } = useI18n()

const props = defineProps({
  repositories: Array,
  selectedRepos: Array,
  selectedTimeRange: String,
  showCustomPicker: Boolean,
  customStartDate: String,
  customEndDate: String,
  loading: Boolean
})

const emit = defineEmits([
  'update:selectedRepos',
  'update:selectedTimeRange',
  'update:showCustomPicker',
  'update:customStartDate',
  'update:customEndDate',
  'analyze'
])

const showDropdown = ref(false)

const timeOptions = [
  { labelKey: 'analytics.thisWeek', value: 'week' },
  { labelKey: 'analytics.thisMonth', value: 'month' },
  { labelKey: 'analytics.thisYear', value: 'year' }
]

const allReposSelected = computed(() => {
  return props.repositories.length > 0 &&
    props.selectedRepos.length === props.repositories.length
})

function selectTimeRange(value) {
  emit('update:selectedTimeRange', value)
  emit('update:customStartDate', '')
  emit('update:customEndDate', '')
  emit('update:showCustomPicker', false)
}

function onCustomStart(v) {
  emit('update:customStartDate', v)
  if (props.customEndDate) {
    emit('update:customEndDate', '')
    emit('update:showCustomPicker', false)
  }
}

function onCustomEnd(v) {
  emit('update:customEndDate', v)
  emit('update:selectedTimeRange', 'custom')
  emit('update:showCustomPicker', false)
  emit('analyze')
}

function toggleRepo(repoPath) {
  const idx = props.selectedRepos.indexOf(repoPath)
  if (idx > -1) {
    const next = [...props.selectedRepos]
    next.splice(idx, 1)
    emit('update:selectedRepos', next)
  } else {
    emit('update:selectedRepos', [...props.selectedRepos, repoPath])
  }
}

function toggleAllRepos() {
  if (allReposSelected.value) {
    emit('update:selectedRepos', [])
  } else {
    emit('update:selectedRepos', props.repositories.map(r => r.path))
  }
}

function closeDropdown(e) {
  if (showDropdown.value) {
    const el = e.target.closest('.repo-dropdown')
    if (!el) showDropdown.value = false
  }
}

function closeCustomPicker(e) {
  if (props.showCustomPicker) {
    const picker = e.target.closest('.custom-picker-popup')
    const btn = e.target.closest('.custom-btn')
    if (!picker && !btn) emit('update:showCustomPicker', false)
  }
}

onMounted(() => {
  document.addEventListener('click', closeDropdown)
  document.addEventListener('click', closeCustomPicker)
})

onUnmounted(() => {
  document.removeEventListener('click', closeDropdown)
  document.removeEventListener('click', closeCustomPicker)
})
</script>

<style scoped>
.controls-section {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.controls-inline {
  display: flex;
  gap: 1.5rem;
  align-items: flex-end;
}

.control-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.control-label {
  font-family: var(--font-display);
  font-size: 0.75rem;
  color: var(--color-accent);
  letter-spacing: 2px;
  text-transform: uppercase;
  font-weight: 600;
}

.analyze-btn {
  border-radius: 8px;
  padding: 0.6rem 1.5rem;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s ease;
  outline: none;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  letter-spacing: 1px;
  height: fit-content;
  align-self: flex-end;
  opacity: 0.5;
  cursor: not-allowed;
}

.analyze-btn.can-analyze {
  opacity: 1;
  cursor: pointer;
}

.analyze-btn.can-analyze:hover:not(:disabled) {
  transform: translateY(var(--btn-hover-lift));
}

.analyze-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.analyze-btn svg {
  width: 18px;
  height: 18px;
}

.btn-spinner {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.time-selector-cyber {
  position: relative;
  display: flex;
  flex-direction: row;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.6rem;
}

.quick-options-cyber {
  display: flex;
  gap: 0.6rem;
}

.cyber-time-btn {
  border-radius: 4px;
  min-width: 80px;
}

.cyber-time-btn.custom-btn {
  min-width: 70px;
  font-size: 0.78rem;
}

.custom-picker-popup {
  position: absolute;
  top: 100%;
  right: 0;
  z-index: 100;
  margin-top: 0.5rem;
  padding: 0.25rem;
  animation: slideDown 0.2s ease;
}

.repo-dropdown {
  position: relative;
}

.repo-dropdown-btn {
  min-width: 180px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
}

.btn-text {
  flex: 1;
  text-align: left;
}

.dropdown-icon {
  width: 16px;
  height: 16px;
  transition: transform 0.3s ease;
  color: var(--color-accent);
}

.dropdown-icon.rotated {
  transform: rotate(180deg);
}

.repo-dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  right: 0;
  background: var(--bg-dropdown-full);
  backdrop-filter: blur(var(--blur-card));
  border: 2px solid var(--border-dropdown);
  border-radius: 12px;
  box-shadow: var(--shadow-lg-black), 0 0 30px var(--border-card);
  z-index: 1000;
  overflow: hidden;
  animation: dropdownSlide 0.2s ease;
}

@keyframes dropdownSlide {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--bg-insight-hover);
  background: var(--bg-tag-section);
}

.select-all-btn {
  padding: 0.35rem 0.75rem;
  font-size: 0.8rem;
  border: 1px solid var(--color-accent);
  border-radius: 6px;
}

.selected-count {
  font-size: 0.8rem;
  color: var(--color-text-muted);
  font-family: var(--font-display);
}

.dropdown-list {
  max-height: 300px;
  overflow-y: auto;
  padding: 0.5rem;
}

.dropdown-list::-webkit-scrollbar {
  width: 6px;
}

.dropdown-list::-webkit-scrollbar-track {
  background: var(--bg-detail-msg);
}

.dropdown-list::-webkit-scrollbar-thumb {
  background: var(--border-dropdown);
  border-radius: 3px;
}

.dropdown-list::-webkit-scrollbar-thumb:hover {
  background: var(--border-card-hover);
}

.repo-option {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.6rem 0.75rem;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: var(--color-text-secondary);
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
  font-family: var(--font-body);
}

.repo-option:hover {
  background: var(--bg-insight-hover);
  color: var(--color-input-text);
}

.repo-option.active {
  background: var(--bg-stat);
  color: var(--color-accent);
}

.option-checkbox {
  width: 18px;
  height: 18px;
  border: 2px solid rgba(148, 163, 184, 0.4);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.repo-option.active .option-checkbox {
  border-color: var(--color-accent);
  background: var(--border-card);
}

.option-checkbox svg {
  width: 12px;
  height: 12px;
  color: var(--color-accent);
}

.option-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.single-repo-label {
  background: var(--bg-panel);
  border: 2px solid var(--border-dropdown);
  border-radius: 8px;
  padding: 0.6rem 1rem;
  color: var(--color-accent);
  font-size: 0.9rem;
  font-family: var(--font-body);
  min-width: 180px;
  display: flex;
  align-items: center;
}

</style>
