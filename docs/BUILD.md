# Building Vox

This document describes how to build Vox from source.

## Prerequisites

- Go 1.23 or later
- Make (optional, but recommended)
- Git

## Quick Start

### Using Make (Recommended)

```bash
# Build for all platforms
make all

# Build for specific platform
make windows-amd64
make linux-amd64
make darwin-arm64

# Run tests
make test

# Clean build artifacts
make clean

# Show help
make help
```

### Using Go directly

```bash
# Build for current platform
go build -o vox ./cmd/vox

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o vox-linux-amd64 ./cmd/vox
GOOS=windows GOARCH=amd64 go build -o vox-windows-amd64.exe ./cmd/vox
GOOS=darwin GOARCH=arm64 go build -o vox-darwin-arm64 ./cmd/vox
```

## Build Configuration

### CGO

The project is built with `CGO_ENABLED=0` (pure Go, no C dependencies). This means:
- Static binaries that work on any system
- No runtime dependencies
- Simple cross-compilation
- Smaller Docker images

### Version Information

Version and build time are embedded during build:

```bash
# Using Make
make VERSION=1.0.0 all

# Using Go directly
go build -ldflags "-X main.Version=1.0.0 -X main.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" ./cmd/vox
```

## Supported Platforms

The project supports 6 target platforms:

| OS      | Architecture | Binary Name              |
|---------|--------------|--------------------------|
| Windows | amd64        | vox-windows-amd64.exe    |
| Windows | arm64        | vox-windows-arm64.exe    |
| Linux   | amd64        | vox-linux-amd64          |
| Linux   | arm64        | vox-linux-arm64          |
| macOS   | amd64        | vox-darwin-amd64         |
| macOS   | arm64        | vox-darwin-arm64         |

## Build Output

All binaries are placed in the `dist/` directory:

```
dist/
├── vox-windows-amd64.exe
├── vox-windows-arm64.exe
├── vox-linux-amd64
├── vox-linux-arm64
├── vox-darwin-amd64
└── vox-darwin-arm64
```

## Dependencies

The project uses pure Go libraries without CGO:

- `github.com/ebitengine/oto/v3` - Audio playback (pure Go)
- `golang.design/x/hotkey` - Global hotkeys
- `github.com/getlantern/systray` - System tray integration

To update dependencies:

```bash
go get -u ./...
go mod tidy
```

## Testing

```bash
# Run all tests
go test ./...

# Or using Make
make test

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/state
```

## Docker Build

For consistent builds across environments, use Docker:

```bash
# Build using Docker
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest make all

# Run tests using Docker
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest make test
```

See [DOCKER_BUILD.md](DOCKER_BUILD.md) for more details.

## Troubleshooting

### Build fails with "package not found"

```bash
# Download dependencies
go mod download
go mod tidy
```

### Cross-compilation fails

Ensure `CGO_ENABLED=0` is set:

```bash
export CGO_ENABLED=0
make all
```

### Tests fail on headless system

Some tests require a display. Use Xvfb:

```bash
Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 &
export DISPLAY=:99.0
go test ./...
```

## CI/CD

The project uses GitHub Actions for automated builds. See `.github/workflows/build.yml` for the CI/CD configuration.

Builds are triggered on:
- Push to `main` or `develop` branches
- Pull requests to `main`
- Version tags (`v*`)

## Release Process

1. Tag a new version:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Actions will automatically:
   - Run tests
   - Build for all platforms
   - Create a GitHub Release
   - Upload binaries as release assets

## Additional Resources

- [Makefile](../Makefile) - Build configuration
- [DOCKER_BUILD.md](DOCKER_BUILD.md) - Docker build environment
- [LOCAL_DOCKER_TESTING.md](LOCAL_DOCKER_TESTING.md) - Local Docker testing guide
