package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// view renders the application UI
func view(m Model) string {
	if m.quitting {
		return ""
	}

	// Update styles with current terminal dimensions
	styles := m.updateStyles()

	// Build the main layout
	content := strings.Builder{}

	// Title
	content.WriteString(styles.title.Render("CCPM Calculator"))
	content.WriteString("\n\n")

	// Display area (current calculator state)
	content.WriteString(styles.display.Render(m.calculatorState.displayValue))
	content.WriteString("\n")

	// Input area
	inputText := m.input
	if m.cursorPosition >= 0 && m.cursorPosition < len(m.input) {
		// Show cursor position
		before := m.input[:m.cursorPosition]
		after := m.input[m.cursorPosition:]
		inputText = before + "█" + after
	}
	content.WriteString(styles.input.Render(inputText))
	content.WriteString("\n")

	// Output area (results)
	content.WriteString(styles.output.Render(m.output))
	content.WriteString("\n")

	// Error area
	if m.error != "" {
		content.WriteString(styles.error.Render("Error: " + m.error))
		content.WriteString("\n")
	}

	// Button layout
	content.WriteString("\n")
	content.WriteString(m.renderButtons(styles))

	// History (if any)
	if len(m.history) > 0 {
		content.WriteString("\n")
		content.WriteString(m.renderHistory(styles))
	}

	// Wrap everything in the main container
	return styles.app.Render(content.String())
}

// renderButtons creates the calculator button layout
func (m Model) renderButtons(styles styles) string {
	buttons := strings.Builder{}

	// Button layout: 4x4 grid
	buttonLayout := [][]string{
		{"C", "±", "%", "÷"},
		{"7", "8", "9", "×"},
		{"4", "5", "6", "-"},
		{"1", "2", "3", "+"},
		{"0", ".", "=", "⌫"},
	}

	for _, row := range buttonLayout {
		rowButtons := make([]string, len(row))
		for i, button := range row {
			// Style based on button type
			buttonStyle := styles.button
			switch button {
			case "C", "±", "%", "⌫":
				buttonStyle = styles.inactive
			case "÷", "×", "-", "+", "=":
				buttonStyle = styles.active
			}

			rowButtons[i] = buttonStyle.Render(button)
		}
		buttons.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, rowButtons...))
		buttons.WriteString("\n")
	}

	return styles.buttons.Render(buttons.String())
}

// renderHistory shows calculation history
func (m Model) renderHistory(styles styles) string {
	if len(m.history) == 0 {
		return ""
	}

	history := strings.Builder{}
	history.WriteString("History:\n")

	// Show last 5 entries
	start := 0
	if len(m.history) > 5 {
		start = len(m.history) - 5
	}

	for i := start; i < len(m.history); i++ {
		prefix := "  "
		if i == m.historyIndex {
			prefix = "→ "
		}
		history.WriteString(prefix + m.history[i] + "\n")
	}

	return history.String()
}

// updateStyles updates the styles based on current terminal dimensions
func (m Model) updateStyles() styles {
	styles := m.styles

	// Adjust width based on terminal size
	appWidth := m.getDisplayWidth()
	if appWidth > 80 {
		appWidth = 80
	} else if appWidth < 60 {
		appWidth = 60
	}

	// Adjust height based on terminal size
	appHeight := m.getDisplayHeight()
	if appHeight > 40 {
		appHeight = 40
	} else if appHeight < 25 {
		appHeight = 25
	}

	// Update styles with new dimensions
	styles.app = styles.app.Width(appWidth).Height(appHeight)
	styles.display = styles.display.Width(appWidth - 4)
	styles.input = styles.input.Width(appWidth - 4)
	styles.output = styles.output.Width(appWidth - 4)
	styles.error = styles.error.Width(appWidth - 4)
	styles.buttons = styles.buttons.Width(appWidth - 4)

	return styles
}

// renderWelcomeScreen shows the initial welcome screen
func (m Model) renderWelcomeScreen() string {
	if m.ready {
		return ""
	}

	styles := m.updateStyles()

	welcome := strings.Builder{}
	welcome.WriteString(styles.title.Render("Welcome to CCPM Calculator"))
	welcome.WriteString("\n\n")
	welcome.WriteString(styles.display.Render("Press any key to start..."))
	welcome.WriteString("\n")
	welcome.WriteString(styles.input.Render("Terminal UI initializing..."))

	return styles.app.Render(welcome.String())
}

// renderSplashScreen shows a splash screen during startup
func (m Model) renderSplashScreen() string {
	styles := m.updateStyles()

	splash := strings.Builder{}
	splash.WriteString(styles.title.Render("CCPM Calculator"))
	splash.WriteString("\n\n")
	splash.WriteString(styles.display.Render("Critical Chain Project Management"))
	splash.WriteString("\n")
	splash.WriteString(styles.input.Render("Calculator & Planning Tool"))
	splash.WriteString("\n\n")
	splash.WriteString(styles.output.Render("Loading calculator engine..."))

	return styles.app.Render(splash.String())
}

// renderHelpScreen shows keyboard shortcuts and help
func (m Model) renderHelpScreen() string {
	styles := m.updateStyles()

	help := strings.Builder{}
	help.WriteString(styles.title.Render("Help & Shortcuts"))
	help.WriteString("\n\n")

	helpContent := `Calculator:
  0-9, .  - Numbers and decimal point
  +, -, ×, ÷ - Basic operations
  =        - Calculate result
  C        - Clear
  ±        - Toggle sign
  %        - Percentage
  ⌫        - Backspace

Navigation:
  q, Esc   - Quit
  h        - Toggle help
  ↑, ↓     - Navigate history
  Enter    - Execute calculation

Mouse:
  Click buttons with mouse
  Scroll to navigate history
`

	help.WriteString(styles.display.Render(helpContent))

	return styles.app.Render(help.String())
}