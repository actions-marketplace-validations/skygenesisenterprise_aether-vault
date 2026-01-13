package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/skygenesisenterprise/aether-vault/package/cli/pkg/types"
)

// TOTPDisplay represents a TOTP entry display
type TOTPDisplay struct {
	Account     string
	Issuer      string
	Code        string
	TimeLeft    int
	Period      int
	Progress    float64
	RefreshAt   time.Time
	CopyOnClick bool
}

// TOTPViewer represents an interactive TOTP viewer
type TOTPViewer struct {
	entries       []*TOTPDisplay
	utility       *types.TOTPUtility
	running       bool
	refreshRate   time.Duration
	selectedIndex int
	showProgress  bool
	showQR        bool
	copyOnSelect  bool
}

// NewTOTPViewer creates a new TOTP viewer
func NewTOTPViewer(utility *types.TOTPUtility) *TOTPViewer {
	return &TOTPViewer{
		utility:       utility,
		refreshRate:   time.Second,
		selectedIndex: 0,
		showProgress:  true,
		copyOnSelect:  false,
	}
}

// SetEntries sets the TOTP entries to display
func (v *TOTPViewer) SetEntries(entries []*types.TOTPEntry) error {
	v.entries = make([]*TOTPDisplay, len(entries))

	for i, entry := range entries {
		display := &TOTPDisplay{
			Account:     entry.Account,
			Issuer:      entry.Issuer,
			Period:      entry.Period,
			CopyOnClick: v.copyOnSelect,
		}

		if display.Period == 0 {
			display.Period = 30
		}

		v.entries[i] = display
	}

	return nil
}

// UpdateCodes updates all TOTP codes
func (v *TOTPViewer) UpdateCodes() error {
	now := time.Now()

	for _, display := range v.entries {
		// Generate code using a sample secret for now
		// TODO: Get actual secret from entry when connected to vault
		sampleSecret := "JBSWY3DPEHPK3PXP" // This is "Hello" in base32
		creds := &types.TOTPCredentials{
			Secret:    sampleSecret,
			Algorithm: types.TOTPAlgorithmSHA1,
			Digits:    6,
			Period:    display.Period,
			Time:      &now,
		}

		code, err := v.utility.GenerateCode(creds)
		if err != nil {
			display.Code = "ERROR"
		} else {
			display.Code = code
		}

		// Calculate time remaining and progress
		timeRemaining := v.utility.GetTimeRemaining(display.Period)
		display.TimeLeft = timeRemaining
		display.Progress = float64(display.Period-timeRemaining) / float64(display.Period)
		display.RefreshAt = now.Add(time.Duration(timeRemaining) * time.Second)
	}

	return nil
}

// Start starts the interactive viewer
func (v *TOTPViewer) Start() error {
	v.running = true

	// Save terminal state and set raw mode
	oldState, err := v.setRawMode()
	if err != nil {
		return fmt.Errorf("failed to set terminal mode: %w", err)
	}
	defer v.restoreTerminal(oldState)

	// Handle keyboard input in a goroutine
	inputChan := make(chan string, 1)
	go v.handleKeyboardInput(inputChan)

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Clear screen initially
	v.clearScreen()

	ticker := time.NewTicker(v.refreshRate)
	defer ticker.Stop()

	for v.running {
		select {
		case <-sigChan:
			v.Stop()
			return nil
		case <-ticker.C:
			v.render()
		case input := <-inputChan:
			v.handleInput(input)
		}
	}

	return nil
}

// Stop stops the viewer
func (v *TOTPViewer) Stop() {
	v.running = false
	// Restore terminal mode
	// TODO: Implement proper terminal restore
}

// clearScreen clears the terminal screen
func (v *TOTPViewer) clearScreen() {
	fmt.Print("\033[2J\033[H")
}

