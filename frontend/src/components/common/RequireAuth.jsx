import { Navigate } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'

export default function RequireAuth({ children }) {
  const { token, loading } = useAuth()
  if (loading) return null
  if (!token) return <Navigate to="/login" replace />
  return children
}