import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'
import { login } from '../../features/auth/authService'
import Modal from './Modal'
import Input from './Input'
import Button from './Button'

export default function LoginModal({ isOpen, onClose }) {
	const navigate = useNavigate()
	const { loginUser } = useAuth()
	const [formData, setFormData] = useState({ email: '', password: '' })
	const [error, setError] = useState(null)
	const [loading, setLoading] = useState(false)

	const handleChange = (e) => {
		setFormData({ ...formData, [e.target.name]: e.target.value })
	}

	const handleSubmit = async () => {
		if (!formData.email || !formData.password) return
		setLoading(true)
		setError(null)
		try {
			const data = await login(formData.email, formData.password)
			if (data.needs_2fa) {
				onClose()
				navigate('/login/2fa', { state: { pendingToken: data.pending_token } })
				return
			}
			loginUser(data)
			onClose()
		} catch (err) {
			setError(err.response?.data?.error || 'Something went wrong')
		} finally {
			setLoading(false)
		}
	}

	return (
		<Modal
			isOpen={isOpen}
			onClose={onClose}
			title="Log in to Synk"
			size="sm"
		>
			<div className="flex flex-col gap-4">
				{error && (
					<p className="text-red-500 text-sm bg-red-50 p-3 rounded-lg">{error}</p>
				)}
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
				<Button
					variant="primary"
					size="lg"
					fullWidth
					loading={loading}
					onClick={handleSubmit}
				>
					Log in
				</Button>
				<p className="text-center text-sm text-gray-500">
					Don't have an account?{' '}
					<span
						onClick={() => { onClose(); navigate('/register') }}
						className="text-blue-400 cursor-pointer hover:underline"
					>
						Sign up
					</span>
				</p>
			</div>
		</Modal>
	)
}