package utils

import "github.com/google/uuid"

// GenerateUUID returns a unique UUIDv4 string
func GenerateUUID() string {
	return uuid.New().String()
}
