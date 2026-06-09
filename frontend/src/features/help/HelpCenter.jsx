import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { UserIcon, PencilIcon, ChatBubbleLeftIcon } from '@heroicons/react/24/outline'

const sections = [
    {
        key: 'account',
        title: 'My account',
        icon: UserIcon,
        colors: {
            border: '#ede8fd',
            icon: '#ede8fd',
            iconText: '#534ab7',
            title: '#534ab7',
            itemBorder: '#ede8fd',
            itemHover: '#afa9ec',
            answerBorder: '#f5f0fe',
        },
        faqs: [
            {
                q: 'I lost my password, how do I reset it?',
                a: 'Go to the login page and click "Forgot password". Enter your email address and we\'ll send you a reset link. Check your spam folder if you don\'t see it within a few minutes.',
            },
            {
                q: 'How do I delete my account?',
                a: 'Head to your Profile page and click the "Delete" button in the top right. You\'ll be asked to confirm — this action is permanent and cannot be undone. All your posts and data will be removed.',
            },
            {
                q: 'Why should I create an account?',
                a: "With an account you can post, comment, follow people, join communities, and send messages. It's the best way to connect with people who share your interests on Synk.",
            },
            {
                q: 'How can I create an account?',
                a: 'Click "Sign up" on the login page. Fill in your name, email, and a password. You can also complete your profile by adding a bio and a profile picture afterwards.',
            },
        ],
    },
    {
        key: 'posts',
        title: 'Posts',
        icon: PencilIcon,
        colors: {
            border: '#fde8f0',
            icon: '#fde8f0',
            iconText: '#d4537e',
            title: '#d4537e',
            itemBorder: '#fde8f0',
            itemHover: '#ed93b1',
            answerBorder: '#fef0f5',
        },
        faqs: [
            {
                q: 'How do I write a good post?',
                a: 'Be clear and genuine. Share something you actually care about — a thought, a question, or something you learned. Posts with a personal angle tend to get more engagement than generic ones.',
            },
            {
                q: 'Can I answer my own post?',
                a: "Yes, absolutely. Adding your own comment is a great way to give more context, update people with new information, or kick off the conversation if it's been quiet.",
            },
            {
                q: 'What should I do if no one comments on my post?',
                a: "Don't worry — it happens to everyone. Try engaging with other people's posts first, follow more users, or post in a community related to your topic. Visibility grows as your network grows.",
            },
            {
                q: 'What topics can I post about?',
                a: "Anything you're passionate about — tech, art, travel, daily life, questions, opinions. Synk is for everyone. Just keep it respectful and follow the community guidelines.",
            },
        ],
    },
    {
        key: 'chat',
        title: 'Chat & friends',
        icon: ChatBubbleLeftIcon,
        colors: {
            border: '#e8f0fd',
            icon: '#e8f0fd',
            iconText: '#185fa5',
            title: '#185fa5',
            itemBorder: '#e8f0fd',
            itemHover: '#85b7eb',
            answerBorder: '#f0f4fe',
        },
        faqs: [
            {
                q: 'What is a chat?',
                a: 'Chat is your private messaging space on Synk. You can send direct messages to other users — only you and the other person can see the conversation.',
            },
            {
                q: 'Who should I chat with?',
                a: "Anyone you've connected with — people you follow, friends you've added, or someone whose post caught your attention. Chat is a great way to continue a conversation more privately.",
            },
            {
                q: 'How do I connect with people?',
                a: 'Visit someone\'s profile and click "Follow" to see their posts, or "Add friend" to send them a friend request. Once they accept, you\'ll be connected and can message each other.',
            },
            {
                q: 'How does the friends system work?',
                a: "Send a friend request from someone's profile. If they accept, you become friends — different from following, which is one-way. Friends can see each other's activity more easily and have a closer connection on the platform.",
            },
        ],
    },
]

function FaqItem({ faq, colors }) {
    const [open, setOpen] = useState(false)

    return (
        <div
            onClick={() => setOpen(!open)}
            className="rounded-xl mb-2 overflow-hidden cursor-pointer transition-colors"
            style={{ border: `0.5px solid ${open ? colors.itemHover : colors.itemBorder}` }}
        >
            <div className="flex items-center justify-between px-4 py-3 bg-white">
                <span className="text-sm font-semibold" style={{ color: '#2c2c2a' }}>
                    {faq.q}
                </span>
                <svg
                    className="w-4 h-4 flex-shrink-0 ml-3 transition-transform"
                    style={{
                        color: '#b4b2a9',
                        transform: open ? 'rotate(180deg)' : 'rotate(0deg)',
                    }}
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                    strokeWidth={2}
                >
                    <path strokeLinecap="round" strokeLinejoin="round" d="M19 9l-7 7-7-7" />
                </svg>
            </div>
            {open && (
                <div
                    className="px-4 pb-3 bg-white text-sm leading-relaxed"
                    style={{ color: '#5f5e5a', borderTop: `0.5px solid ${colors.answerBorder}` }}
                >
                    {faq.a}
                </div>
            )}
        </div>
    )
}

export default function HelpCenter() {
    return (
        <div className="min-h-screen bg-transparent px-6 py-8 mx-auto">
            <div className="mb-8">
                <h1 className="text-2xl font-bold mb-1" style={{ color: '#2c2c2a' }}>
                    Help center
                </h1>
                <p className="text-sm" style={{ color: '#b4b2a9' }}>
                    Find answers to the most common questions about Synk.
                </p>
            </div>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                {sections.map((section) => {
                    const Icon = section.icon
                    const c = section.colors
                    return (
                        <div key={section.key} className="mb-8">
                            <div
                                className="flex items-center gap-3 mb-3 pb-2"
                                style={{ borderBottom: `1.5px solid ${c.border}` }}
                            >
                                <div
                                    className="w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0"
                                    style={{ background: c.icon }}
                                >
                                    <Icon className="w-4 h-4" style={{ color: c.iconText }} />
                                </div>
                                <span className="text-sm font-semibold" style={{ color: c.title }}>
                                    {section.title}
                                </span>
                            </div>

                            {section.faqs.map((faq, i) => (
                                <FaqItem key={i} faq={faq} colors={c} />
                            ))}
                        </div>
                    )
                })}{' '}
            </div>
        </div>
    )
}
