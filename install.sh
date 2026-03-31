#!/bin/sh
set -e

REPO="ludleth/hello-cli"
GITHUB_URL="https://github.com/${REPO}/releases/latest/download"

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Linux)  OS_NAME="linux" ;;
  Darwin) OS_NAME="darwin" ;;
  *)      echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64)  ARCH_NAME="amd64" ;;
  aarch64) ARCH_NAME="arm64" ;;
  arm64)   ARCH_NAME="arm64" ;;
  *)       echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

TARBALL="hello-cli_${OS_NAME}_${ARCH_NAME}.tar.gz"
DOWNLOAD_URL="${GITHUB_URL}/${TARBALL}"
CHECKSUM_URL="${GITHUB_URL}/${TARBALL}.sha256"

INSTALL_DIR="$HOME/.local/bin"
TMP_DIR="$(mktemp -d)" || { echo "Error: failed to create temp directory"; exit 1; }
cleanup() {
  if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ] && case "$TMP_DIR" in /tmp/*|/var/folders/*) true;; *) false;; esac; then
    rm -rf "$TMP_DIR"
  fi
}
trap cleanup EXIT

mkdir -p "$INSTALL_DIR"

echo "Downloading ${TARBALL}..."
curl -sL -o "$TMP_DIR/$TARBALL" "$DOWNLOAD_URL"
curl -sL -o "$TMP_DIR/${TARBALL}.sha256" "$CHECKSUM_URL"

echo "Verifying checksum..."
EXPECTED=$(awk '{print $1}' "$TMP_DIR/${TARBALL}.sha256")
if [ -z "$EXPECTED" ]; then
  echo "Error: checksum file is empty"; exit 1
fi
ACTUAL=$(shasum -a 256 "$TMP_DIR/$TARBALL" | awk '{print $1}')
if [ "$EXPECTED" != "$ACTUAL" ]; then
  echo "Error: checksum mismatch (expected $EXPECTED, got $ACTUAL)"; exit 1
fi
echo "Checksum verified."

echo "Extracting..."
tar -xzf "$TMP_DIR/$TARBALL" -C "$TMP_DIR" hello-cli

echo "Installing to $INSTALL_DIR/hello-cli..."
mv "$TMP_DIR/hello-cli" "$INSTALL_DIR/hello-cli"
chmod +x "$INSTALL_DIR/hello-cli"

# Detect user's login shell via $SHELL env var
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
  echo "Adding $INSTALL_DIR to PATH..."

  case "$SHELL" in
    */zsh)  SHELL_RC="$HOME/.zshrc" ;;
    */bash) SHELL_RC="$HOME/.bashrc" ;;
    *)      SHELL_RC="$HOME/.profile" ;;
  esac

  echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$SHELL_RC"
  echo "Added to $SHELL_RC. Please restart your shell or run 'source $SHELL_RC'."
fi

echo "Installation complete. Run 'hello-cli' to test."
