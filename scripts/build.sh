#!/bin/bash

set -e

BINARY_NAME="ccpm"
VERSION=${VERSION:-"dev"}
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo "Building CCPM v$VERSION"
echo "Build time: $BUILD_TIME"
echo "Commit: $COMMIT_HASH"
echo ""

# Create build directory
mkdir -p build

# Build for different platforms
platforms=(
  "linux/amd64"
  "linux/arm64"
  "linux/386"
  "windows/amd64"
  "windows/386"
  "darwin/amd64"
  "darwin/arm64"
)

for platform in "${platforms[@]}"; do
  os="${platform%/*}"
  arch="${platform#*/}"
  output_name="build/$BINARY_NAME-$os-$arch"

  if [ "$os" = "windows" ]; then
    output_name="$output_name.exe"
  fi

  echo "Building for $os/$arch..."

  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build \
    -ldflags="-s -w -X 'main.Version=$VERSION' -X 'main.BuildTime=$BUILD_TIME' -X 'main.CommitHash=$COMMIT_HASH'" \
    -o "$output_name" \
    .

  # Check binary size
  size=$(stat -c%s "$output_name" 2>/dev/null || stat -f%z "$output_name")
  size_mb=$(echo "scale=2; $size / 1048576" | bc)
  echo "  Size: $size_mb MB ($size bytes)"

  # Check if binary is under 5MB limit
  if [ "$size" -gt 5242880 ]; then
    echo "  WARNING: Binary size exceeds 5MB limit"
  fi

  # Make binary executable
  if [ "$os" != "windows" ]; then
    chmod +x "$output_name"
  fi

  echo "  Built: $output_name"
  echo ""
done

# Generate checksums
echo "Generating checksums..."
cd build
sha256sum * > checksums.txt
cd ..

echo "Build complete! Binaries are in the build/ directory."
echo "Checksums: build/checksums.txt"