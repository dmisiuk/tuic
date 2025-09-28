package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ccpm-demo/internal/calculator"
	"ccpm-demo/internal/ui"
)

func main() {
	// Set up graceful shutdown handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initialize the calculator engine
	calcEngine := calculator.NewEngine()

	// Create the initial model
	model := ui.NewModel(calcEngine)

	// Create the Bubble Tea program with options
	opts := []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
		tea.WithOutput(os.Stderr),
	}

	// Check terminal capabilities
	if !ui.IsTerminalCompatible() {
		fmt.Fprintln(os.Stderr, "Error: Terminal not compatible with TUI")
		fmt.Fprintln(os.Stderr, "Required: 256-color support, mouse events")
		os.Exit(1)
	}

	// Create and start the program
	program := tea.NewProgram(model, opts...)

	// Start the program in a goroutine to handle signals
	go func() {
		if _, err := program.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	<-sigChan

	// Gracefully shutdown the program
	if program != nil {
		program.Kill()
	}

	fmt.Println("\nCCPM Calculator TUI - Gracefully shutdown")
}

func init() {
	// Configure lipgloss for better rendering
	lipgloss.SetHasDarkBackground(true)
}