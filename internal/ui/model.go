package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ccpm-demo/internal/calculator"
)

// Model represents the application state following the MVU pattern
type Model struct {
	// Calculator engine reference
	engine *calculator.Engine

	// Terminal dimensions
	width  int
	height int

	// Application state
	calculatorState calculatorState
	input          string
	output         string
	error          string
	cursorPosition int
	history        []string
	historyIndex   int

	// UI state
	ready bool
	quitting bool

	// Styling
	styles styles
}

// calculatorState represents the current calculator state
type calculatorState struct {
	displayValue string
	operator     string
	previousValue float64
	isWaitingForOperand bool
}

// styles contains all the lipgloss styles for the UI
type styles struct {
	app      lipgloss.Style
	title    lipgloss.Style
	display  lipgloss.Style
	input    lipgloss.Style
	output   lipgloss.Style
	error    lipgloss.Style
	buttons  lipgloss.Style
	button   lipgloss.Style
	active   lipgloss.Style
	inactive lipgloss.Style
}

// NewModel creates a new application model
func NewModel(engine *calculator.Engine) Model {
	return Model{
		engine: engine,
		calculatorState: calculatorState{
			displayValue: "0",
			operator:     "",
			previousValue: 0,
			isWaitingForOperand: false,
		},
		input:        "",
		output:       "",
		error:        "",
		cursorPosition: 0,
		history:      []string{},
		historyIndex: -1,
		ready:        false,
		quitting:     false,
		styles:       defaultStyles(),
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return update(m, msg)
}

// View implements tea.Model
func (m Model) View() string {
	return view(m)
}

// defaultStyles returns the default styling for the application
func defaultStyles() styles {
	return styles{
		app: lipgloss.NewStyle().
			Width(60).
			Height(30).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")),

		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Align(lipgloss.Center).
			Padding(0, 1),

		display: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("236")).
			Align(lipgloss.Right).
			Padding(0, 1).
			Width(56).
			Height(3),

		input: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("237")).
			Align(lipgloss.Left).
			Padding(0, 1).
			Width(56),

		output: lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Background(lipgloss.Color("235")).
			Align(lipgloss.Right).
			Padding(0, 1).
			Width(56).
			Height(2),

		error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Background(lipgloss.Color("52")).
			Align(lipgloss.Center).
			Padding(0, 1).
			Width(56),

		buttons: lipgloss.NewStyle().
			Width(56).
			Height(15),

		button: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("240")).
			Align(lipgloss.Center).
			Width(12).
			Height(3),

		active: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("62")).
			Align(lipgloss.Center).
			Width(12).
			Height(3),

		inactive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Background(lipgloss.Color("240")).
			Align(lipgloss.Center).
			Width(12).
			Height(3),
	}
}

// formatValue formats a float value for display
func (m Model) formatValue(value float64) string {
	// Remove trailing .0 for whole numbers
	if value == float64(int(value)) {
		return fmt.Sprintf("%.0f", value)
	}
	return fmt.Sprintf("%.6f", value)
}

// truncateString truncates a string to fit within a width
func (m Model) truncateString(str string, width int) string {
	if len(str) <= width {
		return str
	}
	// If width is too small for "...", just return "..."
	if width <= 3 {
		return "..."
	}
	return str[:width-3] + "..."
}

// addToHistory adds an expression to the history
func (m *Model) addToHistory(expression string) {
	m.history = append(m.history, expression)
	if len(m.history) > 100 { // Keep last 100 entries
		m.history = m.history[1:]
	}
	m.historyIndex = len(m.history) - 1
}

// clearError clears any error message
func (m *Model) clearError() {
	m.error = ""
}

// setError sets an error message
func (m *Model) setError(err error) {
	if err != nil {
		m.error = err.Error()
	} else {
		m.error = ""
	}
}

// getDisplayWidth returns the available display width
func (m Model) getDisplayWidth() int {
	if m.width > 0 {
		return m.width - 4 // Account for borders
	}
	return 56
}

// getDisplayHeight returns the available display height
func (m Model) getDisplayHeight() int {
	if m.height > 0 {
		return m.height - 4 // Account for borders
	}
	return 30
}

// GetInput returns the current input string
func (m Model) GetInput() string {
	return m.input
}

// SetInput sets the input string and updates cursor position
func (m *Model) SetInput(input string) {
	m.input = input
	m.cursorPosition = len(input)
}

// GetOutput returns the current output string
func (m Model) GetOutput() string {
	return m.output
}

// SetOutput sets the output string
func (m *Model) SetOutput(output string) {
	m.output = output
}

// GetCursorPosition returns the current cursor position
func (m Model) GetCursorPosition() int {
	return m.cursorPosition
}

// SetCursorPosition sets the cursor position
func (m *Model) SetCursorPosition(pos int) {
	if pos >= 0 && pos <= len(m.input) {
		m.cursorPosition = pos
	}
}

// GetError returns the current error message
func (m Model) GetError() string {
	return m.error
}

// SetError sets the error message
func (m *Model) SetError(err string) {
	m.error = err
}

// ClearError clears the error message
func (m *Model) ClearError() {
	m.error = ""
}