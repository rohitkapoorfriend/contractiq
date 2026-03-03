package identifier

import "github.com/google/uuid"

// New generates a new random UUID string.
func New() string {
	return uuid.New().String()
}

// IsValid checks if the given string is a valid UUID.
func IsValid(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
