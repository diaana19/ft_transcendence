import { Link } from 'react-router-dom'
import { MENTION_RE, MENTION_ONE_RE } from '../../utils/username'

const SPLIT = new RegExp(`(#\\w+|${MENTION_RE.source})`, 'gi')
const IS_TAG = /^#\w+$/
const IS_MENTION = MENTION_ONE_RE

const accent = { color: '#534ab7', fontWeight: 600 }

export default function RichText({ text, className }) {
    if (!text) return null
    return (
        <span className={className}>
            {text.split(SPLIT).map((part, i) => {
                if (IS_TAG.test(part)) {
                    return (
                        <Link
                            key={i}
                            to={`/tag/${part.slice(1)}`}
                            onClick={(e) => e.stopPropagation()}
                            style={accent}
                        >
                            {part}
                        </Link>
                    )
                }
                if (IS_MENTION.test(part)) {
                    return (
                        <Link
                            key={i}
                            to={`/profile/u/${part.slice(1)}`}
                            onClick={(e) => e.stopPropagation()}
                            style={accent}
                        >
                            {part}
                        </Link>
                    )
                }
                return part
            })}
        </span>
    )
}
