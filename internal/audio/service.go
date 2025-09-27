package audio

import (
	"context"
	"sync"
	"time"

	"github.com/faiface/beep"
)

// AudioEvent is imported from types.go
// AudioConfig is imported from types.go

// DefaultAudioConfig returns the default audio configuration
func DefaultAudioConfig() *AudioConfig {
	return &AudioConfig{
		Enabled:    true,
		Volume:     0.5,
		Muted:      false,
		SampleRate: 44100,
		BufferSize: 512,
		DeviceName: "",
		Profiles:   make(map[string]SoundProfile),
		Mappings:   []SoundProfileMapping{},
	}
}

// AudioService defines the interface for audio functionality
type AudioService interface {
	// Lifecycle management
	Initialize(ctx context.Context, config *AudioConfig) error
	Close() error
	IsInitialized() bool

	// Configuration
	GetConfig() *AudioConfig
	UpdateConfig(config *AudioConfig) error
	SetEnabled(enabled bool) error
	SetVolume(volume float64) error
	SetMuted(muted bool) error

	// Event playback
	PlayEvent(event *AudioEvent) error
	PlayEventAsync(event *AudioEvent) chan error

	// Sound playback
	PlaySound(streamer beep.Streamer) error
	PlaySoundAsync(streamer beep.Streamer) chan error
	PlayTone(frequency float64, duration time.Duration) error
	PlayToneAsync(frequency float64, duration time.Duration) chan error

	// Predefined sounds
	PlayBeep() error
	PlayBeepAsync() chan error
	PlayErrorSound() error
	PlayErrorSoundAsync() chan error

	// Resource management
	LoadSoundFile(path string) (beep.StreamSeekCloser, error)
	UnloadSoundFile(path string) error
	CleanupResources() error

	// Status and diagnostics
	GetStatus() *AudioStatus
	GetStats() *AudioStats
	IsAudioAvailable() bool
	TestAudio() error
}

// AudioStatus represents the current status of the audio service
type AudioStatus struct {
	Initialized   bool    `json:"initialized"`
	Available     bool    `json:"available"`
	Enabled       bool    `json:"enabled"`
	Muted         bool    `json:"muted"`
	Volume        float64 `json:"volume"`
	SampleRate    int     `json:"sampleRate"`
	Error         string  `json:"error,omitempty"`
	LastErrorTime time.Time `json:"lastErrorTime,omitempty"`
}

// AudioStats represents audio service statistics
type AudioStats struct {
	EventsPlayed      int64         `json:"eventsPlayed"`
	SoundsPlayed      int64         `json:"soundsPlayed"`
	ErrorsOccurred    int64         `json:"errorsOccurred"`
	LastEventTime     time.Time     `json:"lastEventTime"`
	LastError         error         `json:"lastError,omitempty"`
	Uptime            time.Duration `json:"uptime"`
	AveragePlayTime   time.Duration `json:"averagePlayTime"`
	MaxPlayTime       time.Duration `json:"maxPlayTime"`
	MinPlayTime       time.Duration `json:"minPlayTime"`
}

// audioServiceImpl implements the AudioService interface
type audioServiceImpl struct {
	config      *AudioConfig
	audioCtx    *AudioContext
	errorHandler *ErrorHandler
	mu          sync.RWMutex
	stats       *AudioStats
	startTime   time.Time
	closed      bool
}

// NewAudioService creates a new audio service instance
func NewAudioService() AudioService {
	return &audioServiceImpl{
		config:      DefaultAudioConfig(),
		audioCtx:    NewAudioContext(),
		errorHandler: DefaultErrorHandler(),
		stats:       &AudioStats{},
		startTime:   time.Now(),
	}
}

// Initialize initializes the audio service
func (s *audioServiceImpl) Initialize(ctx context.Context, config *AudioConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return NewAudioError(ErrContextAlreadyClosed, "audio service is closed")
	}

	// Update configuration if provided
	if config != nil {
		s.config = config
	}

	// Initialize audio context
	err := s.audioCtx.Initialize()
	if err != nil {
		s.stats.ErrorsOccurred++
		s.stats.LastError = err
		return s.errorHandler.HandleError(err)
	}

	s.stats.Uptime = time.Since(s.startTime)
	return nil
}

// Close closes the audio service and releases resources
func (s *audioServiceImpl) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	err := s.audioCtx.Close()
	if err != nil {
		s.errorHandler.HandleError(err)
	}

	s.closed = true
	return nil
}

// IsInitialized checks if the audio service is initialized
func (s *audioServiceImpl) IsInitialized() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.audioCtx.GetBeepIntegration().IsInitialized() && !s.closed
}

