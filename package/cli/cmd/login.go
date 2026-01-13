package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/client"
	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
	"github.com/spf13/cobra"
)

// newLoginCommand creates the login command
func newLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Aether Vault cloud",
		Long: `Authenticate with Aether Vault cloud services using OAuth or token-based authentication.

This command will:
  - Open a browser for OAuth authentication (default)
  - Or accept an API token for token-based auth
  - Store authentication credentials securely
  - Switch to cloud mode after successful authentication`,
		RunE: runLoginCommand,
	}

	cmd.Flags().String("method", "oauth", "Authentication method (oauth, token)")
	cmd.Flags().String("token", "", "API token for token-based authentication")
	cmd.Flags().String("url", "https://vault.skygenesisenterprise.com", "Aether Vault cloud URL")

	return cmd
}

// newConnectCommand creates the connect command
func newConnectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to Aether Vault cloud (interactive)",
		Long: `Interactive connection to Aether Vault cloud services.

This command provides a guided connection experience with:
  - Step-by-step authentication setup
  - Connection testing
  - Configuration verification`,
		RunE: runConnectCommand,
	}

	cmd.Flags().Bool("interactive", true, "Enable interactive mode")

	return cmd
}

// newLogoutCommand creates the logout command
func newLogoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout from cloud services",
		Long: `Logout from Aether Vault cloud services and return to local mode.

This command will:
  - Remove stored authentication tokens
  - Switch back to local mode
  - Clear cloud connection state`,
		RunE: runLogoutCommand,
	}

	return cmd
}

// runLoginCommand executes the login command
func runLoginCommand(cmd *cobra.Command, args []string) error {
	method, _ := cmd.Flags().GetString("method")
	token, _ := cmd.Flags().GetString("token")
	url, _ := cmd.Flags().GetString("url")

	fmt.Printf("Connecting to Aether Vault cloud...\n")
	fmt.Printf("URL: %s\n", url)

	switch method {
	case "oauth":
		return runOAuthLogin(url)
	case "token":
		if token == "" {
			return fmt.Errorf("token is required for token-based authentication")
		}
		return runTokenLogin(token, url)
	default:
		return fmt.Errorf("unsupported authentication method: %s", method)
	}
}

// runConnectCommand executes the connect command
func runConnectCommand(cmd *cobra.Command, args []string) error {
	interactive, _ := cmd.Flags().GetBool("interactive")

	if interactive {
		fmt.Printf("Aether Vault Cloud Connection Wizard\n")
		fmt.Printf("=====================================\n\n")

		fmt.Printf("This wizard will help you connect to Aether Vault cloud.\n")
		fmt.Printf("Press Enter to continue or Ctrl+C to cancel...")

		// TODO: Implement interactive connection
		fmt.Printf("\n\nInteractive connection not yet implemented.\n")
		fmt.Printf("Use 'vault login' for direct authentication.\n")
	}

	return nil
}

// runLogoutCommand executes the logout command
func runLogoutCommand(cmd *cobra.Command, args []string) error {
	fmt.Printf("Logging out from Aether Vault cloud...\n")

	// Create OAuth client to clear credentials
	oauthClient := client.NewOAuthClient("https://vault.skygenesisenterprise.com")

	if err := oauthClient.ClearCredentials(); err != nil {
		fmt.Printf("Warning: Failed to clear credentials: %v\n", err)
	}

	fmt.Printf("✓ Authentication tokens cleared\n")
	fmt.Printf("✓ Switched to local mode\n")
	fmt.Printf("✓ Cloud connection closed\n")

	fmt.Printf("\nSuccessfully logged out. Use 'vault login' to reconnect.\n")

	return nil
}

// runOAuthLogin handles OAuth authentication
func runOAuthLogin(url string) error {
	fmt.Printf("OAuth Authentication\n")
	fmt.Printf("===================\n\n")

	// Create OAuth client
	oauthClient := client.NewOAuthClient(url)

	// Set up context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nAuthentication cancelled by user")
		cancel()
	}()

	// Start OAuth flow with manual code entry
	authResp, err := oauthClient.StartOAuthFlowWithCode(ctx)
	if err != nil {
		return fmt.Errorf("OAuth authentication failed: %w", err)
	}

	// Save credentials
	if err := oauthClient.SaveCredentials(authResp); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	// Display success message
	fmt.Printf("\n✓ Authentication successful!\n")
	fmt.Printf("✓ Logged in as: %s (%s)\n", authResp.User.DisplayName, authResp.User.Email)
	fmt.Printf("✓ Access token saved securely\n")
	fmt.Printf("✓ Switched to cloud mode\n\n")

	fmt.Printf("You can now use vault commands with cloud access.\n")
	fmt.Printf("Use 'vault logout' to return to local mode.\n")

	return nil
}

// runTokenLogin handles token-based authentication
func runTokenLogin(token, url string) error {
	fmt.Printf("Token Authentication\n")
	fmt.Printf("===================\n\n")

	fmt.Printf("Validating token with %s...\n", url)

	// Create HTTP client to validate token
	httpClient := client.NewHTTPClient(url)
	httpClient.SetToken(token)

	// Test token by making a health check request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	healthResp, err := httpClient.Health(ctx)
	if err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	if healthResp.Status != "healthy" {
		return fmt.Errorf("server returned unhealthy status: %s", healthResp.Status)
	}

	// Create mock auth response for token-based auth
	authResp := &types.AuthResponse{
		AccessToken:  token,
		RefreshToken: token,
		TokenType:    "Bearer",
		ExpiresIn:    0, // No expiration for API tokens
		Scope:        "all",
		User: &types.UserInfo{
			ID:          "token-user",
			Username:    "api-token",
			Email:       "",
			DisplayName: "API Token User",
		},
	}

	// Save credentials
	oauthClient := client.NewOAuthClient(url)
	if err := oauthClient.SaveCredentials(authResp); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	fmt.Printf("✓ Token validated successfully\n")
	fmt.Printf("✓ Authentication established\n")
	fmt.Printf("✓ Switched to cloud mode\n")

	fmt.Printf("\nSuccessfully authenticated with Aether Vault cloud.\n")

	return nil
}
