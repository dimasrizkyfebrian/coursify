import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import LandingView from '../views/LandingView.vue'
import RegisterView from '../views/RegisterView.vue'
import LoginView from '../views/LoginView.vue'
import DashboardView from '../views/DashboardView.vue'
import PendingApprovalView from '../views/admin/PendingApprovalView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'landing', component: LandingView, meta: { breadcrumb: 'Home' } },
    { path: '/register', name: 'register', component: RegisterView },
    { path: '/login', name: 'login', component: LoginView },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true, breadcrumb: 'Dashboard' },
    },
    {
      path: '/admin/pending-approval',
      name: 'admin-pending-approval',
      component: PendingApprovalView,
      meta: { requiresAuth: true, breadcrumb: 'Pending Approval' },
    },
  ],
})

// Navigation Guard
router.beforeEach((to, from, next) => {
  const { setBreadcrumbs } = useUserStore()
  const isLoggedIn = !!localStorage.getItem('authToken')

  const crumbs = to.matched
    .filter((item) => item.meta && item.meta.breadcrumb)
    .map((item) => ({
      text: item.meta.breadcrumb as string,
      to: item.path,
    }))
  setBreadcrumbs(crumbs)

  // Check if the route requires authentication and the user is not logged in
  if (to.meta.requiresAuth && !isLoggedIn) {
    next({ name: 'login' }) // Redirect to the login page if not authenticated
  } else {
    next() // Allow navigation if authenticated
  }
})

export default router
