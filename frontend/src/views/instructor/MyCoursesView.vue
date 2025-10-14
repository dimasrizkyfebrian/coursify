<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterView, RouterLink, useRoute } from 'vue-router'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'

// Import components
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { PlusCircle, Pencil } from 'lucide-vue-next'
import CreateCourseDialog from '@/components/instructor/my-courses/CreateCourseDialog.vue'
import EditCourseDialog from '@/components/instructor/my-courses/EditCourseDialog.vue'

// State variables
const courses = ref<any[]>([])
const isLoading = ref(true)
const route = useRoute()
// Modal states
const isCreateModalOpen = ref(false)
const isEditModalOpen = ref(false)
const selectedCourse = ref<any | null>(null)

// Fetch my courses from the API
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
  isCreateModalOpen.value = true
}

function openEditCourseModal(course: any) {
  selectedCourse.value = course
  isEditModalOpen.value = true
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

    <template v-if="!route.params.id">
      <div v-if="isLoading" class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        <Skeleton v-for="i in 3" :key="i" class="h-48 rounded-lg" />
      </div>

      <div v-else>
        <div
          v-if="courses.length === 0"
          class="text-center py-16 border-2 border-dashed rounded-lg"
        >
          <h3 class="text-xl font-semibold">No Courses Found</h3>
          <p class="text-muted-foreground mt-2">
            You haven't created any courses yet. Get started now!
          </p>
          <Button @click="openCreateCourseModal" class="mt-4"> Create Your First Course </Button>
        </div>

        <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <RouterLink
            v-for="course in courses"
            :key="course.id"
            :to="{ name: 'instructor-course-detail', params: { id: course.id } }"
          >
            <Card class="hover:bg-accent transition-colors h-full overflow-hidden">
              <div class="mr-4 ml-4">
                <img
                  src="https://placehold.co/1080x720?text=Course+Image"
                  alt="Course Image Placeholder"
                  class="aspect-video w-full object-cover rounded md:rounded-lg"
                />
              </div>
              <CardHeader>
                <CardTitle>{{ course.title }}</CardTitle>
              </CardHeader>
              <CardContent class="flex-1">
                <CardDescription>{{ course.description }}</CardDescription>
              </CardContent>
              <CardFooter class="flex justify-center">
                <Button @click="openEditCourseModal(course)" variant="outline" size="sm">
                  <Pencil class="w-4 h-4 mr-2" />
                  Edit
                </Button>
              </CardFooter>
            </Card>
          </RouterLink>
        </div>
      </div>
    </template>

    <RouterView v-else />

    <CreateCourseDialog
      :is-open="isCreateModalOpen"
      :refresh-data="fetchMyCourses"
      @update:is-open="isCreateModalOpen = $event"
    />
    <EditCourseDialog
      :is-open="isEditModalOpen"
      :course="selectedCourse"
      :refresh-data="fetchMyCourses"
      @update:is-open="isEditModalOpen = $event"
    />
  </div>
</template>
