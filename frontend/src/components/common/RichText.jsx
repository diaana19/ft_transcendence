/*
** File: RichText.jsx
** Description: Renders user text with #hashtags highlighted.
** Splits the string on a hashtag regex and renders matches as colored spans;
** plain text stays as escaped React string children (no dangerouslySetInnerHTML),
** so user content cannot inject markup.
*/

// Capturing group so String.split keeps the #tags as their own array entries.
const SPLIT = /(#\w+)/g
// Separate, non-global (stateless) matcher for classifying each split part.
const IS_TAG = /^#\w+$/

export default function RichText({ text, className }) {
  if (!text) return null
  return (
    <span className={className}>
      {text.split(SPLIT).map((part, i) =>
        IS_TAG.test(part) ? (
          <span key={i} style={{ color: '#534ab7', fontWeight: 600 }}>{part}</span>
        ) : (
          part
        )
      )}
    </span>
  )
}
