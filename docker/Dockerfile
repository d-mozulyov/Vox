# Vox Builder Image
# Cross-compiles Vox for all 6 target platforms
# Inherits clang cross-compilers, lld, xwin and Go from alpine-cross-go
# Smart compiler wrappers auto-inject -fuse-ld=lld when linking is detected

FROM ghcr.io/powertech-center/alpine/cross-go:latest

# Project-specific: Linux x64 dev libraries for CGo (ALSA for audio, X11 for hotkeys)
# xvfb is needed for running tests in headless environment (hotkey library initializes X11 in init())
RUN apk add --no-cache \
    alsa-lib-dev \
    libx11-dev \
    musl-fts-dev \
    xvfb

# Project-specific: aarch64 headers for Linux ARM64 cross-compilation (ALSA, X11)
RUN mkdir -p /usr/aarch64-alpine-linux-musl/etc/apk && \
    echo 'https://dl-cdn.alpinelinux.org/alpine/latest-stable/main' \
        > /usr/aarch64-alpine-linux-musl/etc/apk/repositories && \
    apk add --no-cache --arch aarch64 \
        --root /usr/aarch64-alpine-linux-musl \
        --allow-untrusted --initdb \
        alsa-lib-dev libx11-dev && \
    rm -rf /usr/aarch64-alpine-linux-musl/lib/apk \
           /usr/aarch64-alpine-linux-musl/etc/apk

WORKDIR /workspace
