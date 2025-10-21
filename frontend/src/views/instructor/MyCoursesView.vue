<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterView, RouterLink, useRoute } from 'vue-router'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'

// Import components
import { Button } from '@/components/ui/button'
import { Card, CardDescription, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { AspectRatio } from '@/components/ui/aspect-ratio'
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

// Fetch my courses when the component mounts
onMounted(() => {
  fetchMyCourses()
})

// Open the create course modal
function openCreateCourseModal() {
  isCreateModalOpen.value = true
}

// Open the edit course modal
function openEditCourseModal(course: any) {
  selectedCourse.value = course
  isEditModalOpen.value = true
}

// Placeholder image URL
const placeholderImage = '/images/covers/placeholder.webp'

// Get the course image URL
function getCourseImageUrl(coverUrl: any) {
  if (coverUrl && coverUrl.Valid && coverUrl.String) {
    return coverUrl.String
  }
  return placeholderImage
}
</script>

<template>
  <div>
    <template v-if="!route.params.id">
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
        <Card v-for="i in 3" :key="i" class="rounded-lg overflow-hidden">
          <AspectRatio :ratio="16 / 9">
            <Skeleton class="w-full h-full" />
          </AspectRatio>
          <div class="p-4">
            <Skeleton class="h-5 w-3/4 mb-2" />
            <Skeleton class="h-4 w-full mb-1" />
            <Skeleton class="h-4 w-1/2 mb-4" />
            <Skeleton class="h-9 w-24 mx-auto" />
          </div>
        </Card>
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
            <Card class="h-full overflow-hidden rounded-lg flex flex-col">
              <AspectRatio :ratio="16 / 9" class="bg-muted mr-4 ml-4 rounded-lg overflow-hidden">
                <img
                  :src="getCourseImageUrl(course.cover_image_url)"
                  :alt="course.title"
                  class="w-full h-full object-cover transition-transform hover:scale-105"
                />
              </AspectRatio>
              <div class="pr-4 pl-4 flex flex-col flex-1">
                <CardTitle class="mb-1 text-lg">{{ course.title }}</CardTitle>
                <CardDescription class="text-sm mb-4 flex-1">
                  {{ course.description }}
                </CardDescription>

                <div class="flex justify-center">
                  <Button
                    @click="openEditCourseModal(course)"
                    variant="outline"
                    size="sm"
                    class="cursor-pointer"
                  >
                    <Pencil class="w-4 h-4 mr-2" />
                    Edit
                  </Button>
                </div>
              </div>
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
