package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newAuditCommand creates the audit command
func newAuditCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Interact with audit devices",
		Long:  `Manage audit devices for logging Vault operations.`,
	}

	cmd.AddCommand(newAuditListCommand())
	cmd.AddCommand(newAuditEnableCommand())
	cmd.AddCommand(newAuditDisableCommand())

	return cmd
}

// newAuditListCommand creates the audit list command
func newAuditListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List enabled audit devices",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Audit devices:")
			fmt.Println("  file/    File-based audit device")
			fmt.Println("  syslog/  Syslog audit device")
			return nil
		},
	}
	return cmd
}

// newAuditEnableCommand creates the audit enable command
func newAuditEnableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable [type] [path]",
		Short: "Enable an audit device",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Enabled audit device %s at %s\n", args[0], args[1])
			return nil
		},
	}
	return cmd
}

// newAuditDisableCommand creates the audit disable command
func newAuditDisableCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable [path]",
		Short: "Disable an audit device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Disabled audit device at %s\n", args[0])
			return nil
		},
	}
	return cmd
}
