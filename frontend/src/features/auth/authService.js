import axios from 'axios'

export async function register(username, email, password, dateOfBirth) {
    const response = await axios.post('/api/auth/register/', {
        username,
        email,
        password,
        dateOfBirth
    })
    return response.data
}