import { useState } from 'react'
import Leaderboard from './LeaderBoard'
import Badges from './Badges'
import Info from './Info'
import { useAuth } from '../../hooks/useAuth'

export default function Gamification() {
    const { user } = useAuth()
    const isLoggedIn = !!user
    const [activeSection, setActiveSection] = useState(isLoggedIn ? 'leaderboard' : 'badges')

    const NAV_TABS = isLoggedIn
        ? [
              {
                  key: 'leaderboard',
                  label: 'Leaderboard',
                  icon: 'ti-trophy',
                  activeColor: '#534ab7',
                  activeBg: '#eeedfe',
                  activeBorder: '#ede8fd',
              },
              {
                  key: 'badges',
                  label: 'Badges',
                  icon: 'ti-award',
                  activeColor: '#3b6d11',
                  activeBg: '#eaf3de',
                  activeBorder: '#dce6cf',
              },
              {
                  key: 'info',
                  label: 'How it works',
                  icon: 'ti-info-circle',
                  activeColor: '#185fa5',
                  activeBg: '#e6f1fb',
                  activeBorder: '#e8f0fd',
              },
          ]
        : [
              {
                  key: 'badges',
                  label: 'Badges',
                  icon: 'ti-award',
                  activeColor: '#3b6d11',
                  activeBg: '#eaf3de',
                  activeBorder: '#dce6cf',
              },
              {
                  key: 'info',
                  label: 'How it works',
                  icon: 'ti-info-circle',
                  activeColor: '#185fa5',
                  activeBg: '#e6f1fb',
                  activeBorder: '#e8f0fd',
              },
          ]

    return (
        <div className="min-h-screen bg-transparent px-6 py-8 w-full">
            <div
                className="flex gap-2 mb-8 p-1 rounded-2xl"
                style={{ background: 'white', border: '0.5px solid #ede8fd' }}
            >
                {NAV_TABS.map((tab) => {
                    const active = activeSection === tab.key
                    return (
                        <button
                            key={tab.key}
                            onClick={() => setActiveSection(tab.key)}
                            className="flex-1 flex items-center justify-center gap-2 py-2.5 rounded-xl text-sm font-semibold transition-all"
                            style={{
                                background: active ? tab.activeBg : 'transparent',
                                color: active ? tab.activeColor : '#b4b2a9',
                                border: active
                                    ? `0.5px solid ${tab.activeBorder}`
                                    : '0.5px solid transparent',
                                cursor: 'pointer',
                            }}
                        >
                            <i
                                className={`ti ${tab.icon}`}
                                aria-hidden="true"
                                style={{ fontSize: '15px' }}
                            />
                            <span className="hidden sm:inline">{tab.label}</span>
                        </button>
                    )
                })}
            </div>

            {activeSection === 'leaderboard' && <Leaderboard />}
            {activeSection === 'badges' && <Badges />}
            {activeSection === 'info' && <Info />}
        </div>
    )
}
