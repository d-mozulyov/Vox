# Vox Builder Docker Image

This Docker image contains all dependencies needed to build Vox for all supported platforms.

## Contents

- Ubuntu 22.04 base
- Go 1.21.6
- musl-tools and musl-dev for static Linux builds
- aarch64-linux-musl-cross for ARM64 Linux cross-compilation
- All required development libraries (ALSA, X11, AppIndicator, etc.)
- Xvfb for headless testing

## Building the Image

### Prerequisites

- Docker installed (on Windows, use WSL2 with Docker)
- Access to GitHub Container Registry

### Build Command

```bash
# From the repository root
docker build -t ghcr.io/YOUR_USERNAME/vox-builder:latest -f docker/builder/Dockerfile .
```

Replace `YOUR_USERNAME` with your GitHub username.

### Testing the Image

```bash
# Run a container to test
docker run --rm -it ghcr.io/YOUR_USERNAME/vox-builder:latest bash

# Inside the container, verify tools
go version
musl-gcc --version
aarch64-linux-musl-gcc --version
```

## Publishing the Image

### 1. Login to GitHub Container Registry

```bash
# Create a Personal Access Token (PAT) with 'write:packages' scope
# Then login:
echo YOUR_PAT | docker login ghcr.io -u YOUR_USERNAME --password-stdin
```

### 2. Push the Image

```bash
docker push ghcr.io/YOUR_USERNAME/vox-builder:latest
```

### 3. Make the Image Public

1. Go to https://github.com/YOUR_USERNAME?tab=packages
2. Find the `vox-builder` package
3. Click on it
4. Go to "Package settings"
5. Scroll down to "Danger Zone"
6. Click "Change visibility" and select "Public"

## Using the Image Locally

You can use this image for local development to ensure consistency with CI:

```bash
# Run build in container
docker run --rm -v $(pwd):/workspace ghcr.io/YOUR_USERNAME/vox-builder:latest \
  go build -o vox ./cmd/vox

# Run tests in container
docker run --rm -v $(pwd):/workspace ghcr.io/YOUR_USERNAME/vox-builder:latest \
  bash -c "Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 & export DISPLAY=:99.0 && sleep 1 && go test ./..."
```

## Updating the Image

When dependencies change:

1. Update the Dockerfile
2. Rebuild the image with a new tag (e.g., `v1.1`)
3. Push both the versioned tag and `latest`
4. Update GitHub Actions workflows if needed

```bash
docker build -t ghcr.io/YOUR_USERNAME/vox-builder:v1.1 -f docker/builder/Dockerfile .
docker tag ghcr.io/YOUR_USERNAME/vox-builder:v1.1 ghcr.io/YOUR_USERNAME/vox-builder:latest
docker push ghcr.io/YOUR_USERNAME/vox-builder:v1.1
docker push ghcr.io/YOUR_USERNAME/vox-builder:latest
```

## Image Size

Expected size: ~1.5-2GB (compressed: ~600-800MB)

## Troubleshooting

### Docker not found in WSL

```bash
# Install Docker in WSL2
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

### Permission denied when pushing

- Ensure your PAT has `write:packages` scope
- Verify you're logged in: `docker login ghcr.io`

### Image too large

- Consider using multi-stage builds
- Remove unnecessary files in the same RUN command
- Use `--squash` flag when building (experimental)
