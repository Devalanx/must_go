# Deployment Guide

This guide explains how to deploy the `must_go` package to GitHub so it can be used with `go get`.

## Prerequisites

1. A GitHub account
2. Git installed on your machine
3. Go 1.24+ installed

## Steps to Deploy

### 1. Create GitHub Repository

1. Go to [GitHub](https://github.com) and create a new repository
2. Name it `must_go`
3. Make it public
4. Don't initialize with README, .gitignore, or license (we already have these)

### 2. Initialize Git and Push

```bash
# Initialize git repository
git init

# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: must_go package"

# Add remote origin (replace with your GitHub username)
git remote add origin https://github.com/Devalanx/must_go.git

# Push to main branch
git push -u origin main
```

### 3. Create First Release

```bash
# Create and push a tag for the first release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### 4. Verify Installation

After pushing to GitHub, you can test the installation:

```bash
# Create a test directory
mkdir test-install
cd test-install

# Initialize a new Go module
go mod init test

# Try to get the package
go get github.com/Devalanx/must_go

# Create a test file
cat > main.go << 'EOF'
package main

import (
    "fmt"
    "net/http"
    "github.com/Devalanx/must_go/pkg/must_go"
)

func main() {
    fmt.Println("Package imported successfully!")
    
    // Test basic functionality
    mux := http.NewServeMux()
    handler := must_go.RecoveryMiddleware(mux)
    fmt.Printf("Handler type: %T\n", handler)
}
EOF

# Run the test
go run main.go
```

## Using the Package

Once deployed, users can install and use the package:

```bash
# Install the package
go get github.com/Devalanx/must_go

# Use in your code
import "github.com/Devalanx/must_go/pkg/must_go"
```

## Future Releases

To create new releases:

1. Make your changes
2. Update the version in `CHANGELOG.md`
3. Commit and push changes
4. Create a new tag:
   ```bash
   git tag -a v1.1.0 -m "Release v1.1.0"
   git push origin v1.1.0
   ```

Or use the provided script:
```bash
./scripts/release.sh v1.1.0
```

## GitHub Actions

The repository includes GitHub Actions workflows that will:

- Run tests on every push and pull request
- Create releases automatically when tags are pushed
- Run linting checks

## Package Discovery

Once published, your package will be available on:

- [pkg.go.dev](https://pkg.go.dev/github.com/Devalanx/must_go)
- [Go Report Card](https://goreportcard.com/report/github.com/Devalanx/must_go)

## Troubleshooting

### If `go get` fails:

1. Make sure the repository is public
2. Verify the module name in `go.mod` matches the GitHub path
3. Check that the tag exists and is pushed
4. Try clearing the module cache: `go clean -modcache`

### If tests fail:

1. Run `go mod tidy` to ensure dependencies are correct
2. Check that all import paths are updated to use the GitHub path
3. Verify Go version compatibility

## Maintenance

- Keep the CHANGELOG.md updated
- Respond to issues and pull requests
- Maintain backward compatibility for patch releases
- Update documentation as needed 