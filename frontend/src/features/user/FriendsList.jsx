import { useEffect, useState } from 'react'
import { getFollowers, getFollowing, getFriends } from '../../features/user/userService'
import FollowersModal from './FollowersModal'

function FriendsList({ userId }) {
    const [followers, setFollowers] = useState([])
    const [following, setFollowing] = useState([])
    const [friends, setFriends] = useState([])
    const [modal, setModal] = useState(null)

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
        <div className="flex gap-4 mt-3">
            <span className="text-sm cursor-pointer hover:underline">
                <strong style={{ color: '#534ab7' }}>{following.length}</strong>
                <span
                    onClick={() => setModal('following')}
                    className="ml-1"
                    style={{ color: '#757575' }}
                >
                    Following
                </span>
            </span>
            <span className="text-sm cursor-pointer hover:underline">
                <strong style={{ color: '#534ab7' }}>{followers.length}</strong>
                <span
                    onClick={() => setModal('followers')}
                    className="ml-1"
                    style={{ color: '#757575' }}
                >
                    Followers
                </span>
            </span>
            <span className="text-sm cursor-pointer hover:underline">
                <strong style={{ color: '#534ab7' }}>{friends.length}</strong>
                <span
                    onClick={() => setModal('friends')}
                    className="ml-1"
                    style={{ color: '#757575' }}
                >
                    Friends
                </span>
            </span>

            {modal && (
                <FollowersModal userId={userId} type={modal} onClose={() => setModal(null)} />
            )}
        </div>
    )
}

export default FriendsList
