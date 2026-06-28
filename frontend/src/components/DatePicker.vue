<template>
  <div class="dp-wrapper">
    <div class="dp-trigger btn" @click.stop="open = !open">
      <span class="dp-label">{{ label }}</span>
      <span class="dp-value" :class="{ empty: !start }">
        {{ start ? `${start} ~ ${end}` : placeholder || t('datePicker.placeholder') }}
      </span>
      <span class="dp-arrow">▾</span>
    </div>
    <Teleport to="body">
      <div v-if="open" class="dp-overlay" @click.stop="open = false">
        <div class="dp-dropdown" :style="dropdownStyle" @click.stop>
          <div class="dp-nav-row">
            <button @click="navMonth(-1)" class="dp-nav-btn btn" :aria-label="t('datePicker.prevMonth')">‹</button>
            <span class="dp-nav-title">{{ leftLabel }} — {{ rightLabel }}</span>
            <button @click="navMonth(1)" class="dp-nav-btn btn" :aria-label="t('datePicker.nextMonth')">›</button>
          </div>
          <div class="dp-months">
            <div class="dp-month">
              <div class="dp-month-label">{{ leftMonthLabel }}</div>
              <div class="dp-days">
                <span v-for="d in dayNames" :key="d" class="dp-day-hd">{{ d }}</span>
              </div>
              <div class="dp-grid" @mouseleave="hoverDate = ''">
                <div v-for="cell in leftGrid" :key="cell.dateStr"
                  class="dp-cell"
                  :class="cellCls(cell)"
                  @mouseenter="cell.day && (hoverDate = cell.dateStr)"
                  @click="cell.day && clickCell(cell.dateStr)"
                >{{ cell.day || '' }}</div>
              </div>
            </div>
            <div class="dp-month">
              <div class="dp-month-label">{{ rightMonthLabel }}</div>
              <div class="dp-days">
                <span v-for="d in dayNames" :key="d" class="dp-day-hd">{{ d }}</span>
              </div>
              <div class="dp-grid" @mouseleave="hoverDate = ''">
                <div v-for="cell in rightGrid" :key="cell.dateStr"
                  class="dp-cell"
                  :class="cellCls(cell)"
                  @mouseenter="cell.day && (hoverDate = cell.dateStr)"
                  @click="cell.day && clickCell(cell.dateStr)"
                >{{ cell.day || '' }}</div>
              </div>
            </div>
          </div>
          <div class="dp-range-info">
            <span class="dp-info-label">{{ t('datePicker.range') }}</span>
            <span class="dp-info-dates">{{ localStart || '____-__-__' }} ~ {{ localEnd || '____-__-__' }}</span>
          </div>
          <div class="dp-actions">
            <button class="dp-clear-btn btn" @click="clearRange">{{ t('datePicker.clear') }}</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, watch, onUnmounted } from 'vue'
import { useI18n } from '../i18n'
import { pad } from '../utils/dates'

const { t } = useI18n()

const props = defineProps({
  start: String,
  end: String,
  placeholder: { type: String, default: '' },
  label: { type: String, default: '' }
})
const emit = defineEmits(['update:start', 'update:end'])

const open = ref(false)
const localStart = ref(props.start || '')
const localEnd = ref(props.end || '')
const selecting = ref(null) // 'start' or 'end'
const hoverDate = ref('')
const viewYear = ref(new Date().getFullYear())
const viewMonth = ref(new Date().getMonth())
const dropdownStyle = ref({})
let positionTimer = null

onUnmounted(() => {
  if (positionTimer) clearTimeout(positionTimer)
})

watch(open, (v) => {
  if (v) {
    localStart.value = props.start || ''
    localEnd.value = props.end || ''
    selecting.value = 'start'
    const refDate = localStart.value || localEnd.value || ''
    const d = refDate ? new Date(refDate + 'T00:00:00') : new Date()
    viewYear.value = d.getFullYear()
    viewMonth.value = d.getMonth()
    positionDropdown()
  }
})

