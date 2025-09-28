package input

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"ccpm-demo/internal/ui"
)

// EventRouter manages the routing and dispatch of all input events
type EventRouter struct {
	keyHandler    *KeyboardHandler
	mouseHandler  *MouseHandler
	validators    []*InputValidator

	// Event processing state
	enabled       bool
	eventQueue    []Event
	processEvents bool
}

// Event represents a unified input event that can be processed by the router
type Event struct {
	Type      EventType
	Source    EventSource
	Data      interface{}
	Timestamp int64
	Priority  EventPriority
}

// EventType represents the type of event
type EventType int

const (
	EventTypeKey EventType = iota
	EventTypeMouse
	EventTypeFocus
	EventTypeSystem
)

// EventSource represents the source of the event
type EventSource int

const (
	EventSourceKeyboard EventSource = iota
	EventSourceMouse
	EventSourceFocus
	EventSourceSystem
)

// EventPriority represents the priority of event processing
type EventPriority int

const (
	PriorityLow EventPriority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// EventValidator defines the interface for input validation
type EventValidator interface {
	// Validate validates an input event and returns true if valid
	Validate(event Event) bool

	// GetValidationError returns an error message if validation fails
	GetValidationError() string
}

// EventHandler defines the interface for handling processed events
type EventHandler interface {
	// HandleEvent processes a validated event and returns updated model and commands
	HandleEvent(model ui.Model, event Event) (ui.Model, tea.Cmd)

	// CanHandle returns true if this handler can process the given event type
	CanHandle(eventType EventType) bool
}

// NewEventRouter creates a new event router with default handlers
func NewEventRouter() *EventRouter {
	return &EventRouter{
		keyHandler:    NewKeyboardHandler(),
		mouseHandler:  NewMouseHandler(),
		validators:    []*InputValidator{},
		enabled:       true,
		eventQueue:    make([]Event, 0),
		processEvents: true,
	}
}

// AddValidator adds an input validator to the router
func (er *EventRouter) AddValidator(validator *InputValidator) {
	er.validators = append(er.validators, validator)
}

// RemoveValidator removes an input validator from the router
func (er *EventRouter) RemoveValidator(validator *InputValidator) {
	for i, v := range er.validators {
		if v == validator {
			er.validators = append(er.validators[:i], er.validators[i+1:]...)
			break
		}
	}
}

// SetEnabled enables or disables event routing
func (er *EventRouter) SetEnabled(enabled bool) {
	er.enabled = enabled
	if !enabled {
		er.eventQueue = make([]Event, 0)
	}
	er.mouseHandler.SetEnabled(enabled)
}

// ProcessMessage processes a tea.Msg and routes it to appropriate handlers
func (er *EventRouter) ProcessMessage(model ui.Model, msg tea.Msg) (ui.Model, tea.Cmd) {
	if !er.enabled || !er.processEvents {
		return model, nil
	}

	var command tea.Cmd
	var events []Event

	// Convert tea.Msg to internal events
	switch m := msg.(type) {
	case tea.KeyMsg:
		events = er.processKeyMessage(m)
	case tea.MouseMsg:
		events = er.processMouseMessage(m)
	default:
		// Handle other message types as system events
		events = er.processSystemMessage(msg)
	}

	// Process each event
	for _, event := range events {
		// Validate the event
		if !er.validateEvent(event) {
			continue
		}

		// Route the event to appropriate handlers
		updatedModel, cmd := er.routeEvent(model, event)
		model = updatedModel

		if cmd != nil {
			command = cmd // Return the last command for simplicity
		}
	}

	return model, command
}

// processKeyMessage converts a tea.KeyMsg to internal key events
func (er *EventRouter) processKeyMessage(msg tea.KeyMsg) []Event {
	// Get the key action from the keyboard handler
	keyEvent := er.keyHandler.GetActionForKey(msg)

	event := Event{
		Type:      EventTypeKey,
		Source:    EventSourceKeyboard,
		Data:      keyEvent,
		Timestamp: time.Now().UnixNano(),
		Priority:  er.getEventPriority(keyEvent.Action),
	}

	return []Event{event}
}

// processMouseMessage converts a tea.MouseMsg to internal mouse events
func (er *EventRouter) processMouseMessage(msg tea.MouseMsg) []Event {
	// Process mouse message through mouse handler
	mouseEvents := er.mouseHandler.HandleMessage(msg)

	var events []Event
	for _, mouseEvent := range mouseEvents {
		event := Event{
			Type:      EventTypeMouse,
			Source:    EventSourceMouse,
			Data:      mouseEvent,
			Timestamp: time.Now().UnixNano(),
			Priority:  er.getMouseEventPriority(mouseEvent.(MouseEvent).Type),
		}
		events = append(events, event)
	}

	return events
}

// processSystemMessage handles system-level messages
func (er *EventRouter) processSystemMessage(msg tea.Msg) []Event {
	event := Event{
		Type:      EventTypeSystem,
		Source:    EventSourceSystem,
		Data:      msg,
		Timestamp: time.Now().UnixNano(),
		Priority:  PriorityNormal,
	}

	return []Event{event}
}

// validateEvent validates an event using all registered validators
func (er *EventRouter) validateEvent(event Event) bool {
	for _, validator := range er.validators {
		if !validator.Validate(event) {
			return false
		}
	}
	return true
}

// routeEvent routes a validated event to the appropriate handler
func (er *EventRouter) routeEvent(model ui.Model, event Event) (ui.Model, tea.Cmd) {
	switch event.Type {
	case EventTypeKey:
		keyEvent := event.Data.(KeyEvent)
		return er.keyHandler.HandleKey(model, tea.KeyMsg{
			Type:  keyEvent.Key,
			Runes: []rune{keyEvent.Rune},
			Alt:   keyEvent.Alt,
		})
	case EventTypeMouse:
		// Mouse events are handled by the mouse handler and converted to commands
		// The actual model updates happen when the commands are executed
		return model, er.handleMouseEvent(event.Data.(MouseEvent))
	case EventTypeSystem:
		return er.handleSystemEvent(model, event)
	default:
		return model, nil
	}
}

// handleMouseEvent processes mouse events and returns appropriate commands
func (er *EventRouter) handleMouseEvent(mouseEvent MouseEvent) tea.Cmd {
	// Convert mouse events to calculator operations
	switch mouseEvent.Action.Type {
	case "number":
		return er.createNumberCommand(mouseEvent.Action.Value)
	case "operator":
		return er.createOperatorCommand(mouseEvent.Action.Value)
	case "equals":
		return er.createEqualsCommand()
	case "clear":
		return er.createClearCommand()
	case "backspace":
		return er.createBackspaceCommand()
	default:
		return nil
	}
}

// handleSystemEvent processes system events
func (er *EventRouter) handleSystemEvent(model ui.Model, event Event) (ui.Model, tea.Cmd) {
	switch event.Data.(type) {
	case tea.WindowSizeMsg:
		// Handle window resize
		return model, nil
	case tea.QuitMsg:
		// Handle quit message
		return model, tea.Quit
	default:
		return model, nil
	}
}

// getEventPriority returns the priority for a given key action
func (er *EventRouter) getEventPriority(action KeyAction) EventPriority {
	switch action {
	case KeyActionQuit:
		return PriorityCritical
	case KeyActionEquals, KeyActionClear:
		return PriorityHigh
	case KeyActionNumber, KeyActionOperator:
		return PriorityNormal
	default:
		return PriorityLow
	}
}

// getMouseEventPriority returns the priority for a given mouse event type
func (er *EventRouter) getMouseEventPriority(eventType MouseEventType) EventPriority {
	switch eventType {
	case MouseEventClick:
		return PriorityHigh
	case MouseEventDoubleClick:
		return PriorityHigh
	case MouseEventPress, MouseEventRelease:
		return PriorityNormal
	default:
		return PriorityLow
	}
}

// Command creation helpers
func (er *EventRouter) createNumberCommand(value string) tea.Cmd {
	return func() tea.Msg {
		return NumberInputMsg{Value: value}
	}
}

func (er *EventRouter) createOperatorCommand(operator string) tea.Cmd {
	return func() tea.Msg {
		return OperatorInputMsg{Operator: operator}
	}
}

func (er *EventRouter) createEqualsCommand() tea.Cmd {
	return func() tea.Msg {
		return EqualsInputMsg{}
	}
}

func (er *EventRouter) createClearCommand() tea.Cmd {
	return func() tea.Msg {
		return ClearInputMsg{}
	}
}

func (er *EventRouter) createBackspaceCommand() tea.Cmd {
	return func() tea.Msg {
		return BackspaceInputMsg{}
	}
}

// Message types for calculator operations
type NumberInputMsg struct {
	Value string
}

type OperatorInputMsg struct {
	Operator string
}

type EqualsInputMsg struct{}

type ClearInputMsg struct{}

type BackspaceInputMsg struct{}

// GetKeyHandler returns the keyboard handler
func (er *EventRouter) GetKeyHandler() *KeyboardHandler {
	return er.keyHandler
}

// GetMouseHandler returns the mouse handler
func (er *EventRouter) GetMouseHandler() *MouseHandler {
	return er.mouseHandler
}

// GetEventQueue returns the current event queue (for testing/debugging)
func (er *EventRouter) GetEventQueue() []Event {
	return er.eventQueue
}

// ClearEventQueue clears the event queue
func (er *EventRouter) ClearEventQueue() {
	er.eventQueue = make([]Event, 0)
}

// SetEventProcessing enables or disables event processing
func (er *EventRouter) SetEventProcessing(enabled bool) {
	er.processEvents = enabled
}