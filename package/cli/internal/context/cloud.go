package context

import (
	"context"
	"fmt"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// CloudContext provides cloud execution context
type CloudContext struct {
	*Context
	ctx     context.Context
	cancel  context.CancelFunc
	timeout time.Duration
}

// NewCloudContext creates a new cloud execution context
func NewCloudContext(cfg *types.Config) (*CloudContext, error) {
	baseCtx, err := New(cfg)
	if err != nil {
		return nil, err
	}

	// Set mode to cloud
	baseCtx.SetMode(types.CloudMode)

	// Create context with timeout
	timeout := 30 * time.Second
	if cfg != nil && cfg.General.Timeout > 0 {
		timeout = cfg.General.Timeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	cloudCtx := &CloudContext{
		Context: baseCtx,
		ctx:     ctx,
		cancel:  cancel,
		timeout: timeout,
	}

	// Initialize cloud environment
	if err := cloudCtx.initializeCloudEnvironment(); err != nil {
		cloudCtx.cancel()
		return nil, err
	}

	return cloudCtx, nil
}

// initializeCloudEnvironment sets up the cloud environment
func (cc *CloudContext) initializeCloudEnvironment() error {
	// Validate configuration
	if cc.Config == nil {
		return fmt.Errorf("configuration is required for cloud mode")
	}

	// Validate cloud configuration
	if err := cc.validateCloudConfig(); err != nil {
		return err
	}

	// Update runtime information
	cc.updateRuntimeInfo()

	return nil
}

// validateCloudConfig validates the cloud configuration
func (cc *CloudContext) validateCloudConfig() error {
	cloud := cc.Config.Cloud

	if cloud.URL == "" {
		return fmt.Errorf("cloud URL is required")
	}

	if cloud.AuthMethod == "" {
		return fmt.Errorf("cloud authentication method is required")
	}

	if cloud.AuthMethod == "oauth" && cloud.OAuth.ClientID == "" {
		return fmt.Errorf("OAuth client ID is required for OAuth authentication")
	}

	return nil
}

// updateRuntimeInfo updates runtime information for cloud context
func (cc *CloudContext) updateRuntimeInfo() {
	if cc.Runtime == nil {
		return
	}

	// Add cloud-specific environment variables
	if cc.Runtime.Env == nil {
		cc.Runtime.Env = make(map[string]string)
	}

	cc.Runtime.Env["VAULT_MODE"] = "cloud"
	cc.Runtime.Env["VAULT_URL"] = cc.Config.Cloud.URL
	cc.Runtime.Env["VAULT_AUTH_METHOD"] = cc.Config.Cloud.AuthMethod
}

// GetContext returns the Go context
func (cc *CloudContext) GetContext() context.Context {
	return cc.ctx
}

// GetTimeout returns the timeout duration
func (cc *CloudContext) GetTimeout() time.Duration {
	return cc.timeout
}

// IsOnline returns true if running in online mode
func (cc *CloudContext) IsOnline() bool {
	return true
}

// GetServerURL returns the configured server URL
func (cc *CloudContext) GetServerURL() string {
	if cc.Config != nil {
		return cc.Config.Cloud.URL
	}
	return ""
}

// GetAuthMethod returns the authentication method
func (cc *CloudContext) GetAuthMethod() string {
	if cc.Config != nil {
		return cc.Config.Cloud.AuthMethod
	}
	return ""
}

// ValidateCloudEnvironment validates the cloud environment
func (cc *CloudContext) ValidateCloudEnvironment() error {
	// TODO: Implement actual connectivity check
	// For now, just validate configuration
	return cc.validateCloudConfig()
}

// GetCloudFeatures returns available cloud features
func (cc *CloudContext) GetCloudFeatures() []string {
	return []string{
		"cloud_storage",
		"oauth_authentication",
		"real_time_sync",
		"collaboration",
		"audit_logging",
		"enterprise_features",
	}
}

// Close closes the cloud context
func (cc *CloudContext) Close() error {
	if cc.cancel != nil {
		cc.cancel()
	}
	return nil
}

// ExtendTimeout extends the context timeout
func (cc *CloudContext) ExtendTimeout(timeout time.Duration) {
	cc.timeout = timeout
	cc.ctx, cc.cancel = context.WithTimeout(context.Background(), timeout)
}

// IsExpired checks if the context is expired
func (cc *CloudContext) IsExpired() bool {
	select {
	case <-cc.ctx.Done():
		return true
	default:
		return false
	}
}

// WithValue returns a new context with the given value
func (cc *CloudContext) WithValue(key, value interface{}) *CloudContext {
	newCtx := &CloudContext{
		Context: cc.Context,
		ctx:     context.WithValue(cc.ctx, key, value),
		cancel:  cc.cancel,
		timeout: cc.timeout,
	}
	return newCtx
}
