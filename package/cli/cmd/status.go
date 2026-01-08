package cmd

import (
	"fmt"
	"os"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/config"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/ui"
	"github.com/spf13/cobra"
)

// newStatusCommand creates the status command
func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show Vault status information",
		Long: `Display comprehensive status information including:
  - Current execution mode (local/cloud)
  - Configuration status
  - Runtime environment
  - Authentication state
  - Connection status`,
		RunE: runStatusCommand,
	}

	cmd.Flags().Bool("verbose", false, "Show detailed status information")
	cmd.Flags().String("format", "table", "Output format (json, yaml, table)")

	return cmd
}

// runStatusCommand executes the status command
func runStatusCommand(cmd *cobra.Command, args []string) error {
	_, _ = cmd.Flags().GetBool("verbose") // Read but not used for now
	format, _ := cmd.Flags().GetString("format")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to load configuration: %v\n", err)
		cfg = config.Defaults()
	}

	// Create context
	ctx, err := context.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}

	// Get status
	status, err := ctx.GetStatus()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}

	// Output based on format
	var outputFormat ui.OutputFormat
	switch format {
	case "json":
		outputFormat = ui.FormatJSON
	case "yaml":
		outputFormat = ui.FormatYAML
	default:
		outputFormat = ui.FormatTable
	}

	return ui.DisplayStatus(status, outputFormat)
}
