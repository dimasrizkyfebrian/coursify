import axios, { isAxiosError } from 'axios'

// Create an Axios instance with the base URL
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
})

// Request interceptor to add the token to the Authorization header
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('authToken') // Retrieve the token from local storage
    if (token) {
      config.headers.Authorization = `Bearer ${token}` // Add the token to the Authorization header
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  },
)

export default api
export { isAxiosError }
