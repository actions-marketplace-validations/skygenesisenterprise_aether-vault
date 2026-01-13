package cmd

import (
	"fmt"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/config"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
	"github.com/spf13/cobra"
)

// newUnwrapCommand creates the unwrap command
func newUnwrapCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unwrap [token]",
		Short: "Unwrap a wrapped secret",
		Long: `Unwrap a secret that was wrapped using the response-wrapping functionality.

Examples:
  vault unwrap s.1234567890abcdef
  vault unwrap -format=json s.1234567890abcdef`,
		RunE: runUnwrapCommand,
	}

	return cmd
}

// runUnwrapCommand executes the unwrap command
func runUnwrapCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("token argument is required")
	}

	token := args[0]

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

	// Unwrap the secret
	data, err := ctx.Unwrap(token)
	if err != nil {
		return fmt.Errorf("failed to unwrap token: %w", err)
	}

	// Display unwrapped data
	fmt.Printf("Successfully unwrapped token\n")
	fmt.Printf("Data: %+v\n", data)

	return nil
}
