package hotkey

import (
	"testing"
	"time"
)

// TestNewHotkeyManager tests the creation of a new hotkey manager
func TestNewHotkeyManager(t *testing.T) {
	manager := NewHotkeyManager()
	if manager == nil {
		t.Fatal("NewHotkeyManager returned nil")
	}
}

// TestHotkeyString tests the string representation of a hotkey
func TestHotkeyString(t *testing.T) {
	tests := []struct {
		name     string
		hotkey   Hotkey
		expected string
	}{
		{
			name:     "Alt+Shift+V",
			hotkey:   Hotkey{Modifiers: []Modifier{ModAlt, ModShift}, Key: KeyV},
			expected: "Alt+Shift+V",
		},
		{
			name:     "Ctrl+C",
			hotkey:   Hotkey{Modifiers: []Modifier{ModCtrl}, Key: KeyC},
			expected: "Ctrl+C",
		},
		{
			name:     "Win+R",
			hotkey:   Hotkey{Modifiers: []Modifier{ModWin}, Key: KeyR},
			expected: "Win+R",
		},
		{
			name:     "No modifiers",
			hotkey:   Hotkey{Modifiers: []Modifier{}, Key: KeyA},
			expected: "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hotkey.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestRegisterHotkey tests successful hotkey registration
func TestRegisterHotkey(t *testing.T) {
	manager := NewHotkeyManager()

	callback := func() {
		// Callback would be triggered by actual hotkey press
	}

	// Use a less common hotkey combination to avoid conflicts
	hotkey := Hotkey{
		Modifiers: []Modifier{ModAlt, ModShift, ModCtrl},
		Key:       KeyZ,
	}

	err := manager.Register(hotkey, callback)
	if err != nil {
		t.Logf("Warning: Failed to register hotkey (may be taken by another app): %v", err)
		t.Skip("Skipping test due to hotkey registration failure")
	}

	// Clean up
	defer func() {
		if err := manager.Unregister(hotkey); err != nil {
			t.Logf("Warning: Failed to unregister hotkey: %v", err)
		}
	}()

	// Note: We can't actually trigger the hotkey programmatically in a test
	// This test only verifies that registration succeeds without error
}

// TestRegisterDuplicateHotkey tests that registering the same hotkey twice fails
func TestRegisterDuplicateHotkey(t *testing.T) {
	manager := NewHotkeyManager()

	hotkey := Hotkey{
		Modifiers: []Modifier{ModAlt, ModShift, ModCtrl},
		Key:       KeyY,
	}

	callback := func() {}

	// First registration should succeed
	err := manager.Register(hotkey, callback)
	if err != nil {
		t.Logf("Warning: Failed to register hotkey (may be taken by another app): %v", err)
		t.Skip("Skipping test due to hotkey registration failure")
	}

	defer func() {
		if err := manager.Unregister(hotkey); err != nil {
			t.Logf("Warning: Failed to unregister hotkey: %v", err)
		}
	}()

	// Second registration should fail
	err = manager.Register(hotkey, callback)
	if err == nil {
		t.Error("Expected error when registering duplicate hotkey, got nil")
	}
}

// TestUnregisterHotkey tests hotkey unregistration
func TestUnregisterHotkey(t *testing.T) {
	manager := NewHotkeyManager()

	hotkey := Hotkey{
		Modifiers: []Modifier{ModAlt, ModShift, ModCtrl},
		Key:       KeyX,
	}

	callback := func() {}

	// Register the hotkey
	err := manager.Register(hotkey, callback)
	if err != nil {
		t.Logf("Warning: Failed to register hotkey (may be taken by another app): %v", err)
		t.Skip("Skipping test due to hotkey registration failure")
	}

	// Unregister should succeed
	err = manager.Unregister(hotkey)
	if err != nil {
		t.Errorf("Failed to unregister hotkey: %v", err)
	}

	// Unregistering again should fail
	err = manager.Unregister(hotkey)
	if err == nil {
		t.Error("Expected error when unregistering non-existent hotkey, got nil")
	}
}

// TestUnregisterNonExistentHotkey tests unregistering a hotkey that was never registered
func TestUnregisterNonExistentHotkey(t *testing.T) {
	manager := NewHotkeyManager()

	hotkey := Hotkey{
		Modifiers: []Modifier{ModAlt},
		Key:       KeyW,
	}

	err := manager.Unregister(hotkey)
	if err == nil {
		t.Error("Expected error when unregistering non-existent hotkey, got nil")
	}
}

// TestUnregisterAll tests unregistering all hotkeys
func TestUnregisterAll(t *testing.T) {
	manager := NewHotkeyManager()

	// Register multiple hotkeys
	hotkeys := []Hotkey{
		{Modifiers: []Modifier{ModAlt, ModShift, ModCtrl}, Key: KeyQ},
		{Modifiers: []Modifier{ModAlt, ModShift, ModCtrl}, Key: KeyW},
		{Modifiers: []Modifier{ModAlt, ModShift, ModCtrl}, Key: KeyE},
	}

	registeredCount := 0
	for _, hk := range hotkeys {
		err := manager.Register(hk, func() {})
		if err != nil {
			t.Logf("Warning: Failed to register hotkey %s: %v", hk.String(), err)
		} else {
			registeredCount++
		}
	}

	if registeredCount == 0 {
		t.Skip("Could not register any hotkeys, skipping test")
	}

	// Unregister all
	err := manager.UnregisterAll()
	if err != nil {
		t.Errorf("Failed to unregister all hotkeys: %v", err)
	}

	// Try to unregister again - should succeed with no errors (map is empty)
	err = manager.UnregisterAll()
	if err != nil {
		t.Errorf("UnregisterAll on empty manager should not error: %v", err)
	}
}

// TestUnregisterAllWithNoHotkeys tests UnregisterAll when no hotkeys are registered
func TestUnregisterAllWithNoHotkeys(t *testing.T) {
	manager := NewHotkeyManager()

	err := manager.UnregisterAll()
	if err != nil {
		t.Errorf("UnregisterAll with no hotkeys should not error: %v", err)
	}
}

// TestConcurrentRegistration tests concurrent hotkey registration
func TestConcurrentRegistration(t *testing.T) {
	manager := NewHotkeyManager()

	// Use different hotkeys to avoid conflicts
	hotkeys := []Hotkey{
		{Modifiers: []Modifier{ModAlt, ModShift, ModCtrl}, Key: KeyA},
		{Modifiers: []Modifier{ModAlt, ModShift, ModCtrl}, Key: KeyB},
		{Modifiers: []Modifier{ModAlt, ModShift, ModCtrl}, Key: KeyC},
		{Modifiers: []Modifier{ModAlt, ModShift, ModCtrl}, Key: KeyD},
	}

	done := make(chan bool, len(hotkeys))

	// Register hotkeys concurrently
	for _, hk := range hotkeys {
		go func(hotkey Hotkey) {
			err := manager.Register(hotkey, func() {})
			if err != nil {
				t.Logf("Warning: Failed to register hotkey %s: %v", hotkey.String(), err)
			}
			done <- true
		}(hk)
	}

	// Wait for all goroutines to complete
	for i := 0; i < len(hotkeys); i++ {
		<-done
	}

	// Clean up
	if err := manager.UnregisterAll(); err != nil {
		t.Logf("Warning: Failed to unregister all hotkeys: %v", err)
	}
}

// TestHotkeyCleanupOnShutdown tests that hotkeys are properly cleaned up
func TestHotkeyCleanupOnShutdown(t *testing.T) {
	manager := NewHotkeyManager()

	hotkey := Hotkey{
		Modifiers: []Modifier{ModAlt, ModShift, ModCtrl},
		Key:       KeyF,
	}

	err := manager.Register(hotkey, func() {})
	if err != nil {
		t.Logf("Warning: Failed to register hotkey: %v", err)
		t.Skip("Skipping test due to hotkey registration failure")
	}

	// Simulate shutdown by unregistering all
	err = manager.UnregisterAll()
	if err != nil {
		t.Errorf("Failed to cleanup hotkeys on shutdown: %v", err)
	}

	// Give goroutines time to stop
	time.Sleep(100 * time.Millisecond)

	// Try to register the same hotkey again - should succeed if cleanup worked
	err = manager.Register(hotkey, func() {})
	if err != nil {
		t.Errorf("Failed to re-register hotkey after cleanup: %v", err)
	}

	// Final cleanup
	if err := manager.UnregisterAll(); err != nil {
		t.Logf("Warning: Failed final cleanup: %v", err)
	}
}

// TestModifierConversion tests the conversion of modifiers
func TestModifierConversion(t *testing.T) {
	tests := []struct {
		name      string
		modifiers []Modifier
	}{
		{
			name:      "Single modifier",
			modifiers: []Modifier{ModAlt},
		},
		{
			name:      "Multiple modifiers",
			modifiers: []Modifier{ModAlt, ModShift, ModCtrl},
		},
		{
			name:      "All modifiers",
			modifiers: []Modifier{ModAlt, ModShift, ModCtrl, ModWin},
		},
		{
			name:      "No modifiers",
			modifiers: []Modifier{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertModifiers(tt.modifiers)
			if len(result) != len(tt.modifiers) {
				t.Errorf("Expected %d modifiers, got %d", len(tt.modifiers), len(result))
			}
		})
	}
}

// TestKeyConversion tests the conversion of keys
func TestKeyConversion(t *testing.T) {
	keys := []Key{KeyA, KeyB, KeyC, KeyV, KeyZ}

	for _, k := range keys {
		t.Run(string(rune('A'+int(k))), func(t *testing.T) {
			result := convertKey(k)
			// Just verify it doesn't panic
			_ = result
		})
	}
}
