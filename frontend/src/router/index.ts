import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'
import LandingView from '../views/LandingView.vue'
import RegisterView from '../views/RegisterView.vue'
import LoginView from '../views/LoginView.vue'
import DashboardLayout from '../layouts/DashboardLayout.vue'
import DashboardView from '@/views/DashboardView.vue'
// --- Admin ---
import PendingApprovalView from '../views/admin/PendingApprovalView.vue'
import PendingApprovalListView from '../views/admin/PendingApprovalListView.vue'
import UserDetailView from '../views/admin/UserDetailView.vue'
import AllUsersView from '@/views/admin/AllUsersView.vue'
import SettingsView from '@/views/admin/SettingsView.vue'
// --- Instructor ---
import MyCoursesView from '../views/instructor/MyCoursesView.vue'
import CourseDetailView from '../views/instructor/CourseDetailView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    // --- Public Routes ---
    { path: '/', name: 'landing', component: LandingView, meta: { breadcrumb: 'Home' } },
    { path: '/register', name: 'register', component: RegisterView },
    { path: '/login', name: 'login', component: LoginView },

    // --- Protected Routes ---
    {
      path: '/dashboard',
      component: DashboardLayout,
      meta: { requiresAuth: true, breadcrumb: 'Dashboard' },
      redirect: { name: 'dashboard-home' },
      children: [
        // --- Admin Routes ---
        {
          path: '',
          name: 'dashboard-home',
          component: DashboardView,
        },
        {
          path: 'admin/pending-approval',
          component: PendingApprovalView,
          meta: { breadcrumb: 'Pending Approval' },
          children: [
            { path: '', name: 'admin-pending-list', component: PendingApprovalListView },
            {
              path: ':id',
              name: 'admin-pending-detail',
              component: UserDetailView,
              meta: { breadcrumb: 'User Detail' },
            },
          ],
        },
        {
          path: 'admin/all-users',
          name: 'admin-all-users',
          component: AllUsersView,
          meta: { breadcrumb: 'All Users' },
        },
        {
          path: 'admin/settings',
          name: 'admin-settings',
          component: SettingsView,
          meta: { breadcrumb: 'Settings' },
        },
        // --- Instructor Routes ---
        {
          path: 'instructor/courses',
          name: 'instructor-my-courses',
          component: MyCoursesView,
          meta: { breadcrumb: 'My Courses' },
          children: [
            {
              path: ':id',
              name: 'instructor-course-detail',
              component: CourseDetailView,
              meta: { breadcrumb: 'Course Detail' },
            },
          ],
        },
      ],
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
