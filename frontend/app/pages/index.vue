<script setup lang="ts">
import { computed, ref } from 'vue'
import { useHolidays } from '~/composables/useHolidays'

const { holidaysMap, isHoliday, getHolidayInfo, monthHolidays, error, refresh, data } = useHolidays()

const currentDate = ref(new Date())
const selectedDate = ref<Date | null>(null)

const currentMonth = computed(() => currentDate.value.getMonth())
const currentYear = computed(() => currentDate.value.getFullYear())

const monthNames = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December'
]
const dayNames = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']

const calendarDays = computed<(number | null)[]>(() => {
  const first = new Date(currentYear.value, currentMonth.value, 1)
  const last = new Date(currentYear.value, currentMonth.value + 1, 0)
  const lead = first.getDay()
  const days: (number | null)[] = Array(lead).fill(null)
  for (let d = 1; d <= last.getDate(); d++) days.push(d)
  return days
})

const thisMonthHolidays = computed(() =>
  monthHolidays(holidaysMap.value, currentYear.value, currentMonth.value)
)

function previousMonth() { currentDate.value = new Date(currentYear.value, currentMonth.value - 1, 1) }
function nextMonth() { currentDate.value = new Date(currentYear.value, currentMonth.value + 1, 1) }
function goToToday() { currentDate.value = new Date() }
function selectDate(day: number | null) { if (day) selectedDate.value = new Date(currentYear.value, currentMonth.value, day) }

function isToday(day: number | null) {
  if (!day) return false
  const t = new Date()
  return day === t.getDate() && currentMonth.value === t.getMonth() && currentYear.value === t.getFullYear()
}
function formatDateIndonesian(date: Date) {
  return date.toLocaleDateString('id-ID', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
}
</script>

<template>
  <div class="container mx-auto space-y-6 p-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold tracking-tight">Dashboard Calendar</h1>
        <p class="text-muted-foreground">Indonesian holidays</p>
      </div>
      <button class="border-input px-3 py-2 rounded-md text-sm hover:bg-accent hover:text-accent-foreground border"
        @click="goToToday">Today</button>
    </div>

    <!-- Error / Loading -->
    <div v-if="error || (data && data.status !== 'success')">
      <div class="border rounded-lg p-6 bg-red-50">
        <p class="text-red-600 font-medium mb-3">
          Error:
          <span class="font-normal">
            {{ (data && data.message) || (error && (error as any).message) || 'Unknown error' }}
          </span>
        </p>
        <button class="border px-3 py-2 rounded-md text-sm" @click="refresh()">Reload</button>
      </div>
    </div>

    <div v-else>
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <!-- Calendar -->
        <div class="lg:col-span-2 border rounded-2xl">
          <div class="border-b p-4">
            <div class="flex items-center justify-between">
              <h2 class="text-2xl font-semibold">
                {{ monthNames[currentMonth] }} {{ currentYear }}
              </h2>
              <div class="flex items-center gap-2">
                <button
                  class="border px-2 py-1 rounded-md transition-colors hover:bg-accent hover:text-accent-foreground"
                  @click="previousMonth" aria-label="Previous month">‹</button>
                <button
                  class="border px-2 py-1 rounded-md transition-colors hover:bg-accent hover:text-accent-foreground"
                  @click="nextMonth" aria-label="Next month">›</button>
              </div>
            </div>
          </div>

          <div class="p-3 max-w-3xl mx-auto">
            <div class="mb-3 grid grid-cols-7 gap-1">
              <div v-for="d in dayNames" :key="d" class="text-muted-foreground p-1.5 text-center text-sm font-medium">{{
                d
                }}</div>
            </div>

            <div class="grid grid-cols-7 gap-1">
              <div v-for="(day, idx) in calendarDays" :key="idx" class="aspect-square">
                <template v-if="day">
                  <button type="button" @click="selectDate(day)"
                    class="relative flex h-full w-full flex-col items-center justify-center rounded-md border p-1 transition-all hover:bg-accent"
                    :class="[
                      isToday(day) ? 'border-border bg-primary/5' : 'border-transparent',
                      selectedDate && selectedDate.getDate() === day &&
                        selectedDate.getMonth() === currentMonth &&
                        selectedDate.getFullYear() === currentYear
                        ? 'ring-2 ring-ring'
                        : ''
                    ]">
                    <span class="mb-1 text-sm font-medium">{{ day }}</span>
                    <div class="flex flex-col items-center space-y-1">
                      <div v-if="isHoliday(holidaysMap, currentYear, currentMonth, day)"
                        class="h-2 w-2 rounded-full bg-destructive"></div>
                    </div>
                  </button>
                </template>
                <template v-else>
                  <div class="h-full w-full"></div>
                </template>
              </div>
            </div>
          </div>
        </div>

        <!-- Holiday List -->
        <div class="border rounded-2xl">
          <div class="border-b p-4">
            <h3 class="text-lg font-semibold">This Month</h3>
          </div>
          <div class="p-4 space-y-3">
            <div class="space-y-2 border-b pb-3">
              <p class="text-sm font-medium">Legend:</p>
              <div class="flex items-center space-x-2 text-xs">
                <div class="h-2 w-2 rounded-full bg-destructive"></div>
                <span class="text-muted-foreground">Holiday</span>
              </div>
            </div>

            <template v-if="thisMonthHolidays.length > 0">
              <div>
                <p class="mb-2 text-sm font-medium">Holidays:</p>
                <div v-for="h in thisMonthHolidays" :key="h.dateStr"
                  class="mb-2 flex items-start space-x-3 rounded-lg border p-2">
                  <div class="flex-shrink-0">
                    <span
                      class="inline-flex items-center justify-center h-6 w-6 rounded-full bg-destructive text-destructive-foreground text-xs font-semibold">
                      {{ h.date }}
                    </span>
                  </div>
                  <div class="min-w-0 flex-1">
                    <p class="text-xs leading-tight font-medium">{{ h.summary }}</p>
                  </div>
                </div>
              </div>
            </template>
            <p v-else class="text-muted-foreground text-sm">No holidays this month</p>
          </div>
        </div>
      </div>

      <!-- Selected Date Info -->
      <div v-if="selectedDate" class="border rounded-2xl mt-6">
        <div class="border-b p-4">
          <h3 class="text-lg font-semibold">{{ formatDateIndonesian(selectedDate) }}</h3>
        </div>
        <div class="p-4">
          <div class="space-y-4">
            <template v-if="getHolidayInfo(holidaysMap, currentYear, currentMonth, selectedDate.getDate())">
              <div class="flex items-center gap-2">
                <span
                  class="inline-flex items-center rounded-md bg-destructive/10 text-destructive px-2 py-1 text-xs font-semibold">Holiday</span>
                <span class="text-sm font-medium">
                  {{ getHolidayInfo(holidaysMap, currentYear, currentMonth, selectedDate.getDate())!.summary }}
                </span>
              </div>
            </template>
            <p v-else class="text-muted-foreground">No holiday on this date</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped></style>
