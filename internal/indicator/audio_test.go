package indicator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/d-mozulyov/vox/internal/state"
)

func TestNewAudioIndicator(t *testing.T) {
	// Create temporary directory for test sounds
	tempDir := t.TempDir()

	// Test with non-existent directory - should succeed (sounds loaded on demand)
	ai, err := NewAudioIndicator(tempDir)
	if err != nil {
		t.Fatalf("NewAudioIndicator failed: %v", err)
	}
	if ai == nil {
		t.Fatal("Expected non-nil AudioIndicator")
	}
}

func TestAudioIndicator_GetSoundFilename(t *testing.T) {
	tempDir := t.TempDir()
	ai, _ := NewAudioIndicator(tempDir)
	impl := ai.(*audioIndicator)

	tests := []struct {
		name     string
		from     state.State
		to       state.State
		expected string
	}{
		{
			name:     "Idle to Recording",
			from:     state.StateIdle,
			to:       state.StateRecording,
			expected: "start_recording.wav",
		},
		{
			name:     "Recording to Processing",
			from:     state.StateRecording,
			to:       state.StateProcessing,
			expected: "stop_recording.wav",
		},
		{
			name:     "Processing to Idle",
			from:     state.StateProcessing,
			to:       state.StateIdle,
			expected: "processing_done.wav",
		},
		{
			name:     "Invalid transition",
			from:     state.StateIdle,
			to:       state.StateProcessing,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := impl.getSoundFilename(tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("getSoundFilename(%v, %v) = %q, want %q",
					tt.from, tt.to, result, tt.expected)
			}
		})
	}
}

func TestAudioIndicator_PlaySound_MissingFile(t *testing.T) {
	// Create temporary directory without sound files
	tempDir := t.TempDir()

	ai, err := NewAudioIndicator(tempDir)
	if err != nil {
		t.Fatalf("NewAudioIndicator failed: %v", err)
	}

	// Playing sound with missing file should not return error
	// (error is logged but not returned)
	err = ai.PlaySound(state.StateIdle, state.StateRecording)
	if err != nil {
		t.Errorf("PlaySound should not return error for missing file, got: %v", err)
	}
}

func TestAudioIndicator_PlaySound_NoSoundForTransition(t *testing.T) {
	tempDir := t.TempDir()

	ai, err := NewAudioIndicator(tempDir)
	if err != nil {
		t.Fatalf("NewAudioIndicator failed: %v", err)
	}

	// Invalid transition should not return error
	err = ai.PlaySound(state.StateIdle, state.StateProcessing)
	if err != nil {
		t.Errorf("PlaySound should not return error for invalid transition, got: %v", err)
	}
}

// createMinimalWavFile creates a minimal valid WAV file for testing
func createMinimalWavFile(path string) error {
	// Minimal WAV file header (44 bytes) + 1 sample
	// RIFF header
	wavData := []byte{
		// "RIFF" chunk descriptor
		0x52, 0x49, 0x46, 0x46, // "RIFF"
		0x24, 0x00, 0x00, 0x00, // File size - 8 (36 bytes)
		0x57, 0x41, 0x56, 0x45, // "WAVE"

		// "fmt " sub-chunk
		0x66, 0x6d, 0x74, 0x20, // "fmt "
		0x10, 0x00, 0x00, 0x00, // Subchunk size (16 bytes)
		0x01, 0x00, // Audio format (1 = PCM)
		0x01, 0x00, // Number of channels (1 = mono)
		0x44, 0xac, 0x00, 0x00, // Sample rate (44100 Hz)
		0x88, 0x58, 0x01, 0x00, // Byte rate (44100 * 1 * 2)
		0x02, 0x00, // Block align (1 * 2)
		0x10, 0x00, // Bits per sample (16)

		// "data" sub-chunk
		0x64, 0x61, 0x74, 0x61, // "data"
		0x00, 0x00, 0x00, 0x00, // Subchunk size (0 bytes of audio data)
	}

	return os.WriteFile(path, wavData, 0644)
}

func TestAudioIndicator_PlaySound_ValidFile(t *testing.T) {
	// Create temporary directory with valid WAV file
	tempDir := t.TempDir()

	// Create a minimal WAV file
	wavPath := filepath.Join(tempDir, "start_recording.wav")
	if err := createMinimalWavFile(wavPath); err != nil {
		t.Fatalf("Failed to create test WAV file: %v", err)
	}

	ai, err := NewAudioIndicator(tempDir)
	if err != nil {
		t.Fatalf("NewAudioIndicator failed: %v", err)
	}

	// Playing sound with valid file should not return error
	err = ai.PlaySound(state.StateIdle, state.StateRecording)
	if err != nil {
		t.Errorf("PlaySound failed with valid file: %v", err)
	}
}
