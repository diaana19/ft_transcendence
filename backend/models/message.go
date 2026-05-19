package models

import (
	"time"
)

// Type is either "dm" or "tweet"
// it can have parentID if it "replie" to an another "tweet"
// replies is all "tweet" under the main "tweet"
type Message struct {
    ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
    CreatedAt time.Time `json:"created_at"`

    SenderID string `json:"sender_id" gorm:"not null"`
    Username string `json:"username" gorm:"not null"`
    RoomID   string `json:"room_id" gorm:"not null"`

    Type     string  `json:"type"`
    Content  string  `json:"content"`
    ParentID *string `json:"parent_id" gorm:"default:null"`
    FileID   *string `json:"file_id,omitempty" gorm:"type:uuid;default:null"`

    Replies []Message `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}
