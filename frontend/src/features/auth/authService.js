import api from '../../services/axiosInstance'

export async function register(username, email, password, dateOfBirth) {
    const response = await axios.post('/api/auth/register', {
        username,
        email,
        password,
        dateOfBirth
    })
    return response.data
}
export async function login(email, password) {
    const response = await api.post('/api/auth/login', {
        email,
        password
    })
    return response.data
}



