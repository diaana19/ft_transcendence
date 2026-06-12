package repositories

import (
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

// NotificationRepositories handles the notifications in the database.
type NotificationRepositories struct {
	db *gorm.DB
}

// NewNotificationRepositories creates a new NotificationRepositories using the given database.
func NewNotificationRepositories(db *gorm.DB) *NotificationRepositories {
	return &NotificationRepositories{db: db}
}

// Create saves a new notification.
func (r *NotificationRepositories) Create(notif *models.Notification) error {
	return r.db.Create(notif).Error
}

func (r *NotificationRepositories) FindByUserID(userID string, limit int) ([]models.Notification, error) {
	var notif []models.Notification
	err := r.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(limit).
		Find(&notif).Error
	return notif, err
}

// FindUnreadByUserID returns the unread notifications of the user.
func (r *NotificationRepositories) FindUnreadByUserID(userID string) ([]models.Notification, error) {
	var notif []models.Notification
	err := r.db.Where("user_id = ? AND read = false", userID).
		Order("created_at desc").
		Find(&notif).Error
	return notif, err
}

// MarkAllReadByUserID marks all the notifications of the user as read.
func (r *NotificationRepositories) MarkAllReadByUserID(userID string) error {
	return r.db.Model(&models.Notification{}).
		Where("user_id = ? AND read = false", userID).
		Update("read", true).Error
}

// MarkReadByID marks one notification of the user as read. Marking a notification
// that no longer exists is a no-op, so a repeated click never errors.
func (r *NotificationRepositories) MarkReadByID(userID, notifID string) error {
	return r.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notifID, userID).
		Update("read", true).Error
}

func (r *NotificationRepositories) DeleteByID(userID, notifID string) error {
	return r.db.Where("id = ? AND user_id = ?", notifID, userID).
		Delete(&models.Notification{}).Error
}

func (r *NotificationRepositories) DeleteAllByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).
		Delete(&models.Notification{}).Error
}

func (r *NotificationRepositories) DeleteByActorAndType(userID, actorID, notifType string) ([]string, error) {
	var ids []string
	err := r.db.Model(&models.Notification{}).
		Where("user_id = ? AND actor_id = ? AND type = ?", userID, actorID, notifType).
		Pluck("id", &ids).Error
	if err != nil || len(ids) == 0 {
		return nil, err
	}
	err = r.db.Where("id IN ?", ids).Delete(&models.Notification{}).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// GetUsernameByID returns the username of the user with this id.
func (r *NotificationRepositories) GetUsernameByID(userID string) (string, error) {
	var user models.User
	err := r.db.Select("username").First(&user, "id = ?", userID).Error
	return user.Username, err
}
