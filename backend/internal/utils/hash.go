package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashString(str string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash string: %w", err)
	}

	return string(hashedPassword), nil
}

func CheckHashString(plain, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
