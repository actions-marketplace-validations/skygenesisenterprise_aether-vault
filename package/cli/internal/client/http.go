package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// HTTPClient implements the Client interface for HTTP communication
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetToken sets the authentication token
func (c *HTTPClient) SetToken(token string) {
	c.token = token
}

// makeRequest makes an HTTP request with authentication
func (c *HTTPClient) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if c.token != "" {
		req.Header.Set("X-Vault-Token", c.token)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

// Login authenticates with the server
func (c *HTTPClient) Login(ctx context.Context, credentials *types.Credentials) (*types.AuthResponse, error) {
	loginData := map[string]interface{}{
		"role":  "default",
		"token": credentials.Token,
	}

	resp, err := c.makeRequest("POST", "/v1/auth/token/login", loginData)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var response struct {
		Auth struct {
			ClientToken string   `json:"client_token"`
			Policies    []string `json:"policies"`
		} `json:"auth"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode login response: %w", err)
	}

	// Set token for subsequent requests
	c.token = response.Auth.ClientToken

	return &types.AuthResponse{
		AccessToken:  response.Auth.ClientToken,
		RefreshToken: response.Auth.ClientToken,
		TokenType:    "Bearer",
		ExpiresIn:    0,
		Scope:        "all",
		User: &types.UserInfo{
			ID:          "local-user",
			Username:    "local",
			Email:       "",
			DisplayName: "Local User",
		},
	}, nil
}

// Logout logs out
func (c *HTTPClient) Logout(ctx context.Context) error {
	c.token = ""
	return nil
}

// RefreshToken refreshes the token
func (c *HTTPClient) RefreshToken(ctx context.Context, refreshToken string) (*types.TokenResponse, error) {
	return &types.TokenResponse{
		AccessToken:  c.token,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    0,
	}, nil
}

// IsAuthenticated returns true if token is set
func (c *HTTPClient) IsAuthenticated() bool {
	return c.token != ""
}

// GetSecret retrieves a secret
func (c *HTTPClient) GetSecret(ctx context.Context, path string) (*types.Secret, error) {
	resp, err := c.makeRequest("GET", "/v1/secret/"+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("secret not found: %s", path)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get secret failed with status: %d", resp.StatusCode)
	}

	var response struct {
		Data     map[string]interface{} `json:"data"`
		Metadata *types.SecretMetadata  `json:"metadata"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode secret response: %w", err)
	}

	return &types.Secret{
		Path:     path,
		Data:     response.Data,
		Metadata: response.Metadata,
		Version:  1,
	}, nil
}

// SetSecret stores a secret
func (c *HTTPClient) SetSecret(ctx context.Context, path string, secret *types.Secret) error {
	requestData := map[string]interface{}{
		"data": secret.Data,
	}

	resp, err := c.makeRequest("POST", "/v1/secret/"+path, requestData)
	if err != nil {
		return fmt.Errorf("failed to set secret: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set secret failed with status: %d", resp.StatusCode)
	}

	return nil
}

// DeleteSecret removes a secret
func (c *HTTPClient) DeleteSecret(ctx context.Context, path string) error {
	resp, err := c.makeRequest("DELETE", "/v1/secret/"+path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete secret failed with status: %d", resp.StatusCode)
	}

	return nil
}

// ListSecrets lists secrets
func (c *HTTPClient) ListSecrets(ctx context.Context, prefix string) ([]*types.SecretMetadata, error) {
	resp, err := c.makeRequest("LIST", "/v1/secret/"+prefix, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list secrets failed with status: %d", resp.StatusCode)
	}

	var response struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode list response: %w", err)
	}

	var secrets []*types.SecretMetadata
	for _, key := range response.Data.Keys {
		secrets = append(secrets, &types.SecretMetadata{
			Path:      key,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
			CreatedBy: "vault-server",
			UpdatedBy: "vault-server",
		})
	}

	return secrets, nil
}

// GetPath gets path information
func (c *HTTPClient) GetPath(ctx context.Context, path string) (*types.PathInfo, error) {
	// TODO: Implement path retrieval
	return &types.PathInfo{
		Path:      path,
		Type:      "folder",
		Metadata:  &types.PathMetadata{},
		Children:  []*types.PathInfo{},
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}, nil
}

// CreatePath creates a new path
func (c *HTTPClient) CreatePath(ctx context.Context, path string, metadata *types.PathMetadata) error {
	// TODO: Implement path creation
	return nil
}

// DeletePath removes a path
func (c *HTTPClient) DeletePath(ctx context.Context, path string) error {
	// TODO: Implement path deletion
	return nil
}

// Sync performs sync operations (no-op for HTTP client)
func (c *HTTPClient) Sync(ctx context.Context, direction types.SyncDirection) (*types.SyncResult, error) {
	return &types.SyncResult{
		Direction:    direction,
		ItemsSynced:  0,
		ItemsCreated: 0,
		ItemsUpdated: 0,
		ItemsDeleted: 0,
		Conflicts:    0,
		Errors:       []string{},
		StartTime:    time.Now().Unix(),
		EndTime:      time.Now().Unix(),
		Duration:     0,
	}, nil
}

// GetSyncStatus gets sync status (no-op for HTTP client)
func (c *HTTPClient) GetSyncStatus(ctx context.Context) (*types.SyncStatus, error) {
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
func (c *HTTPClient) GetConfig() *types.ClientConfig {
	return &types.ClientConfig{
		Type:      "http",
		ServerURL: c.baseURL,
	}
}

// SetConfig updates the client configuration
func (c *HTTPClient) SetConfig(config *types.ClientConfig) error {
	if config.ServerURL != "" {
		c.baseURL = config.ServerURL
	}
	return nil
}

// Health checks the health of the server
func (c *HTTPClient) Health(ctx context.Context) (*types.HealthStatus, error) {
	resp, err := c.makeRequest("GET", "/v1/sys/health", nil)
	if err != nil {
		return &types.HealthStatus{
			Status:  "unhealthy",
			Version: "1.0.0",
			Uptime:  0,
			System:  &types.ClientSystemInfo{},
			Checks:  []types.HealthCheck{},
		}, err
	}
	defer resp.Body.Close()

	healthy := resp.StatusCode == http.StatusOK
	return &types.HealthStatus{
		Status:  map[bool]string{true: "healthy", false: "unhealthy"}[healthy],
		Version: "1.0.0",
		Uptime:  0,
		System:  &types.ClientSystemInfo{},
		Checks:  []types.HealthCheck{},
	}, nil
}

// Status returns the client status
func (c *HTTPClient) Status(ctx context.Context) (*types.ClientStatus, error) {
	return &types.ClientStatus{
		Connected:     true,
		Mode:          "http",
		ServerURL:     c.baseURL,
		LastSync:      nil,
		LocalPath:     "",
		Authenticated: c.IsAuthenticated(),
	}, nil
}
