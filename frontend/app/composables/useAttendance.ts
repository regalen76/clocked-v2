import { ref } from 'vue'

import type { ApiResponse } from '~/composables/useTasks'

export interface Attendance {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: null | string
  user_id: number
  task_id: number | null
  clock_in: string
  clock_out: string | null
}

export function useAttendance() {
  const { token } = useAuth()
  const apiBase = useRuntimeConfig().public.apiBase as string
  const list = ref<Attendance[]>([])
  const loading = ref(false)
  const error = ref<unknown>(null)

  const authHeaders = () =>
    token.value
      ? { Authorization: `Bearer ${token.value}` }
      : undefined

  async function refresh() {
    loading.value = true
    error.value = null
    try {
      const res = await $fetch<ApiResponse<Attendance[]>>(`${apiBase}/attendance`, {
        headers: authHeaders(),
      })
      list.value = res.data
    } catch (e) {
      error.value = e
    } finally {
      loading.value = false
    }
  }

  async function clockIn(clockIn: Date) {
    const res = await $fetch<ApiResponse<Attendance>>(`${apiBase}/attendance/clockin`, {
      method: 'POST',
      headers: authHeaders(),
      body: { clock_in: clockIn.toISOString() },
    })
    list.value.unshift(res.data)
    return res.data
  }

  async function clockOut(payload: {
    clock_out: Date
    task_id?: number
    new_task?: { name: string; subject?: string; description?: string; how_its_done?: string }
    mark_completed?: boolean
  }) {
    const res = await $fetch<ApiResponse<Attendance>>(`${apiBase}/attendance/clockout`, {
      method: 'POST',
      headers: authHeaders(),
      body: {
        clock_out: payload.clock_out.toISOString(),
        task_id: payload.task_id,
        new_task: payload.new_task,
        mark_completed: payload.mark_completed ?? false,
      },
    })
    // update the matching open record
    const idx = list.value.findIndex(a => a.ID === res.data.ID)
    if (idx !== -1) list.value[idx] = res.data
    else list.value.unshift(res.data)
    return res.data
  }

  return { list, loading, error, refresh, clockIn, clockOut }
}
