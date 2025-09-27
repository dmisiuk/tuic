package input

import (
	tea "github.com/charmbracelet/bubbletea"

	"ccpm-demo/internal/ui"
)

// KeyAction represents the action to perform when a key is pressed
type KeyAction int

const (
	KeyActionNone KeyAction = iota
	KeyActionNumber
	KeyActionOperator
	KeyActionEquals
	KeyActionClear
	KeyActionBackspace
	KeyActionNavigate
	KeyActionFocusActivate
	KeyActionQuit
)

// KeyEvent represents a keyboard input event
type KeyEvent struct {
	Key       tea.KeyType
	Rune      rune
	Action    KeyAction
	Value     string
	Modifiers tea.KeyModifier
}

// KeyHandler defines the interface for keyboard input handling
type KeyHandler interface {
	// HandleKey processes a keyboard event and returns the resulting model and command
	HandleKey(model ui.Model, msg tea.KeyMsg) (ui.Model, tea.Cmd)

	// IsDirectKey returns true if the key should be processed regardless of focus
	IsDirectKey(msg tea.KeyMsg) bool

	// GetActionForKey determines the action for a given key press
	GetActionForKey(msg tea.KeyMsg) KeyEvent
}

// KeyboardHandler implements the KeyHandler interface
type KeyboardHandler struct {
	keyBindingManager *KeyBindingManager
}

// NewKeyboardHandler creates a new keyboard handler
func NewKeyboardHandler() *KeyboardHandler {
	return &KeyboardHandler{
		keyBindingManager: NewKeyBindingManager(DefaultKeyBindings()),
	}
}

// HandleKey processes a keyboard event and returns the resulting model and command
func (kh *KeyboardHandler) HandleKey(model ui.Model, msg tea.KeyMsg) (ui.Model, tea.Cmd) {
	// Check if this is a direct key (should be processed regardless of focus)
	if kh.IsDirectKey(msg) {
		return kh.handleDirectKey(model, msg)
	}

	// Check if this is a navigation key
	keyEvent := kh.GetActionForKey(msg)
	if keyEvent.Action == KeyActionNavigate || keyEvent.Action == KeyActionFocusActivate {
		return kh.handleNavigationKey(model, msg)
	}

	// Default to regular key handling
	return kh.handleRegularKey(model, msg)
}

// IsDirectKey returns true if the key should be processed regardless of focus
func (kh *KeyboardHandler) IsDirectKey(msg tea.KeyMsg) bool {
	return kh.keyBindingManager.IsDirectKey(msg)
}

// GetActionForKey determines the action for a given key press
func (kh *KeyboardHandler) GetActionForKey(msg tea.KeyMsg) KeyEvent {
	binding := kh.keyBindingManager.GetActionForKey(msg)
	if binding == nil {
		return KeyEvent{
			Key:       msg.Type,
			Rune:      msg.Runes,
			Action:    KeyActionNone,
			Value:     "",
			Modifiers: msg.Modifiers,
		}
	}

	return KeyEvent{
		Key:       binding.Key,
		Rune:      msg.Runes,
		Action:    binding.Action,
		Value:     binding.Value,
		Modifiers: binding.Modifiers,
	}
}

// handleDirectKey handles keys that should be processed regardless of focus
func (kh *KeyboardHandler) handleDirectKey(model ui.Model, msg tea.KeyMsg) (ui.Model, tea.Cmd) {
	keyEvent := kh.GetActionForKey(msg)

	switch keyEvent.Action {
	case KeyActionNumber:
		return kh.handleNumberInput(model, keyEvent.Value)
	case KeyActionOperator:
		return kh.handleOperatorInput(model, keyEvent.Value)
	case KeyActionEquals:
		return kh.handleEquals(model)
	case KeyActionBackspace:
		return kh.handleBackspace(model)
	case KeyActionQuit:
		return kh.handleQuit(model)
	default:
		return model, nil
	}
}

// handleNavigationKey handles navigation and focus-related keys
func (kh *KeyboardHandler) handleNavigationKey(model ui.Model, msg tea.KeyMsg) (ui.Model, tea.Cmd) {
	keyEvent := kh.GetActionForKey(msg)

	switch keyEvent.Action {
	case KeyActionNavigate:
		return kh.handleNavigation(model, keyEvent.Value)
	case KeyActionFocusActivate:
		return kh.handleFocusActivate(model)
	default:
		return model, nil
	}
}

