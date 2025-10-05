<script setup lang="ts">
import AppSidebar from '@/components/layout/AppSidebar.vue'
import { SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
</script>

<template>
  <SidebarProvider>
    <AppSidebar />
    <main class="flex-1 p-8 overflow-y-auto">
      <SidebarTrigger />
      <slot />
    </main>
  </SidebarProvider>
</template>
