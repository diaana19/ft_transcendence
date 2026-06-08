import axiosInstance from '../../services/axiosInstance'

export const getFollowers = async (userId) => {
    const res = await axiosInstance.get(`/api/users/${userId}/followers`)
    return res.data?.data || res.data
}

export const getFollowing = async (userId) => {
    const res = await axiosInstance.get(`/api/users/${userId}/following`)
    return res.data?.data || res.data
}

export const followUser = async (targetId) => {
    const res = await axiosInstance.post(`/api/friends/follow/${targetId}`)
    return res.data
}

export const unfollowUser = async (targetId) => {
    const res = await axiosInstance.delete(`/api/friends/follow/${targetId}`)
    return res.data
}

export const sendFriendRequest = async (targetId) => {
    const res = await axiosInstance.post(`/api/friends/request/${targetId}`)
    return res.data
}

export const acceptFriendRequest = async (requesterId) => {
    const res = await axiosInstance.post(`/api/friends/accept/${requesterId}`)
    return res.data
}

export const removeFriend = async (friendId) => {
    const res = await axiosInstance.delete(`/api/friends/${friendId}`)
    return res.data
}
