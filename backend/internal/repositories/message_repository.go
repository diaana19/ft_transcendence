package repositories

import (
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

type MessageRepository interface {
	Create(message *models.Message) error
	GetByRoomID(roomID string, since string, limit int) ([]models.Message, error)
	GetReplies(parentID string) ([]models.Message, error)
	PollSince(userID, since string, limit int) ([]models.Message, error)
	ListConversation(userID, peerID, since string, limit int) ([]models.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) GetByRoomID(roomID, since string, limit int) ([]models.Message, error) {
	q := r.db.Where("room_id = ? AND parent_id IS NULL", roomID)
	return runCursorQuery(q, since, limit) // reuses the existing cursor logic
}

// func (r *messageRepository) GetByRoomID(roomID string, limit int) ([]models.Message, error) {
// 	var messages []models.Message
// 	err := r.db.Where("room_id = ? AND parent_id IS NULL", roomID).
// 		Order("created_at desc").
// 		Limit(limit).
// 		Find(&messages).Error
// 	return messages, err
// }

func (r *messageRepository) GetReplies(parentID string) ([]models.Message, error) {
	var replies []models.Message
	err := r.db.Where("parent_id = ?", parentID).
		Order("created_at asc").
		Find(&replies).Error
	return replies, err
}

func (r *messageRepository) PollSince(userID, since string, limit int) ([]models.Message, error) {
	q := r.db.Where("sender_id = ? OR recipient_id = ?", userID, userID)
	return runCursorQuery(q, since, limit)
}

func (r *messageRepository) ListConversation(userID, peerID, since string, limit int) ([]models.Message, error) {
	q := r.db.Where(
		"(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
		userID, peerID, peerID, userID,
	)
	return runCursorQuery(q, since, limit)
}

func runCursorQuery(q *gorm.DB, since string, limit int) ([]models.Message, error) {
	var messages []models.Message
	if since == "" {
		if err := q.Order("id DESC").Limit(limit).Find(&messages).Error; err != nil {
			return nil, err
		}
		for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
			messages[i], messages[j] = messages[j], messages[i]
		}
		return messages, nil
	}
	err := q.Where("id > ?", since).Order("id ASC").Limit(limit).Find(&messages).Error
	return messages, err
}
