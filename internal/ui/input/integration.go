package input

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"ccpm-demo/internal/ui"
)

// InputSystem integrates all input components into a unified system
type InputSystem struct {
	router       *EventRouter
	validator    *InputValidator
	isEnabled    bool
	isProcessing bool

	// Integration state
	currentInput string
	errorState   string
	history      []string
	historyIndex int
}

// NewInputSystem creates a new integrated input system
func NewInputSystem() *InputSystem {
	system := &InputSystem{
		router:       NewEventRouter(),
		validator:    NewInputValidator(),
		isEnabled:    true,
		isProcessing: true,
		currentInput: "",
		errorState:   "",
		history:      []string{},
		historyIndex: -1,
	}

	// Register the validator with the router
	system.router.AddValidator(system.validator)

	return system
}

// Initialize sets up the input system with default configuration
func (is *InputSystem) Initialize() {
	// Configure default validation rules
	is.validator.SetMaxInputLength(20)
	is.validator.SetMaxDecimalPlaces(6)
	is.validator.SetAllowNegative(true)
	is.validator.SetAllowOperators(true)

	// Enable all components by default
	is.SetEnabled(true)
}

// ProcessMessage processes a tea.Msg through the integrated input system
func (is *InputSystem) ProcessMessage(model ui.Model, msg tea.Msg) (ui.Model, tea.Cmd) {
	if !is.isEnabled {
		return model, nil
	}

	var command tea.Cmd
	var err error

	// Update current input state from model
	is.currentInput = model.GetInput()
	is.errorState = model.GetError()

	// Process the message through the router
	updatedModel, cmd := is.router.ProcessMessage(model, msg)
	model = updatedModel
	command = cmd

	// Handle calculator-specific message types
	switch m := msg.(type) {
	case NumberInputMsg:
		model, err = is.handleNumberInput(model, m.Value)
	case OperatorInputMsg:
		model, err = is.handleOperatorInput(model, m.Operator)
	case EqualsInputMsg:
		model, err = is.handleEqualsInput(model)
	case ClearInputMsg:
		model, err = is.handleClearInput(model)
	case BackspaceInputMsg:
		model, err = is.handleBackspaceInput(model)
	}

	// Update error state if there was an error
	if err != nil {
		model.SetError(err.Error())
		is.errorState = err.Error()
	} else {
		model.ClearError()
		is.errorState = ""
	}

	// Update current input state
	is.currentInput = model.GetInput()

	return model, command
}

// handleNumberInput handles number input from both keyboard and mouse
func (is *InputSystem) handleNumberInput(model ui.Model, value string) (ui.Model, error) {
	// Validate the number input
	result := is.validator.ValidateNumberInput(is.currentInput, value)
	if !result.IsValid {
		return model, fmt.Errorf(result.ErrorMsg)
	}

	// Add the number to the current input
	newInput := is.currentInput + value

	// Validate the complete expression
	expressionResult := is.validator.ValidateExpression(newInput)
	if !expressionResult.IsValid {
		return model, fmt.Errorf(expressionResult.ErrorMsg)
	}

	// Update the model with the new input
	model.SetInput(expressionResult.Sanitized)
	return model, nil
}

// handleOperatorInput handles operator input from both keyboard and mouse
func (is *InputSystem) handleOperatorInput(model ui.Model, operator string) (ui.Model, error) {
	// Validate the operator
	if !is.validator.validateOperatorInput(operator) {
		return model, fmt.Errorf(is.validator.GetValidationError())
	}

	// Check if we have a valid current input
	if is.currentInput == "" {
		// If no input, allow negative operator as first character
		if operator == "-" {
			model.SetInput("-")
			return model, nil
		}
		return model, fmt.Errorf("Cannot start with operator")
	}

	// Add operator with proper spacing
	newInput := is.currentInput + " " + operator + " "

	// Validate the complete expression
	expressionResult := is.validator.ValidateExpression(newInput)
	if !expressionResult.IsValid {
		return model, fmt.Errorf(expressionResult.ErrorMsg)
	}

	// Update the model
	model.SetInput(expressionResult.Sanitized)
	return model, nil
}

// handleEqualsInput handles the equals operation
func (is *InputSystem) handleEqualsInput(model ui.Model) (ui.Model, error) {
	// Validate the current expression
	expressionResult := is.validator.ValidateExpression(is.currentInput)
	if !expressionResult.IsValid {
		return model, fmt.Errorf(expressionResult.ErrorMsg)
	}

	// Add to history before evaluation
	if is.currentInput != "" {
		is.addToHistory(is.currentInput)
	}

	// The actual calculation should be done by the calculator engine
	// For now, we'll simulate it by clearing the input
	// In a real implementation, this would call model.engine.Evaluate(expression)

	// Clear the input and prepare for new calculation
	model.SetInput("")
	model.SetOutput("") // This would be set to the calculation result

	return model, nil
}

// handleClearInput handles clear operations
func (is *InputSystem) handleClearInput(model ui.Model) (ui.Model, error) {
	// Clear is always valid
	model.SetInput("")
	model.SetOutput("")
	is.currentInput = ""
	is.errorState = ""
	return model, nil
}

// handleBackspaceInput handles backspace operations
func (is *InputSystem) handleBackspaceInput(model ui.Model) (ui.Model, error) {
	if len(is.currentInput) > 0 {
		// Remove the last character
		newInput := is.currentInput[:len(is.currentInput)-1]

		// Validate the new expression
		expressionResult := is.validator.ValidateExpression(newInput)
		if !expressionResult.IsValid {
			// If validation fails, just clear everything
			newInput = ""
		}

		model.SetInput(expressionResult.Sanitized)
		is.currentInput = expressionResult.Sanitized
	}

	return model, nil
}

