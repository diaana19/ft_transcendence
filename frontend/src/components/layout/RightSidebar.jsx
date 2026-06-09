import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import axiosInstance from '../../services/axiosInstance'

// Trends reflect hashtags used in posts over the last week, counted live by the
// backend. Refetched on mount and every 60s so the list stays current without
// a manual refresh.
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
                console.error('Error loading trends:', err)
            } finally {
                if (active) setLoading(false)
            }
        }
        load()
        const id = setInterval(load, 60000)
        return () => {
            active = false
            clearInterval(id)
        }
    }, [])

    if (loading) {
        return <p className="text-xs text-gray-400">Loading...</p>
    }
    if (trends.length === 0) {
        return <p className="text-xs text-gray-400">Nothing here yet...</p>
    }

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

export default function RightSidebar() {
    return (
        <aside
            className="w-80 h-screen fixed right-0 top-0
                      backdrop-blur-xl
                      border-r border-transparent
                      px-4 py-6 hidden lg:block"
        >
            <div className="bg-white/60 border border-gray-200 rounded-xl p-4 mb-4">
                <h3 className="text-sm font-semibold text-gray-700 mb-2">Trends</h3>
                <Trends />
            </div>

            <div className="bg-white/60 border border-gray-200 rounded-xl p-4">
                <h3 className="text-sm font-semibold text-gray-700 mb-2">Suggestions</h3>
                <p className="text-xs text-gray-400">
                    Future feature: users / posts recommendations
                </p>
            </div>
        </aside>
    )
}
