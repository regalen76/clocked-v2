import type { Login, LoginResponse } from '~/types/auth'

export const useAuth = () => {
  const token = useCookie<string | null>('token')

  const isLoggedIn = computed(() => !!token.value)

  async function login(values: Login) {
    const { data, error } = await useFetch<LoginResponse>(
      'http://localhost:8000/api/auth/login',
      {
        method: 'POST',
        body: values,
      },
    )

    if (error.value) {
      console.error(error.value)
      throw new Error('Login failed')
    }

    if (data.value) {
      token.value = data.value.data
    }

    await navigateTo('/')
  }

  function logout() {
    token.value = null
    navigateTo('/login')
  }

  return { token, isLoggedIn, login, logout }
}
