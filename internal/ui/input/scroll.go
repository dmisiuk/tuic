package input

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ScrollManager handles scroll wheel interactions
type ScrollManager struct {
	state *MouseState

	// Scroll behavior configuration
	scrollSensitivity float64
	scrollThreshold   int
	scrollSmoothing   bool
	smoothingFactor   float64

	// Scroll state tracking
	accumulatedDelta float64
	lastScrollTime   int64
	scrollVelocity   float64

	// Scroll actions
	scrollActions map[string]ScrollAction

	// Momentum scrolling
	momentumEnabled bool
	momentumDecay   float64
	momentumActive  bool
	momentumDelta   float64
}

// ScrollAction represents an action triggered by scrolling
type ScrollAction struct {
	Type        string
	Value       float64
	Handler     func(delta float64) tea.Msg
	MinDelta    float64
	MaxDelta    float64
	Direction   ScrollDirection
}

// ScrollDirection represents the direction of scrolling
type ScrollDirection int

const (
	ScrollVertical ScrollDirection = iota
	ScrollHorizontal
	ScrollBoth
)

// ScrollEvent represents a scroll interaction event
type ScrollEvent struct {
	Type      ScrollEventType
	X         int
	Y         int
	Delta     float64
	Direction ScrollDirection
	Action    ScrollAction
}

// ScrollEventType represents the type of scroll event
type ScrollEventType int

const (
	ScrollStart ScrollEventType = iota
	ScrollContinue
	ScrollEnd
	ScrollMomentum
)

// NewScrollManager creates a new scroll manager instance
func NewScrollManager() *ScrollManager {
	return &ScrollManager{
		state:            NewMouseState(),
		scrollSensitivity: 1.0,
		scrollThreshold:   1,
		scrollSmoothing:   true,
		smoothingFactor:   0.8,
		accumulatedDelta:  0,
		lastScrollTime:    0,
		scrollVelocity:    0,
		scrollActions:     make(map[string]ScrollAction),
		momentumEnabled:   false,
		momentumDecay:     0.95,
		momentumActive:    false,
		momentumDelta:     0,
	}
}

// SetScrollSensitivity adjusts how sensitive the scroll wheel is
func (sm *ScrollManager) SetScrollSensitivity(sensitivity float64) {
	sm.scrollSensitivity = sensitivity
}

// SetScrollThreshold sets the minimum delta before triggering scroll events
func (sm *ScrollManager) SetScrollThreshold(threshold int) {
	sm.scrollThreshold = threshold
}

// EnableSmoothing enables or disables scroll smoothing
func (sm *ScrollManager) EnableSmoothing(enabled bool) {
	sm.scrollSmoothing = enabled
	if !enabled {
		sm.accumulatedDelta = 0
		sm.scrollVelocity = 0
	}
}

// SetSmoothingFactor adjusts the smoothing factor (0-1, higher = smoother)
func (sm *ScrollManager) SetSmoothingFactor(factor float64) {
	if factor >= 0 && factor <= 1 {
		sm.smoothingFactor = factor
	}
}

// EnableMomentum enables or disables momentum scrolling
func (sm *ScrollManager) EnableMomentum(enabled bool) {
	sm.momentumEnabled = enabled
	if !enabled {
		sm.momentumActive = false
		sm.momentumDelta = 0
	}
}

// SetMomentumDecay sets the momentum decay rate (0-1, higher = longer momentum)
func (sm *ScrollManager) SetMomentumDecay(decay float64) {
	if decay >= 0 && decay <= 1 {
		sm.momentumDecay = decay
	}
}

// RegisterScrollAction registers a scroll action for a specific context
func (sm *ScrollManager) RegisterScrollAction(context string, action ScrollAction) {
	sm.scrollActions[context] = action
}

// UnregisterScrollAction removes a scroll action
func (sm *ScrollManager) UnregisterScrollAction(context string) {
	delete(sm.scrollActions, context)
}

// HandleScroll processes scroll wheel events and returns scroll events
func (sm *ScrollManager) HandleScroll(msg tea.MouseMsg, timestamp int64) []tea.Msg {
	var delta float64
	if msg.Type == tea.MouseWheelUp {
		delta = 1.0
	} else if msg.Type == tea.MouseWheelDown {
		delta = -1.0
	} else {
		return nil
	}

	// Apply sensitivity
	delta *= sm.scrollSensitivity

	// Store scroll position
	sm.state.UpdateScroll(msg.X, msg.Y, int(delta))

	// Handle momentum
	if sm.momentumEnabled {
		sm.handleMomentum(delta, timestamp)
	}

	// Smooth scrolling
	if sm.scrollSmoothing {
		return sm.handleSmoothScroll(msg, delta, timestamp)
	}

	return sm.handleDirectScroll(msg, delta, timestamp)
}

