package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/internal/ui"
	"github.com/spf13/cobra"
)

// newLogsCommand creates the logs command
func newLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "View and monitor vault logs",
		Long: `View and monitor vault logs in real-time with filtering options.

This command provides:
  - Real-time log monitoring (tail -f functionality)
  - Log filtering by level, component, or time
  - Log search capabilities
  - Multiple output formats
  - Log export functionality

Log levels: DEBUG, INFO, WARN, ERROR, FATAL
Components: auth, storage, sync, api, cli, all

Examples:
  vault logs                    # Show all logs
  vault logs -f                 # Follow logs in real-time
  vault logs --level error       # Show only error logs
  vault logs --component auth    # Show auth component logs
  vault logs --since "1h"       # Show logs from last hour
  vault logs --search "login"     # Search for "login" in logs`,
		RunE: runLogsCommand,
	}

	// Filtering flags
	cmd.Flags().BoolP("follow", "f", false, "Follow log output in real-time")
	cmd.Flags().StringP("level", "l", "", "Log level to filter (debug, info, warn, error, fatal)")
	cmd.Flags().StringP("component", "c", "", "Component to filter (auth, storage, sync, api, cli, all)")
	cmd.Flags().StringP("since", "s", "", "Show logs since timestamp (e.g., 1h, 30m, 2023-01-01T10:00:00)")
	cmd.Flags().StringP("until", "u", "", "Show logs until timestamp")
	cmd.Flags().StringP("search", "q", "", "Search for specific text in logs")
	cmd.Flags().IntP("limit", "n", 100, "Limit number of log entries to show")
	cmd.Flags().BoolP("reverse", "r", false, "Show logs in reverse order (newest first)")

	// Output flags
	cmd.Flags().StringP("output", "o", "table", "Output format (table, json, raw)")
	cmd.Flags().BoolP("no-color", "", false, "Disable colored output")
	cmd.Flags().StringP("export", "e", "", "Export logs to file")
	cmd.Flags().BoolP("tail", "t", false, "Show last N lines and follow (like tail -f)")

	return cmd
}

