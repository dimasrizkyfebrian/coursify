<script setup lang="ts">
import { computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { RouterLink } from 'vue-router'

import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar'

import { Users, Settings, BookCopy, LayoutDashboard } from 'lucide-vue-next'

const { user } = useUserStore()

const adminLinks = [
  { title: 'User Management', to: '/admin/users', icon: Users },
  { title: 'Settings', to: '/admin/settings', icon: Settings },
]

const instructorLinks = [{ title: 'My Courses', to: '/instructor/courses', icon: BookCopy }]

const studentLinks = [{ title: 'Enrolled Courses', to: '/student/courses', icon: LayoutDashboard }]

const navLinks = computed(() => {
  switch (user.value.role) {
    case 'admin':
      return adminLinks
    case 'instructor':
      return instructorLinks
    case 'student':
      return studentLinks
    default:
      return []
  }
})
</script>

<template>
  <Sidebar>
    <SidebarHeader>
      <h1 class="text-xl font-bold text-center">Coursify</h1>
    </SidebarHeader>

    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupLabel class="uppercase"> {{ user.role }} Menu </SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="link in navLinks" :key="link.title">
              <SidebarMenuButton as-child>
                <RouterLink :to="link.to">
                  <component :is="link.icon" class="w-5 h-5 mr-2" />
                  <span>{{ link.title }}</span>
                </RouterLink>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>
  </Sidebar>
</template>
