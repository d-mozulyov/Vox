# Vox Builder Image
# Minimal Docker image for building pure Go binaries (CGO_ENABLED=0)
# Based on Alpine Linux for small size and fast builds

FROM alpine:3.19

# Install minimal dependencies
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
    alsa-lib-dev

# Install ARM64 cross-compiler
RUN wget https://musl.cc/aarch64-linux-musl-cross.tgz && \
    tar -xzf aarch64-linux-musl-cross.tgz -C /opt && \
    rm aarch64-linux-musl-cross.tgz

ENV PATH="/opt/aarch64-linux-musl-cross/bin:${PATH}"

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
