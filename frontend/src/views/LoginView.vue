<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import axios from 'axios'
import { toast } from 'vue-sonner'

import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

const email = ref('')
const password = ref('')
const router = useRouter()

// Function to handle form submission
async function handleSubmit() {
  try {
    const response = await axios.post('http://localhost:8080/api/login', {
      email: email.value,
      password: password.value,
    })

    // Token received from the server
    const token = response.data.token
    localStorage.setItem('authToken', token)

    toast.success('Login successful!', {
      description: 'Redirecting you to the dashboard...',
    })

    // Redirect to the dashboard after successful login
    router.push('/dashboard')
  } catch (error: any) {
    // Handle different error scenarios
    if (axios.isAxiosError(error) && error.response) {
      if (error.response.status === 403) {
        // Handle forbidden error (e.g., account not approved)
        toast.error('Login Failed', {
          description: error.response.data || 'Your account is pending approval.',
        })
      } else if (error.response.status === 401) {
        // Handle unauthorized error (e.g., invalid credentials)
        toast.error('Login Failed', {
          description: error.response.data || 'Invalid email or password.',
        })
      } else {
        // Handle other error scenarios
        toast.error('An unexpected error occurred.', {
          description: 'Please try again later.',
        })
      }
    } else {
      // Handle other types of errors
      toast.error('Network Error', {
        description: 'Could not connect to the server.',
      })
    }
  }
}
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100 dark:bg-gray-900">
    <form @submit.prevent="handleSubmit" class="w-full max-w-md">
      <Card>
        <CardHeader class="text-center">
          <CardTitle class="text-2xl font-bold"> Login to Coursify </CardTitle>
          <CardDescription> Enter your email and password to access your account. </CardDescription>
        </CardHeader>
        <CardContent class="space-y-4">
          <div class="space-y-2">
            <Label for="email">Email</Label>
            <Input
              id="email"
              v-model="email"
              type="email"
              placeholder="user@example.com"
              required
            />
          </div>
          <div class="space-y-2">
            <Label for="password">Password</Label>
            <Input id="password" v-model="password" type="password" required />
          </div>
          <Button type="submit" class="w-full"> Login </Button>
        </CardContent>
        <CardFooter class="text-sm text-center justify-center">
          Don't have an account?&nbsp;
          <RouterLink to="/register" class="font-semibold underline"> Register </RouterLink>
        </CardFooter>
      </Card>
    </form>
  </div>
</template>
