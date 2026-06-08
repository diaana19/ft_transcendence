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

type AuthService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

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

func (s *AuthService) LogoutAuthUserService(token string, expire time.Duration, rdb *redis.Client) error {
	ctx := context.Background()
	err := rdb.Set(ctx, "blacklist:"+token, "1", expire).Err()
	if err != nil {
		return errors.New("could not logout")
	}
	return nil
}

func (s *AuthService) GetUserByID(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

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

func (s *AuthService) ConsumePendingLogin(pendingToken string, rdb *redis.Client) {
	ctx := context.Background()
	key := "pending_login:" + pendingToken
	rdb.Del(ctx, key)
}
