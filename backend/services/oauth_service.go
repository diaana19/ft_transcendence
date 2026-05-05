package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Transcendence/config"
	"github.com/Transcendence/models"
	"github.com/Transcendence/repositories"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

type GitHubUser struct {
	ID			int64 `json:"id"`
	Login		string `json:"login"`
	Name		string `json:"name"`
	Email		string `json:"email"`
	AvatarURL	string `json:"avatar_url"`
}

type ghEmail struct {
	Email		string	`json:"email"`
	Primary		bool	`json:"primary"`
	Verified	bool	`json:"verified"`
}

type OAuthService struct {
	userRepo	repositories.UserRepository
	redisClient	*redis.Client
	oauthConfig	*oauth2.Config
}

func NewOAuthService(repo repositories.UserRepository, rdb *redis.Client, cfg *config.Config) *OAuthService {
	oauthConfig := &oauth2.Config {
		ClientID: cfg.GithubClientID,
		ClientSecret: cfg.GithubClientSecret,
		RedirectURL: cfg.GithubRedirectURL,
		Scopes: []string{"user:email", "read:user"},
		Endpoint: github.Endpoint,
	}

	return &OAuthService{
		userRepo: repo,
		redisClient: rdb,
		oauthConfig: oauthConfig,
	}
}

func (s *OAuthService) GenerateState(ctx context.Context) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}

	state := base64.URLEncoding.EncodeToString(bytes)
	key := "oauth:state:" + state
	if err := s.redisClient.Set(ctx, key, "1", 10*time.Minute).Err(); err != nil {
		return "", fmt.Errorf("failed to store state in redis: %w", err)
	}

	return state, nil
}

func (s *OAuthService) VerifyAndConsumeState(ctx context.Context, state string) (bool, error) {
	if state == "" {
		return false, nil
	}

	key := "oauth:state:" + state

	_, err := s.redisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("redis get failed: %w", err)
	}

	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		fmt.Printf("warning: failed to delete consumed state: %v\n", err)
	}

	return true, nil
}

func (s *OAuthService) BuildAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state)
}

func (s *OAuthService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	if code == "" {
		return nil, errors.New("authorization code is empty")
	}

	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}

	return token, nil
}

func (s *OAuthService) FetchGitHubUser(ctx context.Context, token *oauth2.Token) (*GitHubUser, error) {
	client := s.oauthConfig.Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github returned status %d when fetching user", resp.StatusCode)
	}

	var ghUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return nil, fmt.Errorf("failed to decode user response: %w", err)
	}
	if ghUser.Email == "" {
		email, err := s.fetchPrimaryVerifiedEmail(ctx, client)
		if err != nil {
			return nil, err
		}
		ghUser.Email = email
	}

	return &ghUser, nil
}

func (s *OAuthService) fetchPrimaryVerifiedEmail(ctx context.Context, client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", fmt.Errorf("failed to fetch emails: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("github returned status %d when fetching emails", resp.StatusCode)
	}

	var emails []ghEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", fmt.Errorf("failed to decode emails response: %w", err)
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	return "", errors.New("no verified primary email found on github account")
}

func (s *OAuthService) FindOrCreateUser(ctx context.Context, ghUser *GitHubUser) (*models.User, error) {
	githubIDStr := strconv.FormatInt(ghUser.ID, 10)
	user, err := s.userRepo.GetByGithubID(githubIDStr)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to lookup user by github id: %w", err)
	}

	if ghUser.Email != "" {
		user, err = s.userRepo.GetByEmail(ghUser.Email)
		if err == nil {
			user.GithubID = &githubIDStr
			user.Provider = "github"
			if err := s.userRepo.LinkGithub(user.ID, githubIDStr); err != nil {
				return nil, fmt.Errorf("failed to link github account: %w", err)
			}
			return user, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to lookup user by email: %w", err)
		}
	}

	username, err := s.findAvailableUsername(ghUser.Login)
	if err != nil {
		return nil, err
	}

	name := ghUser.Name
	if name == "" {
		name = ghUser.Login
	}

	avatar := ghUser.AvatarURL

	newUser := models.User{
		Name:     name,
		Username: username,
		Email:    ghUser.Email,
		Avatar:   &avatar,
		GithubID: &githubIDStr,
		Provider: "github",
	}

	if err := s.userRepo.CreateUser(&newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}

func (s *OAuthService) findAvailableUsername(base string) (string, error) {
	candidate := base
	for suffix := 1; suffix < 1000; suffix++ {
		_, err := s.userRepo.GetByUsername(candidate)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return candidate, nil
		}
		if err != nil {
			return "", fmt.Errorf("failed to check username availability: %w", err)
		}
		candidate = fmt.Sprintf("%s%d", base, suffix)
	}
	return "", errors.New("could not find an available username after 1000 attempts")
}
