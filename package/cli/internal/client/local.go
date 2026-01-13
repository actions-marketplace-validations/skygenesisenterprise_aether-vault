package client

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// LocalClient implements the Client interface for local storage
type LocalClient struct {
	config    *types.ClientConfig
	basePath  string
	masterKey []byte
	gcm       cipher.AEAD
}

// NewLocalClient creates a new local client
func NewLocalClient(options *ClientOptions) (Client, error) {
	if options == nil || options.Config == nil {
		return nil, fmt.Errorf("client options and config are required")
	}

	basePath := options.Config.Options["basePath"].(string)
	if basePath == "" {
		// Default to ~/.aether/vault
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		basePath = filepath.Join(home, ".aether", "vault")
	}

	// Ensure base path exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base path: %w", err)
	}

	// Generate or load master key
	masterKey, err := getOrCreateMasterKey(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to setup encryption: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	return &LocalClient{
		config:    options.Config,
		basePath:  basePath,
		masterKey: masterKey,
		gcm:       gcm,
	}, nil
}

// getOrCreateMasterKey loads or creates a master encryption key
func getOrCreateMasterKey(basePath string) ([]byte, error) {
	keyFile := filepath.Join(basePath, ".master_key")

	// Try to load existing key
	if keyData, err := os.ReadFile(keyFile); err == nil {
		// Decode base64 key
		return base64.StdEncoding.DecodeString(string(keyData))
	}

	// Generate new master key
	masterKey := make([]byte, 32)
	if _, err := rand.Read(masterKey); err != nil {
		return nil, fmt.Errorf("failed to generate master key: %w", err)
	}

	// Save key (base64 encoded)
	encodedKey := base64.StdEncoding.EncodeToString(masterKey)
	if err := os.WriteFile(keyFile, []byte(encodedKey), 0600); err != nil {
		return nil, fmt.Errorf("failed to save master key: %w", err)
	}

	return masterKey, nil
}

// encrypt encrypts data using AES-GCM
func (c *LocalClient) encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt data
	ciphertext := c.gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts data using AES-GCM
func (c *LocalClient) decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := c.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return c.gcm.Open(nil, nonce, ciphertext, nil)
}

// getSecretPath returns the file path for a secret
func (c *LocalClient) getSecretPath(path string) string {
	// Sanitize path and create hash for filename
	safePath := filepath.Clean(path)
	hash := sha256.Sum256([]byte(safePath))
	filename := base64.StdEncoding.EncodeToString(hash[:]) + ".enc"
	return filepath.Join(c.basePath, "secrets", filename)
}

// ensureSecretsDir ensures the secrets directory exists
func (c *LocalClient) ensureSecretsDir() error {
	secretsDir := filepath.Join(c.basePath, "secrets")
	return os.MkdirAll(secretsDir, 0755)
}

// Login authenticates with local storage (always succeeds for local mode)
func (c *LocalClient) Login(ctx context.Context, credentials *types.Credentials) (*types.AuthResponse, error) {
	// Local mode doesn't require authentication
	return &types.AuthResponse{
		AccessToken:  "local-mode",
		RefreshToken: "local-mode",
		TokenType:    "Bearer",
		ExpiresIn:    0, // Never expires
		Scope:        "all",
		User: &types.UserInfo{
			ID:          "local-user",
			Username:    "local",
			Email:       "",
			DisplayName: "Local User",
		},
	}, nil
}

// Logout logs out (no-op for local mode)
func (c *LocalClient) Logout(ctx context.Context) error {
	return nil
}

// RefreshToken refreshes the token (no-op for local mode)
func (c *LocalClient) RefreshToken(ctx context.Context, refreshToken string) (*types.TokenResponse, error) {
	return &types.TokenResponse{
		AccessToken:  "local-mode",
		RefreshToken: "local-mode",
		TokenType:    "Bearer",
		ExpiresIn:    0,
	}, nil
}

// IsAuthenticated returns true for local mode
func (c *LocalClient) IsAuthenticated() bool {
	return true
}

// GetSecret retrieves a secret from local storage
func (c *LocalClient) GetSecret(ctx context.Context, path string) (*types.Secret, error) {
	if err := c.ensureSecretsDir(); err != nil {
		return nil, fmt.Errorf("failed to ensure secrets directory: %w", err)
	}

	secretPath := c.getSecretPath(path)

	// Check if secret exists
	if _, err := os.Stat(secretPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("secret not found: %s", path)
	}

	// Read encrypted data
	encryptedData, err := os.ReadFile(secretPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret file: %w", err)
	}

	// Decrypt data
	data, err := c.decrypt(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt secret: %w", err)
	}

	// Unmarshal secret
	var secret types.Secret
	if err := json.Unmarshal(data, &secret); err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret: %w", err)
	}

	return &secret, nil
}

