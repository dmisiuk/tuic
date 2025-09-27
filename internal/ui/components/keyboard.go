package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
)

// KeyboardHandler manages keyboard input and navigation for the button grid
type KeyboardHandler struct {
	focusManager     *FocusManager
	keyBindings      keyMap
	shortcuts        map[string]key.Binding
	shortcutBindings map[string]string
}

// keyMap defines all keyboard bindings for the button grid
type keyMap struct {
	// Navigation keys
	up    key.Binding
	down  key.Binding
	left  key.Binding
	right key.Binding

	// Action keys
	enter key.Binding
	space key.Binding
	tab   key.Binding
	shiftTab key.Binding

	// Special calculator keys
	escape key.Binding
	clear  key.Binding
}

// NewKeyboardHandler creates a new keyboard handler for button navigation
func NewKeyboardHandler(focusManager *FocusManager) *KeyboardHandler {
	handler := &KeyboardHandler{
		focusManager:     focusManager,
		keyBindings:      newKeyBindings(),
		shortcuts:        make(map[string]key.Binding),
		shortcutBindings:  make(map[string]string),
	}

	// Register default calculator shortcuts
	handler.RegisterCalculatorShortcuts()

	return handler
}

// newKeyBindings creates the default key bindings
func newKeyBindings() keyMap {
	return keyMap{
		// Navigation
		up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "move left"),
		),
		right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "move right"),
		),

		// Actions
		enter: key.NewBinding(
			key.WithKeys("enter", "return"),
			key.WithHelp("Enter", "activate button"),
		),
		space: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("Space", "activate button"),
		),
		tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("Tab", "next button"),
		),
		shiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("Shift+Tab", "previous button"),
		),

		// Special
		escape: key.NewBinding(
			key.WithKeys("esc", "escape"),
			key.WithHelp("Esc", "cancel/blur"),
		),
		clear: key.NewBinding(
			key.WithKeys("c", "C"),
			key.WithHelp("C", "clear input"),
		),
	}
}

// HandleKeyPress processes a keyboard input and returns the action taken
func (kh *KeyboardHandler) HandleKeyPress(msg key.Msg) (ButtonAction, bool) {
	switch {
	case kh.keyBindings.up.Matches(msg):
		return kh.handleNavigation(DirectionUp)
	case kh.keyBindings.down.Matches(msg):
		return kh.handleNavigation(DirectionDown)
	case kh.keyBindings.left.Matches(msg):
		return kh.handleNavigation(DirectionLeft)
	case kh.keyBindings.right.Matches(msg):
		return kh.handleNavigation(DirectionRight)
	case kh.keyBindings.enter.Matches(msg):
		return kh.handleActivation()
	case kh.keyBindings.space.Matches(msg):
		return kh.handleActivation()
	case kh.keyBindings.tab.Matches(msg):
		return kh.handleTabNavigation(false) // Forward
	case kh.keyBindings.shiftTab.Matches(msg):
		return kh.handleTabNavigation(true) // Backward
	case kh.keyBindings.escape.Matches(msg):
		return kh.handleEscape()
	case kh.keyBindings.clear.Matches(msg):
		return kh.handleClearKey()
	}

	// Check for direct number/operator key mappings
	return kh.handleDirectKeyMapping(msg)
}

// handleNavigation processes arrow key navigation
func (kh *KeyboardHandler) handleNavigation(direction Direction) (ButtonAction, bool) {
	if kh.focusManager == nil {
		return ButtonAction{}, false
	}

	err := kh.focusManager.MoveFocus(direction)
	if err != nil {
		// Navigation failed (no buttons or boundary reached)
		return ButtonAction{}, false
	}

	// Return a navigation action
	focusedButton := kh.focusManager.GetFocusedButton()
	if focusedButton != nil {
		return ButtonAction{
			Button: focusedButton,
			Type:   "navigate",
			Value:  fmt.Sprintf("moved_%s", direction),
		}, true
	}

	return ButtonAction{}, false
}

