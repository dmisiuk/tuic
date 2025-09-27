package components

import (
	"errors"
	"fmt"
)

var (
	ErrNoFocusableButtons = errors.New("no focusable buttons available")
	ErrInvalidFocusMove   = errors.New("invalid focus movement")
)

// FocusManager manages focus state and navigation for a collection of buttons
type FocusManager struct {
	buttons        map[Position]*Button
	focusedButton  *Button
	focusPosition  Position
	wrapping       bool
	cycleMode      FocusCycleMode
	focusedHistory []Position
}

// FocusCycleMode defines how focus wraps around the grid boundaries
type FocusCycleMode int

const (
	// CycleNone prevents wrapping (focus stops at boundaries)
	CycleNone FocusCycleMode = iota

	// CycleRow allows wrapping within the same row
	CycleRow

	// CycleColumn allows wrapping within the same column
	CycleColumn

	// CycleBoth allows full grid wrapping
	CycleBoth
)

// NewFocusManager creates a new focus manager for button navigation
func NewFocusManager() *FocusManager {
	return &FocusManager{
		buttons:        make(map[Position]*Button),
		wrapping:       true,
		cycleMode:      CycleBoth,
		focusedHistory: make([]Position, 0),
	}
}

// WithWrapping enables or disables focus wrapping
func (fm *FocusManager) WithWrapping(enabled bool) *FocusManager {
	fm.wrapping = enabled
	return fm
}

// WithCycleMode sets the focus cycling behavior
func (fm *FocusManager) WithCycleMode(mode FocusCycleMode) *FocusManager {
	fm.cycleMode = mode
	return fm
}

// AddButton adds a button to the focus management system
func (fm *FocusManager) AddButton(button *Button) error {
	if button == nil {
		return fmt.Errorf("cannot add nil button to focus manager")
	}

	position := button.GetPosition()

	// Remove focus from existing button at this position
	if existing, exists := fm.buttons[position]; exists && existing.IsFocused() {
		existing.Blur()
	}

	fm.buttons[position] = button

	// If this is the first button added and no focus is set, focus it
	if fm.focusedButton == nil && len(fm.buttons) == 1 {
		return fm.SetFocus(position.Row, position.Column)
	}

	return nil
}

// RemoveButton removes a button from focus management
func (fm *FocusManager) RemoveButton(position Position) error {
	button, exists := fm.buttons[position]
	if !exists {
		return fmt.Errorf("no button found at position %v", position)
	}

	// If this button is focused, move focus
	if button.IsFocused() {
		if err := fm.Blur(); err != nil {
			return err
		}
	}

	delete(fm.buttons, position)
	return nil
}

// SetFocus sets focus to the button at the specified position
func (fm *FocusManager) SetFocus(row, col int) error {
	position := Position{Row: row, Column: col}

	button, exists := fm.buttons[position]
	if !exists {
		return fmt.Errorf("no button found at position %v", position)
	}

	// Blur currently focused button
	if fm.focusedButton != nil {
		if err := fm.focusedButton.Blur(); err != nil {
			return fmt.Errorf("failed to blur current focused button: %w", err)
		}
	}

	// Focus new button
	if err := button.Focus(); err != nil {
		return fmt.Errorf("failed to focus button at %v: %w", position, err)
	}

	// Update focus state
	fm.focusedButton = button
	fm.focusPosition = position

	// Add to history (avoid duplicates)
	fm.addToHistory(position)

	return nil
}

// Blur removes focus from the currently focused button
func (fm *FocusManager) Blur() error {
	if fm.focusedButton == nil {
		return nil // No button is focused
	}

	if err := fm.focusedButton.Blur(); err != nil {
		return fmt.Errorf("failed to blur focused button: %w", err)
	}

	fm.focusedButton = nil
	fm.focusPosition = Position{-1, -1} // Invalid position indicates no focus

	return nil
}

// MoveFocus moves focus in the specified direction
func (fm *FocusManager) MoveFocus(direction Direction) error {
	if fm.focusedButton == nil {
		// If no button is focused, focus the first available button
		return fm.focusFirstAvailable()
	}

	currentPos := fm.focusPosition
	newPos, err := fm.findNextPosition(currentPos, direction)
	if err != nil {
		return err
	}

	return fm.SetFocus(newPos.Row, newPos.Column)
}

