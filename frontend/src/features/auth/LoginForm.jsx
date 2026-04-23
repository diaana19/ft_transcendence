import { useState } from 'react'
import { login } from './authService.js'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'

function LoginForm() {
    const navigate = useNavigate()
    const { loginUser } = useAuth()
    const [formData, setFormData] = useState({
        email: '',
        password: ''
})

const [error, setError] = useState(null)
const [loading, setLoading] = useState(false)

const handleChange = (e) => {
    setFormData({
        ...formData,
        [e.target.name]: e.target.value
    })
}

const handleSubmit = async (e) => {
    e.preventDefault()

    console.log("SUBMIT FIRED")

    setLoading(true)
    setError(null)

    try {
        console.log("calling login...")

        const data = await login(formData.email, formData.password)

        console.log("RESPONSE:", data)

        loginUser(data.token)
        navigate('/')
    } catch (err) {
        console.log("ERROR:", err)
        setError(err.response?.data?.error || 'Something went wrong')
    } finally {
        setLoading(false)
    }
}

return (
   <div className="min-h-screen bg-white flex items-center justify-center">
      <div className="flex flex-col items-center w-full max-w-sm px-8 py-12">

        {/* Twitter Bird azul */}
        <svg viewBox="0 0 24 24" className="w-10 h-10 fill-blue-400 mb-6">
          <path d="M23.953 4.57a10 10 0 01-2.825.775 4.958 4.958 0 002.163-2.723c-.951.555-2.005.959-3.127 1.184a4.92 4.92 0 00-8.384 4.482C7.69 8.095 4.067 6.13 1.64 3.162a4.822 4.822 0 00-.666 2.475c0 1.71.87 3.213 2.188 4.096a4.904 4.904 0 01-2.228-.616v.06a4.923 4.923 0 003.946 4.827 4.996 4.996 0 01-2.212.085 4.936 4.936 0 004.604 3.417 9.867 9.867 0 01-6.102 2.105c-.39 0-.779-.023-1.17-.067a13.995 13.995 0 007.557 2.209c9.053 0 13.998-7.496 13.998-13.985 0-.21 0-.42-.015-.63A9.935 9.935 0 0024 4.59z"/>
        </svg>

        <h1 className="text-black text-2xl font-bold mb-8">
          Log in to Twitter
        </h1>

        {error && (
          <p className="text-red-500 text-sm mb-4 w-full bg-red-50 p-3 rounded-lg">
            {error}
          </p>
        )}

        <form onSubmit={handleSubmit} className="flex flex-col gap-4 w-full">

          <input
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            placeholder="Phone, email or username"
            className="w-full border border-gray-300 rounded px-4 py-3 text-black placeholder-gray-500 focus:outline-none focus:border-blue-400"
          />

          <input
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            placeholder="Password"
            className="w-full border border-gray-300 rounded px-4 py-3 text-black placeholder-gray-500 focus:outline-none focus:border-blue-400"
          />

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-blue-400 hover:bg-blue-500 text-white font-bold py-3 rounded-full transition-colors disabled:opacity-50 mt-2"
          >
            {loading ? 'Loading...' : 'Log in'}
          </button>

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
              Sign up for Twitter
            </span>
          </div>

        </form>
      </div>
    </div>
  )
}

export default LoginForm
