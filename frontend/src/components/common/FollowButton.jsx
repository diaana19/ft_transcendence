import { useState } from 'react'
import { followUser } from '../../features/user/userService'

function FollowButton({ targetId, isFollowing: initialIsFollowing }) {
  const [isFollowing, setIsFollowing] = useState(initialIsFollowing)
  const [loading, setLoading]         = useState(false)

  const handleFollow = async () => {
    setLoading(true)
    try {
      await followUser(targetId)
      setIsFollowing(true)
    } catch (err) {
      console.error('Error following user:', err)
    } finally {
      setLoading(false)
    }
  }

  if (isFollowing) {
    return (
      <button
        disabled
        className="px-4 py-1.5 rounded-full border border-gray-300 text-gray-500 text-sm font-semibold"
      >
        Following
      </button>
    )
  }

  return (
    <button
      onClick={handleFollow}
      disabled={loading}
      className="px-4 py-1.5 rounded-full bg-black text-white text-sm font-bold hover:bg-gray-800 disabled:opacity-50 transition-colors"
    >
      {loading ? '...' : 'Follow'}
    </button>
  )
}

export default FollowButton