// GetConfig returns the current audio configuration
func (s *audioServiceImpl) GetConfig() *AudioConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

// UpdateConfig updates the audio configuration
func (s *audioServiceImpl) UpdateConfig(config *AudioConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return NewAudioError(ErrContextAlreadyClosed, "audio service is closed")
	}

	s.config = config
	return nil
}

// SetEnabled enables or disables audio
func (s *audioServiceImpl) SetEnabled(enabled bool) error {
	s.mu.Lock()
	defer s.mu.RUnlock()

	if s.closed {
		return NewAudioError(ErrContextAlreadyClosed, "audio service is closed")
	}

	s.config.Enabled = enabled
	return nil
}

// SetVolume sets the audio volume
func (s *audioServiceImpl) SetVolume(volume float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return NewAudioError(ErrContextAlreadyClosed, "audio service is closed")
	}

	if volume < 0 || volume > 1 {
		return NewAudioError(ErrInvalidConfig, "volume must be between 0.0 and 1.0")
	}

	s.config.Volume = volume
	return nil
}

// SetMuted mutes or unmutes audio
func (s *audioServiceImpl) SetMuted(muted bool) error {
	s.mu.Lock()
	defer s.mu.RUnlock()

	if s.closed {
		return NewAudioError(ErrContextAlreadyClosed, "audio service is closed")
	}

	s.config.Muted = muted
	return nil
}

// PlayEvent plays an audio event
func (s *audioServiceImpl) PlayEvent(event *AudioEvent) error {
	if !s.shouldPlayAudio() {
		return nil
	}

	startTime := time.Now()
	defer func() {
		s.updatePlayStats(time.Since(startTime))
	}()

	switch event.Type {
	case AudioEventNumber, AudioEventDecimal:
		return s.PlayTone(600, 50*time.Millisecond)
	case AudioEventOperator:
		return s.PlayTone(800, 75*time.Millisecond)
	case AudioEventEquals:
		return s.PlayTone(1000, 100*time.Millisecond)
	case AudioEventClear, AudioEventClearEntry, AudioEventBackspace:
		return s.PlayTone(300, 100*time.Millisecond)
	case AudioEventSignToggle, AudioEventPercent:
		return s.PlayTone(1200, 80*time.Millisecond)
	case AudioEventError:
		return s.PlayErrorSound()
	case AudioEventSuccess:
		return s.PlayTone(1500, 150*time.Millisecond)
	case AudioEventStartup:
		return s.PlayTone(800, 200*time.Millisecond)
	case AudioEventShutdown:
		return s.PlayTone(400, 300*time.Millisecond)
	default:
		return s.PlayBeep()
	}
}

// PlayEventAsync plays an audio event asynchronously
func (s *audioServiceImpl) PlayEventAsync(event *AudioEvent) chan error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		errChan <- s.PlayEvent(event)
	}()

	return errChan
}

// PlaySound plays a sound streamer
func (s *audioServiceImpl) PlaySound(streamer beep.Streamer) error {
	if !s.shouldPlayAudio() {
		return nil
	}

	if streamer == nil {
		return NewAudioError(ErrInvalidResource, "nil streamer provided")
	}

	startTime := time.Now()
	defer func() {
		s.updatePlayStats(time.Since(startTime))
	}()

	// Apply volume if not muted
	var finalStreamer beep.Streamer = streamer
	if !s.config.Muted && s.config.Volume > 0 {
		finalStreamer = s.audioCtx.GetBeepIntegration().SetVolume(streamer, s.config.Volume)
	} else if s.config.Muted {
		// Use silence streamer instead of beep.Silence which requires Len()
		finalStreamer = &silenceStreamer{duration: time.Second}
	}

	err := s.audioCtx.PlaySound(finalStreamer)
	if err != nil {
		s.stats.ErrorsOccurred++
		s.stats.LastError = err
		return s.errorHandler.HandleError(err)
	}

	s.stats.SoundsPlayed++
	s.stats.LastEventTime = time.Now()
	return nil
}

// PlaySoundAsync plays a sound streamer asynchronously
func (s *audioServiceImpl) PlaySoundAsync(streamer beep.Streamer) chan error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		errChan <- s.PlaySound(streamer)
	}()

	return errChan
}

