package types

// Secret represents a secret for the Vault server
type Secret struct {
	// Secret path
	Path string

	// Secret data
	Data map[string]interface{}

	// Metadata
	Metadata *SecretMetadata

	// Version
	Version int64
}

// SecretMetadata contains secret metadata
type SecretMetadata struct {
	// Secret path
	Path string

	// Creation timestamp
	CreatedAt int64

	// Last modification timestamp
	UpdatedAt int64

	// Created by
	CreatedBy string

	// Last modified by
	UpdatedBy string

	// Tags
	Tags []string

	// TTL
	TTL int64
}

// ClientStatus represents client status
type ClientStatus struct {
	// Is connected
	Connected bool

	// Connection mode
	Mode ExecutionMode

	// Server URL (cloud mode)
	ServerURL string

	// Last sync timestamp
	LastSync *int64

	// Local storage path (local mode)
	LocalPath string

	// Authentication status
	Authenticated bool
}

// ExecutionMode represents execution mode
type ExecutionMode string

const (
	// ExecutionModeLocal represents local execution
	ExecutionModeLocal ExecutionMode = "local"
	// ExecutionModeCloud represents cloud execution
	ExecutionModeCloud ExecutionMode = "cloud"
)

// ClientConfig represents client configuration
type ClientConfig struct {
	// Client type (local, cloud)
	Type string

	// Server URL (cloud mode)
	ServerURL string

	// Authentication method
	AuthMethod string

	// Credentials
	Credentials *Credentials

	// Timeout in seconds
	Timeout int

	// Custom options
	Options map[string]interface{}
}

// Credentials represents authentication credentials
type Credentials struct {
	// Authentication method
	Method string

	// Username/password for basic auth
	Username string
	Password string

	// Token for token auth
	Token string

	// OAuth configuration
	OAuth *OAuthCredentials
}

// OAuthCredentials contains OAuth-specific credentials
type OAuthCredentials struct {
	// Client ID
	ClientID string

	// Client secret
	ClientSecret string

	// Authorization code
	AuthCode string

	// Redirect URL
	RedirectURL string

	// Scopes
	Scopes []string
}

// ClientOptions represents client options
type ClientOptions struct {
	// Configuration
	Config *ClientConfig
}
