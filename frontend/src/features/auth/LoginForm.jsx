/*
 ** File: LoginForm.jsx
 ** Description: Login form component
 ** Responsibilities:
 ** - Render login form with email and password fields
 ** - Handle form submission and authentication
 ** - Redirect to 2FA if required
 ** - Navigate to feed on success
 */

import { useState } from 'react'
import { login } from './authService.js'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'
import OAuthButton from './OAuthButton'
import Button from '../../components/common/Button'
import Input from '../../components/common/Input'

function LoginForm() {
    const navigate = useNavigate()
    const { loginUser } = useAuth()
    const [formData, setFormData] = useState({ email: '', password: '' })
    const [error, setError] = useState(null)
    const [loading, setLoading] = useState(false)

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value })
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setLoading(true)
        setError(null)
        try {
            const data = await login(formData.email, formData.password)
            if (data.needs_2fa) {
                navigate('/login/2fa', { state: { pendingToken: data.pending_token } })
                return
            }
            loginUser(data)
            navigate('/')
        } catch (err) {
            const errMsg = err.response?.data?.error || 'Something went wrong'
            if (errMsg.includes('credential') || errMsg.includes('invalid')) {
                setError('Incorrect username or password. Please try again.')
            } else {
                setError(errMsg)
            }
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="min-h-screen bg-white flex items-center justify-center">
            <div className="flex flex-col items-center w-full max-w-sm px-8 py-12">
                <img src="/logo.png" className="w-15 h-15 object-cover mb-6 rounded" />
                <h1 className="text-black text-2xl font-bold mb-8">Log in to Synk</h1>

                {error && (
                    <p className="text-red-500 text-sm mb-4 w-full bg-red-50 p-3 rounded-lg">
                        {error}
                    </p>
                )}

                <form onSubmit={handleSubmit} className="flex flex-col gap-4 w-full">
                    <Input
                        type="text"
                        name="email"
                        value={formData.email}
                        onChange={handleChange}
                        placeholder="Email or username"
                    />
                    <Input
                        type="password"
                        name="password"
                        value={formData.password}
                        onChange={handleChange}
                        placeholder="Password"
                    />
                    <Button type="submit" variant="primary" size="lg" fullWidth loading={loading}>
                        Log in
                    </Button>

                    <div className="flex items-center gap-3 my-2">
                        <div className="h-px bg-gray-300 flex-1"></div>
                        <span className="text-gray-500 text-sm">or</span>
                        <div className="h-px bg-gray-300 flex-1"></div>
                    </div>

                    <OAuthButton />

                    <div className="flex justify-between mt-2">
                        <span
                            onClick={() => navigate('/forgot-password')}
                            className="text-blue-400 text-sm cursor-pointer hover:underline"
                        >
                            Forgot password?
                        </span>
                        <span
                            onClick={() => navigate('/register')}
                            className="text-blue-400 text-sm cursor-pointer hover:underline"
                        >
                            Sign up for Synk
                        </span>
                    </div>
                </form>
            </div>
        </div>
    )
}

export default LoginForm
