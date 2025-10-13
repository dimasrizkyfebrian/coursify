<script setup lang="ts">
import { ref, onMounted } from 'vue'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'

// Import komponen UI
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { PlusCircle } from 'lucide-vue-next'

const courses = ref<any[]>([])
const isLoading = ref(true)

async function fetchMyCourses() {
  try {
    isLoading.value = true
    const response = await api.get('/instructor/courses')
    courses.value = response.data || []
  } catch (error) {
    toast.error('Failed to fetch your courses.')
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchMyCourses()
})

function openCreateCourseModal() {
  toast.info('Create New Course button clicked.')
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold mb-4">My Courses</h1>
        <p class="text-gray-500">Manage all courses you have created.</p>
      </div>
      <Button @click="openCreateCourseModal">
        <PlusCircle class="w-4 h-4 mr-2" />
        Create New Course
      </Button>
    </div>

    <div v-if="isLoading" class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      <Skeleton v-for="i in 3" :key="i" class="h-48 rounded-lg" />
    </div>

    <div v-else>
      <div v-if="courses.length === 0" class="text-center py-16 border-2 border-dashed rounded-lg">
        <h3 class="text-xl font-semibold">No Courses Found</h3>
        <p class="text-muted-foreground mt-2">
          You haven't created any courses yet. Get started now!
        </p>
        <Button @click="openCreateCourseModal" class="mt-4"> Create Your First Course </Button>
      </div>

      <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <Card v-for="course in courses" :key="course.id">
          <CardHeader>
            <CardTitle>{{ course.title }}</CardTitle>
          </CardHeader>
          <CardContent>
            <CardDescription>{{ course.description }}</CardDescription>
          </CardContent>
        </Card>
      </div>
    </div>
  </div>
</template>
