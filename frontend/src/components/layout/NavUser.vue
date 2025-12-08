<script setup lang="ts">
import { useUserStore } from '@/stores/user'
import { useRouter } from 'vue-router'

import { ChevronsUpDown, LogOut } from 'lucide-vue-next'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar'

const { isMobile } = useSidebar()
const { user, logoutUser } = useUserStore()
const router = useRouter()

// Default avatar URL
const defaultAvatar = '/images/avatars/default.webp'

// Helper to get avatar URL or fallback
function getUserAvatarUrl(avatarPath: string | null) {
  if (avatarPath) {
    const backendBaseUrl =
      import.meta.env.VITE_API_BASE_URL?.replace(/\/api\/v1$/, '') || 'http://localhost:8080'
    return `${backendBaseUrl}${avatarPath}`
  }
  return defaultAvatar
}

// Helper to get initials
function getInitials(name: string): string {
  if (!name || name.trim().length === 0) {
    return '?'
  }
  const names = name.trim().split(' ')
  const firstInitial = names[0]?.charAt(0)?.toUpperCase()
  if (names.length === 1) {
    return firstInitial || '?'
  }
  const lastInitial = names[names.length - 1]?.charAt(0)?.toUpperCase()
  if (firstInitial && lastInitial) {
    return firstInitial + lastInitial
  }
  return firstInitial || '?'
}

function handleLogout() {
  logoutUser()
  router.push('/login')
}
</script>

<template>
  <SidebarMenu>
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <SidebarMenuButton
            size="lg"
            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          >
            <Avatar class="h-8 w-8 rounded-lg">
              <AvatarImage :src="getUserAvatarUrl(user.avatarUrl)" :alt="user.fullName" />
              <AvatarFallback class="rounded-lg">
                {{ getInitials(user.fullName) }}
              </AvatarFallback>
            </Avatar>
            <div class="grid flex-1 text-left text-sm leading-tight">
              <span class="truncate font-semibold">{{ user.fullName }}</span>
              <span class="truncate text-xs">{{ user.email }}</span>
            </div>
            <ChevronsUpDown class="ml-auto size-4" />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          class="w-[--reka-dropdown-menu-trigger-width] min-w-56 rounded-lg"
          :side="isMobile ? 'bottom' : 'right'"
          align="end"
          :side-offset="4"
        >
          <DropdownMenuLabel class="p-0 font-normal">
            <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
              <Avatar class="h-8 w-8 rounded-lg">
                <AvatarImage :src="getUserAvatarUrl(user.avatarUrl)" :alt="user.fullName" />
                <AvatarFallback class="rounded-lg">
                  {{ getInitials(user.fullName) }}
                </AvatarFallback>
              </Avatar>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">{{ user.fullName }}</span>
                <span class="truncate text-xs">{{ user.email }}</span>
              </div>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem @click="handleLogout">
            <LogOut />
            Log out
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenuItem>
  </SidebarMenu>
</template>
