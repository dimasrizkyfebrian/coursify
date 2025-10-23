<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'

import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
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
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { AspectRatio } from '@/components/ui/aspect-ratio'
import { Users, PlusCircle, Pencil, Trash2, FolderOpen } from 'lucide-vue-next'

import AddMaterialDialog from '@/components/instructor/my-courses/AddMaterialDialog.vue'
import EditMaterialDialog from '@/components/instructor/my-courses/EditMaterialDialog.vue'

// State variables
const route = useRoute()
const course = ref<any>(null)
const materials = ref<any[]>([])
const isLoading = ref(true)

// Tabs states
const activeTab = ref('materials')
const enrolledStudents = ref<any[]>([])
const isStudentsLoading = ref(false)
const studentsLoaded = ref(false)

// Modal states
const isAddModalOpen = ref(false)
const isEditModalOpen = ref(false)
const selectedMaterial = ref<any | null>(null)
const isDeleteDialogOpen = ref(false)

// Fetch course details and materials
async function fetchCourseDetails() {
  const courseId = route.params.id
  try {
    const courseResponse = await api.get(`/instructor/courses/${courseId}`)
    course.value = courseResponse.data
  } catch (error) {
    toast.error('Failed to load course details.')
  }
}

// Fetch course materials
async function fetchCourseMaterials() {
  if (!course.value?.id) return
  try {
    const materialsResponse = await api.get(`/instructor/courses/${course.value.id}/materials`)
    materials.value = materialsResponse.data || []
  } catch (error) {
    toast.error('Failed to load course materials.')
  }
}

// Fetch enrolled students
async function fetchEnrolledStudents() {
  if (!course.value?.id || studentsLoaded.value) return
  try {
    isStudentsLoading.value = true
    const response = await api.get(`/instructor/courses/${course.value.id}/enrollments`)
    enrolledStudents.value = response.data || []
    studentsLoaded.value = true
  } catch (error) {
    toast.error('Failed to load enrolled students.')
  } finally {
    isStudentsLoading.value = false
  }
}

// Handle delete material confirmation
async function handleDeleteConfirm() {
  if (!selectedMaterial.value) return
  try {
    await api.delete(
      `/instructor/courses/${course.value.id}/materials/${selectedMaterial.value.id}`,
    )
    toast.success('Material deleted successfully.')
    fetchCourseMaterials() // Refresh materials
  } catch (error) {
    toast.error('Failed to delete material.')
  } finally {
    isDeleteDialogOpen.value = false // Close dialog
  }
}

// Watch for changes in the active tab
watch(activeTab, async (newTab) => {
  if (newTab === 'materials' && materials.value.length === 0) {
    // If switch to the materials tab and there is no data yet, fetch it
    await fetchCourseMaterials()
  } else if (newTab === 'people') {
    // If switch to the people tab, fetch the student data (if not done before)
    await fetchEnrolledStudents()
  }
})

// Fetch course data when the component mounts
onMounted(async () => {
  isLoading.value = true
  await fetchCourseDetails()
  if (activeTab.value === 'materials') {
    await fetchCourseMaterials()
  }
  isLoading.value = false
})

// Open the add material modal
function openAddMaterialModal() {
  isAddModalOpen.value = true
}

// Open the edit material modal
function openEditMaterialModal(material: any) {
  selectedMaterial.value = material
  isEditModalOpen.value = true
}

