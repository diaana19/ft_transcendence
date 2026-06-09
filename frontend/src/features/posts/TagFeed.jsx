import { useEffect, useRef, useState, useCallback } from 'react'
import { useParams } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'
import { getPostsByTag } from './postService'
import PostCard from './PostCard'

const PAGE_SIZE = 20

// TagFeed shows the posts of one hashtag with infinite scroll.
function TagFeed() {
    const { tag } = useParams()
    const { user } = useAuth()
    const [posts, setPosts] = useState([])
    const [fetching, setFetching] = useState(true)
    const [loadingMore, setLoadingMore] = useState(false)
    const [hasMore, setHasMore] = useState(true)

    const offsetRef = useRef(0)
    const hasMoreRef = useRef(true)
    const loadingRef = useRef(false)
    const sentinelRef = useRef(null)

    // loadMore fetches the next page of tagged posts, skipping duplicates.
    const loadMore = useCallback(async () => {
        if (loadingRef.current || !hasMoreRef.current) return
        loadingRef.current = true
        setLoadingMore(true)
        try {
            const data = await getPostsByTag(tag, PAGE_SIZE, offsetRef.current)
            offsetRef.current += data.length
            // A full page means there may be more posts to load.
            hasMoreRef.current = data.length === PAGE_SIZE
            setHasMore(hasMoreRef.current)
            setPosts((prev) => {
                const seen = new Set(prev.map((p) => p.id))
                return [...prev, ...data.filter((p) => !seen.has(p.id))]
            })
        } catch (err) {
            console.info('Error fetching tagged posts:', err)
        } finally {
            loadingRef.current = false
            setLoadingMore(false)
            setFetching(false)
        }
    }, [tag])

    // Reset everything and reload when the tag in the URL changes.
    useEffect(() => {
        offsetRef.current = 0
        hasMoreRef.current = true
        loadingRef.current = false
        setPosts([])
        setHasMore(true)
        setFetching(true)
        loadMore()
    }, [tag, loadMore])

    // Load more posts when the sentinel div gets near the viewport.
    useEffect(() => {
        const el = sentinelRef.current
        if (!el || !hasMore) return
        const obs = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting) loadMore()
            },
            { rootMargin: '300px' }
        )
        obs.observe(el)
        return () => obs.disconnect()
    }, [loadMore, hasMore, fetching])

    const reload = useCallback(async () => {
        if (loadingRef.current) return
        loadingRef.current = true
        try {
            const data = await getPostsByTag(tag, PAGE_SIZE, 0)
            offsetRef.current = data.length
            hasMoreRef.current = data.length === PAGE_SIZE
            setHasMore(hasMoreRef.current)
            setPosts(data)
        } catch (err) {
            console.info('Error fetching tagged posts:', err)
        } finally {
            loadingRef.current = false
        }
    }, [tag])

    const handleUpdate = (updatedPost) => {
        setPosts((prev) => prev.map((p) => (p.id === updatedPost.id ? updatedPost : p)))
    }

    return (
        <div className="w-full mx-auto space-y-4 mt-6">
            <h2 className="text-xl font-bold" style={{ color: '#534ab7' }}>
                #{tag}
            </h2>

            {fetching ? (
                <p className="text-center text-gray-400 py-8">Loading...</p>
            ) : posts.length === 0 ? (
                <p className="text-center text-gray-400 py-8">No posts with this tag yet</p>
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

            {!fetching && hasMore && <div ref={sentinelRef} className="h-1" />}
            {loadingMore && <p className="text-center text-gray-400 py-4">Loading more...</p>}
            {!fetching && !hasMore && posts.length > 0 && (
                <p className="text-center text-gray-300 py-6 text-sm">You're all caught up</p>
            )}
        </div>
    )
}

export default TagFeed
