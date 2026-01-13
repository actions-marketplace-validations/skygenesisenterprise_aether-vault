package types

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// TOTPAlgorithm represents the hash algorithm used for TOTP
type TOTPAlgorithm string

const (
	// TOTPAlgorithmSHA1 uses SHA-1 hash
	TOTPAlgorithmSHA1 TOTPAlgorithm = "SHA1"

	// TOTPAlgorithmSHA256 uses SHA-256 hash
	TOTPAlgorithmSHA256 TOTPAlgorithm = "SHA256"

	// TOTPAlgorithmSHA512 uses SHA-512 hash
	TOTPAlgorithmSHA512 TOTPAlgorithm = "SHA512"
)

// TOTPEntry represents a TOTP (Time-based One-Time Password) entry
type TOTPEntry struct {
	// Unique identifier
	ID string

	// Account name (typically email)
	Account string

	// Issuer name
	Issuer string

	// Secret key (base32 encoded)
	Secret string

	// Algorithm
	Algorithm TOTPAlgorithm

	// Digits (6 or 8)
	Digits int

	// Period in seconds
	Period int

	// Metadata
	Metadata *TOTPMetadata

	// Creation timestamp
	CreatedAt time.Time

	// Last update timestamp
	UpdatedAt time.Time

	// Notes
	Notes string

	// Tags
	Tags []string

	// Folder/Path
	Folder string

	// Is favorite
	Favorite bool

	// Usage count
	UsageCount int

	// Last used timestamp
	LastUsed *time.Time
}

// TOTPMetadata contains metadata about a TOTP entry
type TOTPMetadata struct {
	// Created by
	CreatedBy string

	// Last modified by
	UpdatedBy string

	// Counter for HOTP (not used for TOTP)
	Counter int64

	// Whether it's a backup code
	IsBackupCode bool

	// QR code data
	QRCode string

	// Import source
	ImportSource string

	// Original URI
	OriginalURI string
}

// TOTPCredentials represents TOTP generation/verification credentials
type TOTPCredentials struct {
	// Secret key (base32 encoded)
	Secret string

	// Current time (for testing)
	Time *time.Time

	// Algorithm
	Algorithm TOTPAlgorithm

	// Digits
	Digits int

	// Period
	Period int

	// Counter for HOTP
	Counter int64
}

// TOTPGenerator represents a TOTP code generator
type TOTPGenerator struct {
	secret    []byte
	algorithm TOTPAlgorithm
	digits    int
	period    int
}

// TOTPVerification represents TOTP verification result
type TOTPVerification struct {
	// Generated code
	Code string

	// Time remaining (seconds)
	TimeRemaining int

	// Period
	Period int

	// Current time
	CurrentTime time.Time

	// Success
	Success bool

	// Error message
	Error string
}

// TOTPBackupCode represents a backup code
type TOTPBackupCode struct {
	// Code
	Code string

	// Used status
	Used bool

	// Usage timestamp
	UsedAt *time.Time

	// Creation timestamp
	CreatedAt time.Time

	// Notes
	Notes string
}

// TOTPImport represents TOTP import information
type TOTPImport struct {
	// Source format
	Format string // "uri", "qr", "manual", "backup_codes"

	// Data (URI, QR image path, etc.)
	Data string

	// Import options
	Options map[string]interface{}

	// Preview only
	Preview bool
}

// TOTPSync represents TOTP synchronization between devices
type TOTPSync struct {
	// Sync ID
	ID string

	// Device ID
	DeviceID string

	// Sync timestamp
	Timestamp time.Time

	// Entries to sync
	Entries []*TOTPEntry

	// Conflict resolution
	ConflictResolution string // "local", "remote", "merge"
}

// TOTPFilter represents filtering options for TOTP entries
type TOTPFilter struct {
	// Search term
	Search string

	// Issuer filter
	Issuer string

	// Account filter
	Account string

	// Tags filter
	Tags []string

	// Folder path
	Folder string

	// Favorites only
	FavoritesOnly bool

	// Include backup codes
	IncludeBackupCodes bool

	// Created after
	CreatedAfter *time.Time

	// Updated after
	UpdatedAfter *time.Time

	// Limit results
	Limit int

	// Offset for pagination
	Offset int

	// Sort by field
	SortBy string

	// Sort direction
	SortDirection string // "asc", "desc"
}

// TOTPOptions represents TOTP generation options
type TOTPOptions struct {
	// Algorithm
	Algorithm TOTPAlgorithm

	// Digits
	Digits int

	// Period
	Period int

	// Account name
	Account string

	// Issuer name
	Issuer string

	// Generate QR code
	GenerateQR bool

	// Include backup codes
	GenerateBackupCodes bool

	// Number of backup codes
	BackupCodeCount int
}

// TOTPUtility provides TOTP utility functions
type TOTPUtility struct{}

// NewTOTPUtility creates a new TOTP utility
func NewTOTPUtility() *TOTPUtility {
	return &TOTPUtility{}
}

// GenerateSecret generates a new random secret
func (u *TOTPUtility) GenerateSecret(length int) (string, error) {
	// TODO: Implement secure random secret generation
	// This should generate cryptographically secure random bytes
	// and encode them as base32
	return "", fmt.Errorf("generate secret not yet implemented")
}

