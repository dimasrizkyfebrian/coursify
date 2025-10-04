<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'

const userProfile = ref<any>(null)
const router = useRouter()

// Fetch user profile on component mount
onMounted(async () => {
  try {
    const response = await api.get('/profile')
    userProfile.value = response.data
  } catch (error) {
    console.error('Failed to fetch profile:', error)
    toast.error('Authentication failed', {
      description: 'Please login again.',
    })
    router.push('/login')
  }
})

// Function to handle logout
function handleLogout() {
  localStorage.removeItem('authToken')
  toast.success('You have been logged out.')
  router.push('/login')
}
</script>

<template>
  <div class="p-8">
    <div v-if="userProfile" class="max-w-4xl mx-auto">
      <div class="flex justify-between items-center mb-8">
        <h1 class="text-3xl font-bold">Welcome to your Dashboard!</h1>
        <Button @click="handleLogout" variant="destructive"> Logout </Button>
      </div>
      <div class="bg-gray-100 dark:bg-gray-800 p-4 rounded-lg">
        <h2 class="text-xl mb-2">Your Profile Information:</h2>
        <pre>{{ userProfile }}</pre>
        <p class="mt-4 text-sm text-gray-500">Your User ID is: {{ userProfile.user_id }}</p>
      </div>
    </div>

    <div v-else class="max-w-4xl mx-auto">
      <div class="flex justify-between items-center mb-8">
        <Skeleton class="h-8 w-72 rounded-md" /> <Skeleton class="h-10 w-24 rounded-md" />
      </div>
      <div class="bg-gray-100 dark:bg-gray-800 p-4 rounded-lg space-y-4">
        <Skeleton class="h-6 w-52 rounded-md" /> <Skeleton class="h-4 w-full rounded-md" />
        <Skeleton class="h-4 w-[90%] rounded-md" />
        <Skeleton class="h-4 w-[70%] rounded-md" />
      </div>
    </div>
  </div>
</template>
