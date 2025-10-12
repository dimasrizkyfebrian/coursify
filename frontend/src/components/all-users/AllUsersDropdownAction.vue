<script setup lang="ts">
import type { User } from './columns'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Ellipsis, Pencil, Trash2 } from 'lucide-vue-next'

defineProps<{
  user: User
}>()

function editUser() {
  toast.info('Edit button clicked.')
}
function deleteUser() {
  toast.success('Delete button clicked.')
}
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" class="w-8 h-8 p-0">
        <Ellipsis class="w-4 h-4" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end">
      <DropdownMenuLabel>Actions</DropdownMenuLabel>
      <DropdownMenuItem @click="editUser">
        <Pencil class="h-4 w-4 mr-2" />
        Edit User
      </DropdownMenuItem>
      <DropdownMenuItem @click="deleteUser" class="text-red-600 focus:text-red-700">
        <Trash2 class="h-4 w-4 mr-2" />
        Delete User
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
