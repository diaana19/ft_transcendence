import { useState } from 'react'
import { register } from './authService.js'
import { useNavigate } from 'react-router-dom'

function RegisterForm() {
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: '',
        dateOfBirth: ''
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
        await register( //await, espera rpta de go
            formData.username,
            formData.email,
            formData.password,
            formData.dateOfBirth
        )
        navigate ('/') //si go esta ok, redirige a home
    } catch (err) {
        setError(err.response?.data?.error || 'Something went wrong')
    } finally {
        setLoading(false)
    }
}

return ( //lo q ve el usuario
     <div>
      <h1>Your account</h1>

      {error && <p>{error}</p>}

      <form onSubmit={handleSubmit}>
        <div>
          <label>Username</label>
          <input type="text" name="username" value={formData.username} onChange={handleChange} />
        </div>
        <div>
          <label>Email</label>
          <input type="email" name="email" value={formData.email} onChange={handleChange} />
        </div>
        <div>
          <label>Password</label>
          <input type="password" name="password" value={formData.password} onChange={handleChange} />
        </div>
        <div>
          <label>Date of birth</label>
          <input type="date" name="dateOfBirth" value={formData.dateOfBirth} onChange={handleChange} />
        </div>
        <button type="button" onClick={() => navigate('/')}>Back</button>
        <button type="submit" disabled={loading}>
          {loading ? 'Loading...' : 'Next'}
        </button>
      </form>
    </div>
  )
}

export default RegisterForm




