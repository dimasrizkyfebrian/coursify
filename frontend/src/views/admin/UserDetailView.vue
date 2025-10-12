<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Skeleton } from '@/components/ui/skeleton'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'

const route = useRoute()
const router = useRouter()
const user = ref<any>(null)
const isLoading = ref(true)
const isSubmitting = ref(false)

async function fetchUserDetails() {
  const userId = route.params.id
  try {
    isLoading.value = true
    const response = await api.get(`/admin/users/${userId}`)
    user.value = response.data
  } catch (error) {
    console.error('Failed to fetch user details:', error)
    toast.error('Failed to load user details.')
  } finally {
    isLoading.value = false
  }
}

async function handleApprove() {
  isSubmitting.value = true
  try {
    await api.put(`/admin/users/${user.value.id}/approve`)
    toast.success('User has been approved.')
    await fetchUserDetails()
  } catch (error) {
    toast.error('Failed to approve user.')
  } finally {
    isSubmitting.value = false
  }
}

async function handleReject() {
  isSubmitting.value = true
  try {
    await api.put(`/admin/users/${user.value.id}/reject`)
    toast.warning('User has been rejected.')
    await fetchUserDetails()
  } catch (error) {
    toast.error('Failed to reject user.')
  } finally {
    isSubmitting.value = false
  }
}

async function handleDelete() {
  if (!user.value) return
  isSubmitting.value = true
  try {
    await api.delete(`/admin/users/${user.value.id}`)
    toast.success('User has been deleted.')
    router.push({ name: 'admin-all-users' })
  } catch (error) {
    toast.error('Failed to delete user.')
  } finally {
    isSubmitting.value = false
  }
}

const userInitials = computed(() => {
  return user.value?.full_name?.charAt(0).toUpperCase() || 'U'
})

onMounted(() => {
  fetchUserDetails()
})
</script>

<template>
  <div>
    <Card v-if="isLoading">
      <CardHeader class="flex flex-row items-center gap-4">
        <Skeleton class="w-16 h-16 rounded-full" />
        <div class="space-y-2">
          <Skeleton class="h-6 w-48" />
          <Skeleton class="h-4 w-64" />
        </div>
      </CardHeader>
      <CardContent>
        <Skeleton class="w-full h-24" />
      </CardContent>
    </Card>

    <Card v-if="!isLoading && user">
      <CardHeader>
        <div class="flex items-start justify-between gap-4">
          <div class="flex items-center gap-4">
            <Avatar class="w-16 h-16 text-xl">
              <AvatarImage :src="user.avatar_url" :alt="user.full_name" />
              <AvatarFallback>{{ userInitials }}</AvatarFallback>
            </Avatar>
            <div>
              <CardTitle class="text-2xl">{{ user.full_name }}</CardTitle>
              <CardDescription>{{ user.email }}</CardDescription>
            </div>
          </div>
          <div v-if="user.status === 'pending'" class="flex gap-2">
            <Button @click="handleReject" variant="outline" :disabled="isSubmitting">Reject</Button>
            <Button @click="handleApprove" :disabled="isSubmitting">Approve</Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <Separator class="my-4" />
        <div class="grid grid-cols-2 gap-4 text-sm">
          <div>
            <p class="text-muted-foreground">Role</p>
            <Badge variant="secondary">{{ user.role }}</Badge>
          </div>
          <div>
            <p class="text-muted-foreground">Status</p>
            <Badge
              :variant="
                user.status === 'active'
                  ? 'default'
                  : user.status === 'pending'
                    ? 'outline'
                    : 'destructive'
              "
            >
              {{ user.status }}
            </Badge>
          </div>
          <div>
            <p class="text-muted-foreground">Joined On</p>
            <p>{{ new Date(user.created_at).toLocaleDateString() }}</p>
          </div>
        </div>
      </CardContent>
      <CardFooter class="justify-end border-t pt-4">
        <AlertDialog>
          <AlertDialogTrigger as-child>
            <Button variant="destructive">Delete User</Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
              <AlertDialogDescription>
                This action cannot be undone. This will permanently delete the user account.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction @click="handleDelete" :disabled="isSubmitting"
                >Continue</AlertDialogAction
              >
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </CardFooter>
    </Card>
  </div>
</template>
