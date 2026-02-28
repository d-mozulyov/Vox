FROM ubuntu:22.04

# Avoid interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install base tools and dependencies
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    git \
    build-essential \
    musl-tools \
    musl-dev \
    libasound2-dev \
    libx11-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libgl1-mesa-dev \
    libayatana-appindicator3-dev \
    xvfb \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.21.6
RUN wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz && \
    rm go1.21.6.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Install musl cross-compiler for arm64
RUN wget https://musl.cc/aarch64-linux-musl-cross.tgz && \
    tar -xzf aarch64-linux-musl-cross.tgz && \
    mv aarch64-linux-musl-cross /opt/ && \
    rm aarch64-linux-musl-cross.tgz

ENV PATH="/opt/aarch64-linux-musl-cross/bin:${PATH}"

# Set working directory
WORKDIR /workspace

# Verify installations
RUN go version && \
    musl-gcc --version && \
    aarch64-linux-musl-gcc --version

# Reset DEBIAN_FRONTEND
ENV DEBIAN_FRONTEND=
