<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'

import DataTable from '@/components/pending-users/data-table.vue'
import type { User } from '@/components/all-users/columns'
import { columns } from '@/components/all-users/columns'

const users = ref<User[]>([])
const isLoading = ref(true)

async function fetchAllUsers() {
  try {
    isLoading.value = true
    const response = await api.get('/admin/users/all')
    users.value = response.data || []
  } catch (error) {
    toast.error('Failed to fetch users.')
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchAllUsers()
})
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-4">All Registered Users</h1>
    <p class="text-gray-500 mb-6">Manage all users in the platform.</p>

    <DataTable :columns="columns" :data="users" />
  </div>
</template>
