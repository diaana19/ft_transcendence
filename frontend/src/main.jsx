import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { AuthProvider } from './context/AuthProvider'
import { SocketProvider } from './context/SocketProvider'
import { NotificationProvider } from './context/NotificationProvider'
import App from './app/App.jsx'
import './styles/index.css'
import { registerSW } from 'virtual:pwa-register'

// Register the PWA service worker right away.
registerSW({
    immediate: true,
})

// Mount the app wrapped by the router and all context providers.

ReactDOM.createRoot(document.getElementById('root')).render(
    <React.StrictMode>
        <BrowserRouter>
            <AuthProvider>
                <SocketProvider>
                    <NotificationProvider>
                        <App />
                    </NotificationProvider>
                </SocketProvider>
            </AuthProvider>
        </BrowserRouter>
    </React.StrictMode>
)
