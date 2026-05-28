import { useEffect, useState } from 'react'
import { useLocation, useNavigate } from 'react-router-dom'
import { verifyTwoFA } from './authService'
import { useAuth } from '../../hooks/useAuth'

// Login-time second factor. Reached via redirect from LoginForm with
// { state: { pendingToken } }. Submits (pendingToken, 6-digit code) to
// /api/auth/2fa/verify; on success the backend issues the final JWT and
// the user is logged in.
function TwoFAVerify() {
    const location = useLocation()
    const navigate = useNavigate()
    const { loginUser } = useAuth()
    const pendingToken = location.state?.pendingToken

    const [code, setCode] = useState('')
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState(null)

    // No pending token in router state means the user landed here directly
    // (refresh, deep-link, bookmark). Send them back to the login screen.
    useEffect(() => {
        if (!pendingToken) {
            navigate('/login', { replace: true })
        }
    }, [pendingToken, navigate])

    const handleSubmit = async (e) => {
        e.preventDefault()
        if (code.length !== 6) {
            setError('Enter the 6-digit code from your authenticator app.')
            return
        }
        setLoading(true)
        setError(null)
        try {
            const data = await verifyTwoFA(pendingToken, code)
            loginUser(data)
            navigate('/', { replace: true })
        } catch (err) {
            setError(err.response?.data?.error || 'Invalid code')
            setCode('')
        } finally {
            setLoading(false)
        }
    }

    if (!pendingToken) return null

    return (
        <div className="min-h-screen flex items-center justify-center bg-white">
            <form
                onSubmit={handleSubmit}
                className="w-full max-w-sm p-8 border rounded-xl space-y-4"
            >
                <h1 className="text-2xl font-bold">Two-factor verification</h1>
                <p className="text-sm text-gray-600">
                    Enter the 6-digit code from your authenticator app.
                </p>

                <input
                    type="text"
                    inputMode="numeric"
                    pattern="[0-9]*"
                    autoComplete="one-time-code"
                    autoFocus
                    maxLength={6}
                    value={code}
                    onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                    placeholder="123456"
                    className="w-full border p-2 rounded text-center tracking-widest text-lg"
                />

                <button
                    type="submit"
                    disabled={loading || code.length !== 6}
                    className="w-full bg-blue-500 disabled:bg-gray-300 text-white py-2 rounded"
                >
                    {loading ? 'Verifying…' : 'Verify'}
                </button>

                <button
                    type="button"
                    onClick={() => navigate('/login', { replace: true })}
                    className="w-full text-sm text-gray-500 underline"
                >
                    Back to login
                </button>

                {error && <div className="text-red-500 text-sm">{error}</div>}
            </form>
        </div>
    )
}

export default TwoFAVerify
