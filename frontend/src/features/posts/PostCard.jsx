/*
** File: PostCard.jsx
** Description: Displays a single post card in the feed
** Responsibilities:
** - Render post content, author info and timestamp
** - Handle post editing and deletion
** - Show edit/delete buttons only to post author
*/

import { useState } from 'react'
import { deletePost, updatePost } from './postService.js'
import { Pencil, Trash2, MessageCircle, Heart } from 'lucide-react'

function PostCard({ post, onDelete, onUpdate, currentUserId }) {
    const [isEditing, setIsEditing] = useState(false)
    const [editContent, setEditContent] = useState(post.content)
    const [loading, setLoading] = useState(false)

const handleDelete = async() => {
    setLoading(true)
    try {
        await deletePost(post.id)
        onDelete(post.id)
    } catch (err) {
      console.error('Error deleting post:', err) 
    } finally {
        setLoading(false)
    }
}

const handleUpdate = async() => {
    setLoading(true)
    try {
        const updatedPost = await updatePost(post.id, editContent)
        onUpdate(updatedPost)
        setIsEditing(false)
    } catch (err) {
      console.error('Error updating post:', err) 
    } finally {
        setLoading(false)
    }
}

return (
    <div className="bg-white border border-gray-200 rounded-2xl p-4 mb-4 hover:shadow-sm transition">

      <div className="flex gap-3">
        <img
          src={post.author?.avatar || ''}
          alt={post.author?.username || ''}
          className="w-10 h-10 rounded-full"
        />
        <div className="flex-1">

          <div className="flex items-center gap-2">
            <span className="font-bold text-black">{post.author?.name || 'Unknown'}</span>
            <span className="text-gray-500">@{post.author?.username || 'unknown'}</span>
            <span className="text-gray-500">·</span>
            <span className="text-gray-500 text-sm">
              {new Date(post.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
            </span>
          </div>

          {isEditing ? (
            <div className="mt-2">
              <textarea
                value={editContent}
                onChange={(e) => setEditContent(e.target.value)}
                maxLength={280}
                className="w-full border border-gray-300 rounded p-2 text-black focus:outline-none focus:border-blue-400 resize-none"
                rows={3}
              />
              <div className="flex gap-2 mt-2">
                <button
                  onClick={() => setIsEditing(false)}
                  className="px-4 py-1 rounded-full border border-gray-300 text-gray-600 text-sm hover:bg-gray-100"
                >
                  Cancel
                </button>
                <button
                  onClick={handleUpdate}
                  disabled={loading}
                  className="px-4 py-1 rounded-full bg-blue-400 text-white text-sm font-bold hover:bg-blue-500 disabled:opacity-50"
                >
                  {loading ? 'Saving...' : 'Save'}
                </button>
              </div>
            </div>
          ) : (
            <p className="text-black mt-1">{post.content}</p>
          )}

          <div className="flex gap-6 mt-3 text-gray-500 text-sm">
           <span className="flex items-center gap-1"><MessageCircle size={16} />{post.comments_count}</span>
           <span className="flex items-center gap-1"><Heart size={16} />{post.likes_count}</span>
            {post.author_id === currentUserId && (
              <div className="flex gap-3 ml-auto">
                <button
                  onClick={() => setIsEditing(true)}
                  className="hover:text-blue-400"
                >
                 <Pencil size={16} />
                </button>
                <button
                  onClick={handleDelete}
                  disabled={loading}
                  className="hover:text-red-400"
                >
                  <Trash2 size={16} />
                </button>
              </div>
            )}
          </div>

        </div>
      </div>
    </div>
  )
}

export default PostCard