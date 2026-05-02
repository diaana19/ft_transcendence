package repositories

import (
	"github.com/Transcendence/models"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(message *models.Message) error
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
