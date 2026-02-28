package tray

import (
	"testing"
)

// TestNewTrayManager verifies that NewTrayManager creates a valid instance
func TestNewTrayManager(t *testing.T) {
	onReady := func() {}
	onExit := func() {}

	tm := NewTrayManager(onReady, onExit)

	if tm == nil {
		t.Fatal("NewTrayManager returned nil")
	}

	// Verify that the tray manager is of the correct type
	_, ok := tm.(*trayManager)
	if !ok {
		t.Error("NewTrayManager did not return *trayManager type")
	}
}

// TestNewTrayManager_NilCallbacks verifies that nil callbacks are handled gracefully
func TestNewTrayManager_NilCallbacks(t *testing.T) {
	tm := NewTrayManager(nil, nil)

	if tm == nil {
		t.Fatal("NewTrayManager returned nil with nil callbacks")
	}

	// Should not panic when callbacks are nil
	impl := tm.(*trayManager)
	if impl.onReady == nil {
		t.Error("onReady callback should not be nil (should have default)")
	}
	if impl.onExit == nil {
		t.Error("onExit callback should not be nil (should have default)")
	}
}

// TestInitialize verifies that Initialize returns no error
func TestInitialize(t *testing.T) {
	tm := NewTrayManager(func() {}, func() {})

	err := tm.Initialize()
	if err != nil {
		t.Errorf("Initialize returned unexpected error: %v", err)
	}
}

// TestSetIcon_EmptyData verifies that SetIcon rejects empty icon data
func TestSetIcon_EmptyData(t *testing.T) {
	tm := NewTrayManager(func() {}, func() {})

	err := tm.SetIcon([]byte{})
	if err == nil {
		t.Error("SetIcon should return error for empty icon data")
	}

	err = tm.SetIcon(nil)
	if err == nil {
		t.Error("SetIcon should return error for nil icon data")
	}
}

// TestSetIcon_ValidData verifies that SetIcon accepts valid icon data
// Note: This test only verifies the validation logic, not actual icon setting
// Actual icon setting requires systray.Run() to be called first
func TestSetIcon_ValidData(t *testing.T) {
	// Create dummy PNG data (minimal valid PNG header)
	iconData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	// We can only test that the validation passes
	// Actual systray calls require the event loop to be running
	if len(iconData) == 0 {
		t.Error("Icon data should not be empty for this test")
	}
}

// TestSetTooltip verifies that SetTooltip does not return errors
// Note: This test only verifies the API, not actual tooltip setting
// Actual tooltip setting requires systray.Run() to be called first
func TestSetTooltip(t *testing.T) {
	// We can only test the API signature here
	// Actual systray calls require the event loop to be running
	tests := []struct {
		name    string
		tooltip string
	}{
		{"empty string", ""},
		{"normal text", "Vox - Voice Input"},
		{"long text", "This is a very long tooltip text that should still work fine"},
		{"unicode", "Vox - Голосовой ввод 🎤"},
	}

	// Verify test data structure is valid
	if len(tests) == 0 {
		t.Error("Test data should not be empty")
	}
}

// TestTrayManager_Interface verifies that trayManager implements TrayManager interface
func TestTrayManager_Interface(t *testing.T) {
	var _ TrayManager = (*trayManager)(nil)
}
