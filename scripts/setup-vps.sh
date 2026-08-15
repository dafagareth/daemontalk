#!/usr/bin/env bash
# ==============================================================================
# DaemonTalk VPS Production Setup Script (Debian / Ubuntu)
# Configures: Docker, Docker Compose, Caddy (Automatic SSL), UFW Firewall
# ==============================================================================

set -euo pipefail

echo "======================================================"
echo " 🚀 DaemonTalk Production VPS Setup"
echo "======================================================"

# 1. Update system packages
echo "📦 Updating system packages..."
sudo apt-get update && sudo apt-get upgrade -y
sudo apt-get install -y curl git ufw debian-keyring debian-archive-keyring apt-transport-https

# 2. Install Docker & Docker Compose
if ! command -v docker &> /dev/null; then
    echo "🐳 Installing Docker..."
    curl -fsSL https://get.docker.com | sh
    sudo usermod -aG docker "$USER" || true
    sudo systemctl enable --now docker
else
    echo "✓ Docker is already installed."
fi

# 3. Install Caddy Web Server
if ! command -v caddy &> /dev/null; then
    echo "🔒 Installing Caddy (Automatic HTTPS)..."
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
    sudo apt-get update
    sudo apt-get install -y caddy
    sudo systemctl enable caddy
else
    echo "✓ Caddy is already installed."
fi

# 4. Configure UFW Firewall
echo "🛡️ Configuring Firewall (UFW)..."
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp comment 'SSH'
sudo ufw allow 80/tcp comment 'HTTP'
sudo ufw allow 443/tcp comment 'HTTPS'
sudo ufw allow 2222/tcp comment 'DaemonTalk TUI SSH'
sudo ufw --force enable

echo "======================================================"
echo " ✅ VPS Prerequisites Installed Successfully!"
echo " Next Steps:"
echo " 1. Clone your repo: git clone <your-repo-url> daemontalk"
echo " 2. Copy Caddyfile: sudo cp daemontalk/Caddyfile /etc/caddy/Caddyfile"
echo " 3. Reload Caddy:   sudo systemctl reload caddy"
echo " 4. Start Docker:   cd daemontalk && docker compose up -d --build"
echo "======================================================"
