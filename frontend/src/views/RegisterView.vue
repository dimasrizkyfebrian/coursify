<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import axios from 'axios'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { toast } from 'vue-sonner'

const fullName = ref('')
const email = ref('')
const password = ref('')

// Function to handle form submission
async function handleSubmit() {
  try {
    const response = await axios.post('http://localhost:8080/api/register', {
      full_name: fullName.value,
      email: email.value,
      password: password.value,
      role: 'student', // Set the role to 'student'
    })

    toast.success('Registration successful!', {
      description: 'Please wait for your account to be approved by an admin.',
    })

    // Clear the form fields after successful registration
    fullName.value = ''
    email.value = ''
    password.value = ''
  } catch (error) {
    toast.error('Registration failed.', {
      description: 'Please check your details and try again.',
    })
  }
}
</script>

<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100 dark:bg-gray-900">
    <form @submit.prevent="handleSubmit" class="w-full max-w-md">
      <Card>
        <CardHeader class="text-center">
          <CardTitle class="text-2xl font-bold"> Create an Account </CardTitle>
          <CardDescription>
            Enter your information to create a new account in Coursify.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <div class="space-y-2">
              <Label for="full-name">Full Name</Label>
              <Input id="full-name" v-model="fullName" placeholder="Jhon Smith" required />
            </div>
            <div class="space-y-2">
              <Label for="email">Email</Label>
              <Input
                id="email"
                v-model="email"
                type="email"
                placeholder="newuser@example.com"
                required
              />
            </div>
            <div class="space-y-2">
              <Label for="password">Password</Label>
              <Input id="password" v-model="password" type="password" required />
            </div>
            <Button type="submit" class="w-full"> Create Account </Button>
          </div>
        </CardContent>
        <CardFooter class="text-sm text-center justify-center">
          Already have an account?&nbsp;
          <RouterLink to="/login" class="font-semibold underline"> Login </RouterLink>
        </CardFooter>
      </Card>
    </form>
  </div>
</template>
