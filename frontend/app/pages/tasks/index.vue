<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useTasks, type Task } from '~/composables/useTasks'

const { tasks, fetchTasks, createTask, markCompleted } = useTasks()

const form = ref({
  name: '',
  subject: '',
  description: '',
  how_its_done: '',
})
const creating = ref(false)

async function submit() {
  if (!form.value.name.trim()) return
  creating.value = true
  try {
    await createTask({
      name: form.value.name,
      subject: form.value.subject || undefined,
      description: form.value.description || undefined,
      how_its_done: form.value.how_its_done || undefined,
    })
    form.value = { name: '', subject: '', description: '', how_its_done: '' }
  } finally {
    creating.value = false
  }
}

function label(task: Task) {
  return task.completed ? `Completed${task.completed_at ? ` (${new Date(task.completed_at).toLocaleDateString()})` : ''}` : 'Open'
}

onMounted(fetchTasks)
</script>

<template>
  <div class="container mx-auto p-6 space-y-6">
    <div>
      <h1 class="text-2xl font-semibold">Tasks</h1>
      <p class="text-muted-foreground text-sm">Create and track your tasks</p>
    </div>

    <form class="grid gap-3 max-w-xl border rounded-xl p-4" @submit.prevent="submit">
      <div class="grid gap-1">
        <label class="text-sm font-medium">Name<span class="text-red-500">*</span></label>
        <input v-model="form.name" required class="border rounded p-2" placeholder="Task name" />
      </div>
      <div class="grid gap-1">
        <label class="text-sm font-medium">Subject</label>
        <input v-model="form.subject" class="border rounded p-2" placeholder="Subject" />
      </div>
      <div class="grid gap-1">
        <label class="text-sm font-medium">Description</label>
        <textarea v-model="form.description" class="border rounded p-2" placeholder="Describe the task" />
      </div>
      <div class="grid gap-1">
        <label class="text-sm font-medium">How it's done (optional)</label>
        <textarea v-model="form.how_its_done" class="border rounded p-2" placeholder="Notes on how it's done" />
      </div>
      <div class="flex gap-2">
        <button :disabled="creating" class="border px-3 py-2 rounded">{{ creating ? 'Creating…' : 'Create Task' }}</button>
      </div>
    </form>

    <div class="grid gap-3">
      <h2 class="text-lg font-semibold">My Tasks</h2>
      <div v-if="tasks.length === 0" class="text-muted-foreground">No tasks yet</div>
      <div v-else class="grid gap-2">
        <div v-for="t in tasks" :key="t.ID" class="flex items-start justify-between border rounded-lg p-3">
          <div>
            <div class="font-medium">{{ t.name }}</div>
            <div class="text-xs text-muted-foreground" v-if="t.subject">{{ t.subject }}</div>
            <div class="text-sm" v-if="t.description">{{ t.description }}</div>
            <div class="text-xs" v-if="t.how_its_done">How: {{ t.how_its_done }}</div>
            <span :class="t.completed ? 'text-green-600' : 'text-amber-600'" class="text-xs font-medium">{{ label(t) }}</span>
          </div>
          <div class="flex items-center gap-2">
            <button v-if="!t.completed" class="border px-2 py-1 rounded text-xs" @click="markCompleted(t.ID, true)">Mark Completed</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
