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

export async function createPost(content, mediaUrl = null) {
    const response = await axiosInstance.post(
        '/api/posts',
        {content,
        media_url: mediaUrl }
    )
    return response.data
}

export async function updatePost(id, content,  mediaUrl = undefined) {
    const response = await axiosInstance.put(
        `/api/posts/${id}`,
        { content, ...(mediaUrl !== undefined && { media_url: mediaUrl }) }
    )
    return response.data
}

export async function deletePost(id) {
    const response = await axiosInstance.delete(
        `/api/posts/${id}`
    )
    return response.data
}

// React to a post: value 1 likes, -1 dislikes. Pressing the reaction you
// already have clears it. Returns { user_reaction, likes_count, dislikes_count }.
export async function reactToPost(postId, value) {
    const response = await axiosInstance.post(`/api/posts/${postId}/react`, { value })
    return response.data
}

// React to a comment, same contract as reactToPost.
export async function reactToComment(postId, commentId, value) {
    const response = await axiosInstance.post(
        `/api/posts/${postId}/comments/${commentId}/react`,
        { value }
    )
    return response.data
}

export async function createComment(postId, content) {
    const response = await axiosInstance.post(
        `/api/posts/${postId}/comments`,
        { content }
    )

    return response.data
}

export async function getComments(postId) {
    const response = await axiosInstance.get(
        `/api/posts/${postId}/comments`
    )
    return response.data
}
