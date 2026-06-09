import { createContext, useContext, useEffect, useState } from 'react'
import { useAuth } from '../hooks/useAuth'
import { useSocket } from './SocketProvider'
import {
    getUnreadNotifications,
    markAllNotificationsRead,
    markNotificationRead,
} from '../components/notifications/notificationService'

export const NotificationContext = createContext()
// NotificationProvider loads unread notifications and listens new ones over the socket.
export function NotificationProvider({ children }) {
    const { user, loading } = useAuth()
    const { subscribe } = useSocket()
    const [notifications, setNotifications] = useState([])

    useEffect(() => {
        if (loading || !user) return

        // First load the current unread notifications from the API.
        getUnreadNotifications()
            .then((data) => setNotifications(data ?? []))
            .catch((err) => console.info('[Notif] failed to load unread:', err))

        // Then subscribe to the socket to receive new notifications in real time.
        const unsubscribe = subscribe((msg) => {
            if (msg.type === 'notification' && msg.notification) {
                setNotifications((prev) => {
                    // Skip if we already have this notification.
                    if (prev.some((n) => n.id === msg.notification.id)) return prev
                    return [msg.notification, ...prev]
                })
            }
        })

        return unsubscribe
    }, [loading, user?.userId, subscribe])

    // markAllRead marks every notification as read on the server and locally.
    const markAllRead = async () => {
        try {
            await markAllNotificationsRead()
            setNotifications((prev) => prev.map((n) => ({ ...n, read: true })))
        } catch (err) {
            console.info('[Notif] failed to mark all read:', err)
        }
    }

    // markRead marks one notification as read. It updates the UI first then the server.
    const markRead = async (id) => {
        setNotifications((prev) => prev.map((n) => (n.id === id ? { ...n, read: true } : n)))
        try {
            await markNotificationRead(id)
        } catch (err) {
            console.info('[Notif] failed to mark read:', err)
            // Revert the local change if the request fails.
            setNotifications((prev) => prev.map((n) => (n.id === id ? { ...n, read: false } : n)))
        }
    }

    const unreadCount = notifications.filter((n) => !n.read).length

    return (
        <NotificationContext.Provider
            value={{ notifications, setNotifications, markAllRead, markRead, unreadCount }}
        >
            {children}
        </NotificationContext.Provider>
    )
}

// useNotifications gives access to the notification context.
export function useNotifications() {
    return useContext(NotificationContext)
}
