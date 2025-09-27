package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
)

// KeyboardHandler manages keyboard input and navigation for the button grid
type KeyboardHandler struct {
	focusManager *FocusManager
	keyBindings  keyMap
	shortcuts    map[string]key.Binding
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
	return &KeyboardHandler{
		focusManager: focusManager,
		keyBindings:  newKeyBindings(),
		shortcuts:    make(map[string]key.Binding),
	}
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

	return help
}