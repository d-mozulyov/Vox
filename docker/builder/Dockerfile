# Vox Builder Image
# Docker image for cross-compiling Vox for all target platforms
# Based on Alpine Linux with musl for portable Linux binaries

FROM alpine:3.19

# Install build dependencies
RUN apk add --no-cache \
    bash \
    git \
    make \
    ca-certificates \
    xvfb \
    wget \
    tar \
    gcc \
    musl-dev \
    pkgconf \
    alsa-lib-dev \
    libx11-dev \
    gtk+3.0-dev \
    libayatana-appindicator-dev

# Install ARM64 cross-compiler
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
    aarch64-linux-musl-gcc --version

# Default command
CMD ["/bin/bash"]
