import { NavLink, useNavigate } from "react-router-dom"
import { useAuth } from "../../hooks/useAuth"
import {
  HomeIcon,
  UserIcon,
  SparklesIcon,
  QuestionMarkCircleIcon,
  EnvelopeIcon,
  BellIcon,
  ArrowRightOnRectangleIcon
} from "@heroicons/react/24/outline"
import NotificationBell from '../notifications/NotificationBell'

export default function Sidebar() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/login', { replace: true })
  }

  const linkClass = ({ isActive }) =>
    `flex items-center gap-3 px-4 py-2 rounded-xl transition font-medium
    ${isActive
      ? "bg-white/60 text-gray-700"
      : "text-gray-500 hover:bg-gray-100 hover:text-gray-700"}`

  return (
    <>
      {/* Sidebar desktop — oculta en móvil */}
      <aside className="hidden md:flex flex-col w-64 h-screen fixed left-0 top-0 border-r border-gray-200 px-4 py-6">
        {/* Logo */}
        <div className="mb-8 px-2 flex items-center gap-3">
          <img src="/logo.png" alt="Synk logo" className="w-10 h-10 object-contain" />
          <div className="flex flex-col leading-tight">
            <h1 className="text-xl font-bold text-gray-900">Synk</h1>
            <p className="text-xs text-gray-400">Join your people</p>
          </div>
        </div>

        {/* Navigation */}
        <nav className="space-y-2">
          <NavLink to="/" className={linkClass}>
            <HomeIcon className="w-5 h-5" />
            <span>Home</span>
          </NavLink>
          <NavLink to="/profile" className={linkClass}>
            <UserIcon className="w-5 h-5" />
            <span>Profile</span>
          </NavLink>
          <NavLink to="/comunities" className={linkClass}>
            <SparklesIcon className="w-5 h-5" />
            <span>Comunities</span>
          </NavLink>
          {user?.userId && <NotificationBell />}
          {user?.userId && (
            <NavLink to="/messages" className={linkClass}>
              <EnvelopeIcon className="w-5 h-5" />
              <span>Messages</span>
            </NavLink>
          )}
          <NavLink to="/help" className={linkClass}>
            <QuestionMarkCircleIcon className="w-5 h-5" />
            <span>Help</span>
          </NavLink>
        </nav>

        {/* Bottom user card */}
        {user?.userId && (
          <div className="absolute bottom-6 left-0 w-full px-4">
            <div className="p-3 rounded-xl bg-gray-50 border border-gray-200">
              <p className="text-xs text-gray-400">Logged in as @{user.username}</p>
              <p className="text-sm text-gray-900 font-medium truncate mb-2">{user.email}</p>
              <button
                type="button"
                onClick={handleLogout}
                className="w-full flex items-center justify-center gap-2 px-3 py-1.5 rounded-lg text-sm text-gray-600 hover:text-red-600 hover:bg-red-50 transition border border-gray-200"
              >
                <ArrowRightOnRectangleIcon className="w-4 h-4" />
                <span>Log out</span>
              </button>
            </div>
          </div>
        )}
      </aside>

      {/* Bottom bar móvil — solo en móvil */}
      <nav className="md:hidden fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 flex justify-around items-center py-2 z-50">
        <NavLink to="/" className={({ isActive }) => `flex flex-col items-center p-2 ${isActive ? 'text-blue-400' : 'text-gray-500'}`}>
          <HomeIcon className="w-6 h-6" />
        </NavLink>
        <NavLink to="/profile" className={({ isActive }) => `flex flex-col items-center p-2 ${isActive ? 'text-blue-400' : 'text-gray-500'}`}>
          <UserIcon className="w-6 h-6" />
        </NavLink>
        {user?.userId && (
          <NavLink to="/notifications" className={({ isActive }) => `flex flex-col items-center p-2 ${isActive ? 'text-blue-400' : 'text-gray-500'}`}>
            <BellIcon className="w-6 h-6" />
          </NavLink>
        )}
        {user?.userId && (
          <NavLink to="/messages" className={({ isActive }) => `flex flex-col items-center p-2 ${isActive ? 'text-blue-400' : 'text-gray-500'}`}>
            <EnvelopeIcon className="w-6 h-6" />
          </NavLink>
        )}
      </nav>
    </>
  )
}