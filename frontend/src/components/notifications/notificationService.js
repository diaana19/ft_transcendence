import axiosInstance from '../../services/axiosInstance'

// getNotifications returns the recent notifications of the user, read and unread.
export async function getNotifications() {
    const res = await axiosInstance.get('/api/notification')

    return res.data.data
}

// deleteNotification deletes one notification once the user read it.
export async function deleteNotification(id) {
    const res = await axiosInstance.delete(`/api/notification/${id}`)

    return res.data
}

// deleteAllNotifications deletes every notification of the user.
export async function deleteAllNotifications() {
    const res = await axiosInstance.delete('/api/notification')

    return res.data
}
