import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'
import api from '../../services/axiosInstance'
import Modal from '../../components/common/Modal'
import { changePassword } from '../auth/authService'

// SettingsPage shows account, security, data export and account delete options.
export default function SettingsPage() {
    const { user, logout } = useAuth()
    const navigate = useNavigate()
    const [showDeleteModal, setShowDeleteModal] = useState(false)
    const [password, setPassword] = useState('')
    const [error, setError] = useState(null)
    const [loading, setLoading] = useState(false)

    const [showPasswordModal, setShowPasswordModal] = useState(false)
    const [currentPassword, setCurrentPassword] = useState('')
    const [newPassword, setNewPassword] = useState('')
    const [passwordError, setPasswordError] = useState(null)
    const [passwordDone, setPasswordDone] = useState(false)
    const [passwordLoading, setPasswordLoading] = useState(false)

    const closePasswordModal = () => {
        setShowPasswordModal(false)
        setCurrentPassword('')
        setNewPassword('')
        setPasswordError(null)
        setPasswordDone(false)
    }

    // handleChangePassword changes the password and shows a success message.
    const handleChangePassword = async () => {
        if (!currentPassword || !newPassword) {
            setPasswordError('Please fill both fields')
            return
        }
        setPasswordLoading(true)
        setPasswordError(null)
        try {
            await changePassword(currentPassword, newPassword)
            setPasswordDone(true)
            setCurrentPassword('')
            setNewPassword('')
        } catch (err) {
            setPasswordError(err.response?.data?.error || 'Something went wrong')
        } finally {
            setPasswordLoading(false)
        }
    }

    // handleDelete deletes the account after confirming the password, then logs out.
    const handleDelete = async () => {
        if (!password) {
            setError('Please enter your password')
            return
        }
        setLoading(true)
        setError(null)
        try {
            await api.delete(`/api/users/${user.userId}`, {
                data: { password },
            })
            await logout()
            navigate('/login')
        } catch (err) {
            setError(err.response?.data?.error || 'Something went wrong')
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="max-w-2xl mx-auto">
            <div className="sticky top-0 bg-white border-b border-gray-200 px-4 py-3 z-10">
                <h1 className="text-xl font-bold">Settings</h1>
            </div>

            <div className="px-4 py-6 border-b border-gray-200">
                <h2 className="text-base font-bold text-gray-900 mb-1">Account</h2>
                <p className="text-sm text-gray-500 mb-4">Manage your account information</p>

                <div className="flex flex-col gap-3">
                    <div className="flex items-center justify-between py-3 border-b border-gray-100">
                        <div>
                            <p className="text-sm font-medium text-gray-900">Username</p>
                            <p className="text-sm text-gray-500">@{user?.username}</p>
                        </div>
                    </div>
                    <div className="flex items-center justify-between py-3 border-b border-gray-100">
                        <div>
                            <p className="text-sm font-medium text-gray-900">Email</p>
                            <p className="text-sm text-gray-500">{user?.email}</p>
                        </div>
                    </div>
                    <div className="flex items-center justify-between py-3 border-b border-gray-100">
                        <div>
                            <p className="text-sm font-medium text-gray-900">Change password</p>
                            <p className="text-sm text-gray-500">Update your account password</p>
                        </div>
                        <button
                            onClick={() => setShowPasswordModal(true)}
                            className="text-sm text-blue-400 hover:underline font-medium"
                        >
                            Change
                        </button>
                    </div>
                </div>
            </div>

            <div className="px-4 py-6 border-b border-gray-200">
                <h2 className="text-base font-bold text-gray-900 mb-1">Security</h2>
                <p className="text-sm text-gray-500 mb-4">Manage your security settings</p>

                <div className="flex items-center justify-between py-3 border-b border-gray-100">
                    <div>
                        <p className="text-sm font-medium text-gray-900">
                            Two-factor authentication
                        </p>
                        <p className="text-sm text-gray-500">
                            {user?.two_fa_enabled ? 'Enabled ✓' : 'Disabled'}
                        </p>
                    </div>
                    <button
                        onClick={() => navigate('/2fa/setup')}
                        className="text-sm text-blue-400 hover:underline font-medium"
                    >
                        {user?.two_fa_enabled ? 'Manage' : 'Enable'}
                    </button>
                </div>
            </div>

            <div className="px-4 py-6 border-b border-gray-200">
                <h2 className="text-base font-bold text-gray-900 mb-1">Your data</h2>
                <p className="text-sm text-gray-500 mb-4">Manage your personal data</p>

                <div className="flex items-center justify-between py-3 border-b border-gray-100">
                    <div>
                        <p className="text-sm font-medium text-gray-900">Export your data</p>
                        <p className="text-sm text-gray-500">
                            Download all your data in JSON format
                        </p>
                    </div>
                    <button
                        onClick={async () => {
                            // Export the user data and download it as a JSON file.
                            const res = await api.get('/api/gdpr/export')
                            const blob = new Blob([JSON.stringify(res.data, null, 2)], {
                                type: 'application/json',
                            })
                            const url = URL.createObjectURL(blob)
                            const a = document.createElement('a')
                            a.href = url
                            a.download = 'my-data.json'
                            a.click()
                        }}
                        className="text-sm text-blue-400 hover:underline font-medium"
                    >
                        Download
                    </button>
                </div>
            </div>

            <div className="px-4 py-6">
                <h2 className="text-base font-bold text-red-500 mb-1">Danger zone</h2>
                <p className="text-sm text-gray-500 mb-4">
                    These actions are permanent and cannot be undone
                </p>

                <div className="flex items-center justify-between py-3">
                    <div>
                        <p className="text-sm font-medium text-gray-900">Delete account</p>
                        <p className="text-sm text-gray-500">
                            Permanently delete your account and all your data
                        </p>
                    </div>
                    <button
                        onClick={() => setShowDeleteModal(true)}
                        className="text-sm text-red-500 hover:underline font-medium border border-red-300 px-3 py-1.5 rounded-full hover:bg-red-50 transition-colors"
                    >
                        Delete
                    </button>
                </div>
            </div>

            <Modal
                isOpen={showPasswordModal}
                onClose={closePasswordModal}
                title="Change password"
                confirmText={
                    passwordDone ? 'Done' : passwordLoading ? 'Saving...' : 'Change password'
                }
                cancelText={passwordDone ? 'Close' : 'Cancel'}
                onConfirm={passwordDone ? closePasswordModal : handleChangePassword}
                size="sm"
            >
                {passwordDone ? (
                    <p className="text-gray-600 text-sm">
                        Your password was changed. We sent a confirmation to your email.
                    </p>
                ) : (
                    <>
                        <input
                            type="password"
                            value={currentPassword}
                            onChange={(e) => setCurrentPassword(e.target.value)}
                            placeholder="Current password"
                            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-3 focus:outline-none focus:border-blue-400"
                        />
                        <input
                            type="password"
                            value={newPassword}
                            onChange={(e) => setNewPassword(e.target.value)}
                            placeholder="New password"
                            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-blue-400"
                        />
                        {passwordError && (
                            <p className="text-red-500 text-xs mt-2">{passwordError}</p>
                        )}
                    </>
                )}
            </Modal>

            <Modal
                isOpen={showDeleteModal}
                onClose={() => {
                    setShowDeleteModal(false)
                    setPassword('')
                    setError(null)
                }}
                title="Delete account"
                confirmText={loading ? 'Deleting...' : 'Yes, delete'}
                cancelText="Cancel"
                onConfirm={handleDelete}
                confirmVariant="danger"
                size="sm"
            >
                <p className="text-gray-500 text-sm mb-4">
                    This action is permanent and cannot be undone. Please enter your password to
                    confirm.
                </p>
                <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Your password"
                    className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:border-red-400"
                />
                {error && <p className="text-red-500 text-xs mt-2">{error}</p>}
            </Modal>
            <div className="px-4 py-6 md:hidden">
                <button
                    onClick={async () => {
                        await logout()
                        navigate('/login')
                    }}
                    className="w-full py-2 rounded-full text-sm font-bold transition-colors"
                    style={{ background: '#fde8f0', color: '#d4537e', border: '1px solid #fde8f0' }}
                >
                    Log out
                </button>
            </div>
        </div>
    )
}
