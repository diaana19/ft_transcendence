import { createContext, useState, useEffect } from 'react'
import api from '../services/axiosInstance'

export const AuthContext = createContext()

// AuthProvider keeps the logged user and token, and shares login/logout helpers.
export function AuthProvider({ children }) {
    const [token, setToken] = useState(null)
    const [user, setUser] = useState(null)
    const [loading, setLoading] = useState(true)
    const clearLocalSession = () => {
        localStorage.removeItem('token')
        setToken(null)
        setUser(null)
    }
    const logout = async () => {
        try {
            await api.post('/api/auth/logout')
        } catch {}
        clearLocalSession()
    }

    // decodeToken reads the JWT payload to get the user id, username and exp.
    const decodeToken = (token) => {
        const payload = JSON.parse(atob(token.split('.')[1]))
        return {
            userId: payload.sub ?? payload.user_id,
            username: payload.username,
            exp: payload.exp,
        }
    }

    const isExpired = (exp) => exp * 1000 < Date.now()

    // Listen the auth:expired event sent by the axios interceptor on 401.
    useEffect(() => {
        const onExpired = () => clearLocalSession()
        window.addEventListener('auth:expired', onExpired)
        return () => window.removeEventListener('auth:expired', onExpired)
    }, [])

    useEffect(() => {
        const storedToken = localStorage.getItem('token')

        if (storedToken) {
            try {
                const payload = decodeToken(storedToken)

                if (isExpired(payload.exp)) {
                    clearLocalSession()
                } else {
                    setToken(storedToken)
                    setUser({ userId: payload.userId, username: payload.username })
                }
            } catch {
                clearLocalSession()
            }
        }

        let cancelled = false
        api.get('/api/auth/session')
            .then((res) => {
                if (cancelled) return
                const u = res.data?.user
                if (u) {
                    setUser({
                        userId: u.id,
                        username: u.username,
                        email: u.email,
                        avatar: u.avatar,
                        two_fa_enabled: u.two_fa_enabled,
                    })
                } else {
                    clearLocalSession()
                }
            })
            .catch(() => {})
            .finally(() => {
                if (!cancelled) setLoading(false)
            })

        return () => {
            cancelled = true
        }
    }, [])

    // loginUser saves the token and sets the user state after a successful login.
    const loginUser = (data) => {
        try {
            if (!data?.token) {
                throw new Error('No token provided')
            }

            localStorage.setItem('token', data.token)
            const payload = decodeToken(data.token)
            const fullUser = {
                userId: payload.userId,
                username: data.user?.username,
                email: data.user?.email,
                avatar: data.user?.avatar,
                two_fa_enabled: data.user?.two_fa_enabled,
            }

            setToken(data.token)
            setUser(fullUser)
            console.info('FULL USER:', fullUser)
        } catch (err) {
            console.info('Login failed:', err.message)
            clearLocalSession()
        }
    }

    // updateUser merges a partial patch into the current user state.
    const updateUser = (patch) => {
        setUser((prev) => (prev ? { ...prev, ...patch } : prev))
    }

    return (
        <AuthContext.Provider value={{ token, user, loginUser, logout, loading, updateUser }}>
            {children}
        </AuthContext.Provider>
    )
}
