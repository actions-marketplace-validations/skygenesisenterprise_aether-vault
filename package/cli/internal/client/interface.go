package client

import (
	"context"
	"fmt"

	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// Client interface defines the Vault client contract
type Client interface {
	// Authentication
	Login(ctx context.Context, credentials *types.Credentials) (*types.AuthResponse, error)
	Logout(ctx context.Context) error
	RefreshToken(ctx context.Context, refreshToken string) (*types.TokenResponse, error)
	IsAuthenticated() bool
	SetToken(token string)

	// Secret management
	GetSecret(ctx context.Context, path string) (*types.Secret, error)
	SetSecret(ctx context.Context, path string, secret *types.Secret) error
	DeleteSecret(ctx context.Context, path string) error
	ListSecrets(ctx context.Context, prefix string) ([]*types.SecretMetadata, error)

	// Path management
	GetPath(ctx context.Context, path string) (*types.PathInfo, error)
	CreatePath(ctx context.Context, path string, metadata *types.PathMetadata) error
	DeletePath(ctx context.Context, path string) error

	// Sync operations
	Sync(ctx context.Context, direction types.SyncDirection) (*types.SyncResult, error)
	GetSyncStatus(ctx context.Context) (*types.SyncStatus, error)

	// Configuration
	GetConfig() *types.ClientConfig
	SetConfig(config *types.ClientConfig) error

	// Health and status
	Health(ctx context.Context) (*types.HealthStatus, error)
	Status(ctx context.Context) (*types.ClientStatus, error)
}

// ClientOptions contains options for creating a client
type ClientOptions struct {
	Config  *types.ClientConfig
	Timeout int
	Headers map[string]string
}

// NewClient creates a new Vault client
func NewClient(options *ClientOptions) (Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	// Default configuration
	if options.Config == nil {
		options.Config = &types.ClientConfig{
			Type: "local",
		}
	}

	// Create client based on type
	switch options.Config.Type {
	case "local":
		return NewLocalClient(options)
	case "cloud":
		return NewCloudClient(options)
	default:
		return nil, fmt.Errorf("unsupported client type: %s", options.Config.Type)
	}
}

// ValidateClient validates that a client is properly configured
func ValidateClient(client Client) error {
	if client == nil {
		return fmt.Errorf("client cannot be nil")
	}

	config := client.GetConfig()
	if config == nil {
		return fmt.Errorf("client configuration cannot be nil")
	}

	if config.Type == "" {
		return fmt.Errorf("client type cannot be empty")
	}

	return nil
}
