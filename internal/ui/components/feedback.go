package components

import (
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// FeedbackManager manages visual feedback and animations for button interactions
type FeedbackManager struct {
	mu sync.Mutex

	// Animation settings
	pressDuration   time.Duration
	transitionSpeed time.Duration
	focusAnimation  bool

	// Visual effects
	flashEnabled    bool
	flashDuration   time.Duration
	rippleEnabled   bool

	// Active feedback states
	activeAnimations map[string]*ButtonAnimation
	flashQueue      []FlashEffect
	rippleEffects   []RippleEffect

	// Event handlers
	feedbackHandlers map[string][]func(FeedbackEvent)
}

// ButtonAnimation represents an ongoing button animation
type ButtonAnimation struct {
	Button     *Button
	Type       AnimationType
	StartTime  time.Time
	Duration   time.Duration
	Progress   float64
	Completed  bool
	Properties map[string]interface{}
}

// AnimationType defines types of button animations
type AnimationType int

const (
	AnimPress AnimationType = iota
	AnimFocus
	AnimRelease
	AnimHover
	AnimFlash
	AnimRipple
)

// FlashEffect represents a visual flash effect
type FlashEffect struct {
	Button   *Button
	Color    lipgloss.Color
	StartTime time.Time
	Duration time.Duration
}

// RippleEffect represents a ripple animation effect
type RippleEffect struct {
	Button    *Button
	CenterX   int
	CenterY   int
	StartTime time.Time
	Duration  time.Duration
	Radius    float64
	MaxRadius int
}

// FeedbackEvent represents a feedback system event
type FeedbackEvent struct {
	Type      string
	Button    *Button
	Animation *ButtonAnimation
	Timestamp time.Time
	Details   map[string]interface{}
}

// NewFeedbackManager creates a new feedback manager
func NewFeedbackManager() *FeedbackManager {
	return &FeedbackManager{
		pressDuration:     150 * time.Millisecond,
		transitionSpeed:   100 * time.Millisecond,
		focusAnimation:    true,
		flashEnabled:      true,
		flashDuration:     200 * time.Millisecond,
		rippleEnabled:     false, // Disabled by default for terminal UI
		activeAnimations:  make(map[string]*ButtonAnimation),
		flashQueue:        make([]FlashEffect, 0),
		rippleEffects:     make([]RippleEffect, 0),
		feedbackHandlers:  make(map[string][]func(FeedbackEvent)),
	}
}

// WithPressDuration sets the button press animation duration
func (fm *FeedbackManager) WithPressDuration(duration time.Duration) *FeedbackManager {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.pressDuration = duration
	return fm
}

// WithTransitionSpeed sets the state transition animation speed
func (fm *FeedbackManager) WithTransitionSpeed(speed time.Duration) *FeedbackManager {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.transitionSpeed = speed
	return fm
}

// WithFocusAnimation enables or disables focus animations
func (fm *FeedbackManager) WithFocusAnimation(enabled bool) *FeedbackManager {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.focusAnimation = enabled
	return fm
}

// WithFlash enables or disables flash effects
func (fm *FeedbackManager) WithFlash(enabled bool) *FeedbackManager {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.flashEnabled = enabled
	return fm
}

// WithRipple enables or disables ripple effects
func (fm *FeedbackManager) WithRipple(enabled bool) *FeedbackManager {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.rippleEnabled = enabled
	return fm
}

// TriggerPressAnimation triggers a button press animation
func (fm *FeedbackManager) TriggerPressAnimation(button *Button) error {
	if button == nil {
		return fmt.Errorf("cannot animate nil button")
	}

	fm.mu.Lock()
	defer fm.mu.Unlock()

	animation := &ButtonAnimation{
		Button:     button,
		Type:       AnimPress,
		StartTime:  time.Now(),
		Duration:   fm.pressDuration,
		Progress:   0.0,
		Completed:  false,
		Properties: map[string]interface{}{
			"intensity": 1.0,
		},
	}

	key := fm.getAnimationKey(button, AnimPress)
	fm.activeAnimations[key] = animation

	// Trigger feedback event
	fm.triggerFeedbackEvent(FeedbackEvent{
		Type:      "animation_started",
		Button:    button,
		Animation: animation,
		Timestamp: time.Now(),
	})

	return nil
}

// TriggerFocusAnimation triggers a focus/blur animation
func (fm *FeedbackManager) TriggerFocusAnimation(button *Button, focused bool) error {
	if button == nil || !fm.focusAnimation {
		return nil
	}

	fm.mu.Lock()
	defer fm.mu.Unlock()

	animType := AnimFocus
	if !focused {
		animType = AnimRelease
	}

	animation := &ButtonAnimation{
		Button:     button,
		Type:       animType,
		StartTime:  time.Now(),
		Duration:   fm.transitionSpeed,
		Progress:   0.0,
		Completed:  false,
		Properties: map[string]interface{}{
			"focused": focused,
		},
	}

	key := fm.getAnimationKey(button, animType)
	fm.activeAnimations[key] = animation

	// Trigger feedback event
	fm.triggerFeedbackEvent(FeedbackEvent{
		Type:      "animation_started",
		Button:    button,
		Animation: animation,
		Timestamp: time.Now(),
	})

	return nil
}

// TriggerFlashEffect triggers a visual flash on a button
func (fm *FeedbackManager) TriggerFlashEffect(button *Button, color lipgloss.Color) error {
	if button == nil || !fm.flashEnabled {
		return nil
	}

	fm.mu.Lock()
	defer fm.mu.Unlock()

	flash := FlashEffect{
		Button:    button,
		Color:     color,
		StartTime: time.Now(),
		Duration:  fm.flashDuration,
	}

	fm.flashQueue = append(fm.flashQueue, flash)

	// Trigger feedback event
	fm.triggerFeedbackEvent(FeedbackEvent{
		Type:      "flash_triggered",
		Button:    button,
		Timestamp: time.Now(),
		Details:   map[string]interface{}{"color": color},
	})

	return nil
}

// TriggerRippleEffect triggers a ripple animation from a click point
func (fm *FeedbackManager) TriggerRippleEffect(button *Button, centerX, centerY int) error {
	if button == nil || !fm.rippleEnabled {
		return nil
	}

	fm.mu.Lock()
	defer fm.mu.Unlock()

	ripple := RippleEffect{
		Button:    button,
		CenterX:   centerX,
		CenterY:   centerY,
		StartTime: time.Now(),
		Duration:  500 * time.Millisecond,
		Radius:    0.0,
		MaxRadius:  30, // pixels
	}

	fm.rippleEffects = append(fm.rippleEffects, ripple)

	// Trigger feedback event
	fm.triggerFeedbackEvent(FeedbackEvent{
		Type:      "ripple_triggered",
		Button:    button,
		Timestamp: time.Now(),
		Details:   map[string]interface{}{"center_x": centerX, "center_y": centerY},
	})

	return nil
}

// Update updates all active animations and effects
func (fm *FeedbackManager) Update() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	now := time.Now()

	// Update button animations
	for _, anim := range fm.activeAnimations {
		if anim.Completed {
			continue
		}

		elapsed := now.Sub(anim.StartTime)
		anim.Progress = float64(elapsed) / float64(anim.Duration)

		if anim.Progress >= 1.0 {
			anim.Progress = 1.0
			anim.Completed = true

			// Trigger completion event
			fm.triggerFeedbackEvent(FeedbackEvent{
				Type:      "animation_completed",
				Button:    anim.Button,
				Animation: anim,
				Timestamp: now,
			})
		}

		// Apply animation effects
		fm.applyAnimationFrame(anim)
	}

	// Remove completed animations
	for animKey, anim := range fm.activeAnimations {
		if anim.Completed {
			delete(fm.activeAnimations, animKey)
			_ = anim // use anim variable
		}
	}

	// Update flash effects
	fm.flashQueue = fm.updateFlashEffects(fm.flashQueue, now)

	// Update ripple effects
	fm.rippleEffects = fm.updateRippleEffects(fm.rippleEffects, now)
}

