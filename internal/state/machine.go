package state

import (
	"fmt"
	"sync"

	"github.com/d-mozulyov/vox/internal/platform"
)

// stateMachine is the concrete implementation of StateMachine interface
type stateMachine struct {
	current   State
	mutex     sync.RWMutex
	callbacks []func(oldState, newState State)
}

// NewStateMachine creates a new state machine initialized to StateIdle
func NewStateMachine() StateMachine {
	return &stateMachine{
		current:   StateIdle,
		callbacks: make([]func(oldState, newState State), 0),
	}
}

// GetState returns the current application state
func (sm *stateMachine) GetState() State {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return sm.current
}

// Transition attempts to transition to a new state
// Returns error if the transition is invalid
func (sm *stateMachine) Transition(newState State) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	logger := platform.GetLogger()

	// Validate the transition
	if !sm.isValidTransition(sm.current, newState) {
		err := fmt.Errorf("invalid state transition: %s -> %s", sm.current, newState)
		logger.Error("Invalid state transition attempted: %s -> %s", sm.current, newState)
		return err
	}

	oldState := sm.current
	sm.current = newState

	logger.Info("State transition: %s -> %s", oldState, newState)

	// Notify all subscribers (outside the lock to avoid deadlocks)
	sm.mutex.Unlock()
	sm.notifySubscribers(oldState, newState)
	sm.mutex.Lock()

	return nil
}

// Subscribe registers a callback for state changes
func (sm *stateMachine) Subscribe(callback func(oldState, newState State)) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	sm.callbacks = append(sm.callbacks, callback)
}

// isValidTransition checks if a state transition is valid
// Valid transitions:
// - Idle -> Recording (start recording via hotkey)
// - Recording -> Idle (stop recording via hotkey)
func (sm *stateMachine) isValidTransition(from, to State) bool {
	switch from {
	case StateIdle:
		return to == StateRecording
	case StateRecording:
		return to == StateIdle
	default:
		return false
	}
}

// notifySubscribers calls all registered callbacks with the state change
func (sm *stateMachine) notifySubscribers(oldState, newState State) {
	sm.mutex.RLock()
	callbacks := make([]func(oldState, newState State), len(sm.callbacks))
	copy(callbacks, sm.callbacks)
	sm.mutex.RUnlock()

	for _, callback := range callbacks {
		callback(oldState, newState)
	}
}
