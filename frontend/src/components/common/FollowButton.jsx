import { useState } from 'react'
import { followUser, unfollowUser } from '../../features/user/userService'

// FollowButton toggles follow/unfollow for a user and shows the current state.
function FollowButton({ targetId, isFollowing: initialIsFollowing, onFollow }) {
    const [isFollowing, setIsFollowing] = useState(initialIsFollowing)
    const [loading, setLoading] = useState(false)
    const [hovered, setHovered] = useState(false)

    const handleClick = async () => {
        setLoading(true)
        try {
            if (isFollowing) {
                await unfollowUser(targetId)
                setIsFollowing(false)
            } else {
                await followUser(targetId)
                setIsFollowing(true)
            }
            onFollow?.()
        } catch (err) {
            console.info('Error toggling follow:', err)
        } finally {
            setLoading(false)
        }
    }

    if (isFollowing) {
        return (
            <button
                onClick={handleClick}
                disabled={loading}
                onMouseEnter={() => setHovered(true)}
                onMouseLeave={() => setHovered(false)}
                className="px-4 py-1.5 rounded-full text-sm font-semibold transition-colors disabled:opacity-50"
                style={{
                    border: hovered ? '1px solid #fde8f0' : '1px solid #ede8fd',
                    color: hovered ? '#d4537e' : '#534ab7',
                    background: 'white',
                }}
            >
                {loading ? '...' : hovered ? 'Unfollow' : 'Following ✓'}
            </button>
        )
    }

    return (
        <button
            onClick={handleClick}
            disabled={loading}
            className="px-4 py-1.5 rounded-full text-sm font-bold transition-colors disabled:opacity-50"
            style={{ background: '#534ab7', color: 'white', border: 'none' }}
        >
            {loading ? '...' : 'Follow'}
        </button>
    )
}

export default FollowButton