function positionDropdown() {
  positionTimer = setTimeout(() => {
    const el = document.querySelector('.dp-dropdown')
    const trigger = document.querySelector('.dp-trigger')
    if (el && trigger) {
      const rect = trigger.getBoundingClientRect()
      const left = Math.min(rect.left, Math.max(0, window.innerWidth - 540))
      dropdownStyle.value = {
        position: 'fixed',
        top: rect.bottom + 4 + 'px',
        left: left + 'px'
      }
    }
  })
}

function navMonth(delta) {
  viewMonth.value += delta
  if (viewMonth.value < 0) { viewMonth.value = 11; viewYear.value-- }
  if (viewMonth.value > 11) { viewMonth.value = 0; viewYear.value++ }
}

function rightYear() { return viewMonth.value === 11 ? viewYear.value + 1 : viewYear.value }
function rightMonth() { return viewMonth.value === 11 ? 0 : viewMonth.value + 1 }

const dayNames = computed(() => t('datePicker.dayNames'))

const fmtMonth = (y, m) => t('datePicker.monthFormat').replace('{y}', y).replace('{m}', m + 1)
const leftLabel = computed(() => fmtMonth(viewYear.value, viewMonth.value))
const rightLabel = computed(() => fmtMonth(rightYear(), rightMonth()))
const leftMonthLabel = computed(() => fmtMonth(viewYear.value, viewMonth.value))
const rightMonthLabel = computed(() => fmtMonth(rightYear(), rightMonth()))

function buildGrid(year, month) {
  const first = new Date(year, month, 1)
  const daysInMonth = new Date(year, month + 1, 0).getDate()
  let startCol = first.getDay() - 1
  if (startCol < 0) startCol = 6
  const today = new Date()
  const todayStr = `${today.getFullYear()}-${pad(today.getMonth() + 1)}-${pad(today.getDate())}`
  const cells = []
  for (let i = 0; i < startCol; i++) cells.push({ day: null, dateStr: '' })
  for (let d = 1; d <= daysInMonth; d++) {
    const dateStr = `${year}-${pad(month + 1)}-${pad(d)}`
    cells.push({ day: d, dateStr, isToday: dateStr === todayStr })
  }
  return cells
}

const leftGrid = computed(() => buildGrid(viewYear.value, viewMonth.value))
const rightGrid = computed(() => buildGrid(rightYear(), rightMonth()))

function cellCls(cell) {
  if (!cell.day) return { muted: true }
  const ds = cell.dateStr
  const ls = localStart.value
  const le = localEnd.value
  const rangeEnd = le || hoverDate.value
  const inRange = ls && rangeEnd && ds >= Math.min(ls, rangeEnd) && ds <= Math.max(ls, rangeEnd)
  const isStart = ds === ls
  const isEnd = ds === rangeEnd
  return {
    'in-range': inRange,
    'range-start': isStart && (!isEnd || ls === le),
    'range-end': isEnd && (!isStart || ls === le),
    'range-edge': isStart || (isEnd && !isStart),
    'today': cell.isToday,
    'hovering': !!hoverDate.value
  }
}

function clickCell(dateStr) {
  if (selecting.value === 'start') {
    localStart.value = dateStr
    localEnd.value = ''
    selecting.value = 'end'
  } else {
    if (dateStr < localStart.value) {
      localStart.value = dateStr
      localEnd.value = ''
      selecting.value = 'end'
    } else {
      localEnd.value = dateStr
      emit('update:start', localStart.value)
      emit('update:end', dateStr)
      open.value = false
    }
  }
}

function clearRange() {
  localStart.value = ''
  localEnd.value = ''
  hoverDate.value = ''
  selecting.value = 'start'
}
</script>

