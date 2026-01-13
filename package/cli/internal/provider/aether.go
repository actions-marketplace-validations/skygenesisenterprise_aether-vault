package provider

import (
	"context"
	"fmt"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/client"
	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// AetherVaultProvider implements VaultProvider for Aether Vault
type AetherVaultProvider struct {
	config *types.VaultProviderConfig
	client client.Client
}

// NewAetherVaultProvider creates a new Aether Vault provider
func NewAetherVaultProvider(config *types.VaultProviderConfig) (VaultProvider, error) {
	provider := &AetherVaultProvider{
		config: config,
	}

	return provider, nil
}

// GetType returns the provider type
func (p *AetherVaultProvider) GetType() types.VaultType {
	return types.VaultTypeAether
}

// GetName returns the provider name
func (p *AetherVaultProvider) GetName() string {
	if p.config.Name != "" {
		return p.config.Name
	}
	return "Aether Vault"
}

// GetMetadata returns provider metadata
func (p *AetherVaultProvider) GetMetadata() *types.VaultProviderMetadata {
	return &types.VaultProviderMetadata{
		Type:                 types.VaultTypeAether,
		DisplayName:          "Aether Vault",
		Description:          "Next-generation secure vault management system",
		SupportedAuthMethods: []string{"oauth", "token", "password"},
		DefaultURL:           "https://cloud.aethervault.com",
		RequiredFields:       []string{"server_url"},
		OptionalFields:       []string{"organization_id", "project_id"},
	}
}

// GetCapabilities returns provider capabilities
func (p *AetherVaultProvider) GetCapabilities() *types.ProviderCapabilities {
	return &types.ProviderCapabilities{
		SecretManagement: &types.SecretCapabilities{
			SecretTypes: []string{"generic", "ssh", "database", "certificate"},
			Versioning:  true,
			TTL:         true,
			Sharing:     true,
			Encryption:  true,
		},
		Auth: &types.AuthCapabilities{
			Methods:      []types.AuthMethod{types.AuthMethodOAuth, types.AuthMethodToken, types.AuthMethodPassword},
			MFA:          true,
			SSO:          true,
			TokenRefresh: true,
		},
		Sync: &types.SyncCapabilities{
			RealTime:           true,
			Bidirectional:      true,
			ConflictResolution: true,
			SelectiveSync:      true,
		},
		Features: []string{"audit", "compliance", "multi-tenant"},
	}
}

// ValidateConfig validates provider configuration
func (p *AetherVaultProvider) ValidateConfig(config *types.VaultProviderConfig) error {
	if config.ServerURL == "" {
		return fmt.Errorf("server_url is required")
	}

	if config.AuthMethod == "" {
		return fmt.Errorf("auth_method is required")
	}

	return nil
}

// Initialize initializes the provider
func (p *AetherVaultProvider) Initialize(config *types.VaultProviderConfig) error {
	if err := p.ValidateConfig(config); err != nil {
		return err
	}

	p.config = config

	// Create client configuration
	clientConfig := &types.ClientConfig{
		Type:       "cloud",
		ServerURL:  config.ServerURL,
		AuthMethod: config.AuthMethod,
		Timeout:    30,
		Options:    config.Config,
	}

	// Create client (using existing client system)
	clientOptions := &client.ClientOptions{
		Config: clientConfig,
	}

	var err error
	p.client, err = client.NewClient(clientOptions)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	return nil
}

// Connect connects to the vault
func (p *AetherVaultProvider) Connect(ctx context.Context, credentials *types.Credentials) error {
	if p.client == nil {
		return fmt.Errorf("client not initialized")
	}

	_, err := p.client.Login(ctx, credentials)
	return err
}

// Disconnect disconnects from the vault
func (p *AetherVaultProvider) Disconnect() error {
	if p.client == nil {
		return nil
	}

	ctx := context.Background()
	return p.client.Logout(ctx)
}

// IsConnected returns connection status
func (p *AetherVaultProvider) IsConnected() bool {
	if p.client == nil {
		return false
	}

	ctx := context.Background()
	status, err := p.client.Status(ctx)
	return err == nil && status.Connected
}

// IsAuthenticated returns authentication status
func (p *AetherVaultProvider) IsAuthenticated() bool {
	if p.client == nil {
		return false
	}

	return p.client.IsAuthenticated()
}

// Authenticate authenticates with the vault
func (p *AetherVaultProvider) Authenticate(ctx context.Context, method types.AuthMethod, credentials interface{}) error {
	if p.client == nil {
		return fmt.Errorf("client not initialized")
	}

	// Convert credentials to expected type
	var creds *types.Credentials
	switch c := credentials.(type) {
	case *types.Credentials:
		creds = c
	default:
		return fmt.Errorf("invalid credentials type")
	}

	_, err := p.client.Login(ctx, creds)
	return err
}

// RefreshToken refreshes the authentication token
func (p *AetherVaultProvider) RefreshToken(ctx context.Context) error {
	if p.client == nil {
		return fmt.Errorf("client not initialized")
	}

	// TODO: Get refresh token from storage
	refreshToken := ""
	_, err := p.client.RefreshToken(ctx, refreshToken)
	return err
}

// Logout logs out from the vault
func (p *AetherVaultProvider) Logout() error {
	return p.Disconnect()
}

// GetSecret retrieves a secret
func (p *AetherVaultProvider) GetSecret(ctx context.Context, path string) (*types.Secret, error) {
	if p.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	return p.client.GetSecret(ctx, path)
}

// SetSecret stores a secret
func (p *AetherVaultProvider) SetSecret(ctx context.Context, path string, secret *types.Secret) error {
	if p.client == nil {
		return fmt.Errorf("client not initialized")
	}

	return p.client.SetSecret(ctx, path, secret)
}

// DeleteSecret removes a secret
func (p *AetherVaultProvider) DeleteSecret(ctx context.Context, path string) error {
	if p.client == nil {
		return fmt.Errorf("client not initialized")
	}

	return p.client.DeleteSecret(ctx, path)
}

// ListSecrets lists secrets at a path
func (p *AetherVaultProvider) ListSecrets(ctx context.Context, prefix string) ([]*types.SecretMetadata, error) {
	if p.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	return p.client.ListSecrets(ctx, prefix)
}

// GetPathInfo retrieves path information
func (p *AetherVaultProvider) GetPathInfo(ctx context.Context, path string) (*types.PathInfo, error) {
	if p.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	return p.client.GetPath(ctx, path)
}

// CreatePath creates a path
func (p *AetherVaultProvider) CreatePath(ctx context.Context, path string, metadata *types.PathMetadata) error {
	if p.client == nil {
		return fmt.Errorf("client not initialized")
	}

	return p.client.CreatePath(ctx, path, metadata)
}

// DeletePath removes a path
func (p *AetherVaultProvider) DeletePath(ctx context.Context, path string) error {
	if p.client == nil {
		return fmt.Errorf("client not initialized")
	}

	return p.client.DeletePath(ctx, path)
}

// Sync performs sync operations
func (p *AetherVaultProvider) Sync(ctx context.Context, direction types.SyncDirection) (*types.SyncResult, error) {
	if p.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	return p.client.Sync(ctx, direction)
}

// GetSyncStatus gets sync status
func (p *AetherVaultProvider) GetSyncStatus(ctx context.Context) (*types.SyncStatus, error) {
	if p.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	return p.client.GetSyncStatus(ctx)
}

// Health checks the health of the vault
func (p *AetherVaultProvider) Health(ctx context.Context) (*types.HealthStatus, error) {
	if p.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	return p.client.Health(ctx)
}

// Status gets the client status
func (p *AetherVaultProvider) Status(ctx context.Context) (*types.ClientStatus, error) {
	if p.client == nil {
		return &types.ClientStatus{
			Connected:     false,
			Authenticated: false,
		}, nil
	}

	return p.client.Status(ctx)
}

// GetConfig returns the current configuration
func (p *AetherVaultProvider) GetConfig() *types.VaultProviderConfig {
	return p.config
}

// UpdateConfig updates the configuration
func (p *AetherVaultProvider) UpdateConfig(config *types.VaultProviderConfig) error {
	if err := p.ValidateConfig(config); err != nil {
		return err
	}

	p.config = config
	return nil
}

// Close closes the provider
func (p *AetherVaultProvider) Close() error {
	// Client interface doesn't have Close method, so we just set to nil
	p.client = nil
	return nil
}
