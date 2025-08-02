#!/bin/bash

# Release script for must_go package

set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

VERSION=$1

# Validate version format
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must be in format vX.Y.Z (e.g., v1.0.0)"
    exit 1
fi

echo "Releasing version $VERSION..."

# Run tests
echo "Running tests..."
go test -v ./pkg/must_go/...

# Build example
echo "Building example..."
cd cmd/example
go build -o example .
cd ../..

# Create and push tag
echo "Creating tag $VERSION..."
git tag -a $VERSION -m "Release $VERSION"
git push origin $VERSION

echo "Release $VERSION created successfully!"
echo "GitHub Actions will automatically create a release when the tag is pushed." 