import { Outlet } from "react-router-dom"
import Sidebar from "./Sidebar"
import RightSidebar from "./RightSidebar"

export default function Layout() {
  return (
    <div className="flex min-h-screen text-gray-900">

      <Sidebar />

      {/* Content */}
      <main className="ml-64 mr-80 flex-1 p-6">
        <Outlet />
      </main>

      <RightSidebar />
    </div>
  )
}