// PlayTone plays a tone
func (s *audioServiceImpl) PlayTone(frequency float64, duration time.Duration) error {
	if !s.shouldPlayAudio() {
		return nil
	}

	if duration <= 0 || duration > 5*time.Second {
		return NewAudioError(ErrInvalidConfig, "invalid duration")
	}

	startTime := time.Now()
	defer func() {
		s.updatePlayStats(time.Since(startTime))
	}()

	err := s.audioCtx.PlayTone(frequency, duration)
	if err != nil {
		s.stats.ErrorsOccurred++
		s.stats.LastError = err
		return s.errorHandler.HandleError(err)
	}

	s.stats.SoundsPlayed++
	s.stats.LastEventTime = time.Now()
	return nil
}

// PlayToneAsync plays a tone asynchronously
func (s *audioServiceImpl) PlayToneAsync(frequency float64, duration time.Duration) chan error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		errChan <- s.PlayTone(frequency, duration)
	}()

	return errChan
}

// PlayBeep plays a beep sound
func (s *audioServiceImpl) PlayBeep() error {
	return s.PlayTone(800, 100*time.Millisecond)
}

// PlayBeepAsync plays a beep sound asynchronously
func (s *audioServiceImpl) PlayBeepAsync() chan error {
	return s.PlayToneAsync(800, 100*time.Millisecond)
}

// PlayErrorSound plays an error sound
func (s *audioServiceImpl) PlayErrorSound() error {
	return s.PlayTone(200, 200*time.Millisecond)
}

// PlayErrorSoundAsync plays an error sound asynchronously
func (s *audioServiceImpl) PlayErrorSoundAsync() chan error {
	return s.PlayToneAsync(200, 200*time.Millisecond)
}

// LoadSoundFile loads a sound file
func (s *audioServiceImpl) LoadSoundFile(path string) (beep.StreamSeekCloser, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return nil, NewAudioError(ErrContextAlreadyClosed, "audio service is closed")
	}

	return s.audioCtx.GetBeepIntegration().LoadSoundFile(path)
}

// UnloadSoundFile unloads a sound file
func (s *audioServiceImpl) UnloadSoundFile(path string) error {
	// Currently no explicit unloading needed as Beep manages resources
	// This is a placeholder for future resource management
	return nil
}

// CleanupResources cleans up audio resources
func (s *audioServiceImpl) CleanupResources() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return NewAudioError(ErrContextAlreadyClosed, "audio service is closed")
	}

	// Force garbage collection of audio resources
	// This is mostly a no-op with Beep but included for completeness
	return nil
}

// GetStatus returns the current audio status
func (s *audioServiceImpl) GetStatus() *AudioStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status := &AudioStatus{
		Initialized: s.IsInitialized(),
		Available:   s.IsAudioAvailable(),
		Enabled:     s.config.Enabled,
		Muted:       s.config.Muted,
		Volume:      s.config.Volume,
		SampleRate:  s.config.SampleRate,
	}

	if s.stats.LastError != nil {
		status.Error = s.stats.LastError.Error()
		status.LastErrorTime = s.stats.LastEventTime
	}

	return status
}

// GetStats returns audio service statistics
func (s *audioServiceImpl) GetStats() *AudioStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Update uptime
	stats := *s.stats
	stats.Uptime = time.Since(s.startTime)
	return &stats
}

// IsAudioAvailable checks if audio is available
func (s *audioServiceImpl) IsAudioAvailable() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.audioCtx.GetBeepIntegration().IsInitialized()
}

// TestAudio tests the audio system
func (s *audioServiceImpl) TestAudio() error {
	if !s.shouldPlayAudio() {
		return nil
	}

	// Play a test tone
	err := s.PlayTone(440, 200*time.Millisecond) // A4 note
	if err != nil {
		return s.errorHandler.HandleError(err)
	}

	return nil
}

// Helper methods

// shouldPlayAudio determines if audio should be played
func (s *audioServiceImpl) shouldPlayAudio() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.config.Enabled && !s.closed && s.audioCtx.GetBeepIntegration().IsInitialized()
}

// updatePlayStats updates playback statistics
func (s *audioServiceImpl) updatePlayStats(playTime time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stats.EventsPlayed++
	s.stats.LastEventTime = time.Now()

	// Update play time statistics
	if s.stats.MinPlayTime == 0 || playTime < s.stats.MinPlayTime {
		s.stats.MinPlayTime = playTime
	}
	if playTime > s.stats.MaxPlayTime {
		s.stats.MaxPlayTime = playTime
	}

	// Update average play time
	if s.stats.EventsPlayed > 0 {
		totalPlayTime := s.stats.AveragePlayTime * time.Duration(s.stats.EventsPlayed-1)
		s.stats.AveragePlayTime = (totalPlayTime + playTime) / time.Duration(s.stats.EventsPlayed)
	} else {
		s.stats.AveragePlayTime = playTime
	}
}