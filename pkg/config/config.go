package config

import (
	"os"
	"path/filepath"
)

// Config holds application configuration
type Config struct {
	Hotkey  HotkeyConfig
	Audio   AudioConfig
	Logging LoggingConfig
}

// HotkeyConfig holds hotkey configuration
type HotkeyConfig struct {
	Enabled bool
	// Default: Alt+Shift+V
	UseAlt   bool
	UseShift bool
	UseCtrl  bool
	Key      string
}

// AudioConfig holds audio configuration
type AudioConfig struct {
	Enabled bool
	Volume  float64 // 0.0 to 1.0
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level    string // debug, info, warn, error
	FilePath string
}

// Default returns the default configuration
func Default() *Config {
	homeDir, _ := os.UserHomeDir()
	logPath := filepath.Join(homeDir, ".vox", "vox.log")

	return &Config{
		Hotkey: HotkeyConfig{
			Enabled:  true,
			UseAlt:   true,
			UseShift: true,
			UseCtrl:  false,
			Key:      "V",
		},
		Audio: AudioConfig{
			Enabled: true,
			Volume:  0.8,
		},
		Logging: LoggingConfig{
			Level:    "info",
			FilePath: logPath,
		},
	}
}
