package services

import (
	"errors"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/utils"
)

// UserService handles user data operations.
type UserService struct {
	repo repositories.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetUsers returns all users.
func (s *UserService) GetUsers() ([]models.User, error) {
	return s.repo.GetAll()
}

// GetUser returns one user by ID.
func (s *UserService) GetUser(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

// SearchUsers returns a page of users matching the query and the total count.
func (s *UserService) SearchUsers(query string, limit, offset int) ([]models.User, int64, error) {
	return s.repo.SearchUsers(query, limit, offset)
}

// GetUserByUsername returns one user by username.
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.repo.GetByUsername(username)
}

// UpdateUser updates the user with the given input.
func (s *UserService) UpdateUser(id string, input models.UpdateUserInput) (*models.User, error) {
	return s.repo.Update(id, input)
}

// DeleteUser removes the user by ID.
func (s *UserService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}

// VerifyPassword checks if the password matches the user. It fails for accounts without a password.
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
