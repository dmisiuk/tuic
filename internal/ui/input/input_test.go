package input

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"ccpm-demo/internal/calculator"
	"ccpm-demo/internal/ui"
)

// Test helper to create a mock model
func createMockModel() ui.Model {
	engine := calculator.NewEngine()
	return ui.NewModel(engine)
}

// TestEventRouter_BasicFunctionality tests the basic functionality of the event router
func TestEventRouter_BasicFunctionality(t *testing.T) {
	router := NewEventRouter()

	// Test router creation
	if router == nil {
		t.Error("Expected non-nil event router")
	}

	// Test initial state
	if !router.enabled {
		t.Error("Expected router to be enabled by default")
	}

	if !router.processEvents {
		t.Error("Expected event processing to be enabled by default")
	}

	if len(router.eventQueue) != 0 {
		t.Error("Expected empty event queue initially")
	}
}

// TestEventRouter_ValidateKeyEvent tests key event validation
func TestEventRouter_ValidateKeyEvent(t *testing.T) {
	router := NewEventRouter()
	validator := NewInputValidator()
	router.AddValidator(validator)

	// Test valid number key
	keyEvent := KeyEvent{
		Action: KeyActionNumber,
		Value:  "5",
	}
	event := Event{
		Type: EventTypeKey,
		Data: keyEvent,
	}

	if !router.validateEvent(event) {
		t.Error("Expected valid number key event to pass validation")
	}

	// Test invalid number key
	keyEvent.Value = "x"
	if router.validateEvent(event) {
		t.Error("Expected invalid number key event to fail validation")
	}

	// Test valid operator key
	keyEvent.Action = KeyActionOperator
	keyEvent.Value = "+"
	if !router.validateEvent(event) {
		t.Error("Expected valid operator key event to pass validation")
	}

	// Test invalid operator key
	keyEvent.Value = "%"
	if router.validateEvent(event) {
		t.Error("Expected invalid operator key event to fail validation")
	}
}

// TestEventRouter_PriorityHandling tests event priority handling
func TestEventRouter_PriorityHandling(t *testing.T) {
	router := NewEventRouter()

	// Test quit key priority (should be critical)
	priority := router.getEventPriority(KeyActionQuit)
	if priority != PriorityCritical {
		t.Errorf("Expected PriorityCritical for quit key, got %v", priority)
	}

	// Test equals key priority (should be high)
	priority = router.getEventPriority(KeyActionEquals)
	if priority != PriorityHigh {
		t.Errorf("Expected PriorityHigh for equals key, got %v", priority)
	}

	// Test number key priority (should be normal)
	priority = router.getEventPriority(KeyActionNumber)
	if priority != PriorityNormal {
		t.Errorf("Expected PriorityNormal for number key, got %v", priority)
	}
}

// TestInputValidator_BasicValidation tests basic input validation
func TestInputValidator_BasicValidation(t *testing.T) {
	validator := NewInputValidator()

	// Test valid number input
	result := validator.ValidateNumberInput("", "5")
	if !result.IsValid {
		t.Errorf("Expected valid number input, got error: %s", result.ErrorMsg)
	}

	// Test invalid number input
	result = validator.ValidateNumberInput("", "x")
	if result.IsValid {
		t.Error("Expected invalid number input to be rejected")
	}

	// Test decimal point input
	result = validator.ValidateNumberInput("", ".")
	if !result.IsValid {
		t.Errorf("Expected valid decimal point input, got error: %s", result.ErrorMsg)
	}

	// Test multiple decimal points
	result = validator.ValidateNumberInput("123.45", ".")
	if result.IsValid {
		t.Error("Expected multiple decimal points to be rejected")
	}
}

// TestInputValidator_ExpressionValidation tests expression validation
func TestInputValidator_ExpressionValidation(t *testing.T) {
	validator := NewInputValidator()

	// Test valid expression
	result := validator.ValidateExpression("123 + 456")
	if !result.IsValid {
		t.Errorf("Expected valid expression, got error: %s", result.ErrorMsg)
	}

	// Test empty expression
	result = validator.ValidateExpression("")
	if !result.IsValid {
		t.Errorf("Expected empty expression to be valid, got error: %s", result.ErrorMsg)
	}

	// Test expression with invalid operator
	result = validator.ValidateExpression("123 % 456")
	if result.IsValid {
		t.Error("Expected expression with invalid operator to be rejected")
	}

	// Test expression with operator at end
	result = validator.ValidateExpression("123 +")
	if result.IsValid {
		t.Error("Expected expression with operator at end to be rejected")
	}
}

