package services

import (
	"errors"
	"testing"

	"github.com/Transcendence/models"
	"github.com/Transcendence/utils"
)



func ptrStr(s string) *string {
	return &s
}

func TestCreateAuthUserService_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewAuthService(repo)

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: ptrStr("StrongPass123!"),
	}

	resp, err := svc.CreateAuthUserService(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", resp.Username)
	}
	if resp.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", resp.Email)
	}
	if resp.ID == "" {
		t.Error("response should have a generated ID")
	}
}

func TestCreateAuthUserService_GeneratesUUID(t *testing.T) {
	repo := newMockRepo()
	svc := NewAuthService(repo)

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: ptrStr("StrongPass123!"),
	}

	resp, err := svc.CreateAuthUserService(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID == "" {
		t.Error("should generate UUID when ID is empty")
	}
}

func TestCreateAuthUserService_HashesPassword(t *testing.T) {
	repo := newMockRepo()
	svc := NewAuthService(repo)

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: ptrStr("StrongPass123!"),
	}

	_, err := svc.CreateAuthUserService(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stored := repo.users[user.ID]
	if stored.Password == nil {
		t.Fatal("stored password should not be nil after classic registration")
	}
	if *stored.Password == "StrongPass123!" {
		t.Error("password should be hashed, not stored in plaintext")
	}
	if !utils.CheckHashString("StrongPass123!", *stored.Password) {
		t.Error("hashed password should verify against original")
	}
}

func TestCreateAuthUserService_DuplicateEmail(t *testing.T) {
	repo := newMockRepo()
	svc := NewAuthService(repo)

	user1 := &models.User{
		Username: "user1",
		Email:    "dupe@example.com",
		Password: ptrStr("StrongPass123!"),
	}
	svc.CreateAuthUserService(user1)

	user2 := &models.User{
		Username: "user2",
		Email:    "dupe@example.com",
		Password: ptrStr("StrongPass123!"),
	}
	_, err := svc.CreateAuthUserService(user2)
	if err == nil {
		t.Fatal("should fail for duplicate email")
	}
	if err.Error() != "user with this email already exists" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestCreateAuthUserService_DuplicateUsername(t *testing.T) {
	repo := newMockRepo()
	svc := NewAuthService(repo)

	user1 := &models.User{
		Username: "sameuser",
		Email:    "user1@example.com",
		Password: ptrStr("StrongPass123!"),
	}
	svc.CreateAuthUserService(user1)

	user2 := &models.User{
		Username: "sameuser",
		Email:    "user2@example.com",
		Password: ptrStr("StrongPass123!"),
	}
	_, err := svc.CreateAuthUserService(user2)
	if err == nil {
		t.Fatal("should fail for duplicate username")
	}
	if err.Error() != "user with this username already exists" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestCreateAuthUserService_RepoError(t *testing.T) {
	repo := newMockRepo()
	repo.err = errors.New("db connection failed")
	svc := NewAuthService(repo)

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: ptrStr("StrongPass123!"),
	}

	_, err := svc.CreateAuthUserService(user)
	if err == nil {
		t.Fatal("should propagate repository error")
	}
}

func TestCreateAuthUserService_ResponseExcludesPassword(t *testing.T) {
	repo := newMockRepo()
	svc := NewAuthService(repo)

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: ptrStr("StrongPass123!"),
	}

	resp, err := svc.CreateAuthUserService(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}


	if resp.Username == "" || resp.Email == "" {
		t.Error("response should contain user info")
	}
}

func TestLoginAuthUserService_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewAuthService(repo)

	hashed, _ := utils.HashString("StrongPass123!")
	repo.users["user-1"] = &models.User{
		ID:       "user-1",
		Email:    "test@example.com",
		Username: "testuser",
		Password: &hashed,
	}

	user, err := svc.LoginAuthUserService("test@example.com", "StrongPass123!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "user-1" {
		t.Errorf("expected user ID 'user-1', got '%s'", user.ID)
	}
}

func TestLoginAuthUserService_ClearsPassword(t *testing.T) {
	repo := newMockRepo()
	svc := NewAuthService(repo)

	hashed, _ := utils.HashString("StrongPass123!")
	repo.users["user-1"] = &models.User{
		ID:       "user-1",
		Email:    "test@example.com",
		Username: "testuser",
		Password: &hashed,
	}

	user, err := svc.LoginAuthUserService("test@example.com", "StrongPass123!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.Password != nil {
		t.Error("password should be cleared (nil) in login response")
	}
}

func TestLoginAuthUserService_WrongPassword(t *testing.T) {
	repo := newMockRepo()
	svc := NewAuthService(repo)

	hashed, _ := utils.HashString("StrongPass123!")
	repo.users["user-1"] = &models.User{
		ID:       "user-1",
		Email:    "test@example.com",
		Username: "testuser",
		Password: &hashed,
	}

	_, err := svc.LoginAuthUserService("test@example.com", "WrongPassword!")
	if err == nil {
		t.Fatal("should fail with wrong password")
	}
	if err.Error() != "invalid credential" {
		t.Errorf("expected 'invalid credential', got '%s'", err.Error())
	}
}

func TestLoginAuthUserService_UserNotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewAuthService(repo)

	_, err := svc.LoginAuthUserService("nonexistent@example.com", "StrongPass123!")
	if err == nil {
		t.Fatal("should fail when user not found")
	}
	if err.Error() != "invalid credential" {
		t.Errorf("expected 'invalid credential', got '%s'", err.Error())
	}
}

func TestLoginAuthUserService_ByUsername(t *testing.T) {
	repo := newMockRepo()
	svc := NewAuthService(repo)

	hashed, _ := utils.HashString("StrongPass123!")
	repo.users["user-1"] = &models.User{
		ID:       "user-1",
		Email:    "test@example.com",
		Username: "testuser",
		Password: &hashed,
	}

	user, err := svc.LoginAuthUserService("testuser", "StrongPass123!")
	if err != nil {
		t.Fatalf("should allow login by username: %v", err)
	}
	if user.ID != "user-1" {
		t.Errorf("expected user ID 'user-1', got '%s'", user.ID)
	}
}
