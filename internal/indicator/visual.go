package indicator

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/d-mozulyov/vox/internal/platform"
	"github.com/d-mozulyov/vox/internal/state"
)

// VisualIndicator defines the interface for visual state indication
type VisualIndicator interface {
	// UpdateIcon updates the tray icon based on the current state
	UpdateIcon(state state.State) error
}

// visualIndicator implements the VisualIndicator interface
type visualIndicator struct {
	iconSetter IconSetter
	icons      map[state.State][]byte
}

// IconSetter defines the interface for setting tray icons
// This abstraction allows for testing and decoupling from the tray implementation
type IconSetter interface {
	SetIcon(iconData []byte) error
}

// NewVisualIndicator creates a new visual indicator instance
// iconSetter is the component responsible for actually updating the tray icon
// iconsPath is the directory containing icon files
func NewVisualIndicator(iconSetter IconSetter, iconsPath string) (VisualIndicator, error) {
	logger := platform.GetLogger()

	if iconSetter == nil {
		logger.Error("Visual indicator initialization failed: icon setter is nil")
		return nil, fmt.Errorf("icon setter cannot be nil")
	}

	vi := &visualIndicator{
		iconSetter: iconSetter,
		icons:      make(map[state.State][]byte),
	}

	// Load icons for each state
	if err := vi.loadIcons(iconsPath); err != nil {
		logger.Error("Visual indicator initialization failed: %v", err)
		return nil, fmt.Errorf("failed to load icons: %w", err)
	}

	logger.Info("Visual indicator created successfully")

	return vi, nil
}

// loadIcons loads icon files from the specified directory
// On Windows, ICO format is required. On other platforms, PNG is used.
func (vi *visualIndicator) loadIcons(iconsPath string) error {
	logger := platform.GetLogger()

	// Determine icon extension and filename based on platform
	var iconFiles map[state.State]string

	if runtime.GOOS == "windows" {
		// Windows: use ICO files (contain multiple sizes)
		iconFiles = map[state.State]string{
			state.StateIdle:      "idle.ico",
			state.StateRecording: "recording.ico",
		}
	} else if runtime.GOOS == "darwin" {
		// macOS: detect Retina and use appropriate size
		var size string
		if isRetina() {
			size = "44" // Retina: use 44px for @2x
			logger.Info("Retina display detected, using 44px icons")
		} else {
			size = "22" // Non-Retina: use 22px
			logger.Info("Non-Retina display detected, using 22px icons")
		}
		iconFiles = map[state.State]string{
			state.StateIdle:      "idle_" + size + ".png",
			state.StateRecording: "recording_" + size + ".png",
		}
	} else {
		// Linux: use 24px
		iconFiles = map[state.State]string{
			state.StateIdle:      "idle_24.png",
			state.StateRecording: "recording_24.png",
		}
	}

	logger.Info("Loading icons from: %s", iconsPath)

	for state, filename := range iconFiles {
		iconPath := filepath.Join(iconsPath, filename)
		data, err := os.ReadFile(iconPath)
		if err != nil {
			logger.Error("Failed to read icon %s: %v", filename, err)
			return fmt.Errorf("failed to read icon %s: %w", filename, err)
		}
		vi.icons[state] = data
		logger.Info("Icon loaded: %s (%d bytes)", filename, len(data))
	}

	return nil
}

// UpdateIcon updates the tray icon to reflect the current state
func (vi *visualIndicator) UpdateIcon(s state.State) error {
	logger := platform.GetLogger()

	iconData, ok := vi.icons[s]
	if !ok {
		logger.Error("No icon found for state: %s", s.String())
		return fmt.Errorf("no icon found for state: %s", s.String())
	}

	if err := vi.iconSetter.SetIcon(iconData); err != nil {
		logger.Error("Failed to set icon for state %s: %v", s.String(), err)
		return fmt.Errorf("failed to set icon for state %s: %w", s.String(), err)
	}

	logger.Info("Icon updated to state: %s", s.String())

	return nil
}
