# Vox Project Structure

This document describes the organization of the Vox project codebase.

## Directory Layout

```
vox/
├── cmd/                    # Application entry points
│   └── vox/               # Main application
│       └── main.go        # Application entry point
│
├── internal/              # Private application code
│   ├── state/            # State machine implementation
│   ├── tray/             # System tray manager
│   ├── hotkey/           # Global hotkey manager
│   ├── indicator/        # Visual and audio indicators
│   ├── audio/            # Audio playback functionality
│   └── platform/         # Platform-specific code and logging
│
├── pkg/                   # Public library code
│   └── config/           # Configuration structures
│
├── assets/               # Application assets
│   ├── icons/           # System tray icons (idle, recording, processing)
│   └── sounds/          # Audio feedback files
│
├── .kiro/               # Kiro IDE configuration and specs
│   ├── steering/        # Project steering files
│   └── specs/           # Feature specifications
│
├── .vscode/             # VS Code / Kiro IDE configuration
│   ├── launch.json      # Debug configurations
│   ├── settings.json    # Editor settings
│   └── tasks.json       # Build tasks
│
├── .github/             # GitHub configuration
│   └── workflows/       # GitHub Actions CI/CD workflows
│
└── bin/                 # Compiled binaries (gitignored)
```

## Key Components

### cmd/vox
Main application entry point. Handles command-line arguments and initializes the application.

### internal/state
State machine managing application states (Idle, Recording, Processing) and transitions.

### internal/tray
System tray integration using getlantern/systray library. Manages tray icon and context menu.

### internal/hotkey
Global hotkey registration and handling using golang-design/hotkey library.

### internal/indicator
Coordinates visual (icon changes) and audio (sound playback) feedback for state transitions.

### internal/audio
Audio playback functionality using faiface/beep library for playing feedback sounds.

### internal/platform
Platform-specific abstractions and utilities, including logging infrastructure.

### pkg/config
Configuration structures and default values for the application.

### assets/
Static assets including icons and sound files. See README files in subdirectories for specifications.

## Dependencies

- **github.com/getlantern/systray** - Cross-platform system tray support
- **golang.design/x/hotkey** - Global hotkey registration
- **github.com/faiface/beep** - Audio playback

## Build

```bash
# Build for current platform
go build -o bin/vox ./cmd/vox

# Run tests
go test ./...

# Cross-compile (examples)
GOOS=linux GOARCH=amd64 go build -o bin/vox-linux-amd64 ./cmd/vox
GOOS=windows GOARCH=amd64 go build -o bin/vox-windows-amd64.exe ./cmd/vox
GOOS=darwin GOARCH=arm64 go build -o bin/vox-darwin-arm64 ./cmd/vox
```

## Next Steps

This structure provides the foundation for implementing:
1. State machine (Task 2)
2. System tray integration (Task 3)
3. Hotkey management (Task 4)
4. Visual indicators (Task 5)
5. Audio feedback (Task 6)
