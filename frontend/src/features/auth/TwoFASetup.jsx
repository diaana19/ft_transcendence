import { useState } from 'react'
import { setupTwoFA, enableTwoFA } from './authService'

function TwoFASetup() {

    const [qrCode, setQrCode] = useState(null)
    const [secret, setSecret] = useState('')
    const [code, setCode] = useState('')
    const [loading, setLoading] = useState(false)
    const [success, setSuccess] = useState(false)
    const [error, setError] = useState(null)

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

    const handleEnable = async () => {
        try {
            setLoading(true)
            setError(null)

            await enableTwoFA(code)

            setSuccess(true)

        } catch (err) {
            setError(err.response?.data?.error || 'Invalid code')
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="max-w-md mx-auto mt-10 p-6 border rounded-xl">

            <h1 className="text-2xl font-bold mb-4">
                Two-Factor Authentication
            </h1>

            {!qrCode && (
                <button
                    onClick={handleSetup}
                    className="bg-blue-500 text-white px-4 py-2 rounded"
                >
                    Enable 2FA
                </button>
            )}

            {qrCode && !success && (
                <div className="space-y-4">

                    <div>
                        <p className="mb-2">
                            Scan this QR with Google Authenticator
                        </p>

                        <img
                            src={`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(qrCode)}`}
                            alt="QR Code"
                        />
                    </div>

                    <div>
                        <p className="text-sm text-gray-500">
                            Secret:
                        </p>

                        <code className="bg-gray-100 p-2 rounded block">
                            {secret}
                        </code>
                    </div>

                    <input
                        type="text"
                        placeholder="Enter 6-digit code"
                        value={code}
                        onChange={(e) => setCode(e.target.value)}
                        className="w-full border p-2 rounded"
                    />

                    <button
                        onClick={handleEnable}
                        disabled={loading}
                        className="w-full bg-green-500 text-white py-2 rounded"
                    >
                        Verify & Enable
                    </button>
                </div>
            )}

            {success && (
                <div className="text-green-600 font-semibold">
                    2FA enabled successfully
                </div>
            )}

            {error && (
                <div className="text-red-500 mt-4">
                    {error}
                </div>
            )}
        </div>
    )
}

export default TwoFASetup