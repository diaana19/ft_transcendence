import { useState } from 'react'
import { useAuth } from '../../hooks/useAuth'
import axiosInstance from '../../services/axiosInstance'
import { PhotoIcon } from '@heroicons/react/24/outline'

function CommentForm({ postId, onCommentAdded }) {
    const { user } = useAuth()
    const [content, setContent] = useState('')
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState(null)
    const [file, setFile] = useState(null)

    const handleSubmit = async () => {
        if (!content.trim()) return
        setLoading(true)
        setError(null)
        try {
            const formData = new FormData()
            formData.append('content', content)
            if (file) formData.append('file', file)
            const res = await axiosInstance.post(`/api/posts/${postId}/comments`, formData)
            onCommentAdded(res.data)
            setContent('')
            setFile(null)
        } catch (err) {
            setError(err.response?.data?.error || 'Something went wrong')
        } finally {
            setLoading(false)
        }
    }
    const handleFileChange = (e) => {
        const selectedFile = e.target.files[0]
        if (!selectedFile) return
        if (!selectedFile.type.startsWith('image/') && !selectedFile.type.startsWith('video/')) {
            setError('Only image and video files are allowed')
            setFile(null)
            return
        }
        if (selectedFile.size > 25 * 1024 * 1024) {
            setError('File too large. Maximum size is 25MB')
            setFile(null)
            return
        }
        setFile(selectedFile)
        console.log('Selected file:', selectedFile.name, 'size:', selectedFile.size)
        setError(null)
    }

    return (
        <div className="flex gap-3 px-4 py-3" style={{ borderBottom: '0.5px solid #ede8fd' }}>
            <div
                className="w-8 h-8 rounded-full overflow-hidden flex items-center justify-center font-semibold flex-shrink-0"
                style={{ background: '#ede8fd', color: '#534ab7', fontSize: '11px' }}
            >
                {user?.avatar ? (
                    <img
                        src={user.avatar}
                        alt={user.username}
                        className="w-full h-full object-cover"
                    />
                ) : (
                    <span>{user?.username?.charAt(0).toUpperCase() || 'U'}</span>
                )}
            </div>

            <div
                className="flex-1 rounded-xl px-3 py-2"
                style={{ background: 'white', border: '0.5px solid #ede8fd' }}
            >
                <textarea
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    placeholder="Post your reply..."
                    maxLength={280}
                    rows={2}
                    className="w-full focus:outline-none resize-none text-sm bg-transparent placeholder-[#b4b2a9]"
                    style={{ color: '#2c2c2a' }}
                />
                {file && (
                    <div className="relative mt-2 mb-1">
                        {file.type.startsWith('video/') ? (
                            <video
                                src={URL.createObjectURL(file)}
                                controls
                                className="rounded-xl max-h-48 w-full object-cover"
                            />
                        ) : (
                            <img
                                src={URL.createObjectURL(file)}
                                alt="preview"
                                className="rounded-xl max-h-48 w-full object-cover"
                            />
                        )}
                        <button
                            onClick={() => setFile(null)}
                            className="absolute top-2 right-2 w-5 h-5 rounded-full text-xs text-white flex items-center justify-center"
                            style={{ background: 'rgba(83, 74, 183, 0.7)' }}
                        >
                            ✕
                        </button>
                    </div>
                )}
                {error && (
                    <p className="text-xs mb-1" style={{ color: '#d4537e' }}>
                        {error}
                    </p>
                )}

                <div className="flex justify-between items-center mt-1">
                    <span className="text-xs" style={{ color: '#b4b2a9' }}>
                        {content.length}/280
                    </span>
                    <div className="flex gap-3 text-blue-400">
                        <label className="cursor-pointer hover:text-blue-500 flex items-center gap-2 text-gray-500 transition">
                            <PhotoIcon className="w-5 h-5" />
                            <input
                                type="file"
                                accept="image/*,video/*"
                                onChange={handleFileChange}
                                className="hidden"
                            />
                        </label>
                    </div>
                    <button
                        onClick={handleSubmit}
                        disabled={loading || !content.trim()}
                        className="text-sm font-semibold px-4 py-1 rounded-full transition-colors disabled:opacity-50 hover:bg-[#3c3489]"
                        style={{ background: '#534ab7', color: 'white', border: 'none' }}
                    >
                        {loading ? 'Replying...' : 'Reply'}
                    </button>
                </div>
            </div>
        </div>
    )
}

export default CommentForm
