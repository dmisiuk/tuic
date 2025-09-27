package audio

import (
	"math"
	"math/rand"
	"time"
)

// SoundGenerator handles sound generation for different profiles
type SoundGenerator struct {
	sampleRate int
	bufferSize int
}

// NewSoundGenerator creates a new sound generator
func NewSoundGenerator(sampleRate, bufferSize int) *SoundGenerator {
	return &SoundGenerator{
		sampleRate: sampleRate,
		bufferSize: bufferSize,
	}
}

// GenerateSound generates audio data for a given profile
func (sg *SoundGenerator) GenerateSound(profile SoundProfile) ([]float64, error) {
	if profile.Frequency <= 0 || profile.Duration <= 0 {
		return nil, NewAudioError(ErrInvalidResource, "invalid sound profile: frequency and duration must be positive")
	}

	samples := int(float64(sg.sampleRate) * profile.Duration.Seconds())
	data := make([]float64, samples)

	switch profile.WaveType {
	case WaveTypeSine:
		sg.generateSineWave(data, profile.Frequency, profile.Volume)
	case WaveTypeSquare:
		sg.generateSquareWave(data, profile.Frequency, profile.Volume)
	case WaveTypeSawtooth:
		sg.generateSawtoothWave(data, profile.Frequency, profile.Volume)
	case WaveTypeTriangle:
		sg.generateTriangleWave(data, profile.Frequency, profile.Volume)
	case WaveTypeNoise:
		sg.generateNoise(data, profile.Volume)
	default:
		return nil, &AudioError{
			Code:    ErrInvalidResource,
			Message: "unsupported wave type",
		}
	}

	return data, nil
}

// generateSineWave generates a sine wave
func (sg *SoundGenerator) generateSineWave(data []float64, frequency float64, volume float64) {
	for i := range data {
		t := float64(i) / float64(sg.sampleRate)
		data[i] = volume * math.Sin(2*math.Pi*frequency*t)
	}
}

// generateSquareWave generates a square wave
func (sg *SoundGenerator) generateSquareWave(data []float64, frequency float64, volume float64) {
	for i := range data {
		t := float64(i) / float64(sg.sampleRate)
		sample := math.Sin(2*math.Pi*frequency*t)
		if sample >= 0 {
			data[i] = volume
		} else {
			data[i] = -volume
		}
	}
}

// generateSawtoothWave generates a sawtooth wave
func (sg *SoundGenerator) generateSawtoothWave(data []float64, frequency float64, volume float64) {
	for i := range data {
		t := float64(i) / float64(sg.sampleRate)
		phase := math.Mod(t*frequency, 1.0)
		data[i] = volume * (2*phase - 1)
	}
}

// generateTriangleWave generates a triangle wave
func (sg *SoundGenerator) generateTriangleWave(data []float64, frequency float64, volume float64) {
	for i := range data {
		t := float64(i) / float64(sg.sampleRate)
		phase := math.Mod(t*frequency, 1.0)
		var sample float64
		if phase < 0.25 {
			sample = 4 * phase
		} else if phase < 0.75 {
			sample = 2 - 4*phase
		} else {
			sample = 4*phase - 4
		}
		data[i] = volume * sample
	}
}

// generateNoise generates white noise
func (sg *SoundGenerator) generateNoise(data []float64, volume float64) {
	for i := range data {
		// Simple white noise generation
		data[i] = volume * (2*rand.Float64() - 1)
	}
}

