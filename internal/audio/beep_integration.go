package audio

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"os"
)

// BeepIntegration wraps the Beep library for cross-platform audio
type BeepIntegration struct {
	sampleRate beep.SampleRate
	mu         sync.RWMutex
	initialized bool
}

// NewBeepIntegration creates a new Beep integration wrapper
func NewBeepIntegration() *BeepIntegration {
	return &BeepIntegration{
		sampleRate: 44100, // Standard sample rate
	}
}

// Initialize initializes the Beep audio system
func (b *BeepIntegration) Initialize() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.initialized {
		return nil
	}

	// Initialize speaker with reasonable buffer size
	err := speaker.Init(b.sampleRate, b.sampleRate.N(time.Second/30)) // 30ms buffer
	if err != nil {
		return NewAudioErrorWithCause(ErrContextInitialization, "failed to initialize speaker", err)
	}

	b.initialized = true
	return nil
}

// IsInitialized checks if the audio system is initialized
func (b *BeepIntegration) IsInitialized() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.initialized
}

// Close shuts down the audio system
func (b *BeepIntegration) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.initialized {
		return NewAudioError(ErrContextNotInitialized, "audio system not initialized")
	}

	// Beep doesn't have a specific close method for speaker
	// We just mark it as uninitialized
	b.initialized = false
	return nil
}

// GetSampleRate returns the current sample rate
func (b *BeepIntegration) GetSampleRate() beep.SampleRate {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.sampleRate
}

// PlaySound plays a beep.Streamer asynchronously
func (b *BeepIntegration) PlaySound(streamer beep.Streamer) error {
	if !b.IsInitialized() {
		return NewAudioError(ErrContextNotInitialized, "audio system not initialized")
	}

	if streamer == nil {
		return NewAudioError(ErrInvalidResource, "nil streamer provided")
	}

	// Play the sound in a separate goroutine to avoid blocking
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Handle panic during playback
				err := fmt.Errorf("panic during audio playback: %v", r)
				DefaultErrorHandler().HandleError(err)
			}
		}()

		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			// Sound finished playing
			// Note: We don't close the streamer here as beep.Streamer doesn't have Close method
		})))
	}()

	return nil
}

// PlaySoundSync plays a beep.Streamer synchronously
func (b *BeepIntegration) PlaySoundSync(streamer beep.Streamer) error {
	if !b.IsInitialized() {
		return NewAudioError(ErrContextNotInitialized, "audio system not initialized")
	}

	if streamer == nil {
		return NewAudioError(ErrInvalidResource, "nil streamer provided")
	}

	done := make(chan struct{})
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		close(done)
	})))

	select {
	case <-done:
		return nil
	case <-time.After(30 * time.Second): // Timeout after 30 seconds
		return NewAudioError(ErrTimeout, "audio playback timed out")
	}
}

// PlayTone generates and plays a simple tone
func (b *BeepIntegration) PlayTone(frequency float64, duration time.Duration) error {
	if !b.IsInitialized() {
		return NewAudioError(ErrContextNotInitialized, "audio system not initialized")
	}

	oscillator := &toneOscillator{
		freq:  frequency,
		duration: duration,
		sampleRate: b.sampleRate,
		pos: 0,
	}

	return b.PlaySound(oscillator)
}

// PlayBeep plays a simple beep sound
func (b *BeepIntegration) PlayBeep() error {
	return b.PlayTone(800, 100*time.Millisecond)
}

// PlayErrorSound plays an error sound (lower frequency)
func (b *BeepIntegration) PlayErrorSound() error {
	return b.PlayTone(200, 200*time.Millisecond)
}

// LoadSoundFile loads a sound file from disk
func (b *BeepIntegration) LoadSoundFile(filepath string) (beep.StreamSeekCloser, error) {
	if !b.IsInitialized() {
		return nil, NewAudioError(ErrContextNotInitialized, "audio system not initialized")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, NewAudioErrorWithCause(ErrResourceNotFound, fmt.Sprintf("failed to open sound file: %s", filepath), err)
	}
	defer file.Close()

	// Try to decode as WAV first
	streamer, format, err := wav.Decode(file)
	if err == nil {
		// Resample if needed
		if format.SampleRate != b.sampleRate {
			resampled := beep.Resample(4, format.SampleRate, b.sampleRate, streamer)
			return &resampledSeeker{streamer: resampled}, nil
		}
		return streamer, nil
	}

	// Reset file position and try MP3
	if _, err := file.Seek(0, 0); err != nil {
		return nil, NewAudioErrorWithCause(ErrResourceNotFound, fmt.Sprintf("failed to seek in sound file: %s", filepath), err)
	}

	streamer, format, err = mp3.Decode(file)
	if err != nil {
		return nil, NewAudioErrorWithCause(ErrUnsupportedFormat, fmt.Sprintf("failed to decode sound file: %s", filepath), err)
	}

	// Resample if needed
	if format.SampleRate != b.sampleRate {
		resampled := beep.Resample(4, format.SampleRate, b.sampleRate, streamer)
		return &resampledSeeker{streamer: resampled}, nil
	}

	return streamer, nil
}

// SetVolume adjusts the volume of a streamer
func (b *BeepIntegration) SetVolume(streamer beep.Streamer, volume float64) beep.Streamer {
	if volume <= 0 {
		// Return silence instead of using streamer.Len() which doesn't exist
		return &silenceStreamer{duration: time.Second}
	}
	if volume >= 1.0 {
		return streamer
	}
	// For volume control, we need to create a wrapper that applies volume
	return &volumeStreamer{streamer: streamer, volume: volume}
}

// toneOscillator generates a simple sine wave tone
type toneOscillator struct {
	freq       float64
	duration   time.Duration
	sampleRate beep.SampleRate
	pos        int
}

