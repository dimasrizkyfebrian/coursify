<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'
import type { User } from '../../components/pending-users/columns.ts'
import { columns } from '../../components/pending-users/columns.ts'
import DataTable from '../../components/pending-users/data-table.vue'

const users = ref<User[]>([])

async function fetchPendingUsers() {
  try {
    const response = await api.get('/admin/users')
    users.value = response.data || []
  } catch (error) {
    toast.error('Failed to fetch users.')
    console.error(error)
  }
}

onMounted(() => {
  fetchPendingUsers()
})
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-4">User Management</h1>
    <p class="text-gray-500 mb-6">List of users waiting for approval.</p>

    <DataTable :columns="columns" :data="users" />
  </div>
</template>
