# Vox - Simple build Makefile for local development
# Cross-platform builds are handled by GitHub Actions

# Project configuration
PROJECT_NAME := vox
MODULE := github.com/d-mozulyov/vox
CMD_PATH := ./cmd/vox
OUTPUT_DIR := dist

# Version from git tag or default
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)

# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_CYAN := \033[36m

.PHONY: all clean test build run help

# Default target
all: test build

help:
	@echo "$(COLOR_BOLD)Vox Build System$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_CYAN)Available targets:$(COLOR_RESET)"
	@echo "  $(COLOR_GREEN)make$(COLOR_RESET)              - Run tests and build for current platform"
	@echo "  $(COLOR_GREEN)make build$(COLOR_RESET)        - Build for current platform"
	@echo "  $(COLOR_GREEN)make test$(COLOR_RESET)         - Run all tests"
	@echo "  $(COLOR_GREEN)make run$(COLOR_RESET)          - Build and run the application"
	@echo "  $(COLOR_GREEN)make clean$(COLOR_RESET)        - Remove build artifacts"
	@echo "  $(COLOR_GREEN)make help$(COLOR_RESET)         - Show this help message"
	@echo ""
	@echo "$(COLOR_CYAN)Note:$(COLOR_RESET) Cross-platform builds are handled by GitHub Actions."
	@echo "For local development, use 'make' or 'make build' to build for your platform."

# Run tests
test:
	@echo "$(COLOR_CYAN)Running tests...$(COLOR_RESET)"
	@go test -v ./...

# Build for current platform
build:
	@echo "$(COLOR_CYAN)Building for current platform...$(COLOR_RESET)"
	@mkdir -p $(OUTPUT_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(OUTPUT_DIR)/$(PROJECT_NAME) $(CMD_PATH)
	@echo "$(COLOR_GREEN)✓ Build complete: $(OUTPUT_DIR)/$(PROJECT_NAME)$(COLOR_RESET)"

# Build and run
run: build
	@echo "$(COLOR_CYAN)Running application...$(COLOR_RESET)"
	@$(OUTPUT_DIR)/$(PROJECT_NAME)

# Clean build artifacts
clean:
	@echo "$(COLOR_CYAN)Cleaning build artifacts...$(COLOR_RESET)"
	@rm -rf $(OUTPUT_DIR)
	@rm -f $(PROJECT_NAME) $(PROJECT_NAME).exe
	@echo "$(COLOR_GREEN)✓ Clean complete$(COLOR_RESET)"
