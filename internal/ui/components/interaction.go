package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// InteractionHandler manages all user interactions with the button grid
type InteractionHandler struct {
	focusManager    *FocusManager
	keyboardHandler *KeyboardHandler
	clickHandler    *ClickHandler
	interactionLog  []InteractionEvent
	eventHandlers   map[string][]func(InteractionEvent)
}

// InteractionEvent represents a user interaction event
type InteractionEvent struct {
	Type        string
	Button      *Button
	Position    Position
	Timestamp   time.Time
	KeyPressed  string
	MouseEvent  *MouseEvent
	Context     interface{}
}

// MouseEvent represents mouse/touch input
type MouseEvent struct {
	Type      MouseEventType
	X, Y      int
	Button    MouseButton
	Modifiers KeyModifier
}

// MouseEventType defines types of mouse events
type MouseEventType int

const (
	MouseClick MouseEventType = iota
	MouseDoubleClick
	MousePress
	MouseRelease
	MouseHover
	MouseLeave
)

// MouseButton defines mouse buttons
type MouseButton int

const (
	MouseLeft MouseButton = iota
	MouseRight
	MouseMiddle
)

// KeyModifier defines key modifiers
type KeyModifier int

const (
	ModNone KeyModifier = 1 << iota
	ModShift
	ModCtrl
	ModAlt
)

// NewInteractionHandler creates a new interaction handler
func NewInteractionHandler(focusManager *FocusManager) *InteractionHandler {
	keyboardHandler := NewKeyboardHandler(focusManager)
	clickHandler := NewClickHandler(focusManager)

	return &InteractionHandler{
		focusManager:    focusManager,
		keyboardHandler: keyboardHandler,
		clickHandler:    clickHandler,
		interactionLog:  make([]InteractionEvent, 0),
		eventHandlers:   make(map[string][]func(InteractionEvent)),
	}
}

// HandleKeyEvent processes keyboard input
func (ih *InteractionHandler) HandleKeyEvent(keyEvent interface{}) (ButtonAction, bool) {
	// Convert to tea.KeyMsg if needed
	var keyMsg tea.KeyMsg
	switch event := keyEvent.(type) {
	case tea.KeyMsg:
		keyMsg = event
	default:
		return ButtonAction{}, false
	}

	// Handle the key press
	action, handled := ih.keyboardHandler.HandleKeyPress(keyMsg)
	if handled {
		// Log the interaction
		ih.logInteraction(InteractionEvent{
			Type:       "keyboard",
			Button:     action.Button,
			Position:   ih.focusManager.GetFocusPosition(),
			Timestamp:  time.Now(),
			KeyPressed: action.Value,
		})
	}

	return action, handled
}

// HandleMouseEvent processes mouse/touch input
func (ih *InteractionHandler) HandleMouseEvent(mouseEvent MouseEvent, gridLayout *GridLayout) (ButtonAction, bool) {
	// Handle the mouse event
	action, handled := ih.clickHandler.HandleMouseEvent(mouseEvent, gridLayout)
	if handled {
		// Log the interaction
		ih.logInteraction(InteractionEvent{
			Type:       "mouse",
			Button:     action.Button,
			Position:   ih.focusManager.GetFocusPosition(),
			Timestamp:  time.Now(),
			MouseEvent: &mouseEvent,
		})
	}

	return action, handled
}

// HandleDirectButtonPress handles direct button activation (e.g., by value)
func (ih *InteractionHandler) HandleDirectButtonPress(buttonValue string) (ButtonAction, bool) {
	if ih.focusManager == nil {
		return ButtonAction{}, false
	}

	buttons := ih.focusManager.GetAllButtons()
	for pos, button := range buttons {
		if button != nil && button.IsInteractive() && button.GetValue() == buttonValue {
			// Focus the button
			if err := ih.focusManager.SetFocus(pos.Row, pos.Column); err == nil {
				// Activate it
				action := button.Trigger("direct_press")

				// Handle press animation
				if err := button.Press(); err == nil {
					// In a real app, you'd use a timer for release animation
					button.Release()
				}

				// Log the interaction
				ih.logInteraction(InteractionEvent{
					Type:      "direct",
					Button:    button,
					Position:  pos,
					Timestamp: time.Now(),
				})

				return *action, true
			}
		}
	}

	return ButtonAction{}, false
}

// HandlePositionalPress handles button activation by grid position
func (ih *InteractionHandler) HandlePositionalPress(row, col int) (ButtonAction, bool) {
	if ih.focusManager == nil {
		return ButtonAction{}, false
	}

	button := ih.focusManager.GetButtonAtPosition(row, col)
	if button == nil || !button.IsInteractive() {
		return ButtonAction{}, false
	}

	// Focus the button
	if err := ih.focusManager.SetFocus(row, col); err == nil {
		// Activate it
		action := button.Trigger("positional_press")

		// Handle press animation
		if err := button.Press(); err == nil {
			button.Release()
		}

		// Log the interaction
		ih.logInteraction(InteractionEvent{
			Type:      "positional",
			Button:    button,
			Position:  Position{Row: row, Column: col},
			Timestamp: time.Now(),
		})

		return *action, true
	}

	return ButtonAction{}, false
}

