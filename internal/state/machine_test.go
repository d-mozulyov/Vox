package state

import "testing"

// TestStateMachine tests basic state machine functionality
func TestStateMachine(t *testing.T) {
	sm := NewStateMachine()

	// Test initial state
	if sm.GetState() != StateIdle {
		t.Errorf("Expected initial state Idle, got %s", sm.GetState())
	}

	// Test valid transitions: Idle -> Recording -> Idle
	if err := sm.Transition(StateRecording); err != nil {
		t.Errorf("Idle->Recording failed: %v", err)
	}
	if err := sm.Transition(StateIdle); err != nil {
		t.Errorf("Recording->Idle failed: %v", err)
	}

	// Test invalid transition: Idle -> Idle
	if err := sm.Transition(StateIdle); err == nil {
		t.Error("Expected error for invalid Idle->Idle transition")
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

	if err := sm.Transition(StateRecording); err != nil {
		t.Fatalf("Failed to transition to Recording: %v", err)
	}

	if !called {
		t.Error("Callback was not called")
	}
}
