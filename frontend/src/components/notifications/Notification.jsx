/*
** File: Notification.jsx
** Description: Notifications page showing likes and comments on user posts
** Responsibilities:
** - Display list of unread notifications
** - Show notification type (like, comment, friend request)
** - Allow user to mark all notifications as read
** - Navigate to relevant content on click
** - Accept friend requests inline
*/

import { useNavigate } from 'react-router-dom'
import { useNotifications } from '../../context/NotificationProvider'
import { Heart, MessageCircle, UserPlus, Bell } from 'lucide-react'
import { acceptFriendRequest } from '../../features/user/userService'
import { useState } from 'react'

function getNotifIcon(type) {
  if (type === 'like') return <Heart size={16} className="text-red-400" />
  if (type === 'comment') return <MessageCircle size={16} className="text-blue-400" />
  if (type === 'friend_request') return <UserPlus size={16} className="text-green-400" />
  return <Bell size={16} className="text-gray-400" />
}

function NotificationsPage() {
  const { notifications, markAllRead } = useNotifications()
  const navigate = useNavigate()
  const [accepted, setAccepted] = useState({})

  const handleClick = (notif) => {
    if (notif.type === 'like' || notif.type === 'comment') {
      navigate(`/profile/${notif.actor_id}`)
    } else if (notif.type === 'friend_request') {
      navigate(`/profile/${notif.actor_id}`)
    }
  }

  const handleAccept = async (e, notif) => {
    e.stopPropagation()
    try {
      await acceptFriendRequest(notif.actor_id)
      setAccepted(prev => ({ ...prev, [notif.id]: true }))
      await markAllRead() 
    } catch (err) {
      console.error('Error accepting friend request:', err)
    }
  }

  return (
    <div className="max-w-2xl mx-auto">
      {/* Header */}
      <div className="sticky top-0 bg-white border-b border-gray-200 px-4 py-3 z-10 flex items-center justify-between">
        <h1 className="text-xl font-bold">Notifications</h1>
        {notifications.length > 0 && (
          <button
            onClick={markAllRead}
            className="text-sm text-blue-500 hover:underline"
          >
            Mark all read
          </button>
        )}
      </div>

      {/* List */}
      {notifications.length === 0 ? (
        <p className="text-center text-gray-400 py-8 text-sm">No notifications</p>
      ) : (
        notifications.map((notif) => (
          <div
            key={notif.id}
            onClick={() => handleClick(notif)}
            className={`flex items-start gap-3 border-b border-gray-200 px-4 py-4 hover:bg-gray-50 transition-colors cursor-pointer ${notif.read ? 'bg-white' : 'bg-blue-50'
              }`}
          >
            <div className="mt-1">
              {getNotifIcon(notif.type)}
            </div>
            <div className="flex-1">
              <p className="text-sm text-gray-900">{notif.content}</p>
              <p className="text-xs text-gray-400 mt-1">
                {new Date(notif.created_at).toLocaleString()}
              </p>
              {notif.type === 'friend_request' && (
                <button
                  onClick={(e) => handleAccept(e, notif)}
                  disabled={accepted[notif.id]}
                  className="mt-2 px-3 py-1 bg-blue-400 hover:bg-blue-500 text-white text-xs font-bold rounded-full disabled:opacity-50 transition-colors"
                >
                  {accepted[notif.id] ? 'Accepted ✓' : 'Accept'}
                </button>
              )}
            </div>
          </div>
        ))
      )}
    </div>
  )
}

export default NotificationsPage