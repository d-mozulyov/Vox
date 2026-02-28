package main

import (
	"fmt"
	"os"
)

// Version is set during build via -ldflags
var Version = "0.0.0"

func main() {
	fmt.Println("Vox - Voice Input Assistant")
	fmt.Printf("Version: %s\n", Version)

	// TODO: Initialize application
	// - System tray icon
	// - Hotkey registration
	// - Audio recording
	// - AI transcription

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

	// Start main application loop
	fmt.Println("Starting Vox...")
	// TODO: Implement main loop
}

func printHelp() {
	fmt.Println("\nUsage:")
	fmt.Println("  vox           Start the application")
	fmt.Println("  vox version   Show version information")
	fmt.Println("  vox help      Show this help message")
}
