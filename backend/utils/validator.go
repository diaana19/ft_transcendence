package utils

import (
	"regexp"
	"strings"
	"time"
)

func CheckUserAge(birthDate time.Time) bool {
	now := time.Now()

	age := now.Year() - birthDate.Year()

	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return age >= 13
}

func CheckEmailFormat(email string) bool {
	re := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	return re.MatchString(email)
}

func CheckPasswordFormat(password string, username string) (bool, int) {
	lowerPass := strings.ToLower(strings.TrimSpace(password))
	lowerUser := strings.ToLower(strings.TrimSpace(username))

	if len(lowerUser) >= 4 {

		for i := 0; i <= len(lowerUser)-4; i++ {
			sub := lowerUser[i : i+4]
			if strings.Contains(lowerPass, sub) {
				return false, 0
			}
		}
	}

	if len(password) < 8 {
		return false, 1
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case (char >= '!' && char <= '/') || (char >= ':' && char <= '@') || (char >= '[' && char <= '`') || (char >= '{' && char <= '~'):
			hasSpecial = true
		}
	}

	switch {
	case !hasLower:
		return false, 2
	case !hasUpper:
		return false, 3
	case !hasDigit:
		return false, 4
	case !hasSpecial:
		return false, 5
	}

	return true, -1
}
