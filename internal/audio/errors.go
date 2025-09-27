package audio

import (
	"errors"
	"fmt"
)

// AudioError represents an audio-related error
type AudioError struct {
	Code    AudioErrorCode
	Message string
	Cause   error
}

func (e *AudioError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("audio error [%s]: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("audio error [%s]: %s", e.Code, e.Message)
}

func (e *AudioError) Unwrap() error {
	return e.Cause
}

// AudioErrorCode represents different types of audio errors
type AudioErrorCode string

const (
	// Context errors
	ErrContextNotInitialized AudioErrorCode = "context_not_initialized"
	ErrContextAlreadyClosed  AudioErrorCode = "context_already_closed"
	ErrContextInitialization AudioErrorCode = "context_initialization"

	// Device errors
	ErrDeviceNotFound     AudioErrorCode = "device_not_found"
	ErrDeviceUnavailable  AudioErrorCode = "device_unavailable"
	ErrDevicePermission   AudioErrorCode = "device_permission"
	ErrDeviceInUse        AudioErrorCode = "device_in_use"

	// Format errors
	ErrUnsupportedFormat  AudioErrorCode = "unsupported_format"
	ErrInvalidFormat      AudioErrorCode = "invalid_format"
	ErrFormatConversion   AudioErrorCode = "format_conversion"

	// Playback errors
	ErrPlaybackFailed     AudioErrorCode = "playback_failed"
	ErrPlaybackInterrupted AudioErrorCode = "playback_interrupted"
	ErrBufferUnderrun     AudioErrorCode = "buffer_underrun"
	ErrStreamError        AudioErrorCode = "stream_error"

	// Resource errors
	ErrResourceNotFound   AudioErrorCode = "resource_not_found"
	ErrResourceBusy       AudioErrorCode = "resource_busy"
	ErrMemoryLimit        AudioErrorCode = "memory_limit"
	ErrInvalidResource    AudioErrorCode = "invalid_resource"

	// Configuration errors
	ErrInvalidConfig      AudioErrorCode = "invalid_config"
	ErrConfigNotFound     AudioErrorCode = "config_not_found"
	ErrConfigPermission   AudioErrorCode = "config_permission"

	// General errors
	ErrNotSupported       AudioErrorCode = "not_supported"
	ErrTimeout           AudioErrorCode = "timeout"
	ErrUnknown           AudioErrorCode = "unknown"
)

// NewAudioError creates a new AudioError
func NewAudioError(code AudioErrorCode, message string) *AudioError {
	return &AudioError{
		Code:    code,
		Message: message,
	}
}

// NewAudioErrorWithCause creates a new AudioError with a cause
func NewAudioErrorWithCause(code AudioErrorCode, message string, cause error) *AudioError {
	return &AudioError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// IsContextError checks if the error is a context-related error
func IsContextError(err error) bool {
	var audioErr *AudioError
	if errors.As(err, &audioErr) {
		switch audioErr.Code {
		case ErrContextNotInitialized, ErrContextAlreadyClosed, ErrContextInitialization:
			return true
		}
	}
	return false
}

// IsDeviceError checks if the error is a device-related error
func IsDeviceError(err error) bool {
	var audioErr *AudioError
	if errors.As(err, &audioErr) {
		switch audioErr.Code {
		case ErrDeviceNotFound, ErrDeviceUnavailable, ErrDevicePermission, ErrDeviceInUse:
			return true
		}
	}
	return false
}

// IsPlaybackError checks if the error is a playback-related error
func IsPlaybackError(err error) bool {
	var audioErr *AudioError
	if errors.As(err, &audioErr) {
		switch audioErr.Code {
		case ErrPlaybackFailed, ErrPlaybackInterrupted, ErrBufferUnderrun, ErrStreamError:
			return true
		}
	}
	return false
}

// IsRecoverableError checks if the error is recoverable (can be retried)
func IsRecoverableError(err error) bool {
	var audioErr *AudioError
	if errors.As(err, &audioErr) {
		switch audioErr.Code {
		case ErrDeviceInUse, ErrResourceBusy, ErrBufferUnderrun, ErrTimeout:
			return true
		}
	}
	return false
}

// IsFatalError checks if the error is fatal (cannot be recovered)
func IsFatalError(err error) bool {
	var audioErr *AudioError
	if errors.As(err, &audioErr) {
		switch audioErr.Code {
		case ErrContextNotInitialized, ErrContextAlreadyClosed, ErrDeviceNotFound,
			ErrDevicePermission, ErrUnsupportedFormat, ErrInvalidConfig:
			return true
		}
	}
	return false
}

// GracefulDegradationStrategy defines how to handle audio errors
type GracefulDegradationStrategy int

const (
	// FailFast - immediately return errors without fallback
	FailFast GracefulDegradationStrategy = iota
	// LogAndContinue - log errors but continue execution
	LogAndContinue
	// SilentContinue - ignore errors and continue silently
	SilentContinue
	// RetryWithBackoff - retry operation with exponential backoff
	RetryWithBackoff
)

// ErrorHandler handles audio errors according to the configured strategy
type ErrorHandler struct {
	Strategy GracefulDegradationStrategy
	Logger   func(string, ...interface{})
}

// NewErrorHandler creates a new ErrorHandler with the specified strategy
func NewErrorHandler(strategy GracefulDegradationStrategy) *ErrorHandler {
	return &ErrorHandler{
		Strategy: strategy,
		Logger:   func(format string, args ...interface{}) {}, // Default no-op logger
	}
}

// HandleError handles an error according to the configured strategy
func (h *ErrorHandler) HandleError(err error) error {
	if err == nil {
		return nil
	}

	switch h.Strategy {
	case FailFast:
		return err
	case LogAndContinue:
		h.Logger("audio error: %v", err)
		return nil
	case SilentContinue:
		return nil
	case RetryWithBackoff:
		if IsRecoverableError(err) {
			h.Logger("recoverable audio error: %v (will retry)", err)
			return err // Return error to allow retry logic
		}
		h.Logger("non-recoverable audio error: %v", err)
		return nil
	default:
		return err
	}
}

// WithLogger sets a custom logger for the error handler
func (h *ErrorHandler) WithLogger(logger func(string, ...interface{})) *ErrorHandler {
	h.Logger = logger
	return h
}

// DefaultErrorHandler returns a default error handler with graceful degradation
func DefaultErrorHandler() *ErrorHandler {
	return NewErrorHandler(LogAndContinue).
		WithLogger(func(format string, args ...interface{}) {
			// Default to standard library logging
			// Can be overridden by the application
		})
}