#!/bin/sh
set -e

REPO="thumbrise/ghset"
BINARY="ghset"
INSTALL_DIR="/usr/local/bin"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

case "$OS" in
  linux|darwin) ;;
  *)
    echo "Unsupported OS: $OS (use Go install or download from Releases)" >&2
    exit 1
    ;;
esac

URL="https://github.com/${REPO}/releases/latest/download/${BINARY}_${OS}_${ARCH}.tar.gz"

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading ${BINARY} for ${OS}/${ARCH}..."
curl -sfL "$URL" -o "${TMPDIR}/${BINARY}.tar.gz"
tar xz -C "$TMPDIR" -f "${TMPDIR}/${BINARY}.tar.gz" "$BINARY"

sudo mkdir -p "$INSTALL_DIR"
sudo mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"

echo "Installed ${BINARY} to ${INSTALL_DIR}/${BINARY}"
