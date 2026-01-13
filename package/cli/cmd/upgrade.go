package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// GitHubRelease represents a GitHub API release
type GitHubRelease struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	Body       string `json:"body"`
	Draft      bool   `json:"draft"`
	PreRelease bool   `json:"prerelease"`
	Assets     []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// newUpgradeCommand creates the upgrade command
func newUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade CLI to latest version",
		Long: `Check for and upgrade to the latest version of Aether Vault CLI.

This command will:
1. Check the current version
2. Fetch the latest release from GitHub
3. Download and install if a newer version is available

Examples:
  vault upgrade              # Check and upgrade if needed
  vault upgrade --check-only  # Only check for updates
  vault upgrade --force       # Force upgrade even if current`,
		RunE: runUpgradeCommand,
	}

	cmd.Flags().Bool("check-only", false, "Only check for updates, don't upgrade")
	cmd.Flags().Bool("force", false, "Force upgrade even if current version is latest")
	cmd.Flags().String("repo", "https://github.com/skygenesisenterprise/aether-vault", "GitHub repository URL")

	return cmd
}

// runUpgradeCommand executes the upgrade command
func runUpgradeCommand(cmd *cobra.Command, args []string) error {
	checkOnly, _ := cmd.Flags().GetBool("check-only")
	force, _ := cmd.Flags().GetBool("force")
	repoURL, _ := cmd.Flags().GetString("repo")

	// Get the current version
	currentVersion := getCurrentVersion()
	fmt.Printf("Current version: %s\n", currentVersion)

	// Fetch the latest release
	latestRelease, err := fetchLatestRelease(repoURL)
	if err != nil {
		return fmt.Errorf("failed to fetch the latest release: %w", err)
	}

	if latestRelease == nil {
		return fmt.Errorf("no releases found")
	}

	fmt.Printf("Latest version: %s\n", latestRelease.TagName)

	// Compare versions
	if !force && compareVersions(currentVersion, latestRelease.TagName) >= 0 {
		fmt.Println("You are already using the latest version!")
		return nil
	}

	if checkOnly {
		fmt.Printf("Update available: %s -> %s\n", currentVersion, latestRelease.TagName)
		return nil
	}

	// Confirm upgrade
	fmt.Printf("Do you want to upgrade to version %s? [y/N] ", latestRelease.TagName)
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		fmt.Println("Upgrade cancelled")
		return nil
	}

	// Perform the upgrade
	fmt.Println("Starting upgrade...")
	return performUpgrade(latestRelease)
}

// getCurrentVersion returns the current CLI version
func getCurrentVersion() string {
	// TODO: Get from build info or version file
	return "1.0.0"
}

// fetchLatestRelease fetches the latest release from GitHub
func fetchLatestRelease(repoURL string) (*GitHubRelease, error) {
	// Extract owner and repo from GitHub URL
	parts := strings.Split(strings.TrimPrefix(repoURL, "https://github.com/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid GitHub repository URL")
	}
	owner, repo := parts[0], parts[1]

	// API URL for the latest release
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	// Use curl to fetch the release (avoid external dependencies)
	curlCmd := exec.Command("curl", "-s", "-H", "Accept: application/vnd.github.v3+json", apiURL)
	output, err := curlCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from GitHub API: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(output, &release); err != nil {
		return nil, fmt.Errorf("failed to parse release JSON: %w", err)
	}

	// Skip drafts and pre-releases
	if release.Draft || release.PreRelease {
		// Try to get the latest non-prerelease
		return fetchLatestStableRelease(repoURL)
	}

	return &release, nil
}