// SetSecret stores a secret in local storage
func (c *LocalClient) SetSecret(ctx context.Context, path string, secret *types.Secret) error {
	if err := c.ensureSecretsDir(); err != nil {
		return fmt.Errorf("failed to ensure secrets directory: %w", err)
	}

	// Set path and update metadata
	secret.Path = path
	now := time.Now().Unix()

	if secret.Metadata == nil {
		secret.Metadata = &types.SecretMetadata{
			CreatedAt: now,
			CreatedBy: "local",
		}
	}

	secret.Metadata.Path = path
	secret.Metadata.UpdatedAt = now
	secret.Metadata.UpdatedBy = "local"

	// Marshal secret
	data, err := json.Marshal(secret)
	if err != nil {
		return fmt.Errorf("failed to marshal secret: %w", err)
	}

	// Encrypt data
	encryptedData, err := c.encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %w", err)
	}

	// Write to file
	secretPath := c.getSecretPath(path)
	if err := os.WriteFile(secretPath, encryptedData, 0600); err != nil {
		return fmt.Errorf("failed to write secret file: %w", err)
	}

	return nil
}

// DeleteSecret removes a secret from local storage
func (c *LocalClient) DeleteSecret(ctx context.Context, path string) error {
	secretPath := c.getSecretPath(path)

	// Check if secret exists
	if _, err := os.Stat(secretPath); os.IsNotExist(err) {
		return fmt.Errorf("secret not found: %s", path)
	}

	// Delete file
	if err := os.Remove(secretPath); err != nil {
		return fmt.Errorf("failed to delete secret file: %w", err)
	}

	return nil
}

// ListSecrets lists secrets in local storage
func (c *LocalClient) ListSecrets(ctx context.Context, prefix string) ([]*types.SecretMetadata, error) {
	if err := c.ensureSecretsDir(); err != nil {
		return nil, fmt.Errorf("failed to ensure secrets directory: %w", err)
	}

	secretsDir := filepath.Join(c.basePath, "secrets")

	// Read all files in secrets directory
	files, err := os.ReadDir(secretsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets directory: %w", err)
	}

	var secrets []*types.SecretMetadata

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".enc" {
			// Read and decrypt secret
			filePath := filepath.Join(secretsDir, file.Name())
			encryptedData, err := os.ReadFile(filePath)
			if err != nil {
				continue // Skip files that can't be read
			}

			data, err := c.decrypt(encryptedData)
			if err != nil {
				continue // Skip files that can't be decrypted
			}

			var secret types.Secret
			if err := json.Unmarshal(data, &secret); err != nil {
				continue // Skip invalid files
			}

			// Filter by prefix if specified
			if prefix == "" || len(secret.Path) >= len(prefix) && secret.Path[:len(prefix)] == prefix {
				secrets = append(secrets, secret.Metadata)
			}
		}
	}

	return secrets, nil
}

// GetPath gets path information
func (c *LocalClient) GetPath(ctx context.Context, path string) (*types.PathInfo, error) {
	// TODO: Implement actual path retrieval
	return &types.PathInfo{
		Path:      path,
		Type:      "folder",
		Metadata:  &types.PathMetadata{},
		Children:  []*types.PathInfo{},
		CreatedAt: 0,
		UpdatedAt: 0,
	}, nil
}

// CreatePath creates a new path
func (c *LocalClient) CreatePath(ctx context.Context, path string, metadata *types.PathMetadata) error {
	// TODO: Implement actual path creation
	return nil
}

// DeletePath removes a path
func (c *LocalClient) DeletePath(ctx context.Context, path string) error {
	// TODO: Implement actual path deletion
	return nil
}

// Sync performs sync operations (no-op for local mode)
func (c *LocalClient) Sync(ctx context.Context, direction types.SyncDirection) (*types.SyncResult, error) {
	return &types.SyncResult{
		Direction:    direction,
		ItemsSynced:  0,
		ItemsCreated: 0,
		ItemsUpdated: 0,
		ItemsDeleted: 0,
		Conflicts:    0,
		Errors:       []string{},
		StartTime:    0,
		EndTime:      0,
		Duration:     0,
	}, nil
}

// GetSyncStatus gets sync status (no-op for local mode)
func (c *LocalClient) GetSyncStatus(ctx context.Context) (*types.SyncStatus, error) {
	return &types.SyncStatus{
		LastSync:         nil,
		Direction:        types.SyncLocalToCloud,
		InProgress:       false,
		Progress:         0,
		CurrentOperation: "",
		PendingChanges:   0,
		Conflicts:        0,
		LastResult:       nil,
	}, nil
}

// GetConfig returns the client configuration
func (c *LocalClient) GetConfig() *types.ClientConfig {
	return c.config
}

// SetConfig updates the client configuration
func (c *LocalClient) SetConfig(config *types.ClientConfig) error {
	c.config = config
	return nil
}

// Health checks the health of the local storage
func (c *LocalClient) Health(ctx context.Context) (*types.HealthStatus, error) {
	return &types.HealthStatus{
		Status:  "healthy",
		Version: "1.0.0",
		Uptime:  0,
		System:  &types.ClientSystemInfo{},
		Checks:  []types.HealthCheck{},
	}, nil
}

// Status returns the client status
func (c *LocalClient) Status(ctx context.Context) (*types.ClientStatus, error) {
	return &types.ClientStatus{
		Connected:     true,
		Mode:          "local",
		ServerURL:     "",
		LastSync:      nil,
		LocalPath:     c.basePath,
		Authenticated: true,
	}, nil
}
