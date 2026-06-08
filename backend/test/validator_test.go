package test

import (
	"testing"
	"time"

	"ft_transcendence/backend/internal/utils"
)

func TestCheckEmailFormat_Valid(t *testing.T) {
	tests := []string{
		"user@example.com",
		"user.name@domain.com",
		"user+tag@domain.com",
	}
	for _, email := range tests {
		if !utils.CheckEmailFormat(email) {
			t.Errorf("expected valid email: %s", email)
		}
	}
}

func TestCheckEmailFormat_Invalid(t *testing.T) {
	tests := []string{
		"",
		"invalid",
		"@domain.com",
		"user@",
		"user@.com",
	}
	for _, email := range tests {
		if utils.CheckEmailFormat(email) {
			t.Errorf("expected invalid email: %s", email)
		}
	}
}

func TestCheckUsernameFormat_Valid(t *testing.T) {
	tests := []string{
		"user123",
		"test-user",
		"abc",
	}
	for _, username := range tests {
		if !utils.CheckUsernameFormat(username) {
			t.Errorf("expected valid username: %s", username)
		}
	}
}

func TestCheckUsernameFormat_Invalid(t *testing.T) {
	tests := []string{
		"",
		"-user",
		"user-",
		"user--name",
	}
	for _, username := range tests {
		if utils.CheckUsernameFormat(username) {
			t.Errorf("expected invalid username: %s", username)
		}
	}
}

func TestCheckPasswordFormat_Valid(t *testing.T) {
	tests := []struct {
		password string
		username string
	}{
		{"StrongPass123!", "user"},
		{"MyP@ssw0rd", "test"},
	}
	for _, tt := range tests {
		ok, _ := utils.CheckPasswordFormat(tt.password, tt.username)
		if !ok {
			t.Errorf("expected valid password: %s", tt.password)
		}
	}
}

func TestCheckPasswordFormat_Invalid(t *testing.T) {
	tests := []struct {
		password string
		username string
	}{
		{"", "user"},
		{"short", "user"},
		{"password", "user"},
		{"PASSWORD", "user"},
		{"12345678", "user"},
	}
	for _, tt := range tests {
		ok, _ := utils.CheckPasswordFormat(tt.password, tt.username)
		if ok {
			t.Errorf("expected invalid password: %s", tt.password)
		}
	}
}

func TestCheckUserAge_Valid(t *testing.T) {
	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if !utils.CheckUserAge(birthDate) {
		t.Error("expected valid age for 2000-01-01")
	}
}

func TestCheckUserAge_Invalid(t *testing.T) {
	birthDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if utils.CheckUserAge(birthDate) {
		t.Error("expected invalid age for 2020-01-01")
	}
}

func TestCheckUserAge_EdgeCase(t *testing.T) {
	now := time.Now()
	birthDate := time.Date(now.Year()-13, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if !utils.CheckUserAge(birthDate) {
		t.Error("expected valid age for exactly 13 years old")
	}
}