// applyAnimationFrame applies the current frame of an animation
func (fm *FeedbackManager) applyAnimationFrame(anim *ButtonAnimation) {
	if anim.Button == nil {
		return
	}

	switch anim.Type {
	case AnimPress:
		fm.applyPressEffect(anim)
	case AnimFocus, AnimRelease:
		fm.applyFocusEffect(anim)
	}
}

// applyPressEffect applies button press visual feedback
func (fm *FeedbackManager) applyPressEffect(anim *ButtonAnimation) {
	button := anim.Button
	intensity := anim.Properties["intensity"].(float64)

	// Calculate visual intensity based on animation progress
	_ = intensity // use intensity variable
	// For press animation, we want to peak at 50% and then return
	var visualIntensity float64
	if anim.Progress < 0.5 {
		visualIntensity = anim.Progress * 2.0 // 0.0 -> 1.0
	} else {
		visualIntensity = 2.0 - (anim.Progress * 2.0) // 1.0 -> 0.0
	}

	// Apply visual intensity through button state
	// The actual visual changes are handled by the button's render method
	// This just ensures the button is in the pressed state during the animation
	if visualIntensity > 0.1 && !button.IsPressed() {
		button.Press()
	} else if visualIntensity <= 0.1 && button.IsPressed() {
		button.Release()
	}
}

