<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Button } from '@/components/ui/button'
import { PlusCircle } from 'lucide-vue-next'

const route = useRoute()
const course = ref<any>(null)
const materials = ref<any[]>([])
const isLoading = ref(true)

onMounted(async () => {
  const courseId = route.params.id
  try {
    // Fetch course details and materials
    const [courseResponse, materialsResponse] = await Promise.all([
      api.get(`/instructor/courses/${courseId}`),
      api.get(`/instructor/courses/${courseId}/materials`),
    ])

    course.value = courseResponse.data
    materials.value = materialsResponse.data
  } catch (error) {
    toast.error('Failed to load course data.')
    console.error('Failed to fetch course data:', error)
  } finally {
    isLoading.value = false
  }
})

function openAddMaterialModal() {
  toast.info('Add Material button clicked.')
}
</script>

<template>
  <div>
    <div v-if="isLoading">
      <Skeleton class="h-8 w-1/2 mb-4" />
      <Skeleton class="h-24 w-full mb-8" />
      <Skeleton class="h-40 w-full" />
    </div>

    <div v-else-if="course">
      <Card class="mb-6">
        <img
          src="https://placehold.co/1080x720?text=Course+Image"
          alt="Course Image Placeholder"
          class="aspect-video w-full object-cover rounded md:rounded-t-2xl"
        />
        <CardHeader>
          <CardTitle class="text-3xl">{{ course.title }}</CardTitle>
        </CardHeader>
        <CardContent>
          <CardDescription>{{ course.description }}</CardDescription>
        </CardContent>
      </Card>

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
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  </div>
</template>
