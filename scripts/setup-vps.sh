#!/usr/bin/env bash
# DaemonTalk VPS Provisioning Script (Debian / Ubuntu)
# Installs: Docker, Docker Compose, Caddy, UFW Firewall

set -euo pipefail

echo "[info] Updating system packages..."
sudo apt-get update && sudo apt-get upgrade -y
sudo apt-get install -y curl git ufw debian-keyring debian-archive-keyring apt-transport-https

# Install Docker & Docker Compose
if ! command -v docker &> /dev/null; then
    echo "[info] Installing Docker..."
    curl -fsSL https://get.docker.com | sh
    sudo usermod -aG docker "$USER" || true
    sudo systemctl enable --now docker
else
    echo "[ok] Docker is already installed."
fi

# Install Caddy
if ! command -v caddy &> /dev/null; then
    echo "[info] Installing Caddy web server..."
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
    sudo apt-get update
    sudo apt-get install -y caddy
    sudo systemctl enable caddy
else
    echo "[ok] Caddy is already installed."
fi

# Configure Firewall
echo "[info] Configuring UFW firewall rules..."
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp comment 'SSH'
sudo ufw allow 80/tcp comment 'HTTP'
sudo ufw allow 443/tcp comment 'HTTPS'
sudo ufw allow 2222/tcp comment 'DaemonTalk TUI SSH'
sudo ufw --force enable

echo "[ok] VPS setup completed successfully."
echo "[next] Clone repository to /opt/daemontalk and run: docker compose up -d --build"
