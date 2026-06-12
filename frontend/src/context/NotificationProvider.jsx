import { createContext, useContext, useEffect, useState } from 'react'
import { useAuth } from '../hooks/useAuth'
import { useSocket } from './SocketProvider'
import {
    getNotifications,
    deleteNotification,
    deleteAllNotifications,
} from '../components/notifications/notificationService'

export const NotificationContext = createContext()

const STORAGE_KEY = 'notifications'

function readCache() {
    try {
        const raw = localStorage.getItem(STORAGE_KEY)
        return raw ? JSON.parse(raw) : []
    } catch {
        return []
    }
}

// NotificationProvider loads recent notifications and listens new ones over the socket.
export function NotificationProvider({ children }) {
    const { user, loading } = useAuth()
    const { subscribe } = useSocket()
    const [notifications, setNotifications] = useState(readCache)

    useEffect(() => {
        try {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(notifications))
        } catch {
            void 0
        }
    }, [notifications])

    useEffect(() => {
        if (loading) return
        if (!user) {
            setNotifications([])
            return
        }

        getNotifications()
            .then((data) => setNotifications(data ?? []))
            .catch((err) => console.info('[Notif] failed to load notifications:', err))

        // Then subscribe to the socket to receive new notifications in real time.
        const unsubscribe = subscribe((msg) => {
            if (msg.type === 'notification' && msg.notification) {
                setNotifications((prev) => {
                    // Skip if we already have this notification.
                    if (prev.some((n) => n.id === msg.notification.id)) return prev
                    return [msg.notification, ...prev]
                })
            }
            if (msg.type === 'notification_removed' && msg.notification_id) {
                setNotifications((prev) => prev.filter((n) => n.id !== msg.notification_id))
            }
            if (msg.type === 'notifications_cleared') {
                setNotifications([])
            }
        })

        return unsubscribe
    }, [loading, user?.userId, subscribe])

    // removeNotification deletes one read notification. It updates the UI first
    // then the server.
    const removeNotification = async (id) => {
        const previous = notifications
        setNotifications((prev) => prev.filter((n) => n.id !== id))
        try {
            await deleteNotification(id)
        } catch (err) {
            console.info('[Notif] failed to delete:', err)
            // Restore the list if the request fails.
            setNotifications(previous)
        }
    }

    // clearAll deletes every notification. It updates the UI first then the server.
    const clearAll = async () => {
        const previous = notifications
        setNotifications([])
        try {
            await deleteAllNotifications()
        } catch (err) {
            console.info('[Notif] failed to clear:', err)
            // Restore the list if the request fails.
            setNotifications(previous)
        }
    }

    const unreadCount = notifications.length

    return (
        <NotificationContext.Provider
            value={{
                notifications,
                setNotifications,
                removeNotification,
                clearAll,
                unreadCount,
            }}
        >
            {children}
        </NotificationContext.Provider>
    )
}

// useNotifications gives access to the notification context.
export function useNotifications() {
    return useContext(NotificationContext)
}
