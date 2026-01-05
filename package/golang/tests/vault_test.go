package vault

import (
	"context"
	"testing"
	"time"

	"github.com/skygenesisenterprise/aether-vault"
)

func TestVaultClient(t *testing.T) {
	// Test configuration
	cfg := vault.NewConfig("http://localhost:8080", "test-token")
	cfg.Timeout = 5 * time.Second
	cfg.RetryCount = 2

	// Test client creation
	client, err := vault.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create vault client: %v", err)
	}

	// Test default config
	defaultCfg := vault.DefaultConfig()
	if defaultCfg.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout to be 30s, got %v", defaultCfg.Timeout)
	}

	// Test config validation
	invalidCfg := &vault.Config{
		Endpoint: "", // Invalid empty endpoint
	}
	_, err = vault.New(invalidCfg)
	if err == nil {
		t.Error("Expected error for invalid config")
	}

	// Test that all clients are initialized
	if client.Auth == nil {
		t.Error("Auth client should not be nil")
	}
	if client.Secrets == nil {
		t.Error("Secrets client should not be nil")
	}
	if client.TOTP == nil {
		t.Error("TOTP client should not be nil")
	}
	if client.Identity == nil {
		t.Error("Identity client should not be nil")
	}
	if client.Policies == nil {
		t.Error("Policies client should not be nil")
	}
	if client.Audit == nil {
		t.Error("Audit client should not be nil")
	}

	// Test that config and client are accessible
	retrievedConfig := client.GetConfig()
	if retrievedConfig == nil {
		t.Error("Config should not be nil")
	}

	retrievedClient := client.GetClient()
	if retrievedClient == nil {
		t.Error("Client should not be nil")
	}

	// Test close
	err = client.Close()
	if err != nil {
		t.Errorf("Close should not return error, got: %v", err)
	}
}

func TestVaultConfig(t *testing.T) {
	tests := []struct {
		name  string
		cfg   *vault.Config
		valid bool
	}{
		{
			name: "valid config",
			cfg: &vault.Config{
				Endpoint:   "http://localhost:8080",
				Token:      "test-token",
				Timeout:    30 * time.Second,
				RetryCount: 3,
			},
			valid: true,
		},
		{
			name: "invalid empty endpoint",
			cfg: &vault.Config{
				Endpoint: "",
				Token:    "test-token",
			},
			valid: false,
		},
		{
			name: "valid with TLS config",
			cfg: &vault.Config{
				Endpoint: "https://localhost:8080",
				Token:    "test-token",
				TLSConfig: &vault.TLSConfig{
					Enabled:            true,
					InsecureSkipVerify: false,
				},
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := vault.New(tt.cfg)

			if tt.valid && err != nil {
				t.Errorf("Expected valid config to succeed, got error: %v", err)
			}

			if !tt.valid && err == nil {
				t.Error("Expected invalid config to fail")
			}

			if tt.valid && client == nil {
				t.Error("Expected client to be created for valid config")
			}
		})
	}
}

func TestVaultMethods(t *testing.T) {
	// This test would require a running vault server
	// For now, we'll test that the methods exist and have correct signatures

	cfg := vault.NewConfig("http://localhost:8080", "test-token")
	client, err := vault.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test that methods exist (will fail with connection errors if no server)
	_, err = client.Health(ctx)
	if err == nil {
		t.Log("Health check passed (server might be running)")
	}

	_, err = client.Version(ctx)
	if err == nil {
		t.Log("Version check passed (server might be running)")
	}

	_, err = client.Info(ctx)
	if err == nil {
		t.Log("Info check passed (server might be running)")
	}

	t.Log("Method signature tests completed")
}
