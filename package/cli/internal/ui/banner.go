package ui

import (
	"fmt"
	"os"
)

// DisplayBanner shows the Aether Vault CLI banner
func DisplayBanner() error {
	// Check if color is enabled (for now, always use color)
	useColor := true

	// Banner text
	banner := "    ____             __     __        \n" +
		"   / __ \\___  ____ _/ /_ __/ /_  _____\n" +
		"  / /_/ / _ \\/ __ `/ __/ __/ __ \\/ ___/\n" +
		" / ____/  __/ /_/ / /_/ /_/ /_/ / /    \n" +
		"/_/    \\___/\\__,_/\\__/\\__/\\____/_/     \n" +
		"                                       \n" +
		"  Aether Vault CLI - Secure Secrets Management\n"

	// Display banner with color
	if useColor {
		fmt.Fprintf(os.Stdout, "%s%s%s\n", Cyan, banner, Reset)
		fmt.Fprintf(os.Stdout, "%sVersion: 1.0.0%s\n", Green, Reset)
		fmt.Fprintf(os.Stdout, "%sMode: Local%s\n", Yellow, Reset)
	} else {
		fmt.Fprint(os.Stdout, banner)
		fmt.Fprintln(os.Stdout, "Version: 1.0.0")
		fmt.Fprintln(os.Stdout, "Mode: Local")
	}

	fmt.Fprintln(os.Stdout)
	return nil
}
