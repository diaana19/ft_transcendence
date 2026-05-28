import { Navigate } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'

// Check `user`, not `token`. The OAuth path leaves token null because the JWT
// is in an HttpOnly cookie that JS cannot read; AuthProvider populates `user`
// via the /api/auth/me probe on mount. Gating on `user` covers both auth
// modes (localStorage JWT + cookie-only OAuth).
export default function RequireAuth({ children }) {
  const { user, loading } = useAuth()
  if (loading) return null
  if (!user) return <Navigate to="/login" replace />
  return children
}