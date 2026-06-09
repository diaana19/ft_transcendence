import { useState } from 'react'

// MessageInput is the chat text box. Enter sends, Shift+Enter makes a new line.
function MessageInput({ onSend, loading }) {
    const [content, setContent] = useState('')

    const handleSubmit = () => {
        if (!content.trim()) return
        onSend(content)
        setContent('')
    }

    const handleKeyDown = (e) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault()
            handleSubmit()
        }
    }

    return (
        <div className="flex items-center gap-2 border-t border-gray-200 p-3 bg-white">
            <input
                type="text"
                value={content}
                onChange={(e) => setContent(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Write a message..."
                className="flex-1 border border-gray-300 rounded-full px-4 py-2 text-sm text-gray-900 placeholder-gray-400 focus:outline-none focus:border-blue-400"
            />
            <button
                onClick={handleSubmit}
                disabled={loading || !content.trim()}
                className="bg-blue-400 hover:bg-blue-500 text-white font-bold px-4 py-2 rounded-full text-sm disabled:opacity-50 transition-colors"
            >
                Send
            </button>
        </div>
    )
}

export default MessageInput
