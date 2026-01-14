package cmd

import (
	"fmt"
	"os"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/config"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/ui"
	"github.com/spf13/cobra"
)

// NewRootCommand creates the root vault command
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vault",
		Short: "Aether Vault CLI - DevOps & Security Tool",
		Long: `Aether Vault is a comprehensive DevOps and security tool for secret management.

Use 'vault help' to see available commands or 'vault <command> --help' for command-specific help.

Available modes:
  - Local: Offline secret storage and management
  - Cloud: Connected to Aether Vault cloud services

Quick start:
  vault init     Initialize local environment
  vault status   Check current status
  vault login    Connect to cloud services`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check for version flag
			if version, _ := cmd.Flags().GetBool("version"); version {
				return runVersionCommand(cmd, args)
			}

			// If no arguments, show help and status
			if len(args) == 0 {
				return runRootCommand(cmd)
			}
			return fmt.Errorf("unknown command '%s'", args[0])
		},
	}

	// Global flags
	cmd.PersistentFlags().String("format", "table", "Output format (json, yaml, table)")
	cmd.PersistentFlags().Bool("verbose", false, "Enable verbose output")
	cmd.PersistentFlags().String("config", "", "Config file path (default is ~/.aether/vault/config.yaml)")

	// Version flag
	cmd.Flags().Bool("version", false, "Display version information")

	// Add subcommands - Common commands
	cmd.AddCommand(newReadCommand())
	cmd.AddCommand(newWriteCommand())
	cmd.AddCommand(newDeleteCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newLoginCommand())
	cmd.AddCommand(newAgentCommand())
	cmd.AddCommand(newServerCommand())
	cmd.AddCommand(newStatusCommand())
	cmd.AddCommand(newUnwrapCommand())
	cmd.AddCommand(newUpgradeCommand())

	// Add subcommands - Other commands
	cmd.AddCommand(newAuditCommand())
	cmd.AddCommand(newDebugCommand())
	cmd.AddCommand(newEventsCommand())
	cmd.AddCommand(newHcpCommand())
	cmd.AddCommand(newKvCommand())
	cmd.AddCommand(newLeaseCommand())
	cmd.AddCommand(newMonitorCommand())
	cmd.AddCommand(newNamespaceCommand())
	cmd.AddCommand(newOperatorCommand())
	cmd.AddCommand(newPatchCommand())
	cmd.AddCommand(newPathHelpCommand())
	cmd.AddCommand(newPkiCommand())
	cmd.AddCommand(newPluginCommand())
	cmd.AddCommand(newPolicyCommand())
	cmd.AddCommand(newPrintCommand())
	cmd.AddCommand(newProxyCommand())
	cmd.AddCommand(newSecretsCommand())
	cmd.AddCommand(newSshCommand())
	cmd.AddCommand(newTokenCommand())
	cmd.AddCommand(newTransformCommand())
	cmd.AddCommand(newTransitCommand())
	cmd.AddCommand(newVersionHistoryCommand())

	// Add existing commands
	cmd.AddCommand(newVersionCommand())
	cmd.AddCommand(newInitCommand())
	cmd.AddCommand(newHelpCommand())
	cmd.AddCommand(newCapabilityCommand())

	// Add new commands
	cmd.AddCommand(newPasswordCommand())
	cmd.AddCommand(newTOTPCommand())
	cmd.AddCommand(newLogsCommand())
	cmd.AddCommand(newShellCommand())
	cmd.AddCommand(newDataCommand())

	return cmd
}

// runRootCommand executes the root command behavior
func runRootCommand(cmd *cobra.Command) error {
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

	// Display welcome banner
	if err := ui.DisplayBanner(); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to display banner: %v\n", err)
	}

	// Show current status
	status, err := ctx.GetStatus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to get status: %v\n", err)
	} else {
		ui.DisplayStatus(status, ui.FormatTable)
	}

	// Show available commands
	fmt.Println("\nAvailable commands:")
	fmt.Println("  init      Initialize local Vault environment")
	fmt.Println("  login     Connect to Aether Vault cloud")
	fmt.Println("  status    Show current Vault status")
	fmt.Println("  version   Display CLI version information")
	fmt.Println("  help      Show help for commands")

	fmt.Println("\nUse 'vault <command> --help' for detailed help.")

	return nil
}
