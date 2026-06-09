import { useState, useEffect } from 'react'

// CookiesBanner shows a cookie consent banner until the user answers it.
export default function CookiesBanner() {
    const [visible, setVisible] = useState(false)

    // Show the banner only if there is no saved consent yet.
    useEffect(() => {
        const consent = localStorage.getItem('cookie_consent')
        if (!consent) setVisible(true)
    }, [])

    const handleAccept = () => {
        localStorage.setItem('cookie_consent', 'accepted')
        setVisible(false)
    }

    const handleDecline = () => {
        localStorage.setItem('cookie_consent', 'declined')
        setVisible(false)
    }

    if (!visible) return null

    return (
        <div className="fixed bottom-16 left-0 right-0 z-50 p-4 md:bottom-4 md:left-4 md:right-auto md:max-w-sm">
            <div className="bg-white border border-gray-200 rounded-2xl shadow-xl p-4">
                <h3 className="text-sm font-bold text-gray-900 mb-2">🍪 We use cookies</h3>
                <p className="text-xs text-gray-500 mb-4">
                    We use cookies to improve your experience, analyze traffic and personalize
                    content. By continuing to use Synk, you agree to our cookie policy.
                </p>
                <div className="flex gap-2">
                    <button
                        onClick={handleDecline}
                        className="flex-1 border border-gray-300 px-3 py-1.5 rounded-full text-xs font-semibold text-gray-600 hover:bg-gray-100 transition-colors"
                    >
                        Decline
                    </button>
                    <button
                        onClick={handleAccept}
                        className="flex-1 bg-blue-400 hover:bg-blue-500 px-3 py-1.5 rounded-full text-xs font-bold text-white transition-colors"
                    >
                        Accept all
                    </button>
                </div>
            </div>
        </div>
    )
}
