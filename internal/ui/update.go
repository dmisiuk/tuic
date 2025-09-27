package ui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// update handles all incoming messages and updates the model state
func update(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return handleKeyMsg(m, msg)

	case tea.MouseMsg:
		return handleMouseMsg(m, msg)

	case tea.WindowSizeMsg:
		return handleWindowSizeMsg(m, msg)

	case tea.QuitMsg:
		return handleQuitMsg(m)

	default:
		return m, nil
	}
}

// handleKeyMsg processes keyboard input
func handleKeyMsg(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Clear any existing errors
	m.clearError()

	// First, handle special keys that should always work
	switch msg.Type {
	case tea.KeyEsc, tea.KeyCtrlC:
		m.quitting = true
		return m, tea.Quit

	case tea.KeyBackspace:
		return handleBackspaceKey(m)

	case tea.KeyDelete:
		return handleDeleteKey(m)

	case tea.KeyLeft:
		return handleLeftKey(m)

	case tea.KeyRight:
		return handleRightKey(m)

	case tea.KeyUp:
		return handleUpKey(m)

	case tea.KeyDown:
		return handleDownKey(m)

	case tea.KeyEnter:
		// Handle button grid first, then fall back to default
		if action := m.buttonGrid.HandleKeyPress(msg); action != nil {
			return handleButtonGridAction(m, action)
		}
		return handleEnterKey(m)

	case tea.KeySpace:
		// Handle button grid navigation and activation
		if action := m.buttonGrid.HandleKeyPress(msg); action != nil {
			return handleButtonGridAction(m, action)
		}

	case tea.KeyRunes:
		// Handle button grid first for direct input
		if action := m.buttonGrid.HandleKeyPress(msg); action != nil {
			return handleButtonGridAction(m, action)
		}
		return handleRunes(m, msg)

	default:
		// Try button grid for other keys
		if action := m.buttonGrid.HandleKeyPress(msg); action != nil {
			return handleButtonGridAction(m, action)
		}
		return m, nil
	}
}

// handleMouseMsg processes mouse events
func handleMouseMsg(m Model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.MouseLeft:
		// Handle button grid clicks first
		if action := m.buttonGrid.HandleMouse(msg); action != nil {
			return handleButtonGridAction(m, action)
		}
		return handleMouseClick(m, msg)

	case tea.MouseWheelUp:
		return handleMouseWheelUp(m)

	case tea.MouseWheelDown:
		return handleMouseWheelDown(m)

	default:
		return m, nil
	}
}

// handleWindowSizeMsg handles terminal resize events
func handleWindowSizeMsg(m Model, msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	m.ready = true
	return m, nil
}

// handleQuitMsg handles quit messages
func handleQuitMsg(m Model) (tea.Model, tea.Cmd) {
	m.quitting = true
	return m, tea.Quit
}

// handleEnterKey processes Enter key press
func handleEnterKey(m Model) (tea.Model, tea.Cmd) {
	if m.input == "" {
		return m, nil
	}

	// Try to evaluate the input expression
	result, err := m.engine.Evaluate(m.input)
	if err != nil {
		m.setError(err)
		return m, nil
	}

	// Update output and history
	m.output = m.formatValue(result)
	m.addToHistory(fmt.Sprintf("%s = %s", m.input, m.output))

	// Reset input
	m.input = ""
	m.cursorPosition = 0

	// Update calculator state
	m.calculatorState.displayValue = m.output
	m.calculatorState.isWaitingForOperand = true

	return m, nil
}

// handleBackspaceKey processes Backspace key
func handleBackspaceKey(m Model) (tea.Model, tea.Cmd) {
	if m.cursorPosition > 0 {
		m.input = m.input[:m.cursorPosition-1] + m.input[m.cursorPosition:]
		m.cursorPosition--
	}
	return m, nil
}

// handleDeleteKey processes Delete key
func handleDeleteKey(m Model) (tea.Model, tea.Cmd) {
	if m.cursorPosition < len(m.input) {
		m.input = m.input[:m.cursorPosition] + m.input[m.cursorPosition+1:]
	}
	return m, nil
}

// handleLeftKey processes Left arrow key
func handleLeftKey(m Model) (tea.Model, tea.Cmd) {
	if m.cursorPosition > 0 {
		m.cursorPosition--
	}
	return m, nil
}

// handleRightKey processes Right arrow key
func handleRightKey(m Model) (tea.Model, tea.Cmd) {
	if m.cursorPosition < len(m.input) {
		m.cursorPosition++
	}
	return m, nil
}

// handleUpKey processes Up arrow key
func handleUpKey(m Model) (tea.Model, tea.Cmd) {
	if m.historyIndex > 0 {
		m.historyIndex--
		// Set input to the selected history item (extract expression part)
		historyEntry := m.history[m.historyIndex]
		if parts := strings.Split(historyEntry, " = "); len(parts) > 0 {
			m.input = parts[0]
			m.cursorPosition = len(m.input)
		}
	}
	return m, nil
}

// handleDownKey processes Down arrow key
func handleDownKey(m Model) (tea.Model, tea.Cmd) {
	if m.historyIndex < len(m.history)-1 {
		m.historyIndex++
		historyEntry := m.history[m.historyIndex]
		if parts := strings.Split(historyEntry, " = "); len(parts) > 0 {
			m.input = parts[0]
			m.cursorPosition = len(m.input)
		}
	} else if m.historyIndex == len(m.history)-1 {
		m.historyIndex = len(m.history) // Set to end
		m.input = ""
		m.cursorPosition = 0
	}
	return m, nil
}

