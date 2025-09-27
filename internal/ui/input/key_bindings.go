package input

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// KeyBinding represents a mapping from a key to an action
type KeyBinding struct {
	Key         tea.KeyType
	Modifiers   tea.KeyModifier
	Action      KeyAction
	Value       string
	Description string
}

// KeyBindingsConfig holds the configuration for all key bindings
type KeyBindingsConfig struct {
	Bindings []KeyBinding
}

// DefaultKeyBindings returns the default key bindings configuration
func DefaultKeyBindings() *KeyBindingsConfig {
	return &KeyBindingsConfig{
		Bindings: []KeyBinding{
			// Number keys (0-9) - direct input
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionNumber, Value: "0", Description: "Input number 0"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionNumber, Value: "1", Description: "Input number 1"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionNumber, Value: "2", Description: "Input number 2"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionNumber, Value: "3", Description: "Input number 3"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionNumber, Value: "4", Description: "Input number 4"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionNumber, Value: "5", Description: "Input number 5"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionNumber, Value: "6", Description: "Input number 6"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionNumber, Value: "7", Description: "Input number 7"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionNumber, Value: "8", Description: "Input number 8"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionNumber, Value: "9", Description: "Input number 9"},

			// Decimal point
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionNumber, Value: ".", Description: "Decimal point"},

			// Operators
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionOperator, Value: "+", Description: "Addition operator"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionOperator, Value: "-", Description: "Subtraction operator"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionOperator, Value: "*", Description: "Multiplication operator"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionOperator, Value: "/", Description: "Division operator"},

			// Equals
			{Key: tea.KeyEnter, Modifiers: tea.ModNone, Action: KeyActionEquals, Value: "=", Description: "Calculate result"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionEquals, Value: "=", Description: "Calculate result"},

			// Clear operations
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionClear, Value: "c", Description: "Clear input"},
			{Key: tea.KeyRunes, Modifiers: tea.ModShift, Action: KeyActionClear, Value: "C", Description: "Clear input"},
			{Key: tea.KeyRunes, Modifiers: tea.ModNone, Action: KeyActionClear, Value: "C", Description: "Clear input"},

			// Backspace/Delete
			{Key: tea.KeyBackspace, Modifiers: tea.ModNone, Action: KeyActionBackspace, Value: "backspace", Description: "Delete previous character"},
			{Key: tea.KeyDelete, Modifiers: tea.ModNone, Action: KeyActionBackspace, Value: "delete", Description: "Delete next character"},

			// Navigation keys
			{Key: tea.KeyTab, Modifiers: tea.ModNone, Action: KeyActionNavigate, Value: "tab", Description: "Next focus"},
			{Key: tea.KeyTab, Modifiers: tea.ModShift, Action: KeyActionNavigate, Value: "shift_tab", Description: "Previous focus"},
			{Key: tea.KeyUp, Modifiers: tea.ModNone, Action: KeyActionNavigate, Value: "up", Description: "Navigate up"},
			{Key: tea.KeyDown, Modifiers: tea.ModNone, Action: KeyActionNavigate, Value: "down", Description: "Navigate down"},
			{Key: tea.KeyLeft, Modifiers: tea.ModNone, Action: KeyActionNavigate, Value: "left", Description: "Navigate left"},
			{Key: tea.KeyRight, Modifiers: tea.ModNone, Action: KeyActionNavigate, Value: "right", Description: "Navigate right"},

			// Focus activation
			{Key: tea.KeySpace, Modifiers: tea.ModNone, Action: KeyActionFocusActivate, Value: "space", Description: "Activate focused button"},

			// Quit
			{Key: tea.KeyEsc, Modifiers: tea.ModNone, Action: KeyActionQuit, Value: "quit", Description: "Quit application"},
			{Key: tea.KeyCtrlC, Modifiers: tea.ModNone, Action: KeyActionQuit, Value: "quit", Description: "Quit application"},
		},
	}
}

// KeyBindingManager manages key bindings and provides lookup functionality
type KeyBindingManager struct {
	config *KeyBindingsConfig
}

// NewKeyBindingManager creates a new key binding manager
func NewKeyBindingManager(config *KeyBindingsConfig) *KeyBindingManager {
	if config == nil {
		config = DefaultKeyBindings()
	}
	return &KeyBindingManager{
		config: config,
	}
}

