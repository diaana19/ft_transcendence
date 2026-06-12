import { useEffect, useState } from 'react'
import { useAuth } from '../../hooks/useAuth'
import { getFollowers, getFollowing, getFriends } from '../../features/user/userService'

// FriendsList shows the follower and following counts of a user.
function FriendsList({ userId }) {
    const [followers, setFollowers] = useState([])
    const [following, setFollowing] = useState([])
    const [friends, setFriends] = useState([])

    useEffect(() => {
        if (!userId) return
        getFollowers(userId)
            .then(setFollowers)
            .catch(() => {})
        getFollowing(userId)
            .then(setFollowing)
            .catch(() => {})
        getFriends(userId)
            .then(setFriends)
            .catch(() => {})
    }, [userId])

    return (
        <div className="flex gap-4 mt-2 pb-3">
            <span className="text-sm text-gray-700 cursor-pointer hover:underline">
                <strong>{following.length}</strong>
                <span className="text-gray-500 ml-1">Following</span>
            </span>
            <span className="text-sm text-gray-700 cursor-pointer hover:underline">
                <strong>{followers.length}</strong>
                <span className="text-gray-500 ml-1">Followers</span>
            </span>
            <span className="text-sm text-gray-700 cursor-pointer hover:underline">
                <strong>{friends.length}</strong>
                <span className="text-gray-500 ml-1">Friends</span>
            </span>
        </div>
    )
}

export default FriendsList
