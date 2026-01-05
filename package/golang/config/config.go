package config

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Config struct {
	Endpoint   string            `json:"endpoint"`
	Token      string            `json:"token,omitempty"`
	Timeout    time.Duration     `json:"timeout"`
	RetryCount int               `json:"retry_count"`
	UserAgent  string            `json:"user_agent"`
	TLSConfig  *TLSConfig        `json:"tls_config,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Debug      bool              `json:"debug"`
}

type TLSConfig struct {
	Enabled            bool   `json:"enabled"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
	CertFile           string `json:"cert_file,omitempty"`
	KeyFile            string `json:"key_file,omitempty"`
	CAFile             string `json:"ca_file,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		Timeout:    30 * time.Second,
		RetryCount: 3,
		UserAgent:  "aether-vault-go/1.0.0",
		Headers:    make(map[string]string),
		Debug:      false,
	}
}

func NewConfig(endpoint, token string) *Config {
	cfg := DefaultConfig()
	cfg.Endpoint = endpoint
	cfg.Token = token
	return cfg
}

func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	return c
}

func (c *Config) WithRetryCount(count int) *Config {
	c.RetryCount = count
	return c
}

func (c *Config) WithUserAgent(ua string) *Config {
	c.UserAgent = ua
	return c
}

func (c *Config) WithTLS(tlsConfig *TLSConfig) *Config {
	c.TLSConfig = tlsConfig
	return c
}

func (c *Config) WithHeader(key, value string) *Config {
	if c.Headers == nil {
		c.Headers = make(map[string]string)
	}
	c.Headers[key] = value
	return c
}

func (c *Config) WithDebug(debug bool) *Config {
	c.Debug = debug
	return c
}

func (c *Config) Validate() error {
	if c.Endpoint == "" {
		return NewValidationError("endpoint is required")
	}
	if c.Timeout <= 0 {
		return NewValidationError("timeout must be positive")
	}
	if c.RetryCount < 0 {
		return NewValidationError("retry count cannot be negative")
	}
	return nil
}

func (t *TLSConfig) ToTLSConfig() (*tls.Config, error) {
	if !t.Enabled {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: t.InsecureSkipVerify,
	}

	if t.CertFile != "" && t.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
		if err != nil {
			return nil, NewValidationError("failed to load TLS certificate: %v", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}

func (c *Config) ToHTTPClient() (*http.Client, error) {
	client := &http.Client{
		Timeout: c.Timeout,
	}

	if c.TLSConfig != nil {
		tlsConfig, err := c.TLSConfig.ToTLSConfig()
		if err != nil {
			return nil, err
		}
		if tlsConfig != nil {
			client.Transport = &http.Transport{
				TLSClientConfig: tlsConfig,
			}
		}
	}

	return client, nil
}
