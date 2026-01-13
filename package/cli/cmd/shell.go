package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/config"
	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/context"
	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
	"github.com/spf13/cobra"
)

// ShellCommand represents the shell command state
type ShellCommand struct {
	cmd      *cobra.Command
	cfg      *types.Config
	ctx      *context.Context
	replMode bool
	script   string
	history  []string
	vars     map[string]string
}

// DSLCommand represents a DSL command
type DSLCommand struct {
	name        string
	description string
	handler     func(*ShellCommand, []string) error
}

// newShellCommand creates the shell command
func newShellCommand() *cobra.Command {
	shellCmd := &cobra.Command{
		Use:   "shell",
		Short: "Interactive shell with DSL for automation",
		Long: `Start an interactive shell with a Domain Specific Language (DSL) for automating
Vault operations. The shell supports both interactive REPL mode and script execution.

DSL Features:
  - Variable assignment and substitution
  - Conditional statements
  - Loop constructs
  - Command chaining
  - Built-in Vault operations

Usage:
  vault shell              Start interactive shell
  vault shell script.vault  Execute script file
  vault shell -c "command"  Execute single command`,
		RunE: runShellCommand,
	}

	shellCmd.Flags().BoolP("interactive", "i", true, "Start interactive REPL mode")
	shellCmd.Flags().StringP("command", "c", "", "Execute single command and exit")
	shellCmd.Flags().StringP("script", "s", "", "Execute script file")

	return shellCmd
}

// runShellCommand executes the shell command
func runShellCommand(cmd *cobra.Command, args []string) error {
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

	// Parse flags
	interactive, _ := cmd.Flags().GetBool("interactive")
	command, _ := cmd.Flags().GetString("command")
	script, _ := cmd.Flags().GetString("script")

	// Create shell instance
	shell := &ShellCommand{
		cmd:     cmd,
		cfg:     cfg,
		ctx:     ctx,
		vars:    make(map[string]string),
		history: make([]string, 0),
	}

	// Handle script execution
	if script != "" {
		return shell.executeScript(script)
	}

	// Handle single command execution
	if command != "" {
		return shell.executeCommand(command)
	}

	// Start interactive shell
	if interactive {
		return shell.startREPL()
	}

	return nil
}

// startREPL starts the interactive shell
func (s *ShellCommand) startREPL() error {
	fmt.Println("Aether Vault Shell - DSL Automation Mode")
	fmt.Println("Type 'help' for available commands, 'exit' to quit")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	for {
		// Display prompt
		prompt := s.buildPrompt()
		fmt.Print(prompt)

		// Read input
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		// Trim whitespace
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Add to history
		s.history = append(s.history, line)

		// Handle built-in commands
		if s.handleBuiltinCommand(line) {
			continue
		}

		// Execute DSL command
		if err := s.executeCommand(line); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

// executeScript executes a script file
func (s *ShellCommand) executeScript(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read script file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Add to history
		s.history = append(s.history, line)

		// Execute command
		if err := s.executeCommand(line); err != nil {
			return fmt.Errorf("script error at line %d: %w", i+1, err)
		}
	}

	return nil
}

// executeCommand executes a DSL command
func (s *ShellCommand) executeCommand(command string) error {
	// Parse command
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	// Handle variable assignment
	if strings.Contains(parts[0], "=") {
		return s.handleVariableAssignment(command)
	}

	// Handle DSL commands
	dslCommands := s.getDSLCommands()
	for _, dslCmd := range dslCommands {
		if parts[0] == dslCmd.name {
			return dslCmd.handler(s, parts[1:])
		}
	}

	// Handle Vault commands
	return s.handleVaultCommand(parts)
}

// handleBuiltinCommand handles built-in REPL commands
func (s *ShellCommand) handleBuiltinCommand(line string) bool {
	switch line {
	case "exit", "quit":
		fmt.Println("Goodbye!")
		os.Exit(0)
		return true
	case "help":
		s.showHelp()
		return true
	case "clear":
		s.clearScreen()
		return true
	case "history":
		s.showHistory()
		return true
	case "vars":
		s.showVariables()
		return true
	}
	return false
}

// handleVariableAssignment handles variable assignment
func (s *ShellCommand) handleVariableAssignment(command string) error {
	parts := strings.SplitN(command, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid variable assignment")
	}

	varName := strings.TrimSpace(parts[0])
	varValue := strings.TrimSpace(parts[1])

	// Remove quotes if present
	if strings.HasPrefix(varValue, "\"") && strings.HasSuffix(varValue, "\"") {
		varValue = varValue[1 : len(varValue)-1]
	}

	s.vars[varName] = varValue
	fmt.Printf("Set %s = %s\n", varName, varValue)
	return nil
}

// handleVaultCommand handles Vault commands
func (s *ShellCommand) handleVaultCommand(parts []string) error {
	// Substitute variables
	substituted := s.substituteVariables(strings.Join(parts, " "))
	parts = strings.Fields(substituted)

	// Execute Vault command
	vaultCmd := s.cmd.Root()
	args := append([]string{"vault"}, parts...)
	vaultCmd.SetArgs(args)

	return vaultCmd.Execute()
}

// substituteVariables substitutes variables in command
func (s *ShellCommand) substituteVariables(command string) string {
	result := command
	for varName, varValue := range s.vars {
		placeholder := fmt.Sprintf("${%s}", varName)
		result = strings.ReplaceAll(result, placeholder, varValue)
	}
	return result
}

// getDSLCommands returns available DSL commands
func (s *ShellCommand) getDSLCommands() []DSLCommand {
	return []DSLCommand{
		{
			name:        "echo",
			description: "Print a message",
			handler:     dslEcho,
		},
		{
			name:        "sleep",
			description: "Pause execution",
			handler:     dslSleep,
		},
		{
			name:        "if",
			description: "Conditional execution",
			handler:     dslIf,
		},
		{
			name:        "for",
			description: "Loop execution",
			handler:     dslFor,
		},
		{
			name:        "set",
			description: "Set variable",
			handler:     dslSet,
		},
		{
			name:        "get",
			description: "Get variable",
			handler:     dslGet,
		},
		{
			name:        "wait",
			description: "Wait for condition",
			handler:     dslWait,
		},
	}
}

// DSL Command Handlers

func dslEcho(s *ShellCommand, args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}

func dslSleep(s *ShellCommand, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("sleep requires duration argument")
	}

	duration, err := time.ParseDuration(args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	time.Sleep(duration)
	return nil
}

func dslIf(s *ShellCommand, args []string) error {
	// Simple if implementation: if <condition> then <command>
	if len(args) < 3 {
		return fmt.Errorf("if requires condition and command")
	}

	condition := args[0]
	if args[1] != "then" {
		return fmt.Errorf("if syntax: if <condition> then <command>")
	}

	command := strings.Join(args[2:], " ")

	// Evaluate condition (simple string comparison)
	if s.evaluateCondition(condition) {
		return s.executeCommand(command)
	}

	return nil
}

func dslFor(s *ShellCommand, args []string) error {
	// Simple for implementation: for <var> in <values> do <command>
	if len(args) < 5 {
		return fmt.Errorf("for requires variable, values, and command")
	}

	varName := args[0]
	if args[1] != "in" || args[3] != "do" {
		return fmt.Errorf("for syntax: for <var> in <values> do <command>")
	}

	values := strings.Split(args[2], ",")
	command := strings.Join(args[4:], " ")

	// Store original variable value
	originalValue := s.vars[varName]

	// Execute loop
	for _, value := range values {
		s.vars[varName] = strings.TrimSpace(value)
		if err := s.executeCommand(command); err != nil {
			return err
		}
	}

	// Restore original variable value
	if originalValue != "" {
		s.vars[varName] = originalValue
	} else {
		delete(s.vars, varName)
	}

	return nil
}

func dslSet(s *ShellCommand, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("set requires variable and value")
	}

	s.vars[args[0]] = args[1]
	return nil
}

