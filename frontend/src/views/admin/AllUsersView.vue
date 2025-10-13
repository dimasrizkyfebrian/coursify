<script setup lang="ts">
import { ref, onMounted, provide, computed } from 'vue'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'

import DataTable from '@/components/data-table/data-table.vue'
import type { User } from '@/components/admin/all-users/columns'
import { columns } from '@/components/admin/all-users/columns'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

const users = ref<User[]>([])
const isLoading = ref(true)

const dataTable = ref<any>(null)
const table = computed(() => dataTable.value?.table)

const statuses = [
  { value: 'active', label: 'Active' },
  { value: 'pending', label: 'Pending' },
  { value: 'rejected', label: 'Rejected' },
]

function setStatusFilter(statusValue: string) {
  const currentFilter = table.value.getColumn('status')?.getFilterValue()

  if (currentFilter === statusValue) {
    table.value.getColumn('status')?.setFilterValue(undefined)
  } else {
    table.value.getColumn('status')?.setFilterValue(statusValue)
  }
}

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

provide('refreshUsers', fetchAllUsers)
</script>

<template>
  <div>
    <h1 class="text-2xl font-bold mb-4">All Registered Users</h1>
    <p class="text-gray-500">Manage all users in the platform.</p>

    <div class="flex items-center justify-end">
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button variant="outline"> Filter by Status </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent class="w-56">
          <DropdownMenuLabel>Status</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuCheckboxItem
            @select.prevent="setStatusFilter('')"
            :checked="!table?.getColumn('status')?.getFilterValue()"
          >
            All
          </DropdownMenuCheckboxItem>
          <DropdownMenuCheckboxItem
            v-for="status in statuses"
            :key="status.value"
            :checked="table?.getColumn('status')?.getFilterValue() === status.value"
            @select.prevent="setStatusFilter(status.value)"
          >
            {{ status.label }}
          </DropdownMenuCheckboxItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>

    <DataTable ref="dataTable" :columns="columns" :data="users" />
  </div>
</template>
