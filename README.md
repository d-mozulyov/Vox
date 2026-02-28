# Vox - Voice Input Assistant

Vox is an open-source voice-to-text input assistant that helps you type with your voice across any application.

## Features

- **Open Source**: MIT licensed, free to use and modify
- **Privacy-Focused**: Your data stays on your machine
- **Context-Aware**: Intelligent transcription based on active application context
- **Flexible AI Backend**: Works with Mistral AI (free tier available) or any OpenAI-compatible API
- **Cross-Platform**: Supports Windows, Linux, and macOS (x64 and arm64)
- **System Tray Integration**: Runs in the background, always ready when you need it
- **Hotkey Support**: Quick start/stop recording with customizable keyboard shortcuts

## Why Vox?

Unlike traditional transcription tools, Vox uses the `/v1/chat/completions` endpoint instead of classic transcription APIs. This allows for:

- Rich context injection for better accuracy
- Application-specific vocabulary and terminology
- Project-level glossaries
- Analysis of previously entered text
- Significantly improved transcription quality

## Quick Start

### Installation

Download the latest release for your platform from the [Releases](https://github.com/YOUR_USERNAME/vox/releases) page.

### Configuration

1. Launch Vox - it will appear in your system tray
2. Right-click the tray icon and select "Settings"
3. Configure your AI backend:
   - **Mistral AI** (recommended for free tier): Get your API key from [Mistral AI](https://console.mistral.ai/)
   - **Custom OpenAI-compatible API**: Enter your endpoint URL, model name, and API key

### Usage

1. Press the hotkey (default: `Ctrl+Shift+V`) to start recording
2. Speak your text
3. Press the hotkey again to stop recording
4. Vox will transcribe and insert the text at your cursor position

## Building from Source

### Prerequisites

- Go 1.21 or later
- Git

### Build

```bash
git clone https://github.com/YOUR_USERNAME/vox.git
cd vox
go build -o vox ./cmd/vox
```

### Cross-Compilation

```bash
# Windows x64
GOOS=windows GOARCH=amd64 go build -o vox-windows-amd64.exe ./cmd/vox

# Linux x64 (with musl for compatibility)
GOOS=linux GOARCH=amd64 go build -o vox-linux-amd64 ./cmd/vox

# macOS arm64
GOOS=darwin GOARCH=arm64 go build -o vox-darwin-arm64 ./cmd/vox
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Powered by [Mistral AI](https://mistral.ai/) voice models
- Built with Go and love for the open-source community
