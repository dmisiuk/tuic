package integration

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ccpm-demo/internal/ui/components"
	"ccpm-demo/internal/ui/styles"
)

// ButtonGrid represents the complete calculator button grid integration
// It combines buttons, grid layout, styling, and interaction handling
type ButtonGrid struct {
	buttons       map[string]*components.Button
	grid          *components.GridLayout
	themeManager  *styles.ThemeManager
	focusedButton string
	pressedButton string
	dimensions    GridDimensions
}

// GridDimensions defines the size of the button grid
type GridDimensions struct {
	Columns int
	Rows    int
}

// ButtonDefinition defines a button's properties
type ButtonDefinition struct {
	Label    string
	Value    string
	Type     components.ButtonType
	Row      int
	Column   int
	Width    int
	Height   int
}

// ButtonAction represents an action triggered by a button
type ButtonAction struct {
	Button   *components.Button
	Action   string
	Value    string
	ButtonID string
}

// NewButtonGrid creates a new button grid with default calculator layout
func NewButtonGrid() *ButtonGrid {
	themeManager := styles.NewThemeManager()
	grid := components.NewGridLayout()

	buttonGrid := &ButtonGrid{
		buttons:      make(map[string]*components.Button),
		grid:         grid,
		themeManager: themeManager,
		dimensions: GridDimensions{
			Columns: 4,
			Rows:    5,
		},
	}

	// Initialize the calculator button layout
	buttonGrid.initializeCalculatorLayout()

	return buttonGrid
}

// NewButtonGridWithTheme creates a new button grid with a specific theme
func NewButtonGridWithTheme(themeName string) (*ButtonGrid, error) {
	themeManager := styles.NewThemeManager()

	err := themeManager.SetTheme(themeName)
	if err != nil {
		return nil, fmt.Errorf("failed to set theme: %w", err)
	}

	grid := components.NewGridLayout()

	buttonGrid := &ButtonGrid{
		buttons:      make(map[string]*components.Button),
		grid:         grid,
		themeManager: themeManager,
		dimensions: GridDimensions{
			Columns: 4,
			Rows:    5,
		},
	}

	// Initialize the calculator button layout
	buttonGrid.initializeCalculatorLayout()

	return buttonGrid, nil
}

// initializeCalculatorLayout creates the standard calculator button arrangement
func (bg *ButtonGrid) initializeCalculatorLayout() {
	// Standard calculator button layout (4x5 grid)
	buttonDefs := []ButtonDefinition{
		// Row 0 (top row): C, CE, ←, ÷
		{Label: "C", Value: "clear", Type: components.TypeSpecial, Row: 0, Column: 0, Width: 3, Height: 1},
		{Label: "CE", Value: "clear_entry", Type: components.TypeSpecial, Row: 0, Column: 1, Width: 3, Height: 1},
		{Label: "←", Value: "backspace", Type: components.TypeSpecial, Row: 0, Column: 2, Width: 3, Height: 1},
		{Label: "÷", Value: "/", Type: components.TypeOperator, Row: 0, Column: 3, Width: 3, Height: 1},

		// Row 1: 7, 8, 9, ×
		{Label: "7", Value: "7", Type: components.TypeNumber, Row: 1, Column: 0, Width: 3, Height: 1},
		{Label: "8", Value: "8", Type: components.TypeNumber, Row: 1, Column: 1, Width: 3, Height: 1},
		{Label: "9", Value: "9", Type: components.TypeNumber, Row: 1, Column: 2, Width: 3, Height: 1},
		{Label: "×", Value: "*", Type: components.TypeOperator, Row: 1, Column: 3, Width: 3, Height: 1},

		// Row 2: 4, 5, 6, -
		{Label: "4", Value: "4", Type: components.TypeNumber, Row: 2, Column: 0, Width: 3, Height: 1},
		{Label: "5", Value: "5", Type: components.TypeNumber, Row: 2, Column: 1, Width: 3, Height: 1},
		{Label: "6", Value: "6", Type: components.TypeNumber, Row: 2, Column: 2, Width: 3, Height: 1},
		{Label: "-", Value: "-", Type: components.TypeOperator, Row: 2, Column: 3, Width: 3, Height: 1},

		// Row 3: 1, 2, 3, +
		{Label: "1", Value: "1", Type: components.TypeNumber, Row: 3, Column: 0, Width: 3, Height: 1},
		{Label: "2", Value: "2", Type: components.TypeNumber, Row: 3, Column: 1, Width: 3, Height: 1},
		{Label: "3", Value: "3", Type: components.TypeNumber, Row: 3, Column: 2, Width: 3, Height: 1},
		{Label: "+", Value: "+", Type: components.TypeOperator, Row: 3, Column: 3, Width: 3, Height: 1},

		// Row 4 (bottom row): 0, ., =
		{Label: "0", Value: "0", Type: components.TypeNumber, Row: 4, Column: 0, Width: 3, Height: 1},
		{Label: ".", Value: ".", Type: components.TypeNumber, Row: 4, Column: 1, Width: 3, Height: 1},
		{Label: "=", Value: "=", Type: components.TypeSpecial, Row: 4, Column: 2, Width: 3, Height: 1},
		// Empty cell at row 4, column 3 for balance
	}

	// Create buttons from definitions
	for _, def := range buttonDefs {
		buttonID := bg.generateButtonID(def.Row, def.Column)
		button := bg.createButton(def)
		bg.buttons[buttonID] = button

		// Add button to grid
		buttonStyle := bg.getButtonStyle(button)
		bg.grid.AddCell(def.Column, def.Row, def.Label, buttonStyle)
	}

	// Set initial focus on the first button
	if len(bg.buttons) > 0 {
		bg.focusedButton = "button_0_0"  // Focus on the "C" button
		if button, exists := bg.buttons[bg.focusedButton]; exists {
			button.Focus()
		}
	}
}

