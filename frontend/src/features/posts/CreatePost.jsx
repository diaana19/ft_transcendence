/*
** File: CreatePost.jsx
** Description: Component for creating new posts
** Responsibilities:
** - Render post creation form
** - Handle post submission
** - Show character count
** - Notify parent component when post is created
*/

import { useState } from 'react'
import { createPost } from './postService.js'
import { Image } from 'lucide-react'
import axiosInstance from '../../services/axiosInstance'
import { useAuth } from '../../hooks/useAuth'

function CreatePost({ onPostCreated }) {
  const { user } = useAuth()
  const [file, setFile] = useState(null)
  const [content, setContent] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0]
    if (!selectedFile) return
    if (!selectedFile.type.startsWith('image/')) {
      setError('Only image files are allowed')
      setFile(null)
      return
    }
    if (selectedFile.size > 20 * 1024 * 1024) {
      setError('File too large. Maximum size is 20MB')
      setFile(null)
      return
    }
    setFile(selectedFile)
    setError(null)
  }

  const handleSubmit = async () => {
    if (!content.trim()) return
    setLoading(true)
    setError(null)
    try {
      let mediaUrl = null
      if (file) {
        const formData = new FormData()
        formData.append('file', file)
        const res = await axiosInstance.post('/api/upload', formData)
        mediaUrl = res.data.url
      }
      const newPost = await createPost(content, mediaUrl)
      onPostCreated(newPost)
      setContent('')
      setFile(null)
    } catch (err) {
      setError(err.response?.data?.error || 'Something went wrong')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="p-4 mb-4" style={{
      background: 'white',
      borderRadius: '1.25rem',
      border: '1.5px solid #ede8fd',
      boxShadow: '0 4px 24px rgba(167, 139, 250, 0.08)'
    }}>
      <div className="flex gap-3">
        {/* Avatar */}
        <div className="w-10 h-10 rounded-full overflow-hidden flex-shrink-0 flex items-center justify-center font-bold text-sm"
          style={{ background: '#ede8fd', color: '#534ab7' }}
        >
          {user?.avatar
            ? <img src={user.avatar} alt="avatar" className="w-full h-full object-cover" />
            : user?.username?.[0]?.toUpperCase() || '?'
          }
        </div>

        <div className="flex-1">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            maxLength={280}
            placeholder="What's happening?"
            className="w-full text-sm resize-none focus:outline-none pb-2"
            style={{ color: '#2c2c2a', background: 'transparent' }}
            rows={4}
          />

          {file && (
            <div className="relative mb-3">
              {file.type.startsWith('video/') ? (
                <video
                  src={URL.createObjectURL(file)}
                  controls
                  className="rounded-xl max-h-64 w-full object-cover"
                />
              ) : (
                <img
                  src={URL.createObjectURL(file)}
                  alt="preview"
                  className="rounded-xl max-h-64 w-full object-cover"
                />
              )}
              <button
                onClick={() => setFile(null)}
                className="absolute top-2 right-2 w-6 h-6 rounded-full flex items-center justify-center text-xs text-white transition-colors"
                style={{ background: 'rgba(83, 74, 183, 0.7)' }}
              >
                ✕
              </button>
            </div>
          )}

          {error && (
            <p className="text-xs mb-2" style={{ color: '#d4537e' }}>{error}</p>
          )}

          <div className="flex items-center justify-between pt-3" style={{ borderTop: '1px solid #ede8fd' }}>
            <label className="cursor-pointer flex items-center gap-2 transition-colors"
              style={{ color: '#afa9ec' }}
            >
              <Image size={16} />
              <span className="text-xs font-medium">Image / Video</span>
              <input
                type="file"
                accept="image/*"
                onChange={handleFileChange}
                className="hidden"
              />
            </label>

            <div className="flex items-center gap-3">
              <span className="text-xs" style={{ color: content.length > 260 ? '#d4537e' : '#afa9ec' }}>
                {content.length}/280
              </span>
              <button
                onClick={handleSubmit}
                disabled={loading || !content.trim()}
                className="px-4 py-1.5 rounded-full text-xs font-bold transition-colors disabled:opacity-50"
                style={{ background: '#534ab7', color: 'white' }}
              >
                {loading ? 'Posting...' : 'Share Post'}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default CreatePost