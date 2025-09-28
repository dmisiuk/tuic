package audio

import (
	"context"
	"sync"
	"time"

	uiintegration "ccpm-demo/internal/ui/integration"
	"ccpm-demo/internal/ui/components"
)

// Integration handles the connection between the audio service and UI events
type Integration struct {
	audioService AudioService
	eventBuffer  chan *AudioEvent
	errorHandler *ErrorHandler
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	initialized  bool
}

// NewIntegration creates a new audio integration instance
func NewIntegration() *Integration {
	ctx, cancel := context.WithCancel(context.Background())

	return &Integration{
		audioService: NewAudioService(),
		eventBuffer:  make(chan *AudioEvent, 100), // Buffer for audio events
		errorHandler: DefaultErrorHandler(),
		ctx:          ctx,
		cancel:       cancel,
		initialized:  false,
	}
}

// Initialize initializes the audio integration
func (ai *Integration) Initialize() error {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	if ai.initialized {
		return nil
	}

	// Initialize audio service with default configuration
	config := DefaultAudioConfig()

	// Add default sound profiles
	config.Profiles = DefaultSoundProfiles()

	// Add default profile mappings
	config.Mappings = DefaultProfileMappings()

	err := ai.audioService.Initialize(ai.ctx, config)
	if err != nil {
		return ai.errorHandler.HandleError(err)
	}

	// Start event processing goroutine
	go ai.processEvents()

	ai.initialized = true
	return nil
}

// Close shuts down the audio integration
func (ai *Integration) Close() error {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	if !ai.initialized {
		return nil
	}

	// Cancel context to stop event processing
	ai.cancel()

	// Close audio service
	err := ai.audioService.Close()
	if err != nil {
		ai.errorHandler.HandleError(err)
	}

	// Close event buffer
	close(ai.eventBuffer)

	ai.initialized = false
	return nil
}

// IsInitialized checks if the integration is initialized
func (ai *Integration) IsInitialized() bool {
	ai.mu.RLock()
	defer ai.mu.RUnlock()
	return ai.initialized
}

// GetAudioService returns the underlying audio service
func (ai *Integration) GetAudioService() AudioService {
	ai.mu.RLock()
	defer ai.mu.RUnlock()
	return ai.audioService
}

// HandleButtonAction handles a button action from the UI and triggers appropriate audio
func (ai *Integration) HandleButtonAction(action *uiintegration.ButtonAction) error {
	if !ai.IsInitialized() {
		return nil
	}

	// Map button action to audio event
	event, err := ai.mapButtonActionToAudioEvent(action)
	if err != nil {
		return ai.errorHandler.HandleError(err)
	}

	// Queue the audio event
	return ai.QueueAudioEvent(event)
}

// HandleCalculatorEvent handles calculator-specific events
func (ai *Integration) HandleCalculatorEvent(eventType CalculatorEventType, metadata map[string]interface{}) error {
	if !ai.IsInitialized() {
		return nil
	}

	// Map calculator event to audio event
	audioEvent := ai.mapCalculatorEventToAudioEvent(eventType, metadata)

	// Queue the audio event
	return ai.QueueAudioEvent(audioEvent)
}

// QueueAudioEvent queues an audio event for processing
func (ai *Integration) QueueAudioEvent(event *AudioEvent) error {
	if !ai.IsInitialized() {
		return nil
	}

	select {
	case ai.eventBuffer <- event:
		return nil
	default:
		// Buffer is full, drop the event
		return NewAudioError(ErrBufferFull, "audio event buffer is full")
	}
}

// PlayEventImmediately plays an audio event immediately, bypassing the buffer
func (ai *Integration) PlayEventImmediately(event *AudioEvent) error {
	if !ai.IsInitialized() {
		return nil
	}

	return ai.audioService.PlayEvent(event)
}

// processEvents processes audio events from the buffer
func (ai *Integration) processEvents() {
	for {
		select {
		case <-ai.ctx.Done():
			return
		case event, ok := <-ai.eventBuffer:
			if !ok {
				return
			}
			ai.playAudioEvent(event)
		}
	}
}

