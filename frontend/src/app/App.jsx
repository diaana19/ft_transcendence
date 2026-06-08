import { Routes, Route } from 'react-router-dom'
import Layout from '../components/layout/Layout'
import RequireAuth from '../components/common/RequireAuth'
import PostPage from '../features/posts/PostPage.jsx'
import RegisterForm from '../features/auth/RegisterForm.jsx'
import LoginForm from '../features/auth/LoginForm.jsx'
import ForgotPasswordForm from '../features/auth/ForgotPasswordForm.jsx'
import TwoFAVerify from '../features/auth/TwoFAVerify.jsx'
import TwoFASetup from '../features/auth/TwoFASetup.jsx'
import Profile from '../features/user/Profile.jsx'
import Feed from '../features/posts/Feed'
import TagFeed from '../features/posts/TagFeed.jsx'
import NotificationsPage from '../components/notifications/Notification.jsx'
import CookiesBanner from '../components/common/CookiesBanner'
import HelpCenter from '../features/help/HelpCenter.jsx'
import SettingsPage from '../features/settings/SettingsPage.jsx'
import MessagesPage from '../features/chat/MessagesPage'
import Gamification from '../features/gamification/Gamification.jsx'
import TermsOfService from '../features/auth/TermsofService.jsx'
import PrivacyPolicy from '../features/auth/PrivacyPolicy.jsx'

function App() {
    return (
        <>
            <CookiesBanner />
            <Routes>
                <Route element={<Layout />}>
                    <Route path="/" element={<Feed />} />
                    <Route path="/tag/:tag" element={<TagFeed />} />
                    <Route path="/help" element={<HelpCenter />} />
                    <Route path="/badge" element={<Gamification />} />
                </Route>
                <Route path="/register" element={<RegisterForm />} />
                <Route path="/login" element={<LoginForm />} />
                <Route path="/forgot-password" element={<ForgotPasswordForm />} />
                <Route path="/login/2fa" element={<TwoFAVerify />} />
                <Route path="/terms" element={<TermsOfService />} />
                <Route path="/privacy" element={<PrivacyPolicy />} />

                <Route
                    element={
                        <RequireAuth>
                            <Layout />
                        </RequireAuth>
                    }
                >
                    <Route path="/profile" element={<Profile />} />
                    <Route path="/profile/:id" element={<Profile />} />
                    <Route path="/profile/u/:username" element={<Profile />} />
                    <Route path="/notifications" element={<NotificationsPage />} />
                    <Route path="/post/:id" element={<PostPage />} />
                    <Route path="/2fa/setup" element={<TwoFASetup />} />
                    <Route path="/messages" element={<MessagesPage />} />
                    <Route path="/messages/:peerId" element={<MessagesPage />} />
                    <Route path="/settings" element={<SettingsPage />} />
                </Route>
            </Routes>
        </>
    )
}

export default App
