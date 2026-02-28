# Build and Release Scripts

This directory contains scripts for building and releasing Vox.

## Scripts

### release.ps1

PowerShell script for creating and pushing release tags.

**Version format:** `vYY.M.D` or `vYY.M.D.B`
- `YY` - two-digit year (e.g., 26 for 2026)
- `M` - month without leading zero (e.g., 9 for September)
- `D` - day without leading zero (e.g., 2 for 2nd)
- `B` - build number (added if tag already exists, starts from 1)

**Usage:**

```powershell
# Create and push release tag
.\scripts\release.ps1

# Dry run (show what would happen without creating tag)
.\scripts\release.ps1 -DryRun

# Show help
.\scripts\release.ps1 -Help
```

**Examples:**

```
Date: September 2, 2026
First release of the day: v26.9.2
Second release: v26.9.2.1
Third release: v26.9.2.2
```

**What it does:**

1. Determines base version from current date (YY.M.D)
2. Fetches all tags from remote repository
3. Checks if base version tag exists
4. If exists, increments build number until unique tag is found
5. Asks for confirmation
6. Creates annotated git tag
7. Pushes tag to remote
8. GitHub Actions automatically builds and publishes release

### build-all.sh

Bash script for cross-platform compilation.

**Usage:**

```bash
# Build for all platforms
./scripts/build-all.sh

# Build with specific version
VERSION=0.1.0 ./scripts/build-all.sh
```

**Output:** Binaries in `dist/` directory for all 6 platforms:
- Windows (x64, arm64)
- Linux (x64, arm64)
- macOS (x64, arm64)
