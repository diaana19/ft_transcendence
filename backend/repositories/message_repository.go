package repositories

import (
	"github.com/Transcendence/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	DB *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (r *MessageRepository) Save(msg *models.Message) error {
	return r.DB.Create(msg).Error
}

func (r *MessageRepository) GetByRoomID(roomID string, limit int) ([]models.Message, error) {
	var message []models.Message
	err := r.DB.Where("room_id = ? AND parent_id is NULL", roomID).
		Order("created_at desc").
		Limit(limit).
		Find(&message).Error
	return message, err
}

func (r *MessageRepository) GetReplies(parentId string) ([]models.Message, error) {
	var replies []models.Message
	err := r.DB.Where("parent_id = ?", parentId).
		Order("created_at asc").
		Find(&replies).Error
	return replies, err
}

