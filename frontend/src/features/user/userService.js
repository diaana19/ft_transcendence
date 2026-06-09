import axiosInstance from '../../services/axiosInstance'

// getFollowers returns the followers of a user.
export const getFollowers = async (userId) => {
    const res = await axiosInstance.get(`/api/users/${userId}/followers`)
    return res.data?.data || res.data
}

// getFollowing returns the users that a user follows.
export const getFollowing = async (userId) => {
    const res = await axiosInstance.get(`/api/users/${userId}/following`)
    return res.data?.data || res.data
}

// followUser starts following another user.
export const followUser = async (targetId) => {
    const res = await axiosInstance.post(`/api/friends/follow/${targetId}`)
    return res.data
}

// unfollowUser stops following another user.
export const unfollowUser = async (targetId) => {
    const res = await axiosInstance.delete(`/api/friends/follow/${targetId}`)
    return res.data
}

// sendFriendRequest sends a friend request to a user.
export const sendFriendRequest = async (targetId) => {
    const res = await axiosInstance.post(`/api/friends/request/${targetId}`)
    return res.data
}

// acceptFriendRequest accepts a friend request.
export const acceptFriendRequest = async (requesterId) => {
    const res = await axiosInstance.post(`/api/friends/accept/${requesterId}`)
    return res.data
}

// removeFriend removes a friend.
export const removeFriend = async (friendId) => {
    const res = await axiosInstance.delete(`/api/friends/${friendId}`)
    return res.data
}
