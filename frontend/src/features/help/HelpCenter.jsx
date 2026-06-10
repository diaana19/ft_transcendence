import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { UserIcon, PencilIcon, ChatBubbleLeftIcon, SparklesIcon } from '@heroicons/react/24/outline'

// sections holds the help center topics with their questions and answers.
const sections = [
    {
        key: 'account',
        title: 'My account',
        icon: UserIcon,
        colors: {
            border: '#ede8fd', icon: '#ede8fd', iconText: '#534ab7', title: '#534ab7',
            itemBorder: '#ede8fd', itemHover: '#afa9ec', answerBorder: '#f5f0fe',
        },
        faqs: [
            {
                q: 'I lost my password, how do I reset it? ',
                a: 'Go to the login page and click "Forgot password". Enter the email address linked to your account and we\'ll send you a reset link. The link expires after 15 minutes, so use it promptly. If you don\'t see the email, check your spam or junk folder — it sometimes ends up there.',
            },
            {
                q: 'How do I delete my account? ',
                a: 'Head to your Profile page and click the "Settings" button in the top right corner , after you go to danger zone and click in "Delete". You\'ll be asked to confirm your choice — this action is permanent and cannot be undone. All your posts, comments, messages, and personal data will be permanently removed from our servers.',
            },
            {
                q: 'Why should I create an account?',
                a: 'Creating an account unlocks everything Synk has to offer. You can write posts, comment on others\', react with likes and dislikes, follow people whose content you enjoy, send direct messages, add friends, earn badges, and climb the weekly leaderboard. Without an account you can browse public content, but you can\'t interact.',
            },
            {
                q: 'How can I create an account? ',
                a: 'Click "Sign up" on the login page. Fill in your username, email address, password, and date of birth — you must be at least 16 years old to register. You can also sign in with GitHub if you prefer. Once registered, complete your profile by adding a display name, bio, and profile picture.',
            },
            {
                q: 'Can I change my username or email?',
                a: 'Yes. Go to your Profile page, click "Edit profile", and update your display name, email, bio, avatar, or banner image. Note that your username (the @handle) is set at registration and cannot be changed — choose it carefully!',
            },
            {
                q: 'What is two-factor authentication (2FA)? ',
                a: 'Two-factor authentication adds an extra layer of security to your account. Once enabled, you\'ll need to enter a 6-digit code from an authenticator app (like Google Authenticator or Authy) every time you log in. You can enable it in Settings. We strongly recommend turning it on.',
            },
        ],
    },
    {
        key: 'posts',
        title: 'Posts',
        icon: PencilIcon,
        colors: {
            border: '#fde8f0', icon: '#fde8f0', iconText: '#d4537e', title: '#d4537e',
            itemBorder: '#fde8f0', itemHover: '#ed93b1', answerBorder: '#fef0f5',
        },
        faqs: [
            {
                q: 'How do I write a good post? ',
                a: 'Be clear, honest, and specific. Share something you genuinely care about — a question, an opinion, something you learned, or a moment from your day. Posts with a personal angle tend to spark more conversation than generic ones. Keep it under 280 characters and get to the point quickly.',
            },
            {
                q: 'Can I add images or videos to my post? ',
                a: 'Yes! When writing a post, click the photo icon to attach an image or video. Supported formats include JPG, PNG, GIF, MP4, and WebM. The maximum file size is 25MB. Your media will appear inline in the feed.',
            },
            {
                q: 'Can I edit or delete my posts?',
                a: 'Yes to both. On any post you authored, you\'ll see a pencil icon to edit and a trash icon to delete. Editing updates the content immediately. Deletion is permanent — the post and all its comments will be removed.',
            },
            {
                q: 'Can I reply to my own post? ',
                a: 'Absolutely. Commenting on your own post is a great way to add context, share an update, or start the conversation if things are quiet. Your comment appears in the thread just like anyone else\'s.',
            },
            {
                q: 'What should I do if no one engages with my post?',
                a: 'Don\'t be discouraged — it happens to everyone at first. Try engaging with other people\'s posts first: like, comment, and follow users whose content resonates with you. Visibility on Synk grows as your network grows. You can also use hashtags to reach people interested in specific topics.',
            },
            {
                q: 'How do hashtags work? ',
                a: 'Add a hashtag anywhere in your post by typing # followed by a word (e.g. #design, #coding, #travel). Your post will appear in that tag\'s feed and can be discovered by anyone browsing or searching that topic. You can use multiple hashtags in a single post.',
            },
            {
                q: 'What topics can I post about?',
                a: 'Anything you\'re passionate about — tech, art, travel, music, daily life, questions, opinions, projects. Synk is for everyone. The only rule is to keep it respectful and legal. Content that harasses, threatens, or discriminates against others will be removed.',
            },
        ],
    },
    {
        key: 'chat',
        title: 'Chat & friends',
        icon: ChatBubbleLeftIcon,
        colors: {
            border: '#e8f0fd', icon: '#e8f0fd', iconText: '#185fa5', title: '#185fa5',
            itemBorder: '#e8f0fd', itemHover: '#85b7eb', answerBorder: '#f0f4fe',
        },
        faqs: [
            {
                q: 'What is a chat? ',
                a: 'Chat is your private messaging space on Synk. You can send direct messages to any other user — only you and the recipient can see the conversation. Messages are delivered in real time via a live connection, so you\'ll see new messages instantly without refreshing the page.',
            },
            {
                q: 'Who can I message?',
                a: 'You can open a conversation with any registered user on Synk. Just visit their profile and click "Message". You don\'t need to be friends or followers to start a chat — anyone is reachable.',
            },
            {
                q: 'What is the difference between following and friends? ',
                a: 'Following is one-way — you follow someone to see their posts in your feed, and they don\'t have to follow you back. Friendship is mutual — you send a friend request, and if they accept, you both become friends. Friends appear in each other\'s friend list and have a closer connection on the platform.',
            },
            {
                q: 'How do I send a friend request?',
                a: 'Visit someone\'s profile and click "Add friend". They\'ll receive a notification. If they accept, you\'re friends! If you change your mind before they respond, you can cancel the request. You can also see incoming friend requests in your notifications.',
            },
            {
                q: 'How do I connect with new people?',
                a: 'Browse the home feed and engage with posts that interest you — like, comment, and follow the authors. Check out the leaderboard to find the most active users. Use hashtags to discover people who share your interests. The more you interact, the more your network grows.',
            },
        ],
    },
    {
        key: 'gamification',
        title: 'Badges & leaderboard',
        icon: SparklesIcon,
        colors: {
            border: '#cfeacd', icon: '#c5dbc4', iconText: '#3a6d11c3', title: '#3a6d118a',
            itemBorder: '#daefd9', itemHover: '#b9ecb8a0', answerBorder: '#d1dfbe',
        },
        faqs: [
            {
                q: 'What is my score and how is it calculated? ',
                a: 'Your score is the sum of four things: the number of posts you\'ve published, the total likes your posts have received, the number of people following you, and the number of people you follow. Every action counts — there\'s no cap and the score never resets.',
            },
            {
                q: 'What are levels?',
                a: 'Your level is computed from your score using a power-of-two scale: level 1 = 2 points, level 2 = 4, level 3 = 8, and so on up to level 20. Early levels are easy to reach; higher ones require consistent long-term activity. Your level is displayed next to your name on the leaderboard.',
            },
            {
                q: 'What are badges and how do I earn them? ',
                a: 'Badges are permanent achievements that unlock automatically the moment you hit a milestone — no button to press. There are two categories: Social badges (friends and followers) and Post badges (writing and getting likes). Visit the Badges tab in the leaderboard page to see all available badges and which ones you\'ve already earned.',
            },
            {
                q: 'What badges can I earn? ',
                a: 'Social: 👋 Welcome (create your account), 🤝 First bond (add a friend), 🕸️ Network builder (10 friends), 🦋 Social butterfly (50 friends). Posts: ✍️ First words (first post), ✨ First spark (first like received), 📝 Content creator (25 posts), ⭐ Rising star (20 likes received).',
            },
            {
                q: 'How does the weekly leaderboard work?',
                a: 'The leaderboard ranks all Synk users by their activity score. You can filter it by total score, posts, likes, or followers to see who leads in each category. It resets every Monday at midnight UTC so everyone gets a fresh start each week. Your badges and overall level are permanent and never reset.',
            },
            {
                q: 'Can I see another user\'s badges?',
                a: 'Yes! Visit any user\'s profile and click the "Badges" tab to see which achievements they\'ve earned. It\'s a great way to discover active and engaged members of the community.',
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

// HelpCenter shows the help topics and their questions in expandable sections.
export default function HelpCenter() {
    return (
        <div className="min-h-screen bg-transparent px-6 py-8 mx-auto">
            <div className="mb-8">
                <h1 className="text-2xl font-bold mb-1" style={{ color: '#2c2c2a' }}>
                    Help center
                </h1>
				<p className="text-sm leading-relaxed mb-3" style={{ color: '#5f5e5a' }}>
					Welcome to the Synk Help Center. Here you'll find answers to the most common questions about using the platform — from setting up your account and writing your first post, to understanding how the friends system works and earning your first badge.
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
