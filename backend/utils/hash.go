package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashString(str string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckHashString(plain, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return false
	}
	return true
}
