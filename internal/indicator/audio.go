package indicator

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/d-mozulyov/vox/internal/platform"
	"github.com/d-mozulyov/vox/internal/state"
	"github.com/ebitengine/oto/v3"
)

// AudioIndicator defines the interface for audio state indication
type AudioIndicator interface {
	// PlaySound plays an audio feedback for state transition
	PlaySound(from, to state.State) error
}

// audioIndicator implements the AudioIndicator interface
type audioIndicator struct {
	soundsPath string
	context    *oto.Context
}

// NewAudioIndicator creates a new audio indicator instance
// soundsPath is the directory containing sound files
func NewAudioIndicator(soundsPath string) (AudioIndicator, error) {
	logger := platform.GetLogger()

	// Initialize oto context with standard settings
	// 44100 Hz sample rate, 2 channels (stereo), 16-bit samples
	op := &oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}

	ctx, readyChan, err := oto.NewContext(op)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize audio context: %w", err)
	}

	// Wait for the audio context to be ready
	<-readyChan

	logger.Info("Audio indicator created successfully, sounds path: %s", soundsPath)

	return &audioIndicator{
		soundsPath: soundsPath,
		context:    ctx,
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
	// Play start_recording.wav when starting recording (Idle -> Recording)
	if from == state.StateIdle && to == state.StateRecording {
		return "start_recording.wav"
	}
	// Play stop_recording.wav when stopping recording (Recording -> Idle)
	if from == state.StateRecording && to == state.StateIdle {
		return "stop_recording.wav"
	}
	return ""
}

// wavHeader represents a WAV file header
type wavHeader struct {
	ChunkID       [4]byte
	ChunkSize     uint32
	Format        [4]byte
	Subchunk1ID   [4]byte
	Subchunk1Size uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
	Subchunk2ID   [4]byte
	Subchunk2Size uint32
}

// playWavFile loads and plays a WAV file
func (ai *audioIndicator) playWavFile(path string) error {
	logger := platform.GetLogger()

	// Read the entire WAV file
	data, err := os.ReadFile(path)
	if err != nil {
		logger.Error("Failed to read audio file %s: %v", path, err)
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse WAV header
	if len(data) < 44 {
		return fmt.Errorf("file too small to be a valid WAV file")
	}

	reader := bytes.NewReader(data)
	var header wavHeader
	if err := binary.Read(reader, binary.LittleEndian, &header); err != nil {
		return fmt.Errorf("failed to parse WAV header: %w", err)
	}

	// Validate WAV format
	if string(header.ChunkID[:]) != "RIFF" || string(header.Format[:]) != "WAVE" {
		return fmt.Errorf("not a valid WAV file")
	}

	// Get audio data (skip header)
	audioData := data[44:]

	// Convert to stereo if mono
	if header.NumChannels == 1 && header.BitsPerSample == 16 {
		audioData = ai.monoToStereo(audioData)
	}

	// Create a player and play the sound
	player := ai.context.NewPlayer(bytes.NewReader(audioData))
	player.Play()

	// Calculate duration and wait for playback to complete
	duration := time.Duration(float64(len(audioData)) / float64(header.ByteRate) * float64(time.Second))
	time.Sleep(duration)

	logger.Info("Audio feedback played successfully: %s", path)

	return nil
}

// monoToStereo converts mono 16-bit PCM to stereo by duplicating each sample
func (ai *audioIndicator) monoToStereo(mono []byte) []byte {
	stereo := make([]byte, len(mono)*2)
	for i := 0; i < len(mono); i += 2 {
		// Copy left channel
		stereo[i*2] = mono[i]
		stereo[i*2+1] = mono[i+1]
		// Copy right channel (same as left)
		stereo[i*2+2] = mono[i]
		stereo[i*2+3] = mono[i+1]
	}
	return stereo
}
