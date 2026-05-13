package repositories

import (
	"github.com/Transcendence/models"
	"gorm.io/gorm"
)

type NotificationRepositories struct {
	db *gorm.DB
}

func NewNotificationRepositories(db *gorm.DB) *NotificationRepositories {
	return &NotificationRepositories{db: db}
}

func (r *NotificationRepositories) Create(notif *models.Notification) error {
	return r.db.Create(notif).Error
}

func (r *NotificationRepositories) FindUnreadByUserID(userID string) ([]models.Notification, error) {
	var notif []models.Notification
	err := r.db.Where("user_id = ? AND read = false", userID).
		Order("created_at desc").
		Find(&notif).Error
	return notif, err
}

func (r *NotificationRepositories) MarkAllReadByUserID(userID string) error {
	return r.db.Model(&models.Notification{}).
		Where("user_id = ? AND read = false", userID).
		Update("read", true).Error
}

func (r *NotificationRepositories) GetUsernameByID(userID string) (string, error) {
	var user models.User
	err := r.db.Select("username").First(&user, "id = ?", userID).Error
	return user.Username, err
}
