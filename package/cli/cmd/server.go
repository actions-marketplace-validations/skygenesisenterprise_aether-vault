package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/client"
	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
	"github.com/spf13/cobra"
)

// newServerCommand creates the server command
func newServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Manage Vault server service",
		Long: `Manage the Aether Vault server service.

This command controls the local vault server service that provides
secret management functionality similar to 1Password.

Examples:
  vault server              Start server service
  vault server -dev         Start in development mode
  vault server stop         Stop server service
  vault server restart      Restart server service
  vault server status       Show server status
  vault server logs         Show server logs`,
	}

	// Add subcommands
	cmd.AddCommand(newServerStartCommand())
	cmd.AddCommand(newServerStopCommand())
	cmd.AddCommand(newServerRestartCommand())
	cmd.AddCommand(newServerStatusCommand())
	cmd.AddCommand(newServerLogsCommand())

	// Default behavior is start
	cmd.RunE = runServerCommand

	return cmd
}

// newServerStartCommand creates the server start command
func newServerStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the Vault server service",
		Long:  "Start the Aether Vault server service in the background",
		RunE:  runServerStartCommand,
	}

	cmd.Flags().String("config", "", "Configuration file path")
	cmd.Flags().Bool("dev", false, "Start in development mode")
	cmd.Flags().String("log-level", "info", "Log level (trace, debug, info, warn, error)")

	return cmd
}

// newServerStopCommand creates the server stop command
func newServerStopCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the Vault server service",
		Long:  "Stop the running Aether Vault server service",
		RunE:  runServerStopCommand,
	}

	return cmd
}

// newServerRestartCommand creates the server restart command
func newServerRestartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart the Vault server service",
		Long:  "Restart the Aether Vault server service",
		RunE:  runServerRestartCommand,
	}

	cmd.Flags().String("config", "", "Configuration file path")
	cmd.Flags().Bool("dev", false, "Start in development mode")

	return cmd
}

// newServerStatusCommand creates the server status command
func newServerStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show Vault server service status",
		Long:  "Display the current status of the Aether Vault server service",
		RunE:  runServerStatusCommand,
	}

	return cmd
}

// newServerLogsCommand creates the server logs command
func newServerLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Show Vault server logs",
		Long:  "Display the logs from the Aether Vault server service",
		RunE:  runServerLogsCommand,
	}

	cmd.Flags().Int("tail", 50, "Number of lines to show from the end")
	cmd.Flags().Bool("follow", false, "Follow log output")

	return cmd
}

// VaultServer represents the local vault server
type VaultServer struct {
	httpServer *http.Server
	client     client.Client
	router     *mux.Router
	config     *ServerConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host        string
	Port        int
	DevMode     bool
	LogLevel    string
	ConfigPath  string
	StoragePath string
	RootToken   string
}

// runServerCommand executes the server command (default to start)
func runServerCommand(cmd *cobra.Command, args []string) error {
	return runServerStartCommand(cmd, args)
}

// runServerStartCommand executes the server start command
func runServerStartCommand(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")
	devMode, _ := cmd.Flags().GetBool("dev")
	logLevel, _ := cmd.Flags().GetString("log-level")

	// Create service manager
	serviceManager := client.NewServiceManager()

	// Check if we should run in foreground or as service
	if os.Getenv("VAULT_SERVER_FG") == "1" {
		// Run server in foreground (for service mode)
		config := &ServerConfig{
			Host:       "127.0.0.1",
			Port:       8200,
			DevMode:    devMode,
			LogLevel:   logLevel,
			ConfigPath: configPath,
			RootToken:  "dev-token",
		}

		if devMode {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			config.StoragePath = filepath.Join(home, ".aether", "vault", "dev")
		} else {
			config.StoragePath = filepath.Join(filepath.Dir(configPath), "data")
		}

		server, err := NewVaultServer(config)
		if err != nil {
			return fmt.Errorf("failed to create vault server: %w", err)
		}

		return server.Start()
	} else {
		// Run as service (background)
		fmt.Println("Starting Aether Vault server service...")

		if serviceManager.IsRunning() {
			fmt.Println("Vault server is already running")
			status := serviceManager.GetStatus()
			if status.Running {
				fmt.Printf("PID: %d\n", status.PID)
				fmt.Printf("Logs: %s\n", status.LogFile)
				fmt.Printf("Server URL: %s\n", serviceManager.GetServerURL())
			}
			return nil
		}

		err := serviceManager.StartServer(devMode, configPath)
		if err != nil {
			return fmt.Errorf("failed to start server service: %w", err)
		}

		fmt.Println("Vault server service started successfully")
		fmt.Printf("Server URL: %s\n", serviceManager.GetServerURL())
		fmt.Printf("Use 'vault server stop' to stop the service")

		return nil
	}
}