// DefaultSoundProfiles returns the default sound profiles for the calculator
func DefaultSoundProfiles() map[string]SoundProfile {
	return map[string]SoundProfile{
		ProfileNameNumber: {
			Name:        ProfileNameNumber,
			Description: "Number button press",
			Frequency:   800.0,
			Duration:    50 * time.Millisecond,
			Volume:      0.3,
			WaveType:    WaveTypeSine,
		},
		ProfileNameOperator: {
			Name:        ProfileNameOperator,
			Description: "Operator button press",
			Frequency:   600.0,
			Duration:    75 * time.Millisecond,
			Volume:      0.4,
			WaveType:    WaveTypeSquare,
		},
		ProfileNameEquals: {
			Name:        ProfileNameEquals,
			Description: "Equals button press",
			Frequency:   1000.0,
			Duration:    100 * time.Millisecond,
			Volume:      0.5,
			WaveType:    WaveTypeSine,
		},
		ProfileNameClear: {
			Name:        ProfileNameClear,
			Description: "Clear button press",
			Frequency:   400.0,
			Duration:    60 * time.Millisecond,
			Volume:      0.35,
			WaveType:    WaveTypeTriangle,
		},
		ProfileNameError: {
			Name:        ProfileNameError,
			Description: "Error notification",
			Frequency:   200.0,
			Duration:    200 * time.Millisecond,
			Volume:      0.6,
			WaveType:    WaveTypeSawtooth,
		},
		ProfileNameSuccess: {
			Name:        ProfileNameSuccess,
			Description: "Success notification",
			Frequency:   1200.0,
			Duration:    150 * time.Millisecond,
			Volume:      0.4,
			WaveType:    WaveTypeSine,
		},
		ProfileNameStartup: {
			Name:        ProfileNameStartup,
			Description: "Application startup",
			Frequency:   523.25, // C5
			Duration:    200 * time.Millisecond,
			Volume:      0.3,
			WaveType:    WaveTypeSine,
		},
		ProfileNameShutdown: {
			Name:        ProfileNameShutdown,
			Description: "Application shutdown",
			Frequency:   261.63, // C4
			Duration:    300 * time.Millisecond,
			Volume:      0.3,
			WaveType:    WaveTypeSine,
		},
	}
}

// DefaultProfileMappings returns the default button-to-profile mappings
func DefaultProfileMappings() []SoundProfileMapping {
	return []SoundProfileMapping{
		{
			ButtonType:  ButtonTypeNumber,
			ProfileName: ProfileNameNumber,
		},
		{
			ButtonType:  ButtonTypeOperator,
			ProfileName: ProfileNameOperator,
		},
		{
			ButtonType:  ButtonTypeSpecial,
			ProfileName: ProfileNameClear,
		},
		{
			ButtonType:  ButtonTypeSystem,
			ProfileName: ProfileNameError,
		},
	}
}

// CreateCustomProfile creates a custom sound profile with validation
func CreateCustomProfile(name, description string, frequency float64, duration time.Duration, volume float64, waveType WaveType) (*SoundProfile, error) {
	if name == "" {
		return nil, &AudioError{
			Code:    ErrInvalidResource,
			Message: "profile name cannot be empty",
		}
	}

	if frequency <= 0 {
		return nil, &AudioError{
			Code:    ErrInvalidResource,
			Message: "frequency must be positive",
		}
	}

	if duration <= 0 {
		return nil, &AudioError{
			Code:    ErrInvalidResource,
			Message: "duration must be positive",
		}
	}

	if volume < MinVolume || volume > MaxVolume {
		return nil, &AudioError{
			Code:    ErrInvalidResource,
			Message: "volume must be between 0.0 and 1.0",
		}
	}

	if waveType < WaveTypeSine || waveType > WaveTypeNoise {
		return nil, &AudioError{
			Code:    ErrInvalidResource,
			Message: "invalid wave type",
		}
	}

	return &SoundProfile{
		Name:        name,
		Description: description,
		Frequency:   frequency,
		Duration:    duration,
		Volume:      volume,
		WaveType:    waveType,
	}, nil
}

// ValidateProfile validates a sound profile
func ValidateProfile(profile SoundProfile) error {
	if profile.Name == "" {
		return &AudioError{
			Code:    ErrInvalidResource,
			Message: "profile name cannot be empty",
		}
	}

	if profile.Frequency <= 0 {
		return &AudioError{
			Code:    ErrInvalidResource,
			Message: "frequency must be positive",
		}
	}

	if profile.Duration <= 0 {
		return &AudioError{
			Code:    ErrInvalidResource,
			Message: "duration must be positive",
		}
	}

	if profile.Volume < MinVolume || profile.Volume > MaxVolume {
		return &AudioError{
			Code:    ErrInvalidResource,
			Message: "volume must be between 0.0 and 1.0",
		}
	}

	if profile.WaveType < WaveTypeSine || profile.WaveType > WaveTypeNoise {
		return &AudioError{
			Code:    ErrInvalidResource,
			Message: "invalid wave type",
		}
	}

	return nil
}

// Musical note frequencies for convenience
const (
	NoteC4  = 261.63
	NoteD4  = 293.66
	NoteE4  = 329.63
	NoteF4  = 349.23
	NoteG4  = 392.00
	NoteA4  = 440.00
	NoteB4  = 493.88
	NoteC5  = 523.25
	NoteD5  = 587.33
	NoteE5  = 659.25
	NoteF5  = 698.46
	NoteG5  = 783.99
	NoteA5  = 880.00
	NoteB5  = 987.77
	NoteC6  = 1046.50
)