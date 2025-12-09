#!/bin/bash

# Exit on error
set -e

# Ensure we can find 'wails' if it's in the standard Go bin location
export PATH=$PATH:$HOME/go/bin

APP_NAME="craft-launcher"
BUILD_DIR="build/bin"

echo "==========================================="
echo "Building $APP_NAME for all platforms"
echo "==========================================="
echo ""

# Process icons
if [ -f "icons/source/launcher-icon.png" ]; then
    echo "==========================================="
    echo "Processing icons..."
    echo "==========================================="
    node icons/process-icons.js || echo "⚠ Icon processing skipped - continuing build..."
    echo ""
fi

# Generate Integrity Assets
echo "==========================================="
echo "Generating Integrity Assets..."
echo "==========================================="
if [ -f "tools/build_integrity.go" ]; then
    go run tools/build_integrity.go
else
    echo "Warning: Integrity tool not found."
fi
echo ""

# macOS ARM64 (M1/M2/etc)
echo "==========================================="
echo "Building for macOS ARM64 (Apple Silicon)"
echo "==========================================="
wails build -platform darwin/arm64
if [ -d "$BUILD_DIR/$APP_NAME.app" ]; then
    mv "$BUILD_DIR/$APP_NAME.app" "$BUILD_DIR/$APP_NAME-macos-arm64.app"
fi
echo "✓ macOS ARM64 build complete"
echo ""

# macOS x86-64 (Intel)
echo "==========================================="
echo "Building for macOS x86-64 (Intel)"
echo "==========================================="
wails build -platform darwin/amd64
if [ -d "$BUILD_DIR/$APP_NAME.app" ]; then
    mv "$BUILD_DIR/$APP_NAME.app" "$BUILD_DIR/$APP_NAME-macos-amd64.app"
fi
echo "✓ macOS x86-64 build complete"
echo ""

# Windows x86-64
echo "==========================================="
echo "Building for Windows x86-64"
echo "==========================================="
wails build -platform windows/amd64 -o craft-launcher-windows-amd64.exe
echo "✓ Windows x86-64 build complete"
echo ""

# Windows x86 (32-bit)
echo "==========================================="
echo "Building for Windows x86 (32-bit)"
echo "==========================================="
wails build -platform windows/386 -o craft-launcher-windows-386.exe
echo "✓ Windows x86 32-bit build complete"
echo ""

# Windows ARM
echo "==========================================="
echo "Building for Windows ARM"
echo "==========================================="
wails build -platform windows/arm64 -o craft-launcher-windows-arm64.exe
echo "✓ Windows ARM build complete"
echo ""

# Linux x86-64
echo "==========================================="
echo "Building for Linux x86-64"
echo "==========================================="
if wails build -platform linux/amd64 -o craft-launcher-linux-amd64 2>&1; then
    echo "✓ Linux x86-64 build complete"
else
    echo "⚠ Linux x86-64 build skipped (cross-compilation not supported on macOS)"
fi
echo ""

# Linux ARM
echo "==========================================="
echo "Building for Linux ARM"
echo "==========================================="
if wails build -platform linux/arm64 -o craft-launcher-linux-arm64 2>&1; then
    echo "✓ Linux ARM build complete"
else
    echo "⚠ Linux ARM build skipped (cross-compilation not supported on macOS)"
fi
echo ""

# Handle Bundled Data
if [ -d "bundled" ]; then
    echo "==========================================="
    echo "Bundling pre-configured data..."
    echo "==========================================="
    
    # Function to copy bundled data to app/exe/binary location
    bundle_data() {
        DEST="$1"
        if [ -d "$DEST" ] || [ -f "$DEST" ]; then
            # If dest is a file (exe), put data in folder next to it
            TARGET_DIR=$(dirname "$DEST")/data
            # If dest is .app, put in Contents/MacOS/data
            if [[ "$DEST" == *.app ]]; then
                TARGET_DIR="$DEST/Contents/MacOS/data"
            fi
            
            echo "  -> Copying to $TARGET_DIR"
            mkdir -p "$TARGET_DIR"
            cp -r bundled/* "$TARGET_DIR/"
        fi
    }
    
    # Apply to all built artifacts
    ls -1 "$BUILD_DIR" | while read -r file; do
        if [[ "$file" == craft-launcher* || "$file" == *.app ]]; then
             bundle_data "$BUILD_DIR/$file"
        fi
    done
    
    echo "✓ Bundled data included"
    echo ""
fi

echo "==========================================="
echo "All builds complete!"
echo "==========================================="
echo "Build outputs are in: $BUILD_DIR/"
echo ""
echo "Files created:"
ls -1 "$BUILD_DIR" | grep -E '\.(app|exe)$|^craft-launcher-linux' | sed 's/^/  • /'
echo ""
echo "Note: Linux builds require building on a Linux machine"
echo "      or using Docker with a Linux environment."
echo "==========================================="
