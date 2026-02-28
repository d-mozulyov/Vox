#!/bin/bash

# Cross-platform build script for Vox
# Builds binaries for all supported platforms

set -e

VERSION=${VERSION:-"0.1.0"}
OUTPUT_DIR="dist"

echo "Building Vox v${VERSION} for all platforms..."

# Clean previous builds
rm -rf ${OUTPUT_DIR}
mkdir -p ${OUTPUT_DIR}

# Build for each platform
platforms=(
    "windows/amd64"
    "windows/arm64"
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"

    output_name="vox-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi

    echo "Building for ${GOOS}/${GOARCH}..."

    GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-s -w -X main.Version=${VERSION}" \
        -o "${OUTPUT_DIR}/${output_name}" \
        ./cmd/vox

    echo "  ✓ ${output_name}"
done

echo ""
echo "Build complete! Binaries are in ${OUTPUT_DIR}/"
ls -lh ${OUTPUT_DIR}/
