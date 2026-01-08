package runtime

import (
	"os/exec"
	"strings"

	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// IsDockerInstalled checks if Docker is installed
func IsDockerInstalled() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}

// GetDockerInfo returns Docker information
func GetDockerInfo() *types.DockerInfo {
	info := &types.DockerInfo{
		Installed: IsDockerInstalled(),
		Running:   false,
		Version:   "",
		Info:      make(map[string]string),
	}

	if !info.Installed {
		return info
	}

	// Get Docker version
	if version, err := getDockerVersion(); err == nil {
		info.Version = version
	}

	// Check if Docker daemon is running
	if running, err := isDockerRunning(); err == nil {
		info.Running = running
	}

	// Get Docker system info
	if info.Running {
		if dockerInfo, err := getDockerSystemInfo(); err == nil {
			info.Info = dockerInfo
		}
	}

	return info
}

// getDockerVersion returns Docker version
func getDockerVersion() (string, error) {
	cmd := exec.Command("docker", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// isDockerRunning checks if Docker daemon is running
func isDockerRunning() (bool, error) {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	return err == nil, nil
}

// getDockerSystemInfo returns Docker system information
func getDockerSystemInfo() (map[string]string, error) {
	cmd := exec.Command("docker", "info", "--format", "{{json .}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// TODO: Parse JSON output into structured info
	// For now, return basic info
	info := make(map[string]string)
	info["raw"] = strings.TrimSpace(string(output))

	return info, nil
}

// CanRunDocker returns true if Docker commands can be executed
func CanRunDocker() bool {
	return IsDockerInstalled() && isDockerRunningNoError()
}

// isDockerRunningNoError checks if Docker is running without returning error
func isDockerRunningNoError() bool {
	running, _ := isDockerRunning()
	return running
}

// GetDockerImages returns list of Docker images
func GetDockerImages() ([]*types.DockerImage, error) {
	if !CanRunDocker() {
		return nil, nil
	}

	cmd := exec.Command("docker", "images", "--format", "{{.Repository}}:{{.Tag}}|{{.ID}}|{{.CreatedAt}}|{{.Size}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	_ = output // Avoid unused variable warning
	// TODO: Parse output into structured DockerImage slice
	// For now, return empty slice
	return []*types.DockerImage{}, nil
}

// GetDockerContainers returns list of Docker containers
func GetDockerContainers() ([]*types.DockerContainer, error) {
	if !CanRunDocker() {
		return nil, nil
	}

	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}|{{.ID}}|{{.Status}}|{{.Image}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	_ = output // Avoid unused variable warning
	// TODO: Parse output into structured DockerContainer slice
	// For now, return empty slice
	return []*types.DockerContainer{}, nil
}

// IsInDockerContainer returns true if running inside a Docker container
func IsInDockerContainer() bool {
	// Check for .dockerenv file
	if _, err := exec.LookPath("/.dockerenv"); err == nil {
		return true
	}

	// Check for container in cgroup
	// TODO: Implement proper container detection
	return false
}

// GetDockerNetworks returns list of Docker networks
func GetDockerNetworks() ([]*types.DockerNetwork, error) {
	if !CanRunDocker() {
		return nil, nil
	}

	cmd := exec.Command("docker", "network", "ls", "--format", "{{.Name}}|{{.ID}}|{{.Driver}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	_ = output // Avoid unused variable warning
	// TODO: Parse output into structured DockerNetwork slice
	// For now, return empty slice
	return []*types.DockerNetwork{}, nil
}
