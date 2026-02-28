# Audio Feedback Sounds

This directory contains audio feedback sounds for state transitions.

## Required Sounds

- `start_recording.wav` - Sound when recording starts (Idle → Recording)
- `stop_recording.wav` - Sound when recording stops (Recording → Processing)
- `processing_done.wav` - Sound when processing completes (Processing → Idle)

## Sound Specifications

- **Format**: WAV (for minimal dependencies and cross-platform compatibility)
- **Duration**: Maximum 300 milliseconds
- **Characteristics**: Pleasant, non-intrusive sounds (e.g., soft clicks or beeps)
- **Sample rate**: 44.1 kHz or 48 kHz
- **Bit depth**: 16-bit

## Implementation Status

✅ All sound files have been created and meet the specifications:
- `start_recording.wav`: 100ms duration, WAV format, 44.1kHz, 16-bit
- `stop_recording.wav`: 100ms duration, WAV format, 44.1kHz, 16-bit  
- `processing_done.wav`: 150ms duration, WAV format, 44.1kHz, 16-bit

All files are under the 300ms requirement and provide pleasant, non-intrusive audio feedback.
