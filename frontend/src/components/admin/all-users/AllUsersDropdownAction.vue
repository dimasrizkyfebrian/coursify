<script setup lang="ts">
import { ref, inject, watch } from 'vue'
import type { User } from './columns'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Ellipsis, Pencil, Trash2 } from 'lucide-vue-next'

const props = defineProps<{ user: User }>()

const isEditDialogOpen = ref(false)
const isDeleteDialogOpen = ref(false)

const formData = ref({
  full_name: '',
  email: '',
  role: '',
})

const refreshData = inject<() => void>('refreshUsers')

watch(isEditDialogOpen, (newValue) => {
  if (newValue) {
    formData.value = {
      full_name: props.user.full_name,
      email: props.user.email,
      role: props.user.role,
    }
  }
})

async function handleUpdate() {
  try {
    const response = await api.put(`/admin/users/${props.user.id}`, formData.value)
    toast.success(response.data.message)
    isEditDialogOpen.value = false
    if (refreshData) refreshData()
  } catch (error) {
    toast.error('Failed to update user.')
  }
}

async function handleDeleteConfirm() {
  try {
    const response = await api.delete(`/admin/users/${props.user.id}`)
    toast.success(response.data.message || 'User deleted successfully.')
    if (refreshData) refreshData()
  } catch (error) {
    toast.error('Failed to delete user.')
  } finally {
    isDeleteDialogOpen.value = false
  }
}
</script>

<template>
  <div>
    <DropdownMenu>
      <DropdownMenuTrigger as-child>
        <Button variant="ghost" class="w-8 h-8 p-0">
          <Ellipsis class="w-4 h-4" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuLabel>Actions</DropdownMenuLabel>
        <DropdownMenuItem @click="isEditDialogOpen = true">
          <Pencil class="w-4 h-4 mr-2" /> Edit
        </DropdownMenuItem>
        <DropdownMenuItem
          @click="isDeleteDialogOpen = true"
          class="text-red-600 focus:text-red-700"
        >
          <Trash2 class="w-4 h-4 mr-2" /> Delete
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>

    <Dialog v-model:open="isEditDialogOpen">
      <DialogContent class="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Edit User</DialogTitle>
          <DialogDescription>
            Make changes to the user's profile here. Click save when you're done.
          </DialogDescription>
        </DialogHeader>
        <div class="grid gap-4 py-4">
          <div class="grid grid-cols-4 items-center gap-4">
            <Label for="name" class="text-right">Full Name</Label>
            <Input id="name" v-model="formData.full_name" class="col-span-3" />
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label for="email" class="text-right">Email</Label>
            <Input id="email" v-model="formData.email" class="col-span-3" />
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label for="role" class="text-right">Role</Label>
            <Select v-model="formData.role">
              <SelectTrigger class="col-span-3">
                <SelectValue placeholder="Select a role" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value="admin">Admin</SelectItem>
                  <SelectItem value="instructor">Instructor</SelectItem>
                  <SelectItem value="student">Student</SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
        </div>
        <DialogFooter>
          <Button @click="handleUpdate" type="submit">Save changes</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <AlertDialog v-model:open="isDeleteDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This action cannot be undone. This will permanently delete the account for
            <strong>{{ user.full_name }}</strong
            >.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction @click="handleDeleteConfirm">Continue</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
