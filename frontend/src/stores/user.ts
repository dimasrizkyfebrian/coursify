import { ref } from 'vue'
import { jwtDecode } from 'jwt-decode'

const user = ref({
  id: null as string | null,
  role: null as string | null,
  isLoggedIn: false,
})

function setUserFromToken(token: string) {
  try {
    const decoded: { user_id: string; role: string } = jwtDecode(token)
    user.value.id = decoded.user_id
    user.value.role = decoded.role
    user.value.isLoggedIn = true
  } catch (error) {
    console.error('Invalid token:', error)
    logoutUser()
  }
}

function checkUserStatus() {
  const token = localStorage.getItem('authToken')
  if (token) {
    setUserFromToken(token)
  }
}

function logoutUser() {
  localStorage.removeItem('authToken')
  user.value.id = null
  user.value.role = null
  user.value.isLoggedIn = false
}

export function useUserStore() {
  return {
    user,
    checkUserStatus,
    logoutUser,
    setUserFromToken,
  }
}
