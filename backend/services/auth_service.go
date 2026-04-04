package services

import (
	"errors"

	"github.com/Transcendence/models"
	"github.com/Transcendence/repositories"
	"github.com/Transcendence/utils"
	"gorm.io/gorm"
)

func CreateAuthUserService(infos *models.User, DB *gorm.DB) (*models.User, error) {
	var err error
	user := models.User{}
	result := DB.Where("email = ?", infos.Email).First(&user)
	if result.RowsAffected > 0 {
		return nil, errors.New("user already exist")
	}

	user = models.User{}
	result = DB.Where("username = ?", infos.Username).First(&user)
	if result.RowsAffected > 0 {
		return nil, errors.New("user already exist")
	}

	infos.Password, err = utils.HashString(infos.Password)
	if err != nil {
		return nil, err
	}

	err = repositories.CreateUser(DB, infos)
	if err != nil {
		return nil, err
	}

	infos.Password = ""

	return infos, nil
}