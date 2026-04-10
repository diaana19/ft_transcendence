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

func (s *AuthService) CreateAuthUserService(infos *models.User) (*models.User, error) {
	log.Printf("DEBUG: Starting CreateAuthUserService with: %+v\n", infos)

	if infos.ID == "" {
		infos.ID = uuid.New().String()
		log.Printf("DEBUG: Generated UUID: %s\n", infos.ID)
	}

	var err error
	result, err := s.repo.GetByEmail(infos.Email)
	if err == nil && result != nil {
		log.Printf("DEBUG: Email already exists\n")
		return nil, errors.New("user already exist")
	}

	result, err = s.repo.GetByUsername(infos.Username)
	if err == nil && result != nil {
		log.Printf("DEBUG: Username already exists\n")
		return nil, errors.New("user already exist")
	}

	log.Printf("DEBUG: Hashing password\n")
	infos.Password, err = utils.HashString(infos.Password)
	if err != nil {
		log.Printf("DEBUG: Error hashing password: %v\n", err)
		return nil, err
	}

	log.Printf("DEBUG: Creating user in DB with: %+v\n", infos)
	err = s.repo.CreateUser(infos)
	if err != nil {
		log.Printf("DEBUG: Error creating user: %v\n", err)
		return nil, err
	}

	log.Printf("DEBUG: User created successfully\n")
	infos.Password = ""

	return infos, nil
}

func (s *AuthService) LoginAuthUserService(identifier, password string) (*models.User, error) {
	user, err := s.repo.GetByIdentifier(identifier)
	if err != nil {
		return nil, errors.New("invalid credential")
	}
	if !utils.CheckHashString(password, user.Password) {
		return nil, errors.New("invalid credential")
	}

	user.Password = ""
	return user, nil
}
