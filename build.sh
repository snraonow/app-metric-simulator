#!/bin/bash

# Create a builds directory if it doesn't exist
mkdir -p builds

# Set up temporary directory for icon conversion
TEMP_ICONSET="./icon.iconset"
mkdir -p "$TEMP_ICONSET"

# Convert ICO to PNG files at different sizes
convert img/logo.ico -resize 16x16 "$TEMP_ICONSET/icon_16x16.png"
convert img/logo.ico -resize 32x32 "$TEMP_ICONSET/icon_16x16@2x.png"
convert img/logo.ico -resize 32x32 "$TEMP_ICONSET/icon_32x32.png"
convert img/logo.ico -resize 64x64 "$TEMP_ICONSET/icon_32x32@2x.png"
convert img/logo.ico -resize 128x128 "$TEMP_ICONSET/icon_128x128.png"
convert img/logo.ico -resize 256x256 "$TEMP_ICONSET/icon_128x128@2x.png"
convert img/logo.ico -resize 256x256 "$TEMP_ICONSET/icon_256x256.png"
convert img/logo.ico -resize 512x512 "$TEMP_ICONSET/icon_256x256@2x.png"
convert img/logo.ico -resize 512x512 "$TEMP_ICONSET/icon_512x512.png"
convert img/logo.ico -resize 1024x1024 "$TEMP_ICONSET/icon_512x512@2x.png"

# Convert iconset to ICNS
iconutil -c icns "$TEMP_ICONSET"

# Clean up temporary directory
rm -rf "$TEMP_ICONSET"

# Create macOS app bundle directory structure
mkdir -p builds/DEXSimulator.app/Contents/MacOS
mkdir -p builds/DEXSimulator.app/Contents/Resources

# Copy Info.plist and icon to the app bundle
cp Info.plist builds/DEXSimulator.app/Contents/
cp icon.icns builds/DEXSimulator.app/Contents/Resources/

# Build for macOS (amd64)
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.Version=1.0.0'" -o builds/DEXSimulator.app/Contents/MacOS/dex-simulator
cp builds/DEXSimulator.app/Contents/MacOS/dex-simulator builds/dex-simulator-mac
echo "✅ macOS (amd64) build complete"

# Build for macOS (arm64 - M1/M2)
echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.Version=1.0.0'" -o builds/dex-simulator-mac-arm64
echo "✅ macOS (arm64) build complete"

# Build for Windows (amd64)
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.Version=1.0.0'" -o builds/dex-simulator.exe
echo "✅ Windows build complete"

echo "\nAll builds completed! The binaries are in the 'builds' directory:"
echo "- builds/dex-simulator-mac     (macOS Intel)"
echo "- builds/dex-simulator-mac-arm64 (macOS Apple Silicon)"
echo "- builds/dex-simulator.exe     (Windows)"