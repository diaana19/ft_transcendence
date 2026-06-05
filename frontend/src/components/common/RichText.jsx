/*
** File: RichText.jsx
** Description: Renders user text with #hashtags and @mentions highlighted.
** Splits the string on a token regex and renders matches specially; plain text
** stays as escaped React string children (no dangerouslySetInnerHTML), so user
** content cannot inject markup.
**   #tag   -> highlighted span (display only)
**   @user  -> link to that user's profile (/profile/u/:username)
*/

import { Link } from 'react-router-dom'

// Capturing group so String.split keeps the tokens as their own array entries.
const SPLIT = /(#\w+|@\w+)/g
// Separate, non-global (stateless) matchers for classifying each split part.
const IS_TAG = /^#\w+$/
const IS_MENTION = /^@\w+$/

const accent = { color: '#534ab7', fontWeight: 600 }

export default function RichText({ text, className }) {
  if (!text) return null
  return (
    <span className={className}>
      {text.split(SPLIT).map((part, i) => {
        if (IS_TAG.test(part)) {
          return <span key={i} style={accent}>{part}</span>
        }
        if (IS_MENTION.test(part)) {
          return (
            // stopPropagation so the click doesn't also trigger an enclosing
            // card/post click handler.
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
