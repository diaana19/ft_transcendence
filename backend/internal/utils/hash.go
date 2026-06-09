package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashString returns the bcrypt hash of the given string.
func HashString(str string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash string: %w", err)
	}

	return string(hashedPassword), nil
}

// CheckHashString returns true if the plain text is equal to the hash.
func CheckHashString(plain, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
