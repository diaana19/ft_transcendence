/*
** File: MessagesList.jsx
** Description: List of all users to start a conversation with
** Responsibilities:
** - Fetch all users from API
** - Display user list
** - Navigate to chat on click
*/

import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import axiosInstance from '../../services/axiosInstance'
import { useAuth } from '../../hooks/useAuth'

function MessagesList() {
	const [users, setUsers] = useState([])
	const [loading, setLoading] = useState(true)
	const navigate = useNavigate()
	const { user: currentUser } = useAuth()

	useEffect(() => {
		const fetchUsers = async () => {
			try {
				const res = await axiosInstance.get('/api/users')
				const filtered = res.data.filter(u => u.id !== currentUser?.userId)
				setUsers(filtered)
			} catch (err) {
				console.error('Error fetching users:', err)
			} finally {
				setLoading(false)
			}
		}
		fetchUsers()
	}, [])

	return (
		<div className="max-w-2xl mx-auto">
			<div className="sticky top-0 bg-white border-b border-gray-200 px-4 py-3 z-10">
				<h1 className="text-xl font-bold">Messages</h1>
			</div>
			{loading ? (
				<p className="text-center text-gray-400 py-8 text-sm">Loading...</p>
			) : users.length === 0 ? (
				<p className="text-center text-gray-400 py-8 text-sm">No users found</p>
			) : (
				users.map((u) => (
					<div
						key={u.id}
						onClick={() => navigate(`/messages/${u.id}`)}
						className="flex items-center gap-3 px-4 py-4 border-b border-gray-200 hover:bg-gray-50 cursor-pointer transition-colors"
					>
						<div className="w-10 h-10 rounded-full bg-gray-300 overflow-hidden flex items-center justify-center text-white font-bold text-sm">
							{u.avatar
								? <img src={u.avatar} alt={u.username} className="w-full h-full object-cover" />
								: u.username?.[0]?.toUpperCase()
							}
						</div>
						<div className="flex-1">
							<p className="text-sm font-bold text-gray-900">{u.username}</p>
							<p className="text-xs text-gray-400">{u.name || ''}</p>
						</div>
					</div>
				))
			)}
		</div>
	)
}

export default MessagesList