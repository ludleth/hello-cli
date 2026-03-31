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

INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

echo "Downloading ${DOWNLOAD_URL}..."
curl -sL -o /tmp/hello-cli.tar.gz "$DOWNLOAD_URL"

echo "Extracting..."
tar -xzf /tmp/hello-cli.tar.gz -C /tmp hello-cli

echo "Installing to $INSTALL_DIR/hello-cli..."
mv /tmp/hello-cli "$INSTALL_DIR/hello-cli"
chmod +x "$INSTALL_DIR/hello-cli"

rm /tmp/hello-cli.tar.gz

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
