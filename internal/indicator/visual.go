package indicator

import (
	"fmt"
	"os"
	"path/filepath"

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
// Icons are loaded in multiple sizes (16x16, 32x32, 64x64) for different DPI settings
// The system will automatically select the appropriate size based on display settings
func (vi *visualIndicator) loadIcons(iconsPath string) error {
	logger := platform.GetLogger()

	// Use 32x32 as the default size for system tray icons
	// This is a good balance for most displays and DPI settings
	iconFiles := map[state.State]string{
		state.StateIdle:       "idle_32.png",
		state.StateRecording:  "recording_32.png",
		state.StateProcessing: "processing_32.png",
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
