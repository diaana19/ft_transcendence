import { createContext, useContext, useEffect, useRef, useCallback } from 'react'
import { useAuth } from '../hooks/useAuth'

const SocketContext = createContext(null)

// SocketProvider opens one WebSocket for the logged user and shares subscribe/send.
export function SocketProvider({ children }) {
    const { token, user, loading } = useAuth()
    const wsRef = useRef(null)
    const handlersRef = useRef(new Set())
    const queueRef = useRef([])

    // Open the WebSocket when a user is logged, and close it on cleanup.
    useEffect(() => {
        if (loading || !user) return

        // Use wss when the page is https, ws otherwise.
        const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
        const tokenParam = token ? `?token=${encodeURIComponent(token)}` : ''
        const ws = new WebSocket(`${protocol}://${window.location.host}/api/ws/chat${tokenParam}`)
        wsRef.current = ws

        // When connected, flush the frames that were queued while offline.
        ws.onopen = () => {
            queueRef.current.forEach((frame) => ws.send(frame))
            queueRef.current = []
        }

        // Parse each message and call every subscribed handler.
        ws.onmessage = (event) => {
            let msg
            try {
                msg = JSON.parse(event.data)
            } catch {
                return
            }
            handlersRef.current.forEach((handler) => {
                try {
                    handler(msg)
                } catch (err) {
                    console.info('[Socket] handler error:', err)
                }
            })
        }

        ws.onerror = (e) => console.info('[Socket] WebSocket error:', e)
        ws.onclose = () => {
            wsRef.current = null
        }

        return () => {
            ws.close()
            wsRef.current = null
        }
    }, [loading, user?.userId, token])

    // subscribe registers a message handler and returns a function to remove it.
    const subscribe = useCallback((handler) => {
        handlersRef.current.add(handler)
        return () => handlersRef.current.delete(handler)
    }, [])

    // send writes a JSON message, or queues it when the socket is not open yet.
    const send = useCallback((obj) => {
        const frame = JSON.stringify(obj)
        const ws = wsRef.current
        if (ws && ws.readyState === WebSocket.OPEN) {
            ws.send(frame)
        } else {
            queueRef.current.push(frame)
        }
    }, [])

    return <SocketContext.Provider value={{ subscribe, send }}>{children}</SocketContext.Provider>
}

// useSocket gives access to the socket context.
export function useSocket() {
    return useContext(SocketContext)
}
