<script setup lang="ts">
import { ref, watch } from 'vue'
import api from '@/lib/axios'
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

const props = defineProps<{
  course: any | null
  isOpen: boolean
  refreshData?: () => void
}>()

const emit = defineEmits(['update:isOpen'])

const formData = ref({
  title: '',
  description: '',
})
const isSubmitting = ref(false)

// Fill out form with course data when modal is opened
watch(
  () => props.course,
  (newCourse) => {
    if (newCourse) {
      formData.value = {
        title: newCourse.title,
        description: newCourse.description,
      }
    }
  },
)

// Update course details on form submit
async function handleUpdateCourse() {
  if (!props.course) return
  isSubmitting.value = true
  try {
    const response = await api.put(`/instructor/courses/${props.course.id}`, formData.value)
    toast.success(response.data.message || 'Course updated successfully!')
    emit('update:isOpen', false) // Close the modal
    if (props.refreshData) props.refreshData() // Refresh the course list
  } catch (error) {
    toast.error('Failed to update course.')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <Dialog :open="isOpen" @update:open="(value) => emit('update:isOpen', value)">
    <DialogContent class="sm:max-w-xl">
      <DialogHeader>
        <DialogTitle>Edit Course</DialogTitle>
        <DialogDescription> Make changes to your course details below. </DialogDescription>
      </DialogHeader>
      <div class="grid gap-4 py-4">
        <div class="grid items-center gap-2">
          <Label for="title">Title</Label>
          <Input id="title" v-model="formData.title" />
        </div>
        <div class="grid items-center gap-2">
          <Label for="description">Description</Label>
          <Textarea id="description" v-model="formData.description" />
        </div>
      </div>
      <DialogFooter>
        <Button @click="handleUpdateCourse" :disabled="isSubmitting" type="submit">
          {{ isSubmitting ? 'Saving...' : 'Save Changes' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
