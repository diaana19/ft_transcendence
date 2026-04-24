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

        <img
          src="logo.png"
          className="w-15 h-15 object-cover mb-6 rounded"
        />

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
