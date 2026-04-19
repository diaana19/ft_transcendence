import { Routes, Route } from 'react-router-dom'
import RegisterForm from '../features/auth/RegisterForm.jsx'
import LoginForm from '../features/auth/LoginForm.jsx'
import Profile from '../features/user/Profile.jsx'
import Feed from '../features/posts/Feed'

function App() {
    return (
            <Routes>
                <Route path="/register" element={<RegisterForm />} />
                <Route path="/login" element={<LoginForm />} />
                <Route path="/" element={<Feed />} />
                <Route path="/profile" element={<Profile />} />
            </Routes>
    )
}

export default App


