<script setup lang="ts">
import { inject } from 'vue'
import { defineProps } from 'vue'
import { RouterLink } from 'vue-router'
import type { User } from './columns.ts'
import { useUserStore } from '@/stores/user'
import api from '@/lib/axios'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Ellipsis, CircleCheck, CircleX } from 'lucide-vue-next'

defineProps<{
  user: User
}>()

const { fetchPendingUserCount } = useUserStore()
const refreshUsers = inject<() => void>('refreshUsers')

async function approveUser(userId: string) {
  try {
    const response = await api.put(`/admin/users/${userId}/approve`)
    toast.success(response.data.message || 'User approved successfully.')
    fetchPendingUserCount()
    if (refreshUsers) refreshUsers()
  } catch (error) {
    toast.error('Failed to approve user.')
  }
}

async function rejectUser(userId: string) {
  try {
    const response = await api.put(`/admin/users/${userId}/reject`)
    toast.success(response.data.message || 'User rejected successfully.')
    fetchPendingUserCount()
    if (refreshUsers) refreshUsers()
  } catch (error) {
    toast.error('Failed to reject user.')
  }
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" class="w-8 h-8 p-0">
        <span class="sr-only">Open menu</span>
        <Ellipsis class="w-4 h-4" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end">
      <DropdownMenuLabel>Actions</DropdownMenuLabel>
      <DropdownMenuItem @click="approveUser(user.id)" class="text-green-600 focus:text-green-700">
        <CircleCheck class="w-4 h-4 mr-2" /> Approve User
      </DropdownMenuItem>
      <DropdownMenuItem @click="rejectUser(user.id)" class="text-red-600 focus:text-red-700">
        <CircleX class="w-4 h-4 mr-2" /> Reject User
      </DropdownMenuItem>
      <DropdownMenuSeparator />
      <DropdownMenuItem as-child>
        <RouterLink :to="{ name: 'admin-pending-detail', params: { id: user.id } }">
          View details
        </RouterLink>
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
