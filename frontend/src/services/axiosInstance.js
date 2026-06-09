import axios from 'axios'

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL ?? '',
    withCredentials: true,
})

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
        if (status === 401 && !isAuthEndpoint) {
            localStorage.removeItem('token')
            window.dispatchEvent(new CustomEvent('auth:expired'))
        }
        return Promise.reject(error)
    }
)

export default api
