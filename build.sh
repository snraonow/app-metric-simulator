#!/bin/bash

# Create a builds directory if it doesn't exist
mkdir -p builds

# Build for macOS (amd64)
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o builds/dex-simulator-mac
echo "✅ macOS (amd64) build complete"

# Build for macOS (arm64 - M1/M2)
echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o builds/dex-simulator-mac-arm64
echo "✅ macOS (arm64) build complete"

# Build for Windows (amd64)
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o builds/dex-simulator.exe
echo "✅ Windows build complete"

echo "\nAll builds completed! The binaries are in the 'builds' directory:"
echo "- builds/dex-simulator-mac     (macOS Intel)"
echo "- builds/dex-simulator-mac-arm64 (macOS Apple Silicon)"
echo "- builds/dex-simulator.exe     (Windows)"