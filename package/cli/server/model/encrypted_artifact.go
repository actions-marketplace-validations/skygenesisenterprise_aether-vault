package model

import (
	"time"

	"github.com/google/uuid"
)

// EncryptedArtifact represents an encrypted artifact
type EncryptedArtifact struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	FilePath     string    `json:"file_path"`
	OriginalPath string    `json:"original_path"`
	Algorithm    string    `json:"algorithm"`
	Version      string    `json:"version"`
	DataKeyHash  string    `json:"data_key_hash"`
	ContentSize  int64     `json:"content_size"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	AccessMethods []AccessMethod     `json:"access_methods"`
	Policies      []EncryptionPolicy `json:"policies"`
}

// AccessMethod represents an access method for decrypting an artifact
type AccessMethod struct {
	ID           uuid.UUID        `json:"id"`
	ArtifactID   uuid.UUID        `json:"artifact_id"`
	Type         AccessMethodType `json:"type"`
	Name         string           `json:"name"`
	Config       string           `json:"config"` // JSON config
	EncryptedKey string           `json:"encrypted_key"`
	KeyID        string           `json:"key_id"` // For certificate methods
	IsActive     bool             `json:"is_active"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

// AccessMethodType represents the type of access method
type AccessMethodType string

const (
	AccessMethodTypePassphrase  AccessMethodType = "passphrase"
	AccessMethodTypeCertificate AccessMethodType = "certificate"
	AccessMethodTypeRuntime     AccessMethodType = "runtime"
	AccessMethodTypePolicy      AccessMethodType = "policy"
)

// EncryptionPolicy represents a policy for encrypted artifacts
type EncryptionPolicy struct {
	ID         uuid.UUID  `json:"id"`
	ArtifactID uuid.UUID  `json:"artifact_id"`
	Name       string     `json:"name"`
	Type       PolicyType `json:"type"`
	Rules      string     `json:"rules"` // JSON rules
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// PolicyType represents the type of policy
type PolicyType string

const (
	PolicyTypeTTL         PolicyType = "ttl"
	PolicyTypeEnvironment PolicyType = "environment"
	PolicyTypeInstance    PolicyType = "instance"
	PolicyTypeRegion      PolicyType = "region"
	PolicyTypeMultiFactor PolicyType = "multi_factor"
)

// DecryptionAttempt represents a decryption attempt
type DecryptionAttempt struct {
	ID         uuid.UUID        `json:"id"`
	ArtifactID uuid.UUID        `json:"artifact_id"`
	UserID     uuid.UUID        `json:"user_id"`
	MethodType AccessMethodType `json:"method_type"`
	Success    bool             `json:"success"`
	Reason     string           `json:"reason"`
	IPAddress  string           `json:"ip_address"`
	UserAgent  string           `json:"user_agent"`
	CreatedAt  time.Time        `json:"created_at"`
}

// ArtifactMetadata represents metadata for encrypted artifacts
type ArtifactMetadata struct {
	ID            uuid.UUID `json:"id"`
	ArtifactID    uuid.UUID `json:"artifact_id"`
	OriginalName  string    `json:"original_name"`
	OriginalSize  int64     `json:"original_size"`
	OriginalMode  string    `json:"original_mode"`
	OriginalMtime string    `json:"original_mtime"`
	Checksum      string    `json:"checksum"`
	Compression   string    `json:"compression"`
	CreatedAt     time.Time `json:"created_at"`
}

// EncryptionRequest represents a request to encrypt data
type EncryptionRequest struct {
	SourcePath    string                   `json:"source_path"`
	OutputPath    string                   `json:"output_path"`
	AccessMethods []AccessMethodConfig     `json:"access_methods"`
	Policies      []EncryptionPolicyConfig `json:"policies"`
	Description   string                   `json:"description"`
	Compression   bool                     `json:"compression"`
}

// AccessMethodConfig represents configuration for an access method
type AccessMethodConfig struct {
	Type   AccessMethodType       `json:"type"`
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// EncryptionPolicyConfig represents configuration for an encryption policy
type EncryptionPolicyConfig struct {
	Type  PolicyType             `json:"type"`
	Name  string                 `json:"name"`
	Rules map[string]interface{} `json:"rules"`
}

// DecryptionRequest represents a request to decrypt data
type DecryptionRequest struct {
	ArtifactPath string             `json:"artifact_path"`
	OutputPath   string             `json:"output_path"`
	AccessMethod AccessMethodConfig `json:"access_method"`
	Force        bool               `json:"force"`
}

// EncryptionResult represents the result of an encryption operation
type EncryptionResult struct {
	ArtifactID    uuid.UUID `json:"artifact_id"`
	FilePath      string    `json:"file_path"`
	OriginalSize  int64     `json:"original_size"`
	EncryptedSize int64     `json:"encrypted_size"`
	Algorithm     string    `json:"algorithm"`
	AccessMethods []string  `json:"access_methods"`
	Policies      []string  `json:"policies"`
	CreatedAt     time.Time `json:"created_at"`
}

// DecryptionResult represents the result of a decryption operation
type DecryptionResult struct {
	ArtifactID    uuid.UUID `json:"artifact_id"`
	FilePath      string    `json:"file_path"`
	OriginalSize  int64     `json:"original_size"`
	DecryptedSize int64     `json:"decrypted_size"`
	MethodType    string    `json:"method_type"`
	Success       bool      `json:"success"`
	Reason        string    `json:"reason"`
	CreatedAt     time.Time `json:"created_at"`
}
