import { useEffect, useRef, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import axiosInstance from '../../services/axiosInstance'
import { useAuth } from '../../hooks/useAuth'
import { useSocket } from '../../context/SocketProvider'
import MessageList from './MessageList'
import MessageInput from './MessageInput'

// MessagesPage is the chat page: a user list on the left and the conversation on the right.
export default function MessagesPage() {
    const { peerId } = useParams()
    const navigate = useNavigate()
    const { user: currentUser } = useAuth()
    const { subscribe, send } = useSocket()

    const [users, setUsers] = useState([])
    const [search, setSearch] = useState('')
    const [messages, setMessages] = useState([])
    const [fetching, setFetching] = useState(false)
    const [selectedUser, setSelectedUser] = useState(null)
    const [onlineUsers, setOnlineUsers] = useState({})
    const [userOrder, setUserOrder] = useState([])
    const lastMessageRef = useRef(null)

    // fetchOnlineStatus asks the online state of every user in the list.
    const fetchOnlineStatus = async (usersList) => {
        const statuses = {}
        await Promise.all(
            usersList.map(async (u) => {
                try {
                    const res = await axiosInstance.get(`/api/users/${u.id}/online`)
                    statuses[u.id] = res.data?.online ?? false
                } catch {
                    statuses[u.id] = false
                }
            })
        )
        setOnlineUsers(statuses)
    }

    // Load the user list once (without me) and their online status.
    useEffect(() => {
        const fetchUsers = async () => {
            try {
                const res = await axiosInstance.get('/api/users')
                const filtered = res.data.filter((u) => u.id !== currentUser?.userId)
                setUsers(filtered)
                setUserOrder(filtered.map((u) => u.id))
                fetchOnlineStatus(filtered)
            } catch (err) {
                console.error(err)
            }
        }
        fetchUsers()
    }, [])

    // Keep the selected user in sync with the peerId in the URL.
    useEffect(() => {
        if (peerId && users.length > 0) {
            const found = users.find((u) => u.id === peerId)
            setSelectedUser(found || null)
        }
    }, [peerId, users])

    // Open the conversation: ask for history and listen new messages on the socket.
    useEffect(() => {
        if (!peerId) return
        setMessages([])
        setFetching(true)
        // Tell the server which conversation we opened so it sends the history.
        send({ action: 'open', peer_id: peerId })

        const me = currentUser?.userId
        const unsubscribe = subscribe((msg) => {
            // History frame: replace the whole message list.
            if (msg.type === 'history' && msg.peer_id === peerId) {
                setMessages(msg.messages || [])
                setFetching(false)
                return
            }
            if (msg.type === 'message' && msg.message) {
                const m = msg.message
                // Keep only messages that belong to this conversation.
                const belongs =
                    (m.sender_id === peerId && m.recipient_id === me) ||
                    (m.sender_id === me && m.recipient_id === peerId)
                if (!belongs) return
                setMessages((prev) => (prev.some((x) => x.id === m.id) ? prev : [...prev, m]))
                // Move this peer to the top of the conversation list.
                const conversationPeer = m.sender_id === me ? m.recipient_id : m.sender_id
                setUserOrder((prev) => [
                    conversationPeer,
                    ...prev.filter((id) => id !== conversationPeer),
                ])
            }
        })
        return unsubscribe
    }, [peerId, subscribe, send, currentUser?.userId])

    // Scroll to the last message whenever the list changes.
    useEffect(() => {
        lastMessageRef.current?.scrollIntoView({ behavior: 'smooth' })
    }, [messages])

    // handleSend sends a new message over the socket and bumps the peer to the top.
    const handleSend = (content) => {
        send({ action: 'message', recipient_id: peerId, content })
        setUserOrder((prev) => [peerId, ...prev.filter((id) => id !== peerId)])
    }

    const filteredUsers = users
        .filter(
            (u) =>
                u.username?.toLowerCase().includes(search.toLowerCase()) ||
                u.displayname?.toLowerCase().includes(search.toLowerCase())
        )
        .sort((a, b) => {
            const ai = userOrder.indexOf(a.id)
            const bi = userOrder.indexOf(b.id)
            if (ai === -1) return 1
            if (bi === -1) return -1
            return ai - bi
        })

    return (
        <div className="flex h-screen overflow-hidden">
            <div
                className={`${peerId ? 'hidden md:flex' : 'flex'} w-full md:w-80 flex-col`}
                style={{
                    borderRight: '1px solid #ede8fd',
                    background: 'linear-gradient(135deg, #faf5ff 0%, #eef2ff 100%)',
                }}
            >
                <div className="px-4 py-4" style={{ borderBottom: '1px solid #ede8fd' }}>
                    <h1 className="text-xl font-bold mb-3" style={{ color: '#2c2c2a' }}>
                        Messages
                    </h1>
                    <input
                        type="text"
                        placeholder="Search..."
                        value={search}
                        onChange={(e) => setSearch(e.target.value)}
                        className="w-full px-3 py-2 text-sm rounded-full focus:outline-none"
                        style={{
                            background: '#f5f3ff',
                            border: '1px solid #ede8fd',
                            color: '#2c2c2a',
                        }}
                    />
                </div>
                <div className="flex-1 overflow-y-auto">
                    {filteredUsers.length === 0 ? (
                        <p className="text-center py-8 text-sm" style={{ color: '#afa9ec' }}>
                            No users found
                        </p>
                    ) : (
                        filteredUsers.map((u) => (
                            <div
                                key={u.id}
                                onClick={() => navigate(`/messages/${u.id}`)}
                                className="flex items-center gap-3 px-4 py-3 cursor-pointer transition-colors"
                                style={{
                                    background:
                                        peerId === u.id
                                            ? 'rgba(237, 232, 253, 0.8)'
                                            : 'transparent',
                                    borderRight:
                                        peerId === u.id
                                            ? '3px solid #534ab7'
                                            : '3px solid transparent',
                                }}
                            >
                                <div className="relative w-10 h-10 flex-shrink-0">
                                    <div
                                        className="w-10 h-10 rounded-full overflow-hidden flex items-center justify-center font-bold text-sm"
                                        style={{ background: '#ede8fd', color: '#534ab7' }}
                                    >
                                        {u.avatar ? (
                                            <img
                                                src={u.avatar}
                                                alt={u.username}
                                                className="w-full h-full object-cover"
                                            />
                                        ) : (
                                            u.username?.[0]?.toUpperCase()
                                        )}
                                    </div>
                                    {onlineUsers[u.id] && (
                                        <div
                                            className="absolute bottom-0 right-0 w-3 h-3 rounded-full border-2 border-white"
                                            style={{ background: '#4ade80' }}
                                        />
                                    )}
                                </div>
                                <div className="flex-1 min-w-0">
                                    <p
                                        className="text-sm font-bold truncate"
                                        style={{ color: '#2c2c2a' }}
                                    >
                                        {u.username}
                                    </p>
                                    <p className="text-xs truncate" style={{ color: '#afa9ec' }}>
                                        {u.displayname || ''}
                                    </p>
                                </div>
                            </div>
                        ))
                    )}
                </div>
            </div>

            <div
                className={`${peerId ? 'flex' : 'hidden md:flex'} flex-1 flex-col`}
                style={{ background: '#faf5ff' }}
            >
                {peerId && selectedUser ? (
                    <>
                        <div
                            className="px-4 py-3 flex items-center gap-3 bg-white"
                            style={{ borderBottom: '1px solid #ede8fd' }}
                        >
                            <button
                                onClick={() => navigate('/messages')}
                                className="md:hidden mr-1 text-sm font-bold"
                                style={{ color: '#534ab7' }}
                            >
                                ←
                            </button>
                            <div className="relative w-9 h-9 flex-shrink-0">
                                <div
                                    className="w-9 h-9 rounded-full overflow-hidden flex items-center justify-center font-bold text-sm"
                                    style={{ background: '#ede8fd', color: '#534ab7' }}
                                >
                                    {selectedUser.avatar ? (
                                        <img
                                            src={selectedUser.avatar}
                                            alt={selectedUser.username}
                                            className="w-full h-full object-cover"
                                        />
                                    ) : (
                                        selectedUser.username?.[0]?.toUpperCase()
                                    )}
                                </div>
                                {onlineUsers[selectedUser.id] && (
                                    <div
                                        className="absolute bottom-0 right-0 w-3 h-3 rounded-full border-2 border-white"
                                        style={{ background: '#4ade80' }}
                                    />
                                )}
                            </div>
                            <div
                                onClick={() => navigate(`/profile/${selectedUser.id}`)}
                                className="cursor-pointer hover:underline"
                            >
                                <p className="text-sm font-bold" style={{ color: '#2c2c2a' }}>
                                    {selectedUser.username}
                                </p>
                                <p
                                    className="text-xs"
                                    style={{
                                        color: onlineUsers[selectedUser.id] ? '#4ade80' : '#afa9ec',
                                    }}
                                >
                                    {onlineUsers[selectedUser.id]
                                        ? 'Online'
                                        : selectedUser.displayname || ''}
                                </p>
                            </div>
                        </div>

                        <div className="flex-1 overflow-y-auto">
                            {fetching ? (
                                <p
                                    className="text-center py-8 text-sm"
                                    style={{ color: '#afa9ec' }}
                                >
                                    Loading...
                                </p>
                            ) : (
                                <MessageList messages={messages} />
                            )}
                            <div ref={lastMessageRef} />
                        </div>

                        <div style={{ borderTop: '1px solid #ede8fd', background: 'white' }}>
                            <MessageInput onSend={handleSend} loading={false} />
                        </div>
                    </>
                ) : (
                    <div className="flex-1 flex flex-col items-center justify-center">
                        <p className="text-4xl mb-4">💬</p>
                        <p className="text-lg font-medium" style={{ color: '#534ab7' }}>
                            Select a conversation
                        </p>
                        <p className="text-sm mt-1" style={{ color: '#afa9ec' }}>
                            Choose a user from the list to start chatting
                        </p>
                    </div>
                )}
            </div>
        </div>
    )
}
