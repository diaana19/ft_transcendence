export default function RightSidebar() {
  return (
    <aside className="w-80 h-screen fixed right-0 top-0
                      backdrop-blur-xl
                      border-l border-gray-200
                      px-4 py-6 hidden lg:block">

      {/* Search placeholder */}
      <div className="mb-6">
        <input
          type="text"
          placeholder="Search..."
          className="w-full px-4 py-2 rounded-xl
                     bg-white/70 border border-gray-200
                     text-gray-700 placeholder-gray-400
                     focus:outline-none focus:ring-2 focus:ring-gray-300"
        />
      </div>

      {/* Placeholder card 1 */}
      <div className="bg-white/60 border border-gray-200 rounded-xl p-4 mb-4">
        <h3 className="text-sm font-semibold text-gray-700 mb-2">
          Trends
        </h3>
        <p className="text-xs text-gray-400">
          Nothing here yet...
        </p>
      </div>

      {/* Placeholder card 2 */}
      <div className="bg-white/60 border border-gray-200 rounded-xl p-4">
        <h3 className="text-sm font-semibold text-gray-700 mb-2">
          Suggestions
        </h3>
        <p className="text-xs text-gray-400">
          Future feature: users / posts recommendations
        </p>
      </div>

    </aside>
  )
}