// fetchLatestStableRelease fetches the latest stable release
func fetchLatestStableRelease(repoURL string) (*GitHubRelease, error) {
	// Extract owner and repo from GitHub URL
	parts := strings.Split(strings.TrimPrefix(repoURL, "https://github.com/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid GitHub repository URL")
	}
	owner, repo := parts[0], parts[1]

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)

	curlCmd := exec.Command("curl", "-s", "-H", "Accept: application/vnd.github.v3+json", apiURL)
	output, err := curlCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases from GitHub API: %w", err)
	}

	var releases []GitHubRelease
	if err := json.Unmarshal(output, &releases); err != nil {
		return nil, fmt.Errorf("failed to parse releases JSON: %w", err)
	}

	// Find the first non-prerelease, non-draft release
	for _, release := range releases {
		if !release.Draft && !release.PreRelease {
			return &release, nil
		}
	}

	return nil, fmt.Errorf("no stable releases found")
}

// compareVersions compares two version strings
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	// Remove 'v' prefix if present
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	v1Parts := strings.Split(v1, ".")
	v2Parts := strings.Split(v2, ".")

	maxLen := len(v1Parts)
	if len(v2Parts) > maxLen {
		maxLen = len(v2Parts)
	}

	for i := 0; i < maxLen; i++ {
		v1Num := 0
		v2Num := 0

		if i < len(v1Parts) {
			fmt.Sscanf(v1Parts[i], "%d", &v1Num)
		}
		if i < len(v2Parts) {
			fmt.Sscanf(v2Parts[i], "%d", &v2Num)
		}

		if v1Num < v2Num {
			return -1
		} else if v1Num > v2Num {
			return 1
		}
	}

	return 0
}

// performUpgrade downloads and installs the new version
func performUpgrade(release *GitHubRelease) error {
	// Find the appropriate binary for this platform
	asset := findBinaryAsset(release)
	if asset == nil {
		return fmt.Errorf("no suitable binary found for your platform (%s/%s)", runtime.GOOS, runtime.GOARCH)
	}

	fmt.Printf("Downloading %s...\n", asset.Name)

	// Download to a temporary file
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, asset.Name)

	// Use curl to download
	curlCmd := exec.Command("curl", "-L", "-o", tempFile, asset.BrowserDownloadURL)
	if err := curlCmd.Run(); err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}

	// Make the binary executable
	if err := os.Chmod(tempFile, 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Get the current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable path: %w", err)
	}

	// Replace the binary
	fmt.Printf("Installing to %s...\n", currentExe)
	if err := os.Rename(tempFile, currentExe); err != nil {
		// Try copy if rename fails
		if err := copyFile(tempFile, currentExe); err != nil {
			return fmt.Errorf("failed to install new binary: %w", err)
		}
	}

	// Clean up
	os.Remove(tempFile)

	fmt.Printf("Successfully upgraded to %s!\n", release.TagName)
	fmt.Println("Please restart your terminal or run the command again to use the new version.")

	return nil
}

// findBinaryAsset finds the appropriate binary asset for the current platform
func findBinaryAsset(release *GitHubRelease) *struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
} {
	// Look for platform-specific binaries based on our naming convention
	switch {
	case runtime.GOOS == "linux" && runtime.GOARCH == "amd64":
		for _, asset := range release.Assets {
			if asset.Name == "vault-linux-amd64" {
				return &asset
			}
		}
	case runtime.GOOS == "linux" && runtime.GOARCH == "arm64":
		for _, asset := range release.Assets {
			if asset.Name == "vault-linux-arm64" {
				return &asset
			}
		}
	case runtime.GOOS == "darwin" && runtime.GOARCH == "amd64":
		for _, asset := range release.Assets {
			if asset.Name == "vault-darwin-amd64" {
				return &asset
			}
		}
	case runtime.GOOS == "darwin" && runtime.GOARCH == "arm64":
		for _, asset := range release.Assets {
			if asset.Name == "vault-darwin-arm64" {
				return &asset
			}
		}
	case runtime.GOOS == "windows" && runtime.GOARCH == "amd64":
		for _, asset := range release.Assets {
			if asset.Name == "vault-windows-amd64.exe" {
				return &asset
			}
		}
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(srcFile)
	return err
}
