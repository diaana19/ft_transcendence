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

// getFriends returns the accepted friends of a user.
export const getFriends = async (userId) => {
    const res = await axiosInstance.get(`/api/users/${userId}/friends`)
    return res.data?.data || res.data
}

// searchUsers returns a page of users matching the query plus the total count.
export const searchUsers = async (query = '', limit = 20, offset = 0) => {
    const q = encodeURIComponent(query)
    const res = await axiosInstance.get(`/api/users?q=${q}&limit=${limit}&offset=${offset}`)
    return { users: res.data.data || [], total: res.data.total ?? 0 }
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

//Cancels a pending friend request by removing the relationship.
export const cancelFriendRequest = async (targetId) => {
    const res = await axiosInstance.delete(`/api/friends/${targetId}`)
    return res.data
}
