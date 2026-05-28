package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// JWTIssuer is asserted in `iss` on every token. Mismatched issuers (or
	// the empty string from a token created without RegisteredClaims.Issuer
	// set) fail validation via jwt.WithIssuer below.
	JWTIssuer = "ft_transcendence"

	jwtTTL = 24 * time.Hour

	// jwtMinSecretLen reflects HS256's required key length. HMAC-SHA256 keys
	// shorter than 32 bytes weaken the signature; an empty secret makes
	// every token trivially forgeable. Treat under-32 as a config error and
	// refuse to sign or verify.
	jwtMinSecretLen = 32
)

// Claims carries application-specific fields on top of the standard claims
// from RFC 7519. The user ID lives in RegisteredClaims.Subject (the `sub`
// claim) — read it via claims.Subject. Username is kept as a custom claim so
// handlers have a display name without an extra DB lookup.
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// loadSecret returns the signing key from JWT_SECRET, or an error if the
// environment value is too short. Treating an empty/short secret as a hard
// error closes the silent fail-open where the system would otherwise sign
// and verify with []byte("") and accept any forged token.
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
