package models

import (
	"testing"
	"time"
)

func ptrStr(s string) *string {
	return &s
}

func TestToResponse_MapsFieldsCorrectly(t *testing.T) {
	now := time.Now()
	avatar := "avatar.jpg"
	wallpaper := "wall.jpg"
	user := User{
		ID:        "abc-123",
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  ptrStr("should-not-appear"),
		Name:      "Test User",
		Bio:       "A bio",
		Avatar:    &avatar,
		Wallpaper: &wallpaper,
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

	if resp.Avatar == nil || *resp.Avatar != *user.Avatar {
		t.Errorf("expected Avatar %v, got %v", user.Avatar, resp.Avatar)
	}
	if resp.Wallpaper == nil || *resp.Wallpaper != *user.Wallpaper {
		t.Errorf("expected Wallpaper %v, got %v", user.Wallpaper, resp.Wallpaper)
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
		Password: ptrStr("secret-password"),
	}

	resp := user.ToResponse()



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



	if resp.Avatar != nil {
		t.Error("empty Avatar should remain nil")
	}
	if resp.Wallpaper != nil {
		t.Error("empty Wallpaper should remain nil")
	}
}
