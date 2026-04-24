# Post API Documentation

## Overview

This document describes the Post CRUD API endpoints. All protected endpoints require authentication via JWT token in the `Authorization` header.

---

## Authentication

All protected endpoints require:
```
Authorization: Bearer <jwt_token>
```

The JWT token is obtained from:
- `POST /api/auth/register` - Returns token on registration
- `POST /api/auth/login` - Returns token on login

---

## Endpoints

### 1. Get All Posts

**Endpoint:** `GET /api/posts`

**Authentication:** Optional (public endpoint)

**Description:** Retrieve all posts with pagination support.

**Query Parameters:**
```
?page=1          (default: 1)
?limit=10        (default: 10, max: 100)
```

**Example Request:**
```bash
curl http://localhost:8000/api/posts?page=1&limit=10
```

**Success Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "content": "Hello world!",
    "author_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "author": {
      "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "username": "alice",
      "email": "alice@example.com",
      "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=Alice",
      "name": "Alice Smith"
    },
    "likes_count": 42,
    "comments_count": 5,
    "created_at": "2026-04-15T10:30:00Z",
    "updated_at": "2026-04-15T10:30:00Z"
  }
]
```

**Error Response:** `500 Internal Server Error`
```json
{
  "error": "Failed to retrieve posts"
}
```

---

### 2. Get Single Post

**Endpoint:** `GET /api/posts/:id`

**Authentication:** Optional (public endpoint)

**Description:** Retrieve a specific post by ID.

**Path Parameters:**
```
:id - Post UUID (required)
```

**Example Request:**
```bash
curl http://localhost:8000/api/posts/550e8400-e29b-41d4-a716-446655440000
```

**Success Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Hello world!",
  "author_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "author": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "username": "alice",
    "email": "alice@example.com",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=Alice",
    "name": "Alice Smith"
  },
  "likes_count": 42,
  "comments_count": 5,
  "created_at": "2026-04-15T10:30:00Z",
  "updated_at": "2026-04-15T10:30:00Z"
}
```

**Error Responses:**

`404 Not Found`
```json
{
  "error": "Post not found"
}
```

`500 Internal Server Error`
```json
{
  "error": "Failed to retrieve post"
}
```

---

### 3. Create Post

**Endpoint:** `POST /api/posts`

**Authentication:** Required ⭐

**Description:** Create a new post. Only authenticated users can create posts.

**Request Body:**
```json
{
  "content": "This is my first post!"
}
```

**Validation Rules:**
- `content` is required
- `content` must not exceed 280 characters
- `content` must not be empty

**Example Request:**
```bash
curl -X POST http://localhost:8000/api/posts \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{"content": "Hello world!"}'
```

**Success Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Hello world!",
  "author_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "author": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "username": "alice",
    "email": "alice@example.com",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=Alice",
    "name": "Alice Smith"
  },
  "likes_count": 0,
  "comments_count": 0,
  "created_at": "2026-04-15T10:30:00Z",
  "updated_at": "2026-04-15T10:30:00Z"
}
```

**Error Responses:**

`400 Bad Request`
```json
{
  "error": "Content is required"
}
```

`400 Bad Request`
```json
{
  "error": "Content must not exceed 280 characters"
}
```

`401 Unauthorized`
```json
{
  "error": "Unauthorized"
}
```

`500 Internal Server Error`
```json
{
  "error": "Failed to create post"
}
```

---

### 4. Update Post

**Endpoint:** `PUT /api/posts/:id`

**Authentication:** Required ⭐

**Description:** Update a post. Only the post author can update their own posts.

**Path Parameters:**
```
:id - Post UUID (required)
```

**Request Body:**
```json
{
  "content": "Updated content"
}
```

**Validation Rules:**
- `content` is required
- `content` must not exceed 280 characters
- `content` must not be empty

**Example Request:**
```bash
curl -X PUT http://localhost:8000/api/posts/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{"content": "Updated content"}'
```

**Success Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Updated content",
  "author_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "author": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "username": "alice",
    "email": "alice@example.com",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=Alice",
    "name": "Alice Smith"
  },
  "likes_count": 42,
  "comments_count": 5,
  "created_at": "2026-04-15T10:30:00Z",
  "updated_at": "2026-04-15T11:00:00Z"
}
```

