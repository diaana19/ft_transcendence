package models

import "time"

// Notification is an alert sent to a user about an action of another user.
type Notification struct {
	ID            string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UserID        string    `gorm:"not null" json:"user_id"`
	UserUsername  string    `gorm:"column:user_username" json:"user_username"`
	ActorID       string    `gorm:"not null" json:"actor_id"`
	ActorUsername string    `gorm:"column:actor_username" json:"actor_username"`
	Type          string    `gorm:"not null" json:"type"`
	Content       string    `json:"content"`
	PostID        string    `gorm:"column:post_id" json:"post_id"`
	Read          bool      `gorm:"default:false" json:"read"`
}
