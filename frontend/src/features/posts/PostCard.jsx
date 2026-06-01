/*
** File: PostCard.jsx
** Description: Displays a single post card in the feed
** Responsibilities:
** - Render post content, author info and timestamp
** - Handle post editing and deletion
** - Show edit/delete buttons only to post author
** - Handle like/unlike toggle
*/

import { useState } from 'react'
import { deletePost, updatePost } from './postService.js'
import { Pencil, Trash2, MessageCircle, Heart } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import axiosInstance from '../../services/axiosInstance'

function PostCard({ post, onDelete, onUpdate, currentUserId }) {
  const navigate = useNavigate()
  const [isEditing, setIsEditing] = useState(false)
  const [editContent, setEditContent] = useState(post.content)
  const [loading, setLoading] = useState(false)
  const [liked, setLiked] = useState(post.liked || false)
  const [likesCount, setLikesCount] = useState(post.likes_count || 0)

  const handleDelete = async () => {
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

  const handleUpdate = async () => {
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

  const handleLike = async () => {
    try {
      const res = await axiosInstance.post(`/api/posts/${post.id}/like`)
      setLiked(res.data.liked)
      setLikesCount(res.data.likes_count)
    } catch (err) {
      console.error('Error liking post:', err)
    }
  }

  return (
    <div className="bg-white border border-gray-200 rounded-2xl p-4 mb-4 hover:shadow-sm transition">
      <div className="flex gap-3">
        <img
          src={post.author?.avatar || ''}
          alt={post.author?.username || ''}
          className="w-10 h-10 rounded-full bg-gray-200"
        />
        <div className="flex-1">
          <div className="flex items-center gap-2">
            <span
              onClick={() => navigate('/profile/' + post.author_id)}
              className="font-bold text-black cursor-pointer hover:underline"
            >
              {post.author?.name || 'Unknown'}
            </span>
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
            <>
              <p className="text-black mt-1">{post.content}</p>
              {post.media_url && (
                <div className="mt-3 overflow-hidden rounded-2xl">
                  <img
                    src={post.media_url}
                    alt="Post media"
                    className="w-full max-h-96 object-cover"
                  />
                </div>
              )}
            </>
          )}

          <div className="flex gap-6 mt-3 text-gray-500 text-sm">
            <span
              onClick={() => navigate(`/post/${post.id}`)}
              className="flex items-center gap-1 cursor-pointer hover:text-blue-400"
            >
              <MessageCircle size={16} />{post.comments_count}
            </span>
            <span
              onClick={post.author_id !== currentUserId ? handleLike : undefined}
              className={`flex items-center gap-1 ${post.author_id !== currentUserId
                  ? 'cursor-pointer hover:text-red-400'
                  : 'cursor-default'
                } ${liked ? 'text-red-500' : ''}`}
            >
              <Heart size={16} fill={liked ? 'currentColor' : 'none'} />
              {likesCount}
            </span>
            {post.author_id === currentUserId && (
              <div className="flex gap-3 ml-auto">
                <button onClick={() => setIsEditing(true)} className="hover:text-blue-400">
                  <Pencil size={16} />
                </button>
                <button onClick={handleDelete} disabled={loading} className="hover:text-red-400">
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