// handleActivation processes Enter/Space key presses
func (kh *KeyboardHandler) handleActivation() (ButtonAction, bool) {
	if kh.focusManager == nil {
		return ButtonAction{}, false
	}

	focusedButton := kh.focusManager.GetFocusedButton()
	if focusedButton == nil {
		return ButtonAction{}, false
	}

	// Trigger the button press animation and action
	action := focusedButton.Trigger("activate")

	// Handle the button press state
	if err := focusedButton.Press(); err == nil {
		// In a real implementation, you'd schedule a release after a delay
		// For now, just release immediately
		focusedButton.Release()
	}

	return *action, true
}

// handleTabNavigation processes Tab/Shift+Tab navigation
func (kh *KeyboardHandler) handleTabNavigation(reverse bool) (ButtonAction, bool) {
	if kh.focusManager == nil {
		return ButtonAction{}, false
	}

	positions := kh.focusManager.GetFocusablePositions()
	if len(positions) == 0 {
		return ButtonAction{}, false
	}

	currentPos := kh.focusManager.GetFocusPosition()
	if currentPos.Row == -1 {
		// No current focus, focus the first available
		if err := kh.focusManager.focusFirstAvailable(); err != nil {
			return ButtonAction{}, false
		}
		return kh.createNavigationAction("tab_first")
	}

	// Find next position in tab order
	var nextPos Position
	if reverse {
		nextPos = kh.findPreviousPosition(positions, currentPos)
	} else {
		nextPos = kh.findNextPosition(positions, currentPos)
	}

	if err := kh.focusManager.SetFocus(nextPos.Row, nextPos.Column); err != nil {
		return ButtonAction{}, false
	}

	return kh.createNavigationAction("tab"), true
}

// findNextPosition finds the next position in tab order (row-major)
func (kh *KeyboardHandler) findNextPosition(positions []Position, currentPos Position) Position {
	// Sort positions by row then column (row-major order)
	sorted := make([]Position, len(positions))
	copy(sorted, positions)
	// Simple bubble sort for small number of positions
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].Row > sorted[j].Row ||
			   (sorted[i].Row == sorted[j].Row && sorted[i].Column > sorted[j].Column) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Find current position in sorted list
	currentIndex := -1
	for i, pos := range sorted {
		if pos == currentPos {
			currentIndex = i
			break
		}
	}

	// Return next position (wrap around to first if at end)
	if currentIndex == -1 || currentIndex == len(sorted)-1 {
		return sorted[0]
	}
	return sorted[currentIndex+1]
}

// findPreviousPosition finds the previous position in tab order
func (kh *KeyboardHandler) findPreviousPosition(positions []Position, currentPos Position) Position {
	// Sort positions by row then column (row-major order)
	sorted := make([]Position, len(positions))
	copy(sorted, positions)
	// Simple bubble sort for small number of positions
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i].Row > sorted[j].Row ||
			   (sorted[i].Row == sorted[j].Row && sorted[i].Column > sorted[j].Column) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Find current position in sorted list
	currentIndex := -1
	for i, pos := range sorted {
		if pos == currentPos {
			currentIndex = i
			break
		}
	}

	// Return previous position (wrap around to last if at beginning)
	if currentIndex <= 0 {
		return sorted[len(sorted)-1]
	}
	return sorted[currentIndex-1]
}

// handleEscape processes Escape key (blur focus)
func (kh *KeyboardHandler) handleEscape() (ButtonAction, bool) {
	if kh.focusManager == nil {
		return ButtonAction{}, false
	}

	if !kh.focusManager.HasFocus() {
		return ButtonAction{}, false
	}

	focusedButton := kh.focusManager.GetFocusedButton()
	action := ButtonAction{
		Button: focusedButton,
		Type:   "blur",
		Value:  "escape",
	}

	if err := kh.focusManager.Blur(); err == nil {
		return action, true
	}

	return ButtonAction{}, false
}

