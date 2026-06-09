const BODY = '[a-z\\d](?:[a-z\\d]|-(?=[a-z\\d])){0,38}'

export const USERNAME_RE = new RegExp(`^${BODY}$`, 'i')

export const MENTION_RE = new RegExp(`@${BODY}`, 'gi')

export const MENTION_ONE_RE = new RegExp(`^@${BODY}$`, 'i')