// findNextPosition finds the next valid position in the specified direction
func (fm *FocusManager) findNextPosition(currentPos Position, direction Direction) (Position, error) {
	var newPos Position
	var valid bool

	switch direction {
	case DirectionUp:
		newPos = Position{Row: currentPos.Row - 1, Column: currentPos.Column}
	case DirectionDown:
		newPos = Position{Row: currentPos.Row + 1, Column: currentPos.Column}
	case DirectionLeft:
		newPos = Position{Row: currentPos.Row, Column: currentPos.Column - 1}
	case DirectionRight:
		newPos = Position{Row: currentPos.Row, Column: currentPos.Column + 1}
	default:
		return Position{}, ErrInvalidFocusMove
	}

	// Check if new position has a button
	if _, exists := fm.buttons[newPos]; exists {
		return newPos, nil
	}

	// Handle wrapping based on cycle mode
	if fm.wrapping {
		return fm.handleWrapping(currentPos, direction)
	}

	// Find nearest available button in the direction
	return fm.findNearestAvailable(currentPos, direction)
}

// handleWrapping handles focus wrapping based on cycle mode
func (fm *FocusManager) handleWrapping(currentPos Position, direction Direction) (Position, error) {
	switch fm.cycleMode {
	case CycleRow:
		return fm.wrapRow(currentPos, direction)
	case CycleColumn:
		return fm.wrapColumn(currentPos, direction)
	case CycleBoth:
		return fm.wrapGrid(currentPos, direction)
	case CycleNone:
		// Find nearest available button without wrapping
		return fm.findNearestAvailable(currentPos, direction)
	default:
		return Position{}, ErrInvalidFocusMove
	}
}

// wrapRow wraps focus within the same row
func (fm *FocusManager) wrapRow(currentPos Position, direction Direction) (Position, error) {
	if direction == DirectionLeft || direction == DirectionRight {
		// Find the leftmost or rightmost button in the same row
		var buttonsInRow []Position
		for pos := range fm.buttons {
			if pos.Row == currentPos.Row {
				buttonsInRow = append(buttonsInRow, pos)
			}
		}

		if len(buttonsInRow) == 0 {
			return Position{}, ErrNoFocusableButtons
		}

		if direction == DirectionLeft {
			// Find rightmost button
			rightmost := buttonsInRow[0]
			for _, pos := range buttonsInRow {
				if pos.Column > rightmost.Column {
					rightmost = pos
				}
			}
			return rightmost, nil
		} else {
			// Find leftmost button
			leftmost := buttonsInRow[0]
			for _, pos := range buttonsInRow {
				if pos.Column < leftmost.Column {
					leftmost = pos
				}
			}
			return leftmost, nil
		}
	}

	// For vertical movement, don't wrap row-wise
	return fm.findNearestAvailable(currentPos, direction)
}

// wrapColumn wraps focus within the same column
func (fm *FocusManager) wrapColumn(currentPos Position, direction Direction) (Position, error) {
	if direction == DirectionUp || direction == DirectionDown {
		// Find the topmost or bottommost button in the same column
		var buttonsInColumn []Position
		for pos := range fm.buttons {
			if pos.Column == currentPos.Column {
				buttonsInColumn = append(buttonsInColumn, pos)
			}
		}

		if len(buttonsInColumn) == 0 {
			return Position{}, ErrNoFocusableButtons
		}

		if direction == DirectionUp {
			// Find bottommost button
			bottommost := buttonsInColumn[0]
			for _, pos := range buttonsInColumn {
				if pos.Row > bottommost.Row {
					bottommost = pos
				}
			}
			return bottommost, nil
		} else {
			// Find topmost button
			topmost := buttonsInColumn[0]
			for _, pos := range buttonsInColumn {
				if pos.Row < topmost.Row {
					topmost = pos
				}
			}
			return topmost, nil
		}
	}

	// For horizontal movement, don't wrap column-wise
	return fm.findNearestAvailable(currentPos, direction)
}

