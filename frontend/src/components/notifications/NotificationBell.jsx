/*
** File: NotificationBell.jsx
** Description: Notification bell icon with unread count badge
** Responsibilities:
** - Display bell icon in sidebar or navbar
** - Show unread notification count badge
** - Navigate to notifications page on click
*/

import { useNavigate } from 'react-router-dom'
import { useNotifications } from '../../context/NotificationProvider'
import { BellIcon } from '@heroicons/react/24/outline'

export default function NotificationBell() {
  const { notifications } = useNotifications()
  const navigate = useNavigate()
  const unreadCount = notifications.filter(n => !n.read).length

  return (
    <button
      onClick={() => navigate('/notifications')}
      className="relative flex items-center gap-3 px-4 py-2 rounded-xl transition font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700 w-full"
    >
      <div className="relative">
        <BellIcon className="w-5 h-5" />
        {unreadCount > 0 && (
          <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-4 h-4 flex items-center justify-center font-bold">
            {unreadCount > 9 ? '9+' : unreadCount}
          </span>
        )}
      </div>
      <span>Notifications</span>
    </button>
  )
}