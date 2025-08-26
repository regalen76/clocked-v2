<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAttendance } from '~/composables/useAttendance'
import { useTasks } from '~/composables/useTasks'

const { list, refresh, clockIn, clockOut } = useAttendance()
const { tasks, fetchTasks, createTask } = useTasks()

const ci = ref<string>('')
const co = ref<string>('')
const mode = ref<'select' | 'create'>('select')
const selectedTaskId = ref<number | null>(null)
const newTask = ref({ name: '', subject: '', description: '', how_its_done: '' })
const markCompleted = ref(true)

async function doClockIn() {
  const when = ci.value ? new Date(ci.value) : new Date()
  await clockIn(when)
  ci.value = ''
  await refresh()
}

async function doClockOut() {
  const when = co.value ? new Date(co.value) : new Date()
  if (mode.value === 'create') {
    if (!newTask.value.name.trim()) return
    await clockOut({
      clock_out: when,
      new_task: {
        name: newTask.value.name,
        subject: newTask.value.subject || undefined,
        description: newTask.value.description || undefined,
        how_its_done: newTask.value.how_its_done || undefined,
      },
      mark_completed: markCompleted.value,
    })
  } else {
    await clockOut({ clock_out: when, task_id: selectedTaskId.value ?? undefined, mark_completed: markCompleted.value })
  }
  co.value = ''
  newTask.value = { name: '', subject: '', description: '', how_its_done: '' }
  selectedTaskId.value = null
  await Promise.all([refresh(), fetchTasks()])
}

const openRecord = computed(() => list.value.find(a => !a.clock_out))

onMounted(async () => {
  await Promise.all([refresh(), fetchTasks()])
})
</script>

<template>
  <div class="container mx-auto p-6 space-y-6">
    <div>
      <h1 class="text-2xl font-semibold">Attendances</h1>
      <p class="text-muted-foreground text-sm">Clock in/out and associate tasks</p>
    </div>

    <div class="grid gap-4 md:grid-cols-2">
      <div class="border rounded-xl p-4 space-y-3">
        <h2 class="font-medium">Clock In</h2>
        <input v-model="ci" type="datetime-local" class="border rounded p-2 w-full" />
        <button class="border px-3 py-2 rounded" @click="doClockIn">Clock In</button>
        <div class="text-xs text-muted-foreground">If empty, current time is used.</div>
        <div v-if="openRecord" class="text-sm text-amber-600">Currently clocked in since {{ new Date(openRecord.clock_in).toLocaleString() }}</div>
      </div>

      <div class="border rounded-xl p-4 space-y-3">
        <h2 class="font-medium">Clock Out</h2>
        <input v-model="co" type="datetime-local" class="border rounded p-2 w-full" />
        <div class="flex items-center gap-3">
          <label class="text-sm font-medium">Task</label>
          <select v-model.number="selectedTaskId" class="border rounded p-2 flex-1" :disabled="mode === 'create'">
            <option :value="null">— Select task —</option>
            <option v-for="t in tasks" :key="t.ID" :value="t.ID">{{ t.name }}</option>
          </select>
          <button class="text-xs underline" @click="mode = mode==='select' ? 'create' : 'select'">{{ mode==='select' ? 'Create new' : 'Choose existing' }}</button>
        </div>
        <div v-if="mode==='create'" class="grid gap-2 border rounded p-3">
          <input v-model="newTask.name" class="border rounded p-2" placeholder="Task name*" />
          <input v-model="newTask.subject" class="border rounded p-2" placeholder="Subject" />
          <textarea v-model="newTask.description" class="border rounded p-2" placeholder="Description" />
          <textarea v-model="newTask.how_its_done" class="border rounded p-2" placeholder="How it's done" />
        </div>
        <label class="flex items-center gap-2 text-sm"><input type="checkbox" v-model="markCompleted" /> Mark selected/new task as completed</label>
        <button class="border px-3 py-2 rounded" @click="doClockOut">Clock Out</button>
        <div class="text-xs text-muted-foreground">If empty, current time is used.</div>
      </div>
    </div>

    <div class="border rounded-xl p-4">
      <h2 class="font-medium mb-2">History</h2>
      <div v-if="list.length===0" class="text-muted-foreground text-sm">No attendance records</div>
      <div v-else class="grid gap-2">
        <div v-for="a in list" :key="a.ID" class="flex items-center justify-between border rounded p-2 text-sm">
          <div>
            <span class="font-medium">{{ new Date(a.clock_in).toLocaleString() }}</span>
            <span v-if="a.clock_out"> → {{ new Date(a.clock_out).toLocaleString() }}</span>
          </div>
          <div class="text-xs text-muted-foreground">Task: {{ a.task_id ?? '—' }}</div>
        </div>
      </div>
    </div>
  </div>
</template>
