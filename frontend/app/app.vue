<script setup lang="ts">
import { useAuth } from '~/composables/useAuth'

const { isLoggedIn } = useAuth()

const layout = computed(() => (isLoggedIn.value ? 'sidebar' : 'default'))
</script>

<template>
  <NuxtLayout :name="layout">
    <NuxtPage />
  </NuxtLayout>
</template>
