# Vox Builder Image
# Docker image for cross-compiling Vox for all 6 target platforms
# Based on Alpine Linux with musl for portable Linux binaries

FROM alpine:3.19

# Install x86_64 build dependencies
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
    clang \
    llvm-dev \
    cmake \
    python3 \
    libxml2-dev \
    openssl-dev \
    patch \
    bsd-compat-headers \
    linux-headers

# Install ARM64 cross-compiler (musl)
RUN wget https://musl.cc/aarch64-linux-musl-cross.tgz && \
    tar -xzf aarch64-linux-musl-cross.tgz -C /opt && \
    rm aarch64-linux-musl-cross.tgz

ENV PATH="/opt/aarch64-linux-musl-cross/bin:${PATH}"

# Create ARM64 sysroot with dev libraries for cross-compilation
RUN mkdir -p /opt/aarch64-sysroot/etc/apk && \
    cp /etc/apk/repositories /opt/aarch64-sysroot/etc/apk/ && \
    cp -r /etc/apk/keys /opt/aarch64-sysroot/etc/apk/ && \
    apk --arch aarch64 --root /opt/aarch64-sysroot \
    --allow-untrusted --initdb --no-cache \
    add alsa-lib-dev libx11-dev gtk+3.0-dev libayatana-appindicator-dev zlib-dev && \
    rm -f /opt/aarch64-sysroot/usr/include/stdc-predef.h

# Build osxcross for macOS cross-compilation
ARG MACOSX_SDK_VERSION=14.5
RUN git clone --depth 1 https://github.com/tpoechtrager/osxcross.git /tmp/osxcross && \
    wget -q -O /tmp/osxcross/tarballs/MacOSX${MACOSX_SDK_VERSION}.sdk.tar.xz \
    https://github.com/joseluisq/macosx-sdks/releases/download/${MACOSX_SDK_VERSION}/MacOSX${MACOSX_SDK_VERSION}.sdk.tar.xz && \
    cd /tmp/osxcross && \
    UNATTENDED=1 TARGET_DIR=/opt/osxcross ./build.sh && \
    rm -rf /tmp/osxcross

ENV PATH="/opt/osxcross/bin:${PATH}"

# Install Go 1.25.5 manually
RUN wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz && \
    rm go1.25.5.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Set working directory
WORKDIR /workspace

# Verify installations
RUN go version && \
    make --version && \
    gcc --version && \
    aarch64-linux-musl-gcc --version && \
    ls /opt/osxcross/bin/*-clang* || echo "osxcross clang not found"

# Default command
CMD ["/bin/bash"]
