package services

import (
	"github.com/Transcendence/models"
	"github.com/Transcendence/repositories"
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

func (s *UserService) UpdateUser(id string, input models.UpdateUserInput) (*models.User, error) {
	return s.repo.Update(id, input)
}

func (s *UserService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
