/*
 ** File: postService.js
 ** Description: Handles all API calls related to posts
 ** Responsibilities:
 ** - Fetch all posts with pagination
 ** - Fetch single post by ID
 ** - Create, update, and delete posts
 */

import axiosInstance from '../../services/axiosInstance'

function getToken() {
    return localStorage.getItem('token')
}

export async function getPosts(limit = 10, offset = 0) {
    const response = await axiosInstance.get(`/api/posts?limit=${limit}&offset=${offset}`)
    return response.data.data
}

export async function getPost(id) {
    const response = await axiosInstance.get(`/api/posts/${id}`)
    return response.data
}

export async function getPostsByAuthor(id) {
    const response = await axiosInstance.get(`/api/posts/user/${id}`)
    return response.data.data
}

export async function getPostsByTag(tag, limit = 10, offset = 0) {
    const name = encodeURIComponent(String(tag).replace(/^#/, ''))
    const response = await axiosInstance.get(
        `/api/posts?tag=${name}&limit=${limit}&offset=${offset}`
    )
    return response.data.data
}

export async function createPost(content, mediaUrl = null, mediaMime = null) {
    const response = await axiosInstance.post('/api/posts', {
        content,
        media_url: mediaUrl,
        media_mime: mediaMime,
    })
    return response.data
}

export async function updatePost(id, content, mediaUrl = undefined, mediaMime = undefined) {
    const response = await axiosInstance.put(`/api/posts/${id}`, {
        content,
        ...(mediaUrl !== undefined && { media_url: mediaUrl }),
        ...(mediaMime !== undefined && { media_mime: mediaMime }),
    })
    return response.data
}

export async function deletePost(id) {
    const response = await axiosInstance.delete(`/api/posts/${id}`)
    return response.data
}

export async function reactToPost(postId, value) {
    const response = await axiosInstance.post(`/api/posts/${postId}/react`, { value })
    return response.data
}

export async function reactToComment(postId, commentId, value) {
    const response = await axiosInstance.post(`/api/posts/${postId}/comments/${commentId}/react`, {
        value,
    })
    return response.data
}

export async function createComment(postId, content) {
    const response = await axiosInstance.post(`/api/posts/${postId}/comments`, { content })

    return response.data
}

export async function getComments(postId) {
    const response = await axiosInstance.get(`/api/posts/${postId}/comments`)
    return response.data
}

export async function updateComment(postId, commentId, content) {
    const res = await axiosInstance.put(`/api/posts/${postId}/comments/${commentId}`, { content })
    return res.data
}

export async function deleteComment(postId, commentId) {
    const res = await axiosInstance.delete(`/api/posts/${postId}/comments/${commentId}`)
    return res.data
}

export const getRepliesByUser = async (userId) => {
  const { data } = await api.get(`/api/posts?repliedBy=${userId}`)
  return data.data
}
