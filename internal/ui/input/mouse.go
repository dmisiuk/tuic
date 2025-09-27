package input

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// MouseHandler handles all mouse interactions for the calculator TUI
type MouseHandler struct {
	state *MouseState

	// Configuration
	enabled bool
	clickTolerance int
	doubleClickDelay time.Duration

	// Button action mappings
	buttonActions map[string]ButtonAction
}

// NewMouseHandler creates a new mouse handler instance
func NewMouseHandler() *MouseHandler {
	return &MouseHandler{
		state:           NewMouseState(),
		enabled:         true,
		clickTolerance:  5, // pixels tolerance for click detection
		doubleClickDelay: 500 * time.Millisecond,
		buttonActions:   make(map[string]ButtonAction),
	}
}

// SetEnabled enables or disables mouse handling
func (mh *MouseHandler) SetEnabled(enabled bool) {
	mh.enabled = enabled
	if !enabled {
		mh.state.ResetPress()
		mh.state.ResetScroll()
	}
}

// RegisterButtonAction registers an action for a specific button
func (mh *MouseHandler) RegisterButtonAction(buttonID string, action ButtonAction) {
	mh.buttonActions[buttonID] = action
}

// UnregisterButtonAction removes a button action
func (mh *MouseHandler) UnregisterButtonAction(buttonID string) {
	delete(mh.buttonActions, buttonID)
	mh.state.UnregisterButton(buttonID)
}

// RegisterButton registers a button with its boundaries and action
func (mh *MouseHandler) RegisterButton(buttonID string, x, y, width, height int, action ButtonAction) {
	mh.state.RegisterButton(buttonID, x, y, width, height)
	mh.buttonActions[buttonID] = action
}

// HandleMessage processes tea mouse messages and returns appropriate events
func (mh *MouseHandler) HandleMessage(msg tea.Msg) []tea.Msg {
	if !mh.enabled {
		return nil
	}

	mouseMsg, ok := msg.(tea.MouseMsg)
	if !ok {
		return nil
	}

	var events []tea.Msg

	switch mouseMsg.Type {
	case tea.MouseMotion:
		events = mh.handleMouseMove(mouseMsg)
	case tea.MousePress:
		events = mh.handleMousePress(mouseMsg)
	case tea.MouseRelease:
		events = mh.handleMouseRelease(mouseMsg)
	case tea.MouseWheelUp, tea.MouseWheelDown:
		events = mh.handleMouseWheel(mouseMsg)
	}

	return events
}

// handleMouseMove processes mouse movement events
func (mh *MouseHandler) handleMouseMove(msg tea.MouseMsg) []tea.Msg {
	hoveredButton := mh.state.UpdateHover(msg.X, msg.Y)

	var events []tea.Msg

	// Send hover enter event
	if hoveredButton != "" && !mh.state.IsHovering {
		events = append(events, MouseEvent{
			Type:     MouseEventMove,
			X:        msg.X,
			Y:        msg.Y,
			ButtonID: hoveredButton,
			Action:   mh.buttonActions[hoveredButton],
		})
	}

	return events
}

// handleMousePress processes mouse press events
func (mh *MouseHandler) handleMousePress(msg tea.MouseMsg) []tea.Msg {
	mh.state.StartPress(msg.X, msg.Y, msg.Button)

	pressedButton := mh.state.GetButtonAtPosition(msg.X, msg.Y)
	if pressedButton != "" {
		return []tea.Msg{
			MouseEvent{
				Type:     MouseEventPress,
				Button:   msg.Button,
				X:        msg.X,
				Y:        msg.Y,
				ButtonID: pressedButton,
				Action:   mh.buttonActions[pressedButton],
			},
		}
	}

	return nil
}

// handleMouseRelease processes mouse release events
func (mh *MouseHandler) handleMouseRelease(msg tea.MouseMsg) []tea.Msg {
	releasedButton := mh.state.GetButtonAtPosition(msg.X, msg.Y)
	clickedButton := mh.state.EndPress(msg.X, msg.Y)

	var events []tea.Msg

	// Send release event
	if releasedButton != "" {
		events = append(events, MouseEvent{
			Type:     MouseEventRelease,
			X:        msg.X,
			Y:        msg.Y,
			ButtonID: releasedButton,
			Action:   mh.buttonActions[releasedButton],
		})
	}

	// Send click event if valid click detected
	if clickedButton != "" {
		clickType := MouseEventClick
		if mh.state.ClickCount > 1 {
			clickType = MouseEventDoubleClick
		}

		events = append(events, MouseEvent{
			Type:     clickType,
			Button:   mh.state.Button,
			X:        msg.X,
			Y:        msg.Y,
			ButtonID: clickedButton,
			Action:   mh.buttonActions[clickedButton],
		})

		// Execute button action if handler is defined
		if action, exists := mh.buttonActions[clickedButton]; exists && action.Handler != nil {
			events = append(events, action.Handler())
		}
	}

	return events
}

// handleMouseWheel processes mouse wheel events
func (mh *MouseHandler) handleMouseWheel(msg tea.MouseMsg) []tea.Msg {
	var delta int
	if msg.Type == tea.MouseWheelUp {
		delta = 1
	} else {
		delta = -1
	}

	mh.state.UpdateScroll(msg.X, msg.Y, delta)

	return []tea.Msg{
		MouseEvent{
			Type:   MouseEventScroll,
			X:      msg.X,
			Y:      msg.Y,
			Scroll: delta,
		},
	}
}

// GetHoveredButton returns the currently hovered button ID
func (mh *MouseHandler) GetHoveredButton() string {
	return mh.state.HoveredButton
}

// IsHovering returns true if mouse is hovering over a button
func (mh *MouseHandler) IsHovering() bool {
	return mh.state.IsHovering
}

// GetPressedButton returns the currently pressed button ID
func (mh *MouseHandler) GetPressedButton() string {
	return mh.state.PressedButton
}

// GetMousePosition returns the current mouse position
func (mh *MouseHandler) GetMousePosition() (int, int) {
	return mh.state.X, mh.state.Y
}

// GetScrollDelta returns the current scroll delta
func (mh *MouseHandler) GetScrollDelta() int {
	return mh.state.ScrollDelta
}

// Reset resets the mouse handler state
func (mh *MouseHandler) Reset() {
	mh.state.ResetPress()
	mh.state.ResetScroll()
	mh.state.HoveredButton = ""
	mh.state.IsHovering = false
	mh.state.X = -1
	mh.state.Y = -1
}

// ClearButtons clears all registered buttons and actions
func (mh *MouseHandler) ClearButtons() {
	mh.state.ButtonBounds = make(map[string]ButtonRect)
	mh.buttonActions = make(map[string]ButtonAction)
}