<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { Button } from '@/components/ui/button'

const router = useRouter()
const { logoutUser } = useUserStore()

function handleLogout() {
  logoutUser()
  router.push('/login')
}
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h1 class="text-3xl font-bold">Admin Dashboard</h1>
    </div>
    <p class="mt-4">Welcome, Admin! Here you can manage users and platform settings.</p>
  </div>
</template>
