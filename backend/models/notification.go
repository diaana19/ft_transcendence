package models

import "time"

type Notification struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    string    `json:"user_id"`  // who receives it
	ActorID   string    `json:"actor_id"` // who triggered it
	Type      string    `json:"type"`     // "friend_request", "like", "message", "reply"
	Content   string    `json:"content"`  // "testuser sent you a friend request"
	Read      bool      `json:"read"`
}
