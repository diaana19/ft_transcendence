import { createContext, useState, useEffect } from 'react'

export const AuthContext = createContext()

export function AuthProvider({ children }) {
    const [token, setToken] = useState(null)
    const [user, setUser] = useState(null)
    const [loading, setLoading] = useState(true)


    const logout = () => {
        localStorage.removeItem('token')
        setToken(null)
        setUser(null)
    }

    const decodeToken = (token) => {
        const payload = JSON.parse(atob(token.split('.')[1]))
        return { userId: payload.userId, exp: payload.exp }
    }

    const isExpired = (exp) => exp * 1000 < Date.now()

    useEffect(() => {
        const storedToken = localStorage.getItem('token')

        if (storedToken) {
            try {
                const payload = decodeToken(storedToken)

                if (isExpired(payload.exp)) {
                    logout()
                } else {
                    setToken(storedToken)
                    setUser({ userId: payload.userId })
                }
            } catch {
                logout()
            }
        }

        setLoading(false)
    }, [])

    const loginUser = (jwt) => {
        try {
            if (!jwt) {
            throw new Error('No token provided')
            }
            const payload = decodeToken(jwt)
            localStorage.setItem('token', jwt)
            setToken(jwt)
            setUser({ userId: payload.userId })
        } catch (err) {
            console.error('Login failed:', err.message)
            logout()
        }
    }

    return (
        <AuthContext.Provider value={{ token, user, loginUser, logout, loading }}>
            {children}
        </AuthContext.Provider>
    )
}