package input

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ClickDetector handles advanced click detection and button press/release logic
type ClickDetector struct {
	state *MouseState

	// Click detection settings
	clickTolerance      int
	doubleClickDelay    time.Duration
	longClickDelay      time.Duration
	clickThreshold      int

	// Click state tracking
	clickCount         int
	lastClickTime      int64
	lastClickPosition  ClickPosition
	longClickActive    bool
	longClickStartTime int64

	// Press tracking
	pressHistory       []PressEvent
	maxPressHistory    int

	// Button mapping
	buttonPressActions map[string]PressAction
	buttonReleaseActions map[string]ReleaseAction
}

// ClickPosition represents a position for click detection
type ClickPosition struct {
	X int
	Y int
}

// PressEvent represents a button press event
type PressEvent struct {
	ButtonID string
	X        int
	Y        int
	Button   tea.MouseButton
	Timestamp int64
}

// ReleaseEvent represents a button release event
type ReleaseEvent struct {
	ButtonID string
	X        int
	Y        int
	Button   tea.MouseButton
	WasClick bool
	WasLongClick bool
	Timestamp int64
}

// PressAction represents an action triggered by button press
type PressAction struct {
	Type        string
	Handler     func(event PressEvent) tea.Msg
	VisualFeedback bool
}

// ReleaseAction represents an action triggered by button release
type ReleaseAction struct {
	Type        string
	Handler     func(event ReleaseEvent) tea.Msg
	VisualFeedback bool
}

// ClickEvent represents a detected click event
type ClickEvent struct {
	Type        ClickType
	ButtonID    string
	X           int
	Y           int
	Button      tea.MouseButton
	ClickCount  int
	Timestamp   int64
}

// ClickType represents the type of click detected
type ClickType int

const (
	ClickSingle ClickType = iota
	ClickDouble
	ClickTriple
	ClickLong
)

// NewClickDetector creates a new click detector instance
func NewClickDetector() *ClickDetector {
	return &ClickDetector{
		state:                NewMouseState(),
		clickTolerance:       5,
		doubleClickDelay:     500 * time.Millisecond,
		longClickDelay:       1000 * time.Millisecond,
		clickThreshold:       1,
		clickCount:          0,
		lastClickTime:       0,
		lastClickPosition:   ClickPosition{X: -1, Y: -1},
		longClickActive:     false,
		longClickStartTime:  0,
		pressHistory:        make([]PressEvent, 0),
		maxPressHistory:     10,
		buttonPressActions:   make(map[string]PressAction),
		buttonReleaseActions: make(map[string]ReleaseAction),
	}
}

// SetClickTolerance sets the pixel tolerance for click detection
func (cd *ClickDetector) SetClickTolerance(tolerance int) {
	cd.clickTolerance = tolerance
}

// SetDoubleClickDelay sets the time threshold for double-click detection
func (cd *ClickDetector) SetDoubleClickDelay(delay time.Duration) {
	cd.doubleClickDelay = delay
}

// SetLongClickDelay sets the time threshold for long-click detection
func (cd *ClickDetector) SetLongClickDelay(delay time.Duration) {
	cd.longClickDelay = delay
}

// SetClickThreshold sets the minimum movement before considering a drag
func (cd *ClickDetector) SetClickThreshold(threshold int) {
	cd.clickThreshold = threshold
}

// RegisterPressAction registers an action for button press
func (cd *ClickDetector) RegisterPressAction(buttonID string, action PressAction) {
	cd.buttonPressActions[buttonID] = action
}

// RegisterReleaseAction registers an action for button release
func (cd *ClickDetector) RegisterReleaseAction(buttonID string, action ReleaseAction) {
	cd.buttonReleaseActions[buttonID] = action
}

// HandleButtonPress processes a button press event
func (cd *ClickDetector) HandleButtonPress(x, y int, button tea.MouseButton, timestamp int64) []tea.Msg {
	buttonID := cd.state.GetButtonAtPosition(x, y)

	// Update state
	cd.state.StartPress(x, y, button)

	// Record press event
	pressEvent := PressEvent{
		ButtonID:  buttonID,
		X:         x,
		Y:         y,
		Button:    button,
		Timestamp: timestamp,
	}

	cd.recordPressEvent(pressEvent)

	// Start long click timer if on a button
	if buttonID != "" {
		cd.longClickStartTime = timestamp
		cd.longClickActive = true
	}

	var events []tea.Msg

	// Send press event if action is registered
	if action, exists := cd.buttonPressActions[buttonID]; exists && action.Handler != nil {
		events = append(events, action.Handler(pressEvent))
	}

	return events
}

// HandleButtonRelease processes a button release event
func (cd *ClickDetector) HandleButtonRelease(x, y int, button tea.MouseButton, timestamp int64) []tea.Msg {
	releasedButton := cd.state.GetButtonAtPosition(x, y)
	clickedButton := cd.state.EndPress(x, y)

	var events []tea.Msg

	// Determine if this was a valid click
	wasClick := cd.isValidClick(clickedButton, x, y)
	wasLongClick := cd.wasLongClick(timestamp)

	// Create release event
	releaseEvent := ReleaseEvent{
		ButtonID:     releasedButton,
		X:            x,
		Y:            y,
		Button:       button,
		WasClick:     wasClick,
		WasLongClick: wasLongClick,
		Timestamp:    timestamp,
	}

	// Send release event if action is registered
	if action, exists := cd.buttonReleaseActions[releasedButton]; exists && action.Handler != nil {
		events = append(events, action.Handler(releaseEvent))
	}

	// Handle click detection
	if wasClick && clickedButton != "" {
		clickEvents := cd.handleClickDetection(clickedButton, x, y, button, timestamp)
		events = append(events, clickEvents...)
	}

	// Reset long click state
	cd.longClickActive = false
	cd.longClickStartTime = 0

	return events
}

// HandleMouseMove processes mouse movement during button press
func (cd *ClickDetector) HandleMouseMove(x, y int, button tea.MouseButton, timestamp int64) []tea.Msg {
	if cd.state.PressedButton == "" {
		return nil
	}

	// Check if this movement constitutes a drag
	dx := abs(x - cd.state.PressedX)
	dy := abs(y - cd.state.PressedY)

	if dx > cd.clickThreshold || dy > cd.clickThreshold {
		// Movement exceeded threshold, consider this a drag
		cd.state.ResetPress()
		return []tea.Msg{
			DragEvent{
				StartX:    cd.state.PressedX,
				StartY:    cd.state.PressedY,
				CurrentX:  x,
				CurrentY:  y,
				ButtonID:  cd.state.PressedButton,
				Button:    button,
				Timestamp: timestamp,
			},
		}
	}

	return nil
}

// UpdateLongClick checks for long click activation
func (cd *ClickDetector) UpdateLongClick(timestamp int64) []tea.Msg {
	if !cd.longClickActive || cd.state.PressedButton == "" {
		return nil
	}

	elapsed := time.Duration(timestamp - cd.longClickStartTime)
	if elapsed >= cd.longClickDelay {
		// Long click detected
		cd.longClickActive = false

		return []tea.Msg{
			ClickEvent{
				Type:       ClickLong,
				ButtonID:   cd.state.PressedButton,
				X:          cd.state.PressedX,
				Y:          cd.state.PressedY,
				Button:     cd.state.Button,
				ClickCount: 1,
				Timestamp:  timestamp,
			},
		}
	}

	return nil
}

// isValidClick checks if a button release constitutes a valid click
func (cd *ClickDetector) isValidClick(buttonID string, x, y int) bool {
	if buttonID == "" {
		return false
	}

	// Check if release position is within tolerance of press position
	dx := abs(x - cd.state.PressedX)
	dy := abs(y - cd.state.PressedY)

	return dx <= cd.clickTolerance && dy <= cd.clickTolerance
}

// wasLongClick checks if this was a long click
func (cd *ClickDetector) wasLongClick(timestamp int64) bool {
	if cd.longClickStartTime == 0 {
		return false
	}
	elapsed := time.Duration(timestamp - cd.longClickStartTime)
	return elapsed >= cd.longClickDelay
}

// handleClickDetection handles click type detection
func (cd *ClickDetector) handleClickDetection(buttonID string, x, y int, button tea.MouseButton, timestamp int64) []tea.Msg {
	var events []tea.Msg

	// Determine click type
	clickType := cd.determineClickType(buttonID, x, y, timestamp)

	// Create click event
	clickEvent := ClickEvent{
		Type:       clickType,
		ButtonID:   buttonID,
		X:          x,
		Y:          y,
		Button:     button,
		ClickCount: cd.clickCount,
		Timestamp:  timestamp,
	}

	events = append(events, clickEvent)

	return events
}

// determineClickType determines the type of click based on timing and position
func (cd *ClickDetector) determineClickType(buttonID string, x, y int, timestamp int64) ClickType {
	timeSinceLastClick := time.Duration(timestamp - cd.lastClickTime)
	distanceFromLastClick := abs(x-cd.lastClickPosition.X) + abs(y-cd.lastClickPosition.Y)

	// Check if this is a multiple click
	if timeSinceLastClick <= cd.doubleClickDelay &&
	   distanceFromLastClick <= cd.clickTolerance &&
	   buttonID == cd.state.PressedButton {

		cd.clickCount++
		if cd.clickCount == 2 {
			return ClickDouble
		} else if cd.clickCount >= 3 {
			return ClickTriple
		}
	} else {
		cd.clickCount = 1
	}

	// Update tracking
	cd.lastClickTime = timestamp
	cd.lastClickPosition = ClickPosition{X: x, Y: y}

	return ClickSingle
}

// recordPressEvent records a press event for drag detection
func (cd *ClickDetector) recordPressEvent(event PressEvent) {
	cd.pressHistory = append(cd.pressHistory, event)

	// Limit history size
	if len(cd.pressHistory) > cd.maxPressHistory {
		cd.pressHistory = cd.pressHistory[1:]
	}
}

// GetPressedButton returns the currently pressed button ID
func (cd *ClickDetector) GetPressedButton() string {
	return cd.state.PressedButton
}

// GetPressPosition returns the position where the current press started
func (cd *ClickDetector) GetPressPosition() (int, int) {
	return cd.state.PressedX, cd.state.PressedY
}

// IsLongClickActive returns true if long click detection is active
func (cd *ClickDetector) IsLongClickActive() bool {
	return cd.longClickActive
}

// Reset resets the click detector state
func (cd *ClickDetector) Reset() {
	cd.state.ResetPress()
	cd.clickCount = 0
	cd.lastClickTime = 0
	cd.lastClickPosition = ClickPosition{X: -1, Y: -1}
	cd.longClickActive = false
	cd.longClickStartTime = 0
	cd.pressHistory = make([]PressEvent, 0)
}

// ClearActions clears all registered button actions
func (cd *ClickDetector) ClearActions() {
	cd.buttonPressActions = make(map[string]PressAction)
	cd.buttonReleaseActions = make(map[string]ReleaseAction)
}

// DragEvent represents a drag interaction event
type DragEvent struct {
	StartX    int
	StartY    int
	CurrentX  int
	CurrentY  int
	ButtonID  string
	Button    tea.MouseButton
	Timestamp int64
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}