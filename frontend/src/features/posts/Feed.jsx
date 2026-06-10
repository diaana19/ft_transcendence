import { useEffect, useState, useCallback } from 'react'
import { useAuth } from '../../hooks/useAuth'
import { getPostsPage } from './postService'
import PostCard from './PostCard'
import CreatePost from './CreatePost'
import api from '../../services/axiosInstance'

const PAGE_SIZE = 20

// BADGES_CHECK is the list of badges with the rule to know if one is earned.
const BADGES_CHECK = [
    { key: 'welcome', name: 'Welcome', check: (s) => true },
    { key: 'first_bond', name: 'First bond', check: (s) => s.followers >= 1 },
    { key: 'first_words', name: 'First words', check: (s) => s.posts >= 1 },
    { key: 'first_spark', name: 'First spark', check: (s) => s.likes >= 1 },
    { key: 'network_builder', name: 'Network builder', check: (s) => s.followers >= 10 },
    { key: 'content_creator', name: 'Content creator', check: (s) => s.posts >= 25 },
    { key: 'rising_star', name: 'Rising star', check: (s) => s.likes >= 20 },
    { key: 'social_butterfly', name: 'Social butterfly', check: (s) => s.followers >= 50 },
]

// Feed is the home timeline with page-by-page pagination and a badge toast.
function Feed() {
    const { user } = useAuth()
    const [posts, setPosts] = useState([])
    const [fetching, setFetching] = useState(true)
    const [page, setPage] = useState(0)
    const [total, setTotal] = useState(0)
    const [newBadge, setNewBadge] = useState(null)

    // Poll the gamification stats and show a toast when a new badge is unlocked.
    useEffect(() => {
        if (!user?.userId) return
        const checkBadges = async () => {
            try {
                const { data } = await api.get(`/api/users/${user.userId}/gamification`)
                const userStats = {
                    posts: data?.posts?.count ?? 0,
                    likes: data?.likes?.count ?? 0,
                    followers: data?.followers?.count ?? 0,
                }
                const earned = BADGES_CHECK.filter((b) => b.check(userStats))
                const prev = JSON.parse(localStorage.getItem('earned_badges') || 'null')
                // First run only saves the current badges, no toast.
                if (prev === null) {
                    localStorage.setItem('earned_badges', JSON.stringify(earned.map((b) => b.key)))
                    return
                }
                const newOnes = earned.filter((b) => !prev.includes(b.key))
                if (newOnes.length > 0) {
                    setNewBadge(newOnes[0])
                    setTimeout(() => setNewBadge(null), 4000)
                    localStorage.setItem('earned_badges', JSON.stringify(earned.map((b) => b.key)))
                }
            } catch (err) {
                console.info(err)
            }
        }
        checkBadges()
        const interval = setInterval(checkBadges, 3000)
        return () => clearInterval(interval)
    }, [user?.userId])

    const loadPage = useCallback(async (p) => {
        setFetching(true)
        try {
            const { posts: data, total: count } = await getPostsPage(PAGE_SIZE, p * PAGE_SIZE)
            setPosts(data)
            setTotal(count)
            setPage(p)
        } catch (err) {
            console.info('Error fetching posts:', err)
        } finally {
            setFetching(false)
        }
    }, [])

    useEffect(() => {
        loadPage(0)
    }, [loadPage])

    const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE))

    const goToPage = (p) => {
        if (p < 0 || p >= totalPages || p === page) return
        loadPage(p)
        window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    const handleCreated = async () => {
        await loadPage(0)
        window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    const handleUpdate = (updatedPost) => {
        setPosts((prev) => prev.map((p) => (p.id === updatedPost.id ? updatedPost : p)))
    }

    return (
        <div className="w-full mx-auto space-y-4 mt-6">
            {user?.userId && <CreatePost onPostCreated={handleCreated} user={user} />}

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

            {newBadge && (
                <div
                    className="fixed bottom-20 left-1/2 z-50 px-4 py-3 rounded-2xl shadow-lg flex items-center gap-3"
                    style={{
                        background: 'white',
                        border: '1.5px solid #ede8fd',
                        transform: 'translateX(-50%)',
                    }}
                >
                    <span className="text-xl">🏆</span>
                    <div>
                        <p className="text-xs font-bold" style={{ color: '#534ab7' }}>
                            New badge unlocked!
                        </p>
                        <p className="text-xs" style={{ color: '#b4b2a9' }}>
                            {newBadge.name}
                        </p>
                    </div>
                </div>
            )}
        </div>
    )
}

export default Feed
