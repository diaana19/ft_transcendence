package models

import (
	"time"

	"gorm.io/gorm"
)

// ─── Post ───────────────────────────────────────────────────────────────────

type Post struct {
	ID            string         `gorm:"primaryKey;type:varchar(36)"`
	AuthorID      string         `gorm:"type:varchar(36);not null"`
	Author        User           `gorm:"foreignKey:AuthorID"`
	Content       string         `gorm:"type:text;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	LikesCount    int            `json:"likes_count"    gorm:"default:0"`
	CommentsCount int            `json:"comments_count" gorm:"default:0"`
}

type UpdatePostInput struct {
	Content string `json:"content" binding:"required,min=1,max=280"`
}

type PostResponse struct {
	ID            string       `json:"id"`
	Content       string       `json:"content"`
	AuthorID      string       `json:"author_id"`
	Author        UserResponse `json:"author"`
	LikesCount    int          `json:"likes_count"`
	CommentsCount int          `json:"comments_count"`
	Liked         bool         `json:"liked"` // true when the requesting user already liked it
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:            p.ID,
		Content:       p.Content,
		AuthorID:      p.AuthorID,
		Author:        p.Author.ToResponse(),
		LikesCount:    p.LikesCount,
		CommentsCount: p.CommentsCount,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// ─── Like ────────────────────────────────────────────────────────────────────

// Like stores a single user↔post like relationship.
// The composite unique index prevents a user from liking the same post twice.
type Like struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	UserID    string    `gorm:"type:varchar(36);not null;uniqueIndex:idx_like_user_post"`
	User      User      `gorm:"foreignKey:UserID"`
	PostID    string    `gorm:"type:varchar(36);not null;uniqueIndex:idx_like_user_post"`
	Post      Post      `gorm:"foreignKey:PostID"`
	CreatedAt time.Time
}

type LikeResponse struct {
	PostID    string    `json:"post_id"`
	Liked     bool      `json:"liked"`
	LikesCount int      `json:"likes_count"`
}

// ─── Comment (Reply) ─────────────────────────────────────────────────────────

type Reply struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)"`
	PostID    string         `gorm:"type:varchar(36);not null;index"`
	Post      Post           `gorm:"foreignKey:PostID"`
	AuthorID  string         `gorm:"type:varchar(36);not null"`
	Author    User           `gorm:"foreignKey:AuthorID"`
	Content   string         `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type CreateCommentInput struct {
	Content string `json:"content" binding:"required,min=1,max=280"`
}

type UpdateCommentInput struct {
	Content string `json:"content" binding:"required,min=1,max=280"`
}

type CommentResponse struct {
	ID        string       `json:"id"`
	PostID    string       `json:"post_id"`
	Content   string       `json:"content"`
	AuthorID  string       `json:"author_id"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (r *Reply) ToResponse() CommentResponse {
	return CommentResponse{
		ID:        r.ID,
		PostID:    r.PostID,
		Content:   r.Content,
		AuthorID:  r.AuthorID,
		Author:    r.Author.ToResponse(),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// ─── Repost ──────────────────────────────────────────────────────────────────

type Repost struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)"`
	PostID    string         `gorm:"type:varchar(36);not null;index"`
	Post      Post           `gorm:"foreignKey:PostID"`
	AuthorID  string         `gorm:"type:varchar(36);not null"`
	Author    User           `gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type RepostResponse struct {
	ID        string       `json:"id"`
	PostID    string       `json:"post_id"`
	AuthorID  string       `json:"author_id"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
}

func (r *Repost) ToResponse() RepostResponse {
	return RepostResponse{
		ID:        r.ID,
		PostID:    r.PostID,
		AuthorID:  r.AuthorID,
		Author:    r.Author.ToResponse(),
		CreatedAt: r.CreatedAt,
	}
}
