import { useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import axiosInstance from "../../services/axiosInstance"
import CommentForm from "./CommentForm"

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
      const res = await axiosInstance.get(`/api/posts/${id}`)

      setPost(res.data)
      setComments(res.data.comments || [])
    } catch (err) {
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleCommentAdded = (newComment) => {
    setComments((prev) => [newComment, ...prev])
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

      {/* POST */}
      <div className="p-4" >
        <div className="rounded-2xl p-4" style={{ background: 'white', border: '0.5px solid #ede8fd' }}>

          {/* Avatar */}
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
                {post.author?.username}
              </h2>

              <span className="text-sm" style={{ color: '#afa9ec' }}>
                @{post.author?.username}
              </span>
            </div>

            <p className="mt-2 text-[15px] whitespace-pre-wrap" style={{ color: '#2c2c2a' }}>
              {post.content}
            </p>
            {post.media_url && (
            <div className="mt-3 overflow-hidden rounded-2xl " style={{ border: '0.5px solid #ede8fd' }}>
              {post.media_url.match(/\.(mp4|webm|ogg)$/i) ? (
                <video
                  src={post.media_url}
                  controls
                  className="w-full max-h-[500px] object-cover"
                />
              ) : (
                <img
                  src={post.media_url}
                  alt="Post media"
                  className="w-full max-h-[500px] object-cover"
                />
              )}
            </div>
          )}
          </div>
        </div>
      </div>

	<div className="px-4 pb-2 pt-4" style={{ borderBottom: '0.5px solid #ede8fd' }} />

      {/* COMMENT FORM */}
      <CommentForm
        postId={id}
        onCommentAdded={handleCommentAdded}
      />

      {/* COMMENTS */}
      <div className="px-4 pt-4">
        {comments.map((comment,index) => (
          <div
            key={comment.id}
            className="flex gap-3 mb-3"
          >
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
				{index < comments.length - 1 && (
				<div className="w-px flex-1 mt-1" style={{ background: '#ede8fd', minHeight: '16px' }} />
				)}
			</div>
            <div className="flex-1 rounded-xl px-3 py-2 mb-1" style={{ background: 'white', border: '0.5px solid #f0ebfe' }} >
              <div className="flex items-center gap-2">
                <span className="font-semibold text-xs" style={{ color: '#2c2c2a' }}>
                  {comment.author?.username}
                </span>

                <span className="text-xs" style={{ color: '#afa9ec' }}>
                  @{comment.author?.username}
                </span>
              </div>

              <p className="text-sm mt-1 whitespace-pre-wrap" style={{ color: '#5f5e5a' }}>
                {comment.content}
              </p>
            </div>
          </div>
        ))}
      </div>

    </div>
  )
}
