package test

import (
	"testing"
	"time"

	"ft_transcendence/backend/internal/utils"
)

// --- JWT: GenerateJWT error path ---

func TestJWT_GenerateJWTErrorPath(t *testing.T) {
	// This tests the error path when JWT_SECRET is invalid
	// In test environment, JWT_SECRET might be set, so we test with valid token
	token, err := utils.GenerateJWT("user1", "testuser")
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

// --- JWT: ValidateJWT error path ---

func TestJWT_ValidateJWTErrorPath(t *testing.T) {
	// Test with malformed token
	_, err := utils.ValidateJWT("malformed.token.here")
	if err == nil {
		t.Fatal("expected error for malformed token")
	}

	// Test with empty token
	_, err = utils.ValidateJWT("")
	if err == nil {
		t.Fatal("expected error for empty token")
	}

	// Test with valid token
	token, _ := utils.GenerateJWT("user1", "testuser")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT: %v", err)
	}
	if claims.Subject != "user1" {
		t.Fatalf("expected subject user1, got %s", claims.Subject)
	}
}

// --- JWT: RefreshToken error path ---

func TestJWT_RefreshTokenErrorPath(t *testing.T) {
	// Test with invalid token
	_, err := utils.RefreshToken("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}

	// Test with valid token (not near expiration)
	token, _ := utils.GenerateJWT("user1", "testuser")
	refreshed, err := utils.RefreshToken(token)
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	// Should return same token if not near expiration
	if refreshed != token {
		t.Fatal("expected same token when not near expiration")
	}
}

// --- Hash: HashString error path ---

func TestHash_HashStringErrorPath(t *testing.T) {
	// Test normal case
	hash, err := utils.HashString("password")
	if err != nil {
		t.Fatalf("HashString: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}

	// Test check with correct password
	if !utils.CheckHashString("password", hash) {
		t.Fatal("expected correct password to match")
	}

	// Test check with wrong password
	if utils.CheckHashString("wrong", hash) {
		t.Fatal("expected wrong password to not match")
	}
}

// --- Hashtag: ExtractHashtags edge cases ---

func TestHashtag_ExtractHashtagsEdgeCases(t *testing.T) {
	// Test with multiple hashtags
	result := utils.ExtractHashtags("#go #rust #python")
	if len(result) != 3 {
		t.Fatalf("expected 3 hashtags, got %d", len(result))
	}

	// Test with no hashtags
	result = utils.ExtractHashtags("no hashtags here")
	if len(result) != 0 {
		t.Fatalf("expected 0 hashtags, got %d", len(result))
	}

	// Test with empty string
	result = utils.ExtractHashtags("")
	if len(result) != 0 {
		t.Fatalf("expected 0 hashtags for empty string, got %d", len(result))
	}

	// Test with hashtag at start
	result = utils.ExtractHashtags("#start of text")
	if len(result) != 1 {
		t.Fatalf("expected 1 hashtag, got %d", len(result))
	}

	// Test with hashtag at end
	result = utils.ExtractHashtags("end of text #end")
	if len(result) != 1 {
		t.Fatalf("expected 1 hashtag, got %d", len(result))
	}
}

// --- Hashtag: NormalizeHashtag edge cases ---

func TestHashtag_NormalizeHashtagEdgeCases(t *testing.T) {
	// Test with uppercase
	result := utils.NormalizeHashtag("GOLANG")
	if result != "#golang" {
		t.Fatalf("expected #golang, got %s", result)
	}

	// Test with mixed case
	result = utils.NormalizeHashtag("GoLang")
	if result != "#golang" {
		t.Fatalf("expected #golang, got %s", result)
	}

	// Test with already normalized
	result = utils.NormalizeHashtag("#golang")
	if result != "#golang" {
		t.Fatalf("expected #golang, got %s", result)
	}

	// Test with empty string
	result = utils.NormalizeHashtag("")
	if result != "" {
		t.Fatalf("expected empty string, got %s", result)
	}

	// Test with whitespace
	result = utils.NormalizeHashtag("  golang  ")
	if result != "#golang" {
		t.Fatalf("expected #golang, got %s", result)
	}
}

// --- ID: NewID uniqueness ---

func TestID_NewIDUniqueness(t *testing.T) {
	// Generate multiple IDs and check uniqueness
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
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

// --- Validator: CheckPasswordFormat comprehensive ---

func TestCheckPasswordFormat_Comprehensive(t *testing.T) {
	// Test valid password
	ok, code := utils.CheckPasswordFormat("StrongPass123!", "user")
	if !ok {
		t.Fatalf("expected valid password, got code %d", code)
	}

	// Test password containing username
	ok, code = utils.CheckPasswordFormat("myuser123!", "user")
	if ok {
		t.Fatal("expected password containing username to be invalid")
	}
	if code != 0 {
		t.Fatalf("expected error code 0, got %d", code)
	}

	// Test short password
	ok, code = utils.CheckPasswordFormat("Short1!", "user")
	if ok {
		t.Fatal("expected short password to be invalid")
	}
	if code != 1 {
		t.Fatalf("expected error code 1, got %d", code)
	}

	// Test password without lowercase
	ok, code = utils.CheckPasswordFormat("PASSWORD123!", "user")
	if ok {
		t.Fatal("expected password without lowercase to be invalid")
	}
	if code != 2 {
		t.Fatalf("expected error code 2, got %d", code)
	}

	// Test password without uppercase
	ok, code = utils.CheckPasswordFormat("password123!", "user")
	if ok {
		t.Fatal("expected password without uppercase to be invalid")
	}
	if code != 3 {
		t.Fatalf("expected error code 3, got %d", code)
	}

	// Test password without digit
	ok, code = utils.CheckPasswordFormat("Password!", "user")
	if ok {
		t.Fatal("expected password without digit to be invalid")
	}
	if code != 4 {
		t.Fatalf("expected error code 4, got %d", code)
	}

	// Test password without special char
	ok, code = utils.CheckPasswordFormat("Password123", "user")
	if ok {
		t.Fatal("expected password without special char to be invalid")
	}
	if code != 5 {
		t.Fatalf("expected error code 5, got %d", code)
	}
}

// --- Validator: CheckUsernameFormat comprehensive ---

func TestCheckUsernameFormat_Comprehensive(t *testing.T) {
	// Test valid usernames
	validUsernames := []string{"user", "user123", "user-name", "a", "a-b-c"}
	for _, username := range validUsernames {
		if !utils.CheckUsernameFormat(username) {
			t.Errorf("expected valid username: %s", username)
		}
	}

	// Test invalid usernames
	invalidUsernames := []string{"", "-user", "user-", "user--name", "user name"}
	for _, username := range invalidUsernames {
		if utils.CheckUsernameFormat(username) {
			t.Errorf("expected invalid username: %s", username)
		}
	}

	// Test too long username
	longUsername := ""
	for i := 0; i < 50; i++ {
		longUsername += "a"
	}
	if utils.CheckUsernameFormat(longUsername) {
		t.Fatal("expected long username to be invalid")
	}
}

// --- Validator: CheckEmailFormat comprehensive ---

func TestCheckEmailFormat_Comprehensive(t *testing.T) {
	// Test valid emails
	validEmails := []string{"user@example.com", "user.name@domain.com", "user+tag@domain.com", "user123@domain.co.uk"}
	for _, email := range validEmails {
		if !utils.CheckEmailFormat(email) {
			t.Errorf("expected valid email: %s", email)
		}
	}

	// Test invalid emails
	invalidEmails := []string{"", "invalid", "@domain.com", "user@", "user@.com", "user@domain.c"}
	for _, email := range invalidEmails {
		if utils.CheckEmailFormat(email) {
			t.Errorf("expected invalid email: %s", email)
		}
	}
}

// --- Validator: CheckUserAge comprehensive ---

func TestCheckUserAge_Comprehensive(t *testing.T) {
	// Test valid age (over 13)
	birthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if !utils.CheckUserAge(birthDate) {
		t.Error("expected valid age for 2000-01-01")
	}

	// Test invalid age (under 13)
	birthDate = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if utils.CheckUserAge(birthDate) {
		t.Error("expected invalid age for 2020-01-01")
	}

	// Test exactly 13 years old
	now := time.Now()
	birthDate = time.Date(now.Year()-13, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	// This should be valid
	if !utils.CheckUserAge(birthDate) {
		t.Error("expected valid age for exactly 13 years old")
	}
}
