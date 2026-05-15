import axiosInstance from '../../services/axiosInstance'

export async function getUnreadNotifications() {
  const res = await api.get('/notifications/unread')
  return res.data.data
}

export async function markAllNotificationsRead() {
  const res = await api.put('/notifications/read-all')
  return res.data
}