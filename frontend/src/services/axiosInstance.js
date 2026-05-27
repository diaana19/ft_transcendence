import axios from 'axios'

// Same-origin by default: in dev Vite proxies /api -> :8000, in production
// nginx proxies /api -> backend. Override with VITE_API_URL when the API is on
// a different origin.
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

export default api