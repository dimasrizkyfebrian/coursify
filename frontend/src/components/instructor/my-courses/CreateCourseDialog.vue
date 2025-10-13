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
  isOpen: boolean
  refreshData?: () => void
}>()

const emit = defineEmits(['update:isOpen'])

const formData = ref({
  title: '',
  description: '',
})
const isSubmitting = ref(false)

watch(
  () => props.isOpen,
  (newValue) => {
    if (newValue) {
      formData.value = { title: '', description: '' }
      isSubmitting.value = false
    }
  },
)

async function handleCreateCourse() {
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
</script>

<template>
  <Dialog :open="isOpen" @update:open="(value) => emit('update:isOpen', value)">
    <DialogContent class="sm:max-w-xl">
      <DialogHeader>
        <DialogTitle>Create a New Course</DialogTitle>
        <DialogDescription>
          Fill in the details below to create a new course. You can add materials later.
        </DialogDescription>
      </DialogHeader>
      <div class="grid gap-4 py-4">
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
          />
        </div>
      </div>
      <DialogFooter>
        <Button @click="handleCreateCourse" :disabled="isSubmitting" type="submit">
          {{ isSubmitting ? 'Creating...' : 'Create Course' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
