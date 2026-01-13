package cmd

import (
	"fmt"
	"strings"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/config"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
	"github.com/spf13/cobra"
)

// newWriteCommand creates the write command
func newWriteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write [path] [data]",
		Short: "Write data, configuration, and secrets",
		Long: `Write data to Vault at the specified path.

Data can be provided as key=value pairs or as JSON/YAML.

Examples:
  vault write secret/data/my-app key=value username=admin
  vault write secret/data/my-app @data.json
  vault write secret/data/my-app - <<EOF
  {
    "username": "admin",
    "password": "secret123"
  }
  EOF`,
		RunE: runWriteCommand,
	}

	cmd.Flags().String("format", "json", "Input format (json, yaml, kv)")
	cmd.Flags().Bool("force", false, "Force write even if path exists")

	return cmd
}

// runWriteCommand executes the write command
func runWriteCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("path argument is required")
	}

	path := args[0]
	force, _ := cmd.Flags().GetBool("force")

	// Parse data from arguments
	var data map[string]interface{}
	if len(args) > 1 {
		data = parseKVData(args[1:])
	} else {
		// If no data provided, read from stdin
		return fmt.Errorf("data argument is required (use key=value pairs or stdin)")
	}

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

	// Write data to Vault
	result, err := ctx.Write(path, data, force)
	if err != nil {
		return fmt.Errorf("failed to write to path '%s': %w", path, err)
	}

	// Display success message
	fmt.Printf("Successfully wrote to %s\n", path)
	if result != nil {
		fmt.Printf("Version: %v\n", result["version"])
	}

	return nil
}

// parseKVData parses key=value pairs into a map
func parseKVData(args []string) map[string]interface{} {
	data := make(map[string]interface{})

	for _, arg := range args {
		if strings.HasPrefix(arg, "@") {
			// TODO: Handle file input
			continue
		}

		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			data[parts[0]] = parts[1]
		}
	}

	return data
}
