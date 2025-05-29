#!/usr/bin/env bash

set -e

APP_NAME="goeditor"
SRC_FILE="main.go"
BUILD_DIR="output"

mkdir -p "$BUILD_DIR"

# Detect platform
OS="$(uname -s)"
ARCH="$(uname -m)"

echo "Detected OS: $OS"
echo "Detected Arch: $ARCH"

# Convert ARCH to GOARCH
case "$ARCH" in
  x86_64)   GOARCH="amd64" ;;
  arm64)    GOARCH="arm64" ;;
  *)        echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Install dependencies
install_deps() {
  echo "Installing dependencies for $1..."
  if [[ "$1" == "Linux" ]]; then
    sudo apt-get update -y
    sudo apt-get install -y golang gcc libgl1-mesa-dev xorg-dev libxkbcommon-dev
  elif [[ "$1" == "Darwin" ]]; then
    xcode-select --install || true
  else
    echo "Dependency installation not supported for OS: $1"
  fi
}

# Build binary for given GOOS and GOARCH
build_target() {
  local GOOS=$1
  local GOARCH=$2
  local OUT_FILE="${BUILD_DIR}/${APP_NAME}-${GOOS}-${GOARCH}"
  [[ "$GOOS" == "windows" ]] && OUT_FILE="${OUT_FILE}.exe"

  echo "Building for $GOOS/$GOARCH..."
  GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=1 go build -o "$OUT_FILE" "$SRC_FILE"
}

# Install local dependencies only
install_deps "$OS"

# Build targets
build_target "darwin" "amd64"
build_target "darwin" "arm64"
build_target "windows" "amd64"
build_target "linux"  "amd64"

echo "✅ Build complete. Binaries are in the '$BUILD_DIR' directory."