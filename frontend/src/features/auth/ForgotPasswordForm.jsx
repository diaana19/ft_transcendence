/*
** File: ForgotPasswordForm.jsx
** Description: Forgot-password form component
** Responsibilities:
** - Render an email field to request a password reset
** - Submit the request to the backend
** - Show a generic confirmation (the API never reveals if the email exists)
** - Link back to login
*/

import { useState } from 'react'
import { forgotPassword } from './authService.js'
import { useNavigate } from 'react-router-dom'
import Button from '../../components/common/Button'
import Input from '../../components/common/Input'

function ForgotPasswordForm() {
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [error, setError] = useState(null)
  const [submitted, setSubmitted] = useState(false)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError(null)
    try {
      await forgotPassword(email)
      setSubmitted(true)
    } catch (err) {
      setError(err.response?.data?.error || 'Something went wrong. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-white flex items-center justify-center">
      <div className="flex flex-col items-center w-full max-w-sm px-8 py-12">
        <img src="/logo.png" className="w-15 h-15 object-cover mb-6 rounded" />
        <h1 className="text-black text-2xl font-bold mb-8">Reset your password</h1>

        {submitted ? (
          <div className="w-full">
            <p className="text-green-700 text-sm mb-6 w-full bg-green-50 p-3 rounded-lg">
              If an account exists for that email, a new password has been sent to it.
              Check your inbox, then log in and change it from your settings.
            </p>
            <Button
              type="button"
              variant="primary"
              size="lg"
              fullWidth
              onClick={() => navigate('/login')}
            >
              Back to log in
            </Button>
          </div>
        ) : (
          <>
            {error && (
              <p className="text-red-500 text-sm mb-4 w-full bg-red-50 p-3 rounded-lg">
                {error}
              </p>
            )}

            <p className="text-gray-600 text-sm mb-6 w-full text-center">
              Enter your account email and we'll send you a new password.
            </p>

            <form onSubmit={handleSubmit} className="flex flex-col gap-4 w-full">
              <Input
                type="email"
                name="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Email"
                required
              />
              <Button type="submit" variant="primary" size="lg" fullWidth loading={loading}>
                Send new password
              </Button>

              <div className="flex justify-center mt-2">
                <span
                  onClick={() => navigate('/login')}
                  className="text-blue-400 text-sm cursor-pointer hover:underline"
                >
                  Back to log in
                </span>
              </div>
            </form>
          </>
        )}
      </div>
    </div>
  )
}

export default ForgotPasswordForm
