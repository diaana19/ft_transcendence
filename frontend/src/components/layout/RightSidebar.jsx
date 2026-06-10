import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import axiosInstance from '../../services/axiosInstance'
import { useAuth } from '../../hooks/useAuth'

// Trends shows the top hashtags and refreshes them every minute.
function Trends() {
    const [trends, setTrends] = useState([])
    const [loading, setLoading] = useState(true)
    useEffect(() => {
        let active = true
        const load = async () => {
            try {
                const { data } = await axiosInstance.get('/api/trends?limit=8')
                if (active) setTrends(data.data || [])
            } catch (err) {
                console.info('Error loading trends:', err)
            } finally {
                if (active) setLoading(false)
            }
        }
        load()
        // Reload trends every 60 seconds.
        const id = setInterval(load, 60000)
        return () => {
            active = false
            clearInterval(id)
        }
    }, [])
    if (loading) return <p className="text-xs text-gray-400">Loading...</p>
    if (trends.length === 0) return <p className="text-xs text-gray-400">Nothing here yet...</p>
    return (
        <ul className="space-y-2">
            {trends.map((t) => (
                <li key={t.tag} className="flex items-center justify-between">
                    <Link
                        to={`/tag/${t.tag.replace(/^#/, '')}`}
                        className="text-sm font-semibold truncate hover:underline"
                        style={{ color: '#534ab7' }}
                    >
                        {t.tag}
                    </Link>
                    <span className="text-xs text-gray-400 flex-shrink-0 ml-2">
                        {t.count} {t.count === 1 ? 'post' : 'posts'}
                    </span>
                </li>
            ))}
        </ul>
    )
}

// Suggestions shows up to 3 users that the current user does not follow yet.
function Suggestions() {
    const { user } = useAuth()
    const navigate = useNavigate()
    const [suggestions, setSuggestions] = useState([])
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        if (!user?.userId) return
        const load = async () => {
            try {
                const [usersRes, followingRes] = await Promise.all([
                    axiosInstance.get('/api/users'),
                    axiosInstance.get(`/api/users/${user.userId}/following`),
                ])
                const following = (followingRes.data?.data || followingRes.data || []).map(
                    (f) => f.id
                )
                // Keep only users that are not me and that I do not follow, max 3.
                const filtered = (usersRes.data || [])
                    .filter((u) => u.id !== user.userId && !following.includes(u.id))
                    .slice(0, 3)
                setSuggestions(filtered)
            } catch (err) {
                console.info(err)
            } finally {
                setLoading(false)
            }
        }
        load()
    }, [user?.userId])

    if (!user?.userId) return <p className="text-xs text-gray-400">Log in to see suggestions</p>
    if (loading) return <p className="text-xs text-gray-400">Loading...</p>
    if (suggestions.length === 0) return <p className="text-xs text-gray-400">No suggestions yet</p>

    return (
        <ul className="space-y-3">
            {suggestions.map((u) => (
                <li
                    key={u.id}
                    className="flex items-center gap-2 cursor-pointer"
                    onClick={() => navigate(`/profile/${u.id}`)}
                >
                    <div
                        className="w-8 h-8 rounded-full overflow-hidden flex-shrink-0 flex items-center justify-center font-bold text-xs"
                        style={{ background: '#ede8fd', color: '#534ab7' }}
                    >
                        {u.avatar ? (
                            <img
                                src={u.avatar}
                                alt={u.username}
                                className="w-full h-full object-cover"
                            />
                        ) : (
                            u.username?.[0]?.toUpperCase()
                        )}
                    </div>
                    <div className="flex-1 min-w-0">
                        <p className="text-xs font-bold truncate" style={{ color: '#2c2c2a' }}>
                            {u.displayname || u.username}
                        </p>
                        <p className="text-xs truncate" style={{ color: '#afa9ec' }}>
                            @{u.username}
                        </p>
                    </div>
                </li>
            ))}
        </ul>
    )
}

// RightSidebar shows the Trends and Who to follow boxes.
export default function RightSidebar() {
    return (
        <aside className="w-80 h-screen fixed right-0 top-0 backdrop-blur-xl border-r border-transparent px-4 py-6 hidden lg:block">
            <div className="bg-white/60 border border-gray-200 rounded-xl p-4 mb-4">
                <h3 className="text-sm font-semibold text-gray-700 mb-2">Trends</h3>
                <Trends />
                <Link
                    to="/tags"
                    className="block text-xs font-semibold mt-3 hover:underline"
                    style={{ color: '#534ab7' }}
                >
                    Search tags →
                </Link>
            </div>
            <div className="bg-white/60 border border-gray-200 rounded-xl p-4">
                <h3 className="text-sm font-semibold text-gray-700 mb-2">Who to follow</h3>
                <Suggestions />
            </div>
        </aside>
    )
}
