# Build Instructions

This document describes how to build Vox for all supported platforms.

## Supported Platforms

Vox supports 6 target platforms:

- **Windows**: x64, arm64
- **Linux**: x64, arm64 (with musl for maximum compatibility)
- **macOS**: x64, arm64

## Prerequisites

### All Platforms

- Go 1.21 or later
- Git (for version tagging)

**Note:** You don't need musl tools or cross-compilers for local development. These are only used by GitHub Actions for release builds.

## Build Methods

### Method 1: Using Makefile (Linux/macOS/WSL)

The Makefile provides a simple interface for local development:

```bash
# Show help
make help

# Run tests and build for current platform
make

# Build for current platform only
make build

# Build and run
make run

# Run tests only
make test

# Clean build artifacts
make clean
```

### Method 2: Using PowerShell Script (Windows/Cross-platform)

The PowerShell script works on Windows, Linux, and macOS:

```powershell
# Build for current platform
.\build.ps1

# Run tests before building
.\build.ps1 -Test

# Build and run
.\build.ps1 -Run

# Clean and build
.\build.ps1 -Clean

# Specify custom version
.\build.ps1 -Version "1.0.0"
```

**Note:** Local build scripts only build for your current platform. Cross-platform compilation is handled automatically by GitHub Actions when you push a version tag.

### Method 3: Manual Build

You can also build manually using Go commands:

```bash
# Build for current platform
go build -o vox ./cmd/vox

# Build with version information
VERSION=$(git describe --tags --always --dirty)
go build -ldflags "-s -w -X main.Version=$VERSION" -o vox ./cmd/vox
```

**Note:** For cross-platform builds, use GitHub Actions by pushing a version tag. Manual cross-compilation is not needed for local development.

## Build Output

### Local Development

Local builds create a single binary in the `dist/` directory for your current platform:

```
dist/
└── vox          # or vox.exe on Windows
```

### GitHub Actions (Release Builds)

When you push a version tag, GitHub Actions automatically builds for all 6 platforms:

```
dist/
├── vox-windows-amd64.exe
├── vox-windows-arm64.exe
├── vox-linux-amd64
├── vox-linux-arm64
├── vox-darwin-amd64
└── vox-darwin-arm64
```

## Version Information

The build scripts automatically embed version information:

- **Version**: Extracted from git tags (e.g., `v1.0.0`) or `dev` if no tag exists
- **Build Time**: UTC timestamp of the build

You can check the version of a built binary:

```bash
./vox --version
```

## CI/CD with GitHub Actions

The project uses GitHub Actions for automated builds and releases.

### Workflow Triggers

- **Push to main/develop**: Runs tests only
- **Pull requests**: Runs tests only
- **Tag push (v*)**: Runs tests, builds all platforms, creates GitHub Release

### Creating a Release

1. Ensure all changes are committed
2. Run the release script:
   ```powershell
   .\release.ps1
   ```
3. The script will:
   - Generate a version tag based on current date (YY.M.D)
   - Create and push the tag to GitHub
   - Trigger GitHub Actions to build and publish the release

### Manual Release

You can also create a release manually:

```bash
# Create and push a tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

GitHub Actions will automatically:
1. Run all tests
2. Build binaries for all 6 platforms
3. Create a GitHub Release with all binaries attached

## Linux musl Builds

Linux binaries in GitHub releases are built with musl libc for maximum compatibility across different Linux distributions. This creates statically-linked binaries that don't depend on specific glibc versions.

**Benefits:**
- Works on any Linux distribution (Ubuntu, Debian, Alpine, CentOS, etc.)
- No dependency on system libraries
- Smaller binary size with `-s -w` linker flags

**Note:** This is handled automatically by GitHub Actions. Local Linux builds use standard dynamic linking for simplicity.

## Troubleshooting

### Build errors

If you encounter build errors, ensure you have:
- Go 1.21 or later installed
- All dependencies downloaded (`go mod download`)
- Proper permissions to write to the `dist/` directory

### Permission denied on Linux/macOS

Make sure the build script is executable:

```bash
chmod +x build.ps1
```

## Development Builds

For development, you can build and run directly:

```bash
# Build and run
go run ./cmd/vox

# Build with race detector (for testing)
go build -race -o vox ./cmd/vox

# Build with debug symbols
go build -gcflags="all=-N -l" -o vox ./cmd/vox
```

## Testing

Run tests before building:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests verbosely
go test -v ./...
```

## Additional Resources

- [Go Cross Compilation](https://golang.org/doc/install/source#environment)
- [musl libc](https://musl.libc.org/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