// GetActionForKey returns the key binding for a given key message
func (kbm *KeyBindingManager) GetActionForKey(msg tea.KeyMsg) *KeyBinding {
	for _, binding := range kbm.config.Bindings {
		if binding.Key == msg.Type && (binding.Modifiers == tea.ModNone || binding.Modifiers == msg.Modifiers) {
			// For rune-based keys, we need to check the rune value
			if binding.Key == tea.KeyRunes && len(msg.Runes) > 0 {
				if string(msg.Runes) == binding.Value {
					return &binding
				}
			} else if binding.Key != tea.KeyRunes {
				return &binding
			}
		}
	}
	return nil
}

// GetAllBindings returns all key bindings
func (kbm *KeyBindingManager) GetAllBindings() []KeyBinding {
	return kbm.config.Bindings
}

// GetBindingsByAction returns all bindings for a specific action
func (kbm *KeyBindingManager) GetBindingsByAction(action KeyAction) []KeyBinding {
	var bindings []KeyBinding
	for _, binding := range kbm.config.Bindings {
		if binding.Action == action {
			bindings = append(bindings, binding)
		}
	}
	return bindings
}

// IsDirectKey returns true if the key should be processed regardless of focus
func (kbm *KeyBindingManager) IsDirectKey(msg tea.KeyMsg) bool {
	binding := kbm.GetActionForKey(msg)
	if binding == nil {
		return false
	}

	// Numbers, decimal, equals, backspace, and quit are direct keys
	switch binding.Action {
	case KeyActionNumber, KeyActionEquals, KeyActionBackspace, KeyActionQuit:
		return true
	default:
		return false
	}
}

// GetKeyDescription returns a human-readable description for a key binding
func (kbm *KeyBindingManager) GetKeyDescription(binding KeyBinding) string {
	var keyName string

	switch binding.Key {
	case tea.KeyRunes:
		keyName = binding.Value
	case tea.KeyEnter:
		keyName = "Enter"
	case tea.KeyTab:
		if binding.Modifiers == tea.ModShift {
			keyName = "Shift+Tab"
		} else {
			keyName = "Tab"
		}
	case tea.KeyBackspace:
		keyName = "Backspace"
	case tea.KeyDelete:
		keyName = "Delete"
	case tea.KeyUp:
		keyName = "↑"
	case tea.KeyDown:
		keyName = "↓"
	case tea.KeyLeft:
		keyName = "←"
	case tea.KeyRight:
		keyName = "→"
	case tea.KeySpace:
		keyName = "Space"
	case tea.KeyEsc:
		keyName = "Esc"
	case tea.KeyCtrlC:
		keyName = "Ctrl+C"
	default:
		keyName = "Unknown"
	}

	return keyName
}

// CreateHelpText generates help text for all key bindings
func (kbm *KeyBindingManager) CreateHelpText() string {
	helpText := "Keyboard Shortcuts:\n\n"

	// Group by action type
	actionGroups := make(map[KeyAction][]KeyBinding)
	for _, binding := range kbm.config.Bindings {
		actionGroups[binding.Action] = append(actionGroups[binding.Action], binding)
	}

	// Define action order for display
	actionOrder := []KeyAction{
		KeyActionNumber,
		KeyActionOperator,
		KeyActionEquals,
		KeyActionClear,
		KeyActionBackspace,
		KeyActionNavigate,
		KeyActionFocusActivate,
		KeyActionQuit,
	}

	for _, action := range actionOrder {
		if bindings, exists := actionGroups[action]; exists {
			helpText += kbm.getActionTitle(action) + ":\n"
			for _, binding := range bindings {
				keyDesc := kbm.GetKeyDescription(binding)
				helpText += fmt.Sprintf("  %-10s - %s\n", keyDesc, binding.Description)
			}
			helpText += "\n"
		}
	}

	return helpText
}

// getActionTitle returns a human-readable title for an action
func (kbm *KeyBindingManager) getActionTitle(action KeyAction) string {
	switch action {
	case KeyActionNumber:
		return "Numbers"
	case KeyActionOperator:
		return "Operators"
	case KeyActionEquals:
		return "Calculation"
	case KeyActionClear:
		return "Clear"
	case KeyActionBackspace:
		return "Delete"
	case KeyActionNavigate:
		return "Navigation"
	case KeyActionFocusActivate:
		return "Activation"
	case KeyActionQuit:
		return "Application Control"
	default:
		return "Other"
	}
}