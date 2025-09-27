package input

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMouseState_NewMouseState(t *testing.T) {
	state := NewMouseState()

	if state.X != -1 {
		t.Errorf("Expected X to be -1, got %d", state.X)
	}
	if state.Y != -1 {
		t.Errorf("Expected Y to be -1, got %d", state.Y)
	}
	if state.IsHovering {
		t.Error("Expected IsHovering to be false")
	}
	if len(state.ButtonBounds) != 0 {
		t.Errorf("Expected ButtonBounds to be empty, got %d items", len(state.ButtonBounds))
	}
}

func TestMouseState_RegisterButton(t *testing.T) {
	state := NewMouseState()

	// Register a button
	state.RegisterButton("test", 10, 20, 30, 40)

	if len(state.ButtonBounds) != 1 {
		t.Errorf("Expected 1 button, got %d", len(state.ButtonBounds))
	}

	bounds := state.ButtonBounds["test"]
	if bounds.X != 10 || bounds.Y != 20 || bounds.Width != 30 || bounds.Height != 40 {
		t.Errorf("Expected bounds {10,20,30,40}, got {%d,%d,%d,%d}", bounds.X, bounds.Y, bounds.Width, bounds.Height)
	}
}

func TestMouseState_GetButtonAtPosition(t *testing.T) {
	state := NewMouseState()

	// Register buttons
	state.RegisterButton("button1", 10, 10, 20, 10)
	state.RegisterButton("button2", 50, 50, 30, 15)

	// Test positions
	tests := []struct {
		name     string
		x, y     int
		expected string
	}{
		{"Inside button1", 15, 15, "button1"},
		{"Inside button2", 60, 55, "button2"},
		{"Outside all", 100, 100, ""},
		{"Edge of button1", 29, 19, "button1"},
		{"Just outside button1", 30, 20, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := state.GetButtonAtPosition(test.x, test.y)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestMouseState_IsPositionInButton(t *testing.T) {
	state := NewMouseState()
	state.RegisterButton("test", 10, 10, 20, 10)

	tests := []struct {
		name     string
		x, y     int
		buttonID string
		expected bool
	}{
		{"Inside correct button", 15, 15, "test", true},
		{"Outside correct button", 5, 5, "test", false},
		{"Inside wrong button", 15, 15, "wrong", false},
		{"Non-existent button", 15, 15, "nonexistent", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := state.IsPositionInButton(test.x, test.y, test.buttonID)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestMouseState_UpdateHover(t *testing.T) {
	state := NewMouseState()
	state.RegisterButton("test", 10, 10, 20, 10)

	// Test hover over button
	result := state.UpdateHover(15, 15)
	if result != "test" {
		t.Errorf("Expected hover over test button, got %q", result)
	}
	if !state.IsHovering {
		t.Error("Expected IsHovering to be true")
	}
	if state.HoveredButton != "test" {
		t.Errorf("Expected HoveredButton to be test, got %q", state.HoveredButton)
	}

	// Test hover outside button
	result = state.UpdateHover(100, 100)
	if result != "" {
		t.Errorf("Expected no hover, got %q", result)
	}
	if state.IsHovering {
		t.Error("Expected IsHovering to be false")
	}
	if state.HoveredButton != "" {
		t.Errorf("Expected HoveredButton to be empty, got %q", state.HoveredButton)
	}
}

func TestMouseState_PressRelease(t *testing.T) {
	state := NewMouseState()
	state.RegisterButton("test", 10, 10, 20, 10)

	// Start press
	state.StartPress(15, 15, tea.MouseButtonLeft)
	if state.PressedButton != "test" {
		t.Errorf("Expected PressedButton to be test, got %q", state.PressedButton)
	}
	if state.PressedX != 15 || state.PressedY != 15 {
		t.Errorf("Expected press position (15,15), got (%d,%d)", state.PressedX, state.PressedY)
	}

	// End press (successful click)
	result := state.EndPress(15, 15)
	if result != "test" {
		t.Errorf("Expected click on test button, got %q", result)
	}

	// End press (invalid click - different position)
	state.StartPress(15, 15, tea.MouseButtonLeft)
	result = state.EndPress(100, 100)
	if result != "" {
		t.Errorf("Expected no click, got %q", result)
	}
}

func TestMouseHandler_NewMouseHandler(t *testing.T) {
	handler := NewMouseHandler()

	if !handler.enabled {
		t.Error("Expected handler to be enabled")
	}
	if handler.clickTolerance != 5 {
		t.Errorf("Expected click tolerance 5, got %d", handler.clickTolerance)
	}
	if handler.doubleClickDelay != 500*time.Millisecond {
		t.Errorf("Expected double click delay 500ms, got %v", handler.doubleClickDelay)
	}
	if len(handler.buttonActions) != 0 {
		t.Errorf("Expected no button actions, got %d", len(handler.buttonActions))
	}
}

func TestMouseHandler_RegisterButton(t *testing.T) {
	handler := NewMouseHandler()
	action := ButtonAction{Type: "test", Value: "value"}

	handler.RegisterButton("test", 10, 10, 20, 10, action)

	if len(handler.state.ButtonBounds) != 1 {
		t.Errorf("Expected 1 button bound, got %d", len(handler.state.ButtonBounds))
	}
	if len(handler.buttonActions) != 1 {
		t.Errorf("Expected 1 button action, got %d", len(handler.buttonActions))
	}
	if _, exists := handler.buttonActions["test"]; !exists {
		t.Error("Expected button action to be registered")
	}
}

func TestMouseHandler_HandleMessage_InvalidType(t *testing.T) {
	handler := NewMouseHandler()
	msg := tea.KeyMsg{} // Not a mouse message

	events := handler.HandleMessage(msg)
	if len(events) != 0 {
		t.Errorf("Expected no events for non-mouse message, got %d", len(events))
	}
}

func TestMouseHandler_HandleMessage_Disabled(t *testing.T) {
	handler := NewMouseHandler()
	handler.SetEnabled(false)
	msg := tea.MouseMsg{Type: tea.MouseMotion, X: 10, Y: 10}

	events := handler.HandleMessage(msg)
	if len(events) != 0 {
		t.Errorf("Expected no events when disabled, got %d", len(events))
	}
}

func TestMouseHandler_HandleMouseMove(t *testing.T) {
	handler := NewMouseHandler()
	handler.RegisterButton("test", 10, 10, 20, 10, ButtonAction{})

	msg := tea.MouseMsg{Type: tea.MouseMotion, X: 15, Y: 15}
	events := handler.HandleMessage(msg)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	event, ok := events[0].(MouseEvent)
	if !ok {
		t.Fatal("Expected MouseEvent")
	}
	if event.Type != MouseEventMove {
		t.Errorf("Expected MouseEventMove, got %v", event.Type)
	}
	if event.ButtonID != "test" {
		t.Errorf("Expected buttonID test, got %q", event.ButtonID)
	}
}

func TestMouseHandler_HandleMousePress(t *testing.T) {
	handler := NewMouseHandler()
	handler.RegisterButton("test", 10, 10, 20, 10, ButtonAction{})

	msg := tea.MouseMsg{Type: tea.MousePress, X: 15, Y: 15, Button: tea.MouseButtonLeft}
	events := handler.HandleMessage(msg)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	event, ok := events[0].(MouseEvent)
	if !ok {
		t.Fatal("Expected MouseEvent")
	}
	if event.Type != MouseEventPress {
		t.Errorf("Expected MouseEventPress, got %v", event.Type)
	}
	if event.Button != tea.MouseButtonLeft {
		t.Errorf("Expected MouseButtonLeft, got %v", event.Button)
	}
}

func TestMouseHandler_HandleMouseRelease(t *testing.T) {
	handler := NewMouseHandler()
	handler.RegisterButton("test", 10, 10, 20, 10, ButtonAction{})

	// First press to set up the state
	handler.HandleMessage(tea.MouseMsg{Type: tea.MousePress, X: 15, Y: 15})

	msg := tea.MouseMsg{Type: tea.MouseRelease, X: 15, Y: 15}
	events := handler.HandleMessage(msg)

	if len(events) != 2 {
		t.Errorf("Expected 2 events (release + click), got %d", len(events))
	}

	// Check release event
	releaseEvent, ok := events[0].(MouseEvent)
	if !ok {
		t.Fatal("Expected MouseEvent")
	}
	if releaseEvent.Type != MouseEventRelease {
		t.Errorf("Expected MouseEventRelease, got %v", releaseEvent.Type)
	}

	// Check click event
	clickEvent, ok := events[1].(MouseEvent)
	if !ok {
		t.Fatal("Expected MouseEvent")
	}
	if clickEvent.Type != MouseEventClick {
		t.Errorf("Expected MouseEventClick, got %v", clickEvent.Type)
	}
}

func TestMouseHandler_HandleMouseWheel(t *testing.T) {
	handler := NewMouseHandler()

	// Test wheel up
	msg := tea.MouseMsg{Type: tea.MouseWheelUp, X: 10, Y: 10}
	events := handler.HandleMessage(msg)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	event, ok := events[0].(MouseEvent)
	if !ok {
		t.Fatal("Expected MouseEvent")
	}
	if event.Type != MouseEventScroll {
		t.Errorf("Expected MouseEventScroll, got %v", event.Type)
	}
	if event.Scroll != 1 {
		t.Errorf("Expected scroll delta 1, got %d", event.Scroll)
	}

	// Test wheel down
	msg = tea.MouseMsg{Type: tea.MouseWheelDown, X: 10, Y: 10}
	events = handler.HandleMessage(msg)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	event = events[0].(MouseEvent)
	if event.Scroll != -1 {
		t.Errorf("Expected scroll delta -1, got %d", event.Scroll)
	}
}

func TestMouseHandler_HelperMethods(t *testing.T) {
	handler := NewMouseHandler()
	handler.RegisterButton("test", 10, 10, 20, 10, ButtonAction{})

	// Test hover state
	handler.state.UpdateHover(15, 15)
	if !handler.IsHovering() {
		t.Error("Expected IsHovering to return true")
	}
	if handler.GetHoveredButton() != "test" {
		t.Errorf("Expected hovered button test, got %q", handler.GetHoveredButton())
	}

	// Test press state
	handler.state.StartPress(15, 15, tea.MouseButtonLeft)
	if handler.GetPressedButton() != "test" {
		t.Errorf("Expected pressed button test, got %q", handler.GetPressedButton())
	}

	// Test position
	x, y := handler.GetMousePosition()
	if x != 15 || y != 15 {
		t.Errorf("Expected position (15,15), got (%d,%d)", x, y)
	}

	// Test scroll
	handler.state.UpdateScroll(20, 20, 5)
	if handler.GetScrollDelta() != 5 {
		t.Errorf("Expected scroll delta 5, got %d", handler.GetScrollDelta())
	}
}

func TestClickDetector_NewClickDetector(t *testing.T) {
	detector := NewClickDetector()

	if detector.clickTolerance != 5 {
		t.Errorf("Expected click tolerance 5, got %d", detector.clickTolerance)
	}
	if detector.doubleClickDelay != 500*time.Millisecond {
		t.Errorf("Expected double click delay 500ms, got %v", detector.doubleClickDelay)
	}
	if detector.longClickDelay != 1000*time.Millisecond {
		t.Errorf("Expected long click delay 1000ms, got %v", detector.longClickDelay)
	}
	if detector.clickCount != 0 {
		t.Errorf("Expected click count 0, got %d", detector.clickCount)
	}
}

func TestClickDetector_IsValidClick(t *testing.T) {
	detector := NewClickDetector()
	detector.state.RegisterButton("test", 10, 10, 20, 10)

	// Valid click (within tolerance)
	detector.state.StartPress(15, 15, tea.MouseButtonLeft)
	valid := detector.isValidClick("test", 17, 17) // Within tolerance
	if !valid {
		t.Error("Expected valid click")
	}

	// Invalid click (outside tolerance)
	detector.state.StartPress(15, 15, tea.MouseButtonLeft)
	valid = detector.isValidClick("test", 25, 25) // Outside tolerance
	if valid {
		t.Error("Expected invalid click")
	}

	// No button
	valid = detector.isValidClick("", 15, 15)
	if valid {
		t.Error("Expected invalid click for no button")
	}
}

func TestScrollManager_NewScrollManager(t *testing.T) {
	manager := NewScrollManager()

	if manager.scrollSensitivity != 1.0 {
		t.Errorf("Expected sensitivity 1.0, got %f", manager.scrollSensitivity)
	}
	if manager.scrollThreshold != 1 {
		t.Errorf("Expected threshold 1, got %d", manager.scrollThreshold)
	}
	if !manager.scrollSmoothing {
		t.Error("Expected smoothing to be enabled")
	}
	if manager.smoothingFactor != 0.8 {
		t.Errorf("Expected smoothing factor 0.8, got %f", manager.smoothingFactor)
	}
}

func TestHoverManager_NewHoverManager(t *testing.T) {
	manager := NewHoverManager()

	if manager.hoverThreshold != 100*time.Millisecond {
		t.Errorf("Expected hover threshold 100ms, got %v", manager.hoverThreshold)
	}
	if len(manager.hoverStyles) != 0 {
		t.Errorf("Expected no hover styles, got %d", len(manager.hoverStyles))
	}
	if len(manager.hoverEffects) != 0 {
		t.Errorf("Expected no hover effects, got %d", len(manager.hoverEffects))
	}
}

func TestHoverManager_UpdateHover(t *testing.T) {
	manager := NewHoverManager()
	manager.state.RegisterButton("test", 10, 10, 20, 10)

	state := manager.UpdateHover(15, 15, 1000)

	if !state.IsHovering {
		t.Error("Expected IsHovering to be true")
	}
	if state.ElementID != "test" {
		t.Errorf("Expected ElementID test, got %q", state.ElementID)
	}
	if state.X != 15 || state.Y != 15 {
		t.Errorf("Expected position (15,15), got (%d,%d)", state.X, state.Y)
	}
}

func TestMouseEventTypes(t *testing.T) {
	// Test that mouse event types have correct values
	if MouseEventMove != 0 {
		t.Errorf("Expected MouseEventMove = 0, got %d", MouseEventMove)
	}
	if MouseEventPress != 1 {
		t.Errorf("Expected MouseEventPress = 1, got %d", MouseEventPress)
	}
	if MouseEventRelease != 2 {
		t.Errorf("Expected MouseEventRelease = 2, got %d", MouseEventRelease)
	}
	if MouseEventClick != 3 {
		t.Errorf("Expected MouseEventClick = 3, got %d", MouseEventClick)
	}
	if MouseEventScroll != 5 {
		t.Errorf("Expected MouseEventScroll = 5, got %d", MouseEventScroll)
	}
}

func TestButtonRect_Equality(t *testing.T) {
	rect1 := ButtonRect{X: 10, Y: 20, Width: 30, Height: 40}
	rect2 := ButtonRect{X: 10, Y: 20, Width: 30, Height: 40}

	if rect1.X != rect2.X || rect1.Y != rect2.Y || rect1.Width != rect2.Width || rect1.Height != rect2.Height {
		t.Error("Expected button rects to be equal")
	}
}

func TestMouseState_ResetPress(t *testing.T) {
	state := NewMouseState()
	state.StartPress(10, 10, tea.MouseButtonLeft)

	state.ResetPress()

	if state.PressedButton != "" {
		t.Errorf("Expected PressedButton to be empty, got %q", state.PressedButton)
	}
	if state.PressedX != -1 || state.PressedY != -1 {
		t.Errorf("Expected press position (-1,-1), got (%d,%d)", state.PressedX, state.PressedY)
	}
}

func TestMouseState_ResetScroll(t *testing.T) {
	state := NewMouseState()
	state.UpdateScroll(10, 10, 5)

	state.ResetScroll()

	if state.ScrollDelta != 0 {
		t.Errorf("Expected ScrollDelta to be 0, got %d", state.ScrollDelta)
	}
	if state.ScrollX != -1 || state.ScrollY != -1 {
		t.Errorf("Expected scroll position (-1,-1), got (%d,%d)", state.ScrollX, state.ScrollY)
	}
}