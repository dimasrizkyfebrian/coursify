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
  material: any | null
  isOpen: boolean
  refreshData?: () => void
}>()

const emit = defineEmits(['update:isOpen'])

const formData = ref({
  title: '',
  text_content: '',
  video_url: '',
})
const isSubmitting = ref(false)

// Fill form with material data when the modal opens
watch(
  () => props.material,
  (newMaterial) => {
    if (newMaterial) {
      formData.value = {
        title: newMaterial.title,
        text_content: newMaterial.text_content,
        video_url: newMaterial.video_url,
      }
    }
  },
)

async function handleUpdateMaterial() {
  if (!props.material) return
  isSubmitting.value = true
  try {
    const response = await api.put(
      `/instructor/courses/${props.material.course_id}/materials/${props.material.id}`,
      formData.value,
    )
    toast.success(response.data.message || 'Material updated successfully!')
    emit('update:isOpen', false) // Close modal
    if (props.refreshData) props.refreshData() // Refresh material list
  } catch (error) {
    toast.error('Failed to update material.')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <Dialog :open="isOpen" @update:open="(value) => emit('update:isOpen', value)">
    <DialogContent class="sm:max-w-xl">
      <DialogHeader>
        <DialogTitle>Edit Material</DialogTitle>
        <DialogDescription> Make changes to your material content below. </DialogDescription>
      </DialogHeader>
      <div class="grid gap-4 py-4">
        <div class="grid items-center gap-2">
          <Label for="title">Material Title</Label>
          <Input id="title" v-model="formData.title" />
        </div>
        <div v-if="material?.content_type === 'text'" class="grid items-center gap-2">
          <Label for="text-content">Text Content</Label>
          <Textarea id="text-content" v-model="formData.text_content" class="min-h-[150px]" />
        </div>
        <div v-if="material?.content_type === 'video'" class="grid items-center gap-2">
          <Label for="video-url">Video URL</Label>
          <Input id="video-url" v-model="formData.video_url" />
        </div>
      </div>
      <DialogFooter>
        <Button @click="handleUpdateMaterial" :disabled="isSubmitting" type="submit">
          {{ isSubmitting ? 'Saving...' : 'Save Changes' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
