package cmd

import (
	"fmt"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/config"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
	"github.com/spf13/cobra"
)

// newDeleteCommand creates the delete command
func newDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [path]",
		Short: "Delete secrets and configuration",
		Long: `Delete data from Vault at the specified path.

Examples:
  vault delete secret/data/my-app
  vault delete secret/data/my-app --versions=2
  vault delete secret/data/my-app --force`,
		RunE: runDeleteCommand,
	}

	cmd.Flags().String("versions", "", "Specify versions to delete")
	cmd.Flags().Bool("force", false, "Force delete without confirmation")
	cmd.Flags().Bool("recursive", false, "Delete recursively")

	return cmd
}

// runDeleteCommand executes the delete command
func runDeleteCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("path argument is required")
	}

	path := args[0]
	versions, _ := cmd.Flags().GetString("versions")
	force, _ := cmd.Flags().GetBool("force")
	recursive, _ := cmd.Flags().GetBool("recursive")

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

	// Confirmation prompt if not forced
	if !force {
		fmt.Printf("Are you sure you want to delete '%s'? [y/N] ", path)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Delete cancelled")
			return nil
		}
	}

	// Delete data from Vault
	result, err := ctx.Delete(path, versions, recursive)
	if err != nil {
		return fmt.Errorf("failed to delete path '%s': %w", path, err)
	}

	// Display success message
	fmt.Printf("Successfully deleted %s\n", path)
	if result != nil {
		if versions, ok := result["versions"]; ok {
			fmt.Printf("Deleted versions: %v\n", versions)
		}
	}

	return nil
}
