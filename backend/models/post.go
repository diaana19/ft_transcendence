package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	AuthorID  string    `gorm:"type:varchar(36);not null"`
	Author    User      `gorm:"foreignKey:AuthorID"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	LikesCount int `json:"likes_count" gorm:"default:0"`
}

type Reply struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	PostID    string    `gorm:"type:varchar(36);not null;index"`
	Post      Post      `gorm:"foreignKey:PostID"`
	AuthorID  string    `gorm:"type:varchar(36);not null"`
	Author    User      `gorm:"foreignKey:AuthorID"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Repost struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	PostID    string    `gorm:"type:varchar(36);not null;index"`
	Post      Post      `gorm:"foreignKey:PostID"`
	AuthorID  string    `gorm:"type:varchar(36);not null"`
	Author    User      `gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UpdatePostInput struct {
	Content string `json:"content" binding:"required,min=1,max=280"`
}

type PostResponse struct {
	ID        string       `json:"id"`
	Content   string       `json:"content"`
	AuthorID  string       `json:"author_id"`
	Author    UserResponse `json:"author"`
	LikesCount int         `json:"likes_count"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (p *Post) ToResponse() PostResponse {
	return PostResponse{
		ID:        p.ID,
		Content:   p.Content,
		AuthorID:  p.AuthorID,
		Author:    p.Author.ToResponse(),
		LikesCount: p.LikesCount,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

type UpdateReplyInput struct {
	Content string `json:"content" binding:"required,min=1,max=280"`
}

type ReplyResponse struct {
	ID        string       `json:"id"`
	PostID    string       `json:"post_id"`
	Content   string       `json:"content"`
	AuthorID  string       `json:"author_id"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func (r *Reply) ToResponse() ReplyResponse {
	return ReplyResponse{
		ID:        r.ID,
		PostID:    r.PostID,
		Content:   r.Content,
		AuthorID:  r.AuthorID,
		Author:    r.Author.ToResponse(),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
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
