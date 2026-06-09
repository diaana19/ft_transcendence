import { describe, it, expect, vi, beforeEach } from 'vitest'

// Mock the shared axios instance so no real HTTP call is made.
vi.mock('../../services/axiosInstance', () => ({
    default: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() },
}))

import axiosInstance from '../../services/axiosInstance'
import {
    getPosts,
    getPost,
    getPostsByAuthor,
    getPostsByTag,
    createPost,
    updatePost,
    deletePost,
    reactToPost,
    reactToComment,
    createComment,
    getComments,
    updateComment,
    deleteComment,
    getRepliesByUser,
} from './postService'

beforeEach(() => {
    vi.clearAllMocks()
})

describe('postService', () => {
    it('getPosts asks the feed with limit/offset and unwraps data.data', async () => {
        axiosInstance.get.mockResolvedValue({ data: { data: ['p1'] } })
        const res = await getPosts(5, 10)
        expect(axiosInstance.get).toHaveBeenCalledWith('/api/posts?limit=5&offset=10')
        expect(res).toEqual(['p1'])
    })

    it('getPosts uses default limit and offset', async () => {
        axiosInstance.get.mockResolvedValue({ data: { data: [] } })
        await getPosts()
        expect(axiosInstance.get).toHaveBeenCalledWith('/api/posts?limit=10&offset=0')
    })

    it('getPost returns the full data object', async () => {
        axiosInstance.get.mockResolvedValue({ data: { id: '1' } })
        const res = await getPost('1')
        expect(axiosInstance.get).toHaveBeenCalledWith('/api/posts/1')
        expect(res).toEqual({ id: '1' })
    })

    it('getPostsByAuthor unwraps data.data', async () => {
        axiosInstance.get.mockResolvedValue({ data: { data: ['a'] } })
        const res = await getPostsByAuthor('u1')
        expect(axiosInstance.get).toHaveBeenCalledWith('/api/posts/user/u1')
        expect(res).toEqual(['a'])
    })

    it('getPostsByTag drops the leading # and encodes the tag', async () => {
        axiosInstance.get.mockResolvedValue({ data: { data: [] } })
        await getPostsByTag('#hello world', 3, 6)
        expect(axiosInstance.get).toHaveBeenCalledWith(
            '/api/posts?tag=hello%20world&limit=3&offset=6'
        )
    })

    it('createPost sends content and media fields', async () => {
        axiosInstance.post.mockResolvedValue({ data: { id: 'x' } })
        const res = await createPost('hi', 'url', 'image/png')
        expect(axiosInstance.post).toHaveBeenCalledWith('/api/posts', {
            content: 'hi',
            media_url: 'url',
            media_mime: 'image/png',
        })
        expect(res).toEqual({ id: 'x' })
    })

    it('updatePost only sends media fields when they are given', async () => {
        axiosInstance.put.mockResolvedValue({ data: {} })
        await updatePost('1', 'new')
        expect(axiosInstance.put).toHaveBeenCalledWith('/api/posts/1', { content: 'new' })

        await updatePost('1', 'new', 'u', 'm')
        expect(axiosInstance.put).toHaveBeenLastCalledWith('/api/posts/1', {
            content: 'new',
            media_url: 'u',
            media_mime: 'm',
        })
    })

    it('deletePost calls delete on the right url', async () => {
        axiosInstance.delete.mockResolvedValue({ data: { ok: true } })
        const res = await deletePost('7')
        expect(axiosInstance.delete).toHaveBeenCalledWith('/api/posts/7')
        expect(res).toEqual({ ok: true })
    })

    it('reactToPost sends the value', async () => {
        axiosInstance.post.mockResolvedValue({ data: {} })
        await reactToPost('1', 1)
        expect(axiosInstance.post).toHaveBeenCalledWith('/api/posts/1/react', { value: 1 })
    })

    it('reactToComment sends the value to the comment url', async () => {
        axiosInstance.post.mockResolvedValue({ data: {} })
        await reactToComment('1', 'c2', -1)
        expect(axiosInstance.post).toHaveBeenCalledWith('/api/posts/1/comments/c2/react', {
            value: -1,
        })
    })

    it('createComment posts the content', async () => {
        axiosInstance.post.mockResolvedValue({ data: { id: 'c1' } })
        await createComment('1', 'nice')
        expect(axiosInstance.post).toHaveBeenCalledWith('/api/posts/1/comments', {
            content: 'nice',
        })
    })

    it('getComments returns the data', async () => {
        axiosInstance.get.mockResolvedValue({ data: ['c1', 'c2'] })
        const res = await getComments('1')
        expect(axiosInstance.get).toHaveBeenCalledWith('/api/posts/1/comments')
        expect(res).toEqual(['c1', 'c2'])
    })

    it('updateComment puts the new content', async () => {
        axiosInstance.put.mockResolvedValue({ data: {} })
        await updateComment('1', 'c2', 'edit')
        expect(axiosInstance.put).toHaveBeenCalledWith('/api/posts/1/comments/c2', {
            content: 'edit',
        })
    })

    it('deleteComment removes the comment', async () => {
        axiosInstance.delete.mockResolvedValue({ data: {} })
        await deleteComment('1', 'c2')
        expect(axiosInstance.delete).toHaveBeenCalledWith('/api/posts/1/comments/c2')
    })

    // Regression: getRepliesByUser used an undefined "api" variable and crashed.
    it('getRepliesByUser uses axiosInstance with the repliedBy filter', async () => {
        axiosInstance.get.mockResolvedValue({ data: { data: ['r1'] } })
        const res = await getRepliesByUser('u42')
        expect(axiosInstance.get).toHaveBeenCalledWith('/api/posts?repliedBy=u42')
        expect(res).toEqual(['r1'])
    })
})
