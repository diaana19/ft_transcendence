import { useEffect, useState } from 'react'
import { useAuth } from '../../hooks/useAuth'
import { getFollowers, getFollowing } from '../../features/user/userService'
import FollowButton from '../../components/common/FollowButton'

function FriendsList({ userId }) {
  const { user: currentUser } = useAuth()
  const [tab, setTab]           = useState('followers') 
  const [followers, setFollowers] = useState([])
  const [following, setFollowing] = useState([])
  const [loading, setLoading]     = useState(false)


  const followingIds = new Set(following.map((u) => u.id))

  useEffect(() => {
    fetchAll()
  }, [userId])

  const fetchAll = async () => {
    setLoading(true)
    try {
      const [followersData, followingData] = await Promise.all([
        getFollowers(userId),
        getFollowing(userId),
      ])
      setFollowers(followersData)
      setFollowing(followingData)
    } catch (err) {
      console.error('Error fetching friends:', err)
    } finally {
      setLoading(false)
    }
  }

  const list = tab === 'followers' ? followers : following

  return (
    <div className="max-w-xl mx-auto">

      {/* Tabs */}
      <div className="flex border-b border-gray-200">
        <button
          onClick={() => setTab('followers')}
          className={`flex-1 py-4 text-sm font-semibold transition-colors ${
            tab === 'followers'
              ? 'border-b-2 border-blue-400 text-black'
              : 'text-gray-500 hover:text-black'
          }`}
        >
          Followers
          <span className="ml-2 text-gray-400 font-normal">{followers.length}</span>
        </button>

        <button
          onClick={() => setTab('following')}
          className={`flex-1 py-4 text-sm font-semibold transition-colors ${
            tab === 'following'
              ? 'border-b-2 border-blue-400 text-black'
              : 'text-gray-500 hover:text-black'
          }`}
        >
          Following
          <span className="ml-2 text-gray-400 font-normal">{following.length}</span>
        </button>
      </div>

      {/* List */}
      {loading ? (
        <p className="text-center text-gray-400 py-8">Loading...</p>
      ) : list.length === 0 ? (
        <p className="text-center text-gray-400 py-8">
          {tab === 'followers' ? 'No followers yet' : 'Not following anyone'}
        </p>
      ) : (
        list.map((person) => (
          <div
            key={person.id}
            className="flex items-center gap-3 px-4 py-3 border-b border-gray-100 hover:bg-gray-50"
          >
            {/* Avatar */}
            <img
              src={person.avatar || `https://i.pravatar.cc/150?u=${person.id}`}
              alt={person.username}
              className="w-10 h-10 rounded-full object-cover"
            />

            {/* Info */}
            <div className="flex-1 min-w-0">
              <p className="font-bold text-black text-sm truncate">{person.name}</p>
              <p className="text-gray-500 text-sm truncate">@{person.username}</p>
            </div>

            {/* Follow button */}
            {currentUser?.userId !== person.id && (
              <FollowButton
                targetId={person.id}
                isFollowing={followingIds.has(person.id)}
              />
            )}
          </div>
        ))
      )}
    </div>
  )
}

export default FriendsList