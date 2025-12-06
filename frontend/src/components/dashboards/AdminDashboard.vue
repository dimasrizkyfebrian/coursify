<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'
import { RouterLink } from 'vue-router'

import { Card, CardContent, CardHeader, CardFooter, CardTitle } from '@/components/ui/card'
import { Users, UserCheck, UserPlus } from 'lucide-vue-next'
import { Skeleton } from '@/components/ui/skeleton'

const stats = ref({
  total_users: 0,
  active_users: 0,
  pending_users: 0,
})
const isLoading = ref(true)

async function fetchStats() {
  try {
    isLoading.value = true
    const response = await api.get('/admin/users/stats')
    stats.value = response.data
  } catch (error) {
    toast.error('Failed to load dashboard statistics.')
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold mb-6">Admin Dashboard</h1>

    <div v-if="isLoading" class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      <Skeleton class="h-32" />
      <Skeleton class="h-32" />
      <Skeleton class="h-32" />
    </div>

    <div v-else class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium"> Total Users </CardTitle>
          <Users class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ stats.total_users }}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium"> Active Users </CardTitle>
          <UserCheck class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ stats.active_users }}
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium"> Pending Approval </CardTitle>
          <UserPlus class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ stats.pending_users }}
          </div>
          <p v-if="stats.pending_users > 0" class="text-xs text-muted-foreground">
            <RouterLink :to="{ name: 'admin-pending-list' }" class="underline">
              View pending users
            </RouterLink>
          </p>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