**Error Responses:**

`400 Bad Request`
```json
{
  "error": "Content is required"
}
```

`401 Unauthorized`
```json
{
  "error": "Unauthorized"
}
```

`403 Forbidden`
```json
{
  "error": "You can only update your own posts"
}
```

`404 Not Found`
```json
{
  "error": "Post not found"
}
```

`500 Internal Server Error`
```json
{
  "error": "Failed to update post"
}
```

---

### 5. Delete Post

**Endpoint:** `DELETE /api/posts/:id`

**Authentication:** Required ⭐

**Description:** Delete a post. Only the post author can delete their own posts. (Soft delete - post will be marked as deleted but data retained)

**Path Parameters:**
```
:id - Post UUID (required)
```

**Example Request:**
```bash
curl -X DELETE http://localhost:8000/api/posts/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer <jwt_token>"
```

**Success Response:** `200 OK`
```json
{
  "message": "Post deleted"
}
```

**Error Responses:**

`401 Unauthorized`
```json
{
  "error": "Unauthorized"
}
```

`403 Forbidden`
```json
{
  "error": "You can only delete your own posts"
}
```

`404 Not Found`
```json
{
  "error": "Post not found"
}
```

`500 Internal Server Error`
```json
{
  "error": "Failed to delete post"
}
```

---

## Data Model

### Post Object

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "string (max 280 characters)",
  "author_id": "uuid-string",
  "author": {
    "id": "uuid-string",
    "username": "string",
    "email": "string",
    "avatar": "url",
    "name": "string"
  },
  "likes_count": "integer",
  "comments_count": "integer",
  "created_at": "ISO 8601 timestamp",
  "updated_at": "ISO 8601 timestamp"
}
```

### Field Descriptions

| Field | Type | Description |
|-------|------|-------------|
| `id` | string (UUID) | Unique post identifier |
| `content` | string | Post text content (max 280 chars) |
| `author_id` | string (UUID) | ID of the user who created the post |
| `author` | object | Full user object (see User API docs) |
| `likes_count` | integer | Number of likes on this post |
| `comments_count` | integer | Number of comments on this post |
| `created_at` | timestamp | When the post was created |
| `updated_at` | timestamp | When the post was last updated |

---

## Status Codes

| Code | Meaning | When |
|------|---------|------|
| `200` | OK | Request succeeded |
| `201` | Created | Post successfully created |
| `400` | Bad Request | Invalid request data |
| `401` | Unauthorized | Missing/invalid JWT token |
| `403` | Forbidden | Authenticated but not authorized (not author) |
| `404` | Not Found | Post doesn't exist |
| `500` | Server Error | Server-side error |

---

## Authentication Example

### Getting JWT Token

```bash
# Register
curl -X POST http://localhost:8000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "Alice123!secure",
    "name": "Alice Smith"
  }'

# Response:
{
  "user": { ... },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Using JWT Token

```bash
curl -X GET http://localhost:8000/api/posts \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

---

## Notes for Frontend Team

1. **Character Limit:** Posts are limited to 280 characters (like Twitter). Show character count in UI.

2. **Ownership Check:** The API will return `403 Forbidden` if user tries to edit/delete another user's post. Handle this in frontend.

3. **Pagination:** Use pagination (page, limit) for listing posts to avoid loading too many at once.

4. **Soft Delete:** Deleted posts are not returned in GET requests, but the data is kept in the database.

5. **Author Information:** The `author` object is always included in responses for easier display of user info.

6. **Timestamp Format:** All timestamps are in ISO 8601 format (UTC).

7. **Loading States:** Consider showing loading states while creating/updating/deleting posts.

8. **Error Handling:** Always check for error responses and display user-friendly messages.