// handleClearKey processes Clear key (C key)
func (kh *KeyboardHandler) handleClearKey() (ButtonAction, bool) {
	// Find and activate a clear button if it exists
	if kh.focusManager == nil {
		return ButtonAction{}, false
	}

	buttons := kh.focusManager.GetAllButtons()
	for pos, button := range buttons {
		if button != nil && (button.GetValue() == "C" || button.GetValue() == "CE") {
			// Focus the clear button first
			if err := kh.focusManager.SetFocus(pos.Row, pos.Column); err == nil {
				return kh.handleActivation()
			}
		}
	}

	return ButtonAction{}, false
}

// handleDirectKeyMapping processes direct key presses for numbers and operators
func (kh *KeyboardHandler) handleDirectKeyMapping(msg key.Msg) (ButtonAction, bool) {
	if kh.focusManager == nil {
		return ButtonAction{}, false
	}

	// Convert key to string
	keyStr := ""
	switch msg.Type {
	case key.KeyRunes:
		if len(msg.Runes) > 0 {
			keyStr = string(msg.Runes[0])
		}
	case key.KeyEnter:
		keyStr = "enter"
	case key.KeySpace:
		keyStr = " "
	case key.KeyBackspace:
		keyStr = "backspace"
	case key.KeyDelete:
		keyStr = "delete"
	case key.KeyUp:
		keyStr = "up"
	case key.KeyDown:
		keyStr = "down"
	case key.KeyLeft:
		keyStr = "left"
	case key.KeyRight:
		keyStr = "right"
	}

	if keyStr == "" {
		return ButtonAction{}, false
	}

	// Look for a button that matches this key
	buttons := kh.focusManager.GetAllButtons()
	for pos, button := range buttons {
		if button != nil && button.IsInteractive() {
			// Check if button value matches the key press
			if button.GetValue() == keyStr || button.GetLabel() == keyStr {
				// Focus and activate the matching button
				if err := kh.focusManager.SetFocus(pos.Row, pos.Column); err == nil {
					return kh.handleActivation()
				}
			}

			// Special mappings for calculator keys
			if kh.isKeyMatch(keyStr, button) {
				if err := kh.focusManager.SetFocus(pos.Row, pos.Column); err == nil {
					return kh.handleActivation()
				}
			}
		}
	}

	return ButtonAction{}, false
}

// isKeyMatch checks if a key press matches a button (with special mappings)
func (kh *KeyboardHandler) isKeyMatch(keyStr string, button *Button) bool {
	buttonValue := button.GetValue()

	// Number mappings
	switch keyStr {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		return buttonValue == keyStr
	}

	// Operator mappings
	switch keyStr {
	case "+", "-", "*", "/":
		return buttonValue == keyStr
	case "x", "X":
		return buttonValue == "*"
	}

	// Special function mappings
	switch keyStr {
	case "=", "enter":
		return buttonValue == "="
	case "c", "C":
		return buttonValue == "C" || buttonValue == "CE"
	case ".":
		return buttonValue == "."
	}

	return false
}

// createNavigationAction creates a navigation action for tracking
func (kh *KeyboardHandler) createNavigationAction(detail string) (ButtonAction, bool) {
	focusedButton := kh.focusManager.GetFocusedButton()
	action := ButtonAction{
		Button: focusedButton,
		Type:   "navigate",
		Value:  detail,
	}
	return action, true
}

// AddShortcut adds a custom keyboard shortcut
func (kh *KeyboardHandler) AddShortcut(keys string, button *Button) error {
	binding := key.NewBinding(key.WithKeys(keys))
	kh.shortcuts[keys] = binding
	return nil
}

// RemoveShortcut removes a custom keyboard shortcut
func (kh *KeyboardHandler) RemoveShortcut(keys string) {
	delete(kh.shortcuts, keys)
}

