package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/Transcendence/models"
	"github.com/Transcendence/repositories"
	"github.com/Transcendence/utils"
)

type AuthService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateAuthUserService(infos *models.User) (*models.UserResponse, error) {
	if infos.ID == "" {
		infos.ID = uuid.New().String()
	}

	if _, err := s.repo.GetByEmail(infos.Email); err == nil {
		return nil, errors.New("user with this email already exists")
	}

	if _, err := s.repo.GetByUsername(infos.Username); err == nil {
		return nil, errors.New("user with this username already exists")
	}

	if infos.Password == nil || *infos.Password == "" {
		return nil, errors.New("password is required")
	}

	hashed, err := utils.HashString(*infos.Password)
	if err != nil {
		return nil, err
	}
	infos.Password = &hashed

	if infos.Provider == "" {
		infos.Provider = "local"
	}

	err = s.repo.CreateUser(infos)
	if err != nil {
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
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate pending token: %w", err)
	}

	pendingToken := base64.URLEncoding.EncodeToString(bytes)

	ctx := context.Background()
	key := "pending_login:" + pendingToken
	if err := rdb.Set(ctx, key, userID, 5*time.Minute).Err(); err != nil {
		return "", fmt.Errorf("failed to store pending login: %w", err)
	}

	return pendingToken, nil
}

func (s *AuthService) ConsumePendingLogin(pendingToken string, rdb *redis.Client) (string, error) {
	ctx := context.Background()
	key := "pending_login:" + pendingToken

	userID, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", errors.New("pending login expired or invalid")
	}

	if err != nil {
		return "", fmt.Errorf("redis error: %w", err)
	}

	rdb.Del(ctx, key)
	return userID, nil
}
