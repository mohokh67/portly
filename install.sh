#!/bin/sh
set -e

REPO="mohokh67/portly"
INSTALL_DIR="/usr/local/bin"
BINARY="portly"

# Detect OS and arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest release tag
LATEST=$(curl -sfL "https://api.github.com/repos/$REPO/releases/latest" \
  | grep '"tag_name"' | sed 's/.*"tag_name": "\(.*\)".*/\1/')

if [ -z "$LATEST" ]; then
  echo "Could not determine latest release. Check https://github.com/$REPO/releases"
  exit 1
fi

URL="https://github.com/$REPO/releases/download/$LATEST/${BINARY}_${LATEST#v}_${OS}_${ARCH}.tar.gz"

echo "Installing portly $LATEST for $OS/$ARCH..."
curl -sfL "$URL" | tar -xz -C /tmp "$BINARY"
install -m 755 /tmp/$BINARY "$INSTALL_DIR/$BINARY"
rm -f /tmp/$BINARY

echo "Installed: $INSTALL_DIR/$BINARY"
portly --version
