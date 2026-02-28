package state

import (
	"sync"
	"testing"
	"time"
)

// TestNewStateMachine verifies that a new state machine starts in Idle state
func TestNewStateMachine(t *testing.T) {
	sm := NewStateMachine()
	if sm.GetState() != StateIdle {
		t.Errorf("Expected initial state to be Idle, got %s", sm.GetState())
	}
}

// TestValidTransition_IdleToRecording tests the valid transition from Idle to Recording
func TestValidTransition_IdleToRecording(t *testing.T) {
	sm := NewStateMachine()
	err := sm.Transition(StateRecording)
	if err != nil {
		t.Errorf("Expected no error for Idle->Recording transition, got: %v", err)
	}
	if sm.GetState() != StateRecording {
		t.Errorf("Expected state to be Recording, got %s", sm.GetState())
	}
}

// TestValidTransition_RecordingToProcessing tests the valid transition from Recording to Processing
func TestValidTransition_RecordingToProcessing(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(StateRecording)
	err := sm.Transition(StateProcessing)
	if err != nil {
		t.Errorf("Expected no error for Recording->Processing transition, got: %v", err)
	}
	if sm.GetState() != StateProcessing {
		t.Errorf("Expected state to be Processing, got %s", sm.GetState())
	}
}

// TestValidTransition_ProcessingToIdle tests the valid transition from Processing to Idle
func TestValidTransition_ProcessingToIdle(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(StateRecording)
	sm.Transition(StateProcessing)
	err := sm.Transition(StateIdle)
	if err != nil {
		t.Errorf("Expected no error for Processing->Idle transition, got: %v", err)
	}
	if sm.GetState() != StateIdle {
		t.Errorf("Expected state to be Idle, got %s", sm.GetState())
	}
}

// TestInvalidTransition_RecordingToIdle tests that Recording->Idle is invalid
func TestInvalidTransition_RecordingToIdle(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(StateRecording)
	err := sm.Transition(StateIdle)
	if err == nil {
		t.Error("Expected error for invalid Recording->Idle transition, got nil")
	}
	if sm.GetState() != StateRecording {
		t.Errorf("Expected state to remain Recording after invalid transition, got %s", sm.GetState())
	}
}

// TestInvalidTransition_IdleToProcessing tests that Idle->Processing is invalid
func TestInvalidTransition_IdleToProcessing(t *testing.T) {
	sm := NewStateMachine()
	err := sm.Transition(StateProcessing)
	if err == nil {
		t.Error("Expected error for invalid Idle->Processing transition, got nil")
	}
	if sm.GetState() != StateIdle {
		t.Errorf("Expected state to remain Idle after invalid transition, got %s", sm.GetState())
	}
}

// TestInvalidTransition_ProcessingToRecording tests that Processing->Recording is invalid
func TestInvalidTransition_ProcessingToRecording(t *testing.T) {
	sm := NewStateMachine()
	sm.Transition(StateRecording)
	sm.Transition(StateProcessing)
	err := sm.Transition(StateRecording)
	if err == nil {
		t.Error("Expected error for invalid Processing->Recording transition, got nil")
	}
	if sm.GetState() != StateProcessing {
		t.Errorf("Expected state to remain Processing after invalid transition, got %s", sm.GetState())
	}
}

// TestSubscribe_SingleCallback tests that a single callback is called on state change
func TestSubscribe_SingleCallback(t *testing.T) {
	sm := NewStateMachine()

	var callbackCalled bool
	var receivedOldState, receivedNewState State

	sm.Subscribe(func(oldState, newState State) {
		callbackCalled = true
		receivedOldState = oldState
		receivedNewState = newState
	})

	sm.Transition(StateRecording)

	// Give callback time to execute
	time.Sleep(10 * time.Millisecond)

	if !callbackCalled {
		t.Error("Expected callback to be called, but it wasn't")
	}
	if receivedOldState != StateIdle {
		t.Errorf("Expected old state to be Idle, got %s", receivedOldState)
	}
	if receivedNewState != StateRecording {
		t.Errorf("Expected new state to be Recording, got %s", receivedNewState)
	}
}

// TestSubscribe_MultipleCallbacks tests that multiple callbacks are called on state change
func TestSubscribe_MultipleCallbacks(t *testing.T) {
	sm := NewStateMachine()

	var callback1Called, callback2Called bool

	sm.Subscribe(func(oldState, newState State) {
		callback1Called = true
	})

	sm.Subscribe(func(oldState, newState State) {
		callback2Called = true
	})

	sm.Transition(StateRecording)

	// Give callbacks time to execute
	time.Sleep(10 * time.Millisecond)

	if !callback1Called {
		t.Error("Expected callback 1 to be called, but it wasn't")
	}
	if !callback2Called {
		t.Error("Expected callback 2 to be called, but it wasn't")
	}
}

// TestSubscribe_NoCallbackOnInvalidTransition tests that callbacks are not called on invalid transitions
func TestSubscribe_NoCallbackOnInvalidTransition(t *testing.T) {
	sm := NewStateMachine()

	var callbackCalled bool

	sm.Subscribe(func(oldState, newState State) {
		callbackCalled = true
	})

	sm.Transition(StateProcessing) // Invalid transition

	// Give callback time to execute (if it would)
	time.Sleep(10 * time.Millisecond)

	if callbackCalled {
		t.Error("Expected callback not to be called on invalid transition, but it was")
	}
}

// TestConcurrentTransitions tests that concurrent transitions are handled safely
func TestConcurrentTransitions(t *testing.T) {
	sm := NewStateMachine()

	var wg sync.WaitGroup
	errors := make(chan error, 100)

	// Try to transition from Idle to Recording concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := sm.Transition(StateRecording)
			if err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	// Only one transition should succeed, the rest should fail
	errorCount := 0
	for range errors {
		errorCount++
	}

	if errorCount != 9 {
		t.Errorf("Expected 9 errors from concurrent transitions, got %d", errorCount)
	}

	if sm.GetState() != StateRecording {
		t.Errorf("Expected final state to be Recording, got %s", sm.GetState())
	}
}

// TestStateString tests the String() method of State
func TestStateString(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{StateIdle, "Idle"},
		{StateRecording, "Recording"},
		{StateProcessing, "Processing"},
		{State(999), "Unknown"},
	}

	for _, tt := range tests {
		if got := tt.state.String(); got != tt.expected {
			t.Errorf("State(%d).String() = %s, want %s", tt.state, got, tt.expected)
		}
	}
}
