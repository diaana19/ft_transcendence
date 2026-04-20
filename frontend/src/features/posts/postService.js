/*
** File: postService.js
** Description: Handles all API calls related to posts
** Responsibilities:
** - Fetch all posts with pagination
** - Fetch single post by ID
** - Create, update, and delete posts
*/

import axios from 'axios'

function getToken() {
    return localStorage.getItem('token')
}

export async function getPosts(page = 1, limit = 10) {
    const response = await axios.get(`/api/posts?page=${page}&limit=${limit}`)
    return response.data.data
}

export async function getPost(id) {
    const response = await axios.get(`/api/posts/${id}`)
    return response.data.data
}

export async function createPost(content) {
    const response = await axios.post(
        '/api/posts',
        { content },
        {
            headers: {
                Authorization : `Bearer ${getToken()}`
            }
        }
    )
    return response.data.data
}

export async function updatePost(id, content) {
    const response = await axios.put(
        `/api/posts/${id}`,
        { content },
        {
            headers: {
                Authorization : `Bearer ${getToken()}`
            }
        }
    )
    return response.data.data
}

export async function deletePost(id) {
    const response = await axios.delete(
        `/api/posts/${id}`,
        {
            headers: {
                Authorization : `Bearer ${getToken()}`
            }
        }
    )
    return response.data.data
}
