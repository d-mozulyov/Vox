# Vox Build Configuration
# This Makefile provides targets for building Vox on all supported platforms

# Global settings
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)

# Output directory
DIST_DIR := dist

# Go build command
GO_BUILD := go build -buildvcs=false -ldflags "$(LDFLAGS)"

# Targets
.PHONY: all clean test windows-amd64 windows-arm64 linux-amd64 linux-arm64 darwin-amd64 darwin-arm64

# Build all platforms
all: windows-amd64 windows-arm64 linux-amd64 linux-arm64 darwin-amd64 darwin-arm64
	@echo "================================"
	@echo "Build complete! Version: $(VERSION)"
	@ls -lh $(DIST_DIR)/

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(DIST_DIR)
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...
	@echo "Tests complete!"

# Create dist directory
$(DIST_DIR):
	@mkdir -p $(DIST_DIR)

# Windows amd64
windows-amd64: $(DIST_DIR)
	@echo "Building for Windows amd64..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO_BUILD) -o $(DIST_DIR)/vox-windows-amd64.exe ./cmd/vox

# Windows arm64
windows-arm64: $(DIST_DIR)
	@echo "Building for Windows arm64..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=arm64 $(GO_BUILD) -o $(DIST_DIR)/vox-windows-arm64.exe ./cmd/vox

# Linux amd64
linux-amd64: $(DIST_DIR)
	@echo "Building for Linux amd64..."
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(DIST_DIR)/vox-linux-amd64 ./cmd/vox

# ARM64 cross-compilation sysroot (populated in Docker builder image)
AARCH64_SYSROOT := /opt/aarch64-sysroot

# Linux arm64
linux-arm64: $(DIST_DIR)
	@echo "Building for Linux arm64..."
	@CGO_ENABLED=1 GOOS=linux GOARCH=arm64 \
	CC=aarch64-linux-musl-gcc \
	PKG_CONFIG_LIBDIR=$(AARCH64_SYSROOT)/usr/lib/pkgconfig:$(AARCH64_SYSROOT)/usr/share/pkgconfig \
	PKG_CONFIG_SYSROOT_DIR=$(AARCH64_SYSROOT) \
	CGO_CFLAGS="-I$(AARCH64_SYSROOT)/usr/include" \
	CGO_LDFLAGS="-L$(AARCH64_SYSROOT)/usr/lib" \
	$(GO_BUILD) -o $(DIST_DIR)/vox-linux-arm64 ./cmd/vox

# macOS amd64
darwin-amd64: $(DIST_DIR)
	@echo "Building for macOS amd64..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO_BUILD) -o $(DIST_DIR)/vox-darwin-amd64 ./cmd/vox

# macOS arm64
darwin-arm64: $(DIST_DIR)
	@echo "Building for macOS arm64..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO_BUILD) -o $(DIST_DIR)/vox-darwin-arm64 ./cmd/vox

# Help
help:
	@echo "Vox Build System"
	@echo "================"
	@echo ""
	@echo "Available targets:"
	@echo "  all            - Build for all platforms (default)"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  windows-amd64  - Build for Windows x64"
	@echo "  windows-arm64  - Build for Windows ARM64"
	@echo "  linux-amd64    - Build for Linux x64"
	@echo "  linux-arm64    - Build for Linux ARM64"
	@echo "  darwin-amd64   - Build for macOS x64"
	@echo "  darwin-arm64   - Build for macOS ARM64"
	@echo ""
	@echo "Environment variables:"
	@echo "  VERSION        - Version string (default: git describe)"
	@echo ""
	@echo "Examples:"
	@echo "  make all                    - Build all platforms"
	@echo "  make windows-amd64          - Build only Windows x64"
	@echo "  make VERSION=1.0.0 all      - Build with specific version"
