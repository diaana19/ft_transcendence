package services

import (
	"errors"
	"testing"

	"github.com/Transcendence/models"
)

func TestGetUsers_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewUserService(repo)

	repo.users["1"] = &models.User{ID: "1", Username: "user1", Email: "u1@test.com"}
	repo.users["2"] = &models.User{ID: "2", Username: "user2", Email: "u2@test.com"}

	users, err := svc.GetUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestGetUsers_Empty(t *testing.T) {
	repo := newMockRepo()
	svc := NewUserService(repo)

	users, err := svc.GetUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

func TestGetUsers_RepoError(t *testing.T) {
	repo := newMockRepo()
	repo.err = errors.New("db error")
	svc := NewUserService(repo)

	_, err := svc.GetUsers()
	if err == nil {
		t.Fatal("should propagate repository error")
	}
}

func TestGetUser_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewUserService(repo)

	repo.users["user-1"] = &models.User{ID: "user-1", Username: "testuser"}

	user, err := svc.GetUser("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "user-1" {
		t.Errorf("expected ID 'user-1', got '%s'", user.ID)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewUserService(repo)

	_, err := svc.GetUser("nonexistent")
	if err == nil {
		t.Fatal("should return error for nonexistent user")
	}
}

func TestUpdateUser_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewUserService(repo)

	repo.users["user-1"] = &models.User{
		ID:       "user-1",
		Username: "oldname",
		Email:    "old@test.com",
	}

	input := models.UpdateUserInput{
		Name: "New Name",
		Bio:  "New bio",
	}

	updated, err := svc.UpdateUser("user-1", input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Name != "New Name" {
		t.Errorf("expected name 'New Name', got '%s'", updated.Name)
	}
	if updated.Bio != "New bio" {
		t.Errorf("expected bio 'New bio', got '%s'", updated.Bio)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewUserService(repo)

	input := models.UpdateUserInput{Name: "test"}
	_, err := svc.UpdateUser("nonexistent", input)
	if err == nil {
		t.Fatal("should return error for nonexistent user")
	}
}

func TestDeleteUser_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewUserService(repo)

	repo.users["user-1"] = &models.User{ID: "user-1", Username: "testuser"}

	err := svc.DeleteUser("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, exists := repo.users["user-1"]; exists {
		t.Error("user should be removed from repository")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewUserService(repo)

	err := svc.DeleteUser("nonexistent")
	if err == nil {
		t.Fatal("should return error for nonexistent user")
	}
}

func TestUpdateUser_PartialUpdate(t *testing.T) {
	repo := newMockRepo()
	svc := NewUserService(repo)

	repo.users["user-1"] = &models.User{
		ID:       "user-1",
		Username: "original",
		Email:    "original@test.com",
		Bio:      "original bio",
	}

	input := models.UpdateUserInput{Bio: "updated bio"}

	updated, err := svc.UpdateUser("user-1", input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Bio != "updated bio" {
		t.Errorf("expected bio 'updated bio', got '%s'", updated.Bio)
	}
	if updated.Username != "original" {
		t.Error("username should not change on partial update")
	}
}