// ParseURI parses a TOTP URI
func (u *TOTPUtility) ParseURI(uri string) (*TOTPEntry, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid URI: %w", err)
	}

	if parsed.Scheme != "otpauth" {
		return nil, fmt.Errorf("invalid scheme: %s", parsed.Scheme)
	}

	if parsed.Host != "totp" {
		return nil, fmt.Errorf("invalid type: %s", parsed.Host)
	}

	entry := &TOTPEntry{
		Algorithm: TOTPAlgorithmSHA1,
		Digits:    6,
		Period:    30,
	}

	// Parse label parts
	pathParts := strings.Split(strings.TrimPrefix(parsed.Path, "/"), ":")
	if len(pathParts) >= 2 {
		entry.Issuer = pathParts[0]
		entry.Account = pathParts[1]
	} else if len(pathParts) == 1 {
		entry.Account = pathParts[0]
	}

	// Parse query parameters
	query := parsed.Query()

	if secret := query.Get("secret"); secret != "" {
		entry.Secret = secret
	}

	if algorithm := query.Get("algorithm"); algorithm != "" {
		entry.Algorithm = TOTPAlgorithm(strings.ToUpper(algorithm))
	}

	if digits := query.Get("digits"); digits != "" {
		if d, err := strconv.Atoi(digits); err == nil {
			entry.Digits = d
		}
	}

	if period := query.Get("period"); period != "" {
		if p, err := strconv.Atoi(period); err == nil {
			entry.Period = p
		}
	}

	if issuer := query.Get("issuer"); issuer != "" {
		entry.Issuer = issuer
	}

	return entry, nil
}

// GenerateURI generates a TOTP URI
func (u *TOTPUtility) GenerateURI(entry *TOTPEntry) (string, error) {
	if entry.Secret == "" {
		return "", fmt.Errorf("secret is required")
	}

	var label string
	if entry.Issuer != "" && entry.Account != "" {
		label = fmt.Sprintf("%s:%s", entry.Issuer, entry.Account)
	} else if entry.Account != "" {
		label = entry.Account
	} else {
		return "", fmt.Errorf("account name is required")
	}

	values := url.Values{}
	values.Set("secret", entry.Secret)
	values.Set("algorithm", string(entry.Algorithm))
	values.Set("digits", strconv.Itoa(entry.Digits))
	values.Set("period", strconv.Itoa(entry.Period))

	if entry.Issuer != "" {
		values.Set("issuer", entry.Issuer)
	}

	uri := fmt.Sprintf("otpauth://totp/%s?%s", url.PathEscape(label), values.Encode())
	return uri, nil
}

// ValidateSecret validates a base32 secret
func (u *TOTPUtility) ValidateSecret(secret string) error {
	_, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return fmt.Errorf("invalid base32 secret: %w", err)
	}
	return nil
}

// CalculateTimeSteps calculates time steps for TOTP
func (u *TOTPUtility) CalculateTimeSteps(t time.Time, period int) int64 {
	return t.Unix() / int64(period)
}

// GenerateCode generates a TOTP code
func (u *TOTPUtility) GenerateCode(creds *TOTPCredentials) (string, error) {
	// Validate secret
	if creds.Secret == "" {
		return "", fmt.Errorf("secret is required")
	}

	// Decode secret
	secret, err := base32.StdEncoding.DecodeString(strings.ToUpper(creds.Secret))
	if err != nil {
		return "", fmt.Errorf("invalid secret: %w", err)
	}

	// Get time
	var t time.Time
	if creds.Time != nil {
		t = *creds.Time
	} else {
		t = time.Now()
	}

	// Calculate time steps
	period := creds.Period
	if period == 0 {
		period = 30
	}

	timeSteps := u.CalculateTimeSteps(t, period)

	// Generate hash
	var h func() hash.Hash
	switch creds.Algorithm {
	case TOTPAlgorithmSHA256:
		h = sha256.New
	case TOTPAlgorithmSHA512:
		h = sha512.New
	default:
		h = sha1.New
	}

	// Generate HMAC
	hmac := hmac.New(h, secret)
	binary.Write(hmac, binary.BigEndian, timeSteps)
	hash := hmac.Sum(nil)

	// Dynamic truncation
	offset := int(hash[len(hash)-1] & 0x0F)
	code := int(hash[offset]&0x7F)<<24 | int(hash[offset+1]&0xFF)<<16 | int(hash[offset+2]&0xFF)<<8 | int(hash[offset+3]&0xFF)

	// Format digits
	digits := creds.Digits
	if digits == 0 {
		digits = 6
	}

	modulo := int(math.Pow10(digits))
	code = code % modulo

	return fmt.Sprintf(fmt.Sprintf("%%0%dd", digits), code), nil
}

// VerifyCode verifies a TOTP code
func (u *TOTPUtility) VerifyCode(creds *TOTPCredentials, code string) (bool, error) {
	// Allow for clock skew by checking current, previous, and next time steps
	var t time.Time
	if creds.Time != nil {
		t = *creds.Time
	} else {
		t = time.Now()
	}

	period := creds.Period
	if period == 0 {
		period = 30
	}

	currentSteps := u.CalculateTimeSteps(t, period)

	// Check current, previous, and next time steps
	for i := -1; i <= 1; i++ {
		testCreds := *creds
		testCreds.Counter = currentSteps + int64(i)

		generatedCode, err := u.GenerateCode(&testCreds)
		if err != nil {
			return false, err
		}

		if generatedCode == code {
			return true, nil
		}
	}

	return false, nil
}

// GetTimeRemaining gets the time remaining until the next code
func (u *TOTPUtility) GetTimeRemaining(period int) int {
	if period == 0 {
		period = 30
	}

	return period - (int(time.Now().Unix()) % period)
}

// GenerateBackupCodes generates backup codes
func (u *TOTPUtility) GenerateBackupCodes(count int) ([]*TOTPBackupCode, error) {
	codes := make([]*TOTPBackupCode, count)

	for i := 0; i < count; i++ {
		// TODO: Generate secure random backup codes
		codes[i] = &TOTPBackupCode{
			Code:      fmt.Sprintf("BACKUP-%08d", i+1),
			Used:      false,
			CreatedAt: time.Now(),
		}
	}

	return codes, nil
}
