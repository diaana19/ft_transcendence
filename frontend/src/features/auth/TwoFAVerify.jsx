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
        <div className="min-h-screen flex items-center justify-center bg-transparent">
            <form
                onSubmit={handleSubmit}
                className="w-full max-w-sm p-8 border rounded-xl space-y-4"
                style={{ background: 'white', border: '0.5px solid #ede8fd' }}
            >
                <h1 className="text-2xl font-bold" style={{ color: '#2c2c2a' }}>
                    Two-factor verification
                </h1>
                <p className="text-sm text-gray-600" style={{ color: '#b4b2a9' }}>
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
                    className="w-full p-2 rounded-xl text-center tracking-widest text-lg focus:outline-none transition-colors"
                    style={{ border: '0.5px solid #ede8fd', color: '#2c2c2a' }}
                />

                <button
                    type="submit"
                    disabled={loading || code.length !== 6}
                    className="w-full py-2 rounded-full text-sm font-semibold transition-colors disabled:opacity-50"
                    style={{ background: '#534ab7', color: 'white', border: 'none' }}
                >
                    {loading ? 'Verifying…' : 'Verify'}
                </button>

                <button
                    type="button"
                    onClick={() => navigate('/login', { replace: true })}
                    className="w-full text-sm transition-colors"
                    style={{ color: '#afa9ec', background: 'transparent', border: 'none' }}
                >
                    Back to login
                </button>

                {error && (
                    <div className="text-center text-sm" style={{ color: '#d4537e' }}>
                        {error}
                    </div>
                )}
            </form>
        </div>
    )
}

export default TwoFAVerify
