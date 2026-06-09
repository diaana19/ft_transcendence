/*
 ** File: username.js
 ** Description: Single source of truth for the username rule, used everywhere on
 ** the frontend (registration validation + @mention detection).
 **
 ** GitHub's rule: alphanumeric + single hyphens, no leading/trailing or
 ** consecutive hyphens, 1-39 chars.
 **   https://github.com/shinnn/github-username-regex
 */

// The username body without anchors (so it can be embedded in other patterns,
// e.g. an @mention token). The lookahead makes a hyphen only match when followed
// by an alphanumeric, which rules out trailing/consecutive hyphens.
const BODY = '[a-z\\d](?:[a-z\\d]|-(?=[a-z\\d])){0,38}'

// Full-string validator (registration).
export const USERNAME_RE = new RegExp(`^${BODY}$`, 'i')

// Matches an @mention whose username is valid, anywhere in a string (global).
export const MENTION_RE = new RegExp(`@${BODY}`, 'gi')

// Anchored single-mention matcher, for classifying one token.
export const MENTION_ONE_RE = new RegExp(`^@${BODY}$`, 'i')