// handleRegularKey handles regular key input that depends on focus
func (kh *KeyboardHandler) handleRegularKey(model ui.Model, msg tea.KeyMsg) (ui.Model, tea.Cmd) {
	keyEvent := kh.GetActionForKey(msg)

	switch keyEvent.Action {
	case KeyActionNumber:
		return kh.handleNumberInput(model, keyEvent.Value)
	case KeyActionOperator:
		return kh.handleOperatorInput(model, keyEvent.Value)
	case KeyActionEquals:
		return kh.handleEquals(model)
	case KeyActionClear:
		return kh.handleClear(model)
	default:
		return model, nil
	}
}

// handleNumberInput handles number and decimal point input
func (kh *KeyboardHandler) handleNumberInput(model ui.Model, value string) (ui.Model, tea.Cmd) {
	// This will be integrated with the model's input handling
	// For now, we'll use a simplified approach that updates the model directly

	// Add the number to the current input
	if value == "." {
		// Handle decimal point
		if model.GetInput() == "" {
			model.SetInput("0.")
		} else {
			// Check if last character is a digit
			input := model.GetInput()
			if len(input) > 0 {
				lastChar := input[len(input)-1]
				if lastChar >= '0' && lastChar <= '9' {
					model.SetInput(input + ".")
				}
			}
		}
	} else {
		// Handle regular numbers
		model.SetInput(model.GetInput() + value)
	}

	return model, nil
}

// handleOperatorInput handles operator keys
func (kh *KeyboardHandler) handleOperatorInput(model ui.Model, operator string) (ui.Model, tea.Cmd) {
	input := model.GetInput()
	if input != "" {
		model.SetInput(input + " " + operator + " ")
	}

	return model, nil
}

// handleEquals handles the equals/enter key
func (kh *KeyboardHandler) handleEquals(model ui.Model) (ui.Model, tea.Cmd) {
	input := model.GetInput()
	if input == "" {
		return model, nil
	}

	// Trigger calculation by updating the model state
	// The actual calculation should be done by the calculator engine
	// For now, we'll just clear the input to simulate calculation completion
	model.SetInput("")
	return model, nil
}

// handleBackspace handles backspace/delete keys
func (kh *KeyboardHandler) handleBackspace(model ui.Model) (ui.Model, tea.Cmd) {
	input := model.GetInput()
	if len(input) > 0 {
		// Remove the last character
		input = input[:len(input)-1]
		model.SetInput(input)
	}

	return model, nil
}

// handleClear handles clear operations
func (kh *KeyboardHandler) handleClear(model ui.Model) (ui.Model, tea.Cmd) {
	model.SetInput("")
	return model, nil
}

// handleNavigation handles navigation keys
func (kh *KeyboardHandler) handleNavigation(model ui.Model, direction string) (ui.Model, tea.Cmd) {
	// Basic navigation support - will be enhanced with focus management
	// For now, we'll provide minimal navigation functionality

	switch direction {
	case "up":
		// Navigate up in button grid
		// This will be implemented when focus management is available
		return model, nil
	case "down":
		// Navigate down in button grid
		return model, nil
	case "left":
		// Navigate left in button grid
		return model, nil
	case "right":
		// Navigate right in button grid
		return model, nil
	case "tab":
		// Navigate to next focusable element
		return model, nil
	case "shift_tab":
		// Navigate to previous focusable element
		return model, nil
	default:
		return model, nil
	}
}

// handleFocusActivate handles focus activation (space key)
func (kh *KeyboardHandler) handleFocusActivate(model ui.Model) (ui.Model, tea.Cmd) {
	// Activate the currently focused button
	// This will be implemented when focus management is available
	// For now, we'll provide a placeholder implementation

	// When focus management is implemented, this will:
	// 1. Get the currently focused button
	// 2. Simulate a click on that button
	// 3. Update the model state accordingly

	return model, nil
}

// handleQuit handles quit operations
func (kh *KeyboardHandler) handleQuit(model ui.Model) (ui.Model, tea.Cmd) {
	return model, tea.Quit
}