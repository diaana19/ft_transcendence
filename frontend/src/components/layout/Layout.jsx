import { Outlet, useLocation } from "react-router-dom"
import Sidebar from "./Sidebar"
import RightSidebar from "./RightSidebar"

export default function Layout() {
	const location = useLocation()

	const showRightSidebar = location.pathname === "/"
  return (
    <div className="flex min-h-screen text-gray-900">

      <Sidebar />

      {/* Content */}
      <main className={`ml-64 flex-1 ${showRightSidebar ? "mr-80" : "mr-0"}`}>
        <Outlet />
      </main>

      {showRightSidebar && <RightSidebar />}
    </div>
  )
}
