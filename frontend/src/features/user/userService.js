import axiosInstance from '../../services/axiosInstance'

export const getFollowers = async (userId) => {
  const res = await axiosInstance.get(`/api/friends/${userId}/followers`)
  return res.data
}

export const getFollowing = async (userId) => {
  const res = await axiosInstance.get(`/api/friends/${userId}/following`)
  return res.data
}

export const followUser = async (targetId) => {
  const res = await axiosInstance.post(`/api/friends/${targetId}/follow`)
  return res.data
}

export const sendFriendRequest = async (targetId) => {
  const res = await axiosInstance.post(`/api/friends/${targetId}`)
  return res.data
}

export const acceptFriendRequest = async (requesterId) => {
  const res = await axiosInstance.put(`/api/friends/${requesterId}/accept`)
  return res.data
}