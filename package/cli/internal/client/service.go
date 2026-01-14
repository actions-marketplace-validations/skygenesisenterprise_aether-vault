package client

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// ServiceManager manages the vault server service
type ServiceManager struct {
	binaryPath string
	pidFile    string
	logFile    string
	configPath string
}

// NewServiceManager creates a new service manager
func NewServiceManager() *ServiceManager {
	home, _ := os.UserHomeDir()
	vaultDir := filepath.Join(home, ".aether", "vault")

	return &ServiceManager{
		binaryPath: os.Args[0], // Path to current binary
		pidFile:    filepath.Join(vaultDir, "vault-server.pid"),
		logFile:    filepath.Join(vaultDir, "vault-server.log"),
		configPath: filepath.Join(vaultDir, "config.yaml"),
	}
}

// StartServer starts the vault server as a service
func (sm *ServiceManager) StartServer(devMode bool, configPath string) error {
	// Check if server is already running
	if sm.IsRunning() {
		return fmt.Errorf("vault server is already running")
	}

	// Ensure vault directory exists
	if err := os.MkdirAll(filepath.Dir(sm.pidFile), 0755); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	// Prepare command arguments
	args := []string{"server", "start"}
	if devMode {
		args = append(args, "--dev")
	}
	if configPath != "" {
		args = append(args, "--config", configPath)
	}

	var cmd *exec.Cmd

	// Start server based on OS
	switch runtime.GOOS {
	case "windows":
		// On Windows, start in background
		cmd = exec.Command(sm.binaryPath, args...)
	default:
		// On Unix-like systems, use nohup to detach from terminal
		allArgs := append([]string{sm.binaryPath}, args...)
		cmd = exec.Command("nohup", allArgs...)
	}

	// Redirect output to log file
	logFile, err := os.OpenFile(sm.logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	// Save PID
	pid := cmd.Process.Pid
	if err := os.WriteFile(sm.pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		// Try to kill the process if we can't save PID
		cmd.Process.Kill()
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	// Wait a moment and check if process is still running
	time.Sleep(2 * time.Second)
	if !sm.IsProcessRunning(pid) {
		os.Remove(sm.pidFile)
		return fmt.Errorf("server failed to start, check logs at %s", sm.logFile)
	}

	// Initialize service if not already done
	if err := sm.initializeService(); err != nil {
		return fmt.Errorf("failed to initialize service: %w", err)
	}

	fmt.Printf("Vault server started successfully (PID: %d)\n", pid)
	fmt.Printf("Logs: %s\n", sm.logFile)

	if devMode {
		fmt.Println("Development mode enabled")
		fmt.Println("Root token: dev-token")
		fmt.Println("Server URL: http://127.0.0.1:8200")
	}

	return nil
}

// StopServer stops the vault server service
func (sm *ServiceManager) StopServer() error {
	if !sm.IsRunning() {
		return fmt.Errorf("vault server is not running")
	}

	// Read PID
	pidData, err := os.ReadFile(sm.pidFile)
	if err != nil {
		return fmt.Errorf("failed to read PID file: %w", err)
	}

	// Parse PID
	var pid int
	if _, err := fmt.Sscanf(string(pidData), "%d", &pid); err != nil {
		return fmt.Errorf("failed to parse PID: %w", err)
	}

	// Find and kill the process
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %w", err)
	}

	// Try graceful shutdown first
	if err := process.Signal(os.Interrupt); err == nil {
		// Wait a bit for graceful shutdown
		time.Sleep(3 * time.Second)
		if !sm.IsProcessRunning(pid) {
			os.Remove(sm.pidFile)
			fmt.Println("Vault server stopped gracefully")
			return nil
		}
	}

	// Force kill if graceful didn't work
	if err := process.Kill(); err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}

	// Wait and verify
	time.Sleep(1 * time.Second)
	if !sm.IsProcessRunning(pid) {
		os.Remove(sm.pidFile)
		fmt.Println("Vault server stopped forcefully")
		return nil
	}

	return fmt.Errorf("failed to stop server process")
}

// RestartServer restarts the vault server service
func (sm *ServiceManager) RestartServer(devMode bool, configPath string) error {
	fmt.Println("Restarting vault server...")

	// Stop if running
	if sm.IsRunning() {
		if err := sm.StopServer(); err != nil {
			return fmt.Errorf("failed to stop server: %w", err)
		}
		time.Sleep(2 * time.Second)
	}

	// Start again
	return sm.StartServer(devMode, configPath)
}

// IsRunning checks if the server is currently running
func (sm *ServiceManager) IsRunning() bool {
	if _, err := os.Stat(sm.pidFile); os.IsNotExist(err) {
		return false
	}

	pidData, err := os.ReadFile(sm.pidFile)
	if err != nil {
		return false
	}

	var pid int
	if _, err := fmt.Sscanf(string(pidData), "%d", &pid); err != nil {
		return false
	}

	return sm.IsProcessRunning(pid)
}

// IsProcessRunning checks if a specific process is running
func (sm *ServiceManager) IsProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists
	err = process.Signal(os.Signal(nil))
	return err == nil
}

// GetStatus returns the current status of the service
func (sm *ServiceManager) GetStatus() *ServiceStatus {
	status := &ServiceStatus{
		Running: sm.IsRunning(),
		PIDFile: sm.pidFile,
		LogFile: sm.logFile,
	}

	if status.Running {
		if pidData, err := os.ReadFile(sm.pidFile); err == nil {
			fmt.Sscanf(string(pidData), "%d", &status.PID)
		}

		// Try to get uptime from log file
		if info, err := os.Stat(sm.logFile); err == nil {
			status.StartTime = info.ModTime()
		}
	}

	return status
}

// ServiceStatus represents the status of the vault server service
type ServiceStatus struct {
	Running   bool
	PID       int
	PIDFile   string
	LogFile   string
	StartTime time.Time
}

// initializeService sets up the service environment
func (sm *ServiceManager) initializeService() error {
	// Create default config if it doesn't exist
	if _, err := os.Stat(sm.configPath); os.IsNotExist(err) {
		defaultConfig := `
# Aether Vault Configuration
server:
  host: "127.0.0.1"
  port: 8200
  
local:
  path: "~/.aether/vault/data"
  
ui:
  color: true
  spinner: true
  table_style: "default"
`
		if err := os.WriteFile(sm.configPath, []byte(defaultConfig), 0644); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
	}

	// Set up socket directory for IPC
	socketDir := filepath.Join(filepath.Dir(sm.pidFile), "ipc")
	if err := os.MkdirAll(socketDir, 0755); err != nil {
		return fmt.Errorf("failed to create IPC directory: %w", err)
	}

	return nil
}

// EnsureServer ensures the server is running
func (sm *ServiceManager) EnsureServer() error {
	if sm.IsRunning() {
		return nil
	}

	fmt.Println("Vault server is not running, starting it...")
	return sm.StartServer(true, "")
}

// GetServerURL returns the server URL
func (sm *ServiceManager) GetServerURL() string {
	return "http://127.0.0.1:8200"
}

// WaitUntilServerReady waits until the server is ready to accept requests
func (sm *ServiceManager) WaitUntilServerReady(ctx context.Context, timeout time.Duration) error {
	client := NewHTTPClient(sm.GetServerURL())

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if _, err := client.Health(ctx); err == nil {
				return nil
			}
			time.Sleep(500 * time.Millisecond)
		}
	}

	return fmt.Errorf("server did not become ready within %v", timeout)
}
