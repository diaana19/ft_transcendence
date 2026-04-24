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
      console.log('DATA:', data)
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
      setPosts([newPost, ...posts])
      setContent('')
    } catch (err) {
      console.error('Error creating post:', err)
    } finally {
      setLoading(false)
    }
  }


  const handleDelete = (postId) => {
    setPosts(posts.filter((p) => p.id !== postId))
  }

  const handleUpdate = (updatedPost) => {
    setPosts(posts.map((p) => (p.id === updatedPost.id ? updatedPost : p)))
  }

  return (
    <div className="max-w-2xl mx-auto space-y-4">

      {/* Composer */}
        <CreatePost onPostCreated={() => fetchPosts()} />
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
