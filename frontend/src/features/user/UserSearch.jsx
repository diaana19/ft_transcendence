import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { searchUsers } from './userService'

const PAGE_SIZE = 20

// UserSearch finds users by username or display name, with pagination.
function UserSearch() {
    const navigate = useNavigate()
    const [query, setQuery] = useState('')
    const [users, setUsers] = useState([])
    const [total, setTotal] = useState(0)
    const [page, setPage] = useState(0)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        setPage(0)
    }, [query])

    useEffect(() => {
        let active = true
        setLoading(true)
        const id = setTimeout(async () => {
            try {
                const data = await searchUsers(query, PAGE_SIZE, page * PAGE_SIZE)
                if (active) {
                    setUsers(data.users)
                    setTotal(data.total)
                }
            } catch (err) {
                console.info('Error searching users:', err)
            } finally {
                if (active) setLoading(false)
            }
        }, 250)
        return () => {
            active = false
            clearTimeout(id)
        }
    }, [query, page])

    const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE))

    const goToPage = (p) => {
        if (p < 0 || p >= totalPages || p === page) return
        setPage(p)
        window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    return (
        <div className="w-full mx-auto space-y-4 mt-6">
            <h2 className="text-xl font-bold" style={{ color: '#534ab7' }}>
                Search users
            </h2>

            <input
                type="text"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Search by name or username..."
                className="w-full px-4 py-2 text-sm rounded-full focus:outline-none"
                style={{ background: '#f5f3ff', border: '1px solid #ede8fd', color: '#2c2c2a' }}
            />

            {loading ? (
                <p className="text-center text-gray-400 py-8">Loading...</p>
            ) : users.length === 0 ? (
                <p className="text-center text-gray-400 py-8">No users found</p>
            ) : (
                <ul className="space-y-2">
                    {users.map((u) => (
                        <li
                            key={u.id}
                            onClick={() => navigate(`/profile/${u.id}`)}
                            className="flex items-center gap-3 bg-white/60 border border-gray-200 rounded-xl px-4 py-3 cursor-pointer"
                        >
                            <div
                                className="w-10 h-10 rounded-full overflow-hidden flex-shrink-0 flex items-center justify-center font-bold text-sm"
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
                            <div className="flex-1 min-w-0">
                                <p
                                    className="text-sm font-bold truncate"
                                    style={{ color: '#2c2c2a' }}
                                >
                                    {u.displayname || u.username}
                                </p>
                                <p className="text-xs truncate" style={{ color: '#afa9ec' }}>
                                    @{u.username}
                                </p>
                            </div>
                        </li>
                    ))}
                </ul>
            )}

            {!loading && total > 0 && (
                <div className="flex items-center justify-center gap-4 py-6">
                    <button
                        onClick={() => goToPage(page - 1)}
                        disabled={page === 0}
                        className="px-4 py-2 rounded-full text-sm font-bold disabled:opacity-40"
                        style={{ background: '#ede8fd', color: '#534ab7' }}
                    >
                        Previous
                    </button>
                    <span className="text-sm font-medium" style={{ color: '#534ab7' }}>
                        Page {page + 1} of {totalPages}
                    </span>
                    <button
                        onClick={() => goToPage(page + 1)}
                        disabled={page + 1 >= totalPages}
                        className="px-4 py-2 rounded-full text-sm font-bold disabled:opacity-40"
                        style={{ background: '#ede8fd', color: '#534ab7' }}
                    >
                        Next
                    </button>
                </div>
            )}
        </div>
    )
}

export default UserSearch
