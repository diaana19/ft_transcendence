package services

import (
	"errors"
	"testing"

	"ft_transcendence/backend/internal/models"
)

func ptrStr(s string) *string {
	return &s
}

// AuthService happy paths (registration, login by email/username, duplicate
// email/username, wrong password, unknown user) are all exercised end-to-end
// through the /auth/register and /auth/login endpoints. Only the branches that
// cannot be reached over HTTP are unit-tested here with a mocked repository.

func TestCreateAuthUserService_RepoError(t *testing.T) {
	repo := newMockRepo()
	repo.err = errors.New("db connection failed")
	svc := NewAuthService(repo)

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: ptrStr("StrongPass123!"),
	}

	if _, err := svc.CreateAuthUserService(user); err == nil {
		t.Fatal("should propagate repository create error")
	}
}

func TestCreateAuthUserService_PasswordRequired(t *testing.T) {
	svc := NewAuthService(newMockRepo())

	// the HTTP binding enforces a non-empty password, so this defensive branch
	// is only reachable below the controller layer.
	user := &models.User{Username: "nopass", Email: "nopass@example.com"}
	if _, err := svc.CreateAuthUserService(user); err == nil {
		t.Fatal("nil password should be rejected")
	}
}
