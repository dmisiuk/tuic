package audio

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// ConfigManager handles audio configuration persistence and management
type ConfigManager struct {
	config     *AudioConfig
	configPath string
	mu         sync.RWMutex
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(configPath string) *ConfigManager {
	cm := &ConfigManager{
		configPath: configPath,
	}
	cm.config = cm.createDefaultConfig()

	// Try to load existing configuration
	if err := cm.Load(); err != nil {
		// If loading fails, save the default configuration
		_ = cm.Save()
	}

	return cm
}

// createDefaultConfig creates a default audio configuration
func (cm *ConfigManager) createDefaultConfig() *AudioConfig {
	return &AudioConfig{
		Enabled:    true,
		Volume:     DefaultVolume,
		Muted:      false,
		DeviceName: "",
		SampleRate: DefaultSampleRate,
		BufferSize: DefaultBufferSize,
		Profiles:   DefaultSoundProfiles(),
		Mappings:   DefaultProfileMappings(),
	}
}

// Load loads configuration from file
func (cm *ConfigManager) Load() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, use default config
			return nil
		}
		return NewAudioErrorWithCause(ErrInvalidConfig, "failed to read config file", err)
	}

	var config AudioConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return NewAudioErrorWithCause(ErrInvalidConfig, "failed to parse config file", err)
	}

	// Validate loaded configuration
	if err := cm.validateConfig(&config); err != nil {
		return err
	}

	cm.config = &config
	return nil
}

// Save saves configuration to file
func (cm *ConfigManager) Save() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Ensure config directory exists
	if err := os.MkdirAll(filepath.Dir(cm.configPath), 0755); err != nil {
		return NewAudioErrorWithCause(ErrInvalidConfig, "failed to create config directory", err)
	}

	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return NewAudioErrorWithCause(ErrInvalidConfig, "failed to marshal config", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return NewAudioErrorWithCause(ErrInvalidConfig, "failed to write config file", err)
	}

	return nil
}

// GetConfig returns a copy of the current configuration
func (cm *ConfigManager) GetConfig() AudioConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return *cm.config
}

// SetConfig updates the configuration and saves it
func (cm *ConfigManager) SetConfig(config AudioConfig) error {
	if err := cm.validateConfig(&config); err != nil {
		return err
	}

	cm.mu.Lock()
	cm.config = &config
	cm.mu.Unlock()

	return cm.Save()
}

// IsEnabled returns whether audio is enabled
func (cm *ConfigManager) IsEnabled() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.Enabled
}

// SetEnabled enables or disables audio
func (cm *ConfigManager) SetEnabled(enabled bool) error {
	cm.mu.Lock()
	cm.config.Enabled = enabled
	cm.mu.Unlock()
	return cm.Save()
}

// GetVolume returns the current volume level
func (cm *ConfigManager) GetVolume() float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.Volume
}

// SetVolume sets the volume level
func (cm *ConfigManager) SetVolume(volume float64) error {
	if volume < MinVolume || volume > MaxVolume {
		return NewAudioError(ErrInvalidConfig, "volume must be between 0.0 and 1.0")
	}

	cm.mu.Lock()
	cm.config.Volume = volume
	cm.mu.Unlock()
	return cm.Save()
}

// IsMuted returns whether audio is muted
func (cm *ConfigManager) IsMuted() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.Muted
}

// SetMuted mutes or unmutes audio
func (cm *ConfigManager) SetMuted(muted bool) error {
	cm.mu.Lock()
	cm.config.Muted = muted
	cm.mu.Unlock()
	return cm.Save()
}

// GetEffectiveVolume returns the effective volume (considering mute state)
func (cm *ConfigManager) GetEffectiveVolume() float64 {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if cm.config.Muted {
		return 0.0
	}
	return cm.config.Volume
}

// GetProfile returns a sound profile by name
func (cm *ConfigManager) GetProfile(name string) (SoundProfile, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	profile, exists := cm.config.Profiles[name]
	return profile, exists
}

// SetProfile updates or adds a sound profile
func (cm *ConfigManager) SetProfile(profile SoundProfile) error {
	if err := ValidateProfile(profile); err != nil {
		return err
	}

	cm.mu.Lock()
	if cm.config.Profiles == nil {
		cm.config.Profiles = make(map[string]SoundProfile)
	}
	cm.config.Profiles[profile.Name] = profile
	cm.mu.Unlock()

	return cm.Save()
}

// RemoveProfile removes a sound profile
func (cm *ConfigManager) RemoveProfile(name string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.config.Profiles[name]; !exists {
		return NewAudioError(ErrInvalidConfig, "profile not found")
	}

	delete(cm.config.Profiles, name)
	cm.mu.Unlock()

	return cm.Save()
}

