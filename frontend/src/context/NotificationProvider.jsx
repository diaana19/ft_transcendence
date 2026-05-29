/*
** File: NotificationProvider.jsx
** Description: Global notification context provider
** Responsibilities:
** - Fetch unread notifications from API on mount
** - Poll for new notifications every 30 seconds
** - Provide markAllRead function to children
*/

import { createContext, useContext, useEffect, useState } from 'react'
import { getUnreadNotifications, markAllNotificationsRead } from '../components/notifications/notificationService'
import { useAuth } from '../hooks/useAuth'

export const NotificationContext = createContext()

export function NotificationProvider({ children }) {
  const { user } = useAuth()
  const [notifications, setNotifications] = useState([])

  const fetchNotifications = async () => {
    try {
      const data = await getUnreadNotifications()
      setNotifications(data || [])
    } catch (err) {
      console.error('Error fetching notifications:', err)
    }
  }

  const markAllRead = async () => {
    try {
      await markAllNotificationsRead()
      setNotifications([])
    } catch (err) {
      console.error('Error marking notifications read:', err)
    }
  }

  useEffect(() => {
    if (!user) return
    fetchNotifications()
    const interval = setInterval(fetchNotifications, 5000)
    return () => clearInterval(interval)
  }, [user])

  return (
    <NotificationContext.Provider value={{ notifications, markAllRead }}>
      {children}
    </NotificationContext.Provider>
  )
}

export function useNotifications() {
  return useContext(NotificationContext)
}