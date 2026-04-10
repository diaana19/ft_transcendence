package services

import (
	"errors"
	"log"

	"github.com/Transcendence/models"
	"github.com/Transcendence/repositories"
	"github.com/Transcendence/utils"
	"github.com/google/uuid"
)

type AuthService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateAuthUserService(infos *models.User) (*models.UserResponse, error) {
	log.Printf("DEBUG: Starting CreateAuthUserService with: %+v\n", infos)

	if infos.ID == "" {
		infos.ID = uuid.New().String()
	}

	if _, err := s.repo.GetByEmail(infos.Email); err == nil {
		return nil, errors.New("user with this email already exists")
	}

	if _, err := s.repo.GetByUsername(infos.Username); err == nil {
		return nil, errors.New("user with this username already exists")
	}

	var err error
	infos.Password, err = utils.HashString(infos.Password)
	if err != nil {
		log.Printf("DEBUG: Error hashing password: %v\n", err)
		return nil, err
	}

	err = s.repo.CreateUser(infos)
	if err != nil {
		log.Printf("DEBUG: Error creating user: %v\n", err)
		return nil, err
	}

	response := models.UserResponse{
		ID:        infos.ID,
		Username:  infos.Username,
		Email:     infos.Email,
		CreatedAt: infos.CreatedAt,
	}

	return &response, nil
}
