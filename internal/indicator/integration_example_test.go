package indicator_test

import (
	"testing"

	"github.com/d-mozulyov/vox/internal/indicator"
	"github.com/d-mozulyov/vox/internal/state"
)

// TestIntegration_StateMachineWithIndicatorManager demonstrates how to integrate
// the State Machine with the Indicator Manager
func TestIntegration_StateMachineWithIndicatorManager(t *testing.T) {
	// Create state machine
	sm := state.NewStateMachine()

	// Create indicator manager
	im := indicator.NewIndicatorManager()

	// Subscribe indicator manager to state changes
	sm.Subscribe(im.OnStateChange)

	// Create mock indicators to verify they are called
	visualCalled := false
	audioCalled := false

	mockVisual := &mockVisualIndicator{
		updateIconFunc: func(s state.State) error {
			visualCalled = true
			return nil
		},
	}

	mockAudio := &mockAudioIndicator{
		playSoundFunc: func(from, to state.State) error {
			audioCalled = true
			return nil
		},
	}

	im.SetVisualIndicator(mockVisual)
	im.SetAudioIndicator(mockAudio)

	// Trigger state transition
	err := sm.Transition(state.StateRecording)
	if err != nil {
		t.Fatalf("Failed to transition state: %v", err)
	}

	// Verify indicators were triggered
	if !visualCalled {
		t.Error("Visual indicator was not called")
	}
	if !audioCalled {
		t.Error("Audio indicator was not called")
	}
}

// mockVisualIndicator is a simple mock for testing
type mockVisualIndicator struct {
	updateIconFunc func(state.State) error
}

func (m *mockVisualIndicator) UpdateIcon(s state.State) error {
	if m.updateIconFunc != nil {
		return m.updateIconFunc(s)
	}
	return nil
}

// mockAudioIndicator is a simple mock for testing
type mockAudioIndicator struct {
	playSoundFunc func(state.State, state.State) error
}

func (m *mockAudioIndicator) PlaySound(from, to state.State) error {
	if m.playSoundFunc != nil {
		return m.playSoundFunc(from, to)
	}
	return nil
}
