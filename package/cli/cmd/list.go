package cmd

import (
	"fmt"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/config"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/ui"
	"github.com/spf13/cobra"
)

// newListCommand creates the list command
func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [path]",
		Short: "List data or secrets",
		Long: `List data and secrets from Vault at the specified path.

Examples:
  vault list secret/
  vault list secret/data/
  vault list secret/ --format=json`,
		RunE: runListCommand,
	}

	cmd.Flags().String("format", "table", "Output format (json, yaml, table)")

	return cmd
}

// runListCommand executes the list command
func runListCommand(cmd *cobra.Command, args []string) error {
	path := ""
	if len(args) > 0 {
		path = args[0]
	}

	format, _ := cmd.Flags().GetString("format")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		cfg = config.Defaults()
	}

	// Create context
	ctx, err := context.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}

	// List data from Vault
	data, err := ctx.List(path)
	if err != nil {
		return fmt.Errorf("failed to list path '%s': %w", path, err)
	}

	// Display data based on format
	switch format {
	case "json":
		return ui.FormatOutput(data, ui.FormatJSON)
	case "yaml":
		return ui.FormatOutput(data, ui.FormatYAML)
	default:
		return ui.FormatOutput(data, ui.FormatTable)
	}
}
