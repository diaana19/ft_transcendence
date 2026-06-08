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

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func loadSecret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if len(s) < jwtMinSecretLen {
		return nil, fmt.Errorf("JWT_SECRET must be set and at least %d characters", jwtMinSecretLen)
	}
	return []byte(s), nil
}

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
