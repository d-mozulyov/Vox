# Vox Builder Docker Image

This Docker image contains all dependencies needed to build Vox for all supported platforms.

**IMPORTANT**: This image is built and published MANUALLY. There is no automated workflow for building this image. When the Dockerfile is updated, follow the instructions below to rebuild and publish.

## Contents

- Alpine Linux 3.19
- Go 1.23
- Make
- Git
- Xvfb (for headless testing)

## Build Strategy

All binaries are built as **static binaries** with `CGO_ENABLED=0`. This means:
- No C dependencies required
- No musl or glibc needed
- Binaries work on any Linux distribution
- Simple cross-compilation for all platforms

The project uses pure Go libraries without CGO:
- `github.com/ebitengine/oto/v3` for audio (pure Go)
- `golang.design/x/hotkey` for hotkeys
- `github.com/getlantern/systray` for system tray

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
make --version
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
  make all

# Run tests in container
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest \
  bash -c "Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 & export DISPLAY=:99.0 && sleep 1 && go test ./..."

# Interactive shell
docker run --rm -it -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest bash
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

When dependencies change (rare - Go version updates):

1. Update the Dockerfile
2. Rebuild the image
3. Push to GitHub Container Registry
4. Commit Dockerfile changes

```bash
docker build -t ghcr.io/d-mozulyov/vox-builder:latest -f docker/builder/Dockerfile .
docker push ghcr.io/d-mozulyov/vox-builder:latest
git add docker/builder/Dockerfile
git commit -m "chore: update Docker builder image"
git push origin main
```

## Image Size

Expected size: ~400-500MB (compressed: ~150-200MB)

Much smaller than the previous musl-based image because we don't need C compilers and development libraries.

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

## Benefits of Pure Go Builds

✓ **Simplicity**: No C dependencies, no cross-compilers needed  
✓ **Speed**: Faster builds without CGO overhead  
✓ **Portability**: Static binaries work everywhere  
✓ **Reliability**: No runtime library dependencies  
✓ **Small image**: Minimal Docker image size  