// TestInputValidator_ExpressionTokenization tests expression tokenization
func TestInputValidator_ExpressionTokenization(t *testing.T) {
	validator := NewInputValidator()

	// Test simple expression
	tokens := validator.tokenizeExpression("123 + 456")
	expected := []string{"123", "+", "456"}
	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
	for i, token := range tokens {
		if token != expected[i] {
			t.Errorf("Expected token %s at index %d, got %s", expected[i], i, token)
		}
	}

	// Test expression with multiple spaces
	tokens = validator.tokenizeExpression("123   +   456")
	expected = []string{"123", "+", "456"}
	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
	}

	// Test negative number
	tokens = validator.tokenizeExpression("-123 + 456")
	expected = []string{"-123", "+", "456"}
	if len(tokens) != len(expected) {
		t.Errorf("Expected %d tokens, got %d", len(expected), len(tokens))
	}
}

// TestInputSystem_Initialization tests input system initialization
func TestInputSystem_Initialization(t *testing.T) {
	system := NewInputSystem()

	// Test system creation
	if system == nil {
		t.Error("Expected non-nil input system")
	}

	// Test initial state
	if !system.isEnabled {
		t.Error("Expected input system to be enabled by default")
	}

	if !system.isProcessing {
		t.Error("Expected input system to be processing by default")
	}

	if system.currentInput != "" {
		t.Error("Expected empty current input initially")
	}

	if system.errorState != "" {
		t.Error("Expected empty error state initially")
	}

	if len(system.history) != 0 {
		t.Error("Expected empty history initially")
	}
}

// TestInputSystem_NumberInput tests number input handling
func TestInputSystem_NumberInput(t *testing.T) {
	system := NewInputSystem()
	model := createMockModel()

	// Test valid number input
	msg := NumberInputMsg{Value: "5"}
	updatedModel, err := system.handleNumberInput(model, "5")
	if err != nil {
		t.Errorf("Expected valid number input to succeed, got error: %v", err)
	}
	if updatedModel.GetInput() != "5" {
		t.Errorf("Expected input to be '5', got '%s'", updatedModel.GetInput())
	}

	// Test invalid number input
	_, err = system.handleNumberInput(model, "x")
	if err == nil {
		t.Error("Expected invalid number input to return error")
	}

	// Test multiple number inputs
	system.currentInput = "12"
	updatedModel, err = system.handleNumberInput(model, "3")
	if err != nil {
		t.Errorf("Expected valid multi-digit input to succeed, got error: %v", err)
	}
	if updatedModel.GetInput() != "123" {
		t.Errorf("Expected input to be '123', got '%s'", updatedModel.GetInput())
	}
}

// TestInputSystem_OperatorInput tests operator input handling
func TestInputSystem_OperatorInput(t *testing.T) {
	system := NewInputSystem()
	model := createMockModel()

	// Test operator with existing input
	system.currentInput = "123"
	msg := OperatorInputMsg{Operator: "+"}
	updatedModel, err := system.handleOperatorInput(model, "+")
	if err != nil {
		t.Errorf("Expected valid operator input to succeed, got error: %v", err)
	}
	expectedInput := "123 + "
	if updatedModel.GetInput() != expectedInput {
		t.Errorf("Expected input to be '%s', got '%s'", expectedInput, updatedModel.GetInput())
	}

	// Test operator without existing input
	system.currentInput = ""
	_, err = system.handleOperatorInput(model, "+")
	if err == nil {
		t.Error("Expected operator without input to return error")
	}

	// Test negative operator at start
	system.currentInput = ""
	updatedModel, err = system.handleOperatorInput(model, "-")
	if err != nil {
		t.Errorf("Expected negative operator at start to succeed, got error: %v", err)
	}
	if updatedModel.GetInput() != "-" {
		t.Errorf("Expected input to be '-', got '%s'", updatedModel.GetInput())
	}
}

// TestInputSystem_EqualsInput tests equals operation handling
func TestInputSystem_EqualsInput(t *testing.T) {
	system := NewInputSystem()
	model := createMockModel()

	// Test equals with valid expression
	system.currentInput = "123 + 456"
	updatedModel, err := system.handleEqualsInput(model)
	if err != nil {
		t.Errorf("Expected equals with valid expression to succeed, got error: %v", err)
	}
	if updatedModel.GetInput() != "" {
		t.Errorf("Expected input to be cleared after equals, got '%s'", updatedModel.GetInput())
	}
	if len(system.history) != 1 {
		t.Errorf("Expected history to contain 1 entry, got %d", len(system.history))
	}
	if system.history[0] != "123 + 456" {
		t.Errorf("Expected history entry to be '123 + 456', got '%s'", system.history[0])
	}

	// Test equals with empty input
	system.currentInput = ""
	_, err = system.handleEqualsInput(model)
	if err != nil {
		t.Errorf("Expected equals with empty input to succeed, got error: %v", err)
	}
}

