/*
** File: MessagesPage.jsx
** Description: Full chat page with conversation list and active chat
** Responsibilities:
** - Display list of users to chat with on the left
** - Display active conversation on the right
** - Handle sending and receiving messages
*/

import { useEffect, useRef, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import axiosInstance from '../../services/axiosInstance'
import { useAuth } from '../../hooks/useAuth'
import { useSocket } from '../../context/SocketProvider'
import MessageList from './MessageList'
import MessageInput from './MessageInput'

export default function MessagesPage() {
	const { peerId } = useParams()
	const navigate = useNavigate()
	const { user: currentUser } = useAuth()
	const { subscribe, send } = useSocket()

	const [users, setUsers] = useState([])
	const [messages, setMessages] = useState([])
	const [fetching, setFetching] = useState(false)
	const [selectedUser, setSelectedUser] = useState(null)
	const lastMessageRef = useRef(null)

	useEffect(() => {
		const fetchUsers = async () => {
			try {
				const res = await axiosInstance.get('/api/users')
				const filtered = res.data.filter(u => u.id !== currentUser?.userId)
				setUsers(filtered)
			} catch (err) {
				console.error(err)
			}
		}
		fetchUsers()
	}, [])

	useEffect(() => {
		if (peerId && users.length > 0) {
			const found = users.find(u => u.id === peerId)
			setSelectedUser(found || null)
		}
	}, [peerId, users])

	// Open the conversation over the socket: the server replies with a "history"
	// frame and pushes each new "message" in real time. No polling.
	useEffect(() => {
		if (!peerId) return
		setMessages([])
		setFetching(true)
		send({ action: 'open', peer_id: peerId })

		const me = currentUser?.userId
		const unsubscribe = subscribe((msg) => {
			if (msg.type === 'history' && msg.peer_id === peerId) {
				setMessages(msg.messages || [])
				setFetching(false)
				return
			}
			if (msg.type === 'message' && msg.message) {
				const m = msg.message
				const belongs =
					(m.sender_id === peerId && m.recipient_id === me) ||
					(m.sender_id === me && m.recipient_id === peerId)
				if (!belongs) return
				setMessages(prev => (prev.some(x => x.id === m.id) ? prev : [...prev, m]))
			}
		})
		return unsubscribe
	}, [peerId, subscribe, send, currentUser?.userId])

	useEffect(() => {
		lastMessageRef.current?.scrollIntoView({ behavior: 'smooth' })
	}, [messages])

	// The sent message is echoed back over the socket (we joined the channel on
	// open), so we render it on arrival rather than appending optimistically.
	const handleSend = (content) => {
		send({ action: 'message', recipient_id: peerId, content })
	}

	return (
		<div className="flex h-screen overflow-hidden">

			{/* Left — user list: oculto en móvil si hay chat abierto */}
			<div className={`${peerId ? 'hidden md:flex' : 'flex'} w-full md:w-80 flex-col`}
				style={{ borderRight: '1px solid #ede8fd', background: 'linear-gradient(135deg, #faf5ff 0%, #eef2ff 100%)' }}>
				<div className="px-4 py-4" style={{ borderBottom: '1px solid #ede8fd' }}>
					<h1 className="text-xl font-bold" style={{ color: '#2c2c2a' }}>Messages</h1>
				</div>
				<div className="flex-1 overflow-y-auto">
					{users.length === 0 ? (
						<p className="text-center py-8 text-sm" style={{ color: '#afa9ec' }}>No users found</p>
					) : (
						users.map((u) => (
							<div
								key={u.id}
								onClick={() => navigate(`/messages/${u.id}`)}
								className="flex items-center gap-3 px-4 py-3 cursor-pointer transition-colors"
								style={{
									background: peerId === u.id ? 'rgba(237, 232, 253, 0.8)' : 'transparent',
									borderRight: peerId === u.id ? '3px solid #534ab7' : '3px solid transparent',
								}}
							>
								<div className="w-10 h-10 rounded-full overflow-hidden flex-shrink-0 flex items-center justify-center font-bold text-sm"
									style={{ background: '#ede8fd', color: '#534ab7' }}
								>
									{u.avatar
										? <img src={u.avatar} alt={u.username} className="w-full h-full object-cover" />
										: u.username?.[0]?.toUpperCase()
									}
								</div>
								<div className="flex-1 min-w-0">
									<p className="text-sm font-bold truncate" style={{ color: '#2c2c2a' }}>{u.username}</p>
									<p className="text-xs truncate" style={{ color: '#afa9ec' }}>{u.name || ''}</p>
								</div>
							</div>
						))
					)}
				</div>
			</div>

			{/* Right — chat: oculto en móvil si no hay chat abierto */}
			<div className={`${peerId ? 'flex' : 'hidden md:flex'} flex-1 flex-col`} style={{ background: '#faf5ff' }}>
				{peerId && selectedUser ? (
					<>
						{/* Chat header con botón back en móvil */}
						<div className="px-4 py-3 flex items-center gap-3 bg-white" style={{ borderBottom: '1px solid #ede8fd' }}>
							<button
								onClick={() => navigate('/messages')}
								className="md:hidden mr-1 text-sm font-bold"
								style={{ color: '#534ab7' }}
							>
								←
							</button>
							<div className="w-9 h-9 rounded-full overflow-hidden flex items-center justify-center font-bold text-sm"
								style={{ background: '#ede8fd', color: '#534ab7' }}
							>
								{selectedUser.avatar
									? <img src={selectedUser.avatar} alt={selectedUser.username} className="w-full h-full object-cover" />
									: selectedUser.username?.[0]?.toUpperCase()
								}
							</div>
							<div>
								<p className="text-sm font-bold" style={{ color: '#2c2c2a' }}>{selectedUser.username}</p>
								<p className="text-xs" style={{ color: '#afa9ec' }}>{selectedUser.name || ''}</p>
							</div>
						</div>

						{/* Messages */}
						<div className="flex-1 overflow-y-auto">
							{fetching ? (
								<p className="text-center py-8 text-sm" style={{ color: '#afa9ec' }}>Loading...</p>
							) : (
								<MessageList messages={messages} />
							)}
							<div ref={lastMessageRef} />
						</div>

						{/* Input */}
						<div style={{ borderTop: '1px solid #ede8fd', background: 'white' }}>
							<MessageInput onSend={handleSend} loading={false} />
						</div>
					</>
				) : (
					<div className="flex-1 flex flex-col items-center justify-center">
						<p className="text-4xl mb-4">💬</p>
						<p className="text-lg font-medium" style={{ color: '#534ab7' }}>Select a conversation</p>
						<p className="text-sm mt-1" style={{ color: '#afa9ec' }}>Choose a user from the list to start chatting</p>
					</div>
				)}
			</div>
		</div>
	)
}