package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newPatchCommand creates the patch command
func newPatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "patch [path] [data]",
		Short: "Patch data, configuration, and secrets",
		Long: `Patch data in Vault at the specified path using JSON merge patch format.

Examples:
  vault patch secret/data/my-app key=new-value
  vault patch secret/data/my-app @patch.json`,
		RunE: runPatchCommand,
	}

	return cmd
}

// runPatchCommand executes the patch command
func runPatchCommand(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("path argument is required")
	}

	path := args[0]

	// TODO: Implement patch logic
	fmt.Printf("Patching data at %s\n", path)
	return nil
}
