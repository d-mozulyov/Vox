package tray_test

import (
	"log"

	"github.com/d-mozulyov/vox/internal/tray"
)

// Example demonstrates basic usage of TrayManager
func Example() {
	// Create callbacks for tray lifecycle events
	onReady := func() {
		log.Println("Tray is ready")
		// Initialize other components here:
		// - Load initial icon
		// - Register hotkeys
		// - Setup state machine
	}

	onExit := func() {
		log.Println("Cleaning up resources")
		// Cleanup resources here:
		// - Unregister hotkeys
		// - Close audio devices
		// - Save state
	}

	// Create tray manager
	trayManager := tray.NewTrayManager(onReady, onExit)

	// Initialize (validation only, actual init happens in Run)
	if err := trayManager.Initialize(); err != nil {
		log.Fatalf("Failed to initialize tray: %v", err)
	}

	// Set initial tooltip
	if err := trayManager.SetTooltip("Vox - Voice Input Assistant"); err != nil {
		log.Printf("Warning: Failed to set tooltip: %v", err)
	}

	// Run the tray event loop (blocking)
	// This must be called from the main goroutine
	trayManager.Run()
}

// Example_withIconUpdate demonstrates updating the tray icon
func Example_withIconUpdate() {
	var trayManager tray.TrayManager

	onReady := func() {
		// Load and set initial icon
		iconData := []byte{0x89, 0x50, 0x4E, 0x47} // PNG header (example)
		if err := trayManager.SetIcon(iconData); err != nil {
			log.Printf("Warning: Failed to set icon: %v", err)
		}
	}

	trayManager = tray.NewTrayManager(onReady, func() {})
	trayManager.Run()
}
