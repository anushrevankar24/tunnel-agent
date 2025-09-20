#!/bin/bash

# Build Windows executable from Linux/macOS
# This script cross-compiles the agent for Windows

set -e

echo "Building Windows executable..."

# Ensure we're in the agent directory
cd "$(dirname "$0")"

# Install dependencies
echo "Installing dependencies..."
go mod tidy

# Build for Windows
echo "Cross-compiling for Windows..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o agent.exe \
    ./cmd/agent

# Check if build was successful
if [ -f "agent.exe" ]; then
    echo "✅ Build successful! Windows executable: agent.exe"
    echo "File size: $(du -h agent.exe | cut -f1)"
else
    echo "❌ Build failed!"
    exit 1
fi

echo ""
echo "Usage:"
echo "  ./agent.exe -server ws://localhost:8080/agent/connect?id=my-test -local http://127.0.0.1:9000"
