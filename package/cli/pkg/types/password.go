package types

import (
	"time"
)

// PasswordType represents different types of password entries
type PasswordType string

const (
	// PasswordTypeLogin represents login credentials
	PasswordTypeLogin PasswordType = "login"

	// PasswordTypeCard represents credit card information
	PasswordTypeCard PasswordType = "card"

	// PasswordTypeIdentity represents identity information
	PasswordTypeIdentity PasswordType = "identity"

	// PasswordTypeSecureNote represents secure notes
	PasswordTypeSecureNote PasswordType = "secure_note"

	// PasswordTypeSSH represents SSH keys
	PasswordTypeSSH PasswordType = "ssh"

	// PasswordTypeDatabase represents database credentials
	PasswordTypeDatabase PasswordType = "database"
)

// PasswordEntry represents a password entry
type PasswordEntry struct {
	// Unique identifier
	ID string

	// Entry type
	Type PasswordType

	// Entry name
	Name string

	// Username
	Username string

	// Password
	Password string

	// URL/Website
	URL string

	// Notes
	Notes string

	// Custom fields
	Fields map[string]interface{}

	// Metadata
	Metadata *PasswordMetadata

	// Creation timestamp
	CreatedAt time.Time

	// Last update timestamp
	UpdatedAt time.Time

	// Expires at
	ExpiresAt *time.Time

	// Tags
	Tags []string

	// Folder/Path
	Folder string

	// Is favorite
	Favorite bool
}

// PasswordMetadata contains metadata about a password entry
type PasswordMetadata struct {
	// Created by
	CreatedBy string

	// Last modified by
	UpdatedBy string

	// Password strength
	PasswordStrength *PasswordStrength

	// Last used timestamp
	LastUsed *time.Time

	// Usage count
	UsageCount int

	// Auto-fill information
	Autofill *AutofillInfo

	// Attachment information
	Attachments []string

	// Version information
	Version int

	// Whether password has been reused
	PasswordReused bool

	// Whether password is compromised
	PasswordCompromised bool
}

// PasswordStrength represents password strength analysis
type PasswordStrength struct {
	// Strength score (0-100)
	Score int

	// Strength level
	Level string // "weak", "fair", "good", "strong"

	// Estimated crack time
	CrackTime string

	// Issues found
	Issues []string

	// Suggestions for improvement
	Suggestions []string
}

// AutofillInfo contains auto-fill information
type AutofillInfo struct {
	// CSS selectors
	UsernameSelector string
	PasswordSelector string

	// Field IDs
	UsernameField string
	PasswordField string

	// Form URL
	FormURL string
}

// CardInfo represents credit card information
type CardInfo struct {
	// Cardholder name
	CardholderName string

	// Card number (encrypted)
	Number string

	// Expiry month
	ExpiryMonth string

	// Expiry year
	ExpiryYear string

	// CVV (encrypted)
	CVV string

	// Card brand
	Brand string // "visa", "mastercard", "amex", etc.

	// Card type
	Type string // "credit", "debit", "prepaid"
}

// IdentityInfo represents personal identity information
type IdentityInfo struct {
	// Personal details
	Title      string
	FirstName  string
	LastName   string
	MiddleName string

	// Contact information
	Email      string
	Phone      string
	Address    string
	City       string
	State      string
	PostalCode string
	Country    string

	// Identification
	SSN            string
	PassportNumber string
	LicenseNumber  string

	// Professional details
	Company    string
	JobTitle   string
	Department string
}

// SSHKeyInfo represents SSH key information
type SSHKeyInfo struct {
	// Key type
	Type string // "rsa", "ed25519", "ecdsa"

	// Public key
	PublicKey string

	// Private key (encrypted)
	PrivateKey string

	// Key fingerprint
	Fingerprint string

	// Key comment
	Comment string

	// Passphrase
	Passphrase string

	// Key size
	KeySize int

	// Algorithm
	Algorithm string
}

// DatabaseInfo represents database credential information
type DatabaseInfo struct {
	// Database type
	Type string // "mysql", "postgresql", "mongodb", etc.

	// Host
	Host string

	// Port
	Port int

	// Database name
	Database string

	// Username
	Username string

	// Password
	Password string

	// Connection string
	ConnectionString string

	// SSL mode
	SSLMode string

	// Additional parameters
	Parameters map[string]string
}

// PasswordFilter represents filtering options for password entries
type PasswordFilter struct {
	// Search term
	Search string

	// Entry type
	Type PasswordType

	// Folder path
	Folder string

	// Tags filter
	Tags []string

	// Favorites only
	FavoritesOnly bool

	// Include archived
	IncludeArchived bool

	// Created after
	CreatedAfter *time.Time

	// Updated after
	UpdatedAfter *time.Time

	// Expires before
	ExpiresBefore *time.Time

	// Limit results
	Limit int

	// Offset for pagination
	Offset int

	// Sort by field
	SortBy string

	// Sort direction
	SortDirection string // "asc", "desc"
}

// PasswordGenerationOptions represents options for password generation
type PasswordGenerationOptions struct {
	// Password length
	Length int

	// Include uppercase letters
	Uppercase bool

	// Include lowercase letters
	Lowercase bool

	// Include numbers
	Numbers bool

	// Include symbols
	Symbols bool

	// Include ambiguous characters
	Ambiguous bool

	// Minimum number of uppercase letters
	MinUppercase int

	// Minimum number of lowercase letters
	MinLowercase int

	// Minimum number of numbers
	MinNumbers int

	// Minimum number of symbols
	MinSymbols int

	// Exclude similar characters
	ExcludeSimilar bool

	// Custom characters to include
	CustomCharacters string

	// Words for passphrase
	Words []string

	// Word separator for passphrase
	WordSeparator string

	// Capitalize words in passphrase
	CapitalizeWords bool

	// Include number in passphrase
	IncludeNumber bool

	// Use passphrase instead of random password
	Passphrase bool
}

// PasswordAnalysis represents password analysis results
type PasswordAnalysis struct {
	// Original password
	Password string

	// Strength analysis
	Strength *PasswordStrength

	// Is it compromised
	Compromised bool

	// Similar passwords found
	SimilarPasswords []string

	// Reused passwords
	ReusedPasswords []string

	// Time to crack
	TimeToCrack string

	// Entropy
	Entropy float64

	// Character usage
	CharacterUsage map[string]int
}

// PasswordExport represents password export format
type PasswordExport struct {
	// Entries
	Entries []*PasswordEntry

	// Export format
	Format string // "csv", "json", "1pif", "xml"

	// Include sensitive data
	IncludePasswords bool

	// Encryption key (if any)
	EncryptionKey string

	// Export timestamp
	ExportedAt time.Time

	// Export version
	Version string
}

// PasswordImport represents password import format
type PasswordImport struct {
	// Source format
	Format string // "csv", "json", "1pif", "xml", "keepass", "lastpass"

	// File path
	FilePath string

	// Import options
	Options map[string]interface{}

	// Preview only (don't actually import)
	Preview bool

	// Merge or replace strategy
	Strategy string // "merge", "replace"
}