// GetProfileForButtonType returns the sound profile for a given button type
func (cm *ConfigManager) GetProfileForButtonType(buttonType ButtonType) (SoundProfile, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Find the mapping for this button type
	for _, mapping := range cm.config.Mappings {
		if mapping.ButtonType == buttonType {
			if mapping.CustomProfile != nil {
				return *mapping.CustomProfile, true
			}

			if profile, exists := cm.config.Profiles[mapping.ProfileName]; exists {
				return profile, true
			}

			break
		}
	}

	// Return default profile as fallback
	switch buttonType {
	case ButtonTypeNumber:
		return cm.config.Profiles[ProfileNameNumber], true
	case ButtonTypeOperator:
		return cm.config.Profiles[ProfileNameOperator], true
	case ButtonTypeSpecial:
		return cm.config.Profiles[ProfileNameClear], true
	case ButtonTypeSystem:
		return cm.config.Profiles[ProfileNameError], true
	default:
		return SoundProfile{}, false
	}
}

// SetProfileMapping sets a custom profile mapping for a button type
func (cm *ConfigManager) SetProfileMapping(buttonType ButtonType, profileName string, customProfile *SoundProfile) error {
	if customProfile != nil {
		if err := ValidateProfile(*customProfile); err != nil {
			return err
		}
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Remove existing mapping for this button type
	for i, mapping := range cm.config.Mappings {
		if mapping.ButtonType == buttonType {
			cm.config.Mappings = append(cm.config.Mappings[:i], cm.config.Mappings[i+1:]...)
			break
		}
	}

	// Add new mapping
	cm.config.Mappings = append(cm.config.Mappings, SoundProfileMapping{
		ButtonType:    buttonType,
		ProfileName:   profileName,
		CustomProfile: customProfile,
	})

	return cm.Save()
}

// ResetToDefaults resets configuration to default values
func (cm *ConfigManager) ResetToDefaults() error {
	cm.mu.Lock()
	cm.config = cm.createDefaultConfig()
	cm.mu.Unlock()
	return cm.Save()
}

// validateConfig validates the entire configuration
func (cm *ConfigManager) validateConfig(config *AudioConfig) error {
	if config.Volume < MinVolume || config.Volume > MaxVolume {
		return NewAudioError(ErrInvalidConfig, "volume must be between 0.0 and 1.0")
	}

	if config.SampleRate <= 0 {
		return NewAudioError(ErrInvalidConfig, "sample rate must be positive")
	}

	if config.BufferSize <= 0 {
		return NewAudioError(ErrInvalidConfig, "buffer size must be positive")
	}

	// Validate all profiles
	for name, profile := range config.Profiles {
		if profile.Name != name {
			return &AudioError{
				Code:    ErrInvalidConfig,
				Message: "profile name mismatch",
			}
		}
		if err := ValidateProfile(profile); err != nil {
			return err
		}
	}

	// Validate mappings
	for _, mapping := range config.Mappings {
		if mapping.ButtonType < ButtonTypeNumber || mapping.ButtonType > ButtonTypeSystem {
			return &AudioError{
				Code:    ErrInvalidConfig,
				Message: "invalid button type in mapping",
			}
		}

		if mapping.CustomProfile != nil {
			if err := ValidateProfile(*mapping.CustomProfile); err != nil {
				return err
			}
		} else if mapping.ProfileName != "" {
			if _, exists := config.Profiles[mapping.ProfileName]; !exists {
				return &AudioError{
					Code:    ErrInvalidConfig,
					Message: "referenced profile not found: " + mapping.ProfileName,
				}
			}
		}
	}

	return nil
}

// ExportConfig exports configuration to a specific path
func (cm *ConfigManager) ExportConfig(path string) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return &AudioError{
			Code:    ErrInvalidConfig,
			Message: "failed to marshal config for export",
			Cause:   err,
		}
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return &AudioError{
			Code:    ErrInvalidConfig,
			Message: "failed to export config",
			Cause:   err,
		}
	}

	return nil
}

// ImportConfig imports configuration from a specific path
func (cm *ConfigManager) ImportConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return &AudioError{
			Code:    ErrInvalidConfig,
			Message: "failed to read import config file",
			Cause:   err,
		}
	}

	var config AudioConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return &AudioError{
			Code:    ErrInvalidConfig,
			Message: "failed to parse import config",
			Cause:   err,
		}
	}

	if err := cm.validateConfig(&config); err != nil {
		return err
	}

	return cm.SetConfig(config)
}