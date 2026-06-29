<template>
  <div class="calendar-area">
    <CalendarWeekView v-if="viewType === 'week'" :dailyStats="dailyStats" :startDate="startDate" />
    <CalendarMonthView v-if="viewType === 'month'" :dailyStats="dailyStats" :startDate="startDate" />
    <CalendarYearView v-if="viewType === 'year'" :monthlyCalendar="monthlyCalendar" :dailyStats="dailyStats" :startDate="startDate" />
  </div>
</template>

<script setup>
import CalendarWeekView from './CalendarWeekView.vue'
import CalendarMonthView from './CalendarMonthView.vue'
import CalendarYearView from './CalendarYearView.vue'

defineProps({
  viewType: { type: String, required: true },
  dailyStats: { type: Array, default: () => [] },
  monthlyCalendar: { type: Array, default: () => [] },
  startDate: { type: String, default: '' }
})
</script>

<style scoped>
.calendar-area {
  animation: slideDown 0.3s ease;
}
</style>
