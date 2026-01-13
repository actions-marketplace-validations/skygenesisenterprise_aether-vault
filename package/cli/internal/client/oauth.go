package client

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// OAuthClient handles OAuth authentication flow
type OAuthClient struct {
	clientID     string
	redirectURI  string
	authURL      string
	tokenURL     string
	httpClient   *http.Client
	state        string
	codeVerifier string
}

// NewOAuthClient creates a new OAuth client
func NewOAuthClient(baseURL string) *OAuthClient {
	return &OAuthClient{
		clientID:    "aether-vault-cli",
		redirectURI: "http://127.0.0.1:8250/callback",
		authURL:     baseURL + "/oauth/authorize",
		tokenURL:    baseURL + "/oauth/token",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// StartOAuthFlow initiates the OAuth authentication flow
func (c *OAuthClient) StartOAuthFlow(ctx context.Context) (*types.AuthResponse, error) {
	// Generate PKCE parameters
	if err := c.generatePKCE(); err != nil {
		return nil, fmt.Errorf("failed to generate PKCE parameters: %w", err)
	}

	// Generate state
	if err := c.generateState(); err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	// Build authorization URL
	authURL, err := c.buildAuthURL()
	if err != nil {
		return nil, fmt.Errorf("failed to build auth URL: %w", err)
	}

	// Start local callback server
	callbackChan := make(chan string, 1)
	callbackServer, err := c.startCallbackServer(callbackChan)
	if err != nil {
		return nil, fmt.Errorf("failed to start callback server: %w", err)
	}
	defer callbackServer.Close()

	// Open browser
	if err := c.openBrowser(authURL); err != nil {
		fmt.Printf("Could not open browser automatically. Please open this URL manually:\n")
		fmt.Printf("%s\n", authURL)
	} else {
		fmt.Printf("Opening browser for authentication...\n")
	}

	// Wait for callback
	select {
	case code := <-callbackChan:
		return c.exchangeCodeForToken(ctx, code)
	case <-time.After(10 * time.Minute):
		return nil, fmt.Errorf("authentication timed out")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// StartOAuthFlowWithCode starts OAuth flow with manual code entry
func (c *OAuthClient) StartOAuthFlowWithCode(ctx context.Context) (*types.AuthResponse, error) {
	// Generate PKCE parameters
	if err := c.generatePKCE(); err != nil {
		return nil, fmt.Errorf("failed to generate PKCE parameters: %w", err)
	}

	// Generate state
	if err := c.generateState(); err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	// Build authorization URL
	authURL, err := c.buildAuthURL()
	if err != nil {
		return nil, fmt.Errorf("failed to build auth URL: %w", err)
	}

	fmt.Printf("OAuth Authentication\n")
	fmt.Printf("===================\n\n")
	fmt.Printf("Please follow these steps to authenticate:\n\n")
	fmt.Printf("1. Open this URL in your browser:\n")
	fmt.Printf("   %s\n\n", authURL)
	fmt.Printf("2. Complete the authentication process\n")
	fmt.Printf("3. Copy the authorization code displayed\n")
	fmt.Printf("4. Enter the code below (or press Ctrl+C to cancel)\n\n")

	// Prompt for authorization code
	var code string
	for {
		fmt.Printf("Enter authorization code: ")
		if _, err := fmt.Scanln(&code); err != nil {
			if strings.Contains(err.Error(), "unexpected newline") {
				continue // Empty line, try again
			}
			return nil, fmt.Errorf("failed to read authorization code: %w", err)
		}

		code = strings.TrimSpace(code)
		if code == "" {
			fmt.Printf("Code cannot be empty. Please try again or press Ctrl+C to cancel.\n")
			continue
		}

		break
	}

	// Exchange code for token
	return c.exchangeCodeForToken(ctx, code)
}

// generatePKCE generates PKCE code verifier and challenge
func (c *OAuthClient) generatePKCE() error {
	// Generate code verifier (random string)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	c.codeVerifier = base64.RawURLEncoding.EncodeToString(bytes)
	return nil
}

// generateState generates random state parameter
func (c *OAuthClient) generateState() error {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	c.state = base64.RawURLEncoding.EncodeToString(bytes)
	return nil
}

// buildAuthURL builds the OAuth authorization URL
func (c *OAuthClient) buildAuthURL() (string, error) {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", c.clientID)
	params.Add("redirect_uri", c.redirectURI)
	params.Add("state", c.state)
	params.Add("scope", "vault.read vault.write vault.admin")
	params.Add("code_challenge", c.codeVerifier)
	params.Add("code_challenge_method", "S256")

	return fmt.Sprintf("%s?%s", c.authURL, params.Encode()), nil
}

// exchangeCodeForToken exchanges authorization code for access token
func (c *OAuthClient) exchangeCodeForToken(ctx context.Context, code string) (*types.AuthResponse, error) {
	// Prepare token request
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", c.clientID)
	data.Set("code", code)
	data.Set("redirect_uri", c.redirectURI)
	data.Set("code_verifier", c.codeVerifier)

	// Make token request
	req, err := http.NewRequestWithContext(ctx, "POST", c.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed: %s - %s", resp.Status, string(body))
	}

	// Parse token response
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
		Scope        string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Get user info
	userInfo, err := c.getUserInfo(ctx, tokenResp.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return &types.AuthResponse{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
		ExpiresIn:    tokenResp.ExpiresIn,
		Scope:        tokenResp.Scope,
		User:         userInfo,
	}, nil
}

// getUserInfo retrieves user information using the access token
func (c *OAuthClient) getUserInfo(ctx context.Context, accessToken string) (*types.UserInfo, error) {
	// This would call the user info endpoint
	// For now, return mock user info
	return &types.UserInfo{
		ID:          "user-123",
		Username:    "vault-user",
		Email:       "user@aethervault.com",
		DisplayName: "Vault User",
	}, nil
}

// startCallbackServer starts a local HTTP server to handle OAuth callback
func (c *OAuthClient) startCallbackServer(callbackChan chan<- string) (*http.Server, error) {
	server := &http.Server{
		Addr:    ":8250",
		Handler: http.HandlerFunc(c.handleCallback(callbackChan)),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Callback server error: %v\n", err)
		}
	}()

	return server, nil
}

// handleCallback handles the OAuth callback
func (c *OAuthClient) handleCallback(callbackChan chan<- string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Handle error response
		if errorParam := r.URL.Query().Get("error"); errorParam != "" {
			errorDesc := r.URL.Query().Get("error_description")
			fmt.Fprintf(w, "Authentication failed: %s - %s", errorParam, errorDesc)
			callbackChan <- "" // Send empty string to indicate error
			return
		}

		// Handle successful response
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		// Verify state
		if state != c.state {
			fmt.Fprintf(w, "Invalid state parameter")
			callbackChan <- "" // Send empty string to indicate error
			return
		}

		// Send success page
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Authentication Successful</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; }
        .success { color: #10b981; font-size: 24px; margin-bottom: 20px; }
        .message { color: #6b7280; margin-bottom: 30px; }
        .close { color: #3b82f6; text-decoration: none; padding: 10px 20px; border: 1px solid #3b82f6; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="success">✓ Authentication Successful</div>
    <div class="message">You can now close this window and return to the terminal.</div>
    <a href="javascript:window.close()" class="close">Close Window</a>
</body>
</html>`)

		// Send code to channel
		callbackChan <- code
	}
}

// openBrowser opens the specified URL in the default browser
func (c *OAuthClient) openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		// Try various Linux browsers
		browsers := []string{"xdg-open", "google-chrome", "firefox", "mozilla"}
		for _, browser := range browsers {
			if _, err := exec.LookPath(browser); err == nil {
				cmd = exec.Command(browser, url)
				break
			}
		}
		if cmd == nil {
			return fmt.Errorf("no suitable browser found")
		}
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

// SaveCredentials saves authentication credentials to local storage
func (c *OAuthClient) SaveCredentials(authResp *types.AuthResponse) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	vaultDir := filepath.Join(home, ".aether", "vault")
	if err := os.MkdirAll(vaultDir, 0700); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	credsFile := filepath.Join(vaultDir, "credentials.json")

	// Save credentials (in a real implementation, this should be encrypted)
	creds := map[string]interface{}{
		"access_token":  authResp.AccessToken,
		"refresh_token": authResp.RefreshToken,
		"token_type":    authResp.TokenType,
		"expires_in":    authResp.ExpiresIn,
		"scope":         authResp.Scope,
		"expires_at":    time.Now().Add(time.Duration(authResp.ExpiresIn) * time.Second).Unix(),
		"user":          authResp.User,
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	if err := os.WriteFile(credsFile, data, 0600); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	return nil
}

// LoadCredentials loads saved authentication credentials
func (c *OAuthClient) LoadCredentials() (*types.AuthResponse, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	credsFile := filepath.Join(home, ".aether", "vault", "credentials.json")

	data, err := os.ReadFile(credsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	var creds map[string]interface{}
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credentials: %w", err)
	}

	// Check if token is expired
	if expiresAt, ok := creds["expires_at"].(float64); ok {
		if time.Now().Unix() > int64(expiresAt) {
			return nil, fmt.Errorf("credentials have expired")
		}
	}

	// Parse user info
	var userInfo *types.UserInfo
	if userMap, ok := creds["user"].(map[string]interface{}); ok {
		userInfo = &types.UserInfo{}
		if id, ok := userMap["id"].(string); ok {
			userInfo.ID = id
		}
		if username, ok := userMap["username"].(string); ok {
			userInfo.Username = username
		}
		if email, ok := userMap["email"].(string); ok {
			userInfo.Email = email
		}
		if displayName, ok := userMap["display_name"].(string); ok {
			userInfo.DisplayName = displayName
		}
	}

	return &types.AuthResponse{
		AccessToken:  creds["access_token"].(string),
		RefreshToken: creds["refresh_token"].(string),
		TokenType:    creds["token_type"].(string),
		ExpiresIn:    int64(creds["expires_in"].(float64)),
		Scope:        creds["scope"].(string),
		User:         userInfo,
	}, nil
}

// ClearCredentials removes saved authentication credentials
func (c *OAuthClient) ClearCredentials() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	credsFile := filepath.Join(home, ".aether", "vault", "credentials.json")

	if err := os.Remove(credsFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove credentials file: %w", err)
	}

	return nil
}
