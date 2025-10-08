<script setup lang="ts">
import type { User } from './columns.ts'
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

function approveUser(userId: string) {
  toast.success(`Approving user: ${userId}`)
  // Nanti di sini kita panggil API backend
}

// 4. Buat fungsi baru untuk rejectUser
function rejectUser(userId: string) {
  toast.error(`Rejecting user: ${userId}`)
  // Nanti di sini kita panggil API backend
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
      <DropdownMenuItem>View details</DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