func dslGet(s *ShellCommand, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("get requires variable name")
	}

	value, exists := s.vars[args[0]]
	if !exists {
		return fmt.Errorf("variable '%s' not found", args[0])
	}

	fmt.Println(value)
	return nil
}

func dslWait(s *ShellCommand, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("wait requires condition and timeout")
	}

	condition := args[0]
	timeout, err := time.ParseDuration(args[1])
	if err != nil {
		return fmt.Errorf("invalid timeout: %w", err)
	}

	start := time.Now()
	for time.Since(start) < timeout {
		if s.evaluateCondition(condition) {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("timeout waiting for condition: %s", condition)
}

// Helper Methods

func (s *ShellCommand) evaluateCondition(condition string) bool {
	// Simple condition evaluation
	condition = s.substituteVariables(condition)

	// Handle equality checks
	if strings.Contains(condition, "==") {
		parts := strings.Split(condition, "==")
		if len(parts) == 2 {
			return strings.TrimSpace(parts[0]) == strings.TrimSpace(parts[1])
		}
	}

	// Handle inequality checks
	if strings.Contains(condition, "!=") {
		parts := strings.Split(condition, "!=")
		if len(parts) == 2 {
			return strings.TrimSpace(parts[0]) != strings.TrimSpace(parts[1])
		}
	}

	// Default: treat non-empty string as true
	return condition != ""
}

func (s *ShellCommand) buildPrompt() string {
	if status, err := s.ctx.GetStatus(); err == nil {
		return fmt.Sprintf("vault:%s> ", status.Mode)
	}
	return "vault> "
}

func (s *ShellCommand) showHelp() {
	fmt.Println("Available DSL Commands:")
	dslCommands := s.getDSLCommands()
	for _, cmd := range dslCommands {
		fmt.Printf("  %-10s %s\n", cmd.name, cmd.description)
	}

	fmt.Println("\nBuilt-in Commands:")
	fmt.Println("  help     Show this help")
	fmt.Println("  clear    Clear screen")
	fmt.Println("  history  Show command history")
	fmt.Println("  vars     Show variables")
	fmt.Println("  exit     Exit shell")

	fmt.Println("\nVariable Assignment:")
	fmt.Println("  name=value    Set variable")
	fmt.Println("  ${name}       Use variable")

	fmt.Println("\nVault Commands:")
	fmt.Println("  Any vault command can be used directly")
	fmt.Println("  Example: read secret/mykey")
}

func (s *ShellCommand) clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (s *ShellCommand) showHistory() {
	fmt.Println("Command History:")
	for i, cmd := range s.history {
		fmt.Printf("  %d: %s\n", i+1, cmd)
	}
}

func (s *ShellCommand) showVariables() {
	fmt.Println("Variables:")
	for name, value := range s.vars {
		fmt.Printf("  %s = %s\n", name, value)
	}
}
