package input

import (
	"testing"
	"ccpm-demo/internal/ui"
)

// MockModel implements ui.Model for testing
type MockModel struct {
	input string
}

func (m MockModel) GetInput() string    { return m.input }
func (m MockModel) GetOutput() string   { return "" }
func (m MockModel) GetError() string    { return "" }
func (m *MockModel) SetInput(input string) { m.input = input }
func (m *MockModel) SetOutput(output string) {}
func (m *MockModel) SetError(err string) {}

func TestFocusManager_BasicNavigation(t *testing.T) {
	// Create test buttons
	buttons := []*Button{
		NewButton("btn1", "1", "1", 0, 0, testAction),
		NewButton("btn2", "2", "2", 0, 1, testAction),
		NewButton("btn3", "3", "3", 1, 0, testAction),
		NewButton("btn4", "4", "4", 1, 1, testAction),
	}

	// Create focus manager
	fm := NewFocusManager()
	for _, btn := range buttons {
		fm.AddFocusable(btn)
	}

	// Test initial focus
	if !fm.SetFocus("btn1") {
		t.Error("Failed to set initial focus")
	}

	if focused := fm.GetFocusedElement(); focused == nil || focused.GetID() != "btn1" {
		t.Error("Expected btn1 to be focused")
	}

	// Test navigation right
	if !fm.Navigate("right") {
		t.Error("Failed to navigate right")
	}

	if focused := fm.GetFocusedElement(); focused == nil || focused.GetID() != "btn2" {
		t.Error("Expected btn2 to be focused after navigating right")
	}

	// Test navigation down
	if !fm.Navigate("down") {
		t.Error("Failed to navigate down")
	}

	if focused := fm.GetFocusedElement(); focused == nil || focused.GetID() != "btn4" {
		t.Error("Expected btn4 to be focused after navigating down")
	}

	// Test navigation left
	if !fm.Navigate("left") {
		t.Error("Failed to navigate left")
	}

	if focused := fm.GetFocusedElement(); focused == nil || focused.GetID() != "btn3" {
		t.Error("Expected btn3 to be focused after navigating left")
	}

	// Test navigation up
	if !fm.Navigate("up") {
		t.Error("Failed to navigate up")
	}

	if focused := fm.GetFocusedElement(); focused == nil || focused.GetID() != "btn1" {
		t.Error("Expected btn1 to be focused after navigating up")
	}
}

func TestFocusManager_TabNavigation(t *testing.T) {
	buttons := []*Button{
		NewButton("btn1", "1", "1", 0, 0, testAction),
		NewButton("btn2", "2", "2", 0, 1, testAction),
		NewButton("btn3", "3", "3", 1, 0, testAction),
	}

	fm := NewFocusManager()
	for _, btn := range buttons {
		fm.AddFocusable(btn)
	}

	// Test tab navigation (next)
	if !fm.Navigate("next") {
		t.Error("Failed to navigate next")
	}

	if focused := fm.GetFocusedElement(); focused == nil || focused.GetID() != "btn1" {
		t.Error("Expected btn1 to be focused first")
	}

	if !fm.Navigate("next") {
		t.Error("Failed to navigate next")
	}

	if focused := fm.GetFocusedElement(); focused == nil || focused.GetID() != "btn2" {
		t.Error("Expected btn2 to be focused after tab")
	}

	if !fm.Navigate("next") {
		t.Error("Failed to navigate next")
	}

	if focused := fm.GetFocusedElement(); focused == nil || focused.GetID() != "btn3" {
		t.Error("Expected btn3 to be focused after tab")
	}

	// Test wrap-around
	if !fm.Navigate("next") {
		t.Error("Failed to navigate next")
	}

	if focused := fm.GetFocusedElement(); focused == nil || focused.GetID() != "btn1" {
		t.Error("Expected btn1 to be focused after wrap-around")
	}

	// Test previous navigation (Shift+Tab)
	if !fm.Navigate("previous") {
		t.Error("Failed to navigate previous")
	}

	if focused := fm.GetFocusedElement(); focused == nil || focused.GetID() != "btn3" {
		t.Error("Expected btn3 to be focused after Shift+Tab")
	}
}

