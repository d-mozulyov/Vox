package indicator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/d-mozulyov/vox/internal/platform"
	"github.com/d-mozulyov/vox/internal/state"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// AudioIndicator defines the interface for audio state indication
type AudioIndicator interface {
	// PlaySound plays an audio feedback for state transition
	PlaySound(from, to state.State) error
}

// audioIndicator implements the AudioIndicator interface
type audioIndicator struct {
	soundsPath string
}

// NewAudioIndicator creates a new audio indicator instance
// soundsPath is the directory containing sound files
func NewAudioIndicator(soundsPath string) (AudioIndicator, error) {
	logger := platform.GetLogger()

	// Initialize speaker once with standard settings
	// 44100 Hz sample rate, 4096 buffer size
	speaker.Init(beep.SampleRate(44100), 4096)

	logger.Info("Audio indicator created successfully, sounds path: %s", soundsPath)

	return &audioIndicator{
		soundsPath: soundsPath,
	}, nil
}

// PlaySound plays an audio feedback for the given state transition
func (ai *audioIndicator) PlaySound(from, to state.State) error {
	logger := platform.GetLogger()

	// Determine which sound file to play
	filename := ai.getSoundFilename(from, to)
	if filename == "" {
		// No sound for this transition
		return nil
	}

	// Load and play the sound
	soundPath := filepath.Join(ai.soundsPath, filename)
	if err := ai.playWavFile(soundPath); err != nil {
		// Log warning but don't return error - audio is non-critical
		logger.Warn("Failed to play audio feedback: %v", err)
	}

	return nil
}

// getSoundFilename returns the sound filename for a state transition
func (ai *audioIndicator) getSoundFilename(from, to state.State) string {
	if from == state.StateIdle && to == state.StateRecording {
		return "start_recording.wav"
	}
	if from == state.StateRecording && to == state.StateProcessing {
		return "stop_recording.wav"
	}
	if from == state.StateProcessing && to == state.StateIdle {
		return "processing_done.wav"
	}
	return ""
}

// playWavFile loads and plays a WAV file
func (ai *audioIndicator) playWavFile(path string) error {
	logger := platform.GetLogger()

	file, err := os.Open(path)
	if err != nil {
		logger.Error("Failed to open audio file %s: %v", path, err)
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	streamer, format, err := wav.Decode(file)
	if err != nil {
		logger.Error("Failed to decode WAV file %s: %v", path, err)
		return fmt.Errorf("failed to decode WAV file: %w", err)
	}
	defer streamer.Close()

	// Resample if needed to match speaker sample rate
	resampled := beep.Resample(4, format.SampleRate, beep.SampleRate(44100), streamer)

	// Play sound asynchronously
	done := make(chan bool)
	speaker.Play(beep.Seq(resampled, beep.Callback(func() {
		done <- true
	})))

	// Wait for completion
	<-done

	logger.Info("Audio feedback played successfully: %s", path)

	return nil
}
