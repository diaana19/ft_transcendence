import { useState } from 'react'
import { register } from './authService.js'
import { useNavigate } from 'react-router-dom'

function RegisterForm() {
    const navigate = useNavigate()
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: '',
        dateOfBirth: ''
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
    setLoading(true)
    setError(null)

    try {
        await register(
            formData.username,
            formData.email,
            formData.password,
            formData.dateOfBirth
        )
        navigate ('/')
    } catch (err) {
        setError(err.response?.data?.error || 'Something went wrong')
    } finally {
        setLoading(false)
    }
}

return (
   <div className="min-h-screen bg-white flex items-center justify-center">
      <div className="flex w-full max-w-4xl">

        <div className="hidden lg:flex flex-1 items-center justify-center">
           <img 
                src="/logo.png" 
                className="w-77 h-77 object-contain mb-1 rounded"
              />
        </div>

        <div className="flex-1 flex flex-col justify-center px-8 py-12">

          <div className="lg:hidden mb-8">
              <img 
                src="/logo.png" 
                className="w-10 h-10 object-cover mb-6 rounded"
              />
          </div>

          <h1 className="text-black text-3xl font-bold mb-8">
            Create your account
          </h1>

          {error && (
            <p className="text-red-500 text-sm mb-4 bg-red-50 p-3 rounded-lg">
              {error}
            </p>
          )}

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">

            <input
              type="text"
              name="username"
              value={formData.username}
              onChange={handleChange}
              placeholder="Username"
              className="w-full border border-gray-300 rounded px-4 py-3 text-black placeholder-gray-500 focus:outline-none focus:border-blue-400"
            />

            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="Email"
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

            <div className="flex flex-col gap-1">
              <label className="text-gray-500 text-sm">Date of birth</label>
              <input
                type="date"
                name="dateOfBirth"
                value={formData.dateOfBirth}
                onChange={handleChange}
                className="w-full border border-gray-300 rounded px-4 py-3 text-black focus:outline-none focus:border-blue-400"
              />
            </div>

            <p className="text-gray-500 text-xs">
              By signing up, you agree to our{' '}
              <span className="text-indigo-400 cursor-pointer hover:underline">Terms of Service</span>
              {' '}and{' '}
              <span className="text-indigo-400 cursor-pointer hover:underline">Privacy Policy</span>
            </p>

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-indigo-200 hover:bg-[#d4dcff] text-white font-bold py-3 rounded-full transition-colors disabled:opacity-50"
            >
              {loading ? 'Loading...' : 'Sign up'}
            </button>

            <p className="text-gray-500 text-sm text-center">
              Already have an account?{' '}
              <span
                onClick={() => navigate('/login')}
                className="text-indigo-400 cursor-pointer hover:underline"
              >
                Log in
              </span>
            </p>

          </form>
        </div>
      </div>
    </div>
  )
}

export default RegisterForm