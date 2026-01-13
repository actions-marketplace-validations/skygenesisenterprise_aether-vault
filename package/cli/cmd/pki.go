package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newPkiCommand creates the pki command
func newPkiCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pki",
		Short: "Interact with Vault's PKI Secrets Engine",
		Long:  `Manage PKI certificates and keys.`,
	}

	cmd.AddCommand(newPkiIssueCommand())
	cmd.AddCommand(newPkiSignCommand())
	cmd.AddCommand(newPkiRevokeCommand())
	cmd.AddCommand(newPkiListCommand())

	return cmd
}

// newPkiIssueCommand creates the pki issue command
func newPkiIssueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue [role] [common-name]",
		Short: "Issue a certificate",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Issued certificate for %s from role %s\n", args[1], args[0])
			return nil
		},
	}
	return cmd
}

// newPkiSignCommand creates the pki sign command
func newPkiSignCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign [role] [csr]",
		Short: "Sign a certificate signing request",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Signed CSR for role %s\n", args[0])
			return nil
		},
	}
	return cmd
}

// newPkiRevokeCommand creates the pki revoke command
func newPkiRevokeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revoke [serial]",
		Short: "Revoke a certificate",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Revoked certificate %s\n", args[0])
			return nil
		},
	}
	return cmd
}

// newPkiListCommand creates the pki list command
func newPkiListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [role]",
		Short: "List issued certificates",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			role := "default"
			if len(args) > 0 {
				role = args[0]
			}
			fmt.Printf("Listed certificates for role %s\n", role)
			return nil
		},
	}
	return cmd
}

// newPluginCommand creates the plugin command
func newPluginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Interact with Vault plugins and catalog",
		Long:  `Manage Vault plugins and plugin catalog.`,
	}

	cmd.AddCommand(newPluginListCommand())
	cmd.AddCommand(newPluginRegisterCommand())
	cmd.AddCommand(newPluginDeregisterCommand())

	return cmd
}

// newPluginListCommand creates the plugin list command
func newPluginListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Plugins:")
			fmt.Println("  database     builtin")
			fmt.Println("  aws          builtin")
			fmt.Println("  custom-auth  unknown")
			return nil
		},
	}
	return cmd
}

// newPluginRegisterCommand creates the plugin register command
func newPluginRegisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register [name] [path]",
		Short: "Register a plugin",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Registered plugin %s from %s\n", args[0], args[1])
			return nil
		},
	}
	return cmd
}

// newPluginDeregisterCommand creates the plugin deregister command
func newPluginDeregisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deregister [name]",
		Short: "Deregister a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Deregistered plugin %s\n", args[0])
			return nil
		},
	}
	return cmd
}
