package models

import "time"

// File visibility values. They control who can access the file.
const (
	FileVisibilityPublic  = "public"
	FileVisibilityFriends = "friends"
	FileVisibilityPrivate = "private"
)

// File holds a file uploaded by a user.
type File struct {
	ID         string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OwnerID    string    `gorm:"type:uuid;not null;index" json:"owner_id"`
	Path       string    `gorm:"type:text;not null" json:"-"`
	Filename   string    `gorm:"type:varchar(255);not null" json:"filename"`
	MimeType   string    `gorm:"type:varchar(100);not null" json:"mime_type"`
	Size       int64     `gorm:"not null" json:"size"`
	Visibility string    `gorm:"type:varchar(20);not null;default:'public'" json:"visibility"`
	CreatedAt  time.Time `json:"created_at"`
}

// FileAccess gives a user permission to access a private file.
type FileAccess struct {
	FileID string `gorm:"primaryKey;type:uuid" json:"file_id"`
	UserID string `gorm:"primaryKey;type:uuid" json:"user_id"`
}
