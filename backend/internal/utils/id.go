package utils

import "github.com/google/uuid"

// NewID returns a time-sortable UUIDv7 string. v7 embeds a millisecond
// timestamp in its high bits, so IDs created later sort lexicographically after
// earlier ones — convenient for chronological ordering and keyset pagination,
// while remaining a standard 128-bit UUID (the DB column stays `uuid`). Falls
// back to a random v4 only if the clock/entropy source ever makes NewV7 fail.
func NewID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New().String()
	}
	return id.String()
}
