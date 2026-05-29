/*
** File: ChatRoom.jsx
** Description: Main chat room component for direct messages between users
** Responsibilities:
** - Fetch conversation history with a specific user
** - Poll for new messages every 3 seconds
** - Send messages via API
** - Render MessageList and MessageInput
*/

import { useEffect, useRef, useState } from 'react'
import { useParams } from 'react-router-dom'
import axiosInstance from '../../services/axiosInstance'
import MessageList from './MessageList'
import MessageInput from './MessageInput'

function ChatRoom({ peerUsername }) {
  const { peerId } = useParams()
  const [messages, setMessages] = useState([])
  const [loading, setLoading] = useState(false)
  const [fetching, setFetching] = useState(true)
  const lastMessageRef = useRef(null)
  const sinceRef = useRef(null)

  const fetchConversation = async () => {
    try {
      const res = await axiosInstance.get(`/api/chat/messages?with=${peerId}`)
      const msgs = res.data?.messages || []
      setMessages(msgs)
      if (msgs.length > 0) {
        sinceRef.current = msgs[msgs.length - 1].created_at
      }
    } catch (err) {
      console.error('Error fetching conversation:', err)
    } finally {
      setFetching(false)
    }
  }

  const pollMessages = async () => {
    try {
      const since = sinceRef.current || ''
      const res = await axiosInstance.get(`/api/chat/poll?since=${since}`)
      const newMsgs = res.data?.messages || []
      if (newMsgs.length > 0) {
        setMessages((prev) => [...prev, ...newMsgs])
        sinceRef.current = newMsgs[newMsgs.length - 1].created_at
      }
    } catch (err) {
      console.error('Error polling messages:', err)
    }
  }

  const handleSend = async (content) => {
    setLoading(true)
    try {
      const res = await axiosInstance.post('/api/chat/messages', {
        recipient_id: peerId,
        content,
      })
      setMessages((prev) => [...prev, res.data])
      sinceRef.current = res.data.created_at
    } catch (err) {
      console.error('Error sending message:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (!peerId) return
    setMessages([])
    setFetching(true)
    sinceRef.current = null
    fetchConversation()
    const interval = setInterval(fetchConversation, 3000)
    return () => clearInterval(interval)
  }, [peerId])

  useEffect(() => {
    lastMessageRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  return (
    <div className="flex flex-col h-full">
      <div className="sticky top-0 bg-white border-b border-gray-200 px-4 py-3 z-10">
        <h2 className="text-lg font-bold text-gray-900">
          {peerUsername || 'Chat'}
        </h2>
      </div>
      <div className="flex-1 overflow-y-auto">
        {fetching ? (
          <p className="text-center text-gray-400 py-8 text-sm">Loading...</p>
        ) : (
          <MessageList messages={messages} />
        )}
        <div ref={lastMessageRef} />
      </div>
      <MessageInput onSend={handleSend} loading={loading} />
    </div>
  )
}

export default ChatRoom