// applyFocusEffect applies focus/blur visual feedback
func (fm *FeedbackManager) applyFocusEffect(anim *ButtonAnimation) {
	button := anim.Button
	focused := anim.Properties["focused"].(bool)

	// Smooth transition animation
	// The button state is already set, this is for additional visual effects
	if focused && !button.IsFocused() {
		button.Focus()
	} else if !focused && button.IsFocused() {
		button.Blur()
	}
}

// updateFlashEffects updates and removes expired flash effects
func (fm *FeedbackManager) updateFlashEffects(effects []FlashEffect, now time.Time) []FlashEffect {
	var activeEffects []FlashEffect

	for _, flash := range effects {
		elapsed := now.Sub(flash.StartTime)
		if elapsed < flash.Duration {
			activeEffects = append(activeEffects, flash)
		}
	}

	return activeEffects
}

// updateRippleEffects updates and removes expired ripple effects
func (fm *FeedbackManager) updateRippleEffects(effects []RippleEffect, now time.Time) []RippleEffect {
	var activeEffects []RippleEffect

	for _, ripple := range effects {
		elapsed := now.Sub(ripple.StartTime)
		if elapsed < ripple.Duration {
			// Update ripple radius
			progress := float64(elapsed) / float64(ripple.Duration)
			ripple.Radius = float64(ripple.MaxRadius) * progress
			activeEffects = append(activeEffects, ripple)
		}
	}

	return activeEffects
}

// getAnimationKey generates a unique key for an animation
func (fm *FeedbackManager) getAnimationKey(button *Button, animType AnimationType) string {
	if button == nil {
		return fmt.Sprintf("unknown_%d", animType)
	}
	return fmt.Sprintf("%s_%d_%p", button.GetLabel(), animType, button)
}

// triggerFeedbackEvent triggers feedback system events
func (fm *FeedbackManager) triggerFeedbackEvent(event FeedbackEvent) {
	handlers := fm.feedbackHandlers[event.Type]
	for _, handler := range handlers {
		handler(event)
	}

	// Also trigger "all" handlers
	allHandlers := fm.feedbackHandlers["all"]
	for _, handler := range allHandlers {
		handler(event)
	}
}

// RegisterFeedbackHandler registers a callback for feedback events
func (fm *FeedbackManager) RegisterFeedbackHandler(eventType string, handler func(FeedbackEvent)) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if fm.feedbackHandlers == nil {
		fm.feedbackHandlers = make(map[string][]func(FeedbackEvent))
	}
	fm.feedbackHandlers[eventType] = append(fm.feedbackHandlers[eventType], handler)
}

// GetActiveAnimations returns all currently active animations
func (fm *FeedbackManager) GetActiveAnimations() []*ButtonAnimation {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	animations := make([]*ButtonAnimation, 0, len(fm.activeAnimations))
	for _, anim := range fm.activeAnimations {
		animations = append(animations, anim)
	}
	return animations
}

// GetActiveFlashEffects returns all active flash effects
func (fm *FeedbackManager) GetActiveFlashEffects() []FlashEffect {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	return fm.flashQueue
}

// GetActiveRippleEffects returns all active ripple effects
func (fm *FeedbackManager) GetActiveRippleEffects() []RippleEffect {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	return fm.rippleEffects
}

// HasActiveAnimations returns true if there are any active animations
func (fm *FeedbackManager) HasActiveAnimations() bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	return len(fm.activeAnimations) > 0 || len(fm.flashQueue) > 0 || len(fm.rippleEffects) > 0
}

// CancelAnimation cancels a specific animation
func (fm *FeedbackManager) CancelAnimation(button *Button, animType AnimationType) {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	key := fm.getAnimationKey(button, animType)
	if anim, exists := fm.activeAnimations[key]; exists {
		anim.Completed = true

		// Trigger cancellation event
		fm.triggerFeedbackEvent(FeedbackEvent{
			Type:      "animation_cancelled",
			Button:    button,
			Animation: anim,
			Timestamp: time.Now(),
		})

		delete(fm.activeAnimations, key)
	}
}

