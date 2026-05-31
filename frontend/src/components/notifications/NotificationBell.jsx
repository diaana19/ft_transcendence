import { useNotifications } from '../../context/NotificationProvider'

export default function NotificationBell() {
  const { unreadCount } = useNotifications()

  if (unreadCount === 0) return null

  return (
    <span className="ml-auto bg-red-500 text-white text-xs font-bold px-1.5 py-0.5 rounded-full min-w-[18px] text-center">
      {unreadCount > 99 ? '99+' : unreadCount}
    </span>
  )
}
