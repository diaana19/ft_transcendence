import { Outlet, useLocation } from 'react-router-dom'
import Sidebar from './Sidebar'
import RightSidebar from './RightSidebar'

// Layout is the page shell with the sidebars and the routed content in the middle.
export default function Layout() {
    const location = useLocation()
    // Show the right sidebar only on the home feed.
    const showRightSidebar =
        location.pathname === '/' ||
        location.pathname.startsWith('/tag/') ||
        location.pathname === '/tags' ||
        location.pathname === '/users'

    return (
        <div className="flex min-h-screen text-gray-900">
            <Sidebar />

            <main
                className={`flex-1 ml-0 md:ml-64 min-w-0 pb-16 md:pb-0 ${showRightSidebar ? 'xl:mr-80' : ''}`}
            >
                <Outlet />
            </main>
            {showRightSidebar && (
                <div className="hidden xl:block">
                    <RightSidebar />
                </div>
            )}
        </div>
    )
}