// CancelAllAnimations cancels all active animations
func (fm *FeedbackManager) CancelAllAnimations() {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	for _, anim := range fm.activeAnimations {
		anim.Completed = true

		// Trigger cancellation event
		fm.triggerFeedbackEvent(FeedbackEvent{
			Type:      "animation_cancelled",
			Button:    anim.Button,
			Animation: anim,
			Timestamp: time.Now(),
		})
	}

	fm.activeAnimations = make(map[string]*ButtonAnimation)
	fm.flashQueue = make([]FlashEffect, 0)
	fm.rippleEffects = make([]RippleEffect, 0)
}

// GetAnimationProgress returns the progress of a specific animation
func (fm *FeedbackManager) GetAnimationProgress(button *Button, animType AnimationType) float64 {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	key := fm.getAnimationKey(button, animType)
	if anim, exists := fm.activeAnimations[key]; exists {
		return anim.Progress
	}
	return 0.0
}

// IsAnimationActive returns true if a specific animation is active
func (fm *FeedbackManager) IsAnimationActive(button *Button, animType AnimationType) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	key := fm.getAnimationKey(button, animType)
	_, exists := fm.activeAnimations[key]
	return exists
}

// EnhancedButtonRenderer extends the button renderer with feedback animations
type EnhancedButtonRenderer struct {
	feedbackManager *FeedbackManager
	baseRenderer    *ButtonRenderer
}

// NewEnhancedButtonRenderer creates an enhanced button renderer with feedback
func NewEnhancedButtonRenderer(feedbackManager *FeedbackManager, baseRenderer *ButtonRenderer) *EnhancedButtonRenderer {
	return &EnhancedButtonRenderer{
		feedbackManager: feedbackManager,
		baseRenderer:    baseRenderer,
	}
}

// RenderWithFeedback renders a button with active feedback effects
func (ebr *EnhancedButtonRenderer) RenderWithFeedback(button *Button) string {
	if button == nil {
		return ""
	}

	// Get base rendering
	baseRender := ebr.baseRenderer.Render(button)

	// Apply feedback effects if any
	if ebr.feedbackManager.HasActiveAnimations() {
		return ebr.applyFeedbackEffects(button, baseRender)
	}

	return baseRender
}

// applyFeedbackEffects applies visual feedback effects to button rendering
func (ebr *EnhancedButtonRenderer) applyFeedbackEffects(button *Button, baseRender string) string {
	// Check for active animations on this button
	animTypes := []AnimationType{AnimPress, AnimFocus, AnimRelease, AnimFlash}

	for _, animType := range animTypes {
		if ebr.feedbackManager.IsAnimationActive(button, animType) {
			progress := ebr.feedbackManager.GetAnimationProgress(button, animType)

			switch animType {
			case AnimPress:
				return ebr.applyPressFeedback(button, baseRender, progress)
			case AnimFocus, AnimRelease:
				return ebr.applyFocusFeedback(button, baseRender, progress, animType == AnimFocus)
			case AnimFlash:
				return ebr.applyFlashFeedback(button, baseRender, progress)
			}
		}
	}

	return baseRender
}

// applyPressFeedback applies press animation visual effects
func (ebr *EnhancedButtonRenderer) applyPressFeedback(button *Button, baseRender string, progress float64) string {
	// Calculate intensity (0.0 to 1.0 to 0.0)
	intensity := progress
	if progress > 0.5 {
		intensity = 2.0 - (progress * 2.0)
	}

	// Apply intensity-based styling modifications
	if intensity > 0.5 {
		// Add visual emphasis during peak press
		style := lipgloss.NewStyle().Bold(true)
		return style.Render(baseRender)
	}

	return baseRender
}

// applyFocusFeedback applies focus/blur animation visual effects
func (ebr *EnhancedButtonRenderer) applyFocusFeedback(button *Button, baseRender string, progress float64, gainingFocus bool) string {
	// Apply subtle animation during focus transitions
	if gainingFocus {
		// Focusing in - slightly brighten
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
		return style.Render(baseRender)
	} else {
		// Focusing out - slightly dim
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
		return style.Render(baseRender)
	}
}

// applyFlashFeedback applies flash animation visual effects
func (ebr *EnhancedButtonRenderer) applyFlashFeedback(button *Button, baseRender string, progress float64) string {
	// Flash effect: bright white overlay that fades
	intensity := 1.0 - progress

	if intensity > 0.1 {
		style := lipgloss.NewStyle().
			Background(lipgloss.Color("15")).
			Foreground(lipgloss.Color("0"))
		return style.Render(baseRender)
	}

	return baseRender
}