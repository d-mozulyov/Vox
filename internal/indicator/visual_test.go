package indicator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/d-mozulyov/vox/internal/state"
)

// mockIconSetter is a mock implementation of IconSetter for testing
type mockIconSetter struct {
	setIconFunc func([]byte) error
	lastIcon    []byte
}

func (m *mockIconSetter) SetIcon(iconData []byte) error {
	m.lastIcon = iconData
	if m.setIconFunc != nil {
		return m.setIconFunc(iconData)
	}
	return nil
}

// setupTestIcons creates temporary icon files for testing
func setupTestIcons(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	// Create test icon files with size suffix (32x32 is the default)
	icons := map[string][]byte{
		"idle_32.png":       []byte("idle_icon_data"),
		"recording_32.png":  []byte("recording_icon_data"),
		"processing_32.png": []byte("processing_icon_data"),
	}

	for filename, data := range icons {
		path := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatalf("Failed to create test icon %s: %v", filename, err)
		}
	}

	return tmpDir
}

func TestNewVisualIndicator_Success(t *testing.T) {
	iconsPath := setupTestIcons(t)
	mockSetter := &mockIconSetter{}

	vi, err := NewVisualIndicator(mockSetter, iconsPath)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if vi == nil {
		t.Error("Expected visual indicator instance, got nil")
	}
}

func TestNewVisualIndicator_NilIconSetter(t *testing.T) {
	iconsPath := setupTestIcons(t)

	vi, err := NewVisualIndicator(nil, iconsPath)

	if err == nil {
		t.Error("Expected error for nil icon setter, got nil")
	}
	if vi != nil {
		t.Error("Expected nil visual indicator, got instance")
	}
}

func TestNewVisualIndicator_MissingIconFile(t *testing.T) {
	tmpDir := t.TempDir()
	mockSetter := &mockIconSetter{}

	// Create only some icons, not all (using new naming convention)
	idlePath := filepath.Join(tmpDir, "idle_32.png")
	if err := os.WriteFile(idlePath, []byte("idle_data"), 0644); err != nil {
		t.Fatalf("Failed to create test icon: %v", err)
	}

	vi, err := NewVisualIndicator(mockSetter, tmpDir)

	if err == nil {
		t.Error("Expected error for missing icon files, got nil")
	}
	if vi != nil {
		t.Error("Expected nil visual indicator, got instance")
	}
}

func TestNewVisualIndicator_InvalidPath(t *testing.T) {
	mockSetter := &mockIconSetter{}
	invalidPath := "/nonexistent/path/to/icons"

	vi, err := NewVisualIndicator(mockSetter, invalidPath)

	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
	if vi != nil {
		t.Error("Expected nil visual indicator, got instance")
	}
}

func TestUpdateIcon_IdleState(t *testing.T) {
	iconsPath := setupTestIcons(t)
	mockSetter := &mockIconSetter{}

	vi, err := NewVisualIndicator(mockSetter, iconsPath)
	if err != nil {
		t.Fatalf("Failed to create visual indicator: %v", err)
	}

	err = vi.UpdateIcon(state.StateIdle)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if string(mockSetter.lastIcon) != "idle_icon_data" {
		t.Errorf("Expected idle icon data, got: %s", string(mockSetter.lastIcon))
	}
}

func TestUpdateIcon_RecordingState(t *testing.T) {
	iconsPath := setupTestIcons(t)
	mockSetter := &mockIconSetter{}

	vi, err := NewVisualIndicator(mockSetter, iconsPath)
	if err != nil {
		t.Fatalf("Failed to create visual indicator: %v", err)
	}

	err = vi.UpdateIcon(state.StateRecording)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if string(mockSetter.lastIcon) != "recording_icon_data" {
		t.Errorf("Expected recording icon data, got: %s", string(mockSetter.lastIcon))
	}
}

func TestUpdateIcon_ProcessingState(t *testing.T) {
	iconsPath := setupTestIcons(t)
	mockSetter := &mockIconSetter{}

	vi, err := NewVisualIndicator(mockSetter, iconsPath)
	if err != nil {
		t.Fatalf("Failed to create visual indicator: %v", err)
	}

	err = vi.UpdateIcon(state.StateProcessing)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if string(mockSetter.lastIcon) != "processing_icon_data" {
		t.Errorf("Expected processing icon data, got: %s", string(mockSetter.lastIcon))
	}
}

func TestUpdateIcon_SetIconError(t *testing.T) {
	iconsPath := setupTestIcons(t)
	expectedError := fmt.Errorf("failed to set icon")
	mockSetter := &mockIconSetter{
		setIconFunc: func([]byte) error {
			return expectedError
		},
	}

	vi, err := NewVisualIndicator(mockSetter, iconsPath)
	if err != nil {
		t.Fatalf("Failed to create visual indicator: %v", err)
	}

	err = vi.UpdateIcon(state.StateIdle)

	if err == nil {
		t.Error("Expected error from SetIcon, got nil")
	}
}

func TestUpdateIcon_AllStatesSequentially(t *testing.T) {
	iconsPath := setupTestIcons(t)
	mockSetter := &mockIconSetter{}

	vi, err := NewVisualIndicator(mockSetter, iconsPath)
	if err != nil {
		t.Fatalf("Failed to create visual indicator: %v", err)
	}

	states := []struct {
		state        state.State
		expectedData string
	}{
		{state.StateIdle, "idle_icon_data"},
		{state.StateRecording, "recording_icon_data"},
		{state.StateProcessing, "processing_icon_data"},
		{state.StateIdle, "idle_icon_data"},
	}

	for _, tc := range states {
		err := vi.UpdateIcon(tc.state)
		if err != nil {
			t.Errorf("Failed to update icon for state %s: %v", tc.state.String(), err)
		}
		if string(mockSetter.lastIcon) != tc.expectedData {
			t.Errorf("State %s: expected %s, got %s",
				tc.state.String(), tc.expectedData, string(mockSetter.lastIcon))
		}
	}
}
