import {createContext, useContext, useEffect, useState } from 'react'
import {getUnreadNotifications, markAllNotificationsRead } from '../components/notifications/notificationService'

export const NotificationContext = createContext()

export function NotificationProvider({ children }) {
  const [notifications, setNotifications] = useState([])

  return (
    <NotificationContext.Provider
      value={{
        notifications,
        setNotifications
      }}
    >
      {children}
    </NotificationContext.Provider>
  )
}

export function useNotifications() {
  return useContext(NotificationContext)
}