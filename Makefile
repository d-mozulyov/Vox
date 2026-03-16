# Vox Build Configuration
# Cross-compilation for all 6 target platforms
# CC is set via CGO_CC_<goos>_<goarch> env vars from alpine-cross-go.
# Smart compiler wrappers auto-inject -fuse-ld=lld when linking is needed.
# No CGO_LDFLAGS required — wrappers handle everything.

# Global settings
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)

# Output directory
DIST_DIR := dist

# Go build command
GO_BUILD := go build -buildvcs=false -ldflags "$(LDFLAGS)"

# Targets
.PHONY: all clean test \
	windows-x64 windows-arm64 \
	linux-x64 linux-arm64 \
	macos-x64 macos-arm64

# Build all platforms
all: linux-x64 linux-arm64 macos-x64 macos-arm64 windows-x64 windows-arm64
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
	@Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 & export DISPLAY=:99.0 && sleep 1 && go test ./internal/... ./pkg/...
	@echo "Tests complete!"

# Create dist directory
$(DIST_DIR):
	@mkdir -p $(DIST_DIR)

# === Linux ===

linux-x64: $(DIST_DIR)
	@echo "Building for Linux x64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=1 \
		CC=$$CGO_CC_linux_amd64 \
		$(GO_BUILD) -o $(DIST_DIR)/vox-linux-x64 ./cmd/vox

linux-arm64: $(DIST_DIR)
	@echo "Building for Linux arm64..."
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=1 \
		CC=$$CGO_CC_linux_arm64 \
		$(GO_BUILD) -o $(DIST_DIR)/vox-linux-arm64 ./cmd/vox

# === macOS ===

macos-x64: $(DIST_DIR)
	@echo "Building for macOS x64..."
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 \
		CC=$$CGO_CC_darwin_amd64 \
		$(GO_BUILD) -o $(DIST_DIR)/vox-macos-x64 ./cmd/vox

macos-arm64: $(DIST_DIR)
	@echo "Building for macOS arm64..."
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 \
		CC=$$CGO_CC_darwin_arm64 \
		$(GO_BUILD) -o $(DIST_DIR)/vox-macos-arm64 ./cmd/vox

# === Windows ===

windows-x64: $(DIST_DIR)
	@echo "Building for Windows x64..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
		CC=$$CGO_CC_windows_amd64 \
		$(GO_BUILD) -o $(DIST_DIR)/vox-windows-x64.exe ./cmd/vox

windows-arm64: $(DIST_DIR)
	@echo "Building for Windows arm64..."
	@GOOS=windows GOARCH=arm64 CGO_ENABLED=1 \
		CC=$$CGO_CC_windows_arm64 \
		$(GO_BUILD) -o $(DIST_DIR)/vox-windows-arm64.exe ./cmd/vox

# Help
help:
	@echo "Vox Build System"
	@echo "================"
	@echo ""
	@echo "Available targets:"
	@echo "  all            - Build for all platforms (default)"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  linux-x64      - Build for Linux x64"
	@echo "  linux-arm64    - Build for Linux arm64"
	@echo "  macos-x64      - Build for macOS x64"
	@echo "  macos-arm64    - Build for macOS arm64"
	@echo "  windows-x64    - Build for Windows x64"
	@echo "  windows-arm64  - Build for Windows arm64"
