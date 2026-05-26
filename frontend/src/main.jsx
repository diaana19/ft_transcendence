import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { AuthProvider } from './context/AuthProvider'
import { NotificationProvider } from './context/NotificationProvider'
import App from './app/App.jsx'
import './styles/index.css'
import { registerSW } from 'virtual:pwa-register'

registerSW({
  immediate: true,
})

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
  <BrowserRouter>
      <AuthProvider>
        <NotificationProvider>
        <App />
        </NotificationProvider>
      </AuthProvider>
    </BrowserRouter>
  </React.StrictMode>
)