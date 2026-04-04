package services

import (
	"errors"

	"github.com/Transcendence/models"
	"github.com/Transcendence/repositories"
	"github.com/Transcendence/utils"
)

type AuthService struct {
	repo repositories.UserRepository
}

func NewAuthService( repo repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateAuthUserService(infos *models.User) (*models.User, error) {
	var err error
	result, err := s.repo.GetByEmail(infos.Email)
	if result != nil {
		return nil, errors.New("user already exist")
	}

	result, err = s.repo.GetByUsername(infos.Username)
	if result != nil {
		return nil, errors.New("user already exist")
	}

	infos.Password, err = utils.HashString(infos.Password)
	if err != nil {
		return nil, err
	}

	err = s.repo.CreateUser(infos)
	if err != nil {
		return nil, err
	}

	infos.Password = ""

	return infos, nil
}