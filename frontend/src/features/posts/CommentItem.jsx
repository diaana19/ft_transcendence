/*
** File: CommentItem.jsx
** Description: Renders a single comment (reply) with like/dislike controls
** Responsibilities:
** - Display comment author, content and the connecting thread line
** - Handle like/dislike toggle for the comment
** - Show login modal if the user is not authenticated
*/

import { useState } from 'react'
import RichText from '../../components/common/RichText'
import LoginModal from '../../components/common/LoginModal'
import { useAuth } from '../../hooks/useAuth'
import { ThumbsUp, ThumbsDown, Pencil, Trash2 } from 'lucide-react'
import { reactToComment, updateComment, deleteComment } from './postService.js'

function CommentItem({ postId, comment, isLast, onDelete, onUpdate }) {
  const { user } = useAuth()
  const [userReaction, setUserReaction] = useState(
    comment.liked ? 1 : comment.disliked ? -1 : 0
  )
	const [likesCount, setLikesCount] = useState(comment.likes_count || 0)
	const [dislikesCount, setDislikesCount] = useState(comment.dislikes_count || 0)
	const [showLoginModal, setShowLoginModal] = useState(false)
	const [isEditing, setIsEditing] = useState(false)
	const [editContent, setEditContent] = useState(comment.content)
	const [editLoading, setEditLoading] = useState(false)
	const [editFileUrl, setEditFileUrl] = useState(null)
	const [editFile, setEditFile] = useState(null)

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

  const handleDelete = async () => {
  try {
    await deleteComment(postId, comment.id)
    onDelete(comment.id)
  } catch (err) {
    console.error(err)
  }
}

	const handleUpdate = async () => {
	if (!editContent.trim()) return
	setEditLoading(true)
	try {
		const updated = await updateComment(postId, comment.id, editContent)
		onUpdate(updated)
		setIsEditing(false)
	} catch (err) {
		console.error(err)
	} finally {
		setEditLoading(false)
	}
	}

	const handleEditFileChange = (e) => {
		const selected = e.target.files[0]
		if (!selected) return
		setEditFile(selected)
		setEditFileUrl(URL.createObjectURL(selected))
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

       {isEditing ? (
		<div className="mt-1">
			<textarea
				value={editContent}
				onChange={(e) => setEditContent(e.target.value)}
				maxLength={280}
				rows={2}
				className="w-full text-sm focus:outline-none resize-none bg-transparent"
				style={{ color: '#2c2c2a', border: '0.5px solid #ede8fd', borderRadius: '8px', padding: '6px' }}
			/>
			{(editFileUrl || comment.file_url) && !editFile?.cleared && (
				<div className="relative mt-2 mb-1">
				<img
					src={editFileUrl || comment.file_url}
					alt="preview"
					className="rounded-xl max-h-48 w-full object-cover"
				/>
				<button
					onClick={() => { setEditFileUrl('removed'); setEditFile({ cleared: true }) }}
					className="absolute top-2 right-2 w-5 h-5 rounded-full text-xs text-white flex items-center justify-center"
					style={{ background: 'rgba(83,74,183,0.7)' }}
				>✕</button>
				</div>
			)}
			<div className="flex items-center gap-2 mt-1">
				<label className="cursor-pointer text-xs flex items-center gap-1 flex-1" style={{ color: '#afa9ec' }}>
				+ Change image
				<input type="file" accept="image/*" onChange={handleEditFileChange} className="hidden" />
				</label>
				<button
				onClick={() => { setIsEditing(false); setEditContent(comment.content); setEditFileUrl(null); setEditFile(null) }}
				className="text-xs px-3 py-1 rounded-full"
				style={{ border: '0.5px solid #ede8fd', color: '#534ab7' }}
				>Cancel</button>
				<button
				onClick={handleUpdate}
				disabled={editLoading}
				className="text-xs px-3 py-1 rounded-full disabled:opacity-50"
				style={{ background: '#534ab7', color: 'white' }}
				>{editLoading ? 'Saving...' : 'Save'}</button>
			</div>
		</div>
		) : (
		<p className="text-sm mt-1 whitespace-pre-wrap" style={{ color: '#5f5e5a' }}>
			<RichText text={comment.content} />
		</p>
		)}


		{!isEditing && comment.file_url && (
		<div className="mt-2 overflow-hidden rounded-xl" style={{ border: '0.5px solid #ede8fd' }}>
			{comment.file_mime?.startsWith('video/') ? (
			<video src={comment.file_url} controls className="w-full max-h-72 object-cover" />
			) : (
			<img src={comment.file_url} alt="reply attachment" className="w-full max-h-72 object-cover" />
			)}
		</div>
		)}
		{!isEditing && (
		<div className="flex gap-5 mt-2 text-xs items-center" style={{ color: '#b4b2a9' }}>
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
		{user?.userId === comment.author_id && (
			<span className="flex gap-2 ml-auto">
			<Pencil size={13} className="cursor-pointer hover:text-blue-400" onClick={() => setIsEditing(true)} />
			<Trash2 size={13} className="cursor-pointer hover:text-red-400" onClick={handleDelete} />
			</span>
		)}
	  </div> )} </div>
      <LoginModal
        isOpen={showLoginModal}
        onClose={() => setShowLoginModal(false)}
      />
    </div>
  )
}
export default CommentItem
