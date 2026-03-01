package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/d-mozulyov/vox/internal/hotkey"
	"github.com/d-mozulyov/vox/internal/indicator"
	"github.com/d-mozulyov/vox/internal/platform"
	"github.com/d-mozulyov/vox/internal/state"
	"github.com/d-mozulyov/vox/internal/tray"
)

// Version is set during build via -ldflags
var Version = "0.0.0"

func main() {
	fmt.Println("Vox - Voice Input Assistant")
	fmt.Printf("Version: %s\n", Version)

	// Handle command-line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Println(Version)
		case "help":
			printHelp()
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			printHelp()
			os.Exit(1)
		}
		return
	}

	// Initialize logger
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get user home directory: %v\n", err)
		os.Exit(1)
	}
	logFilePath := filepath.Join(homeDir, ".vox", "vox.log")
	if err := platform.InitLogger(platform.LogLevelInfo, logFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger := platform.GetLogger()
	logger.Info("Vox starting, version: %s", Version)

	// Start main application
	fmt.Println("Starting Vox...")
	if err := run(); err != nil {
		logger.Fatal("Application error: %v", err)
	}
}

// run initializes and runs the application
// Integration flow:
// 1. Initialize State Machine (manages application state)
// 2. Initialize Hotkey Manager (registers Alt+Shift+V)
// 3. Initialize Indicator Manager (coordinates visual + audio feedback)
// 4. Initialize Tray Manager (system tray icon and menu)
// 5. In onReady callback (when tray is ready):
//    - Initialize Visual Indicator (icon updates)
//    - Initialize Audio Indicator (sound feedback)
//    - Subscribe Indicator Manager to state changes
//    - Register hotkey with callback that transitions states
// 6. Run tray event loop (blocking)
//
// State flow: Hotkey press → State transition → Indicator update (visual + audio)
// Cleanup: defer statements ensure proper resource cleanup on exit
func run() error {
	logger := platform.GetLogger()

	// Initialize State Machine
	stateMachine := state.NewStateMachine()
	logger.Info("State machine initialized")

	// Initialize Hotkey Manager
	hotkeyManager := hotkey.NewHotkeyManager()
	defer func() {
		if err := hotkeyManager.UnregisterAll(); err != nil {
			logger.Error("Error unregistering hotkeys: %v", err)
		} else {
			logger.Info("Hotkeys unregistered")
		}
	}()

	// Initialize Indicator Manager
	indicatorManager := indicator.NewIndicatorManager()
	logger.Info("Indicator manager initialized")

	// Get paths to assets
	assetsPath := getAssetsPath()
	iconsPath := filepath.Join(assetsPath, "icons")
	soundsPath := filepath.Join(assetsPath, "sounds")

	// Variable to hold tray manager (will be initialized in onReady)
	var trayManager tray.TrayManager

	// onReady callback - called when tray is initialized
	onReady := func() {
		logger.Info("Tray is ready, initializing components...")

		// Initialize Visual Indicator
		visualIndicator, err := indicator.NewVisualIndicator(trayManager, iconsPath)
		if err != nil {
			logger.Warn("Failed to initialize visual indicator: %v", err)
		} else {
			indicatorManager.SetVisualIndicator(visualIndicator)
			logger.Info("Visual indicator initialized")

			// Set initial icon to idle state
			if err := visualIndicator.UpdateIcon(state.StateIdle); err != nil {
				logger.Warn("Failed to set initial icon: %v", err)
			}
		}

		// Initialize Audio Indicator
		audioIndicator, err := indicator.NewAudioIndicator(soundsPath)
		if err != nil {
			logger.Warn("Failed to initialize audio indicator: %v", err)
		} else {
			indicatorManager.SetAudioIndicator(audioIndicator)
			logger.Info("Audio indicator initialized")
		}

		// Subscribe Indicator Manager to state changes
		stateMachine.Subscribe(indicatorManager.OnStateChange)
		logger.Info("Indicator manager subscribed to state changes")

		// Register hotkey Alt+Shift+V
		hk := hotkey.Hotkey{
			Modifiers: []hotkey.Modifier{hotkey.ModAlt, hotkey.ModShift},
			Key:       hotkey.KeyV,
		}

		// Hotkey callback - toggles state
		hotkeyCallback := func() {
			logger.Info("Hotkey pressed: %s", hk.String())

			currentState := stateMachine.GetState()
			var nextState state.State

			switch currentState {
			case state.StateIdle:
				nextState = state.StateRecording
			case state.StateRecording:
				nextState = state.StateProcessing
			case state.StateProcessing:
				// Processing state transitions to Idle automatically
				// Hotkey press during processing is ignored
				logger.Info("Hotkey pressed during processing, ignoring")
				return
			}

			if err := stateMachine.Transition(nextState); err != nil {
				logger.Error("Error transitioning state: %v", err)
			} else {
				logger.Info("State transitioned: %s -> %s", currentState, nextState)
			}
		}

		// Register the hotkey
		if err := hotkeyManager.Register(hk, hotkeyCallback); err != nil {
			logger.Warn("Failed to register hotkey %s: %v. Application will work without hotkeys.", hk.String(), err)
		} else {
			logger.Info("Hotkey registered: %s", hk.String())
		}

		logger.Info("Application initialized successfully")
	}

	// onExit callback - called when user selects Exit from menu
	onExit := func() {
		logger.Info("Cleaning up resources...")
		// Cleanup is handled by defer statements in run()
	}

	// Initialize Tray Manager
	trayManager = tray.NewTrayManager(onReady, onExit)
	logger.Info("Tray manager created")

	// Run tray (blocking call)
	// This must be called from the main goroutine
	trayManager.Run()

	return nil
}

// getAssetsPath returns the path to the assets directory
// It tries multiple locations in the following order:
// 1. Current working directory (development mode)
// 2. Next to executable (production mode)
// 3. Parent directory of executable (some build configurations)
func getAssetsPath() string {
	logger := platform.GetLogger()

	// Try current working directory first (development mode)
	cwd, err := os.Getwd()
	if err == nil {
		assetsPath := filepath.Join(cwd, "assets")
		if _, err := os.Stat(assetsPath); err == nil {
			logger.Info("Assets found in working directory: %s", assetsPath)
			return assetsPath
		}
	}

	// Try to find assets directory relative to executable
	exePath, err := os.Executable()
	if err != nil {
		logger.Warn("Failed to get executable path: %v", err)
		return "assets"
	}

	exeDir := filepath.Dir(exePath)

	// Check if assets directory exists next to executable
	assetsPath := filepath.Join(exeDir, "assets")
	if _, err := os.Stat(assetsPath); err == nil {
		logger.Info("Assets found next to executable: %s", assetsPath)
		return assetsPath
	}

	// Check if assets directory exists in parent directory
	assetsPath = filepath.Join(exeDir, "..", "assets")
	if _, err := os.Stat(assetsPath); err == nil {
		logger.Info("Assets found in parent directory: %s", assetsPath)
		return assetsPath
	}

	// Default to "assets" relative to current directory
	logger.Warn("Assets directory not found, using default: assets")
	return "assets"
}


func printHelp() {
	fmt.Println("\nUsage:")
	fmt.Println("  vox           Start the application")
	fmt.Println("  vox version   Show version information")
	fmt.Println("  vox help      Show this help message")
}
