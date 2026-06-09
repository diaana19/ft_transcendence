import { useAuth } from '../../hooks/useAuth'

// MessageList shows the chat bubbles, aligning my own messages to the right.
function MessageList({ messages }) {
    const { user } = useAuth()

    if (messages.length === 0) {
        return <p className="text-center text-gray-400 py-8 text-sm">No messages yet. Say hello!</p>
    }

    return (
        <div className="flex flex-col gap-3 px-4 py-4">
            {messages.map((msg) => {
                const isMine = msg.sender_id === user?.userId

                return (
                    <div
                        key={msg.id}
                        className={`flex flex-col ${isMine ? 'items-end' : 'items-start'}`}
                    >
                        <div
                            className={`max-w-xs px-4 py-2 rounded-2xl text-sm ${
                                isMine ? 'bg-blue-400 text-white' : 'bg-gray-100 text-gray-900'
                            }`}
                        >
                            {msg.content}
                        </div>
                        <span className="text-xs text-gray-400 mt-1">
                            {new Date(msg.created_at).toLocaleTimeString([], {
                                hour: '2-digit',
                                minute: '2-digit',
                            })}
                        </span>
                    </div>
                )
            })}
        </div>
    )
}

export default MessageList
