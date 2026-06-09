import { useContext } from 'react'
import { AuthContext } from '../context/AuthProvider'

// useAuth returns the auth context. It throws when used outside AuthProvider.
export function useAuth() {
    const context = useContext(AuthContext)

    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider')
    }

    return context
}
