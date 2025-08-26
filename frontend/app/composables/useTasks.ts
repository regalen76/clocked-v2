import { ref, watch } from 'vue'

export interface Task {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: null | string
  user_id: number
  name: string
  subject: string
  description: string
  how_its_done: string
  completed: boolean
  completed_at: string | null
}

export interface ApiResponse<T> {
  status: string
  message: string
  data: T
}

export function useTasks() {
  const { token } = useAuth()
  const apiBase = useRuntimeConfig().public.apiBase as string

  const tasks = ref<Task[]>([])
  const loading = ref(false)
  const error = ref<unknown>(null)

  const authHeaders = () =>
    token.value
      ? { Authorization: `Bearer ${token.value}` }
      : undefined

  async function fetchTasks() {
    loading.value = true
    error.value = null
    try {
      const res = await $fetch<ApiResponse<Task[]>>(`${apiBase}/tasks`, {
        headers: authHeaders(),
      })
      tasks.value = res.data
    } catch (e) {
      error.value = e
    } finally {
      loading.value = false
    }
  }

  async function fetchTasksByDay(dateISO: string) {
    loading.value = true
    error.value = null
    try {
      const res = await $fetch<ApiResponse<Task[]>>(`${apiBase}/tasks/day`, {
        headers: authHeaders(),
        query: { date: dateISO.slice(0, 10) },
      })
      return res.data
    } catch (e) {
      error.value = e
      return [] as Task[]
    } finally {
      loading.value = false
    }
  }

  async function createTask(payload: {
    name: string
    subject?: string
    description?: string
    how_its_done?: string
  }) {
    const res = await $fetch<ApiResponse<Task>>(`${apiBase}/tasks`, {
      method: 'POST',
      headers: authHeaders(),
      body: payload,
    })
    tasks.value.unshift(res.data)
    return res.data
  }

  async function markCompleted(id: number, completed = true, when?: Date) {
    const res = await $fetch<ApiResponse<Task>>(`${apiBase}/tasks/${id}`, {
      method: 'PATCH',
      headers: authHeaders(),
      body: {
        completed,
        completed_at: when ? when.toISOString() : undefined,
      },
    })
    const idx = tasks.value.findIndex(t => t.ID === id)
    if (idx !== -1) tasks.value[idx] = res.data
    return res.data
  }

  // keep tasks fresh when token appears/disappears
  watch(token, () => {
    tasks.value = []
  })

  return { tasks, loading, error, fetchTasks, fetchTasksByDay, createTask, markCompleted }
}
