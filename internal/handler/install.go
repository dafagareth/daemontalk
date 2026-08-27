package handler

import (
	"net/http"
)

const installScriptContent = `#!/bin/sh
# DaemonTalk CLI & TUI Installer
# Usage: curl -sL https://daemontalk.com/install.sh | bash

set -e

echo "------------------------------------------------------"
echo " ⚡ DaemonTalk CLI & Interactive TUI Installer"
echo "------------------------------------------------------"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "✓ Detected platform: $OS ($ARCH)"

INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

TARGET="$INSTALL_DIR/daemontalk"
DOWNLOAD_URL="https://daemontalk.com/download/cli/${OS}_${ARCH}/daemontalk"

echo "📦 Fetching latest DaemonTalk binary..."
if curl -fsSL "$DOWNLOAD_URL" -o "$TARGET" 2>/dev/null; then
    chmod +x "$TARGET"
    echo "======================================================"
    echo " ✅ Installation Successful!"
    echo " Installed to: $TARGET"
    echo " Run 'daemontalk' in your terminal to start exploring!"
    echo "======================================================"
else
    echo "ℹ️ Direct binary download endpoint is preparing."
    echo "Meanwhile, you can launch the full TUI instantly via SSH:"
    echo "  👉 ssh daemontalk.com -p 2222"
    echo "Or read the daily tech briefing via curl:"
    echo "  👉 curl -sL daemontalk.com/daily"
fi
`

// InstallScript serves the automated shell installer script
func (h *Handler) InstallScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	_, _ = w.Write([]byte(installScriptContent))
}
