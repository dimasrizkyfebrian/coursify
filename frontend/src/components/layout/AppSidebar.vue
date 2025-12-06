<script setup lang="ts">
import { computed } from 'vue'
import { useUserStore } from '@/stores/user'
import type { SidebarProps } from '@/components/ui/sidebar'
import NavMain from '@/components/layout/NavMain.vue'
import NavUser from '@/components/layout/NavUser.vue'

import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarFooter,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarRail,
} from '@/components/ui/sidebar'

import { CpuIcon, Users, Settings2, BookCopy, LayoutDashboard } from 'lucide-vue-next'

const { user } = useUserStore()
const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: 'icon',
  variant: 'floating',
})

const adminLinks = [
  {
    title: 'User Management',
    icon: Users,
    isActive: true,
    children: [
      { title: 'Pending Approval', name: 'admin-pending-list' },
      { title: 'All Users', name: 'admin-all-users' },
    ],
  },
  { title: 'Settings', name: 'admin-settings', icon: Settings2 },
]

const instructorLinks = [{ title: 'My Courses', name: 'instructor-my-courses', icon: BookCopy }]

const studentLinks = [{ title: 'Enrolled Courses', to: '/student/courses', icon: LayoutDashboard }]

const navItems = computed(() => {
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
  <Sidebar v-bind="props">
    <SidebarHeader>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton size="lg" as-child>
            <a href="#">
              <div
                class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"
              >
                <CpuIcon class="size-4" />
              </div>
              <div class="flex flex-col gap-0.5 leading-none">
                <span class="font-semibold">Coursify</span>
                <span class="">v1.0.0</span>
              </div>
            </a>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarHeader>

    <SidebarContent class="flex-1">
      <NavMain :items="navItems" />
    </SidebarContent>

    <SidebarFooter>
      <NavUser />
    </SidebarFooter>
    <SidebarRail />
  </Sidebar>
</template>