// handleSmoothScroll processes scroll with smoothing
func (sm *ScrollManager) handleSmoothScroll(msg tea.MouseMsg, delta float64, timestamp int64) []tea.Msg {
	sm.accumulatedDelta += delta

	// Calculate scroll velocity
	timeDelta := float64(timestamp - sm.lastScrollTime)
	if timeDelta > 0 {
		sm.scrollVelocity = delta / timeDelta
	}

	sm.lastScrollTime = timestamp

	var events []tea.Msg

	// Only trigger events when accumulated delta exceeds threshold
	if abs(sm.accumulatedDelta) >= float64(sm.scrollThreshold) {
		scrollDelta := sm.accumulatedDelta * sm.smoothingFactor
		sm.accumulatedDelta -= scrollDelta

		scrollEvents := sm.createScrollEvents(msg.X, msg.Y, scrollDelta, ScrollContinue)
		events = append(events, scrollEvents...)

		// Apply smoothing decay
		sm.accumulatedDelta *= sm.smoothingFactor
	}

	return events
}

// handleDirectScroll processes scroll without smoothing
func (sm *ScrollManager) handleDirectScroll(msg tea.MouseMsg, delta float64, timestamp int64) []tea.Msg {
	sm.lastScrollTime = timestamp
	return sm.createScrollEvents(msg.X, msg.Y, delta, ScrollContinue)
}

// handleMomentum manages momentum scrolling physics
func (sm *ScrollManager) handleMomentum(delta float64, timestamp int64) {
	if delta != 0 {
		// User is actively scrolling, update momentum
		sm.momentumActive = true
		sm.momentumDelta = delta * 2.0 // Amplify for momentum effect
	}
}

// createScrollEvents creates scroll events based on current delta
func (sm *ScrollManager) createScrollEvents(x, y int, delta float64, eventType ScrollEventType) []tea.Msg {
	if abs(delta) < float64(sm.scrollThreshold) {
		return nil
	}

	var events []tea.Msg

	// Create vertical scroll event
	verticalAction := sm.scrollActions["vertical"]
	if verticalAction.Handler != nil {
		events = append(events, ScrollEvent{
			Type:      eventType,
			X:         x,
			Y:         y,
			Delta:     delta,
			Direction: ScrollVertical,
			Action:    verticalAction,
		})

		// Execute action handler
		events = append(events, verticalAction.Handler(delta))
	}

	// Create horizontal scroll event (if Shift is held)
	// Note: This would need to be detected from keyboard state
	horizontalAction := sm.scrollActions["horizontal"]
	if horizontalAction.Handler != nil {
		events = append(events, ScrollEvent{
			Type:      eventType,
			X:         x,
			Y:         y,
			Delta:     delta,
			Direction: ScrollHorizontal,
			Action:    horizontalAction,
		})

		// Execute action handler
		events = append(events, horizontalAction.Handler(delta))
	}

	return events
}

// UpdateMomentum updates momentum scrolling state
func (sm *ScrollManager) UpdateMomentum(timestamp int64) []tea.Msg {
	if !sm.momentumEnabled || !sm.momentumActive {
		return nil
	}

	var events []tea.Msg

	// Apply momentum decay
	sm.momentumDelta *= sm.momentumDecay

	// Stop momentum when it becomes too small
	if abs(sm.momentumDelta) < 0.1 {
		sm.momentumActive = false
		sm.momentumDelta = 0
		events = append(events, ScrollEvent{
			Type:      ScrollEnd,
			X:         sm.state.ScrollX,
			Y:         sm.state.ScrollY,
			Delta:     0,
			Direction: ScrollVertical,
		})
		return events
	}

	// Create momentum scroll event
	x, y := sm.GetScrollPosition()
	momentumEvents := sm.createScrollEvents(x, y, sm.momentumDelta, ScrollMomentum)
	events = append(events, momentumEvents...)

	return events
}

// GetScrollDelta returns the current scroll delta
func (sm *ScrollManager) GetScrollDelta() float64 {
	return float64(sm.state.ScrollDelta)
}

// GetScrollPosition returns the current scroll position
func (sm *ScrollManager) GetScrollPosition() (int, int) {
	return sm.state.ScrollX, sm.state.ScrollY
}

// GetScrollVelocity returns the current scroll velocity
func (sm *ScrollManager) GetScrollVelocity() float64 {
	return sm.scrollVelocity
}

// GetMomentumDelta returns the current momentum delta
func (sm *ScrollManager) GetMomentumDelta() float64 {
	return sm.momentumDelta
}

// IsMomentumActive returns true if momentum scrolling is active
func (sm *ScrollManager) IsMomentumActive() bool {
	return sm.momentumActive
}

// Reset resets the scroll manager state
func (sm *ScrollManager) Reset() {
	sm.accumulatedDelta = 0
	sm.lastScrollTime = 0
	sm.scrollVelocity = 0
	sm.momentumActive = false
	sm.momentumDelta = 0
	sm.state.ResetScroll()
}

// ClearActions clears all registered scroll actions
func (sm *ScrollManager) ClearActions() {
	sm.scrollActions = make(map[string]ScrollAction)
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}