// GetKeyBindings returns all key bindings for help display
func (kh *KeyboardHandler) GetKeyBindings() []key.Binding {
	return []key.Binding{
		kh.keyBindings.up,
		kh.keyBindings.down,
		kh.keyBindings.left,
		kh.keyBindings.right,
		kh.keyBindings.enter,
		kh.keyBindings.space,
		kh.keyBindings.tab,
		kh.keyBindings.shiftTab,
		kh.keyBindings.escape,
		kh.keyBindings.clear,
	}
}

// GetShortcuts returns all custom shortcuts
func (kh *KeyboardHandler) GetShortcuts() map[string]key.Binding {
	return kh.shortcuts
}

// SetFocusManager sets or updates the focus manager
func (kh *KeyboardHandler) SetFocusManager(fm *FocusManager) {
	kh.focusManager = fm
}

// GetHelpText returns formatted help text for keyboard controls
func (kh *KeyboardHandler) GetHelpText() string {
	help := "Keyboard Controls:\n"
	help += fmt.Sprintf("  %s\n", kh.keyBindings.up.Help())
	help += fmt.Sprintf("  %s\n", kh.keyBindings.down.Help())
	help += fmt.Sprintf("  %s\n", kh.keyBindings.left.Help())
	help += fmt.Sprintf("  %s\n", kh.keyBindings.right.Help())
	help += fmt.Sprintf("  %s\n", kh.keyBindings.enter.Help())
	help += fmt.Sprintf("  %s\n", kh.keyBindings.space.Help())
	help += fmt.Sprintf("  %s\n", kh.keyBindings.tab.Help())
	help += fmt.Sprintf("  %s\n", kh.keyBindings.shiftTab.Help())
	help += fmt.Sprintf("  %s\n", kh.keyBindings.escape.Help())
	help += fmt.Sprintf("  %s\n", kh.keyBindings.clear.Help())

	// Add special key mappings
	help += "\nSpecial Key Mappings:\n"
	help += "  0-9, .: Direct number input\n"
	help += "  +, -, *, /: Direct operator input\n"
	help += "  x, X: Multiplication operator\n"
	help += "  =, Enter: Equals operation\n"
	help += "  c, C, Esc: Clear/Cancel\n"
	help += "  Backspace: Clear last digit\n"
	help += "  Home/End: Navigate to first/last button\n"
	help += "  PageUp/PageDown: Navigate row by row\n"

	return help
}

// HandleSpecialNavigation handles special navigation keys (Home, End, PageUp, PageDown)
func (kh *KeyboardHandler) HandleSpecialNavigation(msg key.Msg) (ButtonAction, bool) {
	if kh.focusManager == nil {
		return ButtonAction{}, false
	}

	// Check for special navigation keys
	switch {
	case kh.isHomeKey(msg):
		return kh.handleHomeKey()
	case kh.isEndKey(msg):
		return kh.handleEndKey()
	case kh.isPageUpKey(msg):
		return kh.handlePageUpKey()
	case kh.isPageDownKey(msg):
		return kh.handlePageDownKey()
	}

	return ButtonAction{}, false
}

// isHomeKey checks if the key is Home
func (kh *KeyboardHandler) isHomeKey(msg key.Msg) bool {
	return msg.Type == key.KeyHome
}

// isEndKey checks if the key is End
func (kh *KeyboardHandler) isEndKey(msg key.Msg) bool {
	return msg.Type == key.KeyEnd
}

// isPageUpKey checks if the key is PageUp
func (kh *KeyboardHandler) isPageUpKey(msg key.Msg) bool {
	return msg.Type == key.KeyPageUp
}

// isPageDownKey checks if the key is PageDown
func (kh *KeyboardHandler) isPageDownKey(msg key.Msg) bool {
	return msg.Type == key.KeyPageDown
}

// handleHomeKey navigates to the first button
func (kh *KeyboardHandler) handleHomeKey() (ButtonAction, bool) {
	if err := kh.focusManager.focusFirstAvailable(); err != nil {
		return ButtonAction{}, false
	}
	return kh.createNavigationAction("home")
}

