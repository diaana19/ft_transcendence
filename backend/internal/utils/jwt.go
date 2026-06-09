package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTIssuer = "ft_transcendence"

	jwtTTL = 24 * time.Hour

	jwtMinSecretLen = 32
)

// Claims is the data we keep inside the JWT token.
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// loadSecret reads the JWT secret from the env. It must have at least 32 chars.
func loadSecret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if len(s) < jwtMinSecretLen {
		return nil, fmt.Errorf("JWT_SECRET must be set and at least %d characters", jwtMinSecretLen)
	}
	return []byte(s), nil
}

// GenerateJWT creates a new signed token for the user. It is valid for 24 hours.
func GenerateJWT(userId string, username string) (string, error) {
	secret, err := loadSecret()
	if err != nil {
		return "", err
	}
	now := time.Now()
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			Issuer:    JWTIssuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}
	return signed, nil
}

// ValidateJWT checks the token and returns its claims. It fails if the token is not valid.
func ValidateJWT(tokenStr string) (*Claims, error) {
	secret, err := loadSecret()
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return secret, nil
		},
		jwt.WithIssuer(JWTIssuer),
	)
	if err != nil {
		return nil, fmt.Errorf("parse jwt: %w", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid Token")
	}
	return claims, nil
}

// RefreshToken returns a new token when the old one is near to expire. If not, it returns the same token.
func RefreshToken(tokenStr string) (string, error) {
	claims, err := ValidateJWT(tokenStr)
	if err != nil {
		return "", err
	}
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return tokenStr, nil
	}
	return GenerateJWT(claims.Subject, claims.Username)
}
