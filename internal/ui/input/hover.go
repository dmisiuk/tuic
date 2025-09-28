package input

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// HoverManager manages hover states for UI elements
type HoverManager struct {
	state *MouseState

	// Hover state tracking
	hoverHistory map[string]int64
	hoverThreshold time.Duration

	// Visual feedback
	hoverStyles map[string]string
	hoverEffects map[string]bool
}

// HoverState represents the hover state of an element
type HoverState struct {
	ElementID   string
	IsHovering  bool
	HoverTime   int64
	EnterTime   int64
	ExitTime    int64
	X           int
	Y           int
}

// HoverEffect represents visual effects for hover states
type HoverEffect struct {
	Type        string
	Style       string
	Animation   string
	Priority    int
}

// NewHoverManager creates a new hover manager instance
func NewHoverManager() *HoverManager {
	return &HoverManager{
		state:          NewMouseState(),
		hoverHistory:   make(map[string]int64),
		hoverThreshold: 100 * time.Millisecond, // 100ms threshold before hover is considered "active"
		hoverStyles:    make(map[string]string),
		hoverEffects:   make(map[string]bool),
	}
}

// SetHoverThreshold sets the time threshold for hover detection
func (hm *HoverManager) SetHoverThreshold(duration time.Duration) {
	hm.hoverThreshold = duration
}

// RegisterHoverStyle registers a style for hover effects
func (hm *HoverManager) RegisterHoverStyle(elementID, style string) {
	hm.hoverStyles[elementID] = style
}

// RegisterHoverEffect enables hover effects for an element
func (hm *HoverManager) RegisterHoverEffect(elementID string, enabled bool) {
	hm.hoverEffects[elementID] = enabled
}

// UpdateHover updates the hover state based on mouse position
func (hm *HoverManager) UpdateHover(x, y int, timestamp int64) *HoverState {
	hoveredButton := hm.state.UpdateHover(x, y)

	// Create hover state object
	hoverState := &HoverState{
		ElementID:  hoveredButton,
		IsHovering: hoveredButton != "",
		X:          x,
		Y:          y,
	}

	// Update hover history
	if hoveredButton != "" {
		if _, exists := hm.hoverHistory[hoveredButton]; !exists {
			// New hover - set enter time
			hoverState.EnterTime = timestamp
		}
		hm.hoverHistory[hoveredButton] = timestamp
		hoverState.HoverTime = timestamp
	} else {
		// Clear old hover history
		hm.cleanupOldHoverHistory(timestamp)
	}

	return hoverState
}

// IsHovering returns true if an element is currently being hovered
func (hm *HoverManager) IsHovering(elementID string) bool {
	return hm.state.HoveredButton == elementID && hm.state.IsHovering
}

// GetHoverDuration returns how long an element has been hovered
func (hm *HoverManager) GetHoverDuration(elementID string, currentTimestamp int64) time.Duration {
	hoverTime, exists := hm.hoverHistory[elementID]
	if !exists {
		return 0
	}
	return time.Duration(currentTimestamp - hoverTime)
}

// IsHoverActive returns true if hover has been active for longer than threshold
func (hm *HoverManager) IsHoverActive(elementID string, currentTimestamp int64) bool {
	duration := hm.GetHoverDuration(elementID, currentTimestamp)
	return duration >= hm.hoverThreshold
}

// GetHoveredElement returns the currently hovered element ID
func (hm *HoverManager) GetHoveredElement() string {
	return hm.state.HoveredButton
}

// GetHoverPosition returns the current mouse position during hover
func (hm *HoverManager) GetHoverPosition() (int, int) {
	return hm.state.X, hm.state.Y
}

// GetHoverStyle returns the registered hover style for an element
func (hm *HoverManager) GetHoverStyle(elementID string) string {
	if style, exists := hm.hoverStyles[elementID]; exists {
		return style
	}
	return ""
}

// HasHoverEffect returns true if an element has hover effects enabled
func (hm *HoverManager) HasHoverEffect(elementID string) bool {
	return hm.hoverEffects[elementID]
}

// GetHoveredElements returns all elements currently being hovered
func (hm *HoverManager) GetHoveredElements() []string {
	var hovered []string
	if hm.state.IsHovering && hm.state.HoveredButton != "" {
		hovered = append(hovered, hm.state.HoveredButton)
	}
	return hovered
}

// ProcessHoverEvent processes a mouse move event and returns hover-related events
func (hm *HoverManager) ProcessHoverEvent(msg tea.MouseMsg, timestamp int64) []tea.Msg {
	oldHovered := hm.state.HoveredButton
	hoverState := hm.UpdateHover(msg.X, msg.Y, timestamp)

	var events []tea.Msg

	// Hover enter event
	if hoverState.IsHovering && hoverState.ElementID != oldHovered {
		events = append(events, HoverEvent{
			Type:      HoverEnter,
			ElementID: hoverState.ElementID,
			X:         msg.X,
			Y:         msg.Y,
			Timestamp: timestamp,
		})

		// If there was a previously hovered element, send exit event
		if oldHovered != "" {
			events = append(events, HoverEvent{
				Type:      HoverExit,
				ElementID: oldHovered,
				X:         msg.X,
				Y:         msg.Y,
				Timestamp: timestamp,
			})
		}
	}

	// Hover exit event
	if !hoverState.IsHovering && oldHovered != "" {
		events = append(events, HoverEvent{
			Type:      HoverExit,
			ElementID: oldHovered,
			X:         msg.X,
			Y:         msg.Y,
			Timestamp: timestamp,
		})
	}

	return events
}

// cleanupOldHoverHistory removes hover entries older than threshold
func (hm *HoverManager) cleanupOldHoverHistory(currentTimestamp int64) {
	threshold := int64(hm.hoverThreshold)
	for elementID, hoverTime := range hm.hoverHistory {
		if currentTimestamp-hoverTime > threshold {
			delete(hm.hoverHistory, elementID)
		}
	}
}

// Reset resets the hover manager state
func (hm *HoverManager) Reset() {
	hm.state.HoveredButton = ""
	hm.state.IsHovering = false
	hm.state.X = -1
	hm.state.Y = -1
	hm.hoverHistory = make(map[string]int64)
}

// ClearStyles clears all registered hover styles
func (hm *HoverManager) ClearStyles() {
	hm.hoverStyles = make(map[string]string)
}

// ClearEffects clears all registered hover effects
func (hm *HoverManager) ClearEffects() {
	hm.hoverEffects = make(map[string]bool)
}

// HoverEventType represents the type of hover event
type HoverEventType int

const (
	HoverEnter HoverEventType = iota
	HoverExit
	HoverMove
)

// HoverEvent represents a hover interaction event
type HoverEvent struct {
	Type      HoverEventType
	ElementID string
	X         int
	Y         int
	Timestamp int64
}