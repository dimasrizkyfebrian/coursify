<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import api from '@/lib/axios'
import { courseCoverImages } from '@/lib/courseAssets'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { AspectRatio } from '@/components/ui/aspect-ratio'

const props = defineProps<{
  isOpen: boolean
  refreshData?: () => void
}>()

const emit = defineEmits(['update:isOpen'])

const formData = ref({
  title: '',
  description: '',
  cover_image_url: '',
})
const isSubmitting = ref(false)

const isFormInvalid = computed(() => {
  return (
    !formData.value.title.trim() ||
    !formData.value.description.trim() ||
    !formData.value.cover_image_url
  )
})

watch(
  () => props.isOpen,
  (newValue) => {
    if (newValue) {
      formData.value = { title: '', description: '', cover_image_url: '' }
      isSubmitting.value = false
    }
  },
)

async function handleCreateCourse() {
  // Check if the form is invalid
  if (isFormInvalid.value) {
    toast.error('Validation Error', {
      description: 'Title, description, and a cover image are required.',
    })
    return // Exit the function if the form is invalid
  }

  isSubmitting.value = true
  try {
    const response = await api.post('/instructor/courses', formData.value)
    toast.success('Course created successfully!')
    emit('update:isOpen', false)
    if (props.refreshData) props.refreshData()
  } catch (error) {
    toast.error('Failed to create course.')
  } finally {
    isSubmitting.value = false
  }
}

const selectedImagePreview = computed(() => {
  return formData.value.cover_image_url || '/images/covers/placeholder.webp'
})
</script>

<template>
  <Dialog :open="isOpen" @update:open="(value) => emit('update:isOpen', value)">
    <DialogContent class="sm:max-w-3xl">
      <DialogHeader>
        <DialogTitle>Create a New Course</DialogTitle>
        <DialogDescription> Fill in the details below to create a new course. </DialogDescription>
      </DialogHeader>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 py-4">
        <div class="grid gap-4">
          <div class="grid items-center gap-2">
            <Label for="title">Title</Label>
            <Input id="title" v-model="formData.title" placeholder="e.g., Introduction to Go" />
          </div>
          <div class="grid items-center gap-2">
            <Label for="description">Description</Label>
            <Textarea
              id="description"
              v-model="formData.description"
              placeholder="Provide a brief description of your course..."
              class="min-h-[100px]"
            />
          </div>

          <div class="grid items-center gap-2">
            <Label>Select Cover Image</Label>
            <RadioGroup v-model="formData.cover_image_url" class="grid grid-cols-4 gap-2">
              <div v-for="image in courseCoverImages" :key="image.id">
                <RadioGroupItem :value="image.url" :id="image.id" class="sr-only" />
                <Label
                  :for="image.id"
                  class="cursor-pointer rounded-md border-2 border-transparent transition-all hover:opacity-75"
                  :class="{
                    'border-primary': formData.cover_image_url === image.url,
                  }"
                >
                  <AspectRatio :ratio="16 / 9">
                    <img
                      :src="image.url"
                      :alt="image.alt"
                      class="rounded-md object-cover w-full h-full"
                    />
                  </AspectRatio>
                </Label>
              </div>
            </RadioGroup>
          </div>
        </div>

        <div class="grid items-start gap-2">
          <Label>Image Preview</Label>
          <AspectRatio :ratio="16 / 9" class="bg-muted rounded-md">
            <img
              :src="selectedImagePreview"
              alt="Selected cover image preview"
              class="rounded-md object-cover w-full h-full"
            />
          </AspectRatio>
        </div>
      </div>

      <DialogFooter>
        <Button @click="handleCreateCourse" :disabled="isSubmitting || isFormInvalid" type="submit">
          {{ isSubmitting ? 'Creating...' : 'Create Course' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
