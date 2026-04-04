package repositories

import (
	"gorm.io/gorm"
	"github.com/Transcendence/models"
)

func CreateUser(DB *gorm.DB, user *models.User) error {
	return DB.Create(user).Error
}