// Open the delete material dialog
function openDeleteMaterialDialog(material: any) {
  selectedMaterial.value = material
  isDeleteDialogOpen.value = true
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
    <div v-if="isLoading">
      <AspectRatio :ratio="21 / 9" class="bg-muted rounded-lg overflow-hidden mb-6">
        <Skeleton class="w-full h-full" />
      </AspectRatio>
      <div class="max-w-4xl mx-auto px-4">
        <Skeleton class="h-10 w-3/4 mb-4" />
        <Skeleton class="h-5 w-full mb-2" />
        <Skeleton class="h-5 w-2/3 mb-8" />
        <Skeleton class="h-12 w-full mb-4" /> {/* Untuk TabsList */}
        <Skeleton class="h-20 w-full" /> {/* Untuk TabsContent */}
      </div>
    </div>

    <div v-else-if="course">
      <div class="mb-6 rounded-lg overflow-hidden shadow-lg">
        <AspectRatio :ratio="21 / 9" class="bg-muted">
          <img
            :src="getCourseImageUrl(course.cover_image_url)"
            :alt="course.title"
            class="w-full h-full object-cover"
          />
        </AspectRatio>
      </div>

      <div class="max-w-4xl mx-auto px-4">
        <h1 class="text-3xl md:text-4xl font-extrabold tracking-tight mb-2">
          {{ course.title }}
        </h1>
        <p class="text-lg text-muted-foreground mb-10">
          {{ course.description }}
        </p>

        <Tabs v-model="activeTab" default-value="materials" class="w-full">
          <TabsList class="grid w-full grid-cols-2 mb-6">
            <TabsTrigger value="materials" class="cursor-pointer">Materials</TabsTrigger>
            <TabsTrigger value="people" class="cursor-pointer">People</TabsTrigger>
          </TabsList>

          <TabsContent value="materials">
            <div>
              <div class="flex items-center justify-between mb-4">
                <h2 class="text-xl font-semibold">Course Materials</h2>
                <Button @click="openAddMaterialModal" size="sm">
                  <PlusCircle class="w-4 h-4 mr-2" />
                  Add Material
                </Button>
              </div>

              <Card v-if="materials.length === 0" class="text-center py-12 border-2 border-dashed">
                <CardContent>
                  <FolderOpen
                    class="mx-auto h-12 w-12 text-muted-foreground mb-3"
                    :stroke-width="1.5"
                  />
                  <h3 class="text-lg font-semibold">No Materials Yet</h3>
                  <p class="text-sm text-muted-foreground mt-1">
                    Start building your course by adding the first material.
                  </p>
                </CardContent>
              </Card>

              <div v-else class="space-y-4">
                <Card v-for="material in materials" :key="material.id">
                  <CardContent class="p-4 flex items-center justify-between">
                    <div>
                      <p class="font-semibold">{{ material.title }}</p>
                      <span class="text-xs text-muted-foreground uppercase">{{
                        material.content_type
                      }}</span>
                    </div>
                    <div class="flex gap-2">
                      <Button @click="openEditMaterialModal(material)" variant="outline" size="sm">
                        <Pencil class="w-4 h-4 mr-2" />
                        Edit
                      </Button>
                      <Button
                        @click="openDeleteMaterialDialog(material)"
                        variant="destructive"
                        size="sm"
                      >
                        <Trash2 class="w-4 h-4" />
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </div>
          </TabsContent>

          <TabsContent value="people">
            <div>
              <div class="mb-8">
                <h2 class="text-xl font-semibold mb-3 pb-2">Instructor</h2>
                <Card v-if="course.instructor_name">
                  <CardContent class="p-4 flex items-center">
                    <div>
                      <p class="font-semibold">{{ course.instructor_name }}</p>
                    </div>
                  </CardContent>
                </Card>
                <Card v-else>
                  <CardContent class="p-4">
                    <p class="font-semibold">Instructor ID:</p>
                    <p class="text-sm text-muted-foreground">{{ course.instructor_id }}</p>
                  </CardContent>
                </Card>
              </div>

              <h2 class="text-xl font-semibold mb-4">Enrolled Students</h2>

              <div v-if="isStudentsLoading" class="space-y-3">
                <Skeleton class="h-10 w-full" v-for="i in 3" :key="`skel-${i}`" />
              </div>

              <Card
                v-else-if="enrolledStudents.length === 0"
                class="text-center py-12 border-2 border-dashed"
              >
                <CardContent>
                  <Users class="mx-auto h-12 w-12 text-muted-foreground mb-3" :stroke-width="1.5" />
                  <h3 class="text-lg font-semibold">No Students Enrolled</h3>
                  <p class="text-sm text-muted-foreground mt-1">
                    No students have enrolled in this course yet.
                  </p>
                </CardContent>
              </Card>

              <div v-else class="space-y-3">
                <Card v-for="student in enrolledStudents" :key="student.id">
                  <CardContent class="p-4 flex items-center justify-between">
                    <div>
                      <p class="font-semibold">{{ student.full_name }}</p>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </div>

    <AddMaterialDialog
      :is-open="isAddModalOpen"
      :course-id="course?.id || null"
      :refresh-data="fetchCourseMaterials"
      @update:is-open="isAddModalOpen = $event"
    />

    <EditMaterialDialog
      :is-open="isEditModalOpen"
      :material="selectedMaterial"
      :refresh-data="fetchCourseMaterials"
      @update:is-open="isEditModalOpen = $event"
    />

    <AlertDialog v-model:open="isDeleteDialogOpen">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This action cannot be undone. This will permanently delete the material titled
            <strong>"{{ selectedMaterial?.title }}"</strong>.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction @click="handleDeleteConfirm">Continue</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