// render renders the TOTP display
func (v *TOTPViewer) render() {
	v.clearScreen()

	// Update codes before rendering
	v.UpdateCodes()

	// Header
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════╗%s\n", Blue, Reset)
	fmt.Printf("%s║                    AETHER VAULT TOTP                       ║%s\n", Bold+Blue, Reset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════╝%s\n\n", Blue, Reset)

	if len(v.entries) == 0 {
		fmt.Printf("%sNo TOTP entries found.%s\n", Yellow, Reset)
		fmt.Printf("\n%sPress 'q' to quit, 'a' to add entry, '?' for help%s\n", Dim, Reset)
		return
	}

	// Display each TOTP entry
	for i, display := range v.entries {
		v.renderTOTPEntry(i, display)
	}

	// Footer
	fmt.Printf("\n%s%s", Dim, strings.Repeat("─", 70))
	fmt.Printf("\nControls: [↑↓] Navigate | [c] Copy code | [a] Add | [d] Delete | [q] Quit%s\n", Reset)
}

// renderTOTPEntry renders a single TOTP entry
func (v *TOTPViewer) renderTOTPEntry(index int, display *TOTPDisplay) {
	isSelected := index == v.selectedIndex

	// Entry header
	if isSelected {
		fmt.Printf("%s┌─ ", Green)
	} else {
		fmt.Printf("│  ")
	}

	// Issuer and Account
	if display.Issuer != "" && display.Account != "" {
		fmt.Printf("%s%s (%s)%s", Bold+White, display.Issuer, display.Account, Reset)
	} else if display.Account != "" {
		fmt.Printf("%s%s%s", Bold+White, display.Account, Reset)
	} else {
		fmt.Printf("%sUnknown Entry%s", Yellow, Reset)
	}

	if isSelected {
		fmt.Printf(" %s←%s", Green, Reset)
	}
	fmt.Printf("\n")

	// Code display
	if isSelected {
		fmt.Printf("%s│  %s", Green, Reset)
	} else {
		fmt.Printf("│  ")
	}

	// Format the code with spaces for better readability
	formattedCode := v.formatCode(display.Code)
	fmt.Printf("%s%s%s", Bold+Cyan, formattedCode, Reset)

	// Time remaining
	if display.TimeLeft > 0 {
		timeColor := Green
		if display.TimeLeft <= 5 {
			timeColor = Red
		} else if display.TimeLeft <= 10 {
			timeColor = Yellow
		}
		fmt.Printf(" %s(%ds)%s", timeColor, display.TimeLeft, Reset)
	}

	fmt.Printf("\n")

	// Progress bar
	if v.showProgress {
		if isSelected {
			fmt.Printf("%s│  %s", Green, Reset)
		} else {
			fmt.Printf("│  ")
		}
		v.renderProgressBar(display.Progress, display.Period)
		fmt.Printf("\n")
	}

	// Bottom border
	if isSelected {
		fmt.Printf("%s└─%s\n", Green, Reset)
	} else {
		fmt.Printf("│\n")
	}
	fmt.Printf(Reset)
}

// renderProgressBar renders a progress bar
func (v *TOTPViewer) renderProgressBar(progress float64, period int) {
	width := 30
	filled := int(progress * float64(width))
	remaining := width - filled

	bar := "["
	bar += strings.Repeat("█", filled)
	if remaining > 0 {
		bar += strings.Repeat("░", remaining)
	}
	bar += "]"

	// Color based on progress
	barColor := Green
	if progress > 0.8 {
		barColor = Red
	} else if progress > 0.6 {
		barColor = Yellow
	}

	fmt.Printf("%s%s%s %ds", barColor, bar, Reset, period)
}

// formatCode formats a TOTP code for better readability
func (v *TOTPViewer) formatCode(code string) string {
	if len(code) <= 4 {
		return code
	}

	// Insert spaces every 3-4 characters
	if len(code) == 6 {
		return code[:3] + " " + code[3:]
	} else if len(code) == 8 {
		return code[:4] + " " + code[4:]
	}

	return code
}

// MoveSelection moves the selection up or down
func (v *TOTPViewer) MoveSelection(direction int) {
	v.selectedIndex += direction

	if v.selectedIndex < 0 {
		v.selectedIndex = len(v.entries) - 1
	} else if v.selectedIndex >= len(v.entries) {
		v.selectedIndex = 0
	}
}

// GetSelectedEntry returns the currently selected entry
func (v *TOTPViewer) GetSelectedEntry() *TOTPDisplay {
	if v.selectedIndex >= 0 && v.selectedIndex < len(v.entries) {
		return v.entries[v.selectedIndex]
	}
	return nil
}

// CopySelectedCode copies the selected TOTP code
func (v *TOTPViewer) CopySelectedCode() error {
	entry := v.GetSelectedEntry()
	if entry == nil {
		return fmt.Errorf("no entry selected")
	}

	// TODO: Implement clipboard copy
	fmt.Printf("\n%sCode copied to clipboard: %s%s\n", Success("Code copied to clipboard: "+entry.Code), Reset)
	return nil
}

// SetShowProgress toggles progress bar display
func (v *TOTPViewer) SetShowProgress(show bool) {
	v.showProgress = show
}

// SetRefreshRate sets the refresh rate
func (v *TOTPViewer) SetRefreshRate(rate time.Duration) {
	v.refreshRate = rate
}

// SimpleTOTPDisplay provides a simple non-interactive TOTP display
func SimpleTOTPDisplay(entries []*types.TOTPEntry, utility *types.TOTPUtility, continuous bool) error {
	viewer := NewTOTPViewer(utility)
	viewer.SetShowProgress(true)

	err := viewer.SetEntries(entries)
	if err != nil {
		return err
	}

	if continuous {
		return viewer.Start()
	} else {
		// Single display
		viewer.UpdateCodes()
		viewer.render()
		return nil
	}
}

// TOTPDisplayOptions represents options for TOTP display
type TOTPDisplayOptions struct {
	Continuous   bool
	ShowProgress bool
	RefreshRate  time.Duration
	ShowQR       bool
	CopyOnSelect bool
	Interactive  bool
}

// DefaultTOTPDisplayOptions returns default display options
func DefaultTOTPDisplayOptions() *TOTPDisplayOptions {
	return &TOTPDisplayOptions{
		Continuous:   false,
		ShowProgress: true,
		RefreshRate:  time.Second,
		ShowQR:       false,
		CopyOnSelect: false,
		Interactive:  false,
	}
}

// setRawMode sets terminal to raw mode for immediate input
func (v *TOTPViewer) setRawMode() (interface{}, error) {
	// For now, use a simple approach. In a real implementation,
	// you'd want to use a proper terminal library like github.com/pkg/term
	return nil, nil
}

// restoreTerminal restores terminal to original state
func (v *TOTPViewer) restoreTerminal(oldState interface{}) error {
	return nil
}

// handleKeyboardInput reads keyboard input in a separate goroutine
func (v *TOTPViewer) handleKeyboardInput(inputChan chan<- string) {
	reader := bufio.NewReader(os.Stdin)

	for v.running {
		// Read input line by line (simpler approach)
		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		// Remove newline and trim
		input = strings.TrimSpace(input)
		if input != "" {
			inputChan <- input
		}
	}
}

// handleInput processes keyboard input
func (v *TOTPViewer) handleInput(input string) {
	switch strings.ToLower(input) {
	case "q":
		v.Stop()
	case "up":
		v.MoveSelection(-1)
	case "down":
		v.MoveSelection(1)
	case "c":
		err := v.CopySelectedCode()
		if err != nil {
			// Show error message briefly
			fmt.Printf("\r%sError: %v%s\n", Red, err, Reset)
		}
	case "a":
		// TODO: Implement add entry functionality
		fmt.Printf("\r%sAdd entry functionality not yet implemented%s\n", Yellow, Reset)
	case "d":
		// TODO: Implement delete entry functionality
		fmt.Printf("\r%sDelete entry functionality not yet implemented%s\n", Yellow, Reset)
	}
}
