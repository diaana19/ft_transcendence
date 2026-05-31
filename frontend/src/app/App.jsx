import { Routes, Route } from 'react-router-dom'
import Layout from '../components/layout/Layout'
import RequireAuth from '../components/common/RequireAuth'
import PostPage from '../features/posts/PostPage.jsx'
import RegisterForm from '../features/auth/RegisterForm.jsx'
import LoginForm from '../features/auth/LoginForm.jsx'
import TwoFAVerify from '../features/auth/TwoFAVerify.jsx'
import TwoFASetup from '../features/auth/TwoFASetup.jsx'
import Profile from '../features/user/Profile.jsx'
import Feed from '../features/posts/Feed'
import NotificationsPage from '../components/notifications/Notification.jsx'
import ChatRoom from '../features/chat/ChatRoom'
import MessagesList from '../features/chat/MessagesList'
import CookiesBanner from '../components/common/CookiesBanner'
import HelpCenter from '../features/help/HelpCenter.jsx'
import SettingsPage from '../features/settings/SettingsPage.jsx'

function App() {
    return (
        <>
            <CookiesBanner />
            <Routes>
                {/* PUBLIC */}
                <Route element={<Layout />}>
                    <Route path="/" element={<Feed />} />
                    <Route path="/help" element={<HelpCenter />} />
                </Route>
                <Route path="/register" element={<RegisterForm />} />
                <Route path="/login" element={<LoginForm />} />
                <Route path="/login/2fa" element={<TwoFAVerify />} />

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
                    <Route path="/2fa/setup" element={<TwoFASetup />} />
                    <Route path="/messages" element={<MessagesList />} />
                    <Route path="/messages/:peerId" element={<ChatRoom />} />
                    <Route path="/settings" element={<SettingsPage />} />
                </Route>
            </Routes>
        </>
    )
}

export default App
