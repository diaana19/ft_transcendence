import axiosInstance from '../../services/axiosInstance'

function getToken() {
    return localStorage.getItem('token')
}

// getPosts returns a page of feed posts using limit/offset pagination.
export async function getPosts(limit = 10, offset = 0) {
    const response = await axiosInstance.get(`/api/posts?limit=${limit}&offset=${offset}`)
    return response.data.data
}

// getPost returns one post by id.
export async function getPost(id) {
    const response = await axiosInstance.get(`/api/posts/${id}`)
    return response.data
}

// getPostsByAuthor returns the posts of one user.
export async function getPostsByAuthor(id) {
    const response = await axiosInstance.get(`/api/posts/user/${id}`)
    return response.data.data
}

// getPostsByTag returns a page of posts for one hashtag.
export async function getPostsByTag(tag, limit = 10, offset = 0) {
    // Drop a leading # and encode the tag for the query string.
    const name = encodeURIComponent(String(tag).replace(/^#/, ''))
    const response = await axiosInstance.get(
        `/api/posts?tag=${name}&limit=${limit}&offset=${offset}`
    )
    return response.data.data
}

// createPost creates a new post with optional media.
export async function createPost(content, mediaUrl = null, mediaMime = null) {
    const response = await axiosInstance.post('/api/posts', {
        content,
        media_url: mediaUrl,
        media_mime: mediaMime,
    })
    return response.data
}

// updatePost edits a post. Media fields are only sent when given.
export async function updatePost(id, content, mediaUrl = undefined, mediaMime = undefined) {
    const response = await axiosInstance.put(`/api/posts/${id}`, {
        content,
        ...(mediaUrl !== undefined && { media_url: mediaUrl }),
        ...(mediaMime !== undefined && { media_mime: mediaMime }),
    })
    return response.data
}

// deletePost removes a post by id.
export async function deletePost(id) {
    const response = await axiosInstance.delete(`/api/posts/${id}`)
    return response.data
}

// reactToPost sends a reaction (like/dislike) to a post.
export async function reactToPost(postId, value) {
    const response = await axiosInstance.post(`/api/posts/${postId}/react`, { value })
    return response.data
}

// reactToComment sends a reaction to a comment.
export async function reactToComment(postId, commentId, value) {
    const response = await axiosInstance.post(`/api/posts/${postId}/comments/${commentId}/react`, {
        value,
    })
    return response.data
}

// createComment adds a comment to a post.
export async function createComment(postId, content) {
    const response = await axiosInstance.post(`/api/posts/${postId}/comments`, { content })

    return response.data
}

// getComments returns all comments of a post.
export async function getComments(postId) {
    const response = await axiosInstance.get(`/api/posts/${postId}/comments`)
    return response.data
}

// updateComment edits a comment.
export async function updateComment(postId, commentId, content) {
    const res = await axiosInstance.put(`/api/posts/${postId}/comments/${commentId}`, { content })
    return res.data
}

// deleteComment removes a comment.
export async function deleteComment(postId, commentId) {
    const res = await axiosInstance.delete(`/api/posts/${postId}/comments/${commentId}`)
    return res.data
}

// getRepliesByUser returns the posts where a user has replied.
export const getRepliesByUser = async (userId) => {
    const { data } = await axiosInstance.get(`/api/posts?repliedBy=${userId}`)
    return data.data
}
