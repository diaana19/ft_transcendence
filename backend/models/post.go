package models

import (
	"time"
	"gorm.io/gorm"
)

type Post struct {
	ID            string         `gorm:"primaryKey;type:uuid"`
	AuthorID      string         `gorm:"type:uuid;not null"`
	Author        User           `gorm:"foreignKey:AuthorID;references:ID"`
	Content       string         `gorm:"type:text;not null"`
	MediaURL      *string        `gorm:"type:text" json:"media_url,omitempty"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	LikesCount    int            `json:"likes_count" gorm:"default:0"`
	CommentsCount int            `json:"comments_count" gorm:"default:0"`
}

type UpdatePostInput struct {
	Content string `json:"content" binding:"required,min=1,max=280"`
	MediaURL      *string        `gorm:"type:text" json:"media_url,omitempty"`
}

type PostResponse struct {
	ID            string       `json:"id"`
	Content       string       `json:"content"`
	MediaURL      *string        `gorm:"type:text" json:"media_url,omitempty"`
	AuthorID      string       `json:"author_id"`
	Author        UserResponse `json:"author"`
	LikesCount    int          `json:"likes_count"`
	CommentsCount int          `json:"comments_count"`
	Liked         bool         `json:"liked"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:            p.ID,
		Content:       p.Content,
		MediaURL:      p.MediaURL,
		AuthorID:      p.AuthorID,
		Author:        p.Author.ToResponse(),
		LikesCount:    p.LikesCount,
		CommentsCount: p.CommentsCount,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

type Like struct {
	ID        string    `gorm:"primaryKey;type:uuid"`
	UserID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_like_user_post"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	PostID    string    `gorm:"type:uuid;not null;uniqueIndex:idx_like_user_post"`
	Post      Post      `gorm:"foreignKey:PostID;references:ID"`
	CreatedAt time.Time
}

type LikeResponse struct {
	PostID     string `json:"post_id"`
	Liked      bool   `json:"liked"`
	LikesCount int    `json:"likes_count"`
}

type Reply struct {
	ID        string         `gorm:"primaryKey;type:uuid"`
	PostID    string         `gorm:"type:uuid;not null;index"`
	Post      Post           `gorm:"foreignKey:PostID;references:ID"`
	AuthorID  string         `gorm:"type:uuid;not null"`
	Author    User           `gorm:"foreignKey:AuthorID;references:ID"`
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

type Repost struct {
	ID        string         `gorm:"primaryKey;type:uuid"`
	PostID    string         `gorm:"type:uuid;not null;index"`
	Post      Post           `gorm:"foreignKey:PostID;references:ID"`
	AuthorID  string         `gorm:"type:uuid;not null"`
	Author    User           `gorm:"foreignKey:AuthorID;references:ID"`
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