func (t *toneOscillator) Stream(samples [][2]float64) (n int, ok bool) {
	if t.pos >= int(t.sampleRate.N(t.duration)) {
		return 0, false
	}

	for i := range samples {
		if t.pos >= int(t.sampleRate.N(t.duration)) {
			return i, true
		}

		// Generate sine wave
		sample := math.Sin(2 * math.Pi * float64(t.pos) * t.freq / float64(t.sampleRate))
		samples[i][0] = sample * 0.3 // Left channel, reduce volume
		samples[i][1] = sample * 0.3 // Right channel, reduce volume
		t.pos++
	}

	return len(samples), true
}

func (t *toneOscillator) Err() error {
	return nil
}

func (t *toneOscillator) Len() int {
	return int(t.sampleRate.N(t.duration))
}

func (t *toneOscillator) Position() int {
	return t.pos
}

func (t *toneOscillator) Seek(p int) error {
	if p < 0 || p > t.Len() {
		return NewAudioError(ErrInvalidResource, "invalid seek position")
	}
	t.pos = p
	return nil
}

func (t *toneOscillator) Close() error {
	t.pos = t.Len()
	return nil
}

// silenceStreamer generates silence
type silenceStreamer struct {
	duration time.Duration
	pos      int
}

func (s *silenceStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		samples[i][0] = 0
		samples[i][1] = 0
		s.pos++
	}
	return len(samples), true
}

func (s *silenceStreamer) Err() error {
	return nil
}

func (s *silenceStreamer) Len() int {
	return int(44100 * s.duration.Seconds())
}

func (s *silenceStreamer) Position() int {
	return s.pos
}

func (s *silenceStreamer) Seek(p int) error {
	if p < 0 || p > s.Len() {
		return NewAudioError(ErrInvalidResource, "invalid seek position")
	}
	s.pos = p
	return nil
}

func (s *silenceStreamer) Close() error {
	return nil
}

// volumeStreamer applies volume control to another streamer
type volumeStreamer struct {
	streamer beep.Streamer
	volume   float64
}

func (v *volumeStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = v.streamer.Stream(samples)
	for i := 0; i < n; i++ {
		samples[i][0] *= v.volume
		samples[i][1] *= v.volume
	}
	return n, ok
}

func (v *volumeStreamer) Err() error {
	return v.streamer.Err()
}

// resampledSeeker wraps a Resampler to implement StreamSeekCloser
type resampledSeeker struct {
	streamer *beep.Resampler
}

func (r *resampledSeeker) Stream(samples [][2]float64) (n int, ok bool) {
	return r.streamer.Stream(samples)
}

func (r *resampledSeeker) Err() error {
	return r.streamer.Err()
}

func (r *resampledSeeker) Len() int {
	// Resampler doesn't expose Len, return an estimate
	return 44100 // 1 second of audio as default
}

func (r *resampledSeeker) Position() int {
	// Resampler doesn't expose position, return 0
	return 0
}

func (r *resampledSeeker) Seek(p int) error {
	// Resampler doesn't support seeking
	return NewAudioError(ErrNotSupported, "seeking not supported on resampled audio")
}

func (r *resampledSeeker) Close() error {
	return nil
}

// AudioContext manages the audio context and lifecycle
type AudioContext struct {
	beep       *BeepIntegration
	errorHandler *ErrorHandler
	mu         sync.RWMutex
	closed     bool
}

// NewAudioContext creates a new audio context
func NewAudioContext() *AudioContext {
	return &AudioContext{
		beep:        NewBeepIntegration(),
		errorHandler: DefaultErrorHandler(),
	}
}

// Initialize initializes the audio context
func (ctx *AudioContext) Initialize() error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if ctx.closed {
		return NewAudioError(ErrContextAlreadyClosed, "audio context is closed")
	}

	err := ctx.beep.Initialize()
	if err != nil {
		return ctx.errorHandler.HandleError(err)
	}

	return nil
}

// Close closes the audio context and releases resources
func (ctx *AudioContext) Close() error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if ctx.closed {
		return nil
	}

	err := ctx.beep.Close()
	if err != nil {
		ctx.errorHandler.HandleError(err)
	}

	ctx.closed = true
	return nil
}

// IsClosed checks if the context is closed
func (ctx *AudioContext) IsClosed() bool {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.closed
}

// GetBeepIntegration returns the underlying Beep integration
func (ctx *AudioContext) GetBeepIntegration() *BeepIntegration {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()
	return ctx.beep
}

// SetErrorHandler sets a custom error handler
func (ctx *AudioContext) SetErrorHandler(handler *ErrorHandler) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.errorHandler = handler
}

// PlaySound plays a sound using the audio context
func (ctx *AudioContext) PlaySound(streamer beep.Streamer) error {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	if ctx.closed {
		return NewAudioError(ErrContextAlreadyClosed, "audio context is closed")
	}

	err := ctx.beep.PlaySound(streamer)
	return ctx.errorHandler.HandleError(err)
}

// PlayTone plays a tone using the audio context
func (ctx *AudioContext) PlayTone(frequency float64, duration time.Duration) error {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	if ctx.closed {
		return NewAudioError(ErrContextAlreadyClosed, "audio context is closed")
	}

	err := ctx.beep.PlayTone(frequency, duration)
	return ctx.errorHandler.HandleError(err)
}

// PlayBeep plays a beep sound using the audio context
func (ctx *AudioContext) PlayBeep() error {
	return ctx.PlayTone(800, 100*time.Millisecond)
}

// PlayErrorSound plays an error sound using the audio context
func (ctx *AudioContext) PlayErrorSound() error {
	return ctx.PlayTone(200, 200*time.Millisecond)
}