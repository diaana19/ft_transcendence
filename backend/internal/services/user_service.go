package services

import (
	"errors"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/utils"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetUser(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *UserService) UpdateUser(id string, input models.UpdateUserInput) (*models.User, error) {
	return s.repo.Update(id, input)
}

func (s *UserService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}

func (s *UserService) VerifyPassword(userID, password string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}
	if user.Password == nil || *user.Password == "" {
		return errors.New("password verification not supported for this account")
	}
	if !utils.CheckHashString(password, *user.Password) {
		return errors.New("invalid password")
	}
	return nil
}
