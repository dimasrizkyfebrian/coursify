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
      <h1 class="text-3xl font-bold">Instructor Dashboard</h1>
      <Button @click="handleLogout" variant="destructive"> Logout </Button>
    </div>
    <p class="mt-4">Welcome, Instructor! Here you can manage your courses and students.</p>
  </div>
</template>
