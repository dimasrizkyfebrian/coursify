import { ref } from 'vue'
import { jwtDecode } from 'jwt-decode'
import api from '@/lib/axios'

interface BreadcrumbItem {
  text: string
  to: string
}

const pendingUserCount = ref(0)
const breadcrumbs = ref<BreadcrumbItem[]>([])

// Reactive state
const user = ref({
  id: null as string | null,
  fullName: '' as string,
  email: '' as string,
  role: null as string | null,
  isLoggedIn: false,
})

async function fetchPendingUserCount() {
  if (user.value.role !== 'admin') return
  try {
    const response = await api.get('/admin/users/pending/count')
    pendingUserCount.value = response.data.count
  } catch (error) {
    console.error('Failed to fetch pending user count', error)
  }
}

function setBreadcrumbs(crumbs: BreadcrumbItem[]) {
  breadcrumbs.value = crumbs
}

// Functions fetchUserProfile
async function fetchUserProfile() {
  if (!user.value.isLoggedIn) return
  try {
    const response = await api.get('/profile')
    user.value.fullName = response.data.full_name
    user.value.email = response.data.email
    fetchPendingUserCount()
  } catch (error) {
    console.error('Error fetching user profile:', error)
    logoutUser()
  }
}

// Functions setUserFromToken
function setUserFromToken(token: string) {
  try {
    const decoded: { user_id: string; role: string } = jwtDecode(token)
    user.value.id = decoded.user_id
    user.value.role = decoded.role
    user.value.isLoggedIn = true
    fetchUserProfile()
  } catch (error) {
    console.error('Invalid token:', error)
    logoutUser()
  }
}

// Function checkUserStatus
function checkUserStatus() {
  const token = localStorage.getItem('authToken')
  if (token) {
    setUserFromToken(token)
  }
}

// Function logoutUser
function logoutUser() {
  localStorage.removeItem('authToken')
  user.value.id = null
  user.value.role = null
  user.value.isLoggedIn = false
}

// Export function
export function useUserStore() {
  return {
    user,
    pendingUserCount,
    fetchPendingUserCount,
    fetchUserProfile,
    checkUserStatus,
    setUserFromToken,
    logoutUser,
    breadcrumbs,
    setBreadcrumbs,
  }
}
