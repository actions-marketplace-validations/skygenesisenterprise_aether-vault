package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateID() string {
	return uuid.New().String()
}

func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}

	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}

	return string(bytes), nil
}

func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(ctx, timeout)
	}
	return ctx, func() {}
}

func ContextWithMetadata(ctx context.Context, metadata map[string]string) context.Context {
	if metadata == nil {
		return ctx
	}

	newCtx := ctx
	for key, value := range metadata {
		newCtx = context.WithValue(newCtx, contextKey(key), value)
	}

	return newCtx
}

func GetMetadataFromContext(ctx context.Context, key string) (string, bool) {
	if val, ok := ctx.Value(contextKey(key)).(string); ok {
		return val, true
	}
	return "", false
}

type contextKey string

func RetryWithBackoff(ctx context.Context, maxRetries int, baseDelay, maxDelay time.Duration, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := calculateBackoff(attempt-1, baseDelay, maxDelay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		if err := fn(); err != nil {
			lastErr = err
			if attempt == maxRetries {
				return lastErr
			}
			continue
		}

		return nil
	}

	return lastErr
}

func calculateBackoff(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	backoff := time.Duration(float64(baseDelay) * float64(attempt+1))
	if backoff > maxDelay {
		backoff = maxDelay
	}
	return backoff
}

func ValidateSecretName(name string) error {
	if name == "" {
		return fmt.Errorf("secret name cannot be empty")
	}

	if len(name) > 255 {
		return fmt.Errorf("secret name cannot exceed 255 characters")
	}

	return nil
}

func ValidatePolicyEffect(effect string) error {
	switch effect {
	case "allow", "deny":
		return nil
	default:
		return fmt.Errorf("policy effect must be 'allow' or 'deny'")
	}
}

func SanitizeTags(tags map[string]string) map[string]string {
	if tags == nil {
		return nil
	}

	sanitized := make(map[string]string)
	for key, value := range tags {
		if key != "" && value != "" {
			sanitized[key] = value
		}
	}

	return sanitized
}

func TimePtr(t time.Time) *time.Time {
	return &t
}

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
}

func ParseDuration(s string) (time.Duration, error) {
	if s == "" {
		return 0, nil
	}
	return time.ParseDuration(s)
}

func FormatDuration(d time.Duration) string {
	if d == 0 {
		return ""
	}
	return d.String()
}

func IsExpired(expiry *time.Time) bool {
	if expiry == nil {
		return false
	}
	return time.Now().After(*expiry)
}

func TimeUntil(expiry *time.Time) time.Duration {
	if expiry == nil {
		return time.Duration(0)
	}
	return time.Until(*expiry)
}
