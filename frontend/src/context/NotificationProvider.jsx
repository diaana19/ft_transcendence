import { createContext, useContext, useEffect, useRef, useState } from 'react'
import { useAuth } from '../hooks/useAuth'
import { getUnreadNotifications, markAllNotificationsRead } from '../components/notifications/notificationService'

export const NotificationContext = createContext()

// On mount:
// Load all notification while offline
// Then Connect to WebSocket to get notification in real time
export function NotificationProvider({ children }) {
  const { token, user, loading } = useAuth()
  const [notifications, setNotifications] = useState([])
  const wsRef = useRef(null)

  useEffect(() => {
    if (loading || !user) return

    getUnreadNotifications()
      .then(data => setNotifications(data ?? []))
      .catch(err => console.error('[Notif] failed to load unread:', err))

    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const tokenParam = token ? `?token=${encodeURIComponent(token)}` : ''
    const ws = new WebSocket(`${protocol}://${window.location.host}/api/ws/chat${tokenParam}`)
    wsRef.current = ws

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        if (msg.type === 'notification' && msg.notification) {
          setNotifications(prev => {
            if (prev.some(n => n.id === msg.notification.id)) return prev
            return [msg.notification, ...prev]
          })
        }
      } catch {
        // non-JSON frames — ignore
      }
    }

    ws.onerror = (e) => console.error('[Notif] WebSocket error:', e)
    ws.onclose = () => { wsRef.current = null }

    return () => {
      ws.close()
      wsRef.current = null
    }
  }, [loading, user?.userId, token])

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
