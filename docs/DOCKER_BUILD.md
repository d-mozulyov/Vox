# Docker Build Environment

This document explains how the project uses Docker for CI/CD optimization and how to customize it for your needs.

## Overview

The project uses a pre-built Docker image (`ghcr.io/d-mozulyov/vox-builder:latest`) that contains all build dependencies:
- Go 1.23
- Make
- Git
- Xvfb (for headless testing)

This approach speeds up CI/CD by avoiding dependency installation on each run.

## Build Strategy

All binaries are built as **static binaries** with `CGO_ENABLED=0`:
- No C dependencies
- No musl or glibc needed
- Works on any Linux distribution
- Simple cross-compilation

The project uses pure Go libraries:
- `github.com/ebitengine/oto/v3` - audio (pure Go, no ALSA/CGO)
- `golang.design/x/hotkey` - hotkeys
- `github.com/getlantern/systray` - system tray

## Using the Official Image

By default, GitHub Actions workflows use the official image: `ghcr.io/d-mozulyov/vox-builder:latest`

No configuration needed - just push your code and CI/CD will work.

## For Forks: Using Your Own Docker Image

If you fork this project and want to use your own Docker image:

### Option 1: Use the Official Image (Recommended)

The easiest option - just use the official image. It's public and works out of the box.

### Option 2: Build and Use Your Own Image

If you need custom dependencies or want full control:

#### Step 1: Customize the Dockerfile

Edit `docker/builder/Dockerfile` to add your dependencies:

```dockerfile
# Example: Add additional tools
RUN apk add --no-cache \
    your-custom-package
```

#### Step 2: Build and Publish Your Image

**Build Locally:**

```bash
# Build the image
docker build -t ghcr.io/YOUR_USERNAME/vox-builder:latest -f docker/builder/Dockerfile .

# Create a Personal Access Token (PAT) on GitHub:
# https://github.com/settings/tokens
# Scopes needed: write:packages

# Login to GitHub Container Registry
echo YOUR_PAT | docker login ghcr.io -u YOUR_USERNAME --password-stdin

# Push the image
docker push ghcr.io/YOUR_USERNAME/vox-builder:latest
```

#### Step 3: Make the Image Public

1. Go to: `https://github.com/YOUR_USERNAME?tab=packages`
2. Find the `vox-builder` package
3. Click on it
4. Go to "Package settings"
5. Scroll to "Danger Zone"
6. Click "Change visibility" → "Public"

#### Step 4: Configure Your Fork to Use Your Image

**Option A: Set Repository Variable (Recommended)**

1. Go to your fork's settings: `https://github.com/YOUR_USERNAME/Vox/settings/variables/actions`
2. Click "New repository variable"
3. Name: `BUILDER_IMAGE`
4. Value: `ghcr.io/YOUR_USERNAME/vox-builder:latest`
5. Click "Add variable"

The workflows will automatically use your image.

**Option B: Edit Workflow Files**

Edit `.github/workflows/build.yml` and change the default image:

```yaml
# Find this line:
image: ${{ vars.BUILDER_IMAGE || 'ghcr.io/d-mozulyov/vox-builder:latest' }}

# Change to:
image: ${{ vars.BUILDER_IMAGE || 'ghcr.io/YOUR_USERNAME/vox-builder:latest' }}
```

## Updating the Docker Image

When you need to update dependencies (rare - Go version updates):

### Manual Build and Publish

1. Edit `docker/builder/Dockerfile` to update dependencies
2. Build the image locally:
   ```bash
   docker build -t ghcr.io/d-mozulyov/vox-builder:latest -f docker/builder/Dockerfile .
   ```
3. Test the image (see next section)
4. Login to GitHub Container Registry:
   ```bash
   echo YOUR_PAT | docker login ghcr.io -u d-mozulyov --password-stdin
   ```
5. Push the image:
   ```bash
   docker push ghcr.io/d-mozulyov/vox-builder:latest
   ```
