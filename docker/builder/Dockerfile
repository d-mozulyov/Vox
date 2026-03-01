# Vox Builder Image
# Docker image for cross-compiling Vox for all 6 target platforms
# Based on Alpine Linux with musl for portable Linux binaries
# Uses Zig as a universal C cross-compiler (replaces osxcross + musl-cross)

FROM alpine:3.19

# Install build dependencies (host platform)
RUN apk add --no-cache \
    bash \
    git \
    make \
    ca-certificates \
    xvfb \
    wget \
    tar \
    xz \
    gcc \
    musl-dev \
    pkgconf \
    alsa-lib-dev \
    libx11-dev \
    gtk+3.0-dev \
    libayatana-appindicator-dev \
    linux-headers \
    musl-fts-dev

# Install Zig — universal C/C++ cross-compiler for all 6 targets
ARG ZIG_VERSION=0.15.2
RUN wget https://ziglang.org/download/${ZIG_VERSION}/zig-x86_64-linux-${ZIG_VERSION}.tar.xz && \
    tar -xJf zig-x86_64-linux-${ZIG_VERSION}.tar.xz -C /opt && \
    ln -s /opt/zig-x86_64-linux-${ZIG_VERSION}/zig /usr/local/bin/zig && \
    rm zig-x86_64-linux-${ZIG_VERSION}.tar.xz

# Create zig-cc wrapper scripts for non-macOS targets
# Go's CC env var works best with a single executable path
RUN for target in \
    x86_64-linux-musl \
    aarch64-linux-musl \
    x86_64-windows-gnu \
    aarch64-windows-gnu; \
    do \
    printf '#!/bin/sh\nexec zig cc -target %s "$@"\n' "$target" > /usr/local/bin/zig-cc-${target} && \
    chmod +x /usr/local/bin/zig-cc-${target}; \
    done

# Download macOS SDK for framework headers (Cocoa, AppKit, etc.)
# Only extracting the SDK — no osxcross build needed
ARG MACOSX_SDK_VERSION=14.5
RUN wget -q https://github.com/joseluisq/macosx-sdks/releases/download/${MACOSX_SDK_VERSION}/MacOSX${MACOSX_SDK_VERSION}.sdk.tar.xz && \
    tar -xJf MacOSX${MACOSX_SDK_VERSION}.sdk.tar.xz -C /opt && \
    rm MacOSX${MACOSX_SDK_VERSION}.sdk.tar.xz

# Create zig-cc wrapper scripts for macOS targets (with SDK sysroot + framework path)
# -Wno-nullability-completeness: macOS SDK headers use Apple nullability annotations
#   that trigger -Werror in CGo builds
RUN SDK=/opt/MacOSX${MACOSX_SDK_VERSION}.sdk && \
    for target in x86_64-macos aarch64-macos; do \
    printf '#!/bin/sh\nexec zig cc -target %s --sysroot=%s -F%s/System/Library/Frameworks -I%s/usr/include -Wno-nullability-completeness "$@"\n' \
    "$target" "$SDK" "$SDK" "$SDK" \
    > /usr/local/bin/zig-cc-${target} && \
    chmod +x /usr/local/bin/zig-cc-${target}; \
    done

# Create ARM64 sysroot with dev libraries for Linux ARM64 cross-compilation
# Zig provides musl libc, but platform-specific libs (GTK, ALSA) need a sysroot
RUN mkdir -p /opt/aarch64-sysroot/etc/apk && \
    cp /etc/apk/repositories /opt/aarch64-sysroot/etc/apk/ && \
    cp -r /etc/apk/keys /opt/aarch64-sysroot/etc/apk/ && \
    apk --arch aarch64 --root /opt/aarch64-sysroot \
    --allow-untrusted --initdb --no-cache \
    add alsa-lib-dev libx11-dev gtk+3.0-dev libayatana-appindicator-dev zlib-dev && \
    rm -f /opt/aarch64-sysroot/usr/include/stdc-predef.h

# Install Go
ARG GO_VERSION=1.25.5
RUN wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

WORKDIR /workspace

# Verify installations
RUN go version && \
    zig version && \
    make --version && \
    gcc --version && \
    echo "=== Zig CC wrappers ===" && \
    for f in /usr/local/bin/zig-cc-*; do echo "--- $f ---"; cat "$f"; done && \
    echo "=== macOS SDK check ===" && \
    ls /opt/MacOSX*.sdk/System/Library/Frameworks/Cocoa.framework/Headers/Cocoa.h

CMD ["/bin/bash"]