// logInteraction logs an interaction event
func (ih *InteractionHandler) logInteraction(event InteractionEvent) {
	ih.interactionLog = append(ih.interactionLog, event)

	// Limit log size
	if len(ih.interactionLog) > 1000 {
		ih.interactionLog = ih.interactionLog[1:]
	}

	// Trigger event handlers
	ih.triggerEventHandlers(event)
}

// triggerEventHandlers calls registered event handlers for an event
func (ih *InteractionHandler) triggerEventHandlers(event InteractionEvent) {
	handlers := ih.eventHandlers[event.Type]
	for _, handler := range handlers {
		handler(event)
	}

	// Also trigger "all" handlers
	allHandlers := ih.eventHandlers["all"]
	for _, handler := range allHandlers {
		handler(event)
	}
}

// RegisterEventHandler registers a callback for specific interaction events
func (ih *InteractionHandler) RegisterEventHandler(eventType string, handler func(InteractionEvent)) {
	ih.eventHandlers[eventType] = append(ih.eventHandlers[eventType], handler)
}

// UnregisterEventHandler removes a registered event handler
func (ih *InteractionHandler) UnregisterEventHandler(eventType string, handler func(InteractionEvent)) {
	handlers := ih.eventHandlers[eventType]
	for i, h := range handlers {
		if &h == &handler { // Note: This is a simple pointer comparison
			ih.eventHandlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// GetInteractionLog returns the interaction history
func (ih *InteractionHandler) GetInteractionLog() []InteractionEvent {
	return ih.interactionLog
}

// ClearInteractionLog clears the interaction history
func (ih *InteractionHandler) ClearInteractionLog() {
	ih.interactionLog = make([]InteractionEvent, 0)
}

// GetKeyboardHandler returns the keyboard handler
func (ih *InteractionHandler) GetKeyboardHandler() *KeyboardHandler {
	return ih.keyboardHandler
}

// GetClickHandler returns the click handler
func (ih *InteractionHandler) GetClickHandler() *ClickHandler {
	return ih.clickHandler
}

// SetFocusManager sets or updates the focus manager
func (ih *InteractionHandler) SetFocusManager(fm *FocusManager) {
	ih.focusManager = fm
	ih.keyboardHandler.SetFocusManager(fm)
	ih.clickHandler.SetFocusManager(fm)
}

// GetLastInteraction returns the most recent interaction
func (ih *InteractionHandler) GetLastInteraction() *InteractionEvent {
	if len(ih.interactionLog) == 0 {
		return nil
	}
	return &ih.interactionLog[len(ih.interactionLog)-1]
}

// GetInteractionsByType returns interactions of a specific type
func (ih *InteractionHandler) GetInteractionsByType(eventType string) []InteractionEvent {
	var result []InteractionEvent
	for _, event := range ih.interactionLog {
		if event.Type == eventType {
			result = append(result, event)
		}
	}
	return result
}

// GetInteractionsByButton returns interactions for a specific button
func (ih *InteractionHandler) GetInteractionsByButton(button *Button) []InteractionEvent {
	var result []InteractionEvent
	for _, event := range ih.interactionLog {
		if event.Button == button {
			result = append(result, event)
		}
	}
	return result
}

// GetInteractionsSince returns interactions since a specific time
func (ih *InteractionHandler) GetInteractionsSince(since time.Time) []InteractionEvent {
	var result []InteractionEvent
	for _, event := range ih.interactionLog {
		if event.Timestamp.After(since) {
			result = append(result, event)
		}
	}
	return result
}

// ClickHandler handles mouse/touch interactions with buttons
type ClickHandler struct {
	focusManager    *FocusManager
	clickThreshold time.Duration
	dragThreshold  int
}

// NewClickHandler creates a new click handler
func NewClickHandler(focusManager *FocusManager) *ClickHandler {
	return &ClickHandler{
		focusManager:    focusManager,
		clickThreshold: 300 * time.Millisecond,
		dragThreshold:  5, // pixels
	}
}

// WithClickThreshold sets the click threshold
func (ch *ClickHandler) WithClickThreshold(threshold time.Duration) *ClickHandler {
	ch.clickThreshold = threshold
	return ch
}

// WithDragThreshold sets the drag threshold
func (ch *ClickHandler) WithDragThreshold(threshold int) *ClickHandler {
	ch.dragThreshold = threshold
	return ch
}

// HandleMouseEvent processes mouse events
func (ch *ClickHandler) HandleMouseEvent(mouseEvent MouseEvent, gridLayout *GridLayout) (ButtonAction, bool) {
	if ch.focusManager == nil {
		return ButtonAction{}, false
	}

	switch mouseEvent.Type {
	case MouseClick:
		return ch.handleClick(mouseEvent, gridLayout)
	case MouseDoubleClick:
		return ch.handleDoubleClick(mouseEvent, gridLayout)
	case MousePress:
		return ch.handlePress(mouseEvent, gridLayout)
	case MouseRelease:
		return ch.handleRelease(mouseEvent, gridLayout)
	case MouseHover:
		return ch.handleHover(mouseEvent, gridLayout)
	default:
		return ButtonAction{}, false
	}
}

// handleClick processes single clicks
func (ch *ClickHandler) handleClick(mouseEvent MouseEvent, gridLayout *GridLayout) (ButtonAction, bool) {
	// Find which button was clicked
	button, position := ch.findButtonAtPosition(mouseEvent.X, mouseEvent.Y, gridLayout)
	if button == nil || !button.IsInteractive() {
		return ButtonAction{}, false
	}

	// Focus the button
	if err := ch.focusManager.SetFocus(position.Row, position.Column); err != nil {
		return ButtonAction{}, false
	}

	// Create and return the action
	action := button.Trigger("click")

	// Handle press animation
	if err := button.Press(); err == nil {
		// In a real app, use a timer for release
		button.Release()
	}

	return *action, true
}

// handleDoubleClick processes double clicks
func (ch *ClickHandler) handleDoubleClick(mouseEvent MouseEvent, gridLayout *GridLayout) (ButtonAction, bool) {
	button, position := ch.findButtonAtPosition(mouseEvent.X, mouseEvent.Y, gridLayout)
	if button == nil || !button.IsInteractive() {
		return ButtonAction{}, false
	}

	// Focus the button
	if err := ch.focusManager.SetFocus(position.Row, position.Column); err != nil {
		return ButtonAction{}, false
	}

	// Double-click action (might be different from single click)
	action := button.Trigger("double_click")

	// Handle press animation (twice for visual effect)
	if err := button.Press(); err == nil {
		button.Release()
		button.Press()
		button.Release()
	}

	return *action, true
}

// handlePress processes mouse button press
func (ch *ClickHandler) handlePress(mouseEvent MouseEvent, gridLayout *GridLayout) (ButtonAction, bool) {
	button, position := ch.findButtonAtPosition(mouseEvent.X, mouseEvent.Y, gridLayout)
	if button == nil || !button.IsInteractive() {
		return ButtonAction{}, false
	}

	// Focus the button
	if err := ch.focusManager.SetFocus(position.Row, position.Column); err != nil {
		return ButtonAction{}, false
	}

	// Press the button
	if err := button.Press(); err == nil {
		action := button.Trigger("press")
		return *action, true
	}

	return ButtonAction{}, false
}

// handleRelease processes mouse button release
func (ch *ClickHandler) handleRelease(mouseEvent MouseEvent, gridLayout *GridLayout) (ButtonAction, bool) {
	button, position := ch.findButtonAtPosition(mouseEvent.X, mouseEvent.Y, gridLayout)
	if button == nil {
		return ButtonAction{}, false
	}
	_ = position // use position variable

	// Release the button
	if button.IsPressed() {
		if err := button.Release(); err == nil {
			action := button.Trigger("release")
			return *action, true
		}
	}

	return ButtonAction{}, false
}

// handleHover processes mouse hover
func (ch *ClickHandler) handleHover(mouseEvent MouseEvent, gridLayout *GridLayout) (ButtonAction, bool) {
	button, position := ch.findButtonAtPosition(mouseEvent.X, mouseEvent.Y, gridLayout)
	if button == nil {
		return ButtonAction{}, false
	}

	// For hover, we might want to focus without activating
	if err := ch.focusManager.SetFocus(position.Row, position.Column); err == nil {
		action := button.Trigger("hover")
		return *action, true
	}

	return ButtonAction{}, false
}

// findButtonAtPosition finds the button at the given screen coordinates
func (ch *ClickHandler) findButtonAtPosition(x, y int, gridLayout *GridLayout) (*Button, Position) {
	if gridLayout == nil {
		return nil, Position{}
	}

	// Use the grid layout to find which cell was clicked
	cellWidth, _ := gridLayout.CalculateDimensions(80) // Use default width
	col, row, found := gridLayout.GetCellAtPosition(x, y, cellWidth)

	if !found {
		return nil, Position{}
	}

	// Get the button at that position from focus manager
	button := ch.focusManager.GetButtonAtPosition(row, col)
	position := Position{Row: row, Column: col}

	return button, position
}

// SetFocusManager sets or updates the focus manager
func (ch *ClickHandler) SetFocusManager(fm *FocusManager) {
	ch.focusManager = fm
}