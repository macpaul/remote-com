#!/bin/bash

# Remote-COM Release Automation Script
# Usage: ./release.sh v0.1.0

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Error: Version tag required (e.g., ./release.sh v0.1.0)"
    exit 1
fi

# Ensure build directory exists
mkdir -p build/win64/bin

echo "Building Remote-COM for Windows (win64)..."
/Users/macpaul/go/bin/wails build -platform windows/amd64 -o ../win64/bin/remote-com.exe

if [ $? -ne 0 ]; then
    echo "Error: Build failed."
    exit 1
fi

echo "Zipping the executable..."
ZIP_NAME="win64-remote-com-$VERSION.zip"
zip -j "build/win64/bin/$ZIP_NAME" "build/win64/bin/remote-com.exe"

if [ $? -ne 0 ]; then
    echo "Error: Zipping failed."
    exit 1
fi

echo "Creating GitHub release $VERSION..."
# Uploading both the zip and the raw exe, both with win64 prefix
gh release create "$VERSION" \
    "./build/win64/bin/$ZIP_NAME" \
    "./build/win64/bin/remote-com.exe#win64-remote-com.exe" \
    --title "Release $VERSION" \
    --notes "Automated release $VERSION of Remote-COM."

if [ $? -eq 0 ]; then
    echo "Success: Release $VERSION published."
else
    echo "Error: Release failed. Make sure you are logged in via 'gh auth login'."
fi