<style scoped>
.dp-wrapper {
  position: relative;
  font-family: var(--font-body);
}
.dp-trigger {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.5rem 0.6rem;
  min-width: 200px;
  white-space: nowrap;
}
.dp-label {
  font-size: 0.65rem;
  color: var(--color-text-muted);
  font-family: var(--font-display);
  letter-spacing: 1px;
  text-transform: uppercase;
}
.dp-value {
  flex: 1;
  font-size: 0.78rem;
  color: var(--color-input-text);
  font-weight: 600;
  font-variant-numeric: tabular-nums;
}
.dp-value.empty {
  color: var(--color-text-muted);
  font-weight: 400;
}
.dp-arrow {
  font-size: 0.6rem;
  color: var(--color-accent);
  transition: transform 0.2s ease;
}
.dp-overlay {
  position: fixed;
  inset: 0;
  z-index: 999;
}
.dp-dropdown {
  background: var(--bg-tooltip);
  backdrop-filter: blur(var(--blur-card));
  border: 1px solid rgba(var(--color-accent-rgb), 0.25);
  border-radius: 8px;
  padding: 0.75rem;
  box-shadow: var(--shadow-lg-black), var(--shadow-glow);
  animation: dpFadeIn 0.15s ease;
  width: 520px;
  z-index: 1000;
}
@keyframes dpFadeIn {
  from { opacity: 0; transform: translateY(-4px); }
  to { opacity: 1; transform: translateY(0); }
}
.dp-nav-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}
.dp-nav-btn {
  width: 28px;
  height: 28px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  border-radius: 4px;
}
.dp-nav-title {
  font-family: var(--font-display);
  font-size: 0.72rem;
  color: var(--color-accent);
  letter-spacing: 1px;
}
.dp-months {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}
.dp-month-label {
  font-family: var(--font-display);
  font-size: 0.68rem;
  color: var(--color-text-secondary);
  text-align: center;
  margin-bottom: 0.4rem;
  letter-spacing: 1px;
}
.dp-days {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 1px;
  margin-bottom: 4px;
}
.dp-day-hd {
  text-align: center;
  font-family: var(--font-display);
  font-size: 0.55rem;
  color: var(--color-text-muted);
  letter-spacing: 1px;
  text-transform: uppercase;
  padding: 3px 0;
}
.dp-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 1px;
}
.dp-cell {
  text-align: center;
  padding: 4px 0;
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.1s ease;
  position: relative;
}
.dp-cell:hover:not(.muted) {
  background: var(--border-cal-today);
  color: var(--color-input-text);
}
.dp-cell.hovering.in-range:not(.range-edge) {
  background: var(--bg-in-range-hover);
}
.dp-cell.today {
  color: var(--color-accent);
  font-weight: 700;
}
.dp-cell.in-range {
  background: var(--bg-in-range);
  color: var(--color-input-text);
  border-radius: 0;
}
.dp-cell.range-edge {
  background: var(--bg-range-edge);
  color: var(--color-white);
  font-weight: 700;
  border-radius: 4px;
}
.dp-cell.in-range:first-child {
  border-radius: 4px 0 0 4px;
}
.dp-cell.muted {
  cursor: default;
}
.dp-range-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  margin-top: 0.6rem;
  padding: 0.3rem 0.5rem;
  background: var(--bg-range-info);
  border-radius: 4px;
  font-size: 0.7rem;
  font-family: var(--font-display);
  color: var(--color-accent);
  letter-spacing: 0.5px;
}
.dp-info-label {
  color: var(--color-text-muted);
  font-size: 0.6rem;
  text-transform: uppercase;
}
.dp-info-dates {
  font-variant-numeric: tabular-nums;
  color: var(--color-input-text);
}
.dp-actions {
  display: flex;
  justify-content: center;
  margin-top: 0.5rem;
}
.dp-clear-btn {
  padding: 0.3rem 1rem;
  font-size: 0.72rem;
  border-radius: 4px;
}
</style>
