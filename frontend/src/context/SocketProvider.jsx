import { createContext, useContext, useEffect, useRef, useCallback } from 'react'
import { useAuth } from '../hooks/useAuth'

const SocketContext = createContext(null)

export function SocketProvider({ children }) {
  const { token, user, loading } = useAuth()
  const wsRef = useRef(null)
  const handlersRef = useRef(new Set())
  const queueRef = useRef([])

  useEffect(() => {
    if (loading || !user) return

    const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
    const tokenParam = token ? `?token=${encodeURIComponent(token)}` : ''
    const ws = new WebSocket(`${protocol}://${window.location.host}/api/ws/chat${tokenParam}`)
    wsRef.current = ws

    ws.onopen = () => {
      queueRef.current.forEach(frame => ws.send(frame))
      queueRef.current = []
    }

    ws.onmessage = (event) => {
      let msg
      try {
        msg = JSON.parse(event.data)
      } catch {
        return // non-JSON frame — ignore
      }
      handlersRef.current.forEach(handler => {
        try {
          handler(msg)
        } catch (err) {
          console.error('[Socket] handler error:', err)
        }
      })
    }

    ws.onerror = (e) => console.error('[Socket] WebSocket error:', e)
    ws.onclose = () => { wsRef.current = null }

    return () => {
      ws.close()
      wsRef.current = null
    }
  }, [loading, user?.userId, token])

  const subscribe = useCallback((handler) => {
    handlersRef.current.add(handler)
    return () => handlersRef.current.delete(handler)
  }, [])

  const send = useCallback((obj) => {
    const frame = JSON.stringify(obj)
    const ws = wsRef.current
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(frame)
    } else {
      queueRef.current.push(frame)
    }
  }, [])

  return (
    <SocketContext.Provider value={{ subscribe, send }}>
      {children}
    </SocketContext.Provider>
  )
}

export function useSocket() {
  return useContext(SocketContext)
}
