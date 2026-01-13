package provider

import (
	"context"
	"fmt"

	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// VaultProvider defines the interface for vault providers
type VaultProvider interface {
	// Provider information
	GetType() types.VaultType
	GetName() string
	GetMetadata() *types.VaultProviderMetadata
	GetCapabilities() *types.ProviderCapabilities

	// Configuration
	ValidateConfig(config *types.VaultProviderConfig) error
	Initialize(config *types.VaultProviderConfig) error

	// Connection management
	Connect(ctx context.Context, credentials *types.Credentials) error
	Disconnect() error
	IsConnected() bool
	IsAuthenticated() bool

	// Authentication
	Authenticate(ctx context.Context, method types.AuthMethod, credentials interface{}) error
	RefreshToken(ctx context.Context) error
	Logout() error

	// Secret operations
	GetSecret(ctx context.Context, path string) (*types.Secret, error)
	SetSecret(ctx context.Context, path string, secret *types.Secret) error
	DeleteSecret(ctx context.Context, path string) error
	ListSecrets(ctx context.Context, prefix string) ([]*types.SecretMetadata, error)

	// Path operations
	GetPathInfo(ctx context.Context, path string) (*types.PathInfo, error)
	CreatePath(ctx context.Context, path string, metadata *types.PathMetadata) error
	DeletePath(ctx context.Context, path string) error

	// Sync operations
	Sync(ctx context.Context, direction types.SyncDirection) (*types.SyncResult, error)
	GetSyncStatus(ctx context.Context) (*types.SyncStatus, error)

	// Health and status
	Health(ctx context.Context) (*types.HealthStatus, error)
	Status(ctx context.Context) (*types.ClientStatus, error)

	// Configuration management
	GetConfig() *types.VaultProviderConfig
	UpdateConfig(config *types.VaultProviderConfig) error

	// Cleanup
	Close() error
}

// ProviderFactory creates vault providers
type ProviderFactory interface {
	// Create a provider instance
	Create(config *types.VaultProviderConfig) (VaultProvider, error)

	// Get supported provider types
	GetSupportedTypes() []types.VaultType

	// Get provider metadata
	GetMetadata(vaultType types.VaultType) (*types.VaultProviderMetadata, error)
}

// Registry manages vault providers
type Registry struct {
	factory   ProviderFactory
	providers map[types.VaultType]VaultProvider
}

// NewRegistry creates a new provider registry
func NewRegistry(factory ProviderFactory) *Registry {
	return &Registry{
		factory:   factory,
		providers: make(map[types.VaultType]VaultProvider),
	}
}

// RegisterProvider registers a vault provider
func (r *Registry) RegisterProvider(config *types.VaultProviderConfig) error {
	provider, err := r.factory.Create(config)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}

	if err := provider.Initialize(config); err != nil {
		return fmt.Errorf("failed to initialize provider: %w", err)
	}

	r.providers[config.Type] = provider
	return nil
}

// GetProvider gets a vault provider by type
func (r *Registry) GetProvider(vaultType types.VaultType) (VaultProvider, error) {
	provider, exists := r.providers[vaultType]
	if !exists {
		return nil, fmt.Errorf("provider not found for type: %s", vaultType)
	}
	return provider, nil
}

// GetProviders returns all registered providers
func (r *Registry) GetProviders() map[types.VaultType]VaultProvider {
	return r.providers
}

// GetSupportedTypes returns supported provider types
func (r *Registry) GetSupportedTypes() []types.VaultType {
	return r.factory.GetSupportedTypes()
}

// GetProviderMetadata gets metadata for a provider type
func (r *Registry) GetProviderMetadata(vaultType types.VaultType) (*types.VaultProviderMetadata, error) {
	return r.factory.GetMetadata(vaultType)
}

// ConnectProvider connects to a specific provider
func (r *Registry) ConnectProvider(vaultType types.VaultType, credentials *types.Credentials) error {
	provider, err := r.GetProvider(vaultType)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return provider.Connect(ctx, credentials)
}

// DisconnectProvider disconnects from a specific provider
func (r *Registry) DisconnectProvider(vaultType types.VaultType) error {
	provider, err := r.GetProvider(vaultType)
	if err != nil {
		return err
	}

	return provider.Disconnect()
}

// Close closes all providers
func (r *Registry) Close() error {
	for _, provider := range r.providers {
		if err := provider.Close(); err != nil {
			// Log error but continue closing others
			fmt.Printf("Error closing provider: %v\n", err)
		}
	}
	return nil
}

// DefaultProviderFactory implements ProviderFactory
type DefaultProviderFactory struct {
	creators map[types.VaultType]ProviderCreator
}

// ProviderCreator creates a vault provider
type ProviderCreator func(config *types.VaultProviderConfig) (VaultProvider, error)

// NewDefaultProviderFactory creates a new default provider factory
func NewDefaultProviderFactory() *DefaultProviderFactory {
	return &DefaultProviderFactory{
		creators: make(map[types.VaultType]ProviderCreator),
	}
}

// RegisterCreator registers a provider creator
func (f *DefaultProviderFactory) RegisterCreator(vaultType types.VaultType, creator ProviderCreator) {
	f.creators[vaultType] = creator
}

// Create creates a provider instance
func (f *DefaultProviderFactory) Create(config *types.VaultProviderConfig) (VaultProvider, error) {
	creator, exists := f.creators[config.Type]
	if !exists {
		return nil, fmt.Errorf("no creator registered for provider type: %s", config.Type)
	}

	return creator(config)
}

// GetSupportedTypes returns supported provider types
func (f *DefaultProviderFactory) GetSupportedTypes() []types.VaultType {
	types := make([]types.VaultType, 0, len(f.creators))
	for vaultType := range f.creators {
		types = append(types, vaultType)
	}
	return types
}

// GetMetadata gets metadata for a provider type
func (f *DefaultProviderFactory) GetMetadata(vaultType types.VaultType) (*types.VaultProviderMetadata, error) {
	creator, exists := f.creators[vaultType]
	if !exists {
		return nil, fmt.Errorf("no creator registered for provider type: %s", vaultType)
	}

	// Create a temporary instance to get metadata
	provider, err := creator(&types.VaultProviderConfig{Type: vaultType})
	if err != nil {
		return nil, fmt.Errorf("failed to create provider for metadata: %w", err)
	}

	defer provider.Close()
	return provider.GetMetadata(), nil
}