// runLogsCommand executes the logs command
func runLogsCommand(cmd *cobra.Command, args []string) error {
	follow, _ := cmd.Flags().GetBool("follow")
	level, _ := cmd.Flags().GetString("level")
	component, _ := cmd.Flags().GetString("component")
	since, _ := cmd.Flags().GetString("since")
	until, _ := cmd.Flags().GetString("until")
	search, _ := cmd.Flags().GetString("search")
	limit, _ := cmd.Flags().GetInt("limit")
	reverse, _ := cmd.Flags().GetBool("reverse")
	output, _ := cmd.Flags().GetString("output")
	noColor, _ := cmd.Flags().GetBool("no-color")
	exportFile, _ := cmd.Flags().GetString("export")
	tail, _ := cmd.Flags().GetBool("tail")

	// Create log viewer options
	options := &LogViewerOptions{
		Follow:    follow || tail,
		Level:     strings.ToUpper(level),
		Component: strings.ToLower(component),
		Since:     since,
		Until:     until,
		Search:    search,
		Limit:     limit,
		Reverse:   reverse,
		Output:    output,
		NoColor:   noColor,
		Export:    exportFile,
	}

	// Create and start log viewer
	viewer := NewLogViewer(options)
	return viewer.Start()
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Component string                 `json:"component"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Source    string                 `json:"source,omitempty"`
}

// LogViewerOptions represents options for log viewing
type LogViewerOptions struct {
	Follow    bool
	Level     string
	Component string
	Since     string
	Until     string
	Search    string
	Limit     int
	Reverse   bool
	Output    string
	NoColor   bool
	Export    string
}

// LogViewer handles log viewing and monitoring
type LogViewer struct {
	options *LogViewerOptions
	entries []*LogEntry
	running bool
	logFile *os.File
}

// NewLogViewer creates a new log viewer
func NewLogViewer(options *LogViewerOptions) *LogViewer {
	return &LogViewer{
		options: options,
		entries: make([]*LogEntry, 0),
		running: true,
	}
}

// Start starts the log viewer
func (v *LogViewer) Start() error {
	// Load existing logs
	if err := v.loadLogs(); err != nil {
		return fmt.Errorf("failed to load logs: %w", err)
	}

	// Display initial logs
	if err := v.displayLogs(); err != nil {
		return err
	}

	// Start following if requested
	if v.options.Follow {
		return v.followLogs()
	}

	return nil
}

// loadLogs loads log entries from various sources
func (v *LogViewer) loadLogs() error {
	// For now, generate sample logs
	// In a real implementation, this would read from actual log files
	v.generateSampleLogs()

	// Apply filters
	v.entries = v.filterLogs()

	return nil
}

// generateSampleLogs generates sample log entries for demonstration
func (v *LogViewer) generateSampleLogs() {
	now := time.Now()
	sampleLogs := []*LogEntry{
		{
			Timestamp: now.Add(-5 * time.Minute),
			Level:     "INFO",
			Component: "auth",
			Message:   "User authentication successful",
			Fields:    map[string]interface{}{"user": "alice@example.com", "method": "oauth"},
			Source:    "auth.service",
		},
		{
			Timestamp: now.Add(-4 * time.Minute),
			Level:     "DEBUG",
			Component: "storage",
			Message:   "Cache hit for secret key",
			Fields:    map[string]interface{}{"key": "github_token", "cache": "redis"},
			Source:    "storage.local",
		},
		{
			Timestamp: now.Add(-3 * time.Minute),
			Level:     "WARN",
			Component: "sync",
			Message:   "Sync conflict detected",
			Fields:    map[string]interface{}{"conflict_id": "12345", "resolved": false},
			Source:    "sync.service",
		},
		{
			Timestamp: now.Add(-2 * time.Minute),
			Level:     "ERROR",
			Component: "api",
			Message:   "Failed to connect to external service",
			Fields:    map[string]interface{}{"service": "github", "error": "timeout"},
			Source:    "api.gateway",
		},
		{
			Timestamp: now.Add(-1 * time.Minute),
			Level:     "INFO",
			Component: "cli",
			Message:   "Vault initialization completed",
			Fields:    map[string]interface{}{"version": "1.0.0", "duration": "2.3s"},
			Source:    "cli.main",
		},
		{
			Timestamp: now,
			Level:     "INFO",
			Component: "totp",
			Message:   "TOTP code generated successfully",
			Fields:    map[string]interface{}{"issuer": "GitHub", "account": "alice@example.com"},
			Source:    "totp.service",
		},
	}

	v.entries = sampleLogs
}

// filterLogs applies filters to log entries
func (v *LogViewer) filterLogs() []*LogEntry {
	filtered := make([]*LogEntry, 0, len(v.entries))

	for _, entry := range v.entries {
		if v.shouldIncludeEntry(entry) {
			filtered = append(filtered, entry)
		}
	}

	// Apply limit
	if v.options.Limit > 0 && len(filtered) > v.options.Limit {
		filtered = filtered[:v.options.Limit]
	}

	// Apply reverse order if requested
	if v.options.Reverse {
		for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
			filtered[i], filtered[j] = filtered[j], filtered[i]
		}
	}

	return filtered
}

// shouldIncludeEntry checks if a log entry should be included based on filters
func (v *LogViewer) shouldIncludeEntry(entry *LogEntry) bool {
	// Level filter
	if v.options.Level != "" && entry.Level != v.options.Level {
		return false
	}

	// Component filter
	if v.options.Component != "" && v.options.Component != "all" && entry.Component != v.options.Component {
		return false
	}

	// Search filter
	if v.options.Search != "" {
		if !strings.Contains(strings.ToLower(entry.Message), strings.ToLower(v.options.Search)) {
			return false
		}
	}

	return true
}

// displayLogs displays log entries based on output format
func (v *LogViewer) displayLogs() error {
	switch v.options.Output {
	case "json":
		return v.displayJSON()
	case "raw":
		return v.displayRaw()
	default:
		return v.displayTable()
	}
}

// displayTable displays logs in table format
func (v *LogViewer) displayTable() error {
	if len(v.entries) == 0 {
		fmt.Printf("%sNo log entries found matching the criteria.%s\n", ui.Yellow, ui.Reset)
		return nil
	}

	fmt.Printf("%s%-20s %-8s %-12s %-40s %s%s\n",
		ui.Bold, "TIMESTAMP", "LEVEL", "COMPONENT", "MESSAGE", "SOURCE", ui.Reset)
	fmt.Println(strings.Repeat("-", 90))

	for _, entry := range v.entries {
		timestamp := entry.Timestamp.Format("15:04:05")

		// Color by level
		levelColor := ui.Reset
		switch entry.Level {
		case "DEBUG":
			levelColor = ui.Gray
		case "INFO":
			levelColor = ui.Blue
		case "WARN":
			levelColor = ui.Yellow
		case "ERROR", "FATAL":
			levelColor = ui.Red
		}

		// Truncate message if too long
		message := entry.Message
		if len(message) > 40 {
			message = message[:37] + "..."
		}

		if v.options.NoColor {
			fmt.Printf("%-20s %-8s %-12s %-40s %s\n",
				timestamp, entry.Level, entry.Component, message, entry.Source)
		} else {
			fmt.Printf("%-20s %s%-8s%s %-12s %-40s %s\n",
				timestamp, levelColor, entry.Level, ui.Reset, entry.Component, message, entry.Source)
		}
	}

	return nil
}

// displayJSON displays logs in JSON format
func (v *LogViewer) displayJSON() error {
	for i, entry := range v.entries {
		if i > 0 {
			fmt.Println(",")
		}
		fmt.Printf("  {")
		fmt.Printf("\n    \"timestamp\": \"%s\",", entry.Timestamp.Format(time.RFC3339))
		fmt.Printf("\n    \"level\": \"%s\",", entry.Level)
		fmt.Printf("\n    \"component\": \"%s\",", entry.Component)
		fmt.Printf("\n    \"message\": \"%s\",", entry.Message)
		fmt.Printf("\n    \"source\": \"%s\"", entry.Source)
		if entry.Fields != nil {
			fmt.Printf(",\n    \"fields\": {")
			first := true
			for k, v := range entry.Fields {
				if !first {
					fmt.Println(",")
				}
				fmt.Printf("\n      \"%s\": \"%v\"", k, v)
				first = false
			}
			fmt.Printf("\n    }")
		}
		fmt.Printf("\n  }")
	}
	fmt.Println()
	return nil
}

// displayRaw displays logs in raw format
func (v *LogViewer) displayRaw() error {
	for _, entry := range v.entries {
		fmt.Printf("%s [%s] %s: %s\n",
			entry.Timestamp.Format(time.RFC3339), entry.Level, entry.Component, entry.Message)
	}
	return nil
}

// followLogs follows logs in real-time
func (v *LogViewer) followLogs() error {
	fmt.Printf("\n%sFollowing logs... (Press Ctrl+C to stop)%s\n", ui.Dim, ui.Reset)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for v.running {
		select {
		case <-time.After(1 * time.Second):
			// Simulate new log entry
			newEntry := &LogEntry{
				Timestamp: time.Now(),
				Level:     "INFO",
				Component: "cli",
				Message:   "Real-time log monitoring active",
				Source:    "logs.follow",
			}

			// Check if entry should be included
			if v.shouldIncludeEntry(newEntry) {
				if v.options.Output == "table" {
					timestamp := newEntry.Timestamp.Format("15:04:05")
					fmt.Printf("%s%s %-8s %-12s %-40s %s%s\n",
						ui.Green, timestamp, newEntry.Level, newEntry.Component, newEntry.Message, newEntry.Source, ui.Reset)
				} else {
					v.displayEntry(newEntry)
				}
			}
		}
	}

	return nil
}

// displayEntry displays a single log entry
func (v *LogViewer) displayEntry(entry *LogEntry) error {
	switch v.options.Output {
	case "json":
		// JSON implementation for single entry
		fmt.Printf("{\"timestamp\":\"%s\",\"level\":\"%s\",\"component\":\"%s\",\"message\":\"%s\",\"source\":\"%s\"}\n",
			entry.Timestamp.Format(time.RFC3339), entry.Level, entry.Component, entry.Message, entry.Source)
	case "raw":
		fmt.Printf("%s [%s] %s: %s\n",
			entry.Timestamp.Format(time.RFC3339), entry.Level, entry.Component, entry.Message)
	}
	return nil
}
