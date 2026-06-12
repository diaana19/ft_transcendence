import { useState } from 'react'
import { deletePost, updatePost, reactToPost } from './postService.js'
import { Pencil, Trash2, MessageCircle, ThumbsUp, ThumbsDown } from 'lucide-react'
import { useNavigate } from 'react-router-dom'
import axiosInstance from '../../services/axiosInstance'
import LoginModal from '../../components/common/LoginModal'
import RichText from '../../components/common/RichText'
import { useAuth } from '../../hooks/useAuth'

// PostCard shows one post with react, comment, edit and delete actions.
function PostCard({ post, onDelete, onUpdate, currentUserId }) {
    const navigate = useNavigate()
    const { user } = useAuth()
    const [isEditing, setIsEditing] = useState(false)
    const [editContent, setEditContent] = useState(post.content)
    const [editMediaUrl, setEditMediaUrl] = useState(post.media_url || null)
    const [editFile, setEditFile] = useState(null)
    const [editError, setEditError] = useState(null)
    const [loading, setLoading] = useState(false)
    const [userReaction, setUserReaction] = useState(post.liked ? 1 : post.disliked ? -1 : 0)
    const [likesCount, setLikesCount] = useState(post.likes_count || 0)
    const [dislikesCount, setDislikesCount] = useState(post.dislikes_count || 0)
    const [showLoginModal, setShowLoginModal] = useState(false)

    const handleDelete = async () => {
        setLoading(true)
        try {
            await deletePost(post.id)
            onDelete(post.id)
        } catch (err) {
            console.info('Error deleting post:', err)
        } finally {
            setLoading(false)
        }
    }

    // handleUpdate saves the edit, uploading a new file first if one was picked.
    const handleUpdate = async () => {
        setLoading(true)
        try {
            let mediaUrl = editMediaUrl
            let mediaMime = post.media_mime
            if (editFile) {
                const formData = new FormData()
                formData.append('file', editFile)
                const res = await axiosInstance.post('/api/upload', formData)
                mediaUrl = res.data.url
                mediaMime = res.data.mime_type
            }
            const updatedPost = await updatePost(post.id, editContent, mediaUrl, mediaMime)
            onUpdate({ ...updatedPost, media_url: mediaUrl, media_mime: mediaMime })
            setIsEditing(false)
            setEditFile(null)
        } catch (err) {
            console.info('Error updating post:', err)
        } finally {
            setLoading(false)
        }
    }

    // handleReact sends a like/dislike, or opens the login modal for guests.
    const handleReact = async (value) => {
        if (!user) {
            setShowLoginModal(true)
            return
        }
        try {
            const res = await reactToPost(post.id, value)
            setUserReaction(res.user_reaction)
            setLikesCount(res.likes_count)
            setDislikesCount(res.dislikes_count)
        } catch (err) {
            console.info('Error reacting to post:', err)
        }
    }

    const handleEditFileChange = (e) => {
        const selected = e.target.files[0]
        if (!selected) {
            return
        }
        setEditError(null)
        setEditFile(selected)
        setEditMediaUrl(URL.createObjectURL(selected))
    }

    // handleComment opens the post page, or the login modal for guests.
    const handleComment = () => {
        if (!user) {
            setShowLoginModal(true)
            return
        }
        navigate(`/post/${post.id}`)
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
                            className="font-bold cursor-pointer hover:underline"
                            style={{ color: '#2c2c2a' }}
                        >
                            {post.author?.displayname}
                        </span>
                        <span className="text-sm" style={{ color: '#afa9ec' }}>
                            @{post.author?.username}
                        </span>
                        <span className="text-gray-500">·</span>
                        <span className="text-gray-500 text-sm">
                            {new Date(post.created_at).toLocaleDateString('en-US', {
                                month: 'short',
                                day: 'numeric',
                            })}
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
                            {editMediaUrl && (
                                <div className="relative mt-2 mb-2">
                                    {(post.media_mime?.startsWith('video/') && !editFile) ||
                                    editFile?.type?.startsWith('video/') ? (
                                        <video
                                            src={editMediaUrl}
                                            controls
                                            className="rounded-xl max-h-48 w-full object-contain"
                                        />
                                    ) : (
                                        <img
                                            src={editMediaUrl}
                                            alt="preview"
                                            className="rounded-xl max-h-48 w-full object-contain"
                                        />
                                    )}
                                    <button
                                        onClick={() => {
                                            setEditMediaUrl(null)
                                            setEditFile(null)
                                        }}
                                        className="absolute top-2 right-2 w-6 h-6 rounded-full text-xs text-white flex items-center justify-center"
                                        style={{ background: 'rgba(83,74,183,0.7)' }}
                                    >
                                        ✕
                                    </button>
                                </div>
                            )}
                            <div className="flex items-center gap-2 mt-2">
                                <label className="cursor-pointer text-xs text-gray-400 hover:text-gray-600 flex items-center gap-1 flex-1">
                                    + Change media
                                    <input
                                        type="file"
                                        accept="image/*,video/*"
                                        onChange={handleEditFileChange}
                                        className="hidden"
                                    />
                                </label>
                                {editError && (
                                    <span className="text-xs" style={{ color: '#d4537e' }}>
                                        {editError}
                                    </span>
                                )}
                                <button
                                    onClick={() => {
                                        setIsEditing(false)
                                        setEditMediaUrl(post.media_url || null)
                                        setEditFile(null)
                                    }}
                                    className="px-4 py-1 rounded-full text-sm transition-colors"
                                    style={{
                                        border: '0.5px solid #ede8fd',
                                        color: '#534ab7',
                                        background: 'white',
                                    }}
                                >
                                    Cancel
                                </button>
                                <button
                                    onClick={handleUpdate}
                                    disabled={loading}
                                    className="px-4 py-1 rounded-full text-sm font-semibold transition-colors disabled:opacity-50"
                                    style={{
                                        background: '#534ab7',
                                        color: 'white',
                                        border: 'none',
                                    }}
                                >
                                    {loading ? 'Saving...' : 'Save'}
                                </button>{' '}
                            </div>
                        </div>
                    ) : (
                        <>
                            <p className="text-black mt-1 break-all">
                                <RichText text={post.content} />
                            </p>
                            {post.media_url && (
                                <div className="mt-3 overflow-hidden rounded-2xl">
                                    {post.media_mime?.startsWith('video/') ? (
                                        <video
                                            src={post.media_url}
                                            controls
                                            className="w-full max-h-96 object-contain"
                                        />
                                    ) : (
                                        <img
                                            src={post.media_url}
                                            alt="Post media"
                                            className="w-full max-h-96 object-contain"
                                        />
                                    )}
                                </div>
                            )}
                        </>
                    )}
                    {!isEditing && (
                        <div className="flex gap-6 mt-3 text-gray-500 text-sm">
                            <span
                                onClick={handleComment}
                                className="flex items-center gap-1 cursor-pointer hover:text-blue-400"
                            >
                                <MessageCircle size={16} />
                                {post.comments_count}
                            </span>
                            <span
                                onClick={() => handleReact(1)}
                                className={`flex items-center gap-1 cursor-pointer hover:text-green-500 ${userReaction === 1 ? 'text-green-600' : ''}`}
                            >
                                <ThumbsUp
                                    size={16}
                                    fill={userReaction === 1 ? 'currentColor' : 'none'}
                                />
                                {likesCount}
                            </span>
                            <span
                                onClick={() => handleReact(-1)}
                                className={`flex items-center gap-1 cursor-pointer hover:text-red-400 ${userReaction === -1 ? 'text-red-500' : ''}`}
                            >
                                <ThumbsDown
                                    size={16}
                                    fill={userReaction === -1 ? 'currentColor' : 'none'}
                                />
                                {dislikesCount}
                            </span>
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
                    )}
                </div>
            </div>

            <LoginModal isOpen={showLoginModal} onClose={() => setShowLoginModal(false)} />
        </div>
    )
}

export default PostCard
