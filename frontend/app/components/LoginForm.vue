<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'

import { Button } from '~/components/ui/button'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form'
import { Input } from '~/components/ui/input'
import { loginSchema, type Login } from '~/types/auth'
import { useAuth } from '~/composables/useAuth'

const { login } = useAuth()

const formSchema = toTypedSchema(loginSchema)

const { handleSubmit } = useForm<Login>({
  validationSchema: formSchema,
})

const onSubmit = handleSubmit(async (values) => {
  await login(values)
})
</script>

<template>
  <form class="flex flex-col gap-6" @submit="onSubmit">
    <div class="flex flex-col items-center gap-2 text-center">
      <h1 class="text-2xl font-bold">
        Login to your account
      </h1>
      <p class="text-muted-foreground text-sm text-balance">
        Enter your email below to login to your account
      </p>
    </div>
    <div class="grid gap-6">
      <FormField v-slot="{ componentField }" name="email">
        <FormItem>
          <FormLabel>Email</FormLabel>
          <FormControl>
            <Input type="email" placeholder="m@example.com" v-bind="componentField" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>
      <FormField v-slot="{ componentField }" name="password">
        <FormItem>
          <div class="flex items-center">
            <FormLabel>Password</FormLabel>
            <a href="#" class="ml-auto text-sm underline-offset-4 hover:underline">
              Forgot your password?
            </a>
          </div>
          <FormControl>
            <Input type="password" v-bind="componentField" />
          </FormControl>
          <FormMessage />
        </FormItem>
      </FormField>
      <Button type="submit" class-name="w-full">
        Login
      </Button>
    </div>
    <div class="text-center text-sm">
      Don't have an account?
      <a href="#" class="underline underline-offset-4">
        Sign up
      </a>
    </div>
  </form>
</template>