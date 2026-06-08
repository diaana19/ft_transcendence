import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { getFollowers, getFollowing } from './userService'

export default function FollowersModal({ userId, type, onClose }) {
  const navigate = useNavigate()
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetch = type === 'followers' ? getFollowers : getFollowing
    fetch(userId)
      .then(setUsers)
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [userId, type])

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center"
      style={{ background: 'rgba(237,232,253,0.4)', backdropFilter: 'blur(4px)' }}>
      <div className="absolute inset-0" onClick={onClose} />
      <div className="relative z-10 bg-white rounded-2xl w-full max-w-sm overflow-hidden"
        style={{ border: '2px solid #ede8fd', boxShadow: '0 8px 32px rgba(167,139,250,0.15)', maxHeight: '70vh' }}>

        <div className="flex items-center gap-3 px-4 py-3 border-b" style={{ borderColor: '#ede8fd' }}>
          <button onClick={onClose} className="rounded-full p-1 hover:bg-[#faf8f5]" style={{ color: '#534ab7' }}>✕</button>
          <h3 className="font-bold text-base capitalize" style={{ color: '#2c2c2a' }}>{type}</h3>
        </div>

        <div className="overflow-y-auto" style={{ maxHeight: 'calc(70vh - 56px)' }}>
          {loading ? (
            <p className="text-center py-8 text-sm" style={{ color: '#b4b2a9' }}>Loading...</p>
          ) : users.length === 0 ? (
            <p className="text-center py-8 text-sm" style={{ color: '#b4b2a9' }}>No {type} yet</p>
          ) : users.map(u => (
            <div key={u.id}
              className="flex items-center gap-3 px-4 py-3 cursor-pointer hover:bg-[#faf8f5] transition-colors"
              onClick={() => { navigate(`/profile/${u.id}`); onClose() }}>
              <div className="w-10 h-10 rounded-full overflow-hidden flex-shrink-0 flex items-center justify-center font-bold text-sm"
                style={{ background: '#ede8fd', color: '#534ab7' }}>
                {u.avatar
                  ? <img src={u.avatar} alt={u.username} className="w-full h-full object-cover" />
                  : u.username?.[0]?.toUpperCase()}
              </div>
              <div>
                <p className="font-semibold text-sm" style={{ color: '#2c2c2a' }}>{u.displayname}</p>
                <p className="text-xs" style={{ color: '#afa9ec' }}>@{u.username}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