func TestFocusManager_DisabledButtons(t *testing.T) {
	buttons := []*Button{
		NewButton("btn1", "1", "1", 0, 0, testAction),
		NewButton("btn2", "2", "2", 0, 1, testAction),
		NewButton("btn3", "3", "3", 1, 0, testAction),
	}

	// Disable btn2
	buttons[1].SetEnabled(false)

	fm := NewFocusManager()
	for _, btn := range buttons {
		fm.AddFocusable(btn)
	}

	// Test navigation skips disabled buttons
	if !fm.Navigate("next") {
		t.Error("Failed to navigate next")
	}

	if focused := fm.GetFocusedElement(); focused == nil || focused.GetID() != "btn1" {
		t.Error("Expected btn1 to be focused")
	}

	if !fm.Navigate("next") {
		t.Error("Failed to navigate next")
	}

	if focused := fm.GetFocusedElement(); focused == nil || focused.GetID() != "btn3" {
		t.Error("Expected btn3 to be focused (skipping disabled btn2)")
	}
}

func TestNavigationController_BasicIntegration(t *testing.T) {
	controller := NewNavigationController()
	nav := NewFocusNavigation(controller)

	// Create test buttons
	buttons := []*Button{
		NewButton("btn1", "1", "1", 0, 0, testAction),
		NewButton("btn2", "2", "2", 0, 1, testAction),
	}

	fm := NewFocusManager()
	for _, btn := range buttons {
		fm.AddFocusable(btn)
	}

	controller.SetFocusManager(fm)

	// Test navigation through controller
	if !nav.NavigateNext() {
		t.Error("Failed to navigate next through controller")
	}

	id, label, focused := nav.GetFocusedInfo()
	if !focused || id != "btn1" || label != "1" {
		t.Error("Expected btn1 to be focused through controller")
	}

	// Test right navigation
	if !nav.NavigateRight() {
		t.Error("Failed to navigate right through controller")
	}

	id, label, focused = nav.GetFocusedInfo()
	if !focused || id != "btn2" || label != "2" {
		t.Error("Expected btn2 to be focused after right navigation")
	}
}

func TestSetupFocusManager_CalculatorButtons(t *testing.T) {
	fm, buttons := SetupFocusManager()

	// Verify all buttons are created
	if len(buttons) != 16 {
		t.Errorf("Expected 16 calculator buttons, got %d", len(buttons))
	}

	// Verify focus manager has all buttons
	focusables := fm.GetFocusables()
	if len(focusables) != 16 {
		t.Errorf("Expected 16 focusable elements, got %d", len(focusables))
	}

	// Verify initial focus is set
	focused := fm.GetFocusedElement()
	if focused == nil {
		t.Error("Expected initial focus to be set")
	}

	// Test button actions
	var mockModel MockModel
	btn := buttons[0] // Clear button
	if btn.IsEnabled() != true {
		t.Error("Expected clear button to be enabled")
	}

	// Test button activation
	updatedModel, err := btn.Activate(mockModel)
	if err != nil {
		t.Errorf("Button activation failed: %v", err)
	}

	if updatedModel.GetInput() != "" {
		t.Error("Expected clear button to clear input")
	}
}

func TestFocusNavigation_GridDimensions(t *testing.T) {
	controller := NewNavigationController()
	nav := NewFocusNavigation(controller)

	// Create test buttons in a 3x2 grid
	buttons := []*Button{
		NewButton("btn1", "1", "1", 0, 0, testAction),
		NewButton("btn2", "2", "2", 0, 1, testAction),
		NewButton("btn3", "3", "3", 1, 0, testAction),
		NewButton("btn4", "4", "4", 1, 1, testAction),
		NewButton("btn5", "5", "5", 2, 0, testAction),
		NewButton("btn6", "6", "6", 2, 1, testAction),
	}

	fm := NewFocusManager()
	for _, btn := range buttons {
		fm.AddFocusable(btn)
	}

	controller.SetFocusManager(fm)

	rows, cols := nav.GetGridDimensions()
	if rows != 3 || cols != 2 {
		t.Errorf("Expected grid dimensions 3x2, got %dx%d", rows, cols)
	}

	// Test focus summary
	summary := nav.CreateFocusSummary()
	if summary == "" {
		t.Error("Expected focus summary to be generated")
	}
}

// testAction is a simple test action for buttons
func testAction(model ui.Model, value string) (ui.Model, error) {
	return model, nil
}