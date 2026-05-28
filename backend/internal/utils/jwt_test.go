package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// testSecret is at least jwtMinSecretLen (32) bytes so loadSecret accepts it.
const testSecret = "test-secret-key-for-unit-tests-min-32-bytes"

func setupJWTSecret(t *testing.T) {
	t.Helper()
	os.Setenv("JWT_SECRET", testSecret)
	t.Cleanup(func() {
		os.Unsetenv("JWT_SECRET")
	})
}

func TestGenerateJWT_ReturnsToken(t *testing.T) {
	setupJWTSecret(t)

	token, err := GenerateJWT("user-123", "user-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Fatal("token should not be empty")
	}
}

func TestGenerateJWT_DifferentUsersGetDifferentTokens(t *testing.T) {
	setupJWTSecret(t)

	token1, _ := GenerateJWT("user-1", "user-1")
	token2, _ := GenerateJWT("user-2", "user-2")
	if token1 == token2 {
		t.Fatal("different users should get different tokens")
	}
}

func TestValidateJWT_ValidToken(t *testing.T) {
	setupJWTSecret(t)

	token, _ := GenerateJWT("user-123", "user-123")
	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.Subject != "user-123" {
		t.Errorf("expected sub 'user-123', got '%s'", claims.Subject)
	}
	if claims.Issuer != JWTIssuer {
		t.Errorf("expected iss %q, got %q", JWTIssuer, claims.Issuer)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	setupJWTSecret(t)
	secret := os.Getenv("JWT_SECRET")

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			Issuer:    JWTIssuer,
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
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			Issuer:    JWTIssuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("another-secret-also-at-least-32-bytes-long-ok"))

	_, err := ValidateJWT(tokenStr)
	if err == nil {
		t.Fatal("token with wrong signature should fail validation")
	}
}

// TestValidateJWT_WrongIssuer ensures that a token signed with the right
// secret but a different issuer is rejected — defends against accidental
// cross-service token reuse.
func TestValidateJWT_WrongIssuer(t *testing.T) {
	setupJWTSecret(t)
	secret := os.Getenv("JWT_SECRET")

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			Issuer:    "some-other-service",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))

	if _, err := ValidateJWT(tokenStr); err == nil {
		t.Fatal("token with wrong issuer should fail validation")
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

// TestLoadSecret_EmptyEnvFailsFast guards against the silent fail-open where
// an unset JWT_SECRET would otherwise sign / verify with []byte("") and
// accept any forged token.
func TestLoadSecret_EmptyEnvFailsFast(t *testing.T) {
	os.Unsetenv("JWT_SECRET")

	if _, err := GenerateJWT("x", "x"); err == nil {
		t.Fatal("GenerateJWT must refuse an empty JWT_SECRET")
	}
	if _, err := ValidateJWT("anything"); err == nil {
		t.Fatal("ValidateJWT must refuse an empty JWT_SECRET")
	}
}

// TestLoadSecret_ShortEnvFailsFast covers the under-32-byte case. HS256 keys
// shorter than the hash output are weak; refuse them outright.
func TestLoadSecret_ShortEnvFailsFast(t *testing.T) {
	os.Setenv("JWT_SECRET", "too-short")
	t.Cleanup(func() { os.Unsetenv("JWT_SECRET") })

	if _, err := GenerateJWT("x", "x"); err == nil {
		t.Fatal("GenerateJWT must refuse an under-32-byte JWT_SECRET")
	}
}

func TestRefreshToken_FreshToken_ReturnsSameToken(t *testing.T) {
	setupJWTSecret(t)

	token, _ := GenerateJWT("user-123", "user-123")
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
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			Issuer:    JWTIssuer,
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
	if newClaims.Subject != "user-123" {
		t.Error("refreshed token should preserve sub")
	}
}

func TestRefreshToken_ExpiredToken_ReturnsError(t *testing.T) {
	setupJWTSecret(t)
	secret := os.Getenv("JWT_SECRET")

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-123",
			Issuer:    JWTIssuer,
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
