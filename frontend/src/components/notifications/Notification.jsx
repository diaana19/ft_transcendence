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
  if (type === 'like') return (
    <div className="w-8 h-8 rounded-full flex items-center justify-center" style={{ background: '#fde8f0' }}>
      <Heart size={14} style={{ color: '#d4537e' }} />
    </div>
  )
  if (type === 'comment') return (
    <div className="w-8 h-8 rounded-full flex items-center justify-center" style={{ background: '#e8f0fd' }}>
      <MessageCircle size={14} style={{ color: '#534ab7' }} />
    </div>
  )
  if (type === 'friend_request') return (
    <div className="w-8 h-8 rounded-full flex items-center justify-center" style={{ background: '#f0fde8' }}>
      <UserPlus size={14} style={{ color: '#4a7b34' }} />
    </div>
  )
  return (
    <div className="w-8 h-8 rounded-full flex items-center justify-center" style={{ background: '#ede8fd' }}>
      <Bell size={14} style={{ color: '#534ab7' }} />
    </div>
  )
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
      <div className="sticky top-0 bg-white px-4 py-3 z-10 flex items-center justify-between" style={{ borderBottom: '1px solid #ede8fd' }}>
        <h1 className="text-xl font-bold" style={{ color: '#2c2c2a' }}>Notifications</h1>
        {notifications.length > 0 && (
          <button
            onClick={markAllRead}
            className="text-xs font-semibold px-3 py-1.5 rounded-full transition-colors"
            style={{ background: '#ede8fd', color: '#534ab7' }}
          >
            Mark all read
          </button>
        )}
      </div>

      {/* List */}
      {notifications.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-16">
          <div className="w-16 h-16 rounded-full flex items-center justify-center mb-4" style={{ background: '#ede8fd' }}>
            <Bell size={28} style={{ color: '#534ab7' }} />
          </div>
          <p className="text-base font-medium" style={{ color: '#534ab7' }}>No notifications yet</p>
          <p className="text-sm mt-1" style={{ color: '#afa9ec' }}>We'll notify you when something happens</p>
        </div>
      ) : (
        notifications.map((notif) => (
          <div
            key={notif.id}
            onClick={() => handleClick(notif)}
            className="flex items-start gap-3 px-4 py-4 cursor-pointer transition-colors"
            style={{
              borderBottom: '1px solid #ede8fd',
              background: notif.read ? 'white' : 'rgba(237, 232, 253, 0.3)',
            }}
          >
            <div className="mt-1 flex-shrink-0">
              {getNotifIcon(notif.type)}
            </div>
            <div className="flex-1">
              <p className="text-sm" style={{ color: '#2c2c2a' }}>{notif.content}</p>
              <p className="text-xs mt-1" style={{ color: '#afa9ec' }}>
                {new Date(notif.created_at).toLocaleString()}
              </p>
              {notif.type === 'friend_request' && (
                <button
                  onClick={(e) => handleAccept(e, notif)}
                  disabled={accepted[notif.id]}
                  className="mt-2 px-3 py-1 text-xs font-bold rounded-full disabled:opacity-50 transition-colors"
                  style={{ background: '#534ab7', color: 'white' }}
                >
                  {accepted[notif.id] ? 'Accepted ✓' : 'Accept'}
                </button>
              )}
            </div>
            {!notif.read && (
              <div className="w-2 h-2 rounded-full mt-2 flex-shrink-0" style={{ background: '#534ab7' }} />
            )}
          </div>
        ))
      )}
    </div>
  )
}

export default NotificationsPage