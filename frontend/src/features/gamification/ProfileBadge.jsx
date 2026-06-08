import { useEffect, useState } from 'react'
import api from '../../services/axiosInstance'

const BADGES = [
    {
        key: 'welcome',
        icon: 'ti-user-plus',
        name: 'Welcome',
        condition: 'Day one',
        check: (s) => true,
    },
    {
        key: 'first_bond',
        icon: 'ti-heart-handshake',
        name: 'First bond',
        condition: '1 friend',
        check: (s) => s.followers >= 1,
    },
    {
        key: 'network_builder',
        icon: 'ti-topology-star',
        name: 'Network builder',
        condition: '10 friends',
        check: (s) => s.followers >= 10,
    },
    {
        key: 'social_butterfly',
        icon: 'ti-butterfly',
        name: 'Social butterfly',
        condition: '50 friends',
        check: (s) => s.followers >= 50,
    },
    {
        key: 'first_words',
        icon: 'ti-feather',
        name: 'First words',
        condition: '1 post',
        check: (s) => s.posts >= 1,
    },
    {
        key: 'first_spark',
        icon: 'ti-heart',
        name: 'First spark',
        condition: '1 like',
        check: (s) => s.likes >= 1,
    },
    {
        key: 'content_creator',
        icon: 'ti-writing',
        name: 'Content creator',
        condition: '25 posts',
        check: (s) => s.posts >= 25,
    },
    {
        key: 'rising_star',
        icon: 'ti-star',
        name: 'Rising star',
        condition: '20 likes',
        check: (s) => s.likes >= 20,
    },
]

export default function ProfileBadges({ userId }) {
    const [stats, setStats] = useState(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        if (!userId) return
        api.get(`/api/users/${userId}/gamification`)
            .then(({ data }) => setStats(data))
            .catch(console.error)
            .finally(() => setLoading(false))
    }, [userId])

    if (loading)
        return (
            <p className="text-center py-8 text-sm" style={{ color: '#b4b2a9' }}>
                Loading...
            </p>
        )

    const userStats = {
        posts: stats?.posts?.count ?? 0,
        likes: stats?.likes?.count ?? 0,
        followers: stats?.followers?.count ?? 0,
        messages: 0,
    }

    const earned = BADGES.filter((b) => b.check(userStats))

    if (earned.length === 0)
        return (
            <p className="text-center py-8 text-sm" style={{ color: '#b4b2a9' }}>
                No badges yet
            </p>
        )

    return (
        <div className="px-4 py-4">
            <div className="flex flex-wrap gap-3">
                {earned.map((badge) => (
                    <div
                        key={badge.key}
                        className="flex flex-col items-center gap-1 p-3 rounded-2xl transition-all cursor-pointer hover:scale-105"
                        style={{
                            background: 'white',
                            border: '0.5px solid #ede8fd',
                            minWidth: '80px',
                        }}
                        title={badge.condition}
                    >
                        <div
                            className="w-10 h-10 rounded-full flex items-center justify-center"
                            style={{ background: '#ede8fd' }}
                        >
                            <i
                                className={`ti ${badge.icon}`}
                                aria-hidden="true"
                                style={{ color: '#534ab7', fontSize: '18px' }}
                            />
                        </div>
                        <span
                            className="text-xs font-semibold text-center"
                            style={{ color: '#534ab7' }}
                        >
                            {badge.name}
                        </span>
                    </div>
                ))}
            </div>
        </div>
    )
}
