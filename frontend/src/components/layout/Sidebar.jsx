import { NavLink } from "react-router-dom"

export default function Sidebar() {

  const linkClass = ({ isActive }) =>
    `flex items-center gap-3 px-4 py-2 rounded-xl transition font-medium
     ${isActive 
        ? "bg-white/60 text-gray-700" 
        : "text-gray-500 hover:bg-gray-100 hover:text-gray-700"}`
  
  return (
    <aside className="w-64 h-screen fixed left-0 top-0 
                     border-r border-gray-200
                      px-4 py-6">

      {/* Logo */}
      <div className="mb-8 px-2">
        <h1 className="text-xl font-bold text-gray-900">
          Transcendence
        </h1>
        <p className="text-xs text-gray-400">
          University Community
        </p>
      </div>

      {/* Navigation */}
      <nav className="space-y-2">

        <NavLink to="/" className={linkClass}>
          <span>Home</span>
        </NavLink>

        <NavLink to="/profile" className={linkClass}>
          <span>Profile</span>
        </NavLink>

        <NavLink to="/messages" className={linkClass}>
          <span>Messages</span>
        </NavLink>

        <NavLink to="/search" className={linkClass}>
          <span>Search</span>
        </NavLink>

      </nav>

      {/* Bottom user card */}
      <div className="absolute bottom-6 left-0 w-full px-4">
        <div className="p-3 rounded-xl bg-gray-50 border border-gray-200">
          <p className="text-xs text-gray-400">Logged in as</p>
          <p className="text-sm text-gray-900 font-medium truncate">
            user@email.com
          </p>
        </div>
      </div>

    </aside>
  )
}