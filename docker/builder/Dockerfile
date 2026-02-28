FROM alpine:3.19

# Install base tools and dependencies
RUN apk add --no-cache \
    bash \
    wget \
    curl \
    git \
    build-base \
    pkgconfig \
    ca-certificates \
    xvfb \
    # Development libraries for musl-based builds
    alsa-lib-dev \
    libx11-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    mesa-dev \
    gtk+3.0-dev \
    libayatana-appindicator-dev

# Install Go 1.21.6
RUN wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz && \
    rm go1.21.6.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Install cross-compiler for arm64 (musl-based)
RUN wget https://musl.cc/aarch64-linux-musl-cross.tgz && \
    tar -xzf aarch64-linux-musl-cross.tgz && \
    mv aarch64-linux-musl-cross /opt/ && \
    rm aarch64-linux-musl-cross.tgz

ENV PATH="/opt/aarch64-linux-musl-cross/bin:${PATH}"

# Set working directory
WORKDIR /workspace

# Verify installations
RUN go version && \
    gcc --version && \
    aarch64-linux-musl-gcc --version
