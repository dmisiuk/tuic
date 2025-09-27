package components

import "fmt"

// ButtonState represents the visual and interactive state of a button
type ButtonState int

const (
	// StateNormal represents the default, non-interactive button state
	StateNormal ButtonState = iota

	// StateFocused represents when the button has keyboard focus
	StateFocused

	// StatePressed represents when the button is being pressed/clicked
	StatePressed

	// StateDisabled represents when the button is non-interactive
	StateDisabled
)

// String returns a string representation of the button state
func (s ButtonState) String() string {
	switch s {
	case StateNormal:
		return "normal"
	case StateFocused:
		return "focused"
	case StatePressed:
		return "pressed"
	case StateDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}

// ButtonType categorizes buttons by their function and styling
type ButtonType int

const (
	// TypeNumber represents numeric buttons (0-9, decimal point)
	TypeNumber ButtonType = iota

	// TypeOperator represents arithmetic operation buttons (+, -, *, /)
	TypeOperator

	// TypeSpecial represents special function buttons (C, CE, =, etc.)
	TypeSpecial
)

// String returns a string representation of the button type
func (t ButtonType) String() string {
	switch t {
	case TypeNumber:
		return "number"
	case TypeOperator:
		return "operator"
	case TypeSpecial:
		return "special"
	default:
		return "unknown"
	}
}

// ButtonConfig contains the configuration for a button
type ButtonConfig struct {
	// Label is the text displayed on the button
	Label string

	// Type categorizes the button for styling purposes
	Type ButtonType

	// Value is the actual value/action the button represents
	Value string

	// Width is the display width of the button (in characters)
	Width int

	// Height is the display height of the button (in lines)
	Height int

	// Position represents the button's location in a grid
	Position Position
}

// Position represents the grid position of a button
type Position struct {
	Row    int
	Column int
}

// ButtonStyle defines the visual appearance for different button states and types
type ButtonStyle struct {
	Normal   string
	Focused  string
	Pressed  string
	Disabled string
}

// ButtonStateManager manages the state transitions and validation for buttons
type ButtonStateManager struct {
	currentState ButtonState
	buttonType   ButtonType
	config       ButtonConfig
}

// NewButtonStateManager creates a new state manager for a button
func NewButtonStateManager(buttonType ButtonType, config ButtonConfig) *ButtonStateManager {
	return &ButtonStateManager{
		currentState: StateNormal,
		buttonType:   buttonType,
		config:       config,
	}
}

// State returns the current button state
func (sm *ButtonStateManager) State() ButtonState {
	return sm.currentState
}

// SetState attempts to transition to a new state
func (sm *ButtonStateManager) SetState(newState ButtonState) error {
	if !sm.isValidTransition(newState) {
		return &InvalidStateTransitionError{
			From: sm.currentState,
			To:   newState,
		}
	}

	sm.currentState = newState
	return nil
}

// Focus sets the button to focused state
func (sm *ButtonStateManager) Focus() error {
	return sm.SetState(StateFocused)
}

// Press sets the button to pressed state
func (sm *ButtonStateManager) Press() error {
	return sm.SetState(StatePressed)
}

// Release returns the button to normal state from pressed
func (sm *ButtonStateManager) Release() error {
	if sm.currentState == StatePressed {
		return sm.SetState(StateNormal)
	}
	return nil
}

// Blur returns the button to normal state from focused
func (sm *ButtonStateManager) Blur() error {
	if sm.currentState == StateFocused {
		return sm.SetState(StateNormal)
	}
	return nil
}

// Disable sets the button to disabled state
func (sm *ButtonStateManager) Disable() error {
	return sm.SetState(StateDisabled)
}

// Enable returns the button to normal state from disabled
func (sm *ButtonStateManager) Enable() error {
	if sm.currentState == StateDisabled {
		return sm.SetState(StateNormal)
	}
	return nil
}

// isValidTransition checks if a state transition is valid
func (sm *ButtonStateManager) isValidTransition(newState ButtonState) bool {
	// Define valid transitions based on current state
	switch sm.currentState {
	case StateNormal:
		return newState == StateFocused || newState == StatePressed || newState == StateDisabled

	case StateFocused:
		return newState == StateNormal || newState == StatePressed || newState == StateDisabled

	case StatePressed:
		return newState == StateNormal || newState == StateFocused || newState == StateDisabled

	case StateDisabled:
		return newState == StateNormal || newState == StateFocused

	default:
		return false
	}
}

// IsInteractive returns true if the button is in an interactive state
func (sm *ButtonStateManager) IsInteractive() bool {
	return sm.currentState != StateDisabled
}

// IsFocused returns true if the button has focus
func (sm *ButtonStateManager) IsFocused() bool {
	return sm.currentState == StateFocused
}

// IsPressed returns true if the button is being pressed
func (sm *ButtonStateManager) IsPressed() bool {
	return sm.currentState == StatePressed
}

// GetType returns the button type
func (sm *ButtonStateManager) GetType() ButtonType {
	return sm.buttonType
}

// GetConfig returns the button configuration
func (sm *ButtonStateManager) GetConfig() ButtonConfig {
	return sm.config
}

// InvalidStateTransitionError represents an invalid state transition
type InvalidStateTransitionError struct {
	From ButtonState
	To   ButtonState
}

// Error implements the error interface
func (e *InvalidStateTransitionError) Error() string {
	return fmt.Sprintf("invalid state transition from %s to %s", e.From, e.To)
}

// Is compares error types for error handling
func (e *InvalidStateTransitionError) Is(target error) bool {
	_, ok := target.(*InvalidStateTransitionError)
	return ok
}