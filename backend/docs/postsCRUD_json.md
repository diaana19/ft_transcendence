# Post API - Quick Reference

## Routes

| Method | Route | Auth | Description |
|--------|-------|------|-------------|
| GET | `/api/posts` | No | Get all posts |
| GET | `/api/posts/:id` | No | Get single post |
| POST | `/api/posts` | Yes | Create post |
| PUT | `/api/posts/:id` | Yes | Update post |
| DELETE | `/api/posts/:id` | Yes | Delete post |

---

## Request & Response Format

### 1. GET /api/posts

**Request:**
```
GET /api/posts?page=1&limit=10
```

**Response: 200 OK**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "content": "Hello world!",
    "author_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "author": {
      "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
      "username": "alice",
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

---

### 2. GET /api/posts/:id

**Request:**
```
GET /api/posts/550e8400-e29b-41d4-a716-446655440000
```

**Response: 200 OK**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Hello world!",
  "author_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "author": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "username": "alice",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=Alice",
    "name": "Alice Smith"
  },
  "likes_count": 42,
  "comments_count": 5,
  "created_at": "2026-04-15T10:30:00Z",
  "updated_at": "2026-04-15T10:30:00Z"
}
```

**Response: 404 Not Found**
```json
{
  "error": "Post not found"
}
```

---

### 3. POST /api/posts

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "content": "This is my first post!"
}
```

**Response: 201 Created**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "This is my first post!",
  "author_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "author": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "username": "alice",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=Alice",
    "name": "Alice Smith"
  },
  "likes_count": 0,
  "comments_count": 0,
  "created_at": "2026-04-15T10:30:00Z",
  "updated_at": "2026-04-15T10:30:00Z"
}
```

**Response: 400 Bad Request**
```json
{
  "error": "Content is required"
}
```

**Response: 401 Unauthorized**
```json
{
  "error": "Unauthorized"
}
```

---

### 4. PUT /api/posts/:id

**Headers:**
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "content": "Updated content"
}
```

**Response: 200 OK**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "content": "Updated content",
  "author_id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
  "author": {
    "id": "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
    "username": "alice",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=Alice",
    "name": "Alice Smith"
  },
  "likes_count": 42,
  "comments_count": 5,
  "created_at": "2026-04-15T10:30:00Z",
  "updated_at": "2026-04-15T11:00:00Z"
}
```

**Response: 403 Forbidden**
```json
{
  "error": "You can only update your own posts"
}
```

**Response: 404 Not Found**
```json
{
  "error": "Post not found"
}
```

---

### 5. DELETE /api/posts/:id

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response: 200 OK**
```json
{
  "message": "Post deleted"
}
```

**Response: 403 Forbidden**
```json
{
  "error": "You can only delete your own posts"
}
```

**Response: 404 Not Found**
```json
{
  "error": "Post not found"
}
```

---

## Post Object Fields

| Field | Type | Example |
|-------|------|---------|
| `id` | string (UUID) | `550e8400-e29b-41d4-a716-446655440000` |
| `content` | string | `"Hello world!"` |
| `author_id` | string (UUID) | `6ba7b810-9dad-11d1-80b4-00c04fd430c8` |
| `author.id` | string (UUID) | `6ba7b810-9dad-11d1-80b4-00c04fd430c8` |
| `author.username` | string | `alice` |
| `author.avatar` | string (URL) | `https://api.dicebear.com/...` |
| `author.name` | string | `Alice Smith` |
| `likes_count` | integer | `42` |
| `comments_count` | integer | `5` |
| `created_at` | timestamp | `2026-04-15T10:30:00Z` |
| `updated_at` | timestamp | `2026-04-15T10:30:00Z` |

---

## Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized |
| 403 | Forbidden |
| 404 | Not Found |
| 500 | Server Error |
