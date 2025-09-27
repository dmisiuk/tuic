package audio

import (
	"fmt"
	"strings"
	"time"

	uiintegration "ccpm-demo/internal/ui/integration"
	"ccpm-demo/internal/ui/components"
)

// CalculatorEventType represents different types of calculator events
type CalculatorEventType int

const (
	CalculatorEventNumber CalculatorEventType = iota
	CalculatorEventOperator
	CalculatorEventEquals
	CalculatorEventClear
	CalculatorEventClearEntry
	CalculatorEventBackspace
	CalculatorEventError
	CalculatorEventSuccess
	CalculatorEventStartup
	CalculatorEventShutdown
)

// String returns the string representation of a calculator event type
func (et CalculatorEventType) String() string {
	switch et {
	case CalculatorEventNumber:
		return "number"
	case CalculatorEventOperator:
		return "operator"
	case CalculatorEventEquals:
		return "equals"
	case CalculatorEventClear:
		return "clear"
	case CalculatorEventClearEntry:
		return "clear_entry"
	case CalculatorEventBackspace:
		return "backspace"
	case CalculatorEventError:
		return "error"
	case CalculatorEventSuccess:
		return "success"
	case CalculatorEventStartup:
		return "startup"
	case CalculatorEventShutdown:
		return "shutdown"
	default:
		return "unknown"
	}
}

// CalculatorEvent represents a calculator event with audio mapping
type CalculatorEvent struct {
	Type      CalculatorEventType
	Timestamp time.Time
	Value     string
	Metadata  map[string]interface{}
}

// EventHandler handles calculator events and maps them to audio events
type EventHandler struct {
	integration *Integration
	eventHistory []CalculatorEvent
	maxHistory  int
}

// NewEventHandler creates a new calculator event handler
func NewEventHandler(integration *Integration) *EventHandler {
	return &EventHandler{
		integration:  integration,
		eventHistory: make([]CalculatorEvent, 0),
		maxHistory:   100, // Keep last 100 events
	}
}

// HandleButtonPress handles a button press event and triggers appropriate audio
func (eh *EventHandler) HandleButtonPress(action *uiintegration.ButtonAction) error {
	event := eh.createButtonPressEvent(action)

	// Add to history
	eh.addToHistory(event)

	// Map to audio event and play
	return eh.integration.HandleButtonAction(action)
}

// HandleCalculationResult handles the result of a calculation
func (eh *EventHandler) HandleCalculationResult(result string, isError bool) error {
	eventType := CalculatorEventSuccess
	if isError {
		eventType = CalculatorEventError
	}

	metadata := map[string]interface{}{
		"result": result,
		"is_error": isError,
	}

	return eh.integration.HandleCalculatorEvent(eventType, metadata)
}

// HandleClearEvent handles clear operations
func (eh *EventHandler) HandleClearEvent(clearType string) error {
	var eventType CalculatorEventType

	switch clearType {
	case "clear":
		eventType = CalculatorEventClear
	case "clear_entry":
		eventType = CalculatorEventClearEntry
	case "backspace":
		eventType = CalculatorEventBackspace
	default:
		eventType = CalculatorEventClear
	}

	metadata := map[string]interface{}{
		"clear_type": clearType,
	}

	return eh.integration.HandleCalculatorEvent(eventType, metadata)
}

// HandleStartupEvent handles application startup
func (eh *EventHandler) HandleStartupEvent() error {
	return eh.integration.HandleCalculatorEvent(CalculatorEventStartup, map[string]interface{}{
		"action": "startup",
	})
}

// HandleShutdownEvent handles application shutdown
func (eh *EventHandler) HandleShutdownEvent() error {
	return eh.integration.HandleCalculatorEvent(CalculatorEventShutdown, map[string]interface{}{
		"action": "shutdown",
	})
}

// createButtonPressEvent creates a calculator event from a button action
func (eh *EventHandler) createButtonPressEvent(action *uiintegration.ButtonAction) CalculatorEvent {
	button := action.Button
	if button == nil {
		return CalculatorEvent{
			Type:      CalculatorEventNumber,
			Timestamp: time.Now(),
			Value:     action.Value,
			Metadata:  map[string]interface{}{},
		}
	}

	buttonType := button.GetType()
	eventType := eh.mapButtonTypeToCalculatorEventType(buttonType, action.Value)

	return CalculatorEvent{
		Type:      eventType,
		Timestamp: time.Now(),
		Value:     action.Value,
		Metadata: map[string]interface{}{
			"button_id":    action.ButtonID,
			"button_label":  button.GetLabel(),
			"button_type":  buttonType.String(),
		},
	}
}

// mapButtonTypeToCalculatorEventType maps button types to calculator event types
func (eh *EventHandler) mapButtonTypeToCalculatorEventType(buttonType components.ButtonType, value string) CalculatorEventType {
	switch buttonType {
	case components.TypeNumber:
		return CalculatorEventNumber
	case components.TypeOperator:
		return CalculatorEventOperator
	case components.TypeSpecial:
		return eh.mapSpecialButtonToEvent(value)
	default:
		return CalculatorEventNumber
	}
}

