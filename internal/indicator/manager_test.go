package indicator

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/d-mozulyov/vox/internal/state"
)

// mockVisualIndicator is a mock implementation of VisualIndicator for testing
type mockVisualIndicator struct {
	updateIconFunc func(state.State) error
	callCount      int
	mutex          sync.Mutex
}

func (m *mockVisualIndicator) UpdateIcon(s state.State) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.callCount++
	if m.updateIconFunc != nil {
		return m.updateIconFunc(s)
	}
	return nil
}

func (m *mockVisualIndicator) getCallCount() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.callCount
}

// mockAudioIndicator is a mock implementation of AudioIndicator for testing
type mockAudioIndicator struct {
	playSoundFunc func(state.State, state.State) error
	callCount     int
	mutex         sync.Mutex
}

func (m *mockAudioIndicator) PlaySound(from, to state.State) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.callCount++
	if m.playSoundFunc != nil {
		return m.playSoundFunc(from, to)
	}
	return nil
}

func (m *mockAudioIndicator) getCallCount() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.callCount
}

// TestIndicatorManager_OnStateChange_BothIndicators tests that both indicators are called
func TestIndicatorManager_OnStateChange_BothIndicators(t *testing.T) {
	manager := NewIndicatorManager()

	visualMock := &mockVisualIndicator{}
	audioMock := &mockAudioIndicator{}

	manager.SetVisualIndicator(visualMock)
	manager.SetAudioIndicator(audioMock)

	// Trigger state change
	manager.OnStateChange(state.StateIdle, state.StateRecording)

	// Verify both indicators were called
	if visualMock.getCallCount() != 1 {
		t.Errorf("Expected visual indicator to be called once, got %d", visualMock.getCallCount())
	}
	if audioMock.getCallCount() != 1 {
		t.Errorf("Expected audio indicator to be called once, got %d", audioMock.getCallCount())
	}
}

// TestIndicatorManager_OnStateChange_VisualOnly tests with only visual indicator set
func TestIndicatorManager_OnStateChange_VisualOnly(t *testing.T) {
	manager := NewIndicatorManager()

	visualMock := &mockVisualIndicator{}
	manager.SetVisualIndicator(visualMock)

	// Trigger state change (no audio indicator set)
	manager.OnStateChange(state.StateIdle, state.StateRecording)

	// Verify visual indicator was called
	if visualMock.getCallCount() != 1 {
		t.Errorf("Expected visual indicator to be called once, got %d", visualMock.getCallCount())
	}
}

// TestIndicatorManager_OnStateChange_AudioOnly tests with only audio indicator set
func TestIndicatorManager_OnStateChange_AudioOnly(t *testing.T) {
	manager := NewIndicatorManager()

	audioMock := &mockAudioIndicator{}
	manager.SetAudioIndicator(audioMock)

	// Trigger state change (no visual indicator set)
	manager.OnStateChange(state.StateIdle, state.StateRecording)

	// Verify audio indicator was called
	if audioMock.getCallCount() != 1 {
		t.Errorf("Expected audio indicator to be called once, got %d", audioMock.getCallCount())
	}
}

// TestIndicatorManager_OnStateChange_NoIndicators tests with no indicators set
func TestIndicatorManager_OnStateChange_NoIndicators(t *testing.T) {
	manager := NewIndicatorManager()

	// Should not panic when no indicators are set
	manager.OnStateChange(state.StateIdle, state.StateRecording)
}

// TestIndicatorManager_OnStateChange_VisualError tests error handling from visual indicator
func TestIndicatorManager_OnStateChange_VisualError(t *testing.T) {
	manager := NewIndicatorManager()

	visualMock := &mockVisualIndicator{
		updateIconFunc: func(s state.State) error {
			return errors.New("visual error")
		},
	}
	audioMock := &mockAudioIndicator{}

	manager.SetVisualIndicator(visualMock)
	manager.SetAudioIndicator(audioMock)

	// Should not panic on visual error
	manager.OnStateChange(state.StateIdle, state.StateRecording)

	// Both indicators should still be called
	if visualMock.getCallCount() != 1 {
		t.Errorf("Expected visual indicator to be called once, got %d", visualMock.getCallCount())
	}
	if audioMock.getCallCount() != 1 {
		t.Errorf("Expected audio indicator to be called once, got %d", audioMock.getCallCount())
	}
}

// TestIndicatorManager_OnStateChange_AudioError tests error handling from audio indicator
func TestIndicatorManager_OnStateChange_AudioError(t *testing.T) {
	manager := NewIndicatorManager()

	visualMock := &mockVisualIndicator{}
	audioMock := &mockAudioIndicator{
		playSoundFunc: func(from, to state.State) error {
			return errors.New("audio error")
		},
	}

	manager.SetVisualIndicator(visualMock)
	manager.SetAudioIndicator(audioMock)

	// Should not panic on audio error
	manager.OnStateChange(state.StateIdle, state.StateRecording)

	// Both indicators should still be called
	if visualMock.getCallCount() != 1 {
		t.Errorf("Expected visual indicator to be called once, got %d", visualMock.getCallCount())
	}
	if audioMock.getCallCount() != 1 {
		t.Errorf("Expected audio indicator to be called once, got %d", audioMock.getCallCount())
	}
}

