package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
)

// OutputFormat represents the output format
type OutputFormat string

const (
	FormatTable OutputFormat = "table"
	FormatJSON  OutputFormat = "json"
	FormatYAML  OutputFormat = "yaml"
)

// DisplayStatus shows the status information in the specified format
func DisplayStatus(status *context.Status, format OutputFormat) error {
	switch format {
	case FormatJSON:
		return displayStatusJSON(status)
	case FormatYAML:
		return displayStatusYAML(status)
	case FormatTable:
		return displayStatusTable(status)
	default:
		return displayStatusTable(status)
	}
}

// displayStatusJSON displays status in JSON format
func displayStatusJSON(status *context.Status) error {
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal status to JSON: %w", err)
	}
	fmt.Fprintln(os.Stdout, string(data))
	return nil
}

// displayStatusYAML displays status in YAML format
func displayStatusYAML(status *context.Status) error {
	data, err := yaml.Marshal(status)
	if err != nil {
		return fmt.Errorf("failed to marshal status to YAML: %w", err)
	}
	fmt.Fprint(os.Stdout, string(data))
	return nil
}

// displayStatusTable displays status in table format
func displayStatusTable(status *context.Status) error {
	// Header
	fmt.Fprintf(os.Stdout, "%sAether Vault Status%s\n", Bold, Reset)
	fmt.Fprintln(os.Stdout, strings.Repeat("─", 50))

	// Mode
	modeColor := Green
	if status.Mode == "cloud" {
		modeColor = Blue
	}
	fmt.Fprintf(os.Stdout, "%sMode:%s %s%s%s\n", Bold, Reset, modeColor, status.Mode, Reset)

	// Configuration status
	configColor := Green
	if status.ConfigStatus != "loaded" {
		configColor = Red
	}
	fmt.Fprintf(os.Stdout, "%sConfiguration:%s %s%s%s\n", Bold, Reset, configColor, status.ConfigStatus, Reset)

	// Authentication status
	authColor := Red
	if status.AuthStatus == "authenticated" {
		authColor = Green
	}
	fmt.Fprintf(os.Stdout, "%sAuthentication:%s %s%s%s\n", Bold, Reset, authColor, status.AuthStatus, Reset)

	// Runtime information
	if status.Runtime != nil {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "%sRuntime Information:%s\n", Bold, Reset)
		fmt.Fprintf(os.Stdout, "  OS: %s%s%s\n", Cyan, status.Runtime.OS, Reset)
		fmt.Fprintf(os.Stdout, "  Architecture: %s%s%s\n", Cyan, status.Runtime.Arch, Reset)
		fmt.Fprintf(os.Stdout, "  Go Version: %s%s%s\n", Cyan, status.Runtime.GoVersion, Reset)
		fmt.Fprintf(os.Stdout, "  CLI Version: %s%s%s\n", Cyan, status.Runtime.Version, Reset)
	}

	// Details
	if len(status.Details) > 0 {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintf(os.Stdout, "%sDetails:%s\n", Bold, Reset)
		for key, value := range status.Details {
			fmt.Fprintf(os.Stdout, "  %s%s:%s %v\n", Yellow, key, Reset, value)
		}
	}

	// Last updated
	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "%sLast Updated:%s %s%s%s\n", Bold, Reset, Gray, status.LastUpdated.Format(time.RFC3339), Reset)

	return nil
}

// FormatOutput formats data in the specified format
func FormatOutput(data interface{}, format OutputFormat) error {
	switch format {
	case FormatJSON:
		return formatJSON(data)
	case FormatYAML:
		return formatYAML(data)
	case FormatTable:
		return formatTable(data)
	default:
		return formatTable(data)
	}
}

// formatJSON formats data as JSON
func formatJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}
	fmt.Fprintln(os.Stdout, string(jsonData))
	return nil
}

// formatYAML formats data as YAML
func formatYAML(data interface{}) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data to YAML: %w", err)
	}
	fmt.Fprint(os.Stdout, string(yamlData))
	return nil
}

// formatTable formats data as a simple table
func formatTable(data interface{}) error {
	// For now, just use JSON format for table
	// TODO: Implement proper table formatting based on data type
	return formatJSON(data)
}
