
import { useEffect, useState } from "react"
import { useAuth } from "../../hooks/useAuth"
import api from "../../services/axiosInstance"
const BADGES = [
  {
    category: "social",
    items: [
      { key: "welcome", icon: "ti-user-plus", name: "Welcome", desc: "Created your account on Synk", condition: "Day one", check: (s) => true },
      { key: "first_bond", icon: "ti-heart-handshake", name: "First bond",       desc: "Added your very first friend",         condition: "1 friend",     check: (s) => s.followers >= 1 },
      { key: "network_builder", icon: "ti-topology-star",   name: "Network builder",  desc: "Reached 10 friends on the platform",   condition: "10 friends",   check: (s) => s.followers >= 10 },
      { key: "social_butterfly", icon: "ti-butterfly",       name: "Social butterfly", desc: "Connected with more than 50 friends",  condition: "50 friends",   check: (s) => s.followers >= 50 },
    ]
  },
  {
    category: "posts",
    items: [
      { key: "first_words",     icon: "ti-feather",  name: "First words",     desc: "Published your very first post",       condition: "1 post",       check: (s) => s.posts >= 1 },
      { key: "first_spark",     icon: "ti-heart",    name: "First spark",     desc: "Received your first like on a post",   condition: "1 like",       check: (s) => s.likes >= 1 },
      { key: "content_creator", icon: "ti-writing",  name: "Content creator", desc: "Published 25 posts in total",          condition: "25 posts",     check: (s) => s.posts >= 25 },
      { key: "rising_star",     icon: "ti-star",     name: "Rising star",     desc: "Collected 20 likes across your posts", condition: "20 likes",     check: (s) => s.likes >= 20 },
    ]
  },
  {
    category: "chat",
    items: [
      { key: "chatter",   icon: "ti-message-circle", name: "Chatter",   desc: "Sent your very first message", condition: "1 message",    check: (s) => s.messages >= 1 },
      { key: "messenger", icon: "ti-messages",       name: "Messenger", desc: "Sent 100 messages total",      condition: "100 messages", check: (s) => s.messages >= 100 },
    ]
  }
]
const CATEGORY_STYLE = {
  social: { border: "#ede8fd", icon: "#ede8fd", iconText: "#534ab7", title: "#534ab7", cardBorder: "#ede8fd", pillBg: "#eeedfe", pillText: "#534ab7", emojiBg: "#ede8fd" },
  posts:  { border: "#c5dbc4", icon: "#c5dbc4", iconText: "#3b6d11", title: "#3b6d11", cardBorder: "#dce6cf", pillBg: "#eaf3de", pillText: "#3b6d11", emojiBg: "#c5dbc4" },
  chat:   { border: "#d0edef", icon: "#d0edef", iconText: "#0f6e56", title: "#0f6e56", cardBorder: "#d0edef", pillBg: "#e1f5ee", pillText: "#0f6e56", emojiBg: "#d0edef" },
}
const CATEGORY_ICON  = { social: "ti-users",   posts: "ti-pencil",   chat: "ti-message-2" }
const CATEGORY_LABEL = { social: "Social",      posts: "Posts",       chat: "Chat" }
const WEEKLY_TABS = [
  { key: "total",     label: "Score" },
  { key: "posts",     label: "Posts" },
  { key: "likes",     label: "Likes" },
  { key: "followers", label: "Followers" },
]
const NAV_TABS = [
  { key: "leaderboard", label: "Leaderboard", icon: "ti-trophy",      activeColor: "#534ab7", activeBg: "#eeedfe", activeBorder: "#ede8fd" },
  { key: "badges",      label: "Badges",      icon: "ti-award",       activeColor: "#3b6d11", activeBg: "#eaf3de", activeBorder: "#dce6cf" },
  { key: "info",        label: "How it works", icon: "ti-info-circle", activeColor: "#185fa5", activeBg: "#e6f1fb", activeBorder: "#e8f0fd" },
]
function LevelBar({ level }) {
  const pct = Math.min((level / 20) * 100, 100)
  return (
    <div className="w-full rounded-full overflow-hidden" style={{ height: '4px', background: '#ede8fd' }}>
      <div className="rounded-full h-full transition-all" style={{ width: `${pct}%`, background: '#534ab7' }} />
    </div>
  )
}
export default function Gamification() {
  const { user: authUser } = useAuth()
  const [activeSection, setActiveSection] = useState("leaderboard")
  const [stats, setStats]         = useState(null)
  const [leaderboard, setLeaderboard] = useState([])
  const [activeTab, setActiveTab] = useState("total")
  const [loadingStats, setLoadingStats]   = useState(true)
  const [loadingBoard, setLoadingBoard]   = useState(true)
  useEffect(() => { fetchStats(); fetchLeaderboard() }, [])
  const fetchStats = async () => {
    try { const { data } = await api.get("/api/gamification"); setStats(data) }
    catch (err) { console.error(err) }
    finally { setLoadingStats(false) }
  }
  const fetchLeaderboard = async () => {
    try { const { data } = await api.get("/api/leaderboard"); setLeaderboard(data) }
    catch (err) { console.error(err) }
    finally { setLoadingBoard(false) }
  }
  const sortedBoard = [...leaderboard].sort((a, b) => {
    const va = activeTab === "total" ? b.stats?.total          : b.stats?.[activeTab]?.count
    const vb = activeTab === "total" ? a.stats?.total          : a.stats?.[activeTab]?.count
    return (va ?? 0) - (vb ?? 0)
  }).slice(0, 10)
  const userStats = { posts: stats?.posts?.count ?? 0, likes: stats?.likes?.count ?? 0, followers: stats?.followers?.count ?? 0, messages: 0 }
  const earnedCount = BADGES.flatMap(s => s.items).filter(b => b.check(userStats)).length
  const totalBadges = BADGES.flatMap(s => s.items).length
  return (
    <div className="min-h-screen bg-transparent px-6 py-8 w_full ">
      {/* ── NAV TABS ── */}
      <div className="flex gap-2 mb-8 p-1 rounded-2xl" style={{ background: 'white', border: '0.5px solid #ede8fd' }}>
        {NAV_TABS.map(tab => {
          const active = activeSection === tab.key
          return (
            <button
              key={tab.key}
              onClick={() => setActiveSection(tab.key)}
              className="flex-1 flex items-center justify-center gap-2 py-2.5 rounded-xl text-sm font-semibold transition-all"
              style={{
                background:  active ? tab.activeBg    : 'transparent',
                color:       active ? tab.activeColor : '#b4b2a9',
                border:      active ? `0.5px solid ${tab.activeBorder}` : '0.5px solid transparent',
                cursor: 'pointer'
              }}
            >
              <i className={`ti ${tab.icon}`} aria-hidden="true" style={{ fontSize: '15px' }} />
              <span className="hidden sm:inline">{tab.label}</span>
            </button>
          )
        })}
      </div>
      {/* ── LEADERBOARD ── */}
      {activeSection === "leaderboard" && (
        <div>
          <div className="flex items-center gap-3 mb-3 pb-2" style={{ borderBottom: '1.5px solid #ede8fd' }}>
            <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ background: '#ede8fd' }}>
              <i className="ti ti-trophy" aria-hidden="true" style={{ color: '#534ab7', fontSize: '15px' }} />
            </div>
            <span className="text-sm font-semibold" style={{ color: '#534ab7' }}>Weekly leaderboard</span>
          </div>
          <div className="flex gap-1 mb-3 p-1 rounded-xl" style={{ background: '#f5f0fe' }}>
            {WEEKLY_TABS.map(tab => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key)}
                className="flex-1 text-xs font-semibold py-1.5 rounded-lg transition-colors"
                style={{
                  background: activeTab === tab.key ? 'white' : 'transparent',
                  color:      activeTab === tab.key ? '#534ab7' : '#afa9ec',
                  border: 'none', cursor: 'pointer'
                }}
              >
                {tab.label}
              </button>
            ))}
          </div>
          <div className="rounded-2xl overflow-hidden" style={{ border: '0.5px solid #ede8fd' }}>
            {loadingBoard ? (
              <p className="text-center py-6 text-sm" style={{ color: '#b4b2a9' }}>Loading...</p>
            ) : sortedBoard.length === 0 ? (
              <p className="text-center py-6 text-sm" style={{ color: '#b4b2a9' }}>No data yet</p>
            ) : sortedBoard.map((entry, i) => {
              const isMe = entry.id === authUser?.userId
              const value = activeTab === "total" ? entry.stats?.total : entry.stats?.[activeTab]?.count
              const medals = ["🥇", "🥈", "🥉"]
              return (
                <div
                  key={entry.id}
                  className="flex items-center gap-3 px-4 py-3"
                  style={{
                    borderBottom: i < sortedBoard.length - 1 ? '0.5px solid #f5f0fe' : 'none',
                    background: isMe ? '#faf8ff' : 'white'
                  }}
                >
                  <span className="text-sm w-6 text-center" style={{ color: '#b4b2a9' }}>
                    {i < 3 ? medals[i] : i + 1}
                  </span>
                  <div className="w-8 h-8 rounded-full flex items-center justify-center font-semibold flex-shrink-0 overflow-hidden" style={{ background: '#ede8fd', color: '#534ab7', fontSize: '11px' }}>
                    {entry.avatar
                      ? <img src={entry.avatar} alt={entry.username} className="w-full h-full object-cover" />
                      : entry.username?.charAt(0).toUpperCase()
                    }
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-sm font-semibold truncate" style={{ color: '#2c2c2a' }}>{entry.username}</span>
                      {isMe && <span className="text-xs px-2 py-0.5 rounded-full flex-shrink-0" style={{ background: '#eeedfe', color: '#534ab7' }}>you</span>}
                    </div>
                    <LevelBar level={entry.stats?.level ?? 0} />
                  </div>
                  <div className="text-right flex-shrink-0">
                    <div className="text-sm font-semibold" style={{ color: '#534ab7' }}>{value ?? 0}</div>
                    <div className="text-xs" style={{ color: '#b4b2a9' }}>lv. {entry.stats?.level ?? 0}</div>
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      )}
      {/* ── BADGES ── */}
      {activeSection === "badges" && (
        <div>
          <div className="flex items-center gap-3 mb-4 pb-2" style={{ borderBottom: '1.5px solid #fde8f0' }}>
            <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ background: '#fde8f0' }}>
              <i className="ti ti-award" aria-hidden="true" style={{ color: '#d4537e', fontSize: '15px' }} />
            </div>
            <span className="text-sm font-semibold" style={{ color: '#d4537e' }}>Badges</span>
          </div>
          {BADGES.map(section => {
            const c = CATEGORY_STYLE[section.category]
            return (
              <div key={section.category} className="mb-6">
                <div className="flex items-center gap-2 mb-3">
                  <div className="w-5 h-5 rounded flex items-center justify-center" style={{ background: c.icon }}>
                    <i className={`ti ${CATEGORY_ICON[section.category]}`} aria-hidden="true" style={{ color: c.iconText, fontSize: '11px' }} />
                  </div>
                  <span className="text-xs font-semibold" style={{ color: c.title }}>{CATEGORY_LABEL[section.category]}</span>
                </div>
                <div className="grid gap-3" style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))' }}>
                  {section.items.map(badge => {
                    const earned = badge.check(userStats)
                    return (
                      <div
                        key={badge.key}
                        className="rounded-2xl p-4 flex flex-col items-center text-center gap-2 transition-all cursor-pointer hover:scale-105 "
                        style={{ background: 'white', border: `0.5px solid ${c.cardBorder}` }}
                      >
                        <div className="w-12 h-12 rounded-full flex items-center justify-center" style={{ background: c.emojiBg }}>
                          {badge.img ? (
                            <img
                            src={badge.img}
                            alt={badge.name}
                            className="w-full h-full object-cover"
                            onError={(e) => {
                                e.target.style.display = 'none'
                                e.target.nextSibling.style.display = 'flex'
                            }}
                            />
                        ) : null}
                        <i
                            className={`ti ${badge.icon}`}
                            aria-hidden="true"
                            style={{ color: c.iconText, fontSize: '22px', display: badge.img ? 'none' : 'flex' }}
                        />
                        </div>
                        <div className="text-xs font-semibold" style={{ color: earned ? c.title : '#888780' }}>{badge.name}</div>
                        <div className="text-xs leading-relaxed" style={{ color: '#888780' }}>{badge.desc}</div>
                        {earned ? (
                          <span className="text-xs px-2 py-0.5 rounded-full font-semibold" style={{ background: c.pillBg, color: c.pillText }}>{badge.condition}</span>
                        ) : (
                          <span className="text-xs px-2 py-0.5 rounded-full font-semibold" style={{ background: '#f1efe8', color: '#888780' }}>
                            <i className="ti ti-lock" aria-hidden="true" style={{ fontSize: '10px', marginRight: '3px' }} />{badge.condition}
                          </span>
                        )}
                      </div>
                    )
                  })}
                </div>
              </div>
            )
          })}
        </div>
      )}
      {/* ── INFO ── */}
      {activeSection === "info" && (
        <div>
          <div className="flex items-center gap-3 mb-4 pb-2" style={{ borderBottom: '1.5px solid #e8f0fd' }}>
            <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ background: '#e8f0fd' }}>
              <i className="ti ti-info-circle" aria-hidden="true" style={{ color: '#185fa5', fontSize: '15px' }} />
            </div>
            <span className="text-sm font-semibold" style={{ color: '#185fa5' }}>How it works</span>
          </div>
          <div className="flex flex-col gap-3">
            {[
              { icon: "ti-calculator",   title: "Your score",    text: "Score = posts + likes received + followers + following. Every action counts toward your level." },
              { icon: "ti-chart-line",   title: "Levels",        text: "Levels follow a power-of-two scale — each level requires double the score of the previous one. Level 0 starts at 0 points." },
              { icon: "ti-award",        title: "Badges",        text: "Badges unlock automatically when you hit a milestone. No need to claim them — just keep being active." },
              { icon: "ti-calendar-week",title: "Weekly reset",  text: "The leaderboard resets every Monday. Your badges and level are permanent and never reset." },
            ].map(item => (
              <div key={item.title} className="flex gap-3 items-start rounded-xl px-4 py-3" style={{ background: 'white', border: '0.5px solid #e8f0fd' }}>
                <div className="w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0" style={{ background: '#e8f0fd' }}>
                  <i className={`ti ${item.icon}`} aria-hidden="true" style={{ color: '#185fa5', fontSize: '15px' }} />
                </div>
                <div>
                  <div className="text-xs font-semibold mb-1" style={{ color: '#0c447c' }}>{item.title}</div>
                  <div className="text-xs leading-relaxed" style={{ color: '#888780' }}>{item.text}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