// mapSpecialButtonToEvent maps special button values to event types
func (eh *EventHandler) mapSpecialButtonToEvent(value string) CalculatorEventType {
	switch value {
	case "=":
		return CalculatorEventEquals
	case "clear":
		return CalculatorEventClear
	case "clear_entry":
		return CalculatorEventClearEntry
	case "backspace":
		return CalculatorEventBackspace
	default:
		return CalculatorEventClear
	}
}

// addToHistory adds an event to the event history
func (eh *EventHandler) addToHistory(event CalculatorEvent) {
	eh.eventHistory = append(eh.eventHistory, event)

	// Trim history if it exceeds max size
	if len(eh.eventHistory) > eh.maxHistory {
		eh.eventHistory = eh.eventHistory[1:]
	}
}

// GetEventHistory returns the event history
func (eh *EventHandler) GetEventHistory() []CalculatorEvent {
	return eh.eventHistory
}

// GetRecentEvents returns the most recent events
func (eh *EventHandler) GetRecentEvents(count int) []CalculatorEvent {
	if count <= 0 || count > len(eh.eventHistory) {
		return eh.eventHistory
	}

	start := len(eh.eventHistory) - count
	return eh.eventHistory[start:]
}

// ClearHistory clears the event history
func (eh *EventHandler) ClearHistory() {
	eh.eventHistory = make([]CalculatorEvent, 0)
}

// GetEventStats returns statistics about handled events
func (eh *EventHandler) GetEventStats() *EventStats {
	stats := &EventStats{
		TotalEvents: len(eh.eventHistory),
		EventCounts:  make(map[CalculatorEventType]int),
	}

	for _, event := range eh.eventHistory {
		stats.EventCounts[event.Type]++
	}

	return stats
}

// EventStats represents statistics about calculator events
type EventStats struct {
	TotalEvents int                         `json:"totalEvents"`
	EventCounts map[CalculatorEventType]int `json:"eventCounts"`
}

// String returns a string representation of event stats
func (es *EventStats) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("EventStats{Total: %d, Counts: {", es.TotalEvents))

	first := true
	for eventType, count := range es.EventCounts {
		if !first {
			builder.WriteString(", ")
		}
		builder.WriteString(fmt.Sprintf("%s: %d", eventType.String(), count))
		first = false
	}

	builder.WriteString("}}")
	return builder.String()
}

// ValidateButtonAction validates a button action
func (eh *EventHandler) ValidateButtonAction(action *uiintegration.ButtonAction) error {
	if action == nil {
		return NewAudioError(ErrInvalidResource, "button action is nil")
	}

	if action.Button == nil {
		return NewAudioError(ErrInvalidResource, "button is nil")
	}

	if action.Value == "" {
		return NewAudioError(ErrInvalidResource, "button value is empty")
	}

	return nil
}

// IsAudioEnabled checks if audio is enabled for the integration
func (eh *EventHandler) IsAudioEnabled() bool {
	if eh.integration == nil {
		return false
	}

	status := eh.integration.GetStatus()
	return status.Initialized && status.AudioStatus.Enabled && !status.AudioStatus.Muted
}

// EnableAudio enables audio for the integration
func (eh *EventHandler) EnableAudio() error {
	if eh.integration == nil {
		return NewAudioError(ErrNotInitialized, "integration is not initialized")
	}

	return eh.integration.SetEnabled(true)
}

// DisableAudio disables audio for the integration
func (eh *EventHandler) DisableAudio() error {
	if eh.integration == nil {
		return NewAudioError(ErrNotInitialized, "integration is not initialized")
	}

	return eh.integration.SetEnabled(false)
}

// SetAudioVolume sets the audio volume
func (eh *EventHandler) SetAudioVolume(volume float64) error {
	if eh.integration == nil {
		return NewAudioError(ErrNotInitialized, "integration is not initialized")
	}

	return eh.integration.SetVolume(volume)
}

// MuteAudio mutes the audio
func (eh *EventHandler) MuteAudio() error {
	if eh.integration == nil {
		return NewAudioError(ErrNotInitialized, "integration is not initialized")
	}

	return eh.integration.SetMuted(true)
}

// UnmuteAudio unmutes the audio
func (eh *EventHandler) UnmuteAudio() error {
	if eh.integration == nil {
		return NewAudioError(ErrNotInitialized, "integration is not initialized")
	}

	return eh.integration.SetMuted(false)
}

// TestAudio tests the audio system
func (eh *EventHandler) TestAudio() error {
	if eh.integration == nil {
		return NewAudioError(ErrNotInitialized, "integration is not initialized")
	}

	return eh.integration.TestAudio()
}