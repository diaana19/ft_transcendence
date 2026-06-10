import { useEffect, useState, useCallback } from 'react'
import { useParams } from 'react-router-dom'
import { useAuth } from '../../hooks/useAuth'
import { getPostsByTagPage } from './postService'
import PostCard from './PostCard'

const PAGE_SIZE = 20

// TagFeed shows the posts of one hashtag with page-by-page pagination.
function TagFeed() {
    const { tag } = useParams()
    const { user } = useAuth()
    const [posts, setPosts] = useState([])
    const [fetching, setFetching] = useState(true)
    const [page, setPage] = useState(0)
    const [total, setTotal] = useState(0)

    const loadPage = useCallback(
        async (p) => {
            setFetching(true)
            try {
                const { posts: data, total: count } = await getPostsByTagPage(
                    tag,
                    PAGE_SIZE,
                    p * PAGE_SIZE
                )
                setPosts(data)
                setTotal(count)
                setPage(p)
            } catch (err) {
                console.info('Error fetching tagged posts:', err)
            } finally {
                setFetching(false)
            }
        },
        [tag]
    )

    useEffect(() => {
        loadPage(0)
    }, [loadPage])

    const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE))

    const goToPage = (p) => {
        if (p < 0 || p >= totalPages || p === page) return
        loadPage(p)
        window.scrollTo({ top: 0, behavior: 'smooth' })
    }

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
                        onDelete={() => loadPage(page)}
                        onUpdate={handleUpdate}
                    />
                ))
            )}

            {!fetching && total > 0 && (
                <div className="flex items-center justify-center gap-4 py-6">
                    <button
                        onClick={() => goToPage(page - 1)}
                        disabled={page === 0}
                        className="px-4 py-2 rounded-full text-sm font-bold disabled:opacity-40"
                        style={{ background: '#ede8fd', color: '#534ab7' }}
                    >
                        Previous
                    </button>
                    <span className="text-sm font-medium" style={{ color: '#534ab7' }}>
                        Page {page + 1} of {totalPages}
                    </span>
                    <button
                        onClick={() => goToPage(page + 1)}
                        disabled={page + 1 >= totalPages}
                        className="px-4 py-2 rounded-full text-sm font-bold disabled:opacity-40"
                        style={{ background: '#ede8fd', color: '#534ab7' }}
                    >
                        Next
                    </button>
                </div>
            )}
        </div>
    )
}

export default TagFeed