// handleEndKey navigates to the last button
func (kh *KeyboardHandler) handleEndKey() (ButtonAction, bool) {
	// Find bottom-right button
	buttons := kh.focusManager.GetAllButtons()
	if len(buttons) == 0 {
		return ButtonAction{}, false
	}

	// Find bottommost row
	maxRow := -1
	for pos := range buttons {
		if pos.Row > maxRow {
			maxRow = pos.Row
		}
	}

	// Find rightmost button in that row
	maxCol := -1
	targetPos := Position{}
	for pos, button := range buttons {
		if pos.Row == maxRow && button.IsInteractive() && pos.Column > maxCol {
			maxCol = pos.Column
			targetPos = pos
		}
	}

	if maxCol == -1 {
		return ButtonAction{}, false
	}

	if err := kh.focusManager.SetFocus(targetPos.Row, targetPos.Column); err != nil {
		return ButtonAction{}, false
	}

	return kh.createNavigationAction("end")
}

// handlePageUpKey navigates up by entire rows
func (kh *KeyboardHandler) handlePageUpKey() (ButtonAction, bool) {
	currentPos := kh.focusManager.GetFocusPosition()
	if currentPos.Row == -1 {
		return kh.handleHomeKey()
	}

	// Move up 3 rows or to first row
	targetRow := currentPos.Row - 3
	if targetRow < 0 {
		targetRow = 0
	}

	// Try to stay in same column
	for row := targetRow; row <= currentPos.Row; row++ {
		button := kh.focusManager.GetButtonAtPosition(row, currentPos.Column)
		if button != nil && button.IsInteractive() {
			if err := kh.focusManager.SetFocus(row, currentPos.Column); err == nil {
				return kh.createNavigationAction("page_up")
			}
		}
	}

	// If no button in same column, find first available in target rows
	for row := targetRow; row <= currentPos.Row; row++ {
		for col := 0; col < 10; col++ {
			button := kh.focusManager.GetButtonAtPosition(row, col)
			if button != nil && button.IsInteractive() {
				if err := kh.focusManager.SetFocus(row, col); err == nil {
					return kh.createNavigationAction("page_up")
				}
			}
		}
	}

	return ButtonAction{}, false
}

// handlePageDownKey navigates down by entire rows
func (kh *KeyboardHandler) handlePageDownKey() (ButtonAction, bool) {
	currentPos := kh.focusManager.GetFocusPosition()
	if currentPos.Row == -1 {
		return kh.handleHomeKey()
	}

	// Move down 3 rows or to last row
	targetRow := currentPos.Row + 3

	// Find max row
	maxRow := -1
	buttons := kh.focusManager.GetAllButtons()
	for pos := range buttons {
		if pos.Row > maxRow {
			maxRow = pos.Row
		}
	}

	if targetRow > maxRow {
		targetRow = maxRow
	}

	// Try to stay in same column
	for row := currentPos.Row; row <= targetRow; row++ {
		button := kh.focusManager.GetButtonAtPosition(row, currentPos.Column)
		if button != nil && button.IsInteractive() {
			if err := kh.focusManager.SetFocus(row, currentPos.Column); err == nil {
				return kh.createNavigationAction("page_down")
			}
		}
	}

	// If no button in same column, find first available in target rows
	for row := currentPos.Row; row <= targetRow; row++ {
		for col := 0; col < 10; col++ {
			button := kh.focusManager.GetButtonAtPosition(row, col)
			if button != nil && button.IsInteractive() {
				if err := kh.focusManager.SetFocus(row, col); err == nil {
					return kh.createNavigationAction("page_down")
				}
			}
		}
	}

	return ButtonAction{}, false
}

