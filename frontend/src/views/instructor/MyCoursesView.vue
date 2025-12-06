<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { RouterView, RouterLink, useRoute } from 'vue-router'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'

// Import components
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardDescription, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { AspectRatio } from '@/components/ui/aspect-ratio'
import { PlusCircle, FolderOpen, Pencil, Trash2, MoreVertical, Search } from 'lucide-vue-next'
import CreateCourseDialog from '@/components/instructor/my-courses/CreateCourseDialog.vue'
import EditCourseDialog from '@/components/instructor/my-courses/EditCourseDialog.vue'

// State variables
const courses = ref<any[]>([])
const isLoading = ref(true)
const route = useRoute()
const searchQuery = ref('')
// Modal states
const isCreateModalOpen = ref(false)
const isEditModalOpen = ref(false)
const isDeleteModalOpen = ref(false)
const selectedCourse = ref<any | null>(null)

// Computed properties
const filteredCourses = computed(() => {
  // If no search query, return all courses
  if (!searchQuery.value) {
    return courses.value
  }
  // Filter courses based on search query
  return courses.value.filter((course) =>
    course.title.toLowerCase().includes(searchQuery.value.toLowerCase()),
  )
})

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

// Open the delete course confirmation dialog
function openDeleteCourseDialog(course: any) {
  selectedCourse.value = course
  isDeleteModalOpen.value = true
}

// Handle delete course confirmation
async function handleDeleteCourse() {
  if (!selectedCourse.value) return
  try {
    await api.delete(`/instructor/courses/${selectedCourse.value.id}`)
    toast.success('Course deleted successfully.')
    fetchMyCourses()
  } catch (error) {
    toast.error('Failed to delete course.')
  } finally {
    isDeleteModalOpen.value = false
  }
}

// Stop event propagation
function stopPropagation(event: Event) {
  event.preventDefault()
  event.stopPropagation()
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
        <Button @click="openCreateCourseModal" class="cursor-pointer">
          <PlusCircle class="w-4 h-4 mr-2" />
          Create New Course
        </Button>
      </div>

      <div class="rounded mb-6">
        <div class="relative w-full max-w-sm">
          <Input
            id="search"
            type="text"
            placeholder="Search by title..."
            class="pl-10"
            v-model="searchQuery"
          />
          <span class="absolute inset-y-0 left-0 flex items-center pl-3">
            <Search class="h-5 w-5 text-muted-foreground" />
          </span>
        </div>
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
          <FolderOpen class="mx-auto h-12 w-12 text-muted-foreground mb-3" :stroke-width="1.5" />
          <h3 class="text-xl font-semibold">No Courses Found</h3>
          <p class="text-muted-foreground mt-2">
            You haven't created any courses yet. Get started now!
          </p>
          <Button @click="openCreateCourseModal" class="mt-4"> Create Your First Course </Button>
        </div>

        <div v-else>
          <div
            v-if="filteredCourses.length === 0"
            class="text-center py-16 border-2 border-dashed rounded-lg"
          >
            <FolderOpen class="mx-auto h-12 w-12 text-muted-foreground mb-3" :stroke-width="1.5" />
            <h3 class="text-xl font-semibold">No Courses Found</h3>
            <p class="text-muted-foreground mt-2">
              No courses match your search query "{{ searchQuery }}".
            </p>
          </div>

          <div v-else class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            <div
              v-for="course in filteredCourses"
              :key="course.id"
              class="hover:bg-accent transition-colors h-full overflow-hidden rounded-lg flex flex-col group"
            >
              <RouterLink :to="{ name: 'instructor-course-detail', params: { id: course.id } }">
                <Card class="h-full overflow-hidden rounded-lg flex flex-col">
                  <AspectRatio
                    :ratio="16 / 9"
                    class="bg-muted mr-4 ml-4 rounded-lg overflow-hidden"
                  >
                    <img
                      :src="getCourseImageUrl(course.cover_image_url)"
                      :alt="course.title"
                      class="w-full h-full object-cover transition-transform hover:scale-105"
                    />
                  </AspectRatio>
                  <div class="pr-4 pl-4 flex flex-col flex-1">
                    <RouterLink
                      :to="{ name: 'instructor-course-detail', params: { id: course.id } }"
                    >
                      <CardTitle class="mb-1 text-lg hover:underline">{{ course.title }}</CardTitle>
                    </RouterLink>
                    <CardDescription class="text-sm mb-4 flex-1">
                      {{ course.description }}
                    </CardDescription>
                    <DropdownMenu>
                      <DropdownMenuTrigger as-child>
                        <Button
                          @click="stopPropagation"
                          variant="outline"
                          size="sm"
                          class="cursor-pointer w-24"
                        >
                          <MoreVertical class="w-4 h-4 mr-2" />
                          Actions
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="center">
                        <DropdownMenuItem @click="openEditCourseModal(course)">
                          <Pencil class="w-4 h-4 mr-2" />
                          Edit
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          @click="openDeleteCourseDialog(course)"
                          class="text-red-600 focus:text-red-600 focus:bg-red-50"
                        >
                          <Trash2 class="w-4 h-4 mr-2" />
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </Card>
              </RouterLink>
            </div>
          </div>
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

    <AlertDialog v-model:open="isDeleteModalOpen">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This action cannot be undone. This will permanently delete the course
            <strong>"{{ selectedCourse?.title }}"</strong>
            and all its materials.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction @click="handleDeleteCourse" class="bg-red-600 hover:bg-red-700">
            Continue
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
