package types

// Client represents the Vault client interface
type Client interface {
	// GetSecret retrieves a secret
	GetSecret(path string) (*Secret, error)

	// SetSecret stores a secret
	SetSecret(path string, secret *Secret) error

	// DeleteSecret removes a secret
	DeleteSecret(path string) error

	// ListSecrets lists secrets at a path
	ListSecrets(path string) ([]string, error)

	// GetStatus returns client status
	GetStatus() (*ClientStatus, error)

	// Authenticate performs authentication
	Authenticate(method string, credentials interface{}) error

	// Close closes the client
	Close() error
}

// Secret represents a secret
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

// SecretFilter represents secret filtering options
type SecretFilter struct {
	// Path pattern
	Path string

	// Tags filter
	Tags []string

	// Created after
	CreatedAfter *int64

	// Updated after
	UpdatedAfter *int64

	// Limit results
	Limit int

	// Offset for pagination
	Offset int
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

// AuthResponse represents authentication response
type AuthResponse struct {
	// Access token
	AccessToken string

	// Refresh token
	RefreshToken string

	// Token type
	TokenType string

	// Expires in seconds
	ExpiresIn int64

	// Scope
	Scope string

	// User information
	User *UserInfo
}

// TokenResponse represents token refresh response
type TokenResponse struct {
	// Access token
	AccessToken string

	// Refresh token
	RefreshToken string

	// Token type
	TokenType string

	// Expires in seconds
	ExpiresIn int64
}

// SyncDirection represents sync direction
type SyncDirection string

const (
	// SyncLocalToCloud syncs from local to cloud
	SyncLocalToCloud SyncDirection = "local-to-cloud"

	// SyncCloudToLocal syncs from cloud to local
	SyncCloudToLocal SyncDirection = "cloud-to-local"

	// SyncBidirectional syncs both ways
	SyncBidirectional SyncDirection = "bidirectional"
)

// SyncResult represents sync operation result
type SyncResult struct {
	// Sync direction
	Direction SyncDirection

	// Number of items synced
	ItemsSynced int

	// Number of items created
	ItemsCreated int

	// Number of items updated
	ItemsUpdated int

	// Number of items deleted
	ItemsDeleted int

	// Number of conflicts
	Conflicts int

	// Errors encountered
	Errors []string

	// Start time
	StartTime int64

	// End time
	EndTime int64

	// Duration in milliseconds
	Duration int64
}

// SyncStatus represents sync status
type SyncStatus struct {
	// Last sync time
	LastSync *int64

	// Sync direction
	Direction SyncDirection

	// Is sync in progress
	InProgress bool

	// Sync progress percentage
	Progress int

	// Current operation
	CurrentOperation string

	// Pending changes
	PendingChanges int

	// Conflicts
	Conflicts int

	// Last sync result
	LastResult *SyncResult
}

// PathInfo represents path information
type PathInfo struct {
	// Path
	Path string

	// Path type
	Type string

	// Metadata
	Metadata *PathMetadata

	// Children
	Children []*PathInfo

	// Permissions
	Permissions *PathPermissions

	// Created at
	CreatedAt int64

	// Updated at
	UpdatedAt int64
}

// PathMetadata represents path metadata
type PathMetadata struct {
	// Description
	Description string

	// Tags
	Tags []string

	// TTL
	TTL int64

	// Max versions
	MaxVersions int

	// Auto-delete policy
	AutoDelete string

	// Created by
	CreatedBy string

	// Updated by
	UpdatedBy string
}

// PathPermissions represents path permissions
type PathPermissions struct {
	// Read permissions
	Read []string

	// Write permissions
	Write []string

	// Delete permissions
	Delete []string

	// Admin permissions
	Admin []string
}

// HealthStatus represents health status
type HealthStatus struct {
	// Overall health status
	Status string

	// Server version
	Version string

	// Uptime in seconds
	Uptime int64

	// System information
	System *ClientSystemInfo

	// Checks
	Checks []HealthCheck
}

// HealthCheck represents a health check
type HealthCheck struct {
	// Check name
	Name string

	// Status
	Status string

	// Message
	Message string

	// Duration in milliseconds
	Duration int64

	// Timestamp
	Timestamp int64
}

// ClientSystemInfo represents system information for health status
type ClientSystemInfo struct {
	// OS
	OS string

	// Architecture
	Arch string

	// Memory usage
	MemoryUsage *MemoryInfo

	// Disk usage
	DiskUsage *DiskInfo

	// CPU usage
	CPUUsage float64
}

// MemoryInfo represents memory information
type MemoryInfo struct {
	// Total memory in bytes
	Total int64

	// Used memory in bytes
	Used int64

	// Available memory in bytes
	Available int64

	// Usage percentage
	UsagePercentage float64
}

// DiskInfo represents disk information
type DiskInfo struct {
	// Total space in bytes
	Total int64

	// Used space in bytes
	Used int64

	// Available space in bytes
	Available int64

	// Usage percentage
	UsagePercentage float64

	// Mount point
	MountPoint string
}

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

	// Retry configuration
	Retry *RetryConfig

	// Proxy configuration
	Proxy *ProxyConfig

	// TLS configuration
	TLS *TLSConfig

	// Headers
	Headers map[string]string

	// Custom options
	Options map[string]interface{}
}

// RetryConfig represents retry configuration
type RetryConfig struct {
	// Maximum number of retries
	MaxRetries int

	// Initial retry delay in milliseconds
	InitialDelay int64

	// Maximum retry delay in milliseconds
	MaxDelay int64

	// Backoff multiplier
	Multiplier float64

	// Jitter
	Jitter bool
}

// ProxyConfig represents proxy configuration
type ProxyConfig struct {
	// Proxy URL
	URL string

	// Username
	Username string

	// Password
	Password string

	// No proxy hosts
	NoProxy []string
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	// Enable TLS verification
	Verify bool

	// CA certificate file
	CAFile string

	// Client certificate file
	CertFile string

	// Client key file
	KeyFile string

	// Server name
	ServerName string

	// Insecure skip verify
	InsecureSkipVerify bool
}
