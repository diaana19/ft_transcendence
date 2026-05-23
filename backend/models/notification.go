package models

import "time"

type Notification struct {
	ID            string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UserID        string    `gorm:"not null" json:"user_id"` // who receives it
	UserUsername  string    `gorm:"column:user_username" json:"user_username"`
	ActorID       string    `gorm:"not null" json:"actor_id"` // who triggered it
	ActorUsername string    `gorm:"column:actor_username" json:"actor_username"`
	Type          string    `gorm:"not null" json:"type"` // "friend_request", "like", "message", "reply"
	Content       string    `json:"content"`              // "testuser sent you a friend request"
	Read          bool      `gorm:"default:false" json:"read"`
}
