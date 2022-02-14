package utils

import (
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

// GenerateID generates a unique ID that can be used as an identifier for a domain.
func GenerateID() string {
	return uuid.New().String()
}

// IsValidUUID return true if uuid valid
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// ShortUUID generates a uuid and encode it with base64 encoding
func ShortUUID() string {
	uu := uuid.New()
	return base64.RawURLEncoding.EncodeToString(uu[:])
}

// DecodeUUID returns the original UUID
func DecodeUUID(encodedID string) string {
	decoded, _ := base64.RawURLEncoding.DecodeString(encodedID)
	return fmt.Sprintf("%x", decoded)
}
