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
        return { userId: payload.user_id, exp: payload.exp }
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
                    setUser({ userId: payload.userId,
                        username: payload.username
                     })
                }
            } catch {
                logout()
            }
        }

        setLoading(false)
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
            logout()
        }
    }

    return (
        <AuthContext.Provider value={{ token, 
        user, loginUser, logout, loading }}>
            {children}
        </AuthContext.Provider>
    )
}
