package context

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// LocalContext provides local execution context
type LocalContext struct {
	*Context
}

// NewLocalContext creates a new local execution context
func NewLocalContext(cfg *types.Config) (*LocalContext, error) {
	baseCtx, err := New(cfg)
	if err != nil {
		return nil, err
	}

	// Set mode to local
	baseCtx.SetMode(types.LocalMode)

	localCtx := &LocalContext{
		Context: baseCtx,
	}

	// Initialize local environment
	if err := localCtx.initializeLocalEnvironment(); err != nil {
		return nil, err
	}

	return localCtx, nil
}

// initializeLocalEnvironment sets up the local environment
func (lc *LocalContext) initializeLocalEnvironment() error {
	// Ensure local directories exist
	if lc.Config != nil {
		vaultPath := lc.Config.Local.Path
		if vaultPath != "" {
			if err := os.MkdirAll(vaultPath, 0755); err != nil {
				return err
			}

			// Create keys directory
			keysPath := lc.Config.Local.KeyFile
			if keysPath != "" {
				if err := os.MkdirAll(keysPath[:len(keysPath)-len("/vault.key")], 0700); err != nil {
					return err
				}
			}
		}
	}

	// Set environment-specific runtime info
	lc.updateRuntimeInfo()

	return nil
}

// updateRuntimeInfo updates runtime information for local context
func (lc *LocalContext) updateRuntimeInfo() {
	if lc.Runtime == nil {
		return
	}

	// Add local-specific environment variables
	if lc.Runtime.Env == nil {
		lc.Runtime.Env = make(map[string]string)
	}

	lc.Runtime.Env["VAULT_MODE"] = "local"
	lc.Runtime.Env["VAULT_PATH"] = lc.getVaultPath()
	lc.Runtime.Env["OS"] = runtime.GOOS
	lc.Runtime.Env["ARCH"] = runtime.GOARCH
}

// getVaultPath returns the vault path
func (lc *LocalContext) getVaultPath() string {
	if lc.Config != nil && lc.Config.Local.Path != "" {
		return lc.Config.Local.Path
	}

	// Default path
	home, err := os.UserHomeDir()
	if err != nil {
		return "./.aether-vault"
	}
	return filepath.Join(home, ".aether", "vault")
}

// IsOffline returns true if running in offline mode
func (lc *LocalContext) IsOffline() bool {
	return true
}

// GetStoragePath returns the local storage path
func (lc *LocalContext) GetStoragePath() string {
	return lc.getVaultPath()
}

// ValidateLocalEnvironment validates the local environment
func (lc *LocalContext) ValidateLocalEnvironment() error {
	vaultPath := lc.getVaultPath()

	// Check if vault path exists and is accessible
	if _, err := os.Stat(vaultPath); os.IsNotExist(err) {
		return err
	}

	// Check write permissions
	testFile := filepath.Join(vaultPath, ".write_test")
	if file, err := os.Create(testFile); err != nil {
		return err
	} else {
		file.Close()
		os.Remove(testFile)
	}

	return nil
}

// GetLocalFeatures returns available local features
func (lc *LocalContext) GetLocalFeatures() []string {
	return []string{
		"local_storage",
		"encryption",
		"offline_mode",
		"file_based_secrets",
		"local_sync",
	}
}
