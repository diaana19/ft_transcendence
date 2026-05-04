import { Routes, Route } from 'react-router-dom'
import Layout from '../components/layout/Layout'

import RegisterForm from '../features/auth/RegisterForm.jsx'
import LoginForm from '../features/auth/LoginForm.jsx'
import Profile from '../features/user/Profile.jsx'
import Feed from '../features/posts/Feed'

function App() {
    return (
            <Routes>
                <Route path="/register" element={<RegisterForm />} />
                <Route path="/login" element={<LoginForm />} />
                
                <Route element={<Layout />}>
                <Route path="/" element={<Feed />} />
                <Route path="/profile" element={<Profile />} />
                <Route path="/profile/:id" element={<Profile />} />
                </Route>
            </Routes>
    )
}

export default App


