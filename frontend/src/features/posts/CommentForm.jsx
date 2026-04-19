import { useState } from 'react'
import { useAuth } from '../../hooks/useAuth'
import axiosInstance from '../../services/axiosInstance'

function CommentForm({ postId, onCommentAdded }) {
  const { user } = useAuth()
  const [content, setContent]   = useState('')
  const [loading, setLoading]   = useState(false)
  const [error, setError]       = useState(null)

  const handleSubmit = async () => {
    if (!content.trim()) return
    setLoading(true)
    setError(null)
    try {
      const res = await axiosInstance.post(`/api/posts/${postId}/comments`, { content })
      onCommentAdded(res.data) 
      setContent('')
    } catch (err) {
      setError(err.response?.data?.error || 'Something went wrong')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex gap-3 px-4 py-3 border-b border-gray-200">
      {/* Avatar */}
      <div className="w-9 h-9 rounded-full bg-gray-300 flex-shrink-0" />

      <div className="flex-1">
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder="Post your reply"
          maxLength={280}
          rows={2}
          className="w-full text-black placeholder-gray-500 focus:outline-none resize-none text-sm"
        />

        {error && (
          <p className="text-red-500 text-xs mb-1">{error}</p>
        )}

        <div className="flex justify-between items-center mt-1">
          <span className="text-xs text-gray-400">{content.length}/280</span>
          <button
            onClick={handleSubmit}
            disabled={loading || !content.trim()}
            className="bg-blue-400 hover:bg-blue-500 text-white text-sm font-bold px-4 py-1.5 rounded-full disabled:opacity-50 transition-colors"
          >
            {loading ? 'Replying...' : 'Reply'}
          </button>
        </div>
      </div>
    </div>
  )
}

export default CommentForm