package runtime

import (
	"os"
	"os/user"
	"runtime"

	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// DetectEnvironment detects the current runtime environment
func DetectEnvironment() *types.Environment {
	env := &types.Environment{
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		InDocker:     IsInDockerContainer(),
		InCI:         IsRunningInCI(),
		HasDocker:    IsDockerInstalled(),
		DockerInfo:   GetDockerInfo(),
		PlatformInfo: GetPlatformInfo(),
	}

	// Get user information
	if currentUser, err := user.Current(); err == nil {
		env.User = &types.EnvironmentUser{
			Username: currentUser.Username,
			HomeDir:  currentUser.HomeDir,
			UID:      currentUser.Uid,
			GID:      currentUser.Gid,
		}
	}

	// Get shell information
	env.Shell = GetShell()

	// Get terminal information
	env.Terminal = &types.TerminalInfo{
		SupportsColor:   SupportsColor(),
		SupportsUnicode: SupportsUnicode(),
		Width:           GetTerminalWidth(),
		Height:          GetTerminalHeight(),
	}

	return env
}

// GetEnvironmentDetails returns detailed environment information
func GetEnvironmentDetails() *types.EnvironmentDetails {
	env := DetectEnvironment()

	details := &types.EnvironmentDetails{
		Environment: env,
		SystemInfo:  GetRuntimeSystemInfo(),
		Paths: &types.EnvironmentPaths{
			Home:       GetHomeDir(),
			Temp:       GetTempDir(),
			Working:    getCurrentWorkingDir(),
			Executable: getCurrentExecutable(),
		},
		Variables: getEnvironmentVariables(),
	}

	return details
}

// getCurrentWorkingDir returns current working directory
func getCurrentWorkingDir() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return ""
}

// getCurrentExecutable returns current executable path
func getCurrentExecutable() string {
	if exe, err := os.Executable(); err == nil {
		return exe
	}
	return ""
}

// getEnvironmentVariables returns relevant environment variables
func getEnvironmentVariables() map[string]string {
	vars := make(map[string]string)

	relevantVars := []string{
		"PATH", "HOME", "USER", "SHELL", "TERM",
		"LANG", "LC_ALL", "NO_COLOR", "CLICOLOR",
		"VAULT_URL", "VAULT_TOKEN", "VAULT_PATH",
		"DOCKER_HOST", "KUBECONFIG",
	}

	for _, name := range relevantVars {
		if value := os.Getenv(name); value != "" {
			vars[name] = value
		}
	}

	return vars
}

// IsSupportedEnvironment returns true if the current environment is supported
func IsSupportedEnvironment() bool {
	env := DetectEnvironment()

	// Check OS support
	switch env.OS {
	case "linux", "darwin", "windows":
		// Supported
	default:
		return false
	}

	// Check architecture support
	switch env.Arch {
	case "amd64", "arm64":
		// Supported
	default:
		return false
	}

	return true
}

// GetSystemCapabilities returns system capabilities
func GetSystemCapabilities() *types.SystemCapabilities {
	caps := &types.SystemCapabilities{
		SupportsDocker:  IsDockerInstalled() && isDockerRunningNoError(),
		SupportsUnicode: SupportsUnicode(),
		SupportsColor:   SupportsColor(),
		IsUnix:          IsUnix(),
		IsWindows:       IsWindows(),
		HasSudo:         checkSudoAccess(),
		HasInternet:     checkInternetConnection(),
		HasGit:          checkGitInstalled(),
	}

	// Add platform-specific capabilities
	switch runtime.GOOS {
	case "linux":
		caps.SupportsSystemd = checkSystemdInstalled()
		caps.SupportsSelinux = checkSelinuxInstalled()
	case "darwin":
		caps.SupportsHomebrew = checkHomebrewInstalled()
	case "windows":
		caps.SupportsPowershell = checkPowershellInstalled()
		caps.SupportsWsl = checkWslInstalled()
	}

	return caps
}

// checkSudoAccess checks if sudo access is available
func checkSudoAccess() bool {
	// TODO: Implement sudo access check
	return false
}

// checkInternetConnection checks if internet connection is available
func checkInternetConnection() bool {
	// TODO: Implement internet connection check
	return false
}

// checkGitInstalled checks if Git is installed
func checkGitInstalled() bool {
	// TODO: Implement Git installation check
	return false
}

// checkSystemdInstalled checks if systemd is installed
func checkSystemdInstalled() bool {
	// TODO: Implement systemd check
	return false
}

// checkSelinuxInstalled checks if SELinux is installed
func checkSelinuxInstalled() bool {
	// TODO: Implement SELinux check
	return false
}

// checkHomebrewInstalled checks if Homebrew is installed
func checkHomebrewInstalled() bool {
	// TODO: Implement Homebrew check
	return false
}

// checkPowershellInstalled checks if PowerShell is installed
func checkPowershellInstalled() bool {
	// TODO: Implement PowerShell check
	return false
}

// checkWslInstalled checks if WSL is installed
func checkWslInstalled() bool {
	// TODO: Implement WSL check
	return false
}
