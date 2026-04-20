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

function CreatePost({ onPostCreated }) {
    const [file, setFile] = useState(null)
    const [content, setContent] = useState('')
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState(null)

const handleFileChange = (e) => {
  const selectedFile = e.target.files[0]
  if (selectedFile && selectedFile.size > 5 * 1024 * 1024) {
      setError('File too large. Maximum size is 5MB')
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
        <div className="w-10 h-10 rounded-full bg-gray-300 flex items-center justify-center text-white font-bold">
          T
        </div>
        <div className="flex-1">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            maxLength={280}
            placeholder="What's happening?"
            className="w-full text-black placeholder-gray-400 text-lg resize-none focus:outline-none pb-3"
            rows={2}
          />

          {file && (
            <div className="relative mb-2">
              <img
                src={URL.createObjectURL(file)}
                alt="preview"
                className="rounded-2xl max-h-64 w-full object-cover"
              />
              <button
                onClick={() => setFile(null)}
                className="absolute top-2 right-2 bg-black bg-opacity-50 text-white rounded-full w-6 h-6 flex items-center justify-center text-xs hover:bg-opacity-70"
              >
                ✕
              </button>
            </div>
          )}

          {error && (
            <p className="text-red-500 text-sm mb-2">{error}</p>
          )}

          <div className="flex items-center justify-between mt-2 border-t border-gray-100 pt-3">
            <div className="flex gap-3 text-blue-400">
              <label className="cursor-pointer hover:text-blue-500">
                <Image size={20} />
                <input
                  type="file"
                  accept="image/*"
                  onChange={handleFileChange}
                  className="hidden"
                />
              </label>
            </div>
            <div className="flex items-center gap-3">
              <span className={`text-sm ${content.length > 260 ? 'text-red-500' : 'text-gray-400'}`}>
                {content.length}/280
              </span>
              <button
                onClick={handleSubmit}
                disabled={loading || !content.trim()}
                className="px-4 py-2 bg-blue-400 hover:bg-blue-500 text-white font-bold rounded-full disabled:opacity-50 text-sm"
              >
                {loading ? 'Posting...' : 'Post'}
              </button>
            </div>
          </div>

        </div>
      </div>
    </div>
  )
}

export default CreatePost