// EnhancedHandleKeyPress extends HandleKeyPress to include special navigation
func (kh *KeyboardHandler) EnhancedHandleKeyPress(msg key.Msg) (ButtonAction, bool) {
	// First try special navigation
	if action, handled := kh.HandleSpecialNavigation(msg); handled {
		return action, true
	}

	// Then try regular key handling
	return kh.HandleKeyPress(msg)
}

// RegisterCalculatorShortcuts registers common calculator keyboard shortcuts
func (kh *KeyboardHandler) RegisterCalculatorShortcuts() {
	// Common calculator shortcuts
	shortcuts := map[string]string{
		"+":        "+",
		"-":        "-",
		"*":        "*",
		"/":        "/",
		"x":        "*",
		"X":        "*",
		"=":        "=",
		"enter":    "=",
		"return":   "=",
		"escape":   "C",
		"esc":      "C",
		"c":        "C",
		"C":        "C",
		".":        ".",
		"backspace": "backspace",
		"delete":   "CE",
	}

	for keys, value := range shortcuts {
		kh.shortcutBindings[keys] = value
	}
}

// GetShortcutBindings returns all registered shortcut bindings
func (kh *KeyboardHandler) GetShortcutBindings() map[string]string {
	bindings := make(map[string]string)
	for k, v := range kh.shortcutBindings {
		bindings[k] = v
	}
	return bindings
}

// HandleBackspace handles backspace key for clearing input
func (kh *KeyboardHandler) HandleBackspace() (ButtonAction, bool) {
	// Look for a backspace or CE button
	buttons := kh.focusManager.GetAllButtons()
	for pos, button := range buttons {
		if button != nil && button.IsInteractive() {
			if button.GetValue() == "backspace" || button.GetValue() == "CE" {
				if err := kh.focusManager.SetFocus(pos.Row, pos.Column); err == nil {
					return kh.handleActivation()
				}
			}
		}
	}

	// If no dedicated backspace button, try to trigger a clear action
	return kh.handleClearKey()
}

// isKeyMatch checks if a key press matches a button (with enhanced mappings)
func (kh *KeyboardHandler) isKeyMatch(keyStr string, button *Button) bool {
	buttonValue := button.GetValue()

	// Check shortcut bindings first
	if shortcutValue, exists := shortcutBindings[keyStr]; exists {
		return buttonValue == shortcutValue
	}

	// Number mappings
	switch keyStr {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		return buttonValue == keyStr
	}

	// Operator mappings
	switch keyStr {
	case "+", "-", "*", "/":
		return buttonValue == keyStr
	case "x", "X":
		return buttonValue == "*"
	}

	// Special function mappings
	switch keyStr {
	case "=", "enter", "return":
		return buttonValue == "="
	case "c", "C", "escape", "esc":
		return buttonValue == "C" || buttonValue == "CE"
	case ".":
		return buttonValue == "."
	case "backspace":
		return buttonValue == "backspace" || buttonValue == "CE"
	case "delete":
		return buttonValue == "CE"
	}

	return false
}

// GetQuickReference returns a quick reference card for keyboard shortcuts
func (kh *KeyboardHandler) GetQuickReference() string {
	ref := "Quick Reference:\n"
	ref += "┌─────────────────┬─────────────────────────┐\n"
	ref += "│ Action          │ Keys                    │\n"
	ref += "├─────────────────┼─────────────────────────┤\n"
	ref += "│ Navigate        │ Arrow keys, HJKL        │\n"
	ref += "│ Activate        │ Enter, Space, 0-9, ops │\n"
	ref += "│ Next/Prev       │ Tab/Shift+Tab          │\n"
	ref += "│ First/Last      │ Home/End               │\n"
	ref += "│ Page Nav        │ PageUp/PageDown        │\n"
	ref += "│ Clear           │ C, Esc, Backspace      │\n"
	ref += "│ Equals          │ =, Enter               │\n"
	ref += "│ Multiply        │ *, x, X                │\n"
	ref += "└─────────────────┴─────────────────────────┘\n"
	return ref
}