<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, RouterView } from 'vue-router'
import DashboardLayout from '@/layouts/DashboardLayout.vue'
import { Toaster } from '@/components/ui/sonner'
import 'vue-sonner/style.css'

const route = useRoute()
const needsDashboardLayout = computed(() => route.meta.requiresAuth)
</script>

<template>
  <Toaster />
  <DashboardLayout v-if="needsDashboardLayout">
    <RouterView />
  </DashboardLayout>
  <RouterView v-else />
</template>