// wrapGrid wraps focus around the entire grid
func (fm *FocusManager) wrapGrid(currentPos Position, direction Direction) (Position, error) {
	if len(fm.buttons) == 0 {
		return Position{}, ErrNoFocusableButtons
	}

	switch direction {
	case DirectionUp:
		// Find bottommost button in the same column or nearby
		return fm.findBottommostInColumn(currentPos.Column)
	case DirectionDown:
		// Find topmost button in the same column or nearby
		return fm.findTopmostInColumn(currentPos.Column)
	case DirectionLeft:
		// Find rightmost button in the same row or nearby
		return fm.findRightmostInRow(currentPos.Row)
	case DirectionRight:
		// Find leftmost button in the same row or nearby
		return fm.findLeftmostInRow(currentPos.Row)
	default:
		return Position{}, ErrInvalidFocusMove
	}
}

// findNearestAvailable finds the nearest available button in the specified direction
func (fm *FocusManager) findNearestAvailable(currentPos Position, direction Direction) (Position, error) {
	// Search in expanding spirals from the current position
	for distance := 1; distance <= 10; distance++ {
		candidates := fm.getPositionsAtDistance(currentPos, distance, direction)
		for _, pos := range candidates {
			if _, exists := fm.buttons[pos]; exists {
				return pos, nil
			}
		}
	}

	return Position{}, ErrNoFocusableButtons
}

// getPositionsAtDistance returns positions at a specific distance in the given direction
func (fm *FocusManager) getPositionsAtDistance(currentPos Position, distance int, direction Direction) []Position {
	var positions []Position

	switch direction {
	case DirectionUp:
		for col := 0; col < 10; col++ { // Reasonable grid size limit
			pos := Position{Row: currentPos.Row - distance, Column: col}
			if pos.Row >= 0 {
				positions = append(positions, pos)
			}
		}
	case DirectionDown:
		for col := 0; col < 10; col++ {
			pos := Position{Row: currentPos.Row + distance, Column: col}
			positions = append(positions, pos)
		}
	case DirectionLeft:
		for row := 0; row < 10; row++ {
			pos := Position{Row: row, Column: currentPos.Column - distance}
			if pos.Column >= 0 {
				positions = append(positions, pos)
			}
		}
	case DirectionRight:
		for row := 0; row < 10; row++ {
			pos := Position{Row: row, Column: currentPos.Column + distance}
			positions = append(positions, pos)
		}
	}

	return positions
}

// Helper methods for grid wrapping
func (fm *FocusManager) findBottommostInColumn(column int) (Position, error) {
	var candidates []Position
	for pos := range fm.buttons {
		if pos.Column == column {
			candidates = append(candidates, pos)
		}
	}

	if len(candidates) == 0 {
		// If no buttons in this column, find bottommost overall
		return fm.findBottommostOverall()
	}

	bottommost := candidates[0]
	for _, pos := range candidates {
		if pos.Row > bottommost.Row {
			bottommost = pos
		}
	}
	return bottommost, nil
}

func (fm *FocusManager) findTopmostInColumn(column int) (Position, error) {
	var candidates []Position
	for pos := range fm.buttons {
		if pos.Column == column {
			candidates = append(candidates, pos)
		}
	}

	if len(candidates) == 0 {
		// If no buttons in this column, find topmost overall
		return fm.findTopmostOverall()
	}

	topmost := candidates[0]
	for _, pos := range candidates {
		if pos.Row < topmost.Row {
			topmost = pos
		}
	}
	return topmost, nil
}

func (fm *FocusManager) findRightmostInRow(row int) (Position, error) {
	var candidates []Position
	for pos := range fm.buttons {
		if pos.Row == row {
			candidates = append(candidates, pos)
		}
	}

	if len(candidates) == 0 {
		// If no buttons in this row, find rightmost overall
		return fm.findRightmostOverall()
	}

	rightmost := candidates[0]
	for _, pos := range candidates {
		if pos.Column > rightmost.Column {
			rightmost = pos
		}
	}
	return rightmost, nil
}

func (fm *FocusManager) findLeftmostInRow(row int) (Position, error) {
	var candidates []Position
	for pos := range fm.buttons {
		if pos.Row == row {
			candidates = append(candidates, pos)
		}
	}

	if len(candidates) == 0 {
		// If no buttons in this row, find leftmost overall
		return fm.findLeftmostOverall()
	}

	leftmost := candidates[0]
	for _, pos := range candidates {
		if pos.Column < leftmost.Column {
			leftmost = pos
		}
	}
	return leftmost, nil
}

