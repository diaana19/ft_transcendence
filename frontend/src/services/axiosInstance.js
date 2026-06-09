import axios from 'axios'

// api is the shared axios instance used by all services.
const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL ?? '',
    withCredentials: true,
})

// Add the auth token from localStorage on every request.
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token')
    if (token) {
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

api.interceptors.response.use(
    (response) => response,
    (error) => {
        const status = error.response?.status
        const url = error.config?.url || ''
        const isAuthEndpoint = url.includes('/api/auth/')
        // On 401 (not from auth endpoints) clear the session and warn the app.
        if (status === 401 && !isAuthEndpoint) {
            localStorage.removeItem('token')
            window.dispatchEvent(new CustomEvent('auth:expired'))
        }
        return Promise.reject(error)
    }
)

export default api
