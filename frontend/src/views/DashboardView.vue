<script setup lang="ts">
import { useUserStore } from '@/stores/user'
// Dashboard components
import AdminDashboard from '@/components/dashboards/AdminDashboard.vue'
import InstructorDashboard from '@/components/dashboards/InstructorDashboard.vue'
import StudentDashboard from '@/components/dashboards/StudentDashboard.vue'

const { user } = useUserStore()
</script>

<template>
  <AdminDashboard v-if="user.role === 'admin'" />
  <InstructorDashboard v-else-if="user.role === 'instructor'" />
  <StudentDashboard v-else-if="user.role === 'student'" />

  <div v-else>
    <p>Loading user data or role not recognized...</p>
  </div>
</template>
