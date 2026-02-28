# Local Docker Image Testing

This guide describes how to test the Docker builder image locally before publishing it to GitHub Container Registry.

## Prerequisites

- Docker installed in WSL2 (if working on Windows)
- Project accessible in WSL (e.g., `/mnt/c/Projects/Vox` or `~/vox`)

## Installing Docker in WSL2 (if not already installed)

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker $USER

# Restart WSL or logout/login to apply changes
```

## Step 1: Build the Docker Image

```bash
# Navigate to project directory
cd /mnt/c/Projects/Vox  # or your path

# Build the Docker image
docker build -t vox-builder-test:local -f docker/builder/Dockerfile .
```

This creates a local image tagged as `vox-builder-test:local`.

## Step 2: Verify the Image

```bash
# Run container in interactive mode
docker run --rm -it vox-builder-test:local bash

# Inside the container, verify installed tools
go version          # Should show Go 1.23
make --version      # Should show Make
git --version       # Should show Git

# Exit the container
exit
```

## Step 3: Test Project Build

### Full Build for All Platforms

```bash
# Run build for all platforms
docker run --rm -v $(pwd):/workspace vox-builder-test:local make all

# Check results
ls -lh dist/
```

Should produce 6 binaries:
- `vox-windows-amd64.exe`
- `vox-windows-arm64.exe`
- `vox-linux-amd64`
- `vox-linux-arm64`
- `vox-darwin-amd64`
- `vox-darwin-arm64`

### Build for Specific Platform

```bash
# Clean previous artifacts
docker run --rm -v $(pwd):/workspace vox-builder-test:local make clean

# Build only for Linux amd64
docker run --rm -v $(pwd):/workspace vox-builder-test:local make linux-amd64

# Check the binary
file dist/vox-linux-amd64
```

## Step 4: Test Running Tests

```bash
# Run tests with virtual display
docker run --rm -v $(pwd):/workspace vox-builder-test:local \
  bash -c "Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 & export DISPLAY=:99.0 && sleep 1 && go test ./..."
```

Or via Makefile:

```bash
docker run --rm -v $(pwd):/workspace vox-builder-test:local make test
```

## Step 5: Interactive Debugging

If something goes wrong, run the container in interactive mode:

```bash
# Start bash in container with mounted project
docker run --rm -it -v $(pwd):/workspace vox-builder-test:local bash

# Inside the container, run commands manually
cd /workspace
go mod download
make linux-amd64
go test ./...

# Exit when done
exit
```

## Step 6: Publish the Image (after successful testing)

If all tests pass successfully, you can publish the image:

```bash
# Retag the image for publishing
docker tag vox-builder-test:local ghcr.io/d-mozulyov/vox-builder:latest

# Login to GitHub Container Registry
echo YOUR_PAT | docker login ghcr.io -u d-mozulyov --password-stdin

# Push the image
docker push ghcr.io/d-mozulyov/vox-builder:latest
```

## Useful Commands

### Docker Cleanup

```bash
# Remove test image
docker rmi vox-builder-test:local

# Remove all unused images
docker image prune -a

# List all images
docker images
```

### Check Image Size

```bash
# View image size
docker images vox-builder-test:local
```

Expected size: ~400-500MB

### Quick Full Verification

```bash
# One command for complete verification
docker build -t vox-builder-test:local -f docker/builder/Dockerfile . && \
docker run --rm -v $(pwd):/workspace vox-builder-test:local make all && \
ls -lh dist/
```

## Common Issues

### Docker not found in WSL

```bash
# Check if Docker daemon is running
sudo service docker start

# Or install Docker Desktop for Windows with WSL2 integration
```

### "permission denied" error when mounting

```bash
# Ensure you're in the project directory
pwd

# Use full path
docker run --rm -v /mnt/c/Projects/Vox:/workspace vox-builder-test:local make all
```

### Build can't find go.mod

```bash
# Ensure you're mounting the project root where go.mod is located
ls go.mod  # Should exist

# Verify working directory in container is correct
docker run --rm -v $(pwd):/workspace vox-builder-test:local ls -la /workspace
```

## Conclusion

After successfully passing all tests, you can be confident that the image works correctly and is ready for publishing. This ensures that CI/CD on GitHub Actions will also work without issues.