// runServerStopCommand executes the server stop command
func runServerStopCommand(cmd *cobra.Command, args []string) error {
	serviceManager := client.NewServiceManager()

	if !serviceManager.IsRunning() {
		fmt.Println("Vault server is not running")
		return nil
	}

	fmt.Println("Stopping Aether Vault server service...")
	err := serviceManager.StopServer()
	if err != nil {
		return fmt.Errorf("failed to stop server service: %w", err)
	}

	fmt.Println("Vault server service stopped successfully")
	return nil
}

// runServerRestartCommand executes the server restart command
func runServerRestartCommand(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")
	devMode, _ := cmd.Flags().GetBool("dev")

	serviceManager := client.NewServiceManager()

	err := serviceManager.RestartServer(devMode, configPath)
	if err != nil {
		return fmt.Errorf("failed to restart server service: %w", err)
	}

	fmt.Println("Vault server service restarted successfully")
	fmt.Printf("Server URL: %s\n", serviceManager.GetServerURL())

	return nil
}

// runServerStatusCommand executes the server status command
func runServerStatusCommand(cmd *cobra.Command, args []string) error {
	serviceManager := client.NewServiceManager()
	status := serviceManager.GetStatus()

	fmt.Println("Aether Vault Server Status:")
	fmt.Println("===========================")

	if status.Running {
		fmt.Printf("Status: Running\n")
		fmt.Printf("PID: %d\n", status.PID)
		fmt.Printf("Server URL: %s\n", serviceManager.GetServerURL())
		if !status.StartTime.IsZero() {
			fmt.Printf("Started: %s\n", status.StartTime.Format("2006-01-02 15:04:05"))
		}
	} else {
		fmt.Printf("Status: Not running\n")
	}

	fmt.Printf("Log file: %s\n", status.LogFile)
	fmt.Printf("PID file: %s\n", status.PIDFile)

	return nil
}

