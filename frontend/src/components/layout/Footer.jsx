import { Link } from 'react-router-dom'

export default function Footer() {
    return (
        <footer className="border-t border-gray-800 mt-auto bg-black text-gray-400">
            <div className="w-full mx-auto px-4 py-6 flex flex-col sm:flex-row justify-between items-center gap-3">
                <p className="text-sm">© {new Date().getFullYear()} YourApp</p>

                <p className="text-gray-500 text-xs">
                    By signing up, you agree to our{' '}
                    <Link to="/terms" className="text-blue-400 cursor-pointer hover:underline">
                        Terms of Service
                    </Link>{' '}
                    and{' '}
                    <Link to="/privacy" className="text-blue-400 cursor-pointer hover:underline">
                        Privacy Policy
                    </Link>
                </p>
            </div>
        </footer>
    )
}
