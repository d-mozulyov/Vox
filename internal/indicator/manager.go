package indicator

import (
	"sync"

	"github.com/d-mozulyov/vox/internal/platform"
	"github.com/d-mozulyov/vox/internal/state"
)

// IndicatorManager defines the interface for coordinating visual and audio indicators
type IndicatorManager interface {
	// OnStateChange handles state transitions and triggers indicators
	OnStateChange(oldState, newState state.State)

	// SetVisualIndicator sets the visual indicator implementation
	SetVisualIndicator(indicator VisualIndicator)

	// SetAudioIndicator sets the audio indicator implementation
	SetAudioIndicator(indicator AudioIndicator)
}

// indicatorManager implements the IndicatorManager interface
type indicatorManager struct {
	visualIndicator VisualIndicator
	audioIndicator  AudioIndicator
	mutex           sync.RWMutex
}

// NewIndicatorManager creates a new indicator manager instance
func NewIndicatorManager() IndicatorManager {
	return &indicatorManager{}
}

// SetVisualIndicator sets the visual indicator implementation
func (im *indicatorManager) SetVisualIndicator(indicator VisualIndicator) {
	im.mutex.Lock()
	defer im.mutex.Unlock()
	im.visualIndicator = indicator
}

// SetAudioIndicator sets the audio indicator implementation
func (im *indicatorManager) SetAudioIndicator(indicator AudioIndicator) {
	im.mutex.Lock()
	defer im.mutex.Unlock()
	im.audioIndicator = indicator
}

// OnStateChange handles state transitions by coordinating visual and audio indicators
// Both indicators are triggered in parallel for responsiveness
func (im *indicatorManager) OnStateChange(oldState, newState state.State) {
	im.mutex.RLock()
	visual := im.visualIndicator
	audio := im.audioIndicator
	im.mutex.RUnlock()

	logger := platform.GetLogger()

	// Use WaitGroup to execute indicators in parallel
	var wg sync.WaitGroup

	// Trigger visual indicator
	if visual != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := visual.UpdateIcon(newState); err != nil {
				logger.Warn("Failed to update visual indicator: %v", err)
			}
		}()
	}

	// Trigger audio indicator
	if audio != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := audio.PlaySound(oldState, newState); err != nil {
				logger.Warn("Failed to play audio indicator: %v", err)
			}
		}()
	}

	// Wait for both indicators to complete
	wg.Wait()
}
