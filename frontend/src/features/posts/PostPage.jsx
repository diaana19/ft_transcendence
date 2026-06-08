import { useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import axiosInstance from "../../services/axiosInstance"
import RichText from "../../components/common/RichText"
import CommentForm from "./CommentForm"
import CommentItem from "./CommentItem"

export default function PostPage() {
  const { id } = useParams()

  const [post, setPost] = useState(null)
  const [comments, setComments] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchPost()
  }, [id])

  const fetchPost = async () => {
    try {
      const [postRes, commentsRes] = await Promise.all([
      axiosInstance.get(`/api/posts/${id}`),
      axiosInstance.get(`/api/posts/${id}/comments`)
    ])

      setPost(postRes.data)
      setComments(commentsRes.data.data || [])
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

	const handleCommentAdded = (newComment) => {
	setComments((prev) => [...prev, newComment])
	}


	const handleCommentDeleted = (commentId) => {
		setComments(prev => prev.filter(c => c.id !== commentId))
	}

	const handleCommentUpdated = (updatedComment) => {
		setComments(prev => prev.map(c => c.id === updatedComment.id ? updatedComment : c))
	}

  if (loading) {
    return (
      <div className="p-4 text-gray-500">
        Loading...
      </div>
    )
  }

  if (!post) {
    return (
      <div className="p-4 text-red-500">
        Post not found
      </div>
    )
  }

  return (
    <div className="w-full mx-auto min-h-screen bg-transparent" style={{ borderLeft: '0.5px solid #ede8fd', borderRight: '0.5px solid #ede8fd' }} >

      <div className="p-4" >
        <div className="rounded-2xl p-4" style={{ background: 'white', border: '0.5px solid #ede8fd' }}>

          <div className="w-11 h-11 rounded-full overflow-hidden bg-gray-300 flex items-center justify-center font-bold flex-shrink-0" style={{ background: '#ede8fd', color: '#534ab7' }}>
            {post.author?.avatar ? (
              <img
                src={post.author.avatar}
                alt={post.author.username}
                className="w-full h-full object-cover"
              />
            ) : (
              <span>
                {post.author?.username?.charAt(0).toUpperCase() || 'U'}
              </span>
            )}
          </div>
          <div className="flex-1">
            <div className="flex items-center gap-2">
              <h2 className="font-semibold text-sm" style={{ color: '#2c2c2a' }}>
                {post.author?.displayname}
              </h2>

              <span className="text-sm" style={{ color: '#afa9ec' }}>
                @{post.author?.username}
              </span>
            </div>

            <p className="mt-2 text-[15px] whitespace-pre-wrap" style={{ color: '#2c2c2a' }}>
              <RichText text={post.content} />
            </p>
            {post.media_url && (
            <div className="mt-3 overflow-hidden rounded-2xl " style={{ border: '0.5px solid #ede8fd' }}>
              {post.media_mime?.startsWith('video/') ? (
                <video
                  src={post.media_url}
                  controls
                  className="w-full max-h-[500px] object-contain"
                />
              ) : (
                <img
                  src={post.media_url}
                  alt="Post media"
                  className="w-full max-h-[500px] object-contain"
                />
              )}
            </div>
          )}
          </div>
        </div>
      </div>

	<div className="px-4 pb-2 pt-4" style={{ borderBottom: '0.5px solid #ede8fd' }} />

      <div className="px-4 pt-4">
        {comments.map((comment, index) => (
          <CommentItem
            key={comment.id}
            postId={id}
            comment={comment}
            isLast={index === comments.length - 1}
			onDelete={handleCommentDeleted}
  			onUpdate={handleCommentUpdated}
          />
        ))}
      </div>

      <CommentForm
        postId={id}
        onCommentAdded={handleCommentAdded}
      />

    </div>
  )
}
