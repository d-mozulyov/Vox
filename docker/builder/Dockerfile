# Vox Builder Image
# Minimal Docker image for building pure Go binaries (CGO_ENABLED=0)
# Based on Alpine Linux for small size and fast builds

FROM golang:1.23-alpine3.19

# Install minimal dependencies
RUN apk add --no-cache \
    bash \
    git \
    make \
    ca-certificates \
    xvfb

# Set working directory
WORKDIR /workspace

# Verify installations
RUN go version && \
    make --version

# Default command
CMD ["/bin/bash"]