// TestInputSystem_HistoryNavigation tests history navigation
func TestInputSystem_HistoryNavigation(t *testing.T) {
	system := NewInputSystem()

	// Test navigation with empty history
	_, err := system.NavigateHistory(1)
	if err == nil {
		t.Error("Expected navigation with empty history to return error")
	}

	// Add some history entries
	system.addToHistory("123 + 456")
	system.addToHistory("789 - 123")
	system.addToHistory("456 * 789")

	// Test forward navigation
	entry, err := system.NavigateHistory(1)
	if err != nil {
		t.Errorf("Expected forward navigation to succeed, got error: %v", err)
	}
	if entry != "789 - 123" {
		t.Errorf("Expected entry '789 - 123', got '%s'", entry)
	}

	// Test backward navigation
	entry, err = system.NavigateHistory(-1)
	if err != nil {
		t.Errorf("Expected backward navigation to succeed, got error: %v", err)
	}
	if entry != "123 + 456" {
		t.Errorf("Expected entry '123 + 456', got '%s'", entry)
	}

	// Test wraparound (forward from last entry)
	system.historyIndex = 2 // Last entry
	entry, err = system.NavigateHistory(1)
	if err != nil {
		t.Errorf("Expected wraparound navigation to succeed, got error: %v", err)
	}
	if entry != "123 + 456" {
		t.Errorf("Expected wraparound to first entry, got '%s'", entry)
	}
}

// TestInputSystem_EnabledState tests enabled state management
func TestInputSystem_EnabledState(t *testing.T) {
	system := NewInputSystem()

	// Test initial state
	if !system.IsEnabled() {
		t.Error("Expected input system to be enabled initially")
	}

	// Test disabling
	system.SetEnabled(false)
	if system.IsEnabled() {
		t.Error("Expected input system to be disabled after SetEnabled(false)")
	}

	// Test enabling
	system.SetEnabled(true)
	if !system.IsEnabled() {
		t.Error("Expected input system to be enabled after SetEnabled(true)")
	}
}

// TestInputSystem_ProcessingState tests processing state management
func TestInputSystem_ProcessingState(t *testing.T) {
	system := NewInputSystem()

	// Test initial state
	if !system.IsProcessing() {
		t.Error("Expected input system to be processing initially")
	}

	// Test disabling processing
	system.SetProcessing(false)
	if system.IsProcessing() {
		t.Error("Expected input system to not be processing after SetProcessing(false)")
	}

	// Test enabling processing
	system.SetProcessing(true)
	if !system.IsProcessing() {
		t.Error("Expected input system to be processing after SetProcessing(true)")
	}
}

// TestInputSystem_Configuration tests system configuration
func TestInputSystem_Configuration(t *testing.T) {
	system := NewInputSystem()

	// Test default configuration
	if system.GetValidator().GetMaxInputLength() != 20 {
		t.Errorf("Expected max input length of 20, got %d", system.GetValidator().GetMaxInputLength())
	}

	// Test configuration change
	system.ConfigureValidation(30, 8, true, true)
	if system.GetValidator().GetMaxInputLength() != 30 {
		t.Errorf("Expected max input length of 30, got %d", system.GetValidator().GetMaxInputLength())
	}
	if system.GetValidator().GetMaxDecimalPlaces() != 8 {
		t.Errorf("Expected max decimal places of 8, got %d", system.GetValidator().GetMaxDecimalPlaces())
	}
}

// TestInputSystem_ButtonRegistration tests button registration
func TestInputSystem_ButtonRegistration(t *testing.T) {
	system := NewInputSystem()

	// Test button registration
	system.RegisterButton("btn1", 0, 0, 10, 5, "number", "5")
	system.RegisterButton("btn2", 12, 0, 10, 5, "operator", "+")

	// Test that buttons are registered
	mouseHandler := system.GetRouter().GetMouseHandler()
	if len(mouseHandler.buttonActions) != 2 {
		t.Errorf("Expected 2 button actions, got %d", len(mouseHandler.buttonActions))
	}

	// Test button unregistration
	system.UnregisterButton("btn1")
	if len(mouseHandler.buttonActions) != 1 {
		t.Errorf("Expected 1 button action after unregister, got %d", len(mouseHandler.buttonActions))
	}

	// Test clear all buttons
	system.ClearButtons()
	if len(mouseHandler.buttonActions) != 0 {
		t.Errorf("Expected 0 button actions after clear, got %d", len(mouseHandler.buttonActions))
	}
}

