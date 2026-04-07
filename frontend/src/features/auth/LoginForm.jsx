import { useState } from 'react'
import { login } from './authService.js'
import { useNavigate } from 'react-router-dom'

function LoginForm() {
    const navigate = useNavigate()

    const [formData, setFormData] = useState({
        email: '',
        password: ''
})

const [error, setError] = useState(null)
const [loading, setLoading] = useState(false)

const handleChange = (e) => { //name = username... value = lo q escribio el usuario
    setFormData({
        ...formData,
        [e.target.name]: e.target.value //se actualizo solo lo q cambio
    })
}

const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
        const data = await login(formData.email, formData.password) //await, espera rpta de go
        localStorage.setItem('token', data.token) //guarda el JWT en localStorage(memoria del navegador)
        navigate ('/') //si go esta ok, redirige a home
    } catch (err) {
        setError(err.response?.data?.error || 'Something went wrong')
    } finally {
        setLoading(false)
    }
}

return (
    <div>
      <h1>Login</h1>

      {error && <p>{error}</p>}

      <form onSubmit={handleSubmit}>
        <div>
          <label>Email</label>
          <input
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
          />
        </div>
        <div>
          <label>Password</label>
          <input
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
          />
        </div>
        <button type="button" onClick={() => navigate('/register')}> 
          Register
        </button>
        <button type="submit" disabled={loading}>
          {loading ? 'Loading...' : 'Login'}
        </button>
      </form>
    </div>
  )
}

export default LoginForm