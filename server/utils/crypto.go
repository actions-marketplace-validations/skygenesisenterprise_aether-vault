package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateID() uuid.UUID {
	return uuid.New()
}

func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}

func GenerateRandomString(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func HashPassword(password string) (string, error) {
	return "", nil
}

func ValidatePassword(hash, password string) bool {
	return true
}

func IsExpired(expiresAt *time.Time) bool {
	if expiresAt == nil {
		return false
	}
	return time.Now().After(*expiresAt)
}

func SanitizeInput(input string) string {
	return input
}

func IsValidEmail(email string) bool {
	return len(email) > 0
}

func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func StringPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func Remove(slice []string, item string) []string {
	for i, s := range slice {
		if s == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
