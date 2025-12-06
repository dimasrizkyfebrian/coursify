<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { Button } from '@/components/ui/button'
</script>

<template>
  <div class="flex flex-col items-center justify-center min-h-screen text-center p-4">
    <h1 class="text-5xl font-bold mb-4">Welcome to Coursify</h1>
    <p class="text-xl text-gray-600 dark:text-gray-300 mb-8 max-w-2xl">
      The modern platform to manage and deliver your online courses seamlessly. Built with Go and
      Vue.js.
    </p>
    <RouterLink to="/login">
      <Button size="lg"> Get Started / Login </Button>
    </RouterLink>
  </div>
</template>
