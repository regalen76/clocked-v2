<script setup lang="ts">
import {
  Command,
  Frame,
  LifeBuoy,
  Map,
  PieChart,
  Send,
} from 'lucide-vue-next'

import NavProjects from '~/components/NavProjects.vue'
import NavSecondary from '~/components/NavSecondary.vue'
import NavUser from '~/components/NavUser.vue'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  type SidebarProps,
} from '~/components/ui/sidebar'
import { useUser } from '~/composables/useUser'
import { useUserStore } from '~/composables/useUserStore'

const props = withDefaults(defineProps<SidebarProps>(), {
  variant: 'inset',
})

const { fetchUser } = useUser()
const userStore = useUserStore();

onMounted(() => {
  fetchUser()
})

const displayName = computed(() => {
  if (!userStore.user) return ''
  return userStore.user.names || userStore.user.username
})

const userEmail = computed(() => {
  if (!userStore.user) return ''
  return userStore.user.email
})

const data = {
  navSecondary: [
    {
      title: 'Support',
      url: '#',
      icon: LifeBuoy,
    },
    {
      title: 'Feedback',
      url: '#',
      icon: Send,
    },
  ],
  projects: [
    {
      name: 'Attendances',
      url: '/attendances',
      icon: Frame,
    },
    {
      name: 'Tasks',
      url: '/tasks',
      icon: PieChart,
    },
    {
      name: 'Timesheet',
      url: '/timesheet',
      icon: Map,
    },
  ],
}
</script>

<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child>
            <NuxtLink to="/">
              <div
                class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <Command class="size-4" />
              </div>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-medium">Reonify</span>
                <span class="truncate text-xs">Code</span>
              </div>
            </NuxtLink>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>
    <SidebarContent>
      <NavProjects :projects="data.projects" />
      <NavSecondary :items="data.navSecondary" class="mt-auto" />
    </SidebarContent>
    <SidebarFooter>
      <NavUser
        :user="{ name: displayName, email: userEmail, avatar: userStore.user?.avatar ? userStore.avatarUrl : '/avatars/shadcn.jpg' }" />
    </SidebarFooter>
  </Sidebar>
</template>
