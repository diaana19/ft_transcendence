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
    PostID    string    `gorm:"type:varchar(36);not null;index"`  // Quel post on répond
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
    PostID    string    `gorm:"type:varchar(36);not null;index"`  // Quel post on partage
    Post      Post      `gorm:"foreignKey:PostID"`
    AuthorID  string    `gorm:"type:varchar(36);not null"`
    Author    User      `gorm:"foreignKey:AuthorID"`
    CreatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