func (fm *FocusManager) findBottommostOverall() (Position, error) {
	if len(fm.buttons) == 0 {
		return Position{}, ErrNoFocusableButtons
	}

	var bottommost Position
	for pos := range fm.buttons {
		if pos.Row > bottommost.Row || bottommost.Row == -1 {
			bottommost = pos
		}
	}
	return bottommost, nil
}

func (fm *FocusManager) findTopmostOverall() (Position, error) {
	if len(fm.buttons) == 0 {
		return Position{}, ErrNoFocusableButtons
	}

	topmost := Position{Row: 999, Column: 999} // Start with high values
	for pos := range fm.buttons {
		if pos.Row < topmost.Row {
			topmost = pos
		}
	}
	return topmost, nil
}

func (fm *FocusManager) findRightmostOverall() (Position, error) {
	if len(fm.buttons) == 0 {
		return Position{}, ErrNoFocusableButtons
	}

	rightmost := Position{Row: -1, Column: -1}
	for pos := range fm.buttons {
		if pos.Column > rightmost.Column || rightmost.Column == -1 {
			rightmost = pos
		}
	}
	return rightmost, nil
}

func (fm *FocusManager) findLeftmostOverall() (Position, error) {
	if len(fm.buttons) == 0 {
		return Position{}, ErrNoFocusableButtons
	}

	leftmost := Position{Row: 999, Column: 999}
	for pos := range fm.buttons {
		if pos.Column < leftmost.Column {
			leftmost = pos
		}
	}
	return leftmost, nil
}

// focusFirstAvailable focuses the first available button in the grid
func (fm *FocusManager) focusFirstAvailable() error {
	if len(fm.buttons) == 0 {
		return ErrNoFocusableButtons
	}

	// Try to find top-left button first
	topLeft, err := fm.findTopmostOverall()
	if err != nil {
		return err
	}

	// Find leftmost in that row
	leftmost, err := fm.findLeftmostInRow(topLeft.Row)
	if err != nil {
		return err
	}

	return fm.SetFocus(leftmost.Row, leftmost.Column)
}

// addToHistory adds a position to focus history, avoiding duplicates
func (fm *FocusManager) addToHistory(position Position) {
	// Remove if already exists to avoid duplicates
	for i, histPos := range fm.focusedHistory {
		if histPos == position {
			fm.focusedHistory = append(fm.focusedHistory[:i], fm.focusedHistory[i+1:]...)
			break
		}
	}

	fm.focusedHistory = append(fm.focusedHistory, position)

	// Limit history size
	if len(fm.focusedHistory) > 50 {
		fm.focusedHistory = fm.focusedHistory[1:]
	}
}

// GetFocusedButton returns the currently focused button
func (fm *FocusManager) GetFocusedButton() *Button {
	return fm.focusedButton
}

// GetFocusPosition returns the current focus position
func (fm *FocusManager) GetFocusPosition() Position {
	return fm.focusPosition
}

// HasFocus returns true if any button has focus
func (fm *FocusManager) HasFocus() bool {
	return fm.focusedButton != nil
}

// GetButtonAtPosition returns the button at the specified position
func (fm *FocusManager) GetButtonAtPosition(row, col int) *Button {
	pos := Position{Row: row, Column: col}
	return fm.buttons[pos]
}

// GetAllButtons returns all buttons managed by this focus manager
func (fm *FocusManager) GetAllButtons() map[Position]*Button {
	return fm.buttons
}

// GetFocusablePositions returns all positions that have focusable buttons
func (fm *FocusManager) GetFocusablePositions() []Position {
	positions := make([]Position, 0, len(fm.buttons))
	for pos := range fm.buttons {
		if button := fm.buttons[pos]; button != nil && button.IsInteractive() {
			positions = append(positions, pos)
		}
	}
	return positions
}

// GetFocusHistory returns the focus navigation history
func (fm *FocusManager) GetFocusHistory() []Position {
	return fm.focusedHistory
}

// ClearHistory clears the focus navigation history
func (fm *FocusManager) ClearHistory() {
	fm.focusedHistory = make([]Position, 0)
}

// Clear removes all buttons from focus management
func (fm *FocusManager) Clear() error {
	if err := fm.Blur(); err != nil {
		return err
	}

	fm.buttons = make(map[Position]*Button)
	fm.ClearHistory()
	return nil
}