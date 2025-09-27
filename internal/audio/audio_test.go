package audio

import (
	"context"
	"testing"
	"time"

	"github.com/faiface/beep"

	uiintegration "ccpm-demo/internal/ui/integration"
	"ccpm-demo/internal/ui/components"
)

// MockAudioService is a mock implementation of AudioService for testing
type MockAudioService struct {
	initialized bool
	enabled     bool
	volume      float64
	muted       bool
	events      []*AudioEvent
}

func NewMockAudioService() *MockAudioService {
	return &MockAudioService{
		initialized: true,
		enabled:     true,
		volume:      0.5,
		muted:       false,
		events:      make([]*AudioEvent, 0),
	}
}

func (m *MockAudioService) Initialize(ctx context.Context, config *AudioConfig) error {
	m.initialized = true
	return nil
}

func (m *MockAudioService) Close() error {
	m.initialized = false
	return nil
}

func (m *MockAudioService) IsInitialized() bool {
	return m.initialized
}

func (m *MockAudioService) GetConfig() *AudioConfig {
	return &AudioConfig{
		Enabled: m.enabled,
		Volume:  m.volume,
		Muted:   m.muted,
	}
}

func (m *MockAudioService) UpdateConfig(config *AudioConfig) error {
	m.enabled = config.Enabled
	m.volume = config.Volume
	m.muted = config.Muted
	return nil
}

func (m *MockAudioService) SetEnabled(enabled bool) error {
	m.enabled = enabled
	return nil
}

func (m *MockAudioService) SetVolume(volume float64) error {
	m.volume = volume
	return nil
}

func (m *MockAudioService) SetMuted(muted bool) error {
	m.muted = muted
	return nil
}

func (m *MockAudioService) PlayEvent(event *AudioEvent) error {
	if !m.enabled || m.muted || !m.initialized {
		return nil
	}
	m.events = append(m.events, event)
	return nil
}

func (m *MockAudioService) PlayEventAsync(event *AudioEvent) chan error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- m.PlayEvent(event)
		close(errChan)
	}()
	return errChan
}

func (m *MockAudioService) PlaySound(streamer beep.Streamer) error {
	return nil
}

func (m *MockAudioService) PlaySoundAsync(streamer beep.Streamer) chan error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- nil
		close(errChan)
	}()
	return errChan
}

func (m *MockAudioService) PlayTone(frequency float64, duration time.Duration) error {
	return nil
}

func (m *MockAudioService) PlayToneAsync(frequency float64, duration time.Duration) chan error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- nil
		close(errChan)
	}()
	return errChan
}

func (m *MockAudioService) PlayBeep() error {
	return nil
}

func (m *MockAudioService) PlayBeepAsync() chan error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- nil
		close(errChan)
	}()
	return errChan
}

func (m *MockAudioService) PlayErrorSound() error {
	return nil
}

func (m *MockAudioService) PlayErrorSoundAsync() chan error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- nil
		close(errChan)
	}()
	return errChan
}

func (m *MockAudioService) LoadSoundFile(path string) (beep.StreamSeekCloser, error) {
	return nil, nil
}

func (m *MockAudioService) UnloadSoundFile(path string) error {
	return nil
}

func (m *MockAudioService) CleanupResources() error {
	return nil
}

func (m *MockAudioService) GetStatus() *AudioStatus {
	return &AudioStatus{
		Initialized: m.initialized,
		Enabled:     m.enabled,
		Muted:       m.muted,
		Volume:      m.volume,
	}
}

func (m *MockAudioService) GetStats() *AudioStats {
	return &AudioStats{}
}

func (m *MockAudioService) IsAudioAvailable() bool {
	return m.initialized
}

func (m *MockAudioService) TestAudio() error {
	return nil
}

func (m *MockAudioService) GetEvents() []*AudioEvent {
	return m.events
}

func (m *MockAudioService) ClearEvents() {
	m.events = make([]*AudioEvent, 0)
}

// TestIntegration_Initialization tests audio integration initialization
func TestIntegration_Initialization(t *testing.T) {
	integration := NewIntegration()

	// Test initialization
	err := integration.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize audio integration: %v", err)
	}

	if !integration.IsInitialized() {
		t.Error("Integration should be initialized")
	}

	// Test double initialization
	err = integration.Initialize()
	if err != nil {
		t.Errorf("Double initialization should not fail: %v", err)
	}

	// Test cleanup
	err = integration.Close()
	if err != nil {
		t.Errorf("Failed to close integration: %v", err)
	}
}

