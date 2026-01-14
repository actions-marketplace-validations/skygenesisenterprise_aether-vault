package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/server/types"
)

// LocalClient provides local file-based secret storage
type LocalClient struct {
	config   *types.ClientConfig
	basePath string
	secrets  map[string]*types.Secret
	mutex    sync.RWMutex
}

// NewLocalClient creates a new local client
func NewLocalClient(options *types.ClientOptions) (*LocalClient, error) {
	basePath, ok := options.Config.Options["basePath"].(string)
	if !ok {
		basePath = "./data"
	}

	client := &LocalClient{
		config:   options.Config,
		basePath: basePath,
		secrets:  make(map[string]*types.Secret),
	}

	// Load existing secrets
	if err := client.loadSecrets(); err != nil {
		return nil, fmt.Errorf("failed to load secrets: %w", err)
	}

	return client, nil
}

// GetSecret retrieves a secret
func (c *LocalClient) GetSecret(path string) (*types.Secret, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	secret, exists := c.secrets[path]
	if !exists {
		return nil, fmt.Errorf("secret not found: %s", path)
	}

	return secret, nil
}

// SetSecret stores a secret
func (c *LocalClient) SetSecret(path string, secret *types.Secret) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	secret.Path = path
	if secret.Metadata == nil {
		secret.Metadata = &types.SecretMetadata{}
	}

	now := time.Now().Unix()
	if secret.Metadata.CreatedAt == 0 {
		secret.Metadata.CreatedAt = now
		secret.Metadata.CreatedBy = "vault-server"
	}
	secret.Metadata.UpdatedAt = now
	secret.Metadata.UpdatedBy = "vault-server"

	c.secrets[path] = secret

	return c.saveSecrets()
}

// DeleteSecret removes a secret
func (c *LocalClient) DeleteSecret(path string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.secrets[path]; !exists {
		return fmt.Errorf("secret not found: %s", path)
	}

	delete(c.secrets, path)

	return c.saveSecrets()
}

// ListSecrets lists secrets at a path
func (c *LocalClient) ListSecrets(prefix string) ([]string, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var keys []string
	for path := range c.secrets {
		if prefix == "" || strings.HasPrefix(path, prefix) {
			// Remove prefix and get the next part
			relativePath := strings.TrimPrefix(path, prefix)
			relativePath = strings.TrimPrefix(relativePath, "/")

			if relativePath != "" {
				parts := strings.Split(relativePath, "/")
				if len(parts) > 0 {
					keys = append(keys, parts[0])
				}
			}
		}
	}

	// Remove duplicates
	uniqueKeys := make(map[string]bool)
	for _, key := range keys {
		uniqueKeys[key] = true
	}

	var result []string
	for key := range uniqueKeys {
		result = append(result, key)
	}

	return result, nil
}

// GetStatus returns client status
func (c *LocalClient) GetStatus() (*types.ClientStatus, error) {
	return &types.ClientStatus{
		Connected:     true,
		Mode:          types.ExecutionModeLocal,
		LocalPath:     c.basePath,
		Authenticated: true,
	}, nil
}

// Authenticate performs authentication
func (c *LocalClient) Authenticate(method string, credentials interface{}) error {
	// Local client doesn't require authentication
	return nil
}

// Close closes the client
func (c *LocalClient) Close() error {
	return c.saveSecrets()
}

// loadSecrets loads secrets from disk
func (c *LocalClient) loadSecrets() error {
	if err := os.MkdirAll(c.basePath, 0755); err != nil {
		return fmt.Errorf("failed to create base path: %w", err)
	}

	secretsFile := filepath.Join(c.basePath, "secrets.json")
	if _, err := os.Stat(secretsFile); os.IsNotExist(err) {
		return nil // No existing secrets
	}

	data, err := os.ReadFile(secretsFile)
	if err != nil {
		return fmt.Errorf("failed to read secrets file: %w", err)
	}

	if err := json.Unmarshal(data, &c.secrets); err != nil {
		return fmt.Errorf("failed to unmarshal secrets: %w", err)
	}

	return nil
}

// saveSecrets saves secrets to disk
func (c *LocalClient) saveSecrets() error {
	if err := os.MkdirAll(c.basePath, 0755); err != nil {
		return fmt.Errorf("failed to create base path: %w", err)
	}

	secretsFile := filepath.Join(c.basePath, "secrets.json")
	data, err := json.MarshalIndent(c.secrets, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal secrets: %w", err)
	}

	if err := os.WriteFile(secretsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write secrets file: %w", err)
	}

	return nil
}
