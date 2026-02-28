package indicator

import (
	"path/filepath"
	"testing"

	"github.com/d-mozulyov/vox/internal/state"
)

// TestAudioIndicator_Integration tests the audio indicator with real sound files
func TestAudioIndicator_Integration(t *testing.T) {
	// Use actual sounds directory
	soundsPath := filepath.Join("..", "..", "assets", "sounds")

	ai, err := NewAudioIndicator(soundsPath)
	if err != nil {
		t.Fatalf("NewAudioIndicator failed: %v", err)
	}

	// Test all valid transitions
	transitions := []struct {
		name string
		from state.State
		to   state.State
	}{
		{"Start recording", state.StateIdle, state.StateRecording},
		{"Stop recording", state.StateRecording, state.StateProcessing},
		{"Processing done", state.StateProcessing, state.StateIdle},
	}

	for _, tt := range transitions {
		t.Run(tt.name, func(t *testing.T) {
			err := ai.PlaySound(tt.from, tt.to)
			if err != nil {
				t.Errorf("PlaySound(%v, %v) failed: %v", tt.from, tt.to, err)
			}
		})
	}
}
