import axiosInstance from '../../services/axiosInstance'

// getUnreadNotifications returns the unread notifications of the user.
export async function getUnreadNotifications() {
    const res = await axiosInstance.get('/api/notification')

    return res.data.data
}

// markAllNotificationsRead marks every notification as read.
export async function markAllNotificationsRead() {
    const res = await axiosInstance.patch('/api/notification/read')

    return res.data
}

// markNotificationRead marks one notification as read.
export async function markNotificationRead(id) {
    const res = await axiosInstance.patch(`/api/notification/${id}/read`)

    return res.data
}
