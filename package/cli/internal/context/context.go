package context

import (
	"runtime"
	"time"

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

	// Create context
	ctx := &Context{
		Mode:    types.LocalMode, // Default to local mode
		Config:  cfg,
		Runtime: runtimeInfo,
		Auth:    &types.AuthState{Authenticated: false},
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
