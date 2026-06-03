package repositories

import (
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
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

// MarkReadByID marks a single notification read, scoped to its owner so a user
// can only touch their own. Returns ErrRecordNotFound if nothing matched.
func (r *NotificationRepositories) MarkReadByID(userID, notifID string) error {
	res := r.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notifID, userID).
		Update("read", true)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *NotificationRepositories) GetUsernameByID(userID string) (string, error) {
	var user models.User
	err := r.db.Select("username").First(&user, "id = ?", userID).Error
	return user.Username, err
}
