import { useEffect, useState } from 'react'
import { useAuth } from '../../hooks/useAuth'
import api from '../../services/axiosInstance'

const WEEKLY_TABS = [
    { key: 'total', label: 'Score' },
    { key: 'posts', label: 'Posts' },
    { key: 'likes', label: 'Likes' },
    { key: 'followers', label: 'Followers' },
]

// metricValue extracts the number shown for the active tab from an entry.
function metricValue(entry, tab) {
    return (tab === 'total' ? entry.stats?.total : entry.stats?.[tab]?.count) ?? 0
}

// metricLevel is the level for the active tab — the per-metric level for a
// specific tab (e.g. likes), or the overall level on the "Score" tab. Keeps the
// shown level consistent with the count: 0 likes reads as lv. 0, not the global
// level driven by posts/followers.
function metricLevel(entry, tab) {
    return (tab === 'total' ? entry.stats?.level : entry.stats?.[tab]?.level) ?? 0
}

// ScoreBar fills relative to the leader: the top entry is 100%, everyone else
// is their share of the leader's value for the active tab.
function ScoreBar({ pct }) {
    return (
        <div
            className="w-full rounded-full overflow-hidden"
            style={{ height: '4px', background: '#ede8fd' }}
        >
            <div
                className="rounded-full h-full transition-all"
                style={{ width: `${pct}%`, background: '#534ab7' }}
            />
        </div>
    )
}

export default function LeaderBoard() {
    const { user: authUser } = useAuth()
    const [leaderboard, setLeaderboard] = useState([])
    const [activeTab, setActiveTab] = useState('total')
    const [loadingBoard, setLoadingBoard] = useState(true)

    useEffect(() => {
        fetchLeaderboard()
    }, [])

    const fetchLeaderboard = async () => {
        try {
            const { data } = await api.get('/api/leaderboard')
            setLeaderboard(data)
        } catch (err) {
            console.error(err)
        } finally {
            setLoadingBoard(false)
        }
    }

    const sortedBoard = [...leaderboard]
        .sort((a, b) => metricValue(b, activeTab) - metricValue(a, activeTab))
        .slice(0, 10)
    // The leader (first, after the descending sort) defines 100%; every other
    // bar is its share of that value.
    const topValue = sortedBoard.length ? metricValue(sortedBoard[0], activeTab) : 0
    localStorage.getItem('token')
    return (
        <div>
            <div
                className="flex items-center gap-3 mb-3 pb-2"
                style={{ borderBottom: '1.5px solid #ede8fd' }}
            >
                <div
                    className="w-8 h-8 rounded-lg flex items-center justify-center"
                    style={{ background: '#ede8fd' }}
                >
                    <i
                        className="ti ti-trophy"
                        aria-hidden="true"
                        style={{ color: '#534ab7', fontSize: '15px' }}
                    />
                </div>
                <span className="text-sm font-semibold" style={{ color: '#534ab7' }}>
                    Weekly leaderboard
                </span>
            </div>
            <div className="flex gap-1 mb-3 p-1 rounded-xl" style={{ background: '#f5f0fe' }}>
                {WEEKLY_TABS.map((tab) => (
                    <button
                        key={tab.key}
                        onClick={() => setActiveTab(tab.key)}
                        className="flex-1 text-xs font-semibold py-1.5 rounded-lg transition-colors"
                        style={{
                            background: activeTab === tab.key ? 'white' : 'transparent',
                            color: activeTab === tab.key ? '#534ab7' : '#afa9ec',
                            border: 'none',
                            cursor: 'pointer',
                        }}
                    >
                        {tab.label}
                    </button>
                ))}
            </div>
            <div className="rounded-2xl overflow-hidden" style={{ border: '0.5px solid #ede8fd' }}>
                {loadingBoard ? (
                    <p className="text-center py-6 text-sm" style={{ color: '#b4b2a9' }}>
                        Loading...
                    </p>
                ) : sortedBoard.length === 0 ? (
                    <p className="text-center py-6 text-sm" style={{ color: '#b4b2a9' }}>
                        No data yet
                    </p>
                ) : (
                    sortedBoard.map((entry, i) => {
                        const isMe = entry.id === authUser?.userId
                        const value = metricValue(entry, activeTab)
                        const pct = topValue > 0 ? Math.min((value / topValue) * 100, 100) : 0
                        const medals = ['🥇', '🥈', '🥉']
                        return (
                            <div
                                key={entry.id}
                                className="flex items-center gap-3 px-4 py-3"
                                style={{
                                    borderBottom:
                                        i < sortedBoard.length - 1 ? '0.5px solid #f5f0fe' : 'none',
                                    background: isMe ? '#faf8ff' : 'white',
                                }}
                            >
                                <span
                                    className="text-sm w-6 text-center"
                                    style={{ color: '#b4b2a9' }}
                                >
                                    {i < 3 ? medals[i] : i + 1}
                                </span>
                                <div
                                    className="w-8 h-8 rounded-full flex items-center justify-center font-semibold flex-shrink-0 overflow-hidden"
                                    style={{
                                        background: '#ede8fd',
                                        color: '#534ab7',
                                        fontSize: '11px',
                                    }}
                                >
                                    {entry.avatar ? (
                                        <img
                                            src={entry.avatar}
                                            alt={entry.username}
                                            className="w-full h-full object-cover"
                                        />
                                    ) : (
                                        entry.username?.charAt(0).toUpperCase()
                                    )}
                                </div>
                                <div className="flex-1 min-w-0">
                                    <div className="flex items-center gap-2 mb-1">
                                        <span
                                            className="text-sm font-semibold truncate"
                                            style={{ color: '#2c2c2a' }}
                                        >
                                            {entry.username}
                                        </span>
                                        {isMe && (
                                            <span
                                                className="text-xs px-2 py-0.5 rounded-full flex-shrink-0"
                                                style={{ background: '#eeedfe', color: '#534ab7' }}
                                            >
                                                you
                                            </span>
                                        )}
                                    </div>
                                    <ScoreBar pct={pct} />
                                </div>
                                <div className="text-right flex-shrink-0">
                                    <div
                                        className="text-sm font-semibold"
                                        style={{ color: '#534ab7' }}
                                    >
                                        {value ?? 0}
                                    </div>
                                    <div className="text-xs" style={{ color: '#b4b2a9' }}>
                                        lv. {metricLevel(entry, activeTab)}
                                    </div>
                                </div>
                            </div>
                        )
                    })
                )}
            </div>
        </div>
    )
}
