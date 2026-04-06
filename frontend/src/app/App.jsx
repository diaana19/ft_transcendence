import { BrowserRouter, Routes, Route } from 'react-router-dom'
import RegisterForm from '../features/auth/RegisterForm.jsx'

function App() {
    return (
        <BrowserRouter>
         <Routes>
            <Route path="/register" element={<RegisterForm />} />
         </Routes>
        </BrowserRouter>
    )
}

export default App