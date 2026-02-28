//go:build ignore
// +build ignore

package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
)

// generateBeepWav generates a simple beep sound WAV file
func generateBeepWav(filename string, frequency float64, durationMs int) error {
	sampleRate := 44100
	numSamples := (sampleRate * durationMs) / 1000

	// Create WAV file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Calculate sizes
	dataSize := numSamples * 2 // 16-bit samples
	fileSize := 36 + dataSize

	// Write RIFF header
	file.WriteString("RIFF")
	binary.Write(file, binary.LittleEndian, uint32(fileSize))
	file.WriteString("WAVE")

	// Write fmt chunk
	file.WriteString("fmt ")
	binary.Write(file, binary.LittleEndian, uint32(16))        // Chunk size
	binary.Write(file, binary.LittleEndian, uint16(1))         // Audio format (PCM)
	binary.Write(file, binary.LittleEndian, uint16(1))         // Num channels (mono)
	binary.Write(file, binary.LittleEndian, uint32(sampleRate)) // Sample rate
	binary.Write(file, binary.LittleEndian, uint32(sampleRate*2)) // Byte rate
	binary.Write(file, binary.LittleEndian, uint16(2))         // Block align
	binary.Write(file, binary.LittleEndian, uint16(16))        // Bits per sample

	// Write data chunk
	file.WriteString("data")
	binary.Write(file, binary.LittleEndian, uint32(dataSize))

	// Generate sine wave samples with fade in/out
	for i := 0; i < numSamples; i++ {
		t := float64(i) / float64(sampleRate)

		// Generate sine wave
		sample := math.Sin(2 * math.Pi * frequency * t)

		// Apply fade in/out envelope to avoid clicks
		fadeLength := numSamples / 10 // 10% fade
		if i < fadeLength {
			sample *= float64(i) / float64(fadeLength)
		} else if i > numSamples-fadeLength {
			sample *= float64(numSamples-i) / float64(fadeLength)
		}

		// Convert to 16-bit integer
		sample16 := int16(sample * 16000) // Reduced amplitude for pleasant sound
		binary.Write(file, binary.LittleEndian, sample16)
	}

	return nil
}

func main() {
	// Create assets/sounds directory if it doesn't exist
	soundsDir := filepath.Join("assets", "sounds")
	if err := os.MkdirAll(soundsDir, 0755); err != nil {
		log.Fatalf("Failed to create sounds directory: %v", err)
	}

	// Generate placeholder sounds
	sounds := []struct {
		filename  string
		frequency float64
		duration  int
	}{
		{"start_recording.wav", 800.0, 100},  // Higher pitch, short beep
		{"stop_recording.wav", 600.0, 100},   // Lower pitch, short beep
		{"processing_done.wav", 700.0, 150},  // Medium pitch, slightly longer
	}

	for _, sound := range sounds {
		path := filepath.Join(soundsDir, sound.filename)
		fmt.Printf("Generating %s...\n", sound.filename)
		if err := generateBeepWav(path, sound.frequency, sound.duration); err != nil {
			log.Fatalf("Failed to generate %s: %v", sound.filename, err)
		}
	}

	fmt.Println("All placeholder sounds generated successfully!")
}
