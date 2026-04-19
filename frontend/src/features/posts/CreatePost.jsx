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

function CreatePost({ onPostCreated }) {
    const [content, setContent] = useState('')
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState(null)

const handleSubmit = async () => {
    if (!content.trim()) return
    setLoading(true)
    setError(null)
    try {
        const newPost = await createPost(content)
        onPostCreated(newPost)
        setContent('')
    } catch (err) {
        setError(err.response?.data?.error || 'Something went wrong')
    } finally {
        setLoading(false)
    }
}

return (
    <div className="border-b border-gray-200 p-4">
      <div className="flex gap-3">
        <div className="w-10 h-10 rounded-full bg-blue-400 flex items-center justify-center text-white font-bold">
          T
        </div>
        <div className="flex-1">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            maxLength={280}
            placeholder="What's happening?"
            className="w-full text-black placeholder-gray-500 text-lg resize-none focus:outline-none border-b border-gray-200 pb-3"
            rows={2}
          />
          {error && (
            <p className="text-red-500 text-sm mb-2">{error}</p>
          )}
          <div className="flex items-center justify-end gap-3 mt-2">
            <span className={`text-sm ${content.length > 260 ? 'text-red-500' : 'text-gray-500'}`}>
              {280 - content.length}
            </span>
            <button
              onClick={handleSubmit}
              disabled={loading || !content.trim()}
              className="px-4 py-2 bg-blue-400 hover:bg-blue-500 text-white font-bold rounded-full disabled:opacity-50"
            >
              {loading ? 'Posting...' : 'Tweet'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

export default CreatePost
