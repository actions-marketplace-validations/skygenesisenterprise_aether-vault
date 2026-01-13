package cmd

import (
	"fmt"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/config"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/ui"
	"github.com/spf13/cobra"
)

// newReadCommand creates the read command
func newReadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read [path]",
		Short: "Read data and retrieves secrets",
		Long: `Read data from Vault and retrieve secrets from the specified path.

Examples:
  vault read secret/data/my-app
  vault read -field=password secret/data/my-app
  vault read -format=json secret/data/my-app`,
		RunE: runReadCommand,
	}

	cmd.Flags().String("field", "", "Return only the specified field")
	cmd.Flags().String("format", "table", "Output format (json, yaml, table)")

	return cmd
}

// runReadCommand executes the read command
func runReadCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("path argument is required")
	}

	path := args[0]
	field, _ := cmd.Flags().GetString("field")
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

	// Read data from Vault
	data, err := ctx.Read(path)
	if err != nil {
		return fmt.Errorf("failed to read from path '%s': %w", path, err)
	}

	// Handle field extraction
	if field != "" {
		if fieldValue, exists := data[field]; exists {
			fmt.Println(fieldValue)
			return nil
		}
		return fmt.Errorf("field '%s' not found in response", field)
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
