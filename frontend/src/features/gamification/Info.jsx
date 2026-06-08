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
                        text: 'Score = posts + likes received + followers + following. Every action counts toward your level.',
                    },
                    {
                        icon: 'ti-chart-line',
                        title: 'Levels',
                        text: 'Levels follow a power-of-two scale — each level requires double the score of the previous one. Level 0 starts at 0 points.',
                    },
                    {
                        icon: 'ti-award',
                        title: 'Badges',
                        text: 'Badges unlock automatically when you hit a milestone. No need to claim them — just keep being active.',
                    },
                    {
                        icon: 'ti-calendar-week',
                        title: 'Weekly reset',
                        text: 'The leaderboard resets every Monday. Your badges and level are permanent and never reset.',
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