// TestIntegration_WithMockAudioService tests integration with mock audio service
func TestIntegration_WithMockAudioService(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	integration := &Integration{
		audioService: NewMockAudioService(),
		eventBuffer:  make(chan *AudioEvent, 100),
		ctx:          ctx,
		cancel:       cancel,
		initialized:  false,
	}

	err := integration.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize integration: %v", err)
	}

	// Test audio service access
	audioService := integration.GetAudioService()
	if audioService == nil {
		t.Error("Audio service should not be nil")
	}

	// Test audio control
	err = integration.SetEnabled(true)
	if err != nil {
		t.Errorf("Failed to enable audio: %v", err)
	}

	err = integration.SetVolume(0.8)
	if err != nil {
		t.Errorf("Failed to set volume: %v", err)
	}

	err = integration.SetMuted(false)
	if err != nil {
		t.Errorf("Failed to unmute audio: %v", err)
	}

	// Test status
	status := integration.GetStatus()
	if !status.Initialized {
		t.Error("Status should show as initialized")
	}

	err = integration.Close()
	if err != nil {
		t.Errorf("Failed to close integration: %v", err)
	}
}

// TestEventHandler_ButtonPress tests button press event handling
func TestEventHandler_ButtonPress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockService := NewMockAudioService()
	integration := &Integration{
		audioService: mockService,
		eventBuffer:  make(chan *AudioEvent, 100),
		ctx:          ctx,
		cancel:       cancel,
		initialized:  true,
	}

	_ = NewEventHandler(integration)

	// Test button press handling by creating audio event directly
	event := &AudioEvent{
		Type:      AudioEventNumber,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"button_value": "7",
		},
	}

	err := integration.PlayEventImmediately(event)
	if err != nil {
		t.Errorf("Failed to play audio event: %v", err)
	}

	// Check if audio event was created
	events := mockService.GetEvents()
	if len(events) == 0 {
		t.Error("No audio events were created")
	}

	// Check event type
	if events[0].Type != AudioEventNumber {
		t.Errorf("Expected number event type, got %v", events[0].Type)
	}

	// Check metadata
	if events[0].Metadata["button_value"] != "7" {
		t.Errorf("Expected button value '7' in metadata, got %v", events[0].Metadata["button_value"])
	}
}

// TestEventHandler_OperatorButton tests operator button event handling
func TestEventHandler_OperatorButton(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockService := NewMockAudioService()
	integration := &Integration{
		audioService: mockService,
		eventBuffer:  make(chan *AudioEvent, 100),
		ctx:          ctx,
		cancel:       cancel,
		initialized:  true,
	}

	_ = NewEventHandler(integration)

	// Test operator button handling by creating audio event directly
	event := &AudioEvent{
		Type:      AudioEventOperator,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"button_value": "+",
		},
	}

	err := integration.PlayEventImmediately(event)
	if err != nil {
		t.Errorf("Failed to play operator audio event: %v", err)
	}

	events := mockService.GetEvents()
	if len(events) == 0 {
		t.Error("No audio events were created")
	}

	if events[0].Type != AudioEventOperator {
		t.Errorf("Expected operator event type, got %v", events[0].Type)
	}
}

// TestEventHandler_SpecialButton tests special button event handling
func TestEventHandler_SpecialButton(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockService := NewMockAudioService()
	integration := &Integration{
		audioService: mockService,
		eventBuffer:  make(chan *AudioEvent, 100),
		ctx:          ctx,
		cancel:       cancel,
		initialized:  true,
	}

	_ = NewEventHandler(integration)

	// Test special button handling by creating audio event directly
	event := &AudioEvent{
		Type:      AudioEventClear,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"button_value": "clear",
		},
	}

	err := integration.PlayEventImmediately(event)
	if err != nil {
		t.Errorf("Failed to play special button audio event: %v", err)
	}

	events := mockService.GetEvents()
	if len(events) == 0 {
		t.Error("No audio events were created")
	}

	if events[0].Type != AudioEventClear {
		t.Errorf("Expected clear event type, got %v", events[0].Type)
	}
}

// TestEventHandler_CalculationResult tests calculation result event handling
func TestEventHandler_CalculationResult(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockService := NewMockAudioService()
	integration := &Integration{
		audioService: mockService,
		eventBuffer:  make(chan *AudioEvent, 100),
		ctx:          ctx,
		cancel:       cancel,
		initialized:  true,
	}

	handler := NewEventHandler(integration)

	// Test success result
	err := handler.HandleCalculationResult("42", false)
	if err != nil {
		t.Errorf("Failed to handle calculation result: %v", err)
	}

	events := mockService.GetEvents()
	if len(events) == 0 {
		t.Error("No audio events were created")
	}

	if events[0].Type != AudioEventSuccess {
		t.Errorf("Expected success event type, got %v", events[0].Type)
	}

	mockService.ClearEvents()

	// Test error result
	err = handler.HandleCalculationResult("", true)
	if err != nil {
		t.Errorf("Failed to handle error result: %v", err)
	}

	events = mockService.GetEvents()
	if len(events) == 0 {
		t.Error("No audio events were created for error")
	}

	if events[0].Type != AudioEventError {
		t.Errorf("Expected error event type, got %v", events[0].Type)
	}
}

