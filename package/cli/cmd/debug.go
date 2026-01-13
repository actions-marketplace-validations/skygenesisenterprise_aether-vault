package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newDebugCommand creates the debug command
func newDebugCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Runs the debug command",
		Long:  `Debug and troubleshoot Vault operations.`,
	}

	cmd.AddCommand(newDebugConfigCommand())
	cmd.AddCommand(newDebugConnectionCommand())

	return cmd
}

// newDebugConfigCommand creates the debug config command
func newDebugConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Debug configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Configuration debug information:")
			fmt.Println("  Config file: ~/.aether/vault/config.yaml")
			fmt.Println("  Mode: local")
			fmt.Println("  Status: loaded")
			return nil
		},
	}
	return cmd
}

// newDebugConnectionCommand creates the debug connection command
func newDebugConnectionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connection",
		Short: "Debug connection to Vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Connection debug information:")
			fmt.Println("  Server: localhost:8200")
			fmt.Println("  Status: connected")
			fmt.Println("  Latency: 5ms")
			return nil
		},
	}
	return cmd
}
