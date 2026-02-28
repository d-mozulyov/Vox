package indicator

import (
	"testing"

	"github.com/d-mozulyov/vox/internal/state"
)

// mockVisualIndicator is a mock implementation of VisualIndicator for testing
type mockVisualIndicator struct {
	callCount int
}

func (m *mockVisualIndicator) UpdateIcon(s state.State) error {
	m.callCount++
	return nil
}

// mockAudioIndicator is a mock implementation of AudioIndicator for testing
type mockAudioIndicator struct {
	callCount int
}

func (m *mockAudioIndicator) PlaySound(from, to state.State) error {
	m.callCount++
	return nil
}

// TestIndicatorManager_OnStateChange tests basic coordination of indicators
func TestIndicatorManager_OnStateChange(t *testing.T) {
	manager := NewIndicatorManager()

	visualMock := &mockVisualIndicator{}
	audioMock := &mockAudioIndicator{}

	manager.SetVisualIndicator(visualMock)
	manager.SetAudioIndicator(audioMock)

	// Trigger state change
	manager.OnStateChange(state.StateIdle, state.StateRecording)

	// Verify both indicators were called
	if visualMock.callCount != 1 {
		t.Errorf("Expected visual indicator to be called once, got %d", visualMock.callCount)
	}
	if audioMock.callCount != 1 {
		t.Errorf("Expected audio indicator to be called once, got %d", audioMock.callCount)
	}
}

// TestIndicatorManager_NoIndicators tests that manager doesn't panic with no indicators
func TestIndicatorManager_NoIndicators(t *testing.T) {
	manager := NewIndicatorManager()

	// Should not panic when no indicators are set
	manager.OnStateChange(state.StateIdle, state.StateRecording)
}