// runServerLogsCommand executes the server logs command
func runServerLogsCommand(cmd *cobra.Command, args []string) error {
	tail, _ := cmd.Flags().GetInt("tail")
	follow, _ := cmd.Flags().GetBool("follow")

	serviceManager := client.NewServiceManager()
	status := serviceManager.GetStatus()

	if _, err := os.Stat(status.LogFile); os.IsNotExist(err) {
		fmt.Println("No log file found. Server may not have started yet.")
		return nil
	}

	fmt.Printf("Showing logs from %s\n", status.LogFile)
	fmt.Println(strings.Repeat("=", 50))

	if follow {
		// Follow logs (simple implementation)
		// In a real implementation, you'd use tail -f or a proper log follower
		fmt.Println("Following logs (press Ctrl+C to stop)...")
		// This is a simplified implementation
		cmd := exec.Command("tail", "-f", status.LogFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	} else {
		// Show last N lines
		cmd := exec.Command("tail", "-n", fmt.Sprintf("%d", tail), status.LogFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}

// NewVaultServer creates a new vault server instance
func NewVaultServer(config *ServerConfig) (*VaultServer, error) {
	// Create local client
	clientConfig := &types.ClientConfig{
		Type: "local",
		Options: map[string]interface{}{
			"basePath": config.StoragePath,
		},
	}

	clientOptions := &client.ClientOptions{
		Config: clientConfig,
	}

	localClient, err := client.NewLocalClient(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create local client: %w", err)
	}

	router := mux.NewRouter()
	server := &VaultServer{
		client: localClient,
		router: router,
		config: config,
	}

	server.setupRoutes()

	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server, nil
}

// Start starts the vault server
func (s *VaultServer) Start() error {
	fmt.Printf("Starting Aether Vault server on %s...\n", s.httpServer.Addr)

	if s.config.DevMode {
		fmt.Println("Running in development mode")
		fmt.Printf("Root token: %s\n", s.config.RootToken)
		fmt.Printf("Storage path: %s\n", s.config.StoragePath)
	}

	// Ensure storage directory exists
	if err := os.MkdirAll(s.config.StoragePath, 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	fmt.Println("Server started successfully. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	fmt.Println("Server stopped.")
	return nil
}

// setupRoutes configures the API routes
func (s *VaultServer) setupRoutes() {
	// Health check
	s.router.HandleFunc("/v1/sys/health", s.handleHealth).Methods("GET")

	// Authentication
	s.router.HandleFunc("/v1/auth/token/login", s.handleTokenLogin).Methods("POST")

	// Secrets management
	s.router.HandleFunc("/v1/secret/{path:.*}", s.handleSecret).Methods("GET", "POST", "DELETE")
	s.router.HandleFunc("/v1/secret/{path:.*}/list", s.handleSecretList).Methods("LIST")

	// Status
	s.router.HandleFunc("/v1/sys/status", s.handleStatus).Methods("GET")
}

// handleHealth handles health check requests
func (s *VaultServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"initialized":                  true,
		"sealed":                       false,
		"standby":                      false,
		"performance_standby":          false,
		"replication_performance_mode": "none",
		"replication_dr_mode":          "disabled",
		"server_time_utc":              time.Now().UTC().Unix(),
		"version":                      "1.0.0",
		"cluster_name":                 "aether-vault",
		"cluster_id":                   "local",
	}
	json.NewEncoder(w).Encode(response)
}

// handleTokenLogin handles token authentication
func (s *VaultServer) handleTokenLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Role  string `json:"role"`
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Simple token validation
	if req.Token != s.config.RootToken {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": req.Token,
			"policies":     []string{"root"},
			"metadata": map[string]interface{}{
				"role": req.Role,
			},
			"lease_duration": 0,
			"renewable":      false,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// handleSecret handles secret CRUD operations
func (s *VaultServer) handleSecret(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		secret, err := s.client.GetSecret(r.Context(), path)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get secret: %v", err), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"data":     secret.Data,
			"metadata": secret.Metadata,
		}
		json.NewEncoder(w).Encode(response)

	case "POST":
		var req struct {
			Data map[string]interface{} `json:"data"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		secret := &types.Secret{
			Path: path,
			Data: req.Data,
			Metadata: &types.SecretMetadata{
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
				CreatedBy: "vault-server",
				UpdatedBy: "vault-server",
			},
			Version: 1,
		}

		if err := s.client.SetSecret(r.Context(), path, secret); err != nil {
			http.Error(w, fmt.Sprintf("Failed to set secret: %v", err), http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"data": req.Data,
		}
		json.NewEncoder(w).Encode(response)

	case "DELETE":
		if err := s.client.DeleteSecret(r.Context(), path); err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete secret: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// handleSecretList handles secret listing
func (s *VaultServer) handleSecretList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prefix := vars["path"]

	secrets, err := s.client.ListSecrets(r.Context(), prefix)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list secrets: %v", err), http.StatusInternalServerError)
		return
	}

	var keys []string
	for _, secret := range secrets {
		keys = append(keys, secret.Path)
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"keys": keys,
		},
	}
	json.NewEncoder(w).Encode(response)
}

// handleStatus handles status requests
func (s *VaultServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status, err := s.client.Status(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get status: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"initialized": true,
		"sealed":      false,
		"standby":     false,
		"version":     "1.0.0",
		"server_time": time.Now().UTC().Unix(),
		"mode":        string(status.Mode),
		"local_path":  status.LocalPath,
	}

	json.NewEncoder(w).Encode(response)
}
