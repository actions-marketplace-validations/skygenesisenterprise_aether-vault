package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// DefaultConfigPath returns the default configuration file path
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./config.yaml"
	}
	return filepath.Join(home, ".aether", "vault", "config.yaml")
}

// DefaultVaultPath returns the default vault directory path
func DefaultVaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./.aether-vault"
	}
	return filepath.Join(home, ".aether", "vault")
}

// DefaultKeyPath returns the default encryption key file path
func DefaultKeyPath() string {
	return filepath.Join(DefaultVaultPath(), "keys", "vault.key")
}

// CreateDefaultConfig creates a default configuration
func CreateDefaultConfig() *types.Config {
	return &types.Config{
		General: types.GeneralConfig{
			DefaultFormat: "table",
			Verbose:       false,
			Timeout:       30 * time.Second,
		},
		Local: types.LocalConfig{
			Path:            DefaultVaultPath(),
			KeyFile:         DefaultKeyPath(),
			AutoLockTimeout: 10 * time.Minute,
		},
		Cloud: types.CloudConfig{
			URL:        "https://vault.skygenesisenterprise.com",
			AuthMethod: "oauth",
			OAuth: types.OAuthConfig{
				ClientID:    "vault-cli",
				Scopes:      []string{"vault:read", "vault:write"},
				RedirectURL: "http://localhost:8080/callback",
			},
		},
		UI: types.UIConfig{
			Color:      true,
			Spinner:    true,
			TableStyle: "default",
		},
	}
}

// EnsureDefaultDirectories creates the default directory structure
func EnsureDefaultDirectories() error {
	vaultPath := DefaultVaultPath()
	keyPath := filepath.Dir(DefaultKeyPath())

	// Create vault directory
	if err := os.MkdirAll(vaultPath, 0755); err != nil {
		return err
	}

	// Create keys directory
	if err := os.MkdirAll(keyPath, 0700); err != nil {
		return err
	}

	return nil
}

// IsConfigured checks if the CLI is configured
func IsConfigured() bool {
	configPath := DefaultConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetConfigPaths returns all important configuration paths
func GetConfigPaths() map[string]string {
	return map[string]string{
		"config": DefaultConfigPath(),
		"vault":  DefaultVaultPath(),
		"key":    DefaultKeyPath(),
	}
}
