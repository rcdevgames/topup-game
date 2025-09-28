import axios from 'axios'

// Runtime environment variable support for Docker deployment
const getApiUrl = () => {
  // Priority: Runtime window config > Vite build-time env > fallback
  return window.APP_CONFIG?.API_URL || 
         import.meta.env.VITE_API_URL || 
         'http://localhost:3000/api'
}

const api = axios.create({
  baseURL: getApiUrl(),
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const token = sessionStorage.getItem('accessToken')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config

    // If the error is 401 and we haven't retried yet
    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        // Try to refresh token
        const refreshToken = sessionStorage.getItem('refreshToken')
        const response = await axios.post(`${getApiUrl()}/auth/refresh_token`, {
          refresh_token: refreshToken
        })

        const { access_token } = response.data
        sessionStorage.setItem('accessToken', access_token)

        // Retry the original request with new token
        originalRequest.headers.Authorization = `Bearer ${access_token}`
        return api(originalRequest)
      } catch (err) {
        // If refresh token fails, logout user
        sessionStorage.removeItem('accessToken')
        sessionStorage.removeItem('refreshToken')
        window.location.href = '/login'
        return Promise.reject(error)
      }
    }

    return Promise.reject(error)
  }
)

export default api