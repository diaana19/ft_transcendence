package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func setupJWTSecret(t *testing.T) {
	t.Helper()
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-tests")
	t.Cleanup(func() {
		os.Unsetenv("JWT_SECRET")
	})
}

func TestGenerateJWT_ReturnsToken(t *testing.T) {
	setupJWTSecret(t)

	token, err := GenerateJWT("user-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("token should not be empty")
	}
}

func TestGenerateJWT_DifferentUsersGetDifferentTokens(t *testing.T) {
	setupJWTSecret(t)

	token1, _ := GenerateJWT("user-1")
	token2, _ := GenerateJWT("user-2")
	if token1 == token2 {
		t.Fatal("different users should get different tokens")
	}
}

func TestValidateJWT_ValidToken(t *testing.T) {
	setupJWTSecret(t)

	token, _ := GenerateJWT("user-123")
	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.UserId != "user-123" {
		t.Errorf("expected userId 'user-123', got '%s'", claims.UserId)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	setupJWTSecret(t)
	secret := os.Getenv("JWT_SECRET")

	claims := Claims{
		UserId: "user-123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))

	_, err := ValidateJWT(tokenStr)
	if err == nil {
		t.Fatal("expired token should fail validation")
	}
}

func TestValidateJWT_InvalidSignature(t *testing.T) {
	setupJWTSecret(t)

	claims := Claims{
		UserId: "user-123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("wrong-secret"))

	_, err := ValidateJWT(tokenStr)
	if err == nil {
		t.Fatal("token with wrong signature should fail validation")
	}
}

func TestValidateJWT_MalformedToken(t *testing.T) {
	setupJWTSecret(t)

	_, err := ValidateJWT("not.a.valid.token")
	if err == nil {
		t.Fatal("malformed token should fail validation")
	}
}

func TestValidateJWT_EmptyToken(t *testing.T) {
	setupJWTSecret(t)

	_, err := ValidateJWT("")
	if err == nil {
		t.Fatal("empty token should fail validation")
	}
}

func TestRefreshToken_FreshToken_ReturnsSameToken(t *testing.T) {
	setupJWTSecret(t)

	token, _ := GenerateJWT("user-123")
	refreshed, err := RefreshToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if refreshed != token {
		t.Fatal("token with >1h until expiry should be returned as-is")
	}
}

func TestRefreshToken_NearExpiry_ReturnsNewToken(t *testing.T) {
	setupJWTSecret(t)
	secret := os.Getenv("JWT_SECRET")

	claims := Claims{
		UserId: "user-123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-23*time.Hour - 30*time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))

	refreshed, err := RefreshToken(tokenStr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if refreshed == tokenStr {
		t.Fatal("token near expiry should be refreshed to a new token")
	}

	newClaims, err := ValidateJWT(refreshed)
	if err != nil {
		t.Fatalf("refreshed token should be valid: %v", err)
	}
	if newClaims.UserId != "user-123" {
		t.Error("refreshed token should preserve userId")
	}
}

func TestRefreshToken_ExpiredToken_ReturnsError(t *testing.T) {
	setupJWTSecret(t)
	secret := os.Getenv("JWT_SECRET")

	claims := Claims{
		UserId: "user-123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-25 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))

	_, err := RefreshToken(tokenStr)
	if err == nil {
		t.Fatal("expired token should fail refresh")
	}
}