// playAudioEvent plays a single audio event
func (ai *Integration) playAudioEvent(event *AudioEvent) {
	// Use async playback to avoid blocking
	errChan := ai.audioService.PlayEventAsync(event)

	// Handle errors in a goroutine to avoid blocking
	go func() {
		if err := <-errChan; err != nil {
			ai.errorHandler.HandleError(err)
		}
	}()
}

// mapButtonActionToAudioEvent maps a button action to an audio event
func (ai *Integration) mapButtonActionToAudioEvent(action *uiintegration.ButtonAction) (*AudioEvent, error) {
	button := action.Button
	if button == nil {
		return nil, NewAudioError(ErrInvalidResource, "button is nil")
	}

	buttonType := button.GetType()
	var eventType AudioEventType

	switch buttonType {
	case components.TypeNumber:
		eventType = AudioEventNumber
	case components.TypeOperator:
		eventType = AudioEventOperator
	case components.TypeSpecial:
		eventType = ai.mapSpecialButtonValue(action.Value)
	default:
		eventType = AudioEventNumber // Default fallback
	}

	return &AudioEvent{
		Type:      eventType,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"button_id":   action.ButtonID,
			"button_value": action.Value,
			"button_label": button.GetLabel(),
		},
	}, nil
}

// mapSpecialButtonValue maps special button values to audio event types
func (ai *Integration) mapSpecialButtonValue(value string) AudioEventType {
	switch value {
	case "clear", "clear_entry":
		return AudioEventClear
	case "backspace":
		return AudioEventBackspace
	case "=":
		return AudioEventEquals
	default:
		return AudioEventClear // Default for special buttons
	}
}

// mapCalculatorEventToAudioEvent maps calculator events to audio events
func (ai *Integration) mapCalculatorEventToAudioEvent(eventType CalculatorEventType, metadata map[string]interface{}) *AudioEvent {
	var audioEventType AudioEventType

	switch eventType {
	case CalculatorEventError:
		audioEventType = AudioEventError
	case CalculatorEventSuccess:
		audioEventType = AudioEventSuccess
	case CalculatorEventClear:
		audioEventType = AudioEventClear
	case CalculatorEventStartup:
		audioEventType = AudioEventStartup
	case CalculatorEventShutdown:
		audioEventType = AudioEventShutdown
	default:
		audioEventType = AudioEventNumber // Default fallback
	}

	return &AudioEvent{
		Type:      audioEventType,
		Timestamp: time.Now(),
		Metadata:  metadata,
	}
}

// GetStatus returns the current integration status
func (ai *Integration) GetStatus() *IntegrationStatus {
	ai.mu.RLock()
	defer ai.mu.RUnlock()

	audioStatus := ai.audioService.GetStatus()
	bufferSize := len(ai.eventBuffer)

	return &IntegrationStatus{
		Initialized:   ai.initialized,
		AudioStatus:   *audioStatus,
		BufferSize:    bufferSize,
		BufferCapacity: cap(ai.eventBuffer),
	}
}

// SetEnabled enables or disables audio integration
func (ai *Integration) SetEnabled(enabled bool) error {
	return ai.audioService.SetEnabled(enabled)
}

// SetVolume sets the audio volume
func (ai *Integration) SetVolume(volume float64) error {
	return ai.audioService.SetVolume(volume)
}

// SetMuted mutes or unmutes audio
func (ai *Integration) SetMuted(muted bool) error {
	return ai.audioService.SetMuted(muted)
}

// TestAudio tests the audio integration
func (ai *Integration) TestAudio() error {
	if !ai.IsInitialized() {
		return NewAudioError(ErrNotInitialized, "audio integration is not initialized")
	}

	return ai.audioService.TestAudio()
}

// IntegrationStatus represents the current status of the audio integration
type IntegrationStatus struct {
	Initialized    bool          `json:"initialized"`
	AudioStatus    AudioStatus   `json:"audioStatus"`
	BufferSize     int           `json:"bufferSize"`
	BufferCapacity int           `json:"bufferCapacity"`
}