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
export async function login(identifier, password) {
    const response = await api.post('/api/auth/login', {
        email: identifier,
        username: identifier,
        password
    })
    return response.data
}

export async function forgotPassword(email) {
    const response = await api.post('/api/auth/forgot-password', {
        email
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

export async function disableTwoFA(code) {
    const response = await api.post('/api/2fa/disable', { code })
    return response.data
}

export async function verifyTwoFA(pendingToken, code) {
    const response = await api.post('/api/auth/2fa/verify', {
        pending_token: pendingToken,
        code
    })
    return response.data
}