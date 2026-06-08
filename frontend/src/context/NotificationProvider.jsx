import { createContext, useContext, useEffect, useState } from 'react'
import { useAuth } from '../hooks/useAuth'
import { useSocket } from './SocketProvider'
import { getUnreadNotifications, markAllNotificationsRead, markNotificationRead } from '../components/notifications/notificationService'

export const NotificationContext = createContext()
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

  const markRead = async (id) => {
    setNotifications(prev => prev.map(n => (n.id === id ? { ...n, read: true } : n)))
    try {
      await markNotificationRead(id)
    } catch (err) {
      console.error('[Notif] failed to mark read:', err)
      setNotifications(prev => prev.map(n => (n.id === id ? { ...n, read: false } : n)))
    }
  }

  const unreadCount = notifications.filter(n => !n.read).length

  return (
    <NotificationContext.Provider value={{ notifications, setNotifications, markAllRead, markRead, unreadCount }}>
      {children}
    </NotificationContext.Provider>
  )
}

export function useNotifications() {
  return useContext(NotificationContext)
}
