package utils

import "github.com/google/uuid"

// NewID returns a new unique id using UUID v7. If it fails, it returns a UUID v4.
func NewID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}
