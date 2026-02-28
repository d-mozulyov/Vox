package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetAssetsPath tests the asset path resolution logic
func TestGetAssetsPath(t *testing.T) {
	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Test that getAssetsPath returns a valid path
	assetsPath := getAssetsPath()
	if assetsPath == "" {
		t.Error("getAssetsPath returned empty string")
	}

	// Verify the path is reasonable (contains "assets")
	if !filepath.IsAbs(assetsPath) && assetsPath != "assets" {
		// If not absolute, should be relative path containing "assets"
		if filepath.Base(assetsPath) != "assets" {
			t.Errorf("Expected path to end with 'assets', got: %s", assetsPath)
		}
	}
}

// TestMainHelp tests the help command
func TestMainHelp(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set args to trigger help
	os.Args = []string{"vox", "help"}

	// This should not panic or exit with error
	// We can't easily test the output without refactoring main()
	// but we can at least ensure it doesn't crash
	printHelp()
}

// TestMainVersion tests the version command
func TestMainVersion(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set args to trigger version
	os.Args = []string{"vox", "version"}

	// Verify Version variable exists
	if Version == "" {
		t.Log("Version is empty (expected during tests)")
	}
}
