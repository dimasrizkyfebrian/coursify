import { createRouter, createWebHistory } from 'vue-router'
import LandingView from '../views/LandingView.vue'
import RegisterView from '../views/RegisterView.vue'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import AdminDashboard from '@/components/dashboards/AdminDashboard.vue'
import InstructorDashboard from '@/components/dashboards/InstructorDashboard.vue'
import StudentDashboard from '@/components/dashboards/StudentDashboard.vue'
import UserManagementView from '../views/admin/UserManagement.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: LandingView,
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true },
    },
    {
      path: '/admin/users',
      name: 'admin-users',
      component: UserManagementView,
      meta: { requiresAuth: true },
    },
  ],
})

// Navigation Guard
router.beforeEach((to, from, next) => {
  const isLoggedIn = !!localStorage.getItem('authToken')

  // Check if the route requires authentication and the user is not logged in
  if (to.meta.requiresAuth && !isLoggedIn) {
    next({ name: 'login' }) // Redirect to the login page if not authenticated
  } else {
    next() // Allow navigation if authenticated
  }
})

export default router
