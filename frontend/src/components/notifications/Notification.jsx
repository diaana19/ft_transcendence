import { useNotifications } from '../../context/NotificationProvider'

function NotificationsPage() {
  const { notifications, markAllRead } = useNotifications()


  return (
    <div className="max-w-2xl mx-auto">

      {/* Header */}
      <div className="sticky top-0 bg-white border-b border-gray-200 px-4 py-3 z-10 flex items-center justify-between">
        <h1 className="text-xl font-bold">
          Notifications
        </h1>

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
        <p className="text-center text-gray-400 py-8">
          No notifications
        </p>
      ) : (
        notifications.map((notif) => (
          <div
            key={notif.id}
            className={`border-b border-gray-200 px-4 py-4 hover:bg-gray-50 transition-colors ${
              notif.read ? 'bg-white' : 'bg-blue-50'
            }`}
          >
            <p className="text-sm text-gray-900">
              {notif.content}
            </p>

            <p className="text-xs text-gray-400 mt-1">
              {new Date(notif.created_at).toLocaleString()}
            </p>
          </div>
        ))
      )}
    </div>
  )
}

export default NotificationsPage