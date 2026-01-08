package types

// ExecutionMode represents the execution mode of the CLI
type ExecutionMode string

const (
	// LocalMode represents offline/local execution
	LocalMode ExecutionMode = "local"

	// CloudMode represents connected/cloud execution
	CloudMode ExecutionMode = "cloud"
)

// Context represents the execution context
type Context struct {
	// Current execution mode
	Mode ExecutionMode

	// Configuration
	Config *Config

	// Runtime information
	Runtime *RuntimeInfo

	// Authentication state
	Auth *AuthState
}

// RuntimeInfo contains runtime environment information
type RuntimeInfo struct {
	// Operating system
	OS string

	// Architecture
	Arch string

	// Go version
	GoVersion string

	// CLI version
	Version string

	// Build information
	Build *BuildInfo

	// Environment variables
	Env map[string]string
}

// BuildInfo contains build information
type BuildInfo struct {
	// Git commit hash
	Commit string

	// Build timestamp
	Timestamp string

	// Build environment (dev, prod)
	Environment string

	// Build tools version
	ToolsVersion string
}

// AuthState represents authentication state
type AuthState struct {
	// Is authenticated
	Authenticated bool

	// Authentication method
	Method string

	// Token information
	Token *TokenInfo

	// User information
	User *UserInfo

	// Expiration time
	ExpiresAt *int64
}

// TokenInfo contains token information
type TokenInfo struct {
	// Access token
	AccessToken string

	// Refresh token
	RefreshToken string

	// Token type
	Type string

	// Token scope
	Scope string
}

// UserInfo contains user information
type UserInfo struct {
	// User ID
	ID string

	// Username
	Username string

	// Email
	Email string

	// Display name
	DisplayName string

	// Organization
	Organization string
}

// PlatformInfo contains platform-specific information
type PlatformInfo struct {
	// Operating system
	OS string

	// Architecture
	Arch string

	// Go version
	GoVersion string

	// Compiler
	Compiler string

	// Number of CPUs
	NumCPU int

	// Number of goroutines
	NumGoroutine int

	// Timestamp
	Timestamp int64

	// Memory usage
	MemoryUsage *MemoryInfo

	// Disk usage
	DiskUsage *DiskInfo

	// CPU usage
	CPUUsage float64
}

// RuntimeSystemInfo contains runtime system information
type RuntimeSystemInfo struct {
	// OS
	OS string

	// Architecture
	Arch string

	// Go version
	GoVersion string

	// Number of CPUs
	NumCPU int

	// Timestamp
	Timestamp int64

	// Memory usage
	MemoryUsage *MemoryInfo

	// Disk usage
	DiskUsage *DiskInfo

	// CPU usage
	CPUUsage float64
}

// DockerInfo contains Docker-related information
type DockerInfo struct {
	// Is Docker installed
	Installed bool

	// Is Docker daemon running
	Running bool

	// Docker version
	Version string

	// Additional Docker info
	Info map[string]string
}

// DockerImage represents a Docker image
type DockerImage struct {
	// Repository and tag
	Name string

	// Image ID
	ID string

	// Created timestamp
	Created string

	// Size
	Size string

	// Labels
	Labels map[string]string
}

// DockerContainer represents a Docker container
type DockerContainer struct {
	// Container name
	Name string

	// Container ID
	ID string

	// Status
	Status string

	// Image
	Image string

	// Ports
	Ports []string

	// Labels
	Labels map[string]string
}

// DockerNetwork represents a Docker network
type DockerNetwork struct {
	// Network name
	Name string

	// Network ID
	ID string

	// Driver
	Driver string

	// Subnet
	Subnet string

	// Gateway
	Gateway string

	// Labels
	Labels map[string]string
}

// Environment contains runtime environment information
type Environment struct {
	// Operating system
	OS string

	// Architecture
	Arch string

	// Is running in Docker
	InDocker bool

	// Is running in CI
	InCI bool

	// Has Docker available
	HasDocker bool

	// Docker information
	DockerInfo *DockerInfo

	// Platform information
	PlatformInfo *PlatformInfo

	// User information
	User *EnvironmentUser

	// Shell information
	Shell string

	// Terminal information
	Terminal *TerminalInfo
}

// EnvironmentUser contains user information
type EnvironmentUser struct {
	// Username
	Username string

	// Home directory
	HomeDir string

	// User ID
	UID string

	// Group ID
	GID string
}

// TerminalInfo contains terminal information
type TerminalInfo struct {
	// Supports color
	SupportsColor bool

	// Supports Unicode
	SupportsUnicode bool

	// Terminal width
	Width int

	// Terminal height
	Height int
}

// EnvironmentDetails contains detailed environment information
type EnvironmentDetails struct {
	// Basic environment
	Environment *Environment

	// System information
	SystemInfo *RuntimeSystemInfo

	// Path information
	Paths *EnvironmentPaths

	// Environment variables
	Variables map[string]string
}

// EnvironmentPaths contains path information
type EnvironmentPaths struct {
	// Home directory
	Home string

	// Temporary directory
	Temp string

	// Working directory
	Working string

	// Executable path
	Executable string
}

// SystemCapabilities contains system capabilities
type SystemCapabilities struct {
	// Supports Docker
	SupportsDocker bool

	// Supports Unicode
	SupportsUnicode bool

	// Supports color
	SupportsColor bool

	// Is Unix system
	IsUnix bool

	// Is Windows system
	IsWindows bool

	// Has sudo access
	HasSudo bool

	// Has internet connection
	HasInternet bool

	// Has Git installed
	HasGit bool

	// Has systemd installed (Linux)
	SupportsSystemd bool

	// Has SELinux installed (Linux)
	SupportsSelinux bool

	// Has Homebrew installed (macOS)
	SupportsHomebrew bool

	// Has PowerShell installed (Windows)
	SupportsPowershell bool

	// Has WSL installed (Windows)
	SupportsWsl bool
}

// SystemInfo contains system information for runtime
type SystemInfo struct {
	// Operating system
	OS string

	// Architecture
	Arch string

	// Go version
	GoVersion string

	// Number of CPUs
	NumCPU int

	// Timestamp
	Timestamp int64
}
