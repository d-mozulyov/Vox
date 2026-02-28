package state

import "testing"

// TestStateMachine tests basic state machine functionality
func TestStateMachine(t *testing.T) {
	sm := NewStateMachine()

	// Test initial state
	if sm.GetState() != StateIdle {
		t.Errorf("Expected initial state Idle, got %s", sm.GetState())
	}

	// Test valid transitions: Idle -> Recording -> Processing -> Idle
	if err := sm.Transition(StateRecording); err != nil {
		t.Errorf("Idle->Recording failed: %v", err)
	}
	if err := sm.Transition(StateProcessing); err != nil {
		t.Errorf("Recording->Processing failed: %v", err)
	}
	if err := sm.Transition(StateIdle); err != nil {
		t.Errorf("Processing->Idle failed: %v", err)
	}

	// Test invalid transition
	if err := sm.Transition(StateProcessing); err == nil {
		t.Error("Expected error for invalid Idle->Processing transition")
	}
}

// TestSubscribe tests the subscription mechanism
func TestSubscribe(t *testing.T) {
	sm := NewStateMachine()

	var called bool
	sm.Subscribe(func(oldState, newState State) {
		called = true
		if oldState != StateIdle || newState != StateRecording {
			t.Errorf("Expected Idle->Recording, got %s->%s", oldState, newState)
		}
	})

	sm.Transition(StateRecording)

	if !called {
		t.Error("Callback was not called")
	}
}
