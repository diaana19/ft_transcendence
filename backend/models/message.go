package models

import (
	"time"
)

// Type is either "dm" or "tweet"
// it can have parentID if it "replie" to an another "tweet"
// replies is all "tweet" under the main "tweet"
type Message struct {
	Username  string `json:"username" gorm:"not null"`
	Content   string `json:"content" gorm:"not null"`
	Type      string
	CreatedAt time.Time `json:"created_at"`
	RoomID    string    `json:"room_id" gorm:"not null"`
	Replies   []Message `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
	ParentID  *string   `json:"parent_id" gorm:"default:null"`
	SenderID  string    `json:"sender_id" gorm:"not null"`
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
}
