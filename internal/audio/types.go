package audio

import "time"

// AudioEventType represents different types of audio events
type AudioEventType int

const (
	// Number button events
	AudioEventNumber AudioEventType = iota
	AudioEventDecimal

	// Operator button events
	AudioEventOperator
	AudioEventEquals

	// Special button events
	AudioEventClear
	AudioEventClearEntry
	AudioEventBackspace
	AudioEventSignToggle
	AudioEventPercent

	// System events
	AudioEventError
	AudioEventSuccess
	AudioEventStartup
	AudioEventShutdown
)

// AudioEvent represents an audio event that should trigger sound playback
type AudioEvent struct {
	Type      AudioEventType
	Timestamp time.Time
	Metadata  map[string]interface{}
}

// SoundProfile defines the characteristics of a sound
type SoundProfile struct {
	Name        string
	Description string
	Frequency   float64 // Hz
	Duration    time.Duration
	Volume      float64 // 0.0 to 1.0
	WaveType    WaveType
}

// WaveType represents the type of sound wave
type WaveType int

const (
	WaveTypeSine WaveType = iota
	WaveTypeSquare
	WaveTypeSawtooth
	WaveTypeTriangle
	WaveTypeNoise
)

// ButtonType categorizes calculator buttons for sound mapping
type ButtonType int

const (
	ButtonTypeNumber ButtonType = iota
	ButtonTypeOperator
	ButtonTypeSpecial
	ButtonTypeSystem
)

// SoundProfileMapping maps button types to sound profiles
type SoundProfileMapping struct {
	ButtonType    ButtonType
	ProfileName   string
	CustomProfile *SoundProfile
}

// AudioConfig represents the audio configuration
type AudioConfig struct {
	Enabled      bool                   `json:"enabled"`
	Volume       float64                `json:"volume"`       // Master volume 0.0 to 1.0
	Muted        bool                   `json:"muted"`
	Profiles     map[string]SoundProfile `json:"profiles"`
	Mappings     []SoundProfileMapping  `json:"mappings"`
	DeviceName   string                 `json:"deviceName"`   // Specific audio device
	SampleRate   int                    `json:"sampleRate"`   // Audio sample rate
	BufferSize   int                    `json:"bufferSize"`   // Audio buffer size
}

// Default configuration values
const (
	DefaultVolume     = 0.7
	DefaultSampleRate = 44100
	DefaultBufferSize = 512
	MinVolume         = 0.0
	MaxVolume         = 1.0
)

// Sound profile presets
const (
	ProfileNameNumber      = "number"
	ProfileNameOperator    = "operator"
	ProfileNameEquals      = "equals"
	ProfileNameClear       = "clear"
	ProfileNameError       = "error"
	ProfileNameSuccess     = "success"
	ProfileNameStartup     = "startup"
	ProfileNameShutdown    = "shutdown"
)

// Button type to event type mapping
var ButtonTypeToEventType = map[ButtonType]AudioEventType{
	ButtonTypeNumber:   AudioEventNumber,
	ButtonTypeOperator: AudioEventOperator,
	ButtonTypeSpecial:  AudioEventClear, // Default special button event
	ButtonTypeSystem:   AudioEventError, // Default system event
}

// Event type to profile name mapping
var EventTypeToProfileName = map[AudioEventType]string{
	AudioEventNumber:      ProfileNameNumber,
	AudioEventDecimal:     ProfileNameNumber,
	AudioEventOperator:    ProfileNameOperator,
	AudioEventEquals:      ProfileNameEquals,
	AudioEventClear:       ProfileNameClear,
	AudioEventClearEntry:  ProfileNameClear,
	AudioEventBackspace:   ProfileNameClear,
	AudioEventSignToggle:  ProfileNameOperator,
	AudioEventPercent:     ProfileNameOperator,
	AudioEventError:       ProfileNameError,
	AudioEventSuccess:     ProfileNameSuccess,
	AudioEventStartup:     ProfileNameStartup,
	AudioEventShutdown:    ProfileNameShutdown,
}