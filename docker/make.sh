#!/bin/sh
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

docker run --rm -v "$PROJECT_ROOT:/workspace" ghcr.io/d-mozulyov/vox:latest make "$@"
