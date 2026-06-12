import { useState } from 'react'
import { register } from './authService.js'
import { useNavigate, Link } from 'react-router-dom'
import Button from '../../components/common/Button'
import Input from '../../components/common/Input'
import { USERNAME_RE } from '../../utils/username'

// RegisterForm is the sign up page. It validates input and creates the account.
function RegisterForm() {
    const navigate = useNavigate()
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: '',
        dateOfBirth: '',
    })
    const [error, setError] = useState(null)
    const [loading, setLoading] = useState(false)

    const handleChange = (e) => {
        setFormData({ ...formData, [e.target.name]: e.target.value })
    }

    // parseError turns a server error message into a friendly text for the user.
    const parseError = (errMsg) => {
        if (errMsg.includes('Password') || errMsg.includes('password')) {
            return 'Password must be at least 8 characters and contain uppercase, lowercase and a number'
        } else if (errMsg.includes('Username') || errMsg.includes('username')) {
            return 'Username is required and must be valid'
        } else if (errMsg.includes('Email') || errMsg.includes('email')) {
            return 'Please enter a valid email address'
        } else if (errMsg.includes('older than 13')) {
            return 'You must be older than 13 to register'
        } else if (errMsg.includes('already')) {
            return 'Username or email already exists'
        } else if (errMsg.includes('date')) {
            return 'Please enter a valid date of birth (YYYY-MM-DD)'
        }
        return errMsg
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        // Check the username format before sending the request.
        if (!USERNAME_RE.test(formData.username)) {
            setError(
                'Username may only contain alphanumeric characters or single hyphens, ' +
                    'cannot begin or end with a hyphen, and must be 1-25 characters.'
            )
            return
        }
        setLoading(true)
        setError(null)
        try {
            await register(
                formData.username,
                formData.email,
                formData.password,
                formData.dateOfBirth
            )
            navigate('/login')
        } catch (err) {
            const errMsg = err.response?.data?.error || 'Something went wrong'
            setError(parseError(errMsg))
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="min-h-screen bg-white flex items-center justify-center">
            <div className="flex w-full max-w-4xl">
                <div className="hidden lg:flex flex-1 items-center justify-center">
                    <img src="/logo512.png" className="w-77 h-77 object-contain mb-1 rounded" />
                </div>
                <div className="flex-1 flex flex-col justify-center px-8 py-12">
                    <div className="lg:hidden mb-8">
                        <img src="/logo.png" className="w-10 h-10 object-cover mb-6 rounded" />
                    </div>
                    <h1 className="text-black text-3xl font-bold mb-8">Create your account</h1>

                    {error && (
                        <p className="text-red-500 text-sm mb-4 bg-red-50 p-3 rounded-lg">
                            {error}
                        </p>
                    )}

                    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
                        <Input
                            type="text"
                            name="username"
                            value={formData.username}
                            onChange={handleChange}
                            placeholder="Username"
                        />
                        <Input
                            type="email"
                            name="email"
                            value={formData.email}
                            onChange={handleChange}
                            placeholder="Email"
                        />
                        <Input
                            type="password"
                            name="password"
                            value={formData.password}
                            onChange={handleChange}
                            placeholder="Password"
                        />
                        <Input
                            type="date"
                            name="dateOfBirth"
                            value={formData.dateOfBirth}
                            onChange={handleChange}
                            label="Date of birth"
                        />

                        <p className="text-gray-500 text-xs">
                            By signing up, you agree to our{' '}
                            <Link
                                to="/terms"
                                className="text-blue-400 cursor-pointer hover:underline"
                            >
                                Terms of Service
                            </Link>{' '}
                            and{' '}
                            <Link
                                to="/privacy"
                                className="text-blue-400 cursor-pointer hover:underline"
                            >
                                Privacy Policy
                            </Link>
                        </p>

                        <Button
                            type="submit"
                            variant="primary"
                            size="lg"
                            fullWidth
                            loading={loading}
                        >
                            Sign up
                        </Button>

                        <p className="text-gray-500 text-sm text-center">
                            Already have an account?{' '}
                            <span
                                onClick={() => navigate('/login')}
                                className="text-blue-400 cursor-pointer hover:underline"
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
