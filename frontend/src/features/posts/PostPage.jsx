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
    <div className="max-w-2xl mx-auto border-x border-gray-200 min-h-screen bg-white">
      
      {/* POST */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex gap-3">
          
          {/* Avatar */}
          <div className="w-11 h-11 rounded-full overflow-hidden bg-gray-300 flex items-center justify-center text-white font-bold flex-shrink-0">
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
              <h2 className="font-bold text-sm">
                {post.author?.username}
              </h2>

              <span className="text-gray-500 text-sm">
                @{post.author?.username}
              </span>
            </div>

            <p className="mt-2 text-[15px] text-gray-900 whitespace-pre-wrap">
              {post.content}
            </p>
            {post.media_url && (
            <div className="mt-3 overflow-hidden rounded-2xl border border-gray-200">
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

      {/* COMMENT FORM */}
      <CommentForm
        postId={id}
        onCommentAdded={handleCommentAdded}
      />

      {/* COMMENTS */}
      <div>
        {comments.map((comment) => (
          <div
            key={comment._id}
            className="flex gap-3 p-4 border-b border-gray-200"
          >
            {/* Avatar */}
           <div className="w-9 h-9 rounded-full overflow-hidden bg-gray-300 flex items-center justify-center text-white font-bold flex-shrink-0">
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
            <div className="flex-1">
              <div className="flex items-center gap-2">
                <span className="font-bold text-sm">
                  {comment.author?.username}
                </span>

                <span className="text-gray-500 text-sm">
                  @{comment.author?.username}
                </span>
              </div>

              <p className="text-sm text-gray-900 mt-1 whitespace-pre-wrap">
                {comment.content}
              </p>
            </div>
          </div>
        ))}
      </div>

    </div>
  )
}