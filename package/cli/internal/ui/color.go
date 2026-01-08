package ui

import (
	"fmt"
	"os"
)

// Color codes for terminal output
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
)

// IsColorSupported returns true if color output is supported
func IsColorSupported() bool {
	// Check if NO_COLOR environment variable is set
	if os.Getenv("NO_COLOR") != "" {
		return false
	}

	// Check if TERM environment variable indicates color support
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}

	return true
}

// Colorize applies color to text if color is supported
func Colorize(text, color string) string {
	if !IsColorSupported() {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

// RedText colors text red
func RedText(text string) string {
	return Colorize(text, Red)
}

// GreenText colors text green
func GreenText(text string) string {
	return Colorize(text, Green)
}

// YellowText colors text yellow
func YellowText(text string) string {
	return Colorize(text, Yellow)
}

// BlueText colors text blue
func BlueText(text string) string {
	return Colorize(text, Blue)
}

// CyanText colors text cyan
func CyanText(text string) string {
	return Colorize(text, Cyan)
}

// BoldText makes text bold
func BoldText(text string) string {
	return Colorize(text, Bold)
}

// DimText makes text dim
func DimText(text string) string {
	return Colorize(text, Dim)
}

// Success formats a success message
func Success(text string) string {
	return Colorize("✓ "+text, Green)
}

// Error formats an error message
func Error(text string) string {
	return Colorize("✗ "+text, Red)
}

// Warning formats a warning message
func Warning(text string) string {
	return Colorize("⚠ "+text, Yellow)
}

// Info formats an info message
func Info(text string) string {
	return Colorize("ℹ "+text, Blue)
}
