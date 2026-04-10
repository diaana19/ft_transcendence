import { BrowserRouter, Routes, Route } from 'react-router-dom'
import RegisterForm from '../features/auth/RegisterForm.jsx'
import LoginForm from '../features/auth/LoginForm.jsx'

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/register" element={<RegisterForm />} />
                <Route path="/login" element={<LoginForm />} />
            </Routes>
        </BrowserRouter>
    )
}

export default App