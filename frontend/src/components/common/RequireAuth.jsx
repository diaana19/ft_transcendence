import { Navigate } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'

// RequireAuth blocks its children and sends to /login when there is no user.
export default function RequireAuth({ children }) {
    const { user, loading } = useAuth()
    // Wait while the session is still loading.
    if (loading) return null
    if (!user) return <Navigate to="/login" replace />
    return children
}
