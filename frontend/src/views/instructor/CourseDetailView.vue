<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'

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
import { PlusCircle, Pencil, Trash2, FolderOpen } from 'lucide-vue-next'

import AddMaterialDialog from '@/components/instructor/my-courses/AddMaterialDialog.vue'
import EditMaterialDialog from '@/components/instructor/my-courses/EditMaterialDialog.vue'

// State variables
const route = useRoute()
const course = ref<any>(null)
const materials = ref<any[]>([])
const isLoading = ref(true)

const isAddModalOpen = ref(false)
const isEditModalOpen = ref(false)
const selectedMaterial = ref<any | null>(null)
const isDeleteDialogOpen = ref(false)

onMounted(async () => {
  const courseId = route.params.id
  try {
    const [courseResponse, materialsResponse] = await Promise.all([
      api.get(`/instructor/courses/${courseId}`),
      api.get(`/instructor/courses/${courseId}/materials`),
    ])

    course.value = courseResponse.data
    materials.value = materialsResponse.data || []
  } catch (error) {
    toast.error('Failed to load course data.')
  } finally {
    isLoading.value = false
  }
})

// Fetch course data from the API
async function fetchCourseData() {
  const courseId = route.params.id
  try {
    isLoading.value = true
    const [courseResponse, materialsResponse] = await Promise.all([
      api.get(`/instructor/courses/${courseId}`),
      api.get(`/instructor/courses/${courseId}/materials`),
    ])
    course.value = courseResponse.data
    materials.value = materialsResponse.data || []
  } catch (error) {
    toast.error('Failed to load course data.')
  } finally {
    isLoading.value = false
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
    fetchCourseData() // Refresh the course data
  } catch (error) {
    toast.error('Failed to delete material.')
  } finally {
    isDeleteDialogOpen.value = false // Close dialog
  }
}

// Fetch course data when the component mounts
onMounted(() => {
  fetchCourseData()
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
        <Skeleton class="h-12 w-full" />
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
      </div>
    </div>

    <AddMaterialDialog
      :is-open="isAddModalOpen"
      :course-id="course?.id || null"
      :refresh-data="fetchCourseData"
      @update:is-open="isAddModalOpen = $event"
    />

    <EditMaterialDialog
      :is-open="isEditModalOpen"
      :material="selectedMaterial"
      :refresh-data="fetchCourseData"
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
