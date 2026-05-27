import axiosInstance from '../../services/axiosInstance'

export async function getUnreadNotifications() {
  const res = await axiosInstance.get('/api/notification')

  return res.data.data
}

export async function markAllNotificationsRead() {
  const res = await axiosInstance.patch('/api/notification/read')

  return res.data
}