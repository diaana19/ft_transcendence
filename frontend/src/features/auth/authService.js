import api from '../../services/axiosInstance'

// register creates a new account.
export async function register(username, email, password, dateOfBirth) {
    const response = await api.post('/api/auth/register', {
        username,
        email,
        password,
        dateOfBirth,
    })
    return response.data
}
// login signs in using an email or username as identifier.
export async function login(identifier, password) {
    const response = await api.post('/api/auth/login', {
        email: identifier,
        username: identifier,
        password,
    })
    return response.data
}

// forgotPassword asks the server to send a password reset email.
export async function forgotPassword(email) {
    const response = await api.post('/api/auth/forgot-password', {
        email,
    })
    return response.data
}

// setupTwoFA starts 2FA setup and returns the secret/QR code.
export async function setupTwoFA() {
    const response = await api.post('/api/2fa/setup')
    return response.data
}

// enableTwoFA confirms the code and turns 2FA on.
export async function enableTwoFA(code) {
    const response = await api.post('/api/2fa/enable', {
        code,
    })
    return response.data
}

// disableTwoFA turns 2FA off after checking the code.
export async function disableTwoFA(code) {
    const response = await api.post('/api/2fa/disable', { code })
    return response.data
}

// verifyTwoFA finishes login by checking the 2FA code with the pending token.
export async function verifyTwoFA(pendingToken, code) {
    const response = await api.post('/api/auth/2fa/verify', {
        pending_token: pendingToken,
        code,
    })
    return response.data
}
