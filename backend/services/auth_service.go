package services

import (
	"github.com/Transcendence/models"
	"gorm.io/gorm"
)

func AuthService(infos *models.User, DB *gorm.DB) (response models.User, err error) {

	// TODO: check in db if email already exist, if username already exist once all done send it to repositories/users.go

	return response, nil
}