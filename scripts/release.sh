
#!/bin/bash

# Usage: ./release.sh v1.0.0 "Release title" "Release description"

VERSION=$1
TITLE=$2
DESCRIPTION=$3

# Build all binaries
echo "Building binaries..."
GOOS=linux GOARCH=amd64 go build -o bin/viren_linux_amd64 ./cmd/viren
GOOS=linux GOARCH=arm64 go build -o bin/viren_linux_arm64 ./cmd/viren
GOOS=darwin GOARCH=amd64 go build -o bin/viren_darwin_amd64 ./cmd/viren
GOOS=darwin GOARCH=arm64 go build -o bin/viren_darwin_arm64 ./cmd/viren
GOOS=windows GOARCH=amd64 go build -o bin/viren_windows_amd64.exe ./cmd/viren

# Create release
echo "Creating release $VERSION..."
gh release create $VERSION \
  --title "$TITLE" \
  --notes "$DESCRIPTION" \
  bin/viren_linux_amd64 \
  bin/viren_linux_arm64 \
  bin/viren_darwin_amd64 \
  bin/viren_darwin_arm64 \
  bin/viren_windows_amd64.exe

echo "✅ Release created successfully!"
echo "🔗 View at: https://github.com/fraol163/viren/releases/tag/$VERSION"
