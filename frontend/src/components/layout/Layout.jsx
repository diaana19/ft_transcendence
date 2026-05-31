import { Outlet, useLocation } from "react-router-dom"
import Sidebar from "./Sidebar"
import RightSidebar from "./RightSidebar"

export default function Layout() {
  const location = useLocation()
  const showRightSidebar = location.pathname === "/"

  return (
    <div className="flex min-h-screen text-gray-900">
      <Sidebar />
      {/* Content — se adapta según sidebar */}
      <main className={`flex-1 ml-0 md:ml-64 min-w-0 ${showRightSidebar ? 'xl:mr-80' : ''}`}>
        <Outlet />
      </main>
      {/* RightSidebar solo en home y pantallas grandes */}
      {showRightSidebar && (
        <div className="hidden xl:block">
          <RightSidebar />
        </div>
      )}
    </div>
  )
}
