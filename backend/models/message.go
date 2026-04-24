package models

import "time"

type Message struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	SenderID  string    `json:"sender_id" gorm:"not null"`
	RoomID    string    `json:"room_id" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	ParentID  *string   `json:"parent_id" gorm:"default:null"`
	Replies   []Message `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}
