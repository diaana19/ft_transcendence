// Info explains how the score, levels and badges work.
export default function Info() {
    return (
        <div>
            <div
                className="flex items-center gap-3 mb-4 pb-2"
                style={{ borderBottom: '1.5px solid #e8f0fd' }}
            >
                <div
                    className="w-8 h-8 rounded-lg flex items-center justify-center"
                    style={{ background: '#e8f0fd' }}
                >
                    <i
                        className="ti ti-info-circle"
                        aria-hidden="true"
                        style={{ color: '#185fa5', fontSize: '15px' }}
                    />
                </div>
                <span className="text-sm font-semibold" style={{ color: '#185fa5' }}>
                    How it works
                </span>
            </div>
            <div className="flex flex-col gap-3">
                {[
                    {
                        icon: 'ti-calculator',
                        title: 'Your score',
                        text: 'Your score is the sum of everything you do on Synk — every post you publish, every like your posts receive, every person who follows you, and every person you follow. The more active you are, the higher your score climbs. There is no cap.',
                    },
                    {
                        icon: 'ti-chart-line',
                        title: 'Levels',
                        text: 'Levels are based on a power-of-two scale: level 1 requires 2 points, level 2 requires 4, level 3 requires 8, and so on. This means early levels are easy to reach but higher levels require sustained activity. Your level is always displayed on the leaderboard next to your score.',
                    },
                    {
                        icon: 'ti-award',
                        title: 'Badges',
                        text: 'Badges are permanent achievements that unlock automatically when you hit a milestone — no button to press, no claim needed. There are two categories: 🤝 Social (friends & followers) and  ✍️ Posts (writing & getting likes). Locked badges show you exactly what you need to do next.',
                    },
					{
						icon: 'ti-users',
						title: 'Social badges',
						text: '👋 Welcome — create your account. 🤝 First bond — add your first friend. 🕸️ Network builder — reach 10 friends. 🦋 Social butterfly — connect with 50+ friends. Follow people and send friend requests to unlock these.',
					},
					{
						icon: 'ti-pencil',
						title: 'Post badges',
						text: '✍️ First words — publish your first post. ✨ First spark — receive your first like. 📝 Content creator — publish 25 posts. ⭐ Rising star — collect 20 likes across all your posts. Quality and consistency both count.',
					},
                    {
                        icon: 'ti-calendar-week',
                        title: 'Weekly reset',
                        text: 'The leaderboard resets every Monday at midnight UTC so everyone gets a fresh start each week. Your badges and overall level are permanent — they never reset and accumulate over your entire time on Synk.',
                    },
                ].map((item) => (
                    <div
                        key={item.title}
                        className="flex gap-3 items-start rounded-xl px-4 py-3"
                        style={{ background: 'white', border: '0.5px solid #e8f0fd' }}
                    >
                        <div
                            className="w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0"
                            style={{ background: '#e8f0fd' }}
                        >
                            <i
                                className={`ti ${item.icon}`}
                                aria-hidden="true"
                                style={{ color: '#185fa5', fontSize: '15px' }}
                            />
                        </div>
                        <div>
                            <div
                                className="text-xs font-semibold mb-1"
                                style={{ color: '#0c447c' }}
                            >
                                {item.title}
                            </div>
                            <div className="text-xs leading-relaxed" style={{ color: '#888780' }}>
                                {item.text}
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    )
}
