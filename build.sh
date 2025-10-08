#!/bin/bash

# Create a builds directory if it doesn't exist
mkdir -p builds

# Check if ImageMagick is installed
if ! command -v convert &> /dev/null; then
    echo "⚠️  ImageMagick not found. Icon conversion will be skipped."
    echo "To install ImageMagick: brew install imagemagick"
else
    # Set up temporary directory for icon conversion
    TEMP_ICONSET="./icon.iconset"
    mkdir -p "$TEMP_ICONSET"

    # Check if the source icon exists
    if [ -f "img/logo.ico" ]; then
        echo "Converting icon..."
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
        echo "✅ Icon conversion complete"
    else
        echo "⚠️  Icon file 'img/logo.ico' not found. Icon conversion will be skipped."
    fi

    # Clean up temporary directory
    rm -rf "$TEMP_ICONSET"
fi

# Create macOS app bundle directory structure
mkdir -p builds/DEXSimulator.app/Contents/MacOS
mkdir -p builds/DEXSimulator.app/Contents/Resources

# Copy Info.plist and icon to the app bundle
cp Info.plist builds/DEXSimulator.app/Contents/
cp icon.icns builds/DEXSimulator.app/Contents/Resources/

# Build for macOS (Intel/AMD64)
echo "Building for macOS (Intel/AMD64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.Version=1.0.2'" -o builds/dex-simulator-mac
echo "✅ macOS (Intel/AMD64) build complete"

# Create a copy for the app bundle
cp builds/dex-simulator-mac builds/DEXSimulator.app/Contents/MacOS/dex-simulator

# Build for macOS (ARM64 - M1/M2)
echo "Building for macOS (ARM64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "-X 'main.Version=1.0.2'" -o builds/dex-simulator-mac-arm64
echo "✅ macOS (ARM64) build complete"

# Build for Windows (AMD64) - Console version (for command line use)
echo "Building Windows Console version (AMD64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.Version=1.0.2' -w -s" -trimpath -o builds/dex-simulator-console.exe
echo "✅ Windows Console build complete"

# Build for Windows (AMD64) - GUI version (for double-click use)
echo "Building Windows GUI version (AMD64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.Version=1.0.2' -H=windowsgui -w -s" -trimpath -o builds/dex-simulator.exe
echo "✅ Windows GUI build complete"

# Attempt to use UPX to compress the Windows executables if available
if command -v upx &> /dev/null; then
    echo "Compressing Windows executables with UPX..."
    upx -9 builds/dex-simulator.exe
    upx -9 builds/dex-simulator-console.exe
    echo "✅ UPX compression complete"
else
    echo "⚠️  UPX not found. To install: brew install upx"
    echo "    Skipping Windows executable compression"
fi

# Create Windows installer ZIP to help avoid antivirus warnings
echo "Creating Windows ZIP packages..."
cd builds
# Create complete package with both versions
zip -r dex-simulator-windows-complete.zip dex-simulator.exe dex-simulator-console.exe
# Create separate packages for GUI and Console versions
zip -r dex-simulator-windows-gui.zip dex-simulator.exe
zip -r dex-simulator-windows-console.zip dex-simulator-console.exe
# Create packages with generic filenames to help avoid detection
cp dex-simulator.exe app-gui.exe
cp dex-simulator-console.exe app-console.exe
zip -r app-simulator-complete.zip app-gui.exe app-console.exe
zip -r app-simulator-gui.zip app-gui.exe
zip -r app-simulator-console.zip app-console.exe
rm app-gui.exe app-console.exe
cd ..
echo "✅ Windows ZIP packages created"

# Create universal binary for macOS app bundle if both builds succeeded
if [ -f "builds/dex-simulator-mac" ] && [ -f "builds/dex-simulator-mac-arm64" ]; then
    echo "Creating universal binary for app bundle..."
    lipo -create -output builds/DEXSimulator.app/Contents/MacOS/dex-simulator-universal builds/dex-simulator-mac builds/dex-simulator-mac-arm64
    mv builds/DEXSimulator.app/Contents/MacOS/dex-simulator-universal builds/DEXSimulator.app/Contents/MacOS/dex-simulator
    echo "✅ Universal binary created"
else
    echo "⚠️  Skipping universal binary creation - one or more builds failed"
    # Copy at least one version if available
    if [ -f "builds/dex-simulator-mac" ]; then
        cp builds/dex-simulator-mac builds/DEXSimulator.app/Contents/MacOS/dex-simulator
    elif [ -f "builds/dex-simulator-mac-arm64" ]; then
        cp builds/dex-simulator-mac-arm64 builds/DEXSimulator.app/Contents/MacOS/dex-simulator
    fi
fi

echo "\nAll builds completed! The binaries are in the 'builds' directory:"
echo "- builds/dex-simulator-mac         (macOS Intel)"
echo "- builds/dex-simulator-mac-arm64   (macOS Apple Silicon)"
echo "- builds/dex-simulator.exe         (Windows GUI version)"
echo "- builds/dex-simulator-console.exe (Windows Console version)"
echo "\nWindows ZIP packages created:"
echo "- builds/dex-simulator-windows-complete.zip  (Both GUI and Console versions)"
echo "- builds/dex-simulator-windows-gui.zip      (GUI version only)"
echo "- builds/dex-simulator-windows-console.zip  (Console version only)"
echo "- builds/app-simulator-*.zip                (Generic filename versions)"