// createButton creates a button from a definition
func (bg *ButtonGrid) createButton(def ButtonDefinition) *components.Button {
	config := components.ButtonConfig{
		Label:  def.Label,
		Type:   def.Type,
		Value:  def.Value,
		Width:  def.Width,
		Height: def.Height,
		Position: components.Position{
			Row:    def.Row,
			Column: def.Column,
		},
	}

	return components.NewButton(config)
}

// generateButtonID generates a unique ID for a button based on its position
func (bg *ButtonGrid) generateButtonID(row, col int) string {
	return fmt.Sprintf("button_%d_%d", row, col)
}

// getButtonStyle returns the appropriate style for a button based on its type and state
func (bg *ButtonGrid) getButtonStyle(button *components.Button) lipgloss.Style {
	buttonType := button.GetType()
	state := button.GetState()

	var style lipgloss.Style

	switch buttonType {
	case components.TypeNumber:
		style = bg.themeManager.GetButtonStyle("number", state.String())
	case components.TypeOperator:
		style = bg.themeManager.GetButtonStyle("operator", state.String())
	case components.TypeSpecial:
		style = bg.themeManager.GetButtonStyle("special", state.String())
	default:
		style = bg.themeManager.GetButtonStyle("number", state.String())
	}

	// Apply button dimensions
	config := button.GetConfig()
	if config.Width > 0 {
		style = style.Width(config.Width)
	}
	if config.Height > 0 {
		style = style.Height(config.Height)
	}

	return style
}

// Render renders the entire button grid
func (bg *ButtonGrid) Render(termWidth int) string {
	// Update grid styling based on current theme
	bg.updateGridStyling()

	// Render the grid
	return bg.grid.Render(termWidth)
}

// updateGridStyling updates the grid layout with current theme styling
func (bg *ButtonGrid) updateGridStyling() {
	theme := bg.themeManager.GetCurrentTheme()

	// Apply grid theme styling
	bg.grid.
		WithBorderStyle(theme.Styles.Grid.Container).
		WithFocusedStyle(theme.Styles.Grid.CellFocused).
		WithPressedStyle(theme.Styles.Grid.CellPressed)
}

// HandleKeyPress handles keyboard input for button navigation and activation
func (bg *ButtonGrid) HandleKeyPress(msg tea.KeyMsg) *ButtonAction {
	switch msg.Type {
	case tea.KeyEnter, tea.KeySpace:
		// Activate focused button
		if bg.focusedButton != "" {
			return bg.activateButton(bg.focusedButton)
		}

	case tea.KeyUp, tea.KeyDown, tea.KeyLeft, tea.KeyRight:
		// Navigate between buttons
		return bg.navigateButtons(msg.Type)

	case tea.KeyRunes:
		// Direct key input for numbers and operators
		char := string(msg.Runes)
		return bg.handleDirectInput(char)
	}

	return nil
}

// HandleMouse handles mouse input for button interaction
func (bg *ButtonGrid) HandleMouse(msg tea.MouseMsg) *ButtonAction {
	if msg.Type != tea.MouseLeft {
		return nil
	}

	// Find which button was clicked
	cellWidth, _ := bg.grid.CalculateDimensions(80) // Use default width for calculation
	col, row, found := bg.grid.GetCellAtPosition(msg.X, msg.Y, cellWidth)

	if found {
		buttonID := bg.generateButtonID(row, col)
		return bg.activateButton(buttonID)
	}

	return nil
}

