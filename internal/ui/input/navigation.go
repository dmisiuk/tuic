package input

import (
	"fmt"
	"ccpm-demo/internal/ui"
)

// Button represents a calculator button that can be focused
type Button struct {
	ID          string
	Label       string
	Value       string
	Row         int
	Col         int
	Enabled     bool
	Focused     bool
	Action      func(ui.Model, string) (ui.Model, error)
}

// NewButton creates a new calculator button
func NewButton(id, label, value string, row, col int, action func(ui.Model, string) (ui.Model, error)) *Button {
	return &Button{
		ID:      id,
		Label:   label,
		Value:   value,
		Row:     row,
		Col:     col,
		Enabled: true,
		Focused: false,
		Action:  action,
	}
}

// GetID returns the button's unique identifier
func (b *Button) GetID() string {
	return b.ID
}

// GetPosition returns the button's position in the grid
func (b *Button) GetPosition() (row, col int) {
	return b.Row, b.Col
}

// GetLabel returns the button's display label
func (b *Button) GetLabel() string {
	return b.Label
}

// IsEnabled returns true if the button can receive focus
func (b *Button) IsEnabled() bool {
	return b.Enabled
}

// OnFocus is called when the button receives focus
func (b *Button) OnFocus() {
	b.Focused = true
}

// OnBlur is called when the button loses focus
func (b *Button) OnBlur() {
	b.Focused = false
}

// Activate performs the button's action
func (b *Button) Activate(model ui.Model) (ui.Model, error) {
	if b.Action != nil {
		return b.Action(model, b.Value)
	}
	return model, nil
}

// SetEnabled sets whether the button is enabled
func (b *Button) SetEnabled(enabled bool) {
	b.Enabled = enabled
}

// NavigationController provides high-level navigation control
type NavigationController struct {
	focusManager *FocusManager
}

// NewNavigationController creates a new navigation controller
func NewNavigationController() *NavigationController {
	return &NavigationController{
		focusManager: NewFocusManager(),
	}
}

// SetFocusManager sets the focus manager for navigation
func (nc *NavigationController) SetFocusManager(fm *FocusManager) {
	nc.focusManager = fm
}

// GetFocusManager returns the current focus manager
func (nc *NavigationController) GetFocusManager() *FocusManager {
	return nc.focusManager
}

// Navigate moves focus in the specified direction
func (nc *NavigationController) Navigate(direction string) bool {
	if nc.focusManager == nil {
		return false
	}
	return nc.focusManager.Navigate(direction)
}

// Activate activates the currently focused element
func (nc *NavigationController) Activate(model ui.Model) (ui.Model, error) {
	if nc.focusManager == nil {
		return model, nil
	}
	return nc.focusManager.Activate(model)
}

// GetFocusedElement returns the currently focused element
func (nc *NavigationController) GetFocusedElement() Focusable {
	if nc.focusManager == nil {
		return nil
	}
	return nc.focusManager.GetFocusedElement()
}

// GetFocusedID returns the ID of the currently focused element
func (nc *NavigationController) GetFocusedID() string {
	if nc.focusManager == nil {
		return ""
	}
	return nc.focusManager.GetFocusedID()
}

// CreateCalculatorButtons creates the standard calculator button grid
func CreateCalculatorButtons() []*Button {
	buttons := make([]*Button, 0)

	// Define button actions
	numberAction := func(model ui.Model, value string) (ui.Model, error) {
		// Handle number input
		input := model.GetInput()
		model.SetInput(input + value)
		return model, nil
	}

	operatorAction := func(model ui.Model, value string) (ui.Model, error) {
		// Handle operator input
		input := model.GetInput()
		if input != "" {
			model.SetInput(input + " " + value + " ")
		}
		return model, nil
	}

	equalsAction := func(model ui.Model, value string) (ui.Model, error) {
		// Handle equals/calculate
		// For now, just clear the input - will be integrated with calculator engine
		model.SetInput("")
		return model, nil
	}

	clearAction := func(model ui.Model, value string) (ui.Model, error) {
		// Handle clear
		model.SetInput("")
		return model, nil
	}

	backspaceAction := func(model ui.Model, value string) (ui.Model, error) {
		// Handle backspace
		input := model.GetInput()
		if len(input) > 0 {
			model.SetInput(input[:len(input)-1])
		}
		return model, nil
	}

	// Row 0: Clear and basic operations
	buttons = append(buttons, NewButton("clear", "C", "clear", 0, 0, clearAction))
	buttons = append(buttons, NewButton("paren_left", "(", "(", 0, 1, operatorAction))
	buttons = append(buttons, NewButton("paren_right", ")", ")", 0, 2, operatorAction))
	buttons = append(buttons, NewButton("backspace", "←", "backspace", 0, 3, backspaceAction))

	// Row 1: Numbers 7-9 and division
	buttons = append(buttons, NewButton("seven", "7", "7", 1, 0, numberAction))
	buttons = append(buttons, NewButton("eight", "8", "8", 1, 1, numberAction))
	buttons = append(buttons, NewButton("nine", "9", "9", 1, 2, numberAction))
	buttons = append(buttons, NewButton("divide", "÷", "/", 1, 3, operatorAction))

	// Row 2: Numbers 4-6 and multiplication
	buttons = append(buttons, NewButton("four", "4", "4", 2, 0, numberAction))
	buttons = append(buttons, NewButton("five", "5", "5", 2, 1, numberAction))
	buttons = append(buttons, NewButton("six", "6", "6", 2, 2, numberAction))
	buttons = append(buttons, NewButton("multiply", "×", "*", 2, 3, operatorAction))

	// Row 3: Numbers 1-3 and subtraction
	buttons = append(buttons, NewButton("one", "1", "1", 3, 0, numberAction))
	buttons = append(buttons, NewButton("two", "2", "2", 3, 1, numberAction))
	buttons = append(buttons, NewButton("three", "3", "3", 3, 2, numberAction))
	buttons = append(buttons, NewButton("subtract", "-", "-", 3, 3, operatorAction))

	// Row 4: Decimal, zero, equals, and addition
	buttons = append(buttons, NewButton("decimal", ".", ".", 4, 0, numberAction))
	buttons = append(buttons, NewButton("zero", "0", "0", 4, 1, numberAction))
	buttons = append(buttons, NewButton("equals", "=", "=", 4, 2, equalsAction))
	buttons = append(buttons, NewButton("add", "+", "+", 4, 3, operatorAction))

	return buttons
}

