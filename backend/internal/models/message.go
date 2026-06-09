package models

import (
	"time"
)

// Message is a chat message sent by a user.
type Message struct {
	ID        string    `gorm:"primaryKey;type:varchar(36);index:idx_msg_sender,priority:2;index:idx_msg_recipient,priority:2" json:"id"` //nolint:lll
	CreatedAt time.Time `json:"created_at"`

	SenderID string `gorm:"type:varchar(36);not null;index:idx_msg_sender,priority:1" json:"sender_id"`
	Username string `json:"username" gorm:"default:''"`

	Type        string  `json:"type,omitempty"`
	Content     string  `gorm:"type:text;not null" json:"content"`
	FileID      *string `json:"file_id,omitempty" gorm:"type:uuid;default:null"`
	RoomID      string  `json:"room_id,omitempty" gorm:"default:null"`
	ParentID    *string `json:"parent_id,omitempty" gorm:"default:null"`
	RecipientID string  `gorm:"type:varchar(36);default:null;index:idx_msg_recipient,priority:1" json:"recipient_id,omitempty"` //nolint:lll

	Replies []Message `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

// MessageResponse is the message data we send to the client.
type MessageResponse struct {
	ID          string    `json:"id"`
	SenderID    string    `json:"sender_id"`
	RecipientID string    `json:"recipient_id,omitempty"`
	Content     string    `json:"content"`
	FileID      *string   `json:"file_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToResponse converts the Message into a MessageResponse.
func (m *Message) ToResponse() MessageResponse {
	return MessageResponse{
		ID:          m.ID,
		SenderID:    m.SenderID,
		RecipientID: m.RecipientID,
		Content:     m.Content,
		FileID:      m.FileID,
		CreatedAt:   m.CreatedAt,
	}
}
