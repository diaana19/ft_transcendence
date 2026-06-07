/*
** File: CommentItem.jsx
** Description: Renders a single comment (reply) with like/dislike controls
** Responsibilities:
** - Display comment author, content and the connecting thread line
** - Handle like/dislike toggle for the comment
** - Show login modal if the user is not authenticated
*/

import { useState } from 'react'
import { ThumbsUp, ThumbsDown } from 'lucide-react'
import { reactToComment } from './postService.js'
import RichText from '../../components/common/RichText'
import LoginModal from '../../components/common/LoginModal'
import { useAuth } from '../../hooks/useAuth'

function CommentItem({ postId, comment, isLast }) {
  const { user } = useAuth()
  const [userReaction, setUserReaction] = useState(
    comment.liked ? 1 : comment.disliked ? -1 : 0
  )
  const [likesCount, setLikesCount] = useState(comment.likes_count || 0)
  const [dislikesCount, setDislikesCount] = useState(comment.dislikes_count || 0)
  const [showLoginModal, setShowLoginModal] = useState(false)

  const handleReact = async (value) => {
    if (!user) {
      setShowLoginModal(true)
      return
    }
    try {
      const res = await reactToComment(postId, comment.id, value)
      setUserReaction(res.user_reaction)
      setLikesCount(res.likes_count)
      setDislikesCount(res.dislikes_count)
    } catch (err) {
      console.error('Error reacting to comment:', err)
    }
  }

  return (
    <div className="flex gap-3 mb-3">
      {/* Avatar */}
      <div className="flex flex-col items-center">
        <div className="w-8 h-8 rounded-full overflow-hidden flex items-center justify-center font-bold flex-shrink-0">
          {comment.author?.avatar ? (
            <img
              src={comment.author.avatar}
              alt={comment.author.username}
              className="w-full h-full object-cover"
            />
          ) : (
            <span>
              {comment.author?.username?.charAt(0).toUpperCase() || 'U'}
            </span>
          )}
        </div>
        {!isLast && (
          <div className="w-px flex-1 mt-1" style={{ background: '#ede8fd', minHeight: '16px' }} />
        )}
      </div>
      <div className="flex-1 rounded-xl px-3 py-2 mb-1" style={{ background: 'white', border: '0.5px solid #f0ebfe' }}>
        <div className="flex items-center gap-2">
          <span className="font-semibold text-xs" style={{ color: '#2c2c2a' }}>
            {comment.author?.displayname}
          </span>

          <span className="text-xs" style={{ color: '#afa9ec' }}>
            @{comment.author?.username}
          </span>
        </div>

        <p className="text-sm mt-1 whitespace-pre-wrap" style={{ color: '#5f5e5a' }}>
          <RichText text={comment.content} />
        </p>

        {comment.file_url && (
          <div className="mt-2 overflow-hidden rounded-xl" style={{ border: '0.5px solid #ede8fd' }}>
            {comment.file_mime?.startsWith('video/') ? (
              <video src={comment.file_url} controls className="w-full max-h-72 object-cover" />
            ) : (
              <img src={comment.file_url} alt="reply attachment" className="w-full max-h-72 object-cover" />
            )}
          </div>
        )}

        <div className="flex gap-5 mt-2 text-xs" style={{ color: '#b4b2a9' }}>
          <span
            onClick={() => handleReact(1)}
            className={`flex items-center gap-1 cursor-pointer hover:text-green-500 ${userReaction === 1 ? 'text-green-600' : ''}`}
          >
            <ThumbsUp size={14} fill={userReaction === 1 ? 'currentColor' : 'none'} />
            {likesCount}
          </span>
          <span
            onClick={() => handleReact(-1)}
            className={`flex items-center gap-1 cursor-pointer hover:text-red-400 ${userReaction === -1 ? 'text-red-500' : ''}`}
          >
            <ThumbsDown size={14} fill={userReaction === -1 ? 'currentColor' : 'none'} />
            {dislikesCount}
          </span>
        </div>
      </div>

      <LoginModal
        isOpen={showLoginModal}
        onClose={() => setShowLoginModal(false)}
      />
    </div>
  )
}

export default CommentItem
