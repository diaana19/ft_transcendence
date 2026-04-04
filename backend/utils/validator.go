package utils

import (
	"regexp"
	"strings"
)

func CheckEmailFormat(email string) bool {
	re := regexp.MustCompile(`(?i)[A-Za-z]+@[A-Za-z]+\\.[A-Za-z]+`)

	// parse email
	if !re.MatchString(email) {
		return false
	}

	return true
}

func CheckPasswordFormat(password string, username string, name string) (bool, int) {
	re := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d).{8,}$`)

	if strings.Contains(password, username) || strings.Contains(password, name) {
		return false, 1
	}

	if len(password) < 8 || !re.MatchString(password) {
		return false, 2
	}

	return true, 0
}