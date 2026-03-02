// Package tray provides a wrapper around getlantern/systray for managing
// system tray icon and menu.
//
// Example usage:
//
//	func main() {
//	    onReady := func() {
//	        log.Println("Tray is ready")
//	        // Initialize other components here
//	    }
//
//	    onExit := func() {
//	        log.Println("Cleaning up...")
//	        // Cleanup resources here
//	    }
//
//	    tray := tray.NewTrayManager(onReady, onExit)
//	    tray.Run() // Blocking call
//	}
package tray

import (
	"fmt"

	"fyne.io/systray"
	"github.com/d-mozulyov/vox/internal/platform"
)

// TrayManager defines the interface for managing system tray icon and menu
type TrayManager interface {
	// Initialize creates system tray icon and menu
	// Returns error if initialization fails (critical error)
	Initialize() error

	// SetIcon updates the tray icon
	// iconData should be PNG image data
	SetIcon(iconData []byte) error

	// SetTooltip updates the tooltip text shown when hovering over the icon
	SetTooltip(text string) error

	// Run starts the tray event loop (blocking call)
	// This should be called in the main goroutine
	Run()

	// Quit removes tray icon and exits the application
	Quit()
}

// trayManager implements the TrayManager interface using getlantern/systray
type trayManager struct {
	onReady func()
	onExit  func()

	// Menu items
	menuSettings *systray.MenuItem
	menuExit     *systray.MenuItem
}

// NewTrayManager creates a new tray manager instance
// onReady is called when the tray is initialized and ready
// onExit is called when the user selects "Exit" from the menu
func NewTrayManager(onReady func(), onExit func()) TrayManager {
	if onReady == nil {
		onReady = func() {}
	}
	if onExit == nil {
		onExit = func() {}
	}

	return &trayManager{
		onReady: onReady,
		onExit:  onExit,
	}
}

// Initialize creates system tray icon and menu
func (tm *trayManager) Initialize() error {
	// Note: systray.Run is blocking, so we can't return from Initialize
	// The actual initialization happens in the onReady callback
	// This method just validates that we can proceed
	return nil
}

// SetIcon updates the tray icon with the provided PNG image data
func (tm *trayManager) SetIcon(iconData []byte) error {
	if len(iconData) == 0 {
		return fmt.Errorf("icon data cannot be empty")
	}

	systray.SetIcon(iconData)
	return nil
}

// SetTooltip updates the tooltip text
func (tm *trayManager) SetTooltip(text string) error {
	systray.SetTooltip(text)
	return nil
}

// Run starts the tray event loop (blocking)
// This must be called from the main goroutine
func (tm *trayManager) Run() {
	systray.Run(tm.onReadyWrapper, tm.onExitWrapper)
}

// onReadyWrapper is called by systray when the tray is ready
func (tm *trayManager) onReadyWrapper() {
	logger := platform.GetLogger()

	// Set initial tooltip
	systray.SetTooltip("Vox - Voice Input Assistant")
	logger.Info("Tray tooltip set")

	// Create menu items
	tm.menuSettings = systray.AddMenuItem("Settings", "Open settings window")
	tm.menuSettings.Disable() // Placeholder - will be enabled in future
	logger.Info("Settings menu item created (disabled)")

	systray.AddSeparator()

	tm.menuExit = systray.AddMenuItem("Exit", "Exit the application")
	logger.Info("Exit menu item created")

	// Start goroutine to handle menu clicks
	go tm.handleMenuClicks()

	logger.Info("Tray initialized successfully")

	// Call user's onReady callback
	tm.onReady()
}

// onExitWrapper is called by systray when exiting
func (tm *trayManager) onExitWrapper() {
	logger := platform.GetLogger()
	logger.Info("Tray manager: cleaning up resources")
	tm.onExit()
}

// handleMenuClicks listens for menu item clicks and handles them
func (tm *trayManager) handleMenuClicks() {
	logger := platform.GetLogger()
	for {
		select {
		case <-tm.menuSettings.ClickedCh:
			// Placeholder for future settings window
			logger.Info("Settings clicked (not implemented yet)")

		case <-tm.menuExit.ClickedCh:
			logger.Info("Exit clicked, quitting application")
			tm.Quit()
			return
		}
	}
}

// Quit removes tray icon and exits the application
func (tm *trayManager) Quit() {
	systray.Quit()
}
