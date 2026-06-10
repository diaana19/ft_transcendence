import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { searchTags } from './postService'

const PAGE_SIZE = 20

// TagSearch lets the user search every hashtag, not only the few in the trends.
function TagSearch() {
    const [query, setQuery] = useState('')
    const [tags, setTags] = useState([])
    const [total, setTotal] = useState(0)
    const [page, setPage] = useState(0)
    const [loading, setLoading] = useState(true)

    // Go back to the first page whenever the query changes.
    useEffect(() => {
        setPage(0)
    }, [query])

    // Reload the matching tags whenever the query or page changes (debounced).
    useEffect(() => {
        let active = true
        setLoading(true)
        const id = setTimeout(async () => {
            try {
                const data = await searchTags(query, PAGE_SIZE, page * PAGE_SIZE)
                if (active) {
                    setTags(data.tags)
                    setTotal(data.total)
                }
            } catch (err) {
                console.info('Error searching tags:', err)
            } finally {
                if (active) setLoading(false)
            }
        }, 250)
        return () => {
            active = false
            clearTimeout(id)
        }
    }, [query, page])

    // totalPages is the number of pages (at least one).
    const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE))

    // goToPage switches page and scrolls to the top.
    const goToPage = (p) => {
        if (p < 0 || p >= totalPages || p === page) return
        setPage(p)
        window.scrollTo({ top: 0, behavior: 'smooth' })
    }

    return (
        <div className="w-full mx-auto space-y-4 mt-6">
            <h2 className="text-xl font-bold" style={{ color: '#534ab7' }}>
                Search tags
            </h2>

            <input
                type="text"
                value={query}
                onChange={(e) => setQuery(e.target.value)}
                placeholder="Search hashtags..."
                className="w-full px-4 py-2 text-sm rounded-full focus:outline-none"
                style={{ background: '#f5f3ff', border: '1px solid #ede8fd', color: '#2c2c2a' }}
            />

            {loading ? (
                <p className="text-center text-gray-400 py-8">Loading...</p>
            ) : tags.length === 0 ? (
                <p className="text-center text-gray-400 py-8">No tags found</p>
            ) : (
                <ul className="space-y-2">
                    {tags.map((t) => (
                        <li
                            key={t.tag}
                            className="flex items-center justify-between bg-white/60 border border-gray-200 rounded-xl px-4 py-3"
                        >
                            <Link
                                to={`/tag/${t.tag.replace(/^#/, '')}`}
                                className="text-sm font-semibold truncate hover:underline"
                                style={{ color: '#534ab7' }}
                            >
                                {t.tag}
                            </Link>
                            <span className="text-xs text-gray-400 flex-shrink-0 ml-2">
                                {t.count} {t.count === 1 ? 'post' : 'posts'}
                            </span>
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

export default TagSearch
