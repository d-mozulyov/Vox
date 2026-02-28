# Docker Build Environment

This document explains how the project uses Docker for CI/CD optimization and how to customize it for your needs.

## Overview

The project uses a pre-built Docker image (`ghcr.io/d-mozulyov/vox-builder:latest`) that contains all build dependencies:
- Go 1.21.6
- musl-tools and musl-dev for static Linux builds
- aarch64-linux-musl-cross for ARM64 cross-compilation
- All required development libraries (ALSA, X11, AppIndicator, etc.)
- Xvfb for headless testing

This approach speeds up CI/CD by 3-4x (from ~10-15 minutes to ~3-5 minutes).

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
RUN apt-get update && apt-get install -y \
    your-custom-package \
    && rm -rf /var/lib/apt/lists/*
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

When you need to update dependencies (rare - once every few months):

### Manual Build and Publish

1. Edit `docker/builder/Dockerfile` to update dependencies
2. Build the image locally:
   ```bash
   docker build -t ghcr.io/d-mozulyov/vox-builder:latest -f docker/builder/Dockerfile .
   ```
3. Test the image (optional but recommended):
   ```bash
   docker run --rm -it ghcr.io/d-mozulyov/vox-builder:latest bash
   go version
   musl-gcc --version
   ```
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
   git commit -m "chore: update Docker build dependencies"
   git push origin main
   ```

**Why manual?** 
- Simple and explicit control
- No complex automation for rare operations
- Dockerfile changes happen once every few months
- Follows KISS principle

### For Forks

Same process, but replace `d-mozulyov` with your GitHub username.

## Local Development with Docker

You can use the Docker image locally to ensure consistency with CI:

```bash
# Build the project
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest \
  go build -o vox ./cmd/vox

# Run tests
docker run --rm -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest \
  bash -c "Xvfb :99 -screen 0 1024x768x24 > /dev/null 2>&1 & export DISPLAY=:99.0 && sleep 1 && go test ./..."

# Interactive shell
docker run --rm -it -v $(pwd):/workspace ghcr.io/d-mozulyov/vox-builder:latest bash
```

## Files Related to Docker Build

- `docker/builder/Dockerfile` - Docker image definition
- `docker/builder/README.md` - Detailed Docker image documentation
- `.github/workflows/build.yml` - Main CI/CD workflow (uses the image)

To update the image, see the "Updating the Docker Image" section above.

## Troubleshooting

### GitHub Actions can't pull the image

- Ensure the image is public (see Step 3 above)
- Check the image name in workflow files
- Verify the image exists: `https://github.com/USERNAME?tab=packages`

### Local docker push fails with "denied"

- Check your PAT has `write:packages` scope
- Verify you're logged in: `docker login ghcr.io`
- Ensure the image name matches: `ghcr.io/USERNAME/vox-builder`

### Image is too large

Current size: ~1.5-2GB (compressed: ~600-800MB)

This is normal for a full build environment. GitHub Actions caches the image, so download time is minimal (~10-20 seconds).

To reduce size:
- Use multi-stage builds
- Remove unnecessary packages
- Combine RUN commands to reduce layers

## Benefits of This Approach

✓ **Speed**: 3-4x faster CI/CD (no dependency installation on each run)  
✓ **Consistency**: Same environment everywhere (CI, local, contributors)  
✓ **Reliability**: No risk of apt-get failures or version mismatches  
✓ **Caching**: Go modules and build cache further speed up builds  
✓ **Flexibility**: Easy to customize for forks  

## Additional Resources

- [GitHub Container Registry Documentation](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Docker Multi-platform Builds](https://docs.docker.com/build/building/multi-platform/)
- [GitHub Actions with Containers](https://docs.github.com/en/actions/using-jobs/running-jobs-in-a-container)

---

**Note:** The Docker image is updated manually when dependencies change (rare). This keeps the process simple and explicit, following the KISS principle.
