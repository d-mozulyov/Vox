package state

// State represents the current state of the application
type State int

const (
	// StateIdle represents the idle state - application is waiting for user action
	StateIdle State = iota
	// StateRecording represents the recording state - voice recording is in progress
	StateRecording
)

// String returns the string representation of the state
func (s State) String() string {
	switch s {
	case StateIdle:
		return "Idle"
	case StateRecording:
		return "Recording"
	default:
		return "Unknown"
	}
}

// StateMachine defines the interface for managing application state
type StateMachine interface {
	// GetState returns the current application state
	GetState() State

	// Transition attempts to transition to a new state
	// Returns error if the transition is invalid
	Transition(newState State) error

	// Subscribe registers a callback for state changes
	// The callback receives the old state and the new state
	Subscribe(callback func(oldState, newState State))
}
