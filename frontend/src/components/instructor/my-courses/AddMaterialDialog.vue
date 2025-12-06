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
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

const props = defineProps<{
  courseId: string | null
  isOpen: boolean
  refreshData?: () => void
}>()

const emit = defineEmits(['update:isOpen'])

const formData = ref({
  title: '',
  content_type: 'text',
  text_content: '',
  video_url: '',
})
const isSubmitting = ref(false)

// Reset form when the modal opens
watch(
  () => props.isOpen,
  (newValue) => {
    if (newValue) {
      formData.value = { title: '', content_type: 'text', text_content: '', video_url: '' }
      isSubmitting.value = false
    }
  },
)

async function handleAddMaterial() {
  if (!props.courseId) return
  isSubmitting.value = true
  try {
    const response = await api.post(
      `/instructor/courses/${props.courseId}/materials`,
      formData.value,
    )
    toast.success('Material added successfully!')
    emit('update:isOpen', false) // Close modal
    if (props.refreshData) props.refreshData() // Refresh material list
  } catch (error) {
    toast.error('Failed to add material.')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <Dialog :open="isOpen" @update:open="(value) => emit('update:isOpen', value)">
    <DialogContent class="sm:max-w-xl">
      <DialogHeader>
        <DialogTitle>Add New Material</DialogTitle>
        <DialogDescription> Select the material type and fill in the details. </DialogDescription>
      </DialogHeader>
      <div class="grid gap-4 py-4">
        <div class="grid items-center gap-2">
          <Label for="title">Material Title</Label>
          <Input
            id="title"
            v-model="formData.title"
            placeholder="e.g., Chapter 1: Getting Started"
          />
        </div>
        <div class="grid items-center gap-2">
          <Label for="content-type">Content Type</Label>
          <Select v-model="formData.content_type">
            <SelectTrigger>
              <SelectValue placeholder="Select a content type" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectItem value="text">Text</SelectItem>
                <SelectItem value="video">Video</SelectItem>
                <SelectItem value="pdf">PDF (Coming Soon)</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
        </div>
        <div v-if="formData.content_type === 'text'" class="grid items-center gap-2">
          <Label for="text-content">Text Content</Label>
          <Textarea
            id="text-content"
            v-model="formData.text_content"
            placeholder="Write your lesson content here..."
            class="min-h-[150px]"
          />
        </div>
        <div v-if="formData.content_type === 'video'" class="grid items-center gap-2">
          <Label for="video-url">Video URL</Label>
          <Input
            id="video-url"
            v-model="formData.video_url"
            placeholder="e.g., https://www.youtube.com/watch?v=..."
          />
        </div>
      </div>
      <DialogFooter>
        <Button @click="handleAddMaterial" :disabled="isSubmitting" type="submit">
          {{ isSubmitting ? 'Adding...' : 'Add Material' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
