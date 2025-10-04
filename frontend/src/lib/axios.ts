import axios, { isAxiosError } from 'axios'

// Create an Axios instance with the base URL
const api = axios.create({
  baseURL: 'http://localhost:8080/api', // Replace with your API base URL
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
