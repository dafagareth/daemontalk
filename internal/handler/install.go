package handler

import (
	"net/http"
)

const installScriptContent = `#!/bin/sh
# DaemonTalk CLI & TUI Installer
# Usage: curl -sL https://daemontalk.com/install.sh | bash

set -e

echo "------------------------------------------------------"
echo " DaemonTalk CLI & Interactive TUI Installer"
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
        ARCH="amd64"
        ;;
esac

echo "Detected platform: $OS ($ARCH)"

INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ] 2>/dev/null; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

TARGET="$INSTALL_DIR/daemontalk"
DOWNLOAD_URL="https://daemontalk.com/download/cli/${OS}_${ARCH}/daemontalk"

installed=0

# Attempt direct standalone binary download
if curl -fsSL "$DOWNLOAD_URL" -o "$TARGET" 2>/dev/null && [ -s "$TARGET" ]; then
    chmod +x "$TARGET"
    installed=1
fi

# Fallback: Install instant SSH TUI launcher
if [ "$installed" -eq 0 ]; then
    cat << 'EOF' > "$TARGET"
#!/bin/sh
# DaemonTalk TUI Launcher
if command -v ssh >/dev/null 2>&1; then
    exec ssh -t ssh.daemontalk.com -p 2222 "$@"
else
    echo "Error: OpenSSH client (ssh) is required to run DaemonTalk." >&2
    exit 1
fi
EOF
    chmod +x "$TARGET"
    installed=1
fi

echo "-------------------------------------------------------"
echo " Installation Successful!"
echo " Installed to: $TARGET"

# Check if INSTALL_DIR is in PATH
case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *)
        echo ""
        echo " Note: $INSTALL_DIR is not currently in your PATH."
        echo " Add it to your shell configuration (~/.zshrc or ~/.bashrc):"
        echo "   export PATH=\"$INSTALL_DIR:\$PATH\""
        echo ""
        ;;
esac

echo " Run 'daemontalk' in your terminal to start exploring!"
echo "-------------------------------------------------------"
`

// InstallScript serves the automated shell installer script
func (h *Handler) InstallScript(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	_, _ = w.Write([]byte(installScriptContent))
}
