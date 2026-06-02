import { createContext, useContext, useEffect, useState } from 'react'
import { useAuth } from '../hooks/useAuth'
import { useSocket } from './SocketProvider'
import { getUnreadNotifications, markAllNotificationsRead } from '../components/notifications/notificationService'

export const NotificationContext = createContext()

// On mount:
// Load all notifications received while offline (REST), then listen on the
// shared chat socket for real-time notification frames.
export function NotificationProvider({ children }) {
  const { user, loading } = useAuth()
  const { subscribe } = useSocket()
  const [notifications, setNotifications] = useState([])

  useEffect(() => {
    if (loading || !user) return

    getUnreadNotifications()
      .then(data => setNotifications(data ?? []))
      .catch(err => console.error('[Notif] failed to load unread:', err))

    const unsubscribe = subscribe((msg) => {
      if (msg.type === 'notification' && msg.notification) {
        setNotifications(prev => {
          if (prev.some(n => n.id === msg.notification.id)) return prev
          return [msg.notification, ...prev]
        })
      }
    })

    return unsubscribe
  }, [loading, user?.userId, subscribe])

  const markAllRead = async () => {
    try {
      await markAllNotificationsRead()
      setNotifications(prev => prev.map(n => ({ ...n, read: true })))
    } catch (err) {
      console.error('[Notif] failed to mark all read:', err)
    }
  }

  const unreadCount = notifications.filter(n => !n.read).length

  return (
    <NotificationContext.Provider value={{ notifications, setNotifications, markAllRead, unreadCount }}>
      {children}
    </NotificationContext.Provider>
  )
}

export function useNotifications() {
  return useContext(NotificationContext)
}
