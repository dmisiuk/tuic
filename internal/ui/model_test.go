package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"ccpm-demo/internal/calculator"
)

func TestNewModel(t *testing.T) {
	engine := calculator.NewEngine()
	model := NewModel(engine)

	// Test initial state
	if model.engine != engine {
		t.Error("Model should have the provided engine")
	}

	if model.calculatorState.displayValue != "0" {
		t.Errorf("Expected display value '0', got '%s'", model.calculatorState.displayValue)
	}

	if model.input != "" {
		t.Errorf("Expected empty input, got '%s'", model.input)
	}

	if model.output != "" {
		t.Errorf("Expected empty output, got '%s'", model.output)
	}

	if model.error != "" {
		t.Errorf("Expected empty error, got '%s'", model.error)
	}

	if len(model.history) != 0 {
		t.Errorf("Expected empty history, got %d entries", len(model.history))
	}

	if model.quitting {
		t.Error("Model should not be quitting initially")
	}
}

func TestModelInit(t *testing.T) {
	engine := calculator.NewEngine()
	model := NewModel(engine)

	cmd := model.Init()
	if cmd != nil {
		t.Error("Init should return nil command")
	}
}

func TestModelUpdateWindowSize(t *testing.T) {
	engine := calculator.NewEngine()
	model := NewModel(engine)

	// Test window size update
	msg := tea.WindowSizeMsg{Width: 80, Height: 40}
	updatedModel, cmd := model.Update(msg)

	if cmd != nil {
		t.Error("Window size update should return nil command")
	}

	// Type assertion correctly
	um, ok := updatedModel.(Model)
	if !ok {
		t.Error("Updated model should be of type Model")
	}

	if um.width != 80 {
		t.Errorf("Expected width 80, got %d", um.width)
	}

	if um.height != 40 {
		t.Errorf("Expected height 40, got %d", um.height)
	}

	if !um.ready {
		t.Error("Model should be ready after window size update")
	}
}

func TestModelUpdateKeyMessages(t *testing.T) {
	engine := calculator.NewEngine()
	model := NewModel(engine)

	tests := []struct {
		name     string
		msg      tea.KeyMsg
		quit     bool
		hasError bool
	}{
		{"Quit with q", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}, true, false},
		{"Quit with Ctrl+C", tea.KeyMsg{Type: tea.KeyCtrlC}, true, false},
		{"Clear with c", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}, false, false},
		{"Enter with empty input", tea.KeyMsg{Type: tea.KeyEnter}, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedModel, cmd := model.Update(tt.msg)

			um, ok := updatedModel.(Model)
			if !ok {
				t.Error("Updated model should be of type Model")
			}

			if tt.quit && !um.quitting {
				t.Error("Model should be quitting")
			}

			if !tt.quit && um.quitting {
				t.Error("Model should not be quitting")
			}

			if tt.quit && cmd == nil {
				t.Error("Quit should return tea.Quit command")
			}
		})
	}
}

func TestModelView(t *testing.T) {
	engine := calculator.NewEngine()
	model := NewModel(engine)

	// Test view rendering
	view := model.View()

	if view == "" {
		t.Error("View should not be empty")
	}

	// Check if view contains expected elements
	if !contains(view, "CCPM Calculator") {
		t.Error("View should contain 'CCPM Calculator'")
	}

	if !contains(view, "0") {
		t.Error("View should contain initial display value '0'")
	}
}

func TestModelFormatting(t *testing.T) {
	engine := calculator.NewEngine()
	model := NewModel(engine)

	tests := []struct {
		value    float64
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{1.0, "1"},
		{1.5, "1.500000"},
		{3.14159, "3.141590"},
	}

	for _, tt := range tests {
		result := model.formatValue(tt.value)
		if result != tt.expected {
			t.Errorf("formatValue(%f) = '%s', expected '%s'", tt.value, result, tt.expected)
		}
	}
}

func TestModelTruncateString(t *testing.T) {
	engine := calculator.NewEngine()
	model := NewModel(engine)

	tests := []struct {
		input    string
		width    int
		expected string
	}{
		{"short", 10, "short"},
		{"exactlyten", 10, "exactlyten"},
		{"toolongstring", 10, "toolong..."},
		{"a", 1, "a"},
		{"ab", 1, "..."},
	}

	for _, tt := range tests {
		result := model.truncateString(tt.input, tt.width)
		if result != tt.expected {
			t.Errorf("truncateString('%s', %d) = '%s', expected '%s'", tt.input, tt.width, result, tt.expected)
		}
	}
}

func TestModelHistory(t *testing.T) {
	engine := calculator.NewEngine()
	model := NewModel(engine)

	// Test adding to history
	model.addToHistory("2 + 2 = 4")
	if len(model.history) != 1 {
		t.Errorf("Expected 1 history entry, got %d", len(model.history))
	}

	if model.history[0] != "2 + 2 = 4" {
		t.Errorf("Expected '2 + 2 = 4', got '%s'", model.history[0])
	}

	if model.historyIndex != 0 {
		t.Errorf("Expected history index 0, got %d", model.historyIndex)
	}

	// Test history limit
	for i := 0; i < 105; i++ {
		model.addToHistory("test")
	}

	if len(model.history) > 100 {
		t.Errorf("History should be limited to 100 entries, got %d", len(model.history))
	}
}

func TestModelErrorHandling(t *testing.T) {
	engine := calculator.NewEngine()
	model := NewModel(engine)

	// Test setting error
	var err error = calculator.CalculatorError("test error")
	model.setError(err)

	if model.error != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", model.error)
	}

	// Test clearing error
	model.clearError()
	if model.error != "" {
		t.Errorf("Expected empty error after clear, got '%s'", model.error)
	}

	// Test setting nil error
	model.setError(nil)
	if model.error != "" {
		t.Errorf("Expected empty error for nil, got '%s'", model.error)
	}
}

func TestModelDisplayDimensions(t *testing.T) {
	engine := calculator.NewEngine()
	model := NewModel(engine)

	// Test with no dimensions set
	width := model.getDisplayWidth()
	expectedWidth := 56 // Default width minus borders
	if width != expectedWidth {
		t.Errorf("Expected display width %d, got %d", expectedWidth, width)
	}

	height := model.getDisplayHeight()
	expectedHeight := 30 // Default height minus borders
	if height != expectedHeight {
		t.Errorf("Expected display height %d, got %d", expectedHeight, height)
	}

	// Test with dimensions set
	model.width = 100
	model.height = 50

	width = model.getDisplayWidth()
	expectedWidth = 96 // 100 - 4
	if width != expectedWidth {
		t.Errorf("Expected display width %d, got %d", expectedWidth, width)
	}

	height = model.getDisplayHeight()
	expectedHeight = 46 // 50 - 4
	if height != expectedHeight {
		t.Errorf("Expected display height %d, got %d", expectedHeight, height)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}