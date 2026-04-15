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

    LikesCount    int    `json:"likes_count" gorm:"default:0"`
    CommentsCount int    `json:"comments_count" gorm:"default:0"`
}

type UpdatePostInput struct {
    Content string `json:"content" binding:"required,max=280"`
}

type PostResponse struct {
    ID        string       `json:"id"`
    Content   string       `json:"content"`
    AuthorID  string       `json:"author_id"`
    Author    UserResponse `json:"author"`
    CreatedAt time.Time    `json:"created_at"`
    UpdatedAt time.Time    `json:"updated_at"`
}
