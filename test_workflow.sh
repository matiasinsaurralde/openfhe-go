#!/bin/bash

# Test script to simulate GitHub workflow steps
set -e

echo "=== Testing OpenFHE Go Build Workflow ==="

# Set up environment (simulating what the Docker container should provide)
export CGO_ENABLED=1
export GO111MODULE=on

# Check if we have Go
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed"
    exit 1
fi

echo "Go version: $(go version)"

# Check if we have CGO enabled
if [ "$CGO_ENABLED" != "1" ]; then
    echo "ERROR: CGO is not enabled"
    exit 1
fi

echo "CGO_ENABLED: $CGO_ENABLED"

# Check if we can download dependencies
echo "=== Downloading dependencies ==="
go mod download

echo "=== Checking if code compiles without OpenFHE (should fail gracefully) ==="
if go build -v ./... 2>&1 | grep -q "fatal error"; then
    echo "✓ Code correctly fails to compile without OpenFHE libraries (expected)"
else
    echo "✗ Code should fail to compile without OpenFHE libraries"
fi

echo "=== Testing with build tags (should also fail without OpenFHE) ==="
if go build -tags="openfhe" -v ./... 2>&1 | grep -q "fatal error"; then
    echo "✓ Code correctly fails to compile with OpenFHE build tags but no libraries (expected)"
else
    echo "✗ Code should fail to compile without OpenFHE libraries even with build tags"
fi

echo "=== Testing go mod tidy ==="
go mod tidy

echo "=== Testing gofmt ==="
if [ -n "$(gofmt -l .)" ]; then
    echo "✗ Code is not properly formatted"
    gofmt -l .
    exit 1
else
    echo "✓ Code is properly formatted"
fi

echo "=== Testing basic Go file syntax ==="
for file in *.go; do
    if [ -f "$file" ]; then
        echo "Checking syntax of $file..."
        if ! go tool compile -o /dev/null "$file" 2>/dev/null; then
            echo "✗ Syntax error in $file"
        else
            echo "✓ $file has valid syntax"
        fi
    fi
done

echo "=== Workflow simulation completed ==="
echo ""
echo "Next steps:"
echo "1. The GitHub workflow should use the OpenFHE buildbox Docker image"
echo "2. That image should provide OpenFHE headers and libraries"
echo "3. The build should succeed in the Docker environment"
echo "4. Tests should pass with the 'openfhe' build tag"
echo ""
echo "The workflow is configured to:"
echo "- Use ghcr.io/matiasinsaurralde/openfhe-buildbox:latest"
echo "- Set proper environment variables for CGO and OpenFHE"
echo "- Build with -tags=\"openfhe\""
echo "- Run tests with -tags=\"openfhe\""
echo "- Run golangci-lint with timeout"