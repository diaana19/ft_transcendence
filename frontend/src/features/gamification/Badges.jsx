import { useState, useEffect } from 'react'
import { useAuth } from '../../hooks/useAuth'
import api from '../../services/axiosInstance'

const BADGES = [
    {
        category: 'social',
        items: [
            {
                key: 'welcome',
                icon: 'ti-user-plus',
                name: 'Welcome',
                desc: 'Created your account on Synk',
                condition: 'Day one',
                check: (s) => true,
            },
            {
                key: 'first_bond',
                icon: 'ti-heart-handshake',
                name: 'First bond',
                desc: 'Added your very first friend',
                condition: '1 friend',
                check: (s) => s.followers >= 1,
            },
            {
                key: 'network_builder',
                icon: 'ti-topology-star',
                name: 'Network builder',
                desc: 'Reached 10 friends on the platform',
                condition: '10 friends',
                check: (s) => s.followers >= 10,
            },
            {
                key: 'social_butterfly',
                icon: 'ti-butterfly',
                name: 'Social butterfly',
                desc: 'Connected with more than 50 friends',
                condition: '50 friends',
                check: (s) => s.followers >= 50,
            },
        ],
    },
    {
        category: 'posts',
        items: [
            {
                key: 'first_words',
                icon: 'ti-feather',
                name: 'First words',
                desc: 'Published your very first post',
                condition: '1 post',
                check: (s) => s.posts >= 1,
            },
            {
                key: 'first_spark',
                icon: 'ti-heart',
                name: 'First spark',
                desc: 'Received your first like on a post',
                condition: '1 like',
                check: (s) => s.likes >= 1,
            },
            {
                key: 'content_creator',
                icon: 'ti-writing',
                name: 'Content creator',
                desc: 'Published 25 posts in total',
                condition: '25 posts',
                check: (s) => s.posts >= 25,
            },
            {
                key: 'rising_star',
                icon: 'ti-star',
                name: 'Rising star',
                desc: 'Collected 20 likes across your posts',
                condition: '20 likes',
                check: (s) => s.likes >= 20,
            },
        ],
    },
]

const CATEGORY_STYLE = {
    social: {
        border: '#ede8fd',
        icon: '#ede8fd',
        iconText: '#534ab7',
        title: '#534ab7',
        cardBorder: '#ede8fd',
        pillBg: '#eeedfe',
        pillText: '#534ab7',
        emojiBg: '#ede8fd',
        progressBg: '#534ab7',
    },
    posts: {
        border: '#c5dbc4',
        icon: '#c5dbc4',
        iconText: '#3b6d11',
        title: '#3b6d11',
        cardBorder: '#dce6cf',
        pillBg: '#eaf3de',
        pillText: '#3b6d11',
        emojiBg: '#c5dbc4',
        progressBg: '#3b6d11',
    },
}
const CATEGORY_ICON = { social: 'ti-users', posts: 'ti-pencil' }
const CATEGORY_LABEL = { social: 'Social', posts: 'Posts' }

const CATEGORY_PROGRESS = {
    social: (s) => ({
        value: Math.min(s.followers, 50),
        max: 50,
        label: `${s.followers}/50 followers`,
    }),
    posts: (s) => ({ value: Math.min(s.posts, 25), max: 25, label: `${s.posts}/25 posts` }),
}

export default function Badges() {
    const { user: authUser } = useAuth()
    const [stats, setStats] = useState(null)

    useEffect(() => {
        if (!authUser?.userId) return
        fetchStats()
        const interval = setInterval(fetchStats, 5000)
        return () => clearInterval(interval)
    }, [authUser?.userId])

    const fetchStats = async () => {
        try {
            const { data } = await api.get(`/api/users/${authUser?.userId}/gamification`)
            setStats(data)
        } catch (err) {
            console.error(err)
        }
    }

    const userStats = {
        posts: stats?.posts?.count ?? 0,
        likes: stats?.likes?.count ?? 0,
        followers: stats?.followers?.count ?? 0,
    }

    return (
        <div>
            <div
                className="flex items-center gap-3 mb-4 pb-2"
                style={{ borderBottom: '1.5px solid #c5dbc4' }}
            >
                <div
                    className="w-8 h-8 rounded-lg flex items-center justify-center"
                    style={{ background: '#c5dbc4' }}
                >
                    <i
                        className="ti ti-award"
                        aria-hidden="true"
                        style={{ color: '#3b6d11', fontSize: '15px' }}
                    />
                </div>
                <span className="text-sm font-semibold" style={{ color: '#3b6d11' }}>
                    Badges
                </span>
            </div>

            {BADGES.map((section) => {
                const c = CATEGORY_STYLE[section.category]
                const progress = CATEGORY_PROGRESS[section.category]?.(userStats)
                const pct = progress ? Math.round((progress.value / progress.max) * 100) : 0

                return (
                    <div key={section.category} className="mb-6">
                        <div className="flex items-center justify-between gap-2 mb-2">
                            <div className="flex items-center gap-2">
                                <div
                                    className="w-5 h-5 rounded flex items-center justify-center"
                                    style={{ background: c.icon }}
                                >
                                    <i
                                        className={`ti ${CATEGORY_ICON[section.category]}`}
                                        aria-hidden="true"
                                        style={{ color: c.iconText, fontSize: '11px' }}
                                    />
                                </div>
                                <span className="text-xs font-semibold" style={{ color: c.title }}>
                                    {CATEGORY_LABEL[section.category]}
                                </span>
                            </div>
                            {progress && (
                                <span className="text-xs" style={{ color: c.title }}>
                                    {progress.label}
                                </span>
                            )}
                        </div>

                        {progress && (
                            <div
                                className="mb-3 rounded-full overflow-hidden"
                                style={{ background: '#f1efe8', height: '6px' }}
                            >
                                <div
                                    className="h-full rounded-full transition-all"
                                    style={{ width: `${pct}%`, background: c.progressBg }}
                                />
                            </div>
                        )}

                        <div
                            className="grid gap-3"
                            style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))' }}
                        >
                            {section.items.map((badge) => (
                                <div
                                    key={badge.key}
                                    className="rounded-2xl p-4 flex flex-col items-center text-center gap-2 transition-all cursor-pointer hover:scale-105"
                                    style={{
                                        background: 'white',
                                        border: `0.5px solid ${c.cardBorder}`,
                                    }}
                                >
                                    <div
                                        className="w-12 h-12 rounded-full flex items-center justify-center"
                                        style={{ background: c.emojiBg }}
                                    >
                                        <i
                                            className={`ti ${badge.icon}`}
                                            aria-hidden="true"
                                            style={{ color: c.iconText, fontSize: '22px' }}
                                        />
                                    </div>
                                    <div
                                        className="text-xs font-semibold"
                                        style={{ color: c.title }}
                                    >
                                        {badge.name}
                                    </div>
                                    <div
                                        className="text-xs leading-relaxed"
                                        style={{ color: '#888780' }}
                                    >
                                        {badge.desc}
                                    </div>
                                    <span
                                        className="text-xs px-2 py-0.5 rounded-full font-semibold"
                                        style={{ background: c.pillBg, color: c.pillText }}
                                    >
                                        {badge.condition}
                                    </span>
                                </div>
                            ))}
                        </div>
                    </div>
                )
            })}
        </div>
    )
}
