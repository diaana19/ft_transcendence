import { useState } from 'react'
import { QRCodeSVG } from 'qrcode.react'
import { setupTwoFA, enableTwoFA } from './authService'
import { useAuth } from '../../hooks/useAuth'

// TwoFASetup shows the QR code and enables 2FA after the user confirms a code.
function TwoFASetup() {
    const { updateUser } = useAuth()
    const [qrCode, setQrCode] = useState(null)
    const [secret, setSecret] = useState('')
    const [code, setCode] = useState('')
    const [loading, setLoading] = useState(false)
    const [success, setSuccess] = useState(false)
    const [error, setError] = useState(null)

    // handleSetup asks the server for the secret and QR code.
    const handleSetup = async () => {
        try {
            setLoading(true)
            setError(null)
            const data = await setupTwoFA()
            setQrCode(data.qr_code_url)
            setSecret(data.secret)
        } catch (err) {
            setError(err.response?.data?.error || 'Failed to setup 2FA')
        } finally {
            setLoading(false)
        }
    }

    // handleEnable confirms the typed code and turns 2FA on.
    const handleEnable = async () => {
        try {
            setLoading(true)
            setError(null)
            await enableTwoFA(code)
            updateUser({ two_fa_enabled: true })
            setSuccess(true)
        } catch (err) {
            setError(err.response?.data?.error || 'Invalid code')
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="w-full mx-auto">
            <div className="sticky top-0 bg-white border-b border-gray-200 px-4 py-3 z-10">
                <h1 className="text-xl font-bold">Two-Factor Authentication</h1>
            </div>

            <div className="px-4 py-6">
                {!qrCode && !success && (
                    <div className="border border-gray-200 rounded-2xl p-6">
                        <h2 className="text-base font-bold text-gray-900 mb-1">Enable 2FA</h2>
                        <p className="text-sm text-gray-500 mb-4">
                            Add an extra layer of security to your account using Google
                            Authenticator or any TOTP app.
                        </p>
                        {error && <p className="text-red-500 text-sm mb-4">{error}</p>}
                        <button
                            onClick={handleSetup}
                            disabled={loading}
                            className="px-4 py-1.5 rounded-full text-sm font-bold transition-colors disabled:opacity-50"
                            style={{ background: '#534ab7', color: 'white' }}
                        >
                            {loading ? 'Loading...' : 'Setup 2FA'}
                        </button>
                    </div>
                )}

                {qrCode && !success && (
                    <div className="border border-gray-200 rounded-2xl p-6 space-y-4">
                        <div>
                            <p className="text-sm font-medium text-gray-900 mb-3">
                                Scan this QR code with Google Authenticator
                            </p>
                            <div className="flex justify-center p-4 bg-gray-50 rounded-xl">
                                <QRCodeSVG value={qrCode} size={180} includeMargin={true} />
                            </div>
                        </div>

                        <div>
                            <p className="text-xs text-gray-400 mb-1">Secret key</p>
                            <code className="block bg-gray-50 border border-gray-200 rounded-lg px-3 py-2 text-sm text-gray-700 break-all">
                                {secret}
                            </code>
                        </div>

                        <input
                            type="text"
                            placeholder="Enter 6-digit code"
                            value={code}
                            onChange={(e) => setCode(e.target.value)}
                            className="w-full border border-gray-200 rounded-lg px-3 py-2 text-sm focus:outline-none"
                            style={{ borderColor: '#ede8fd' }}
                        />

                        {error && <p className="text-red-500 text-sm">{error}</p>}

                        <button
                            onClick={handleEnable}
                            disabled={loading || code.length !== 6}
                            className="w-full py-2 rounded-full text-sm font-bold transition-colors disabled:opacity-50"
                            style={{ background: '#534ab7', color: 'white' }}
                        >
                            {loading ? 'Verifying...' : 'Verify & Enable'}
                        </button>
                    </div>
                )}

                {success && (
                    <div className="border border-gray-200 rounded-2xl p-6 text-center">
                        <p className="text-2xl mb-2">✅</p>
                        <p className="text-base font-bold" style={{ color: '#534ab7' }}>
                            2FA enabled successfully!
                        </p>
                        <p className="text-sm text-gray-500 mt-1">
                            Your account is now protected with two-factor authentication.
                        </p>
                    </div>
                )}
            </div>
        </div>
    )
}

export default TwoFASetup
