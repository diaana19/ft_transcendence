import { useEffect, useState } from 'react'
import { useAuth } from '../../hooks/useAuth'
import { getPosts, createPost } from './postService'
import PostCard from './PostCard'
import CreatePost from './CreatePost'

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
      <CreatePost onPostCreated={(newPost) => setPosts([newPost, ...posts])} />

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