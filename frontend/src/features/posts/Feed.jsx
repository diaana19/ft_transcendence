import { useEffect, useState } from 'react'
import { useAuth } from '../../hooks/useAuth'
import { getPosts, createPost } from './postService'
import PostCard from './PostCard'

function Feed() {
  const { user } = useAuth()
  const [posts, setPosts]       = useState([])
  const [content, setContent]   = useState('')
  const [loading, setLoading]   = useState(false)
  const [fetching, setFetching] = useState(true)

  useEffect(() => {
    fetchPosts()
  }, [])

  const fetchPosts = async () => {
    try {
      setFetching(true)
      const data = await getPosts()
      setPosts(data)
    } catch (err) {
      console.error('Error fetching posts:', err)
    } finally {
      setFetching(false)
    }
  }

  const handleCreate = async () => {
    if (!content.trim()) return
    setLoading(true)
    try {
      const newPost = await createPost(content)
      setPosts([newPost, ...posts]) // aparece arriba sin recargar
      setContent('')
    } catch (err) {
      console.error('Error creating post:', err)
    } finally {
      setLoading(false)
    }
  }

  // Callbacks que le pasa a PostCard
  const handleDelete = (postId) => {
    setPosts(posts.filter((p) => p.id !== postId))
  }

  const handleUpdate = (updatedPost) => {
    setPosts(posts.map((p) => (p.id === updatedPost.id ? updatedPost : p)))
  }

  return (
    <div className="max-w-xl mx-auto">

      {/* Composer */}
      <div className="border-b border-gray-200 p-4 flex gap-3">
        <div className="w-10 h-10 rounded-full bg-gray-300 flex-shrink-0" />
        <div className="flex-1">
          <textarea
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="What's happening?"
            maxLength={280}
            rows={3}
            className="w-full text-black placeholder-gray-500 text-lg focus:outline-none resize-none"
          />
          <div className="flex justify-between items-center mt-2">
            <span className="text-sm text-gray-400">{content.length}/280</span>
            <button
              onClick={handleCreate}
              disabled={loading || !content.trim()}
              className="bg-blue-400 hover:bg-blue-500 text-white font-bold px-5 py-2 rounded-full disabled:opacity-50 transition-colors"
            >
              {loading ? 'Posting...' : 'Post'}
            </button>
          </div>
        </div>
      </div>

      {/* Posts */}
      {fetching ? (
        <p className="text-center text-gray-400 py-8">Loading...</p>
      ) : posts.length === 0 ? (
        <p className="text-center text-gray-400 py-8">No posts yet</p>
      ) : (
        posts.map((post) => (
          <PostCard
            key={post.id}
            post={post}
            currentUserId={user?.userId}
            onDelete={handleDelete}
            onUpdate={handleUpdate}
          />
        ))
      )}

    </div>
  )
}

export default Feed