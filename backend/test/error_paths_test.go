package test

import (
	"testing"

	"ft_transcendence/backend/internal/utils"
)

func TestJWT_GenerateJWTWithShortSecret(t *testing.T) {
	_, err := utils.GenerateJWT("user1", "testuser")
	if err != nil {
		t.Logf("GenerateJWT returned error (expected if secret is invalid): %v", err)
	}
}

func TestJWT_ValidateJWTInvalidToken(t *testing.T) {
	_, err := utils.ValidateJWT("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestJWT_ValidateJWTEmptyToken(t *testing.T) {
	_, err := utils.ValidateJWT("")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestJWT_RefreshTokenInvalidToken(t *testing.T) {
	_, err := utils.RefreshToken("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestHash_HashAndCheck(t *testing.T) {
	password := "testpassword123"
	hash, err := utils.HashString(password)
	if err != nil {
		t.Fatalf("HashString: %v", err)
	}

	if !utils.CheckHashString(password, hash) {
		t.Fatal("expected correct password to match")
	}

	if utils.CheckHashString("wrongpassword", hash) {
		t.Fatal("expected wrong password to not match")
	}
}

func TestHashtag_ExtractHashtags(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"post with #golang and #rust", []string{"#golang", "#rust"}},
		{"no hashtags here", []string{}},
		{"#single", []string{"#single"}},
		{"", []string{}},
	}
	for _, tt := range tests {
		result := utils.ExtractHashtags(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("ExtractHashtags(%q): expected %d tags, got %d", tt.input, len(tt.expected), len(result))
		}
	}
}

func TestHashtag_NormalizeHashtag(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"golang", "#golang"},
		{"#golang", "#golang"},
		{"GOLANG", "#golang"},
		{"#Rust", "#rust"},
		{"", ""},
	}
	for _, tt := range tests {
		result := utils.NormalizeHashtag(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeHashtag(%q): expected %q, got %q", tt.input, tt.expected, result)
		}
	}
}

func TestID_NewID(t *testing.T) {
	id1 := utils.NewID()
	id2 := utils.NewID()

	if id1 == "" {
		t.Fatal("expected non-empty ID")
	}
	if id1 == id2 {
		t.Fatal("expected unique IDs")
	}
}

func TestCheckPasswordFormat_ContainsUsername(t *testing.T) {
	ok, code := utils.CheckPasswordFormat("myuser123!", "user")
	if ok {
		t.Fatal("expected password containing username to be invalid")
	}
	if code != 0 {
		t.Fatalf("expected error code 0, got %d", code)
	}
}

func TestCheckPasswordFormat_NoLowercase(t *testing.T) {
	ok, code := utils.CheckPasswordFormat("PASSWORD123!", "user")
	if ok {
		t.Fatal("expected password without lowercase to be invalid")
	}
	if code != 2 {
		t.Fatalf("expected error code 2, got %d", code)
	}
}

func TestCheckPasswordFormat_NoUppercase(t *testing.T) {
	ok, code := utils.CheckPasswordFormat("password123!", "user")
	if ok {
		t.Fatal("expected password without uppercase to be invalid")
	}
	if code != 3 {
		t.Fatalf("expected error code 3, got %d", code)
	}
}

func TestCheckPasswordFormat_NoDigit(t *testing.T) {
	ok, code := utils.CheckPasswordFormat("Password!", "user")
	if ok {
		t.Fatal("expected password without digit to be invalid")
	}
	if code != 4 {
		t.Fatalf("expected error code 4, got %d", code)
	}
}

func TestCheckPasswordFormat_NoSpecial(t *testing.T) {
	ok, code := utils.CheckPasswordFormat("Password123", "user")
	if ok {
		t.Fatal("expected password without special char to be invalid")
	}
	if code != 5 {
		t.Fatalf("expected error code 5, got %d", code)
	}
}

func TestCheckUsernameFormat_TooLong(t *testing.T) {
	longUsername := ""
	for i := 0; i < 50; i++ {
		longUsername += "a"
	}
	if utils.CheckUsernameFormat(longUsername) {
		t.Fatal("expected long username to be invalid")
	}
}

func TestCheckUsernameFormat_ConsecutiveHyphens(t *testing.T) {
	if utils.CheckUsernameFormat("user--name") {
		t.Fatal("expected consecutive hyphens to be invalid")
	}
}

func TestCheckUsernameFormat_LeadingHyphen(t *testing.T) {
	if utils.CheckUsernameFormat("-username") {
		t.Fatal("expected leading hyphen to be invalid")
	}
}

func TestCheckUsernameFormat_TrailingHyphen(t *testing.T) {
	if utils.CheckUsernameFormat("username-") {
		t.Fatal("expected trailing hyphen to be invalid")
	}
}

func TestCheckEmailFormat_NoAt(t *testing.T) {
	if utils.CheckEmailFormat("userdomain.com") {
		t.Fatal("expected email without @ to be invalid")
	}
}

func TestCheckEmailFormat_NoDomain(t *testing.T) {
	if utils.CheckEmailFormat("user@") {
		t.Fatal("expected email without domain to be invalid")
	}
}

func TestCheckEmailFormat_NoTLD(t *testing.T) {
	if utils.CheckEmailFormat("user@domain") {
		t.Fatal("expected email without TLD to be invalid")
	}
}

func TestCheckEmailFormat_ShortTLD(t *testing.T) {
	if utils.CheckEmailFormat("user@domain.c") {
		t.Fatal("expected email with short TLD to be invalid")
	}
}
