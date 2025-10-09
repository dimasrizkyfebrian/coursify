import { ref } from 'vue'
import { jwtDecode } from 'jwt-decode'
import api from '@/lib/axios'

// Reactive state
const user = ref({
  id: null as string | null,
  fullName: '' as string,
  email: '' as string,
  role: null as string | null,
  isLoggedIn: false,
})

// Functions fetchUserProfile
async function fetchUserProfile() {
  if (!user.value.isLoggedIn) return
  try {
    const response = await api.get('/profile')
    user.value.fullName = response.data.full_name
    user.value.email = response.data.email
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
    checkUserStatus,
    logoutUser,
    setUserFromToken,
    fetchUserProfile,
  }
}
