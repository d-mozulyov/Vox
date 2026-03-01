# Vox Docker Image

Docker image for cross-compiling Vox for all supported platforms.

## Overview

Based on `ghcr.io/powertech-center/alpine-go:latest`, with Vox-specific CGO dependencies:
- ALSA, GTK+3, libayatana-appindicator (audio, UI, system tray)
- Xvfb (headless testing)

## Scripts

Run from the project root:

| Script | Description |
|--------|-------------|
| `./docker/build.sh` | Build the image |
| `./docker/make.sh` | Run `make` in the container (pass targets as args) |
| `./docker/push.sh` | Push the image to ghcr.io |

### Examples

```bash
# Build the image
./docker/build.sh

# Build all platforms
./docker/make.sh all

# Run default make target
./docker/make.sh

# Clean and rebuild
./docker/make.sh clean
./docker/make.sh all

# Push to registry (after docker login ghcr.io)
./docker/push.sh
```

## For Contributors and Forks

**Option 1:** Use the official image `ghcr.io/d-mozulyov/vox:latest` — public and ready to use.

**Option 2:** Build your own from `ghcr.io/powertech-center/alpine-go:latest` or `ghcr.io/d-mozulyov/vox:latest`, customize `docker/Dockerfile`, and publish to your registry.

## For Maintainers

1. Build: `./docker/build.sh`
2. Test: `./docker/make.sh all` → check `dist/`
3. Login: `docker login ghcr.io -u d-mozulyov --password-stdin`
4. Push: `./docker/push.sh`

Ensure the package is set to public visibility in GitHub package settings.

## Troubleshooting

**"cannot execute: required file not found"** — usually caused by CRLF line endings. Fix:

```bash
sed -i 's/\r$//' docker/build.sh docker/make.sh docker/push.sh
```

Or re-checkout to apply `.gitattributes`:

```bash
git checkout -- docker/*.sh
```