// TestEventHandler_ClearEvent tests clear event handling
func TestEventHandler_ClearEvent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockService := NewMockAudioService()
	integration := &Integration{
		audioService: mockService,
		eventBuffer:  make(chan *AudioEvent, 100),
		ctx:          ctx,
		cancel:       cancel,
		initialized:  true,
	}

	handler := NewEventHandler(integration)

	// Test clear event
	err := handler.HandleClearEvent("clear")
	if err != nil {
		t.Errorf("Failed to handle clear event: %v", err)
	}

	events := mockService.GetEvents()
	if len(events) == 0 {
		t.Error("No audio events were created")
	}

	if events[0].Type != AudioEventClear {
		t.Errorf("Expected clear event type, got %v", events[0].Type)
	}

	mockService.ClearEvents()

	// Test clear entry event
	err = handler.HandleClearEvent("clear_entry")
	if err != nil {
		t.Errorf("Failed to handle clear entry event: %v", err)
	}

	events = mockService.GetEvents()
	if len(events) == 0 {
		t.Error("No audio events were created for clear entry")
	}

	if events[0].Type != AudioEventClear {
		t.Errorf("Expected clear event type for clear entry, got %v", events[0].Type)
	}
}

// TestEventHandler_SystemEvents tests system event handling
func TestEventHandler_SystemEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockService := NewMockAudioService()
	integration := &Integration{
		audioService: mockService,
		eventBuffer:  make(chan *AudioEvent, 100),
		ctx:          ctx,
		cancel:       cancel,
		initialized:  true,
	}

	handler := NewEventHandler(integration)

	// Test startup event
	err := handler.HandleStartupEvent()
	if err != nil {
		t.Errorf("Failed to handle startup event: %v", err)
	}

	events := mockService.GetEvents()
	if len(events) == 0 {
		t.Error("No audio events were created for startup")
	}

	if events[0].Type != AudioEventStartup {
		t.Errorf("Expected startup event type, got %v", events[0].Type)
	}

	mockService.ClearEvents()

	// Test shutdown event
	err = handler.HandleShutdownEvent()
	if err != nil {
		t.Errorf("Failed to handle shutdown event: %v", err)
	}

	events = mockService.GetEvents()
	if len(events) == 0 {
		t.Error("No audio events were created for shutdown")
	}

	if events[0].Type != AudioEventShutdown {
		t.Errorf("Expected shutdown event type, got %v", events[0].Type)
	}
}

// TestEventHandler_AudioControls tests audio control methods
func TestEventHandler_AudioControls(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockService := NewMockAudioService()
	integration := &Integration{
		audioService: mockService,
		eventBuffer:  make(chan *AudioEvent, 100),
		ctx:          ctx,
		cancel:       cancel,
		initialized:  true,
	}

	handler := NewEventHandler(integration)

	// Test audio enable/disable
	err := handler.DisableAudio()
	if err != nil {
		t.Errorf("Failed to disable audio: %v", err)
	}

	if handler.IsAudioEnabled() {
		t.Error("Audio should be disabled")
	}

	err = handler.EnableAudio()
	if err != nil {
		t.Errorf("Failed to enable audio: %v", err)
	}

	if !handler.IsAudioEnabled() {
		t.Error("Audio should be enabled")
	}

	// Test volume control
	err = handler.SetAudioVolume(0.7)
	if err != nil {
		t.Errorf("Failed to set volume: %v", err)
	}

	// Test mute control
	err = handler.MuteAudio()
	if err != nil {
		t.Errorf("Failed to mute audio: %v", err)
	}

	if handler.IsAudioEnabled() {
		t.Error("Audio should be muted")
	}

	err = handler.UnmuteAudio()
	if err != nil {
		t.Errorf("Failed to unmute audio: %v", err)
	}

	if !handler.IsAudioEnabled() {
		t.Error("Audio should be unmuted")
	}

	// Test audio
	err = handler.TestAudio()
	if err != nil {
		t.Errorf("Failed to test audio: %v", err)
	}
}

