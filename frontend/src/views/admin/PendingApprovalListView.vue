<script setup lang="ts">
import { ref, onMounted, provide } from 'vue'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'
import type { User } from '../../components/admin/pending-users/columns.ts'
import { columns } from '../../components/admin/pending-users/columns.ts'
import DataTable from '../../components/data-table/data-table.vue'

const users = ref<User[]>([])

async function fetchPendingUsers() {
  try {
    const response = await api.get('/admin/users/pending')
    users.value = response.data || []
  } catch (error) {
    toast.error('Failed to fetch users.')
    console.error(error)
  }
}

onMounted(() => {
  fetchPendingUsers()
})

provide('refreshUsers', fetchPendingUsers)
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-4">User Management</h1>
    <p class="text-gray-500 mb-6">List of users waiting for approval.</p>

    <DataTable :columns="columns" :data="users" />
  </div>
</template>
