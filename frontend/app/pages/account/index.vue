<script setup lang="ts">
import { useUserStore } from '~/composables/useUserStore';

const { fetchUser } = useUser()
const fileInput = ref<HTMLInputElement | null>(null)
const previewUrl = ref<string | null>(null)
const fileRef = ref<File | null>(null)
const userStore = useUserStore();

function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files?.[0]) {
    const file = target.files[0]
    fileRef.value = file
    const reader = new FileReader()
    reader.onload = (e) => {
      previewUrl.value = e.target?.result as string
    }
    reader.readAsDataURL(file)
  }
}

async function handleAvatarUpload() {
  if (!fileRef.value) return

  const formData = new FormData()
  formData.append('avatar', fileRef.value)

  const { error } = await useFetch(
    `http://localhost:8000/api/user/${userStore.user?.ID}/avatar`,
    {
      method: 'POST',
      body: formData,
      headers: {
        Authorization: `Bearer ${useAuth().token.value}`,
      },
    },
  )

  if (error.value) {
    console.error(error.value)
    return
  }

  await fetchUser()
  previewUrl.value = null
}
</script>

<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold">Account</h1>
    <div v-if="userStore.user" class="mt-4 grid gap-6">
      <div class="flex flex-col items-center gap-4">
        <Label for="avatar" class="relative group cursor-pointer">
          <Avatar class="size-44 text-6xl">
            <AvatarImage :src="previewUrl ?? userStore.avatarUrl" :alt="userStore.user.username" />
            <AvatarFallback>{{ userStore.user.username.charAt(0).toUpperCase() }}</AvatarFallback>
          </Avatar>
          <div
            class="absolute inset-0 flex items-center justify-center bg-black bg-opacity-50 rounded-full opacity-0 group-hover:opacity-80 transition-opacity pointer-events-none">
            <span class="text-white text-sm">Change</span>
          </div>
        </Label>
        <Input id="avatar" type="file" ref="fileInput" class="hidden" @change="handleFileChange" />
      </div>

      <div class="flex items-center gap-3">
        <Label for="username" class="w-18 font-medium">Username:</Label>
        <p>{{ userStore.user.username }}</p>
      </div>

      <div class="flex items-center gap-3">
        <Label for="names" class="w-18 font-medium">Name:</Label>
        <p>{{ userStore.user.names }}</p>
      </div>

      <div class="flex items-center gap-3">
        <Label for="email" class="w-18 font-medium">Email:</Label>
        <p>{{ userStore.user.email }}</p>
      </div>

      <Button @click="handleAvatarUpload" :disabled="!previewUrl">Save Avatar</Button>
    </div>
  </div>
</template>
