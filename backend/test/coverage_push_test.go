package test

import (
	"strings"
	"testing"
	"time"

	"ft_transcendence/backend/internal/utils"
)

func TestJWT_GenerateJWTErrorPath(t *testing.T) {
	token, err := utils.GenerateJWT("user1", "testuser")
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestJWT_ValidateJWTErrorPath(t *testing.T) {
	_, err := utils.ValidateJWT("malformed.token.here")
	if err == nil {
		t.Fatal("expected error for malformed token")
	}

	_, err = utils.ValidateJWT("")
	if err == nil {
		t.Fatal("expected error for empty token")
	}

	token, _ := utils.GenerateJWT("user1", "testuser")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT: %v", err)
	}
	if claims.Subject != "user1" {
		t.Fatalf("expected subject user1, got %s", claims.Subject)
	}
}

func TestJWT_RefreshTokenErrorPath(t *testing.T) {
	_, err := utils.RefreshToken("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}

	token, _ := utils.GenerateJWT("user1", "testuser")
	refreshed, err := utils.RefreshToken(token)
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	if refreshed != token {
		t.Fatal("expected same token when not near expiration")
	}
}

func TestHash_HashStringErrorPath(t *testing.T) {
	hash, err := utils.HashString("password")
	if err != nil {
		t.Fatalf("HashString: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}

	if !utils.CheckHashString("password", hash) {
		t.Fatal("expected correct password to match")
	}

	if utils.CheckHashString("wrong", hash) {
		t.Fatal("expected wrong password to not match")
	}
}

func TestHashtag_ExtractHashtagsEdgeCases(t *testing.T) {
	result := utils.ExtractHashtags("#go #rust #python")
	if len(result) != 3 {
		t.Fatalf("expected 3 hashtags, got %d", len(result))
	}

	result = utils.ExtractHashtags("no hashtags here")
	if len(result) != 0 {
		t.Fatalf("expected 0 hashtags, got %d", len(result))
	}

	result = utils.ExtractHashtags("")
	if len(result) != 0 {
		t.Fatalf("expected 0 hashtags for empty string, got %d", len(result))
	}

	result = utils.ExtractHashtags("#start of text")
	if len(result) != 1 {
		t.Fatalf("expected 1 hashtag, got %d", len(result))
	}

	result = utils.ExtractHashtags("end of text #end")
	if len(result) != 1 {
		t.Fatalf("expected 1 hashtag, got %d", len(result))
	}
}

func TestHashtag_NormalizeHashtagEdgeCases(t *testing.T) {
	result := utils.NormalizeHashtag("GOLANG")
	if result != "#golang" {
		t.Fatalf("expected #golang, got %s", result)
	}

	result = utils.NormalizeHashtag("GoLang")
	if result != "#golang" {
		t.Fatalf("expected #golang, got %s", result)
	}

	result = utils.NormalizeHashtag("#golang")
	if result != "#golang" {
		t.Fatalf("expected #golang, got %s", result)
	}

	result = utils.NormalizeHashtag("")
	if result != "" {
		t.Fatalf("expected empty string, got %s", result)
	}

	result = utils.NormalizeHashtag("  golang  ")
	if result != "#golang" {
		t.Fatalf("expected #golang, got %s", result)
	}
}

func TestID_NewIDUniqueness(t *testing.T) {
	ids := make(map[string]bool)
	for range 100 {
		id := utils.NewID()
		if id == "" {
			t.Fatal("expected non-empty ID")
		}
		if ids[id] {
			t.Fatalf("duplicate ID generated: %s", id)
		}
		ids[id] = true
	}
}

func TestCheckPasswordFormat_Comprehensive(t *testing.T) {
	ok, code := utils.CheckPasswordFormat("StrongPass123!", "user")
	if !ok {
		t.Fatalf("expected valid password, got code %d", code)
	}

	ok, code = utils.CheckPasswordFormat("myuser123!", "user")
	if ok {
		t.Fatal("expected password containing username to be invalid")
	}
	if code != 0 {
		t.Fatalf("expected error code 0, got %d", code)
	}

	ok, code = utils.CheckPasswordFormat("Short1!", "user")
	if ok {
		t.Fatal("expected short password to be invalid")
	}
	if code != 1 {
		t.Fatalf("expected error code 1, got %d", code)
	}

	ok, code = utils.CheckPasswordFormat("PASSWORD123!", "user")
	if ok {
		t.Fatal("expected password without lowercase to be invalid")
	}
	if code != 2 {
		t.Fatalf("expected error code 2, got %d", code)
	}

	ok, code = utils.CheckPasswordFormat("password123!", "user")
	if ok {
		t.Fatal("expected password without uppercase to be invalid")
	}
	if code != 3 {
		t.Fatalf("expected error code 3, got %d", code)
	}

	ok, code = utils.CheckPasswordFormat("Password!", "user")
	if ok {
		t.Fatal("expected password without digit to be invalid")
	}
	if code != 4 {
		t.Fatalf("expected error code 4, got %d", code)
	}

	ok, code = utils.CheckPasswordFormat("Password123", "user")
	if ok {
		t.Fatal("expected password without special char to be invalid")
	}
	if code != 5 {
		t.Fatalf("expected error code 5, got %d", code)
	}
}

func TestCheckUsernameFormat_Comprehensive(t *testing.T) {
	validUsernames := []string{"user", "user123", "user-name", "a", "a-b-c"}
	for _, username := range validUsernames {
		if !utils.CheckUsernameFormat(username) {
			t.Errorf("expected valid username: %s", username)
		}
	}

	invalidUsernames := []string{"", "-user", "user-", "user--name", "user name"}
	for _, username := range invalidUsernames {
		if utils.CheckUsernameFormat(username) {
			t.Errorf("expected invalid username: %s", username)
		}
	}

	var longUsername strings.Builder
	for range 50 {
		longUsername.WriteString("a")
	}
	if utils.CheckUsernameFormat(longUsername.String()) {
		t.Fatal("expected long username to be invalid")
	}
}

func TestCheckEmailFormat_Comprehensive(t *testing.T) {
	validEmails := []string{"user@example.com", "user.name@domain.com", "user+tag@domain.com", "user123@domain.co.uk"}
	for _, email := range validEmails {
		if !utils.CheckEmailFormat(email) {
			t.Errorf("expected valid email: %s", email)
		}
	}

	invalidEmails := []string{"", "invalid", "@domain.com", "user@", "user@.com", "user@domain.c"}
	for _, email := range invalidEmails {
		if utils.CheckEmailFormat(email) {
			t.Errorf("expected invalid email: %s", email)
		}
	}
}

func TestCheckUserAge_Comprehensive(t *testing.T) {
	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if !utils.CheckUserAge(birthDate) {
		t.Error("expected valid age for 2000-01-01")
	}

	birthDate = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if utils.CheckUserAge(birthDate) {
		t.Error("expected invalid age for 2020-01-01")
	}

	now := time.Now()
	birthDate = time.Date(now.Year()-13, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	if !utils.CheckUserAge(birthDate) {
		t.Error("expected valid age for exactly 13 years old")
	}
}
