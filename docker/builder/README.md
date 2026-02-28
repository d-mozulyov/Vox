# Vox Builder Docker Image

This Docker image contains all dependencies needed to build Vox for all supported platforms.

**IMPORTANT**: This image is built and published MANUALLY. There is no automated workflow for building this image. When the Dockerfile is updated, follow the instructions below to rebuild and publish.

## Contents

- Alpine Linux 3.19 (musl-based)
- Go 1.21.6
- musl gcc for native Linux builds
- aarch64-linux-musl-cross for ARM64 Linux cross-compilation
- All required development libraries (ALSA, X11, AppIndicator, GTK, etc.)
- Xvfb for headless testing

## Build Strategy

Linux binaries are built with musl using dynamic linking. Alpine Linux is used as the base image because it's natively built on musl.

The project depends on CGO libraries:
- `github.com/hajimehoshi/oto` requires ALSA
- `golang.design/x/hotkey` requires X11
- `github.com/getlantern/systray` requires GTK/AppIndicator

Static linking with all these libraries is complex in Alpine due to missing static packages. Dynamic linking provides working binaries that depend on system libraries (musl-based).

## Building the Image

### Prerequisites

- Docker installed (on Windows, use WSL2 with Docker)
- Access to GitHub Container Registry

### Build Command

```bash
# From the repository root
docker build -t ghcr.io/d-mozulyov/vox-builder:latest -f docker/builder/Dockerfile .
```

For forks, replace `d-mozulyov` with your GitHub username.

### Testing the Image

```bash
# Run a container to test
docker run --rm -it ghcr.io/d-mozulyov/vox-builder:latest bash

# Inside the container, verify tools
go version
gcc --version
aarch64-linux-musl-gcc --version
```

## Publishing the Image

### 1. Login to GitHub Container Registry

```bash
# Create a Personal Access Token (PAT) with 'write:packages' scope
# Then login:
echo YOUR_PAT | docker login ghcr.io -u d-mozulyov --password-stdin
```

### 2. Push the Image

```bash
docker push ghcr.io/d-mozulyov/vox-builder:latest
```

### 3. Make the Image Public

1. Go to https://github.com/d-mozulyov?tab=packages
2. Find the `vox-builder` package
3. Click on it
4. Go to "Package settings"
5. Scroll down to "Danger Zone"
6. Click "Change visibility" and select "Public"

## Using the Image Locally

You can use this image for local development to ensure consistency with CI:

```bash
# Run build in container
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest \
  go build -o vox ./cmd/vox

# Run tests in container
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest \
  bash -c "Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 & export DISPLAY=:99.0 && sleep 1 && go test ./..."
```

## For Forks and Contributors

If you fork this project and want to use your own Docker image:

1. Build and publish your image: `ghcr.io/yourusername/vox-builder:latest`
2. In your fork's GitHub repository settings:
   - Go to: Settings → Secrets and variables → Actions → Variables
   - Create a variable named `BUILDER_IMAGE`
   - Set value to: `ghcr.io/yourusername/vox-builder:latest`
3. The workflow will automatically use your image

The main repository uses: `ghcr.io/d-mozulyov/vox-builder:latest`

## Updating the Image

When dependencies change:

1. Update the Dockerfile
2. Rebuild the image with a new tag (e.g., `v1.1`)
3. Push both the versioned tag and `latest`
4. Update GitHub Actions workflows if needed

```bash
docker build -t ghcr.io/d-mozulyov/vox-builder:v1.1 -f docker/builder/Dockerfile .
docker tag ghcr.io/d-mozulyov/vox-builder:v1.1 ghcr.io/d-mozulyov/vox-builder:latest
docker push ghcr.io/d-mozulyov/vox-builder:v1.1
docker push ghcr.io/d-mozulyov/vox-builder:latest
```

Or use the automated workflow:
- Go to: https://github.com/d-mozulyov/Vox/actions
- Select "Build Docker Image"
- Click "Run workflow"

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