// navigateButtons handles keyboard navigation between buttons
func (bg *ButtonGrid) navigateButtons(keyType tea.KeyType) *ButtonAction {
	if bg.focusedButton == "" {
		return nil
	}

	// Parse current button position
	var currentRow, currentCol int
	fmt.Sscanf(bg.focusedButton, "button_%d_%d", &currentRow, &currentCol)

	var newRow, newCol int

	switch keyType {
	case tea.KeyUp:
		newRow, newCol = currentRow-1, currentCol
	case tea.KeyDown:
		newRow, newCol = currentRow+1, currentCol
	case tea.KeyLeft:
		newRow, newCol = currentRow, currentCol-1
	case tea.KeyRight:
		newRow, newCol = currentRow, currentCol+1
	}

	// Check if new position is valid
	if bg.isValidPosition(newCol, newRow) {
		// Blur current button
		if currentButton, exists := bg.buttons[bg.focusedButton]; exists {
			currentButton.Blur()
		}

		// Focus new button
		bg.focusedButton = bg.generateButtonID(newRow, newCol)
		if newButton, exists := bg.buttons[bg.focusedButton]; exists {
			newButton.Focus()
		}
	}

	return nil
}

// handleDirectInput handles direct keyboard input for numbers and operators
func (bg *ButtonGrid) handleDirectInput(char string) *ButtonAction {
	// Map direct input to buttons
	inputMap := map[string]string{
		"0": "button_4_0", "1": "button_3_0", "2": "button_3_1", "3": "button_3_2",
		"4": "button_2_0", "5": "button_2_1", "6": "button_2_2", "7": "button_1_0",
		"8": "button_1_1", "9": "button_1_2", ".": "button_4_1",
		"+": "button_3_3", "-": "button_2_3", "*": "button_1_3", "/": "button_0_3",
		"=": "button_4_2",
	}

	if buttonID, exists := inputMap[char]; exists {
		return bg.activateButton(buttonID)
	}

	return nil
}

// activateButton activates a button and returns the corresponding action
func (bg *ButtonGrid) activateButton(buttonID string) *ButtonAction {
	button, exists := bg.buttons[buttonID]
	if !exists {
		return nil
	}

	// Press the button
	button.Press()

	// Create action
	action := &ButtonAction{
		Button:   button,
		Action:   "press",
		Value:    button.GetValue(),
		ButtonID: buttonID,
	}

	// Focus the button
	bg.focusedButton = buttonID

	return action
}

// isValidPosition checks if a grid position is valid
func (bg *ButtonGrid) isValidPosition(col, row int) bool {
	return col >= 0 && col < bg.dimensions.Columns &&
	       row >= 0 && row < bg.dimensions.Rows
}

// GetButton returns a button by its ID
func (bg *ButtonGrid) GetButton(buttonID string) (*components.Button, bool) {
	button, exists := bg.buttons[buttonID]
	return button, exists
}

// GetFocusedButton returns the currently focused button
func (bg *ButtonGrid) GetFocusedButton() (*components.Button, bool) {
	if bg.focusedButton == "" {
		return nil, false
	}
	return bg.GetButton(bg.focusedButton)
}

// GetButtonCount returns the total number of buttons in the grid
func (bg *ButtonGrid) GetButtonCount() int {
	return len(bg.buttons)
}

// GetButtons returns all buttons in the grid
func (bg *ButtonGrid) GetButtons() map[string]*components.Button {
	return bg.buttons
}

// SetTheme changes the theme of the button grid
func (bg *ButtonGrid) SetTheme(themeName string) error {
	err := bg.themeManager.SetTheme(themeName)
	if err != nil {
		return err
	}

	// Re-initialize layout with new theme
	bg.initializeCalculatorLayout()

	return nil
}

// GetCurrentTheme returns the current theme name
func (bg *ButtonGrid) GetCurrentTheme() string {
	return bg.themeManager.GetCurrentTheme().Name
}

// GetDimensions returns the grid dimensions
func (bg *ButtonGrid) GetDimensions() GridDimensions {
	return bg.dimensions
}

// String returns a string representation of the button grid
func (bg *ButtonGrid) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("ButtonGrid{Dimensions: %dx%d, Buttons: %d, Theme: %s, Focus: %s}",
		bg.dimensions.Columns, bg.dimensions.Rows, len(bg.buttons), bg.GetCurrentTheme(), bg.focusedButton))
	return builder.String()
}