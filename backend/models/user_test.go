package models

import (
	"testing"
	"time"
)

func TestToResponse_MapsFieldsCorrectly(t *testing.T) {
	now := time.Now()
	user := User{
		ID:        "abc-123",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "should-not-appear",
		Name:      "Test User",
		Bio:       "A bio",
		Avatar:    "avatar.jpg",
		Wallpaper: "wall.jpg",
		CreatedAt: now,
	}

	resp := user.ToResponse()

	if resp.ID != user.ID {
		t.Errorf("expected ID %s, got %s", user.ID, resp.ID)
	}
	if resp.Username != user.Username {
		t.Errorf("expected Username %s, got %s", user.Username, resp.Username)
	}
	if resp.Email != user.Email {
		t.Errorf("expected Email %s, got %s", user.Email, resp.Email)
	}
	if resp.Name != user.Name {
		t.Errorf("expected Name %s, got %s", user.Name, resp.Name)
	}
	if resp.Bio != user.Bio {
		t.Errorf("expected Bio %s, got %s", user.Bio, resp.Bio)
	}
	if resp.Avatar != user.Avatar {
		t.Errorf("expected Avatar %s, got %s", user.Avatar, resp.Avatar)
	}
	if resp.Wallpaper != user.Wallpaper {
		t.Errorf("expected Wallpaper %s, got %s", user.Wallpaper, resp.Wallpaper)
	}
	if !resp.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt %v, got %v", now, resp.CreatedAt)
	}
}

func TestToResponse_ExcludesPassword(t *testing.T) {
	user := User{
		ID:       "abc-123",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "secret-password",
	}

	resp := user.ToResponse()

	// UserResponse struct has no Password field, so this is a compile-time guarantee.
	// We verify the response only contains expected fields.
	if resp.ID == "" || resp.Username == "" || resp.Email == "" {
		t.Error("response should contain ID, Username, and Email")
	}
}

func TestToResponse_EmptyFields(t *testing.T) {
	user := User{
		ID:       "abc-123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	resp := user.ToResponse()

	if resp.Name != "" {
		t.Error("empty Name should remain empty")
	}
	if resp.Bio != "" {
		t.Error("empty Bio should remain empty")
	}
	if resp.Avatar != "" {
		t.Error("empty Avatar should remain empty")
	}
	if resp.Wallpaper != "" {
		t.Error("empty Wallpaper should remain empty")
	}
}
