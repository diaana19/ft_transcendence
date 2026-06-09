import { describe, it, expect } from 'vitest'
import { USERNAME_RE, MENTION_RE, MENTION_ONE_RE } from './username'

// Tests for the username and mention regexes.
describe('username regex', () => {
    it('accepts valid usernames', () => {
        for (const u of ['john', 'john-doe', 'a', 'user123', 'AB-cd-12']) {
            expect(USERNAME_RE.test(u)).toBe(true)
        }
    })

    it('rejects invalid usernames', () => {
        for (const u of ['', '-john', 'john-', 'jo--hn', 'jo_hn', 'jo hn', 'jo.hn']) {
            expect(USERNAME_RE.test(u)).toBe(false)
        }
    })

    it('MENTION_ONE_RE matches exactly one mention', () => {
        expect(MENTION_ONE_RE.test('@john')).toBe(true)
        expect(MENTION_ONE_RE.test('@john-doe')).toBe(true)
        expect(MENTION_ONE_RE.test('john')).toBe(false)
        expect(MENTION_ONE_RE.test('@john extra')).toBe(false)
    })

    it('MENTION_RE finds all mentions in a text', () => {
        const text = '@john said hi to @jane-doe and @bob123'
        expect(text.match(MENTION_RE)).toEqual(['@john', '@jane-doe', '@bob123'])
    })
})