// handleRunes processes character input
func handleRunes(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	char := string(msg.Runes)

	switch char {
	case "q":
		m.quitting = true
		return m, tea.Quit

	case "h":
		// Toggle help - could be implemented later
		return m, nil

	case "c":
		// Clear input
		m.input = ""
		m.cursorPosition = 0
		m.calculatorState.displayValue = "0"
		return m, nil

	case "+", "-", "*", "/":
		// Handle operators
		if m.input != "" {
			m.input += " " + char + " "
			m.cursorPosition = len(m.input)
		}
		return m, nil

	case "×":
		// Handle multiplication symbol
		if m.input != "" {
			m.input += " * "
			m.cursorPosition = len(m.input)
		}
		return m, nil

	case "÷":
		// Handle division symbol
		if m.input != "" {
			m.input += " / "
			m.cursorPosition = len(m.input)
		}
		return m, nil

	case "=":
		// Calculate result
		return handleEnterKey(m)

	case ".":
		// Handle decimal point
		if m.input == "" {
			m.input = "0."
			m.cursorPosition = 2
		} else {
			// Check if last character is a digit
			if len(m.input) > 0 {
				lastChar := m.input[len(m.input)-1]
				if lastChar >= '0' && lastChar <= '9' {
					m.input += "."
					m.cursorPosition++
				}
			}
		}
		return m, nil

	default:
		// Handle numbers and other valid characters
		if char >= "0" && char <= "9" {
			m.input += char
			m.cursorPosition++
		} else if char == " " {
			// Allow spaces for formatting
			m.input += char
			m.cursorPosition++
		}
		// Could add more validation here for other valid characters
		return m, nil
	}
}

// handleMouseClick processes mouse clicks
func handleMouseClick(m Model, msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// This is a simplified implementation
	// In a full implementation, we'd calculate which button was clicked
	// and handle it appropriately

	// For now, just acknowledge the click
	return m, nil
}

// handleMouseWheelUp processes mouse wheel up
func handleMouseWheelUp(m Model) (tea.Model, tea.Cmd) {
	return handleUpKey(m)
}

// handleMouseWheelDown processes mouse wheel down
func handleMouseWheelDown(m Model) (tea.Model, tea.Cmd) {
	return handleDownKey(m)
}

// handleCalculatorButton processes calculator button clicks
func handleCalculatorButton(m Model, button string) (tea.Model, tea.Cmd) {
	m.clearError()

	switch button {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		m.input += button
		m.cursorPosition++

	case ".":
		if m.input == "" {
			m.input = "0."
			m.cursorPosition = 2
		} else {
			// Check if last character is a digit
			if len(m.input) > 0 {
				lastChar := m.input[len(m.input)-1]
				if lastChar >= '0' && lastChar <= '9' {
					m.input += "."
					m.cursorPosition++
				}
			}
		}

	case "+", "-", "*", "/":
		if m.input != "" {
			m.input += " " + button + " "
			m.cursorPosition = len(m.input)
		}

	case "×":
		if m.input != "" {
			m.input += " * "
			m.cursorPosition = len(m.input)
		}

	case "÷":
		if m.input != "" {
			m.input += " / "
			m.cursorPosition = len(m.input)
		}

	case "=":
		return handleEnterKey(m)

	case "C":
		m.input = ""
		m.cursorPosition = 0
		m.calculatorState.displayValue = "0"

	case "±":
		if m.input != "" {
			// Try to parse as number and negate
			if num, err := strconv.ParseFloat(m.input, 64); err == nil {
				m.input = fmt.Sprintf("%g", -num)
				m.cursorPosition = len(m.input)
			}
		}

	case "%":
		if m.input != "" {
			// Convert to percentage
			if num, err := strconv.ParseFloat(m.input, 64); err == nil {
				m.input = fmt.Sprintf("%g", num/100)
				m.cursorPosition = len(m.input)
			}
		}

	case "⌫":
		return handleBackspaceKey(m)
	}

	return m, nil
}

// handleButtonGridAction processes actions from the button grid
func handleButtonGridAction(m Model, action *integration.ButtonAction) (tea.Model, tea.Cmd) {
	m.clearError()

	// Process the button action based on its value
	switch action.Value {
	case "clear":
		// Clear all input and reset calculator state
		m.input = ""
		m.output = ""
		m.cursorPosition = 0
		m.calculatorState.displayValue = "0"
		m.calculatorState.operator = ""
		m.calculatorState.previousValue = 0
		m.calculatorState.isWaitingForOperand = false

	case "clear_entry":
		// Clear current input only
		m.input = ""
		m.cursorPosition = 0
		m.calculatorState.displayValue = "0"

	case "backspace":
		return handleBackspaceKey(m)

	case "+", "-", "*", "/":
		// Handle operators
		if m.input != "" {
			m.input += " " + action.Value + " "
			m.cursorPosition = len(m.input)
		}

	case "=":
		return handleEnterKey(m)

	case ".":
		// Handle decimal point
		if m.input == "" {
			m.input = "0."
			m.cursorPosition = 2
		} else {
			// Check if last character is a digit
			if len(m.input) > 0 {
				lastChar := m.input[len(m.input)-1]
				if lastChar >= '0' && lastChar <= '9' {
					m.input += "."
					m.cursorPosition++
				}
			}
		}

	default:
		// Handle numbers (0-9)
		if len(action.Value) == 1 && action.Value >= "0" && action.Value <= "9" {
			m.input += action.Value
			m.cursorPosition++

			// Update calculator state display
			m.calculatorState.displayValue = m.input
		}
	}

	return m, nil
}