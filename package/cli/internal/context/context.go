package context

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/client"
	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// Status represents the current status of the CLI
type Status struct {
	// Current execution mode
	Mode types.ExecutionMode `json:"mode"`

	// Configuration status
	ConfigStatus string `json:"config_status"`

	// Authentication status
	AuthStatus string `json:"auth_status"`

	// Runtime environment
	Runtime *types.RuntimeInfo `json:"runtime"`

	// Last updated timestamp
	LastUpdated time.Time `json:"last_updated"`

	// Additional status information
	Details map[string]interface{} `json:"details,omitempty"`
}

// Context represents the execution context
type Context struct {
	// Current execution mode
	Mode types.ExecutionMode

	// Configuration
	Config *types.Config

	// Runtime information
	Runtime *types.RuntimeInfo

	// Authentication state
	Auth *types.AuthState

	// Vault client
	Client client.Client
}

// New creates a new execution context
func New(cfg *types.Config) (*Context, error) {
	if cfg == nil {
		cfg = &types.Config{}
	}

	// Create runtime info
	runtimeInfo := &types.RuntimeInfo{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		GoVersion: runtime.Version(),
		Version:   "1.0.0", // TODO: Get from build info
		Env:       make(map[string]string),
		Build: &types.BuildInfo{
			Commit:       "unknown", // TODO: Get from build info
			Timestamp:    time.Now().Format(time.RFC3339),
			Environment:  "development", // TODO: Get from build info
			ToolsVersion: "unknown",     // TODO: Get from build info
		},
	}

	// Try to load cloud credentials first
	oauthClient := client.NewOAuthClient("https://vault.skygenesisenterprise.com")
	var authClient client.Client
	var mode types.ExecutionMode
	var authState *types.AuthState

	if cloudAuth, err := oauthClient.LoadCredentials(); err == nil {
		// Cloud credentials available, use cloud mode
		fmt.Println("Using cloud authentication")
		authClient = client.NewHTTPClient("https://vault.skygenesisenterprise.com")
		authClient.SetToken(cloudAuth.AccessToken)
		mode = types.CloudMode
		authState = &types.AuthState{
			Authenticated: true,
			User:          cloudAuth.User,
		}
	} else {
		// No cloud credentials, use local server
		fmt.Println("No cloud credentials found, using local server")

		// Create service manager to check server status
		serviceManager := client.NewServiceManager()

		// Ensure server is running
		if !serviceManager.IsRunning() {
			fmt.Println("Vault server is not running, starting it...")
			if err := serviceManager.StartServer(true, ""); err != nil {
				return nil, fmt.Errorf("failed to start server: %w", err)
			}

			// Wait for server to be ready
			vaultCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := serviceManager.WaitUntilServerReady(vaultCtx, 10*time.Second); err != nil {
				return nil, fmt.Errorf("server did not become ready: %w", err)
			}
			fmt.Println("Vault server is now ready")
		}

		// Create HTTP client to connect to local server
		authClient = client.NewHTTPClient(serviceManager.GetServerURL())

		// Authenticate with dev token
		credentials := &types.Credentials{
			Method: "token",
			Token:  "dev-token",
		}

		if _, err := authClient.Login(nil, credentials); err != nil {
			return nil, fmt.Errorf("failed to authenticate with server: %w", err)
		}

		mode = types.LocalMode
		authState = &types.AuthState{Authenticated: true}
	}

	// Create context
	ctx := &Context{
		Mode:    mode,
		Config:  cfg,
		Runtime: runtimeInfo,
		Auth:    authState,
		Client:  authClient,
	}

	return ctx, nil
}

// GetStatus returns the current status
func (c *Context) GetStatus() (*Status, error) {
	status := &Status{
		Mode:         c.Mode,
		ConfigStatus: "loaded",
		AuthStatus:   "unauthenticated",
		Runtime:      c.Runtime,
		LastUpdated:  time.Now(),
		Details:      make(map[string]interface{}),
	}

	// Update auth status
	if c.Auth != nil && c.Auth.Authenticated {
		status.AuthStatus = "authenticated"
		if c.Auth.User != nil {
			status.Details["user"] = c.Auth.User.Username
		}
	}

	// Update config status
	if c.Config == nil {
		status.ConfigStatus = "not_loaded"
		status.Details["config_path"] = "default"
	} else {
		status.Details["config_path"] = c.Config.Local.Path
	}

	return status, nil
}

// SetMode sets the execution mode
func (c *Context) SetMode(mode types.ExecutionMode) {
	c.Mode = mode
}

// SetAuth sets the authentication state
func (c *Context) SetAuth(auth *types.AuthState) {
	c.Auth = auth
}

// IsAuthenticated returns true if authenticated
func (c *Context) IsAuthenticated() bool {
	return c.Auth != nil && c.Auth.Authenticated
}

// IsCloudMode returns true if in cloud mode
func (c *Context) IsCloudMode() bool {
	return c.Mode == types.CloudMode
}

// Read reads data from the specified path
func (c *Context) Read(path string) (map[string]interface{}, error) {
	if c.Client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	secret, err := c.Client.GetSecret(nil, path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}

	return map[string]interface{}{
		"data":     secret.Data,
		"metadata": secret.Metadata,
		"path":     path,
	}, nil
}

// Write writes data to the specified path
func (c *Context) Write(path string, data map[string]interface{}, force bool) (map[string]interface{}, error) {
	if c.Client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	secret := &types.Secret{
		Path: path,
		Data: data,
		Metadata: &types.SecretMetadata{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
			CreatedBy: "cli",
			UpdatedBy: "cli",
		},
		Version: 1,
	}

	if err := c.Client.SetSecret(nil, path, secret); err != nil {
		return nil, fmt.Errorf("failed to write secret: %w", err)
	}

	return map[string]interface{}{
		"path":    path,
		"version": secret.Version,
		"written": true,
	}, nil
}

// Delete deletes data from the specified path
func (c *Context) Delete(path string, versions string, recursive bool) (map[string]interface{}, error) {
	if c.Client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	if err := c.Client.DeleteSecret(nil, path); err != nil {
		return nil, fmt.Errorf("failed to delete secret: %w", err)
	}

	return map[string]interface{}{
		"path":      path,
		"deleted":   true,
		"versions":  versions,
		"recursive": recursive,
	}, nil
}

// List lists data from the specified path
func (c *Context) List(path string) (map[string]interface{}, error) {
	if c.Client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	secrets, err := c.Client.ListSecrets(nil, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	var keys []string
	for _, secret := range secrets {
		keys = append(keys, secret.Path)
	}

	return map[string]interface{}{
		"path": path,
		"keys": keys,
		"summary": map[string]interface{}{
			"total_keys": len(keys),
			"folders":    len(keys),
			"files":      0,
		},
	}, nil
}

// Unwrap unwraps a wrapped secret
func (c *Context) Unwrap(token string) (map[string]interface{}, error) {
	// TODO: Implement actual Vault unwrap logic
	// For now, return mock data
	return map[string]interface{}{
		"token": token,
		"data": map[string]interface{}{
			"secret": "unwrapped-data",
			"ttl":    "1h",
		},
	}, nil
}
