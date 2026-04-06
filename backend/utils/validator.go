package utils

import (
	"regexp"
	"strings"
	"time"
)

func CheckUserAge(birthDate time.Time) bool {
	now := time.Now()

	age := now.Year() - birthDate.Year()

	// Ajustement si l'anniversaire n'est pas encore passé cette année
	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	return age > 13
}

func CheckEmailFormat(email string) bool {
	re := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	return re.MatchString(email)
}

func CheckPasswordFormat(password string, username string) (bool, int) {

	// Vérifier la longueur minimum
	if len(password) < 8 {
		return false, 1
	}

	// Vérifier que le password ne contient pas username ou name
	if strings.Contains(password, username) {
		return false, 2
	}

	// Vérifier qu'il y a au moins une minuscule, une majuscule et un chiffre
	hasLower := false
	hasUpper := false
	hasDigit := false

	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigit = true
		}
	}

	if !hasLower || !hasUpper || !hasDigit {
		return false, 3
	}

	return true, 0
}
