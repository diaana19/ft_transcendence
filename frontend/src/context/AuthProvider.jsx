import { createContext, useState, useEffect } from 'react'
import api from '../services/axiosInstance'

export const AuthContext = createContext()

export function AuthProvider({ children }) {
    const [token, setToken] = useState(null)
    const [user, setUser] = useState(null)
    const [loading, setLoading] = useState(true)


    // clearLocalSession wipes the in-memory + localStorage state. Used by
    // recovery paths (token expired, malformed JWT) where calling the backend
    // would just bounce a 401 — there is no useful token to blacklist.
    const clearLocalSession = () => {
        localStorage.removeItem('token')
        setToken(null)
        setUser(null)
    }

    // logout is the user-facing exit. Best-effort POST to /api/auth/logout
    // so the backend can blacklist the JWT in Redis AND clear the auth_token
    // cookie (the latter is critical for the OAuth path — without it the
    // cookie survives, /api/auth/me succeeds on next mount, and the user
    // appears to be logged back in immediately). The local cleanup runs
    // unconditionally so a network failure still drops the SPA out.
    const logout = async () => {
        try {
            await api.post('/api/auth/logout')
        } catch {
            // best-effort: network or 401 (already invalid) — fall through
        }
        clearLocalSession()
    }

    const decodeToken = (token) => {
        // Backend now puts the user id in the standard `sub` claim. Old
        // tokens issued before the switch carried a custom `user_id` field;
        // honour both so a refresh during the rollout doesn't kick users out.
        const payload = JSON.parse(atob(token.split('.')[1]))
        return { userId: payload.sub ?? payload.user_id, exp: payload.exp }
    }

    const isExpired = (exp) => exp * 1000 < Date.now()

    useEffect(() => {
        const storedToken = localStorage.getItem('token')

        if (storedToken) {
            try {
                const payload = decodeToken(storedToken)

                if (isExpired(payload.exp)) {
                    clearLocalSession()
                } else {
                    setToken(storedToken)
                    setUser({ userId: payload.userId,
                        username: payload.username
                     })
                }
                setLoading(false)
                return
            } catch {
                clearLocalSession()
            }
        }

        // No (valid) localStorage token. The OAuth callback only sets an
        // HttpOnly auth_token cookie, so we probe /api/auth/me to recover
        // the session — axiosInstance has withCredentials: true so the
        // cookie travels with this request.
        let cancelled = false
        api.get('/api/auth/me')
            .then((res) => {
                if (cancelled) return
                const u = res.data?.user
                if (u) {
                    setUser({
                        userId: u.id,
                        username: u.username,
                        email: u.email,
                        avatar: u.avatar,
                    })
                }
            })
            .catch(() => { /* not logged in — fall through */ })
            .finally(() => { if (!cancelled) setLoading(false) })

        return () => { cancelled = true }
    }, [])

    const loginUser = (data) => {
        try {
            if (!data?.token) {
            throw new Error('No token provided')
            }

            localStorage.setItem('token', data.token)
            const payload = decodeToken(data.token)
            const fullUser = ({ userId: payload.userId, 
                username: data.user?.username, 
                email: data.user?.email, 
                avatar: data.user?.avatar })
            
            setToken(data.token)
            setUser(fullUser)
            console.log("FULL USER:", fullUser)

        } catch (err) {
            console.error('Login failed:', err.message)
            clearLocalSession()
        }
    }

    return (
        <AuthContext.Provider value={{ token, 
        user, loginUser, logout, loading }}>
            {children}
        </AuthContext.Provider>
    )
}
