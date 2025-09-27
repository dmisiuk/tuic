package ui

import (
	"os"
	"runtime"
	"strings"

	"golang.org/x/term"
)

// IsTerminalCompatible checks if the current terminal supports required features
func IsTerminalCompatible() bool {
	// Check if we're running in a terminal
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}

	// Check for basic terminal capabilities
	return checkTerminalCapabilities()
}

// checkTerminalCapabilities performs basic terminal capability checks
func checkTerminalCapabilities() bool {
	// Check for color support
	if !hasColorSupport() {
		return false
	}

	// Check for minimum terminal size
	if !hasMinimumSize() {
		return false
	}

	// Check for mouse support (optional, but nice to have)
	_ = hasMouseSupport()

	return true
}

// hasColorSupport checks if the terminal supports colors
func hasColorSupport() bool {
	// Check for common color environment variables
	colorVars := []string{
		"TERM",
		"COLORTERM",
		"CLICOLOR",
		"FORCE_COLOR",
	}

	for _, envVar := range colorVars {
		if value := os.Getenv(envVar); value != "" {
			if envVar == "TERM" {
				// Check for color-capable terminals
				if strings.Contains(value, "color") ||
					strings.Contains(value, "256color") ||
					strings.Contains(value, "truecolor") ||
					strings.Contains(value, "direct") {
					return true
				}
			} else {
				// Other color variables typically indicate color support
				return true
			}
		}
	}

	// On Windows, check if we're running in a modern terminal
	if runtime.GOOS == "windows" {
		// Windows 10+ terminals support colors
		return true
	}

	// Default assumption: most modern terminals support colors
	return true
}

// hasMinimumSize checks if the terminal meets minimum size requirements
func hasMinimumSize() bool {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// If we can't get the size, assume it's adequate
		return true
	}

	// Minimum requirements: 60x25
	return width >= 60 && height >= 25
}

// hasMouseSupport checks if the terminal supports mouse events
func hasMouseSupport() bool {
	// Check for mouse support environment variables
	mouseVars := []string{
		"TERM_PROGRAM",
		"TERM",
	}

	for _, envVar := range mouseVars {
		if value := os.Getenv(envVar); value != "" {
			// Common terminal programs that support mouse
			switch value {
			case "iTerm.app", "Apple_Terminal", "vscode", "tmux":
				return true
			}
			// Check for xterm-compatible terminals
			if strings.Contains(value, "xterm") {
				return true
			}
		}
	}

	return false
}

// GetTerminalSize returns the current terminal dimensions
func GetTerminalSize() (width, height int, err error) {
	return term.GetSize(int(os.Stdout.Fd()))
}

// IsRunningInDocker checks if the application is running inside Docker
func IsRunningInDocker() bool {
	// Check for Docker-specific files or environment variables
	dockerFiles := []string{
		"/.dockerenv",
		"/proc/1/cgroup",
	}

	for _, file := range dockerFiles {
		if _, err := os.Stat(file); err == nil {
			// Check if the file contains Docker-specific content
			if content, err := os.ReadFile(file); err == nil {
				if strings.Contains(string(content), "docker") {
					return true
				}
			}
			return true
		}
	}

	// Check Docker environment variables
	dockerEnvVars := []string{
		"DOCKER_CONTAINER",
		"DOCKER_IMAGE_NAME",
	}

	for _, envVar := range dockerEnvVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}

	return false
}

// IsRunningInCI checks if the application is running in a CI environment
func IsRunningInCI() bool {
	ciEnvVars := []string{
		"CI",
		"CONTINUOUS_INTEGRATION",
		"JENKINS_URL",
		"GITHUB_ACTIONS",
		"TRAVIS",
		"CIRCLECI",
		"GITLAB_CI",
		"BUILDKITE",
		"AZURE_PIPELINES",
	}

	for _, envVar := range ciEnvVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}

	return false
}

// GetTerminalInfo returns information about the current terminal
func GetTerminalInfo() map[string]string {
	info := make(map[string]string)

	// Terminal type
	info["TERM"] = os.Getenv("TERM")
	info["SHELL"] = os.Getenv("SHELL")
	info["TERM_PROGRAM"] = os.Getenv("TERM_PROGRAM")

	// Operating system
	info["OS"] = runtime.GOOS
	info["ARCH"] = runtime.GOARCH

	// Environment flags
	info["COLOR_SUPPORT"] = formatBool(hasColorSupport())
	info["MOUSE_SUPPORT"] = formatBool(hasMouseSupport())
	info["IN_DOCKER"] = formatBool(IsRunningInDocker())
	info["IN_CI"] = formatBool(IsRunningInCI())

	// Terminal size
	if width, height, err := GetTerminalSize(); err == nil {
		info["WIDTH"] = formatInt(width)
		info["HEIGHT"] = formatInt(height)
	} else {
		info["WIDTH"] = "unknown"
		info["HEIGHT"] = "unknown"
	}

	return info
}

// formatBool converts a boolean to a string representation
func formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// formatInt converts an integer to a string representation
func formatInt(i int) string {
	// Handle zero case
	if i == 0 {
		return "0"
	}

	// Handle negative numbers
	var result string
	isNegative := false
	if i < 0 {
		isNegative = true
		i = -i
	}

	// Convert digits
	for i > 0 {
		digit := i % 10
		result = string(rune('0'+digit)) + result
		i = i / 10
	}

	// Add negative sign if needed
	if isNegative {
		result = "-" + result
	}

	return result
}

// SafePrint handles printing in environments where stdout might not be available
func SafePrint(message string) {
	// Try to print to stdout
	if term.IsTerminal(int(os.Stdout.Fd())) {
		os.Stdout.WriteString(message + "\n")
		return
	}

	// Fallback to stderr
	if term.IsTerminal(int(os.Stderr.Fd())) {
		os.Stderr.WriteString(message + "\n")
		return
	}

	// Last resort: ignore the message
}

// HandleTerminalResize provides graceful handling of terminal resize events
func HandleTerminalResize(width, height int) error {
	if width < 60 || height < 25 {
		return &TerminalError{
			Message: "Terminal too small (minimum 60x25 required)",
			Code:    "TERMINAL_TOO_SMALL",
		}
	}

	return nil
}

// TerminalError represents a terminal-related error
type TerminalError struct {
	Message string
	Code    string
}

func (e *TerminalError) Error() string {
	return e.Message
}

// IsTerminalError checks if an error is a TerminalError
func IsTerminalError(err error) bool {
	_, ok := err.(*TerminalError)
	return ok
}