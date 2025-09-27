package input

import (
	"github.com/charmbracelet/bubbletea"
)

// MouseState represents the current state of mouse interactions
type MouseState struct {
	// Current mouse position
	X int
	Y int

	// Button that was pressed (if any)
	Button tea.MouseButton

	// Hover state
	HoveredButton string
	IsHovering    bool

	// Click detection
	PressedButton string
	PressedX     int
	PressedY     int
	ClickCount    int
	LastClickTime int64

	// Wheel scrolling
	ScrollDelta int
	ScrollX     int
	ScrollY     int

	// Button boundaries for hit testing
	ButtonBounds map[string]ButtonRect
}

// ButtonRect represents the rectangular area of a button
type ButtonRect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// ButtonAction represents the action to take when a button is clicked
type ButtonAction struct {
	Type    string
	Value   string
	Handler func() tea.Msg
}

// MouseEvent represents a mouse interaction event
type MouseEvent struct {
	Type     MouseEventType
	Button   tea.MouseButton
	X        int
	Y        int
	Scroll   int
	ButtonID string
	Action   ButtonAction
}

// MouseEventType represents the type of mouse event
type MouseEventType int

const (
	MouseEventMove MouseEventType = iota
	MouseEventPress
	MouseEventRelease
	MouseEventClick
	MouseEventDoubleClick
	MouseEventScroll
)

// NewMouseState creates a new mouse state instance
func NewMouseState() *MouseState {
	return &MouseState{
		X:             -1,
		Y:             -1,
		Button:        tea.MouseButtonLeft,
		HoveredButton: "",
		IsHovering:    false,
		PressedButton: "",
		PressedX:      -1,
		PressedY:      -1,
		ClickCount:    0,
		LastClickTime: 0,
		ScrollDelta:   0,
		ScrollX:       -1,
		ScrollY:       -1,
		ButtonBounds:  make(map[string]ButtonRect),
	}
}

// RegisterButton registers a button's boundaries for hit testing
func (ms *MouseState) RegisterButton(id string, x, y, width, height int) {
	ms.ButtonBounds[id] = ButtonRect{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// UnregisterButton removes a button from hit testing
func (ms *MouseState) UnregisterButton(id string) {
	delete(ms.ButtonBounds, id)
}

// GetButtonAtPosition returns the button ID at the given position
func (ms *MouseState) GetButtonAtPosition(x, y int) string {
	for id, bounds := range ms.ButtonBounds {
		if x >= bounds.X && x < bounds.X+bounds.Width &&
			y >= bounds.Y && y < bounds.Y+bounds.Height {
			return id
		}
	}
	return ""
}

// IsPositionInButton checks if a position is within a button's bounds
func (ms *MouseState) IsPositionInButton(x, y int, buttonID string) bool {
	bounds, exists := ms.ButtonBounds[buttonID]
	if !exists {
		return false
	}
	return x >= bounds.X && x < bounds.X+bounds.Width &&
		y >= bounds.Y && y < bounds.Y+bounds.Height
}

// UpdateHover updates the hover state based on mouse position
func (ms *MouseState) UpdateHover(x, y int) string {
	ms.X = x
	ms.Y = y

	newHovered := ms.GetButtonAtPosition(x, y)

	// Update hover state
	ms.HoveredButton = newHovered
	ms.IsHovering = newHovered != ""

	return newHovered
}

// UpdateScroll updates the scroll state
func (ms *MouseState) UpdateScroll(x, y int, delta int) {
	ms.ScrollX = x
	ms.ScrollY = y
	ms.ScrollDelta = delta
}

// StartPress records the start of a button press
func (ms *MouseState) StartPress(x, y int, button tea.MouseButton) {
	ms.PressedX = x
	ms.PressedY = y
	ms.PressedButton = ms.GetButtonAtPosition(x, y)
	ms.Button = button
}

// EndPress handles the end of a button press and returns the clicked button
func (ms *MouseState) EndPress(x, y int) string {
	releasedButton := ms.GetButtonAtPosition(x, y)

	// Check if this is a valid click (pressed and released on same button)
	if ms.PressedButton == releasedButton && ms.PressedButton != "" {
		// Simple click detection - could be enhanced for double-clicks
		ms.ClickCount++
		return releasedButton
	}

	// Reset press state
	ms.PressedButton = ""
	ms.PressedX = -1
	ms.PressedY = -1
	ms.Button = tea.MouseButtonLeft

	return ""
}

// ResetPress resets the press state
func (ms *MouseState) ResetPress() {
	ms.PressedButton = ""
	ms.PressedX = -1
	ms.PressedY = -1
	ms.Button = tea.MouseButtonLeft
}

// ResetScroll resets the scroll state
func (ms *MouseState) ResetScroll() {
	ms.ScrollDelta = 0
	ms.ScrollX = -1
	ms.ScrollY = -1
}