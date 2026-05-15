import api from '../../services/axiosInstance'

export async function register(username, email, password, dateOfBirth) {
    const response = await api.post('/api/auth/register', {
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

export async function setupTwoFA() {
    const response = await api.post('/api/2fa/setup')
    return response.data
}

export async function enableTwoFA(code) {
    const response = await api.post('/api/2fa/enable', {
        code
    })
    return response.data
}

export async function disableTwoFA() {
    const response = await api.post('/api/2fa/disable')
    return response.data
}

export async function verifyTwoFA(code) {
    const response = await api.post('/api/2fa/verify', {
        code
    })
    return response.data
}