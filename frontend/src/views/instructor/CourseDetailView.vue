<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/lib/axios'
import { toast } from 'vue-sonner'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'

const route = useRoute()
const course = ref<any>(null)
const isLoading = ref(true)

onMounted(async () => {
  const courseId = route.params.id
  try {
    const response = await api.get(`/instructor/courses/${courseId}`)
    course.value = response.data
  } catch (error) {
    toast.error('Failed to load course details.')
    console.error('Failed to fetch course details:', error)
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <div>
    <div v-if="isLoading">
      <Skeleton class="h-8 w-1/2 mb-4" />
      <Skeleton class="h-24 w-full" />
    </div>

    <Card v-else-if="course">
      <div class="mr-4 ml-4">
        <img
          src="https://placehold.co/1080x720?text=Course+Image"
          alt="Course Image Placeholder"
          class="aspect-video w-full object-cover rounded md:rounded-lg"
        />
      </div>
      <CardHeader>
        <CardTitle class="text-3xl">{{ course.title }}</CardTitle>
      </CardHeader>
      <CardContent>
        <CardDescription>{{ course.description }}</CardDescription>
        <div class="mt-8">
          <h2 class="text-xl font-semibold">Course Materials</h2>
          <p class="text-sm text-muted-foreground mt-2">
            This is where you will add and manage your course materials (videos, text, PDFs).
          </p>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