6. Commit and push the Dockerfile changes:
   ```bash
   git add docker/builder/Dockerfile
   git commit -m "chore: update Docker builder image"
   git push origin main
   ```

**Why manual?** 
- Simple and explicit control
- No complex automation for rare operations
- Dockerfile changes happen rarely (Go version updates)
- Follows KISS principle

### For Forks

Same process, but replace `d-mozulyov` with your GitHub username.

## Testing the Docker Image Locally

Before publishing the image, test it locally to ensure the build works:

### Step 1: Build the Docker Image

```bash
# From the repository root
docker build -t ghcr.io/d-mozulyov/vox-builder:latest -f docker/builder/Dockerfile .
```

### Step 2: Test Interactive Shell

```bash
# Start an interactive shell in the container
docker run --rm -it -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest bash

# Inside the container, verify tools
go version
make --version
git --version

# Exit the container
exit
```

### Step 3: Test Building the Project

```bash
# Run the full build
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest make all

# Check the output
ls -lh dist/
```

### Step 4: Test Running Tests

```bash
# Run tests with virtual display
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest \
  bash -c "Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 & export DISPLAY=:99.0 && sleep 1 && go test ./..."
```

### Step 5: Test Specific Platform Build

```bash
# Test building for a specific platform
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest make linux-amd64

# Verify the binary
file dist/vox-linux-amd64
```

### On Windows with WSL

If you're on Windows and using WSL:

```bash
# Navigate to your project directory in WSL
cd /mnt/c/path/to/your/project

# Or if your project is in WSL filesystem
cd ~/vox

# Run the same Docker commands as above
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest make all
```

**Note**: Docker in WSL2 can access both Windows filesystem (`/mnt/c/...`) and WSL filesystem (`~/...`).

## Local Development with Docker

You can use the Docker image locally to ensure consistency with CI:

```bash
# Build the project
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest make all

# Build specific platform
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest make windows-amd64

# Run tests
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest make test

# Clean build artifacts
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest make clean

# Interactive shell for debugging
docker run --rm -it -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest bash
```

## Files Related to Docker Build

- `docker/builder/Dockerfile` - Docker image definition
- `docker/builder/README.md` - Detailed Docker image documentation
- `.github/workflows/build.yml` - Main CI/CD workflow (uses the image)
- `Makefile` - Build targets for all platforms

## Troubleshooting

### GitHub Actions can't pull the image

- Ensure the image is public (see Step 3 above)
- Check the image name in workflow files
- Verify the image exists: `https://github.com/USERNAME?tab=packages`

### Local docker push fails with "denied"

- Check your PAT has `write:packages` scope
- Verify you're logged in: `docker login ghcr.io`
- Ensure the image name matches: `ghcr.io/USERNAME/vox-builder`

### Docker not found in WSL

```bash
# Install Docker in WSL2
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Restart WSL or logout/login
```

### Build fails in Docker but works locally

- Check Go version matches between local and Docker
- Verify all dependencies are in go.mod
- Run `go mod tidy` before building

## Benefits of This Approach

✓ **Speed**: Fast CI/CD (no dependency installation on each run)  
✓ **Consistency**: Same environment everywhere (CI, local, contributors)  
✓ **Reliability**: No risk of package manager failures  
✓ **Simplicity**: Pure Go builds, no C dependencies  
✓ **Portability**: Static binaries work everywhere  
✓ **Small image**: ~400-500MB (compressed: ~150-200MB)  

## Additional Resources

- [GitHub Container Registry Documentation](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Docker Multi-platform Builds](https://docs.docker.com/build/building/multi-platform/)
- [GitHub Actions with Containers](https://docs.github.com/en/actions/using-jobs/running-jobs-in-a-container)

---

**Note:** The Docker image is updated manually when dependencies change (rare). This keeps the process simple and explicit, following the KISS principle.