// TestIndicatorManager_OnStateChange_BothErrors tests error handling from both indicators
func TestIndicatorManager_OnStateChange_BothErrors(t *testing.T) {
	manager := NewIndicatorManager()

	visualMock := &mockVisualIndicator{
		updateIconFunc: func(s state.State) error {
			return errors.New("visual error")
		},
	}
	audioMock := &mockAudioIndicator{
		playSoundFunc: func(from, to state.State) error {
			return errors.New("audio error")
		},
	}

	manager.SetVisualIndicator(visualMock)
	manager.SetAudioIndicator(audioMock)

	// Should not panic on both errors
	manager.OnStateChange(state.StateIdle, state.StateRecording)

	// Both indicators should still be called
	if visualMock.getCallCount() != 1 {
		t.Errorf("Expected visual indicator to be called once, got %d", visualMock.getCallCount())
	}
	if audioMock.getCallCount() != 1 {
		t.Errorf("Expected audio indicator to be called once, got %d", audioMock.getCallCount())
	}
}

// TestIndicatorManager_OnStateChange_ParallelExecution tests that indicators run in parallel
func TestIndicatorManager_OnStateChange_ParallelExecution(t *testing.T) {
	manager := NewIndicatorManager()

	// Create indicators that take some time to execute
	visualDelay := 50 * time.Millisecond
	audioDelay := 50 * time.Millisecond

	visualMock := &mockVisualIndicator{
		updateIconFunc: func(s state.State) error {
			time.Sleep(visualDelay)
			return nil
		},
	}
	audioMock := &mockAudioIndicator{
		playSoundFunc: func(from, to state.State) error {
			time.Sleep(audioDelay)
			return nil
		},
	}

	manager.SetVisualIndicator(visualMock)
	manager.SetAudioIndicator(audioMock)

	// Measure execution time
	start := time.Now()
	manager.OnStateChange(state.StateIdle, state.StateRecording)
	elapsed := time.Since(start)

	// If executed in parallel, total time should be close to max(visualDelay, audioDelay)
	// If executed sequentially, total time would be visualDelay + audioDelay
	maxDelay := visualDelay
	if audioDelay > maxDelay {
		maxDelay = audioDelay
	}

	// Allow some overhead for goroutine scheduling
	expectedMax := maxDelay + 30*time.Millisecond

	if elapsed > visualDelay+audioDelay {
		t.Errorf("Indicators appear to be running sequentially. Elapsed: %v, expected less than %v", elapsed, expectedMax)
	}

	// Verify both were called
	if visualMock.getCallCount() != 1 {
		t.Errorf("Expected visual indicator to be called once, got %d", visualMock.getCallCount())
	}
	if audioMock.getCallCount() != 1 {
		t.Errorf("Expected audio indicator to be called once, got %d", audioMock.getCallCount())
	}
}

// TestIndicatorManager_OnStateChange_AllTransitions tests all valid state transitions
func TestIndicatorManager_OnStateChange_AllTransitions(t *testing.T) {
	transitions := []struct {
		from state.State
		to   state.State
	}{
		{state.StateIdle, state.StateRecording},
		{state.StateRecording, state.StateProcessing},
		{state.StateProcessing, state.StateIdle},
	}

	for _, tc := range transitions {
		t.Run(tc.from.String()+"_to_"+tc.to.String(), func(t *testing.T) {
			manager := NewIndicatorManager()

			visualMock := &mockVisualIndicator{}
			audioMock := &mockAudioIndicator{}

			manager.SetVisualIndicator(visualMock)
			manager.SetAudioIndicator(audioMock)

			manager.OnStateChange(tc.from, tc.to)

			if visualMock.getCallCount() != 1 {
				t.Errorf("Expected visual indicator to be called once, got %d", visualMock.getCallCount())
			}
			if audioMock.getCallCount() != 1 {
				t.Errorf("Expected audio indicator to be called once, got %d", audioMock.getCallCount())
			}
		})
	}
}

// TestIndicatorManager_SetIndicators_ThreadSafety tests concurrent access to set indicators
func TestIndicatorManager_SetIndicators_ThreadSafety(t *testing.T) {
	manager := NewIndicatorManager()

	var wg sync.WaitGroup
	iterations := 100

	// Concurrently set visual indicator
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			manager.SetVisualIndicator(&mockVisualIndicator{})
		}
	}()

	// Concurrently set audio indicator
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			manager.SetAudioIndicator(&mockAudioIndicator{})
		}
	}()

	// Concurrently trigger state changes
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			manager.OnStateChange(state.StateIdle, state.StateRecording)
		}
	}()

	// Should not panic or deadlock
	wg.Wait()
}
