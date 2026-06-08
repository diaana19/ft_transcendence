package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	pwLower   = "abcdefghijkmnopqrstuvwxyz"
	pwUpper   = "ABCDEFGHJKLMNPQRSTUVWXYZ"
	pwDigit   = "23456789"
	pwSpecial = "!@#$%^&*-_+=?"
)

func GenerateRandomPassword() (string, error) {
	const length = 16
	all := pwLower + pwUpper + pwDigit + pwSpecial

	chars := []string{pwLower, pwUpper, pwDigit, pwSpecial}
	out := make([]byte, 0, length)
	for _, set := range chars {
		ch, err := randomChar(set)
		if err != nil {
			return "", err
		}
		out = append(out, ch)
	}
	for len(out) < length {
		ch, err := randomChar(all)
		if err != nil {
			return "", err
		}
		out = append(out, ch)
	}

	for i := len(out) - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		j := jBig.Int64()
		out[i], out[j] = out[j], out[i]
	}

	return string(out), nil
}

func randomChar(set string) (byte, error) {
	idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
	if err != nil {
		return 0, err
	}
	return set[idx.Int64()], nil
}