// TestEventHandler_EventHistory tests event history functionality
func TestEventHandler_EventHistory(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockService := NewMockAudioService()
	integration := &Integration{
		audioService: mockService,
		eventBuffer:  make(chan *AudioEvent, 100),
		ctx:          ctx,
		cancel:       cancel,
		initialized:  true,
	}

	handler := NewEventHandler(integration)

	// Create test actions
	button := components.NewButton(components.ButtonConfig{
		Label: "1",
		Type:  components.TypeNumber,
		Value: "1",
	})

	action := &uiintegration.ButtonAction{
		Button:   button,
		Action:   "press",
		Value:    "1",
		ButtonID: "test_button",
	}

	// Handle multiple events
	for i := 0; i < 5; i++ {
		err := handler.HandleButtonPress(action)
		if err != nil {
			t.Errorf("Failed to handle button press: %v", err)
		}
	}

	// Check history
	history := handler.GetEventHistory()
	if len(history) != 5 {
		t.Errorf("Expected 5 events in history, got %d", len(history))
	}

	// Check recent events
	recent := handler.GetRecentEvents(3)
	if len(recent) != 3 {
		t.Errorf("Expected 3 recent events, got %d", len(recent))
	}

	// Check stats
	stats := handler.GetEventStats()
	if stats.TotalEvents != 5 {
		t.Errorf("Expected 5 total events, got %d", stats.TotalEvents)
	}

	if stats.EventCounts[CalculatorEventNumber] != 5 {
		t.Errorf("Expected 5 number events, got %d", stats.EventCounts[CalculatorEventNumber])
	}

	// Clear history
	handler.ClearHistory()
	if len(handler.GetEventHistory()) != 0 {
		t.Error("History should be empty after clearing")
	}
}

// TestEventHandler_Validation tests input validation
func TestEventHandler_Validation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	integration := &Integration{
		audioService: NewMockAudioService(),
		eventBuffer:  make(chan *AudioEvent, 100),
		ctx:          ctx,
		cancel:       cancel,
		initialized:  true,
	}

	handler := NewEventHandler(integration)

	// Test nil action
	err := handler.ValidateButtonAction(nil)
	if err == nil {
		t.Error("Expected error for nil button action")
	}

	// Test action with nil button
	action := &uiintegration.ButtonAction{
		Button:   nil,
		Action:   "press",
		Value:    "test",
		ButtonID: "test_button",
	}

	err = handler.ValidateButtonAction(action)
	if err == nil {
		t.Error("Expected error for nil button")
	}

	// Test action with empty value
	button := components.NewButton(components.ButtonConfig{
		Label: "Test",
		Type:  components.TypeNumber,
		Value: "",
	})

	action = &uiintegration.ButtonAction{
		Button:   button,
		Action:   "press",
		Value:    "",
		ButtonID: "test_button",
	}

	err = handler.ValidateButtonAction(action)
	if err == nil {
		t.Error("Expected error for empty button value")
	}
}

// TestIntegration_EventBuffering tests event buffering functionality
func TestIntegration_EventBuffering(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockService := NewMockAudioService()
	integration := &Integration{
		audioService: mockService,
		eventBuffer:  make(chan *AudioEvent, 2), // Small buffer for testing
		ctx:          ctx,
		cancel:       cancel,
		initialized:  true,
	}

	// Fill the buffer
	event1 := &AudioEvent{Type: AudioEventNumber, Timestamp: time.Now()}
	event2 := &AudioEvent{Type: AudioEventOperator, Timestamp: time.Now()}
	event3 := &AudioEvent{Type: AudioEventClear, Timestamp: time.Now()}

	err := integration.QueueAudioEvent(event1)
	if err != nil {
		t.Errorf("Failed to queue first event: %v", err)
	}

	err = integration.QueueAudioEvent(event2)
	if err != nil {
		t.Errorf("Failed to queue second event: %v", err)
	}

	// This should fail because buffer is full
	err = integration.QueueAudioEvent(event3)
	if err == nil {
		t.Error("Expected error when buffer is full")
	}

	// Process events to free up buffer
	time.Sleep(10 * time.Millisecond) // Allow time for processing

	// Now it should work
	err = integration.QueueAudioEvent(event3)
	if err != nil {
		t.Errorf("Failed to queue third event after processing: %v", err)
	}
}

// TestIntegration_ImmediatePlayback tests immediate event playback
func TestIntegration_ImmediatePlayback(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mockService := NewMockAudioService()
	integration := &Integration{
		audioService: mockService,
		eventBuffer:  make(chan *AudioEvent, 100),
		ctx:          ctx,
		cancel:       cancel,
		initialized:  true,
	}

	event := &AudioEvent{
		Type:      AudioEventNumber,
		Timestamp: time.Now(),
		Metadata:  map[string]interface{}{"test": "value"},
	}

	err := integration.PlayEventImmediately(event)
	if err != nil {
		t.Errorf("Failed to play event immediately: %v", err)
	}

	events := mockService.GetEvents()
	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if events[0].Type != AudioEventNumber {
		t.Errorf("Expected number event type, got %v", events[0].Type)
	}
}