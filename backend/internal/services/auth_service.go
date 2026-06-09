package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/utils"
)

// AuthService handles register, login, logout and password reset.
type AuthService struct {
	repo repositories.UserRepository
}

// NewAuthService creates a new AuthService.
func NewAuthService(repo repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// CreateAuthUserService registers a new user. It checks for duplicate email and username and hashes the password.
func (s *AuthService) CreateAuthUserService(user *models.User) (*models.UserResponse, error) {
	if user.ID == "" {
		user.ID = utils.NewID()
	}

	if _, err := s.repo.GetByEmail(user.Email); err == nil {
		return nil, errors.New("user with this email already exists")
	}

	if _, err := s.repo.GetByUsername(user.Username); err == nil {
		return nil, errors.New("user with this username already exists")
	}

	if user.Password == nil || *user.Password == "" {
		return nil, errors.New("password is required")
	}

	hashed, err := utils.HashString(*user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = &hashed

	if user.Provider == "" {
		user.Provider = "local"
	}

	if user.DisplayName == "" {
		user.DisplayName = user.Username
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	response := models.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
	}

	return &response, nil
}

// LoginAuthUserService checks the credentials and returns the user without the password.
func (s *AuthService) LoginAuthUserService(identifier, password string) (*models.User, error) {
	user, err := s.repo.GetByIdentifier(identifier)
	if err != nil {
		return nil, errors.New("invalid credential")
	}

	if user.Password == nil || *user.Password == "" {
		return nil, errors.New("invalid credential")
	}

	if !utils.CheckHashString(password, *user.Password) {
		return nil, errors.New("invalid credential")
	}

	user.Password = nil
	return user, nil
}

// LogoutAuthUserService puts the token in the redis blacklist until it expires.
func (s *AuthService) LogoutAuthUserService(token string, expire time.Duration, rdb *redis.Client) error {
	ctx := context.Background()
	err := rdb.Set(ctx, "blacklist:"+token, "1", expire).Err()
	if err != nil {
		return errors.New("could not logout")
	}
	return nil
}

// GetUserByID returns one user by ID.
func (s *AuthService) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

// GetUserByEmail returns one user by email.
func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetByEmail(email)
}

// ResetPassword generates a new password, delivers it with deliver and stores the hash.
func (s *AuthService) ResetPassword(userID string, deliver func(newPassword string) error) error {
	newPassword, err := utils.GenerateRandomPassword()
	if err != nil {
		return fmt.Errorf("generate password: %w", err)
	}
	if err = deliver(newPassword); err != nil {
		return fmt.Errorf("deliver password: %w", err)
	}
	hashed, err := utils.HashString(newPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	if err = s.repo.UpdatePassword(userID, hashed); err != nil {
		return fmt.Errorf("persist password: %w", err)
	}
	return nil
}

// CreatePendingLogin stores a pending login in redis and returns its token. It expires in 5 minutes.
func (s *AuthService) CreatePendingLogin(userID string, rdb *redis.Client) (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate pending token: %w", err)
	}

	pendingToken := base64.URLEncoding.EncodeToString(randomBytes)

	ctx := context.Background()
	key := "pending_login:" + pendingToken
	if err := rdb.Set(ctx, key, userID, 5*time.Minute).Err(); err != nil {
		return "", fmt.Errorf("failed to store pending login: %w", err)
	}

	return pendingToken, nil
}

// PeekPendingLogin returns the user ID of a pending login without removing it.
func (s *AuthService) PeekPendingLogin(pendingToken string, rdb *redis.Client) (string, error) {
	ctx := context.Background()
	key := "pending_login:" + pendingToken

	userID, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("pending login expired or invalid")
	}

	if err != nil {
		return "", fmt.Errorf("redis error: %w", err)
	}

	return userID, nil
}

// ConsumePendingLogin removes a pending login from redis.
func (s *AuthService) ConsumePendingLogin(pendingToken string, rdb *redis.Client) {
	ctx := context.Background()
	key := "pending_login:" + pendingToken
	rdb.Del(ctx, key)
}