// SetupFocusManager creates and configures a focus manager with calculator buttons
func SetupFocusManager() (*FocusManager, []*Button) {
	buttons := CreateCalculatorButtons()
	fm := NewFocusManager()

	// Add all buttons as focusable elements
	for _, button := range buttons {
		fm.AddFocusable(button)
	}

	// Set focus to the first button
	if len(buttons) > 0 {
		fm.SetFocus(buttons[0].GetID())
	}

	return fm, buttons
}

// FocusNavigation provides convenience methods for common navigation patterns
type FocusNavigation struct {
	controller *NavigationController
}

// NewFocusNavigation creates a new focus navigation helper
func NewFocusNavigation(controller *NavigationController) *FocusNavigation {
	return &FocusNavigation{
		controller: controller,
	}
}

// NavigateUp moves focus up
func (fn *FocusNavigation) NavigateUp() bool {
	return fn.controller.Navigate("up")
}

// NavigateDown moves focus down
func (fn *FocusNavigation) NavigateDown() bool {
	return fn.controller.Navigate("down")
}

// NavigateLeft moves focus left
func (fn *FocusNavigation) NavigateLeft() bool {
	return fn.controller.Navigate("left")
}

// NavigateRight moves focus right
func (fn *FocusNavigation) NavigateRight() bool {
	return fn.controller.Navigate("right")
}

// NavigateNext moves to next focusable element (Tab)
func (fn *FocusNavigation) NavigateNext() bool {
	return fn.controller.Navigate("next")
}

// NavigatePrevious moves to previous focusable element (Shift+Tab)
func (fn *FocusNavigation) NavigatePrevious() bool {
	return fn.controller.Navigate("previous")
}

// Activate activates the currently focused element
func (fn *FocusNavigation) Activate(model ui.Model) (ui.Model, error) {
	return fn.controller.Activate(model)
}

// GetFocusedInfo returns information about the currently focused element
func (fn *FocusNavigation) GetFocusedInfo() (id, label string, focused bool) {
	element := fn.controller.GetFocusedElement()
	if element == nil {
		return "", "", false
	}

	button, ok := element.(*Button)
	if !ok {
		return element.GetID(), element.GetLabel(), true
	}

	return button.ID, button.Label, button.Focused
}

// PrintFocusState prints the current focus state for debugging
func (fn *FocusNavigation) PrintFocusState() {
	id, label, focused := fn.GetFocusedInfo()
	if focused {
		fmt.Printf("Currently focused: %s (%s)\n", id, label)
	} else {
		fmt.Println("No element is currently focused")
	}
}

// GetGridDimensions returns the grid dimensions based on focusable elements
func (fn *FocusNavigation) GetGridDimensions() (rows, cols int) {
	fm := fn.controller.GetFocusManager()
	if fm == nil {
		return 0, 0
	}

	maxRow, maxCol := 0, 0
	for _, element := range fm.GetFocusables() {
		row, col := element.GetPosition()
		if row > maxRow {
			maxRow = row
		}
		if col > maxCol {
			maxCol = col
		}
	}

	return maxRow + 1, maxCol + 1
}

// CreateFocusSummary returns a string representation of the focus state
func (fn *FocusNavigation) CreateFocusSummary() string {
	fm := fn.controller.GetFocusManager()
	if fm == nil {
		return "No focus manager available"
	}

	rows, cols := fn.GetGridDimensions()
	focusedID := fn.controller.GetFocusedID()

	summary := fmt.Sprintf("Focus Grid (%dx%d):\n", rows, cols)

	// Create a grid representation
	grid := make([][]string, rows)
	for i := range grid {
		grid[i] = make([]string, cols)
		for j := range grid[i] {
			grid[i][j] = " . "
		}
	}

	// Fill in the grid with button labels
	for _, element := range fm.GetFocusables() {
		row, col := element.GetPosition()
		if row < rows && col < cols {
			label := element.GetLabel()
			if len(label) > 3 {
				label = label[:3]
			} else if len(label) < 3 {
				label = label + "   "[:3-len(label)]
			}

			if element.GetID() == focusedID {
				grid[row][col] = "[" + label + "]"
			} else {
				grid[row][col] = " " + label + " "
			}
		}
	}

	// Build the summary string
	for i := range grid {
		summary += fmt.Sprintf("Row %d: %s\n", i, "")
		for j := range grid[i] {
			summary += grid[i][j] + " "
		}
		summary += "\n"
	}

	return summary
}