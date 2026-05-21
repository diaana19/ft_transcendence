import { Routes, Route } from 'react-router-dom'
import Layout from '../components/layout/Layout'
import RequireAuth from '../components/common/RequireAuth'

import PostPage from '../features/posts/PostPage.jsx'
import RegisterForm from '../features/auth/RegisterForm.jsx'
import LoginForm from '../features/auth/LoginForm.jsx'
import Profile from '../features/user/Profile.jsx'
import Feed from '../features/posts/Feed'
import NotificationsPage from '../components/notifications/Notification.jsx'

function App() {
    return (
            <Routes>

                {/* PUBLIC */}
                <Route element={<Layout />}>
                    <Route path="/" element={<Feed />} />
                </Route>

                <Route path="/register" element={<RegisterForm />} />
                <Route path="/login" element={<LoginForm />} />

                {/* PRIVATE */}
                <Route
                    element={
                    <RequireAuth>
                        <Layout />
                    </RequireAuth>
                    }
                >
                    <Route path="/profile" element={<Profile />} />
                    <Route path="/profile/:id" element={<Profile />} />
                    <Route path="/notifications" element={<NotificationsPage />} />
                    <Route path="/post/:id" element={<PostPage />} />
                </Route>

            </Routes>
    )
}

export default App


