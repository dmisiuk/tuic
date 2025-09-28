package input

import (
	"testing"
	"ccpm-demo/internal/ui"
)

// TestFocusManager_Basic tests basic focus management functionality
func TestFocusManager_Basic(t *testing.T) {
	// Create a simple test button that implements Focusable
	button := &TestButton{
		id:    "test-btn",
		label: "Test",
		row:   0,
		col:   0,
	}

	fm := NewFocusManager()
	fm.AddFocusable(button)

	// Test setting focus
	if !fm.SetFocus("test-btn") {
		t.Error("Failed to set focus")
	}

	focused := fm.GetFocusedElement()
	if focused == nil {
		t.Error("No element is focused")
	}

	if focused.GetID() != "test-btn" {
		t.Errorf("Expected focused ID to be 'test-btn', got '%s'", focused.GetID())
	}

	if !button.Focused {
		t.Error("Button should be focused")
	}
}

// TestFocusManager_Navigation tests navigation functionality
func TestFocusManager_Navigation(t *testing.T) {
	// Create multiple test buttons
	buttons := []*TestButton{
		{id: "btn1", label: "1", row: 0, col: 0},
		{id: "btn2", label: "2", row: 0, col: 1},
		{id: "btn3", label: "3", row: 1, col: 0},
	}

	fm := NewFocusManager()
	for _, btn := range buttons {
		fm.AddFocusable(btn)
	}

	// Test initial focus
	if !fm.SetFocus("btn1") {
		t.Error("Failed to set initial focus")
	}

	// Test navigation right
	if !fm.Navigate("right") {
		t.Error("Failed to navigate right")
	}

	focused := fm.GetFocusedElement()
	if focused == nil || focused.GetID() != "btn2" {
		t.Error("Expected btn2 to be focused")
	}

	// Test navigation down
	if !fm.Navigate("down") {
		t.Error("Failed to navigate down")
	}

	focused = fm.GetFocusedElement()
	if focused == nil || focused.GetID() != "btn3" {
		t.Error("Expected btn3 to be focused")
	}
}

// TestFocusManager_TabNavigationSimple tests tab navigation
func TestFocusManager_TabNavigationSimple(t *testing.T) {
	buttons := []*TestButton{
		{id: "btn1", label: "1", row: 0, col: 0},
		{id: "btn2", label: "2", row: 0, col: 1},
		{id: "btn3", label: "3", row: 1, col: 0},
	}

	fm := NewFocusManager()
	for _, btn := range buttons {
		fm.AddFocusable(btn)
	}

	// Test tab navigation (next)
	if !fm.Navigate("next") {
		t.Error("Failed to navigate next")
	}

	focused := fm.GetFocusedElement()
	if focused == nil || focused.GetID() != "btn1" {
		t.Error("Expected btn1 to be focused first")
	}

	// Navigate to next
	if !fm.Navigate("next") {
		t.Error("Failed to navigate next")
	}

	focused = fm.GetFocusedElement()
	if focused == nil || focused.GetID() != "btn2" {
		t.Error("Expected btn2 to be focused")
	}

	// Navigate to next (should wrap around)
	if !fm.Navigate("next") {
		t.Error("Failed to navigate next")
	}

	focused = fm.GetFocusedElement()
	if focused == nil || focused.GetID() != "btn3" {
		t.Error("Expected btn3 to be focused")
	}

	// Navigate to next (should wrap around to first)
	if !fm.Navigate("next") {
		t.Error("Failed to navigate next")
	}

	focused = fm.GetFocusedElement()
	if focused == nil || focused.GetID() != "btn1" {
		t.Error("Expected btn1 to be focused after wrap-around")
	}
}

// TestButton is a simple implementation of Focusable for testing
type TestButton struct {
	id      string
	label   string
	row     int
	col     int
	enabled bool
	Focused bool
}

func (b *TestButton) GetID() string          { return b.id }
func (b *TestButton) GetPosition() (int, int) { return b.row, b.col }
func (b *TestButton) GetLabel() string       { return b.label }
func (b *TestButton) IsEnabled() bool        { return b.enabled }
func (b *TestButton) OnFocus()               { b.Focused = true }
func (b *TestButton) OnBlur()                { b.Focused = false }
func (b *TestButton) Activate(model ui.Model) (ui.Model, error) {
	return model, nil
}