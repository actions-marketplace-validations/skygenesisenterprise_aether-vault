package vault

import (
	"context"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/golang/audit"
	"github.com/skygenesisenterprise/aether-vault/package/golang/auth"
	"github.com/skygenesisenterprise/aether-vault/package/golang/client"
	"github.com/skygenesisenterprise/aether-vault/package/golang/config"
	"github.com/skygenesisenterprise/aether-vault/package/golang/errors"
	"github.com/skygenesisenterprise/aether-vault/package/golang/identity"
	"github.com/skygenesisenterprise/aether-vault/package/golang/policies"
	"github.com/skygenesisenterprise/aether-vault/package/golang/secrets"
	"github.com/skygenesisenterprise/aether-vault/package/golang/totp"
)

type Vault struct {
	config *config.Config
	client *client.Client

	Auth     *auth.AuthClient
	Secrets  *secrets.SecretsClient
	TOTP     *totp.TOTPClient
	Identity *identity.IdentityClient
	Policies *policies.PolicyClient
	Audit    *audit.AuditClient
}

type Config struct {
	Endpoint   string            `json:"endpoint"`
	Token      string            `json:"token,omitempty"`
	Timeout    time.Duration     `json:"timeout"`
	RetryCount int               `json:"retry_count"`
	UserAgent  string            `json:"user_agent"`
	TLSConfig  *config.TLSConfig `json:"tls_config,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Debug      bool              `json:"debug"`
}

func New(cfg *Config) (*Vault, error) {
	vaultConfig := &config.Config{
		Endpoint:   cfg.Endpoint,
		Token:      cfg.Token,
		Timeout:    cfg.Timeout,
		RetryCount: cfg.RetryCount,
		UserAgent:  cfg.UserAgent,
		TLSConfig:  cfg.TLSConfig,
		Headers:    cfg.Headers,
		Debug:      cfg.Debug,
	}

	if vaultConfig.Timeout == 0 {
		vaultConfig.Timeout = 30 * time.Second
	}
	if vaultConfig.RetryCount == 0 {
		vaultConfig.RetryCount = 3
	}
	if vaultConfig.UserAgent == "" {
		vaultConfig.UserAgent = "aether-vault-go/1.0.0"
	}
	if vaultConfig.Headers == nil {
		vaultConfig.Headers = make(map[string]string)
	}

	client, err := client.NewClient(vaultConfig)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to create vault client")
	}

	vault := &Vault{
		config: vaultConfig,
		client: client,

		Auth:     auth.NewAuthClient(client),
		Secrets:  secrets.NewSecretsClient(client),
		TOTP:     totp.NewTOTPClient(client),
		Identity: identity.NewIdentityClient(client),
		Policies: policies.NewPolicyClient(client),
		Audit:    audit.NewAuditClient(client),
	}

	return vault, nil
}

func DefaultConfig() *Config {
	return &Config{
		Timeout:    30 * time.Second,
		RetryCount: 3,
		UserAgent:  "aether-vault-go/1.0.0",
		Debug:      false,
	}
}

func NewConfig(endpoint, token string) *Config {
	cfg := DefaultConfig()
	cfg.Endpoint = endpoint
	cfg.Token = token
	return cfg
}

func (v *Vault) GetConfig() *config.Config {
	return v.config
}

func (v *Vault) GetClient() *client.Client {
	return v.client
}

func (v *Vault) Close() error {
	return nil
}

func (v *Vault) Health(ctx context.Context) error {
	resp, err := v.client.Get(ctx, "/api/v1/health")
	if err != nil {
		return err
	}

	var healthResp struct {
		Status string `json:"status"`
	}
	if err := resp.Decode(&healthResp); err != nil {
		return errors.WrapError(err, errors.ErrCodeInternal, "failed to decode health response")
	}

	if healthResp.Status != "healthy" {
		return errors.NewError(errors.ErrCodeUnavailable, "vault is not healthy")
	}

	return nil
}

func (v *Vault) Version(ctx context.Context) (string, error) {
	resp, err := v.client.Get(ctx, "/api/v1/version")
	if err != nil {
		return "", err
	}

	var versionResp struct {
		Version string `json:"version"`
	}
	if err := resp.Decode(&versionResp); err != nil {
		return "", errors.WrapError(err, errors.ErrCodeInternal, "failed to decode version response")
	}

	return versionResp.Version, nil
}

func (v *Vault) Info(ctx context.Context) (map[string]interface{}, error) {
	resp, err := v.client.Get(ctx, "/api/v1/info")
	if err != nil {
		return nil, err
	}

	var info map[string]interface{}
	if err := resp.Decode(&info); err != nil {
		return nil, errors.WrapError(err, errors.ErrCodeInternal, "failed to decode info response")
	}

	return info, nil
}
