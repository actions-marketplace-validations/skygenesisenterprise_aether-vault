package ui

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Spinner represents a loading spinner
type Spinner struct {
	chars    []string
	delay    time.Duration
	active   bool
	stopChan chan bool
	message  string
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		chars:    []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		delay:    100 * time.Millisecond,
		message:  message,
		stopChan: make(chan bool, 1),
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	if s.active {
		return
	}
	s.active = true

	go func() {
		i := 0
		for {
			select {
			case <-s.stopChan:
				fmt.Fprintf(os.Stdout, "\r%s\r", strings.Repeat(" ", len(s.message)+3))
				return
			default:
				char := s.chars[i%len(s.chars)]
				fmt.Fprintf(os.Stdout, "\r%s %s", char, s.message)
				i++
				time.Sleep(s.delay)
			}
		}
	}()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	if !s.active {
		return
	}
	s.stopChan <- true
	s.active = false
}

// SetMessage updates the spinner message
func (s *Spinner) SetMessage(message string) {
	s.message = message
}

// Success stops spinner and shows success message
func (s *Spinner) Success(message string) {
	s.Stop()
	fmt.Fprintf(os.Stdout, "\r✓ %s\n", message)
}

// Error stops spinner and shows error message
func (s *Spinner) Error(message string) {
	s.Stop()
	fmt.Fprintf(os.Stdout, "\r✗ %s\n", message)
}

// WithSpinner executes a function with a spinner
func WithSpinner(message string, fn func() error) error {
	spinner := NewSpinner(message)
	spinner.Start()
	defer spinner.Stop()

	if err := fn(); err != nil {
		spinner.Error(message + " failed")
		return err
	}

	spinner.Success(message + " completed")
	return nil
}
