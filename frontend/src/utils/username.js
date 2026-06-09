// BODY is the shared pattern for a valid username (letters, digits and dashes).
const BODY = '[a-z\\d](?:[a-z\\d]|-(?=[a-z\\d])){0,38}'

// USERNAME_RE matches a whole string that is a valid username.
export const USERNAME_RE = new RegExp(`^${BODY}$`, 'i')

// MENTION_RE finds all @mentions inside a text.
export const MENTION_RE = new RegExp(`@${BODY}`, 'gi')

// MENTION_ONE_RE matches a string that is exactly one @mention.
export const MENTION_ONE_RE = new RegExp(`^@${BODY}$`, 'i')