// addToHistory adds an expression to the history
func (is *InputSystem) addToHistory(expression string) {
	is.history = append(is.history, expression)
	if len(is.history) > 100 { // Keep last 100 entries
		is.history = is.history[1:]
	}
	is.historyIndex = len(is.history) - 1
}

// GetHistory returns the input history
func (is *InputSystem) GetHistory() []string {
	return is.history
}

// GetHistoryEntry returns a specific history entry
func (is *InputSystem) GetHistoryEntry(index int) (string, error) {
	if index < 0 || index >= len(is.history) {
		return "", fmt.Errorf("Invalid history index")
	}
	return is.history[index], nil
}

// GetCurrentHistoryIndex returns the current history index
func (is *InputSystem) GetCurrentHistoryIndex() int {
	return is.historyIndex
}

// NavigateHistory navigates through the history
func (is *InputSystem) NavigateHistory(direction int) (string, error) {
	if len(is.history) == 0 {
		return "", fmt.Errorf("No history available")
	}

	newIndex := is.historyIndex + direction

	// Clamp the index to valid range
	if newIndex < 0 {
		newIndex = len(is.history) - 1 // Wrap to end
	} else if newIndex >= len(is.history) {
		newIndex = 0 // Wrap to beginning
	}

	is.historyIndex = newIndex
	return is.history[newIndex], nil
}

// SetEnabled enables or disables the input system
func (is *InputSystem) SetEnabled(enabled bool) {
	is.isEnabled = enabled
	is.router.SetEnabled(enabled)
}

// IsEnabled returns whether the input system is enabled
func (is *InputSystem) IsEnabled() bool {
	return is.isEnabled
}

// SetProcessing enables or disables event processing
func (is *InputSystem) SetProcessing(enabled bool) {
	is.isProcessing = enabled
	is.router.SetEventProcessing(enabled)
}

// IsProcessing returns whether the input system is processing events
func (is *InputSystem) IsProcessing() bool {
	return is.isProcessing
}

// GetCurrentInput returns the current input string
func (is *InputSystem) GetCurrentInput() string {
	return is.currentInput
}

// GetErrorState returns the current error state
func (is *InputSystem) GetErrorState() string {
	return is.errorState
}

// ClearError clears the error state
func (is *InputSystem) ClearError() {
	is.errorState = ""
}

// GetRouter returns the event router
func (is *InputSystem) GetRouter() *EventRouter {
	return is.router
}

// GetValidator returns the input validator
func (is *InputSystem) GetValidator() *InputValidator {
	return is.validator
}

// ConfigureValidation configures the validation settings
func (is *InputSystem) ConfigureValidation(maxLength, maxDecimal int, allowNegative, allowOperators bool) {
	is.validator.SetMaxInputLength(maxLength)
	is.validator.SetMaxDecimalPlaces(maxDecimal)
	is.validator.SetAllowNegative(allowNegative)
	is.validator.SetAllowOperators(allowOperators)
}

// RegisterButton registers a button with the mouse handler
func (is *InputSystem) RegisterButton(buttonID string, x, y, width, height int, actionType, actionValue string) {
	action := ButtonAction{
		Type:  actionType,
		Value: actionValue,
		Handler: func() tea.Msg {
			switch actionType {
			case "number":
				return NumberInputMsg{Value: actionValue}
			case "operator":
				return OperatorInputMsg{Operator: actionValue}
			case "equals":
				return EqualsInputMsg{}
			case "clear":
				return ClearInputMsg{}
			case "backspace":
				return BackspaceInputMsg{}
			default:
				return nil
			}
		},
	}

	is.router.GetMouseHandler().RegisterButton(buttonID, x, y, width, height, action)
}

// UnregisterButton unregisters a button from the mouse handler
func (is *InputSystem) UnregisterButton(buttonID string) {
	is.router.GetMouseHandler().UnregisterButtonAction(buttonID)
}

// ClearButtons clears all registered buttons
func (is *InputSystem) ClearButtons() {
	is.router.GetMouseHandler().ClearButtons()
}

// GetKeyBindings returns the current key bindings
func (is *InputSystem) GetKeyBindings() []KeyBinding {
	return is.router.GetKeyHandler().keyBindingManager.GetAllBindings()
}

// AddCustomValidator adds a custom validator to the router
func (is *InputSystem) AddCustomValidator(validator *InputValidator) {
	is.router.AddValidator(validator)
}

// RemoveCustomValidator removes a custom validator from the router
func (is *InputSystem) RemoveCustomValidator(validator *InputValidator) {
	is.router.RemoveValidator(validator)
}

// Reset resets the input system to its initial state
func (is *InputSystem) Reset() {
	is.currentInput = ""
	is.errorState = ""
	is.history = []string{}
	is.historyIndex = -1
	is.router.ClearEventQueue()
	is.router.GetMouseHandler().Reset()
}

// GetSystemState returns the current state of the input system
func (is *InputSystem) GetSystemState() map[string]interface{} {
	return map[string]interface{}{
		"enabled":       is.isEnabled,
		"processing":    is.isProcessing,
		"currentInput":  is.currentInput,
		"errorState":    is.errorState,
		"historyCount":  len(is.history),
		"historyIndex":  is.historyIndex,
		"eventQueueLen": len(is.router.GetEventQueue()),
	}
}

// ValidateCurrentInput validates the current input expression
func (is *InputSystem) ValidateCurrentInput() ValidationResult {
	return is.validator.ValidateExpression(is.currentInput)
}

// SanitizeCurrentInput sanitizes the current input
func (is *InputSystem) SanitizeCurrentInput() string {
	return is.validator.SanitizeInput(is.currentInput)
}