package types

// VaultType represents the type of vault provider
type VaultType string

const (
	// VaultTypeAether represents Aether Vault
	VaultTypeAether VaultType = "aether"

	// VaultTypeBitwarden represents Bitwarden
	VaultTypeBitwarden VaultType = "bitwarden"

	// VaultTypeVaultwarden represents Vaultwarden (Bitwarden compatible)
	VaultTypeVaultwarden VaultType = "vaultwarden"

	// VaultTypeHashicorp represents HashiCorp Vault
	VaultTypeHashicorp VaultType = "hashicorp"
)

// VaultProviderConfig represents configuration for a vault provider
type VaultProviderConfig struct {
	// Type of vault provider
	Type VaultType

	// Name for this vault instance (for identification)
	Name string

	// Server URL
	ServerURL string

	// Authentication method
	AuthMethod string

	// Provider-specific configuration
	Config map[string]interface{}

	// Is this the default vault
	Default bool

	// Priority for selection
	Priority int
}

// VaultProviderMetadata represents metadata about a vault provider
type VaultProviderMetadata struct {
	// Provider type
	Type VaultType

	// Display name
	DisplayName string

	// Description
	Description string

	// Supported authentication methods
	SupportedAuthMethods []string

	// Default server URL
	DefaultURL string

	// Required configuration fields
	RequiredFields []string

	// Optional configuration fields
	OptionalFields []string
}

// VaultConnection represents a connection to a vault
type VaultConnection struct {
	// Connection ID
	ID string

	// Provider configuration
	Config *VaultProviderConfig

	// Connection status
	Connected bool

	// Last connection timestamp
	LastConnected *int64

	// Authentication status
	Authenticated bool

	// User information
	User *UserInfo
}

// MultiVaultConfig represents configuration for multiple vaults
type MultiVaultConfig struct {
	// List of configured vaults
	Vaults []*VaultProviderConfig

	// Active vault ID
	ActiveVault string

	// Auto-switch settings
	AutoSwitch *AutoSwitchConfig

	// Sync settings between vaults
	Sync *VaultSyncConfig
}

// AutoSwitchConfig represents auto-switch configuration
type AutoSwitchConfig struct {
	// Enable auto-switching
	Enabled bool

	// Priority-based switching
	PriorityBased bool

	// Availability-based switching
	AvailabilityBased bool

	// Switch timeout in seconds
	Timeout int
}

// VaultSyncConfig represents sync configuration between vaults
type VaultSyncConfig struct {
	// Enable syncing
	Enabled bool

	// Default sync direction
	DefaultDirection SyncDirection

	// Sync intervals in minutes
	Interval int

	// Conflict resolution strategy
	ConflictResolution string

	// Excluded paths
	ExcludedPaths []string
}

// AuthMethod represents authentication methods supported by vault providers
type AuthMethod string

const (
	// AuthMethodOAuth represents OAuth authentication
	AuthMethodOAuth AuthMethod = "oauth"

	// AuthMethodToken represents token-based authentication
	AuthMethodToken AuthMethod = "token"

	// AuthMethodPassword represents username/password authentication
	AuthMethodPassword AuthMethod = "password"

	// AuthMethodAPIKey represents API key authentication
	AuthMethodAPIKey AuthMethod = "apikey"

	// AuthMethodCertificate represents certificate-based authentication
	AuthMethodCertificate AuthMethod = "certificate"
)

// ProviderCapabilities represents capabilities of a vault provider
type ProviderCapabilities struct {
	// Secret management capabilities
	SecretManagement *SecretCapabilities

	// Authentication capabilities
	Auth *AuthCapabilities

	// Sync capabilities
	Sync *SyncCapabilities

	// Additional features
	Features []string
}

// SecretCapabilities represents secret management capabilities
type SecretCapabilities struct {
	// Support for different secret types
	SecretTypes []string

	// Version support
	Versioning bool

	// TTL support
	TTL bool

	// Sharing capabilities
	Sharing bool

	// Encryption support
	Encryption bool
}

// AuthCapabilities represents authentication capabilities
type AuthCapabilities struct {
	// Supported auth methods
	Methods []AuthMethod

	// MFA support
	MFA bool

	// SSO support
	SSO bool

	// Token refresh
	TokenRefresh bool
}

// SyncCapabilities represents sync capabilities
type SyncCapabilities struct {
	// Real-time sync
	RealTime bool

	// Bidirectional sync
	Bidirectional bool

	// Conflict resolution
	ConflictResolution bool

	// Selective sync
	SelectiveSync bool
}