// TestInputSystem_Reset tests system reset functionality
func TestInputSystem_Reset(t *testing.T) {
	system := NewInputSystem()

	// Set up some state
	system.currentInput = "123 + 456"
	system.errorState = "Test error"
	system.addToHistory("test expression")
	system.historyIndex = 0

	// Reset the system
	system.Reset()

	// Test that state is reset
	if system.currentInput != "" {
		t.Errorf("Expected empty current input after reset, got '%s'", system.currentInput)
	}
	if system.errorState != "" {
		t.Errorf("Expected empty error state after reset, got '%s'", system.errorState)
	}
	if len(system.history) != 0 {
		t.Errorf("Expected empty history after reset, got %d entries", len(system.history))
	}
	if system.historyIndex != -1 {
		t.Errorf("Expected history index -1 after reset, got %d", system.historyIndex)
	}
}

// TestInputSystem_Integration tests the full integration with model
func TestInputSystem_Integration(t *testing.T) {
	system := NewInputSystem()
	model := createMockModel()

	// Test number input integration
	msg := NumberInputMsg{Value: "1"}
	updatedModel, _ := system.ProcessMessage(model, msg)
	if updatedModel.GetInput() != "1" {
		t.Errorf("Expected input '1' after number input, got '%s'", updatedModel.GetInput())
	}

	// Test operator input integration
	msg = OperatorInputMsg{Operator: "+"}
	updatedModel, _ = system.ProcessMessage(updatedModel, msg)
	expectedInput := "1 + "
	if updatedModel.GetInput() != expectedInput {
		t.Errorf("Expected input '%s' after operator input, got '%s'", expectedInput, updatedModel.GetInput())
	}

	// Test backspace input integration
	msg = BackspaceInputMsg{}
	updatedModel, _ = system.ProcessMessage(updatedModel, msg)
	if updatedModel.GetInput() != "1 " {
		t.Errorf("Expected input '1 ' after backspace, got '%s'", updatedModel.GetInput())
	}

	// Test clear input integration
	msg = ClearInputMsg{}
	updatedModel, _ = system.ProcessMessage(updatedModel, msg)
	if updatedModel.GetInput() != "" {
		t.Errorf("Expected empty input after clear, got '%s'", updatedModel.GetInput())
	}
}

// TestInputSystem_Performance tests performance requirements
func TestInputSystem_Performance(t *testing.T) {
	system := NewInputSystem()
	model := createMockModel()

	// Test input processing performance
	start := time.Now()
	for i := 0; i < 1000; i++ {
		msg := NumberInputMsg{Value: "1"}
		system.ProcessMessage(model, msg)
	}
	duration := time.Since(start)

	// Should process 1000 inputs in less than 100ms (requirement: <100ms per input)
	if duration > 100*time.Millisecond {
		t.Errorf("Input processing too slow: %v for 1000 inputs", duration)
	}

	// Test validation performance
	validator := NewInputValidator()
	start = time.Now()
	for i := 0; i < 1000; i++ {
		validator.ValidateExpression("123 + 456")
	}
	duration = time.Since(start)

	// Should validate 1000 expressions in less than 100ms
	if duration > 100*time.Millisecond {
		t.Errorf("Validation too slow: %v for 1000 validations", duration)
	}
}

// Benchmark tests
func BenchmarkInputSystem_ProcessMessage(b *testing.B) {
	system := NewInputSystem()
	model := createMockModel()
	msg := NumberInputMsg{Value: "1"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		system.ProcessMessage(model, msg)
	}
}

func BenchmarkInputValidator_ValidateExpression(b *testing.B) {
	validator := NewInputValidator()
	expression := "123.456 + 789.123 - 456.789"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateExpression(expression)
	}
}

func BenchmarkEventRouter_ValidateEvent(b *testing.B) {
	router := NewEventRouter()
	validator := NewInputValidator()
	router.AddValidator(validator)

	event := Event{
		Type: EventTypeKey,
		Data: KeyEvent{
			Action: KeyActionNumber,
			Value:  "5",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.validateEvent(event)
	}
}