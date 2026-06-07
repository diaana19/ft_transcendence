import { useEffect, useRef, useState, useCallback } from 'react'
import { useAuth } from '../../hooks/useAuth'
import { getPosts } from './postService'
import PostCard from './PostCard'
import CreatePost from './CreatePost'

const PAGE_SIZE = 20

function Feed() {
  const { user } = useAuth()
  const [posts, setPosts] = useState([])
  const [fetching, setFetching] = useState(true)   // initial load
  const [loadingMore, setLoadingMore] = useState(false)
  const [hasMore, setHasMore] = useState(true)

  // Refs so the IntersectionObserver callback always reads current values
  // without re-subscribing, and concurrent loads can't overlap.
  const offsetRef = useRef(0)
  const hasMoreRef = useRef(true)
  const loadingRef = useRef(false)
  const sentinelRef = useRef(null)

  // loadMore appends the next page, skipping any ids already shown (the feed can
  // shift between requests as posts are created/deleted, so offset paging can
  // overlap — dedup keeps it clean).
  const loadMore = useCallback(async () => {
    if (loadingRef.current || !hasMoreRef.current) return
    loadingRef.current = true
    setLoadingMore(true)
    try {
      const data = await getPosts(PAGE_SIZE, offsetRef.current)
      offsetRef.current += data.length
      hasMoreRef.current = data.length === PAGE_SIZE
      setHasMore(hasMoreRef.current)
      setPosts((prev) => {
        const seen = new Set(prev.map((p) => p.id))
        return [...prev, ...data.filter((p) => !seen.has(p.id))]
      })
    } catch (err) {
      console.error('Error fetching posts:', err)
    } finally {
      loadingRef.current = false
      setLoadingMore(false)
      setFetching(false)
    }
  }, [])

  // First page on mount.
  useEffect(() => { loadMore() }, [loadMore])

  // Auto-load when the bottom sentinel scrolls into view. Re-subscribes when
  // hasMore flips so the observer detaches once the feed is exhausted.
  useEffect(() => {
    const el = sentinelRef.current
    if (!el || !hasMore) return
    const obs = new IntersectionObserver(
      (entries) => { if (entries[0].isIntersecting) loadMore() },
      { rootMargin: '300px' },
    )
    obs.observe(el)
    return () => obs.disconnect()
  }, [loadMore, hasMore, fetching])

  // reload refetches the first page and resets paging. Cheaper than reconciling
  // offsets by hand: a create/delete shifts the whole window, so we just start
  // over from the top rather than track the drift.
  const reload = useCallback(async () => {
    if (loadingRef.current) return
    loadingRef.current = true
    try {
      const data = await getPosts(PAGE_SIZE, 0)
      offsetRef.current = data.length
      hasMoreRef.current = data.length === PAGE_SIZE
      setHasMore(hasMoreRef.current)
      setPosts(data)
    } catch (err) {
      console.error('Error fetching posts:', err)
    } finally {
      loadingRef.current = false
    }
  }, [])

  const handleCreated = async () => {
    await reload()
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  // Update in place — no count/order change, so no refetch needed.
  const handleUpdate = (updatedPost) => {
    setPosts((prev) => prev.map((p) => (p.id === updatedPost.id ? updatedPost : p)))
  }

  return (
    <div className="w-full mx-auto space-y-4 mt-6 ">

      {/* Composer */}
      {user?.userId && (
        <CreatePost onPostCreated={handleCreated} user={user} />
      )}

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
            onDelete={reload}
            onUpdate={handleUpdate}
          />
        ))
      )}

      {/* Infinite-scroll trigger + status */}
      {!fetching && hasMore && <div ref={sentinelRef} className="h-1" />}
      {loadingMore && <p className="text-center text-gray-400 py-4">Loading more...</p>}
      {!fetching && !hasMore && posts.length > 0 && (
        <p className="text-center text-gray-300 py-6 text-sm">You're all caught up</p>
      )}
    </div>
  )
}

export default Feed
