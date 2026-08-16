#!/usr/bin/env bash
# DaemonTalk Production Update & Deployment Script
# Usage:
#   ./scripts/deploy.sh         Standard fast deployment (cached layers)
#   ./scripts/deploy.sh --fresh 100% clean rebuild without Docker cache (force recreate)

set -euo pipefail

APP_DIR="${APP_DIR:-/opt/daemontalk}"
cd "${APP_DIR}"

FRESH_BUILD="false"
for arg in "$@"; do
    if [ "$arg" = "--fresh" ] || [ "$arg" = "-f" ]; then
        FRESH_BUILD="true"
    fi
done

echo "[$(date +'%Y-%m-%d %H:%M:%S')] [deploy] Starting production deployment (fresh_mode: ${FRESH_BUILD})..."

# 1. Pre-deployment automatic safety backup
if [ -f "scripts/backup.sh" ]; then
    echo "[info] Creating automatic pre-deploy safety backup..."
    ./scripts/backup.sh || echo "[warn] Pre-deploy backup warning, proceeding..."
fi

# 2. Pull latest code from GitHub
echo "[info] Pulling latest commits from git..."
git fetch origin main
git reset --hard origin/main

# 3. Ensure persistent data directory and permissions
mkdir -p data backups
chmod 775 data || true

# 4. Synchronize Caddy configuration if changed
if [ -f "Caddyfile" ]; then
    echo "[info] Synchronizing Caddy configuration..."
    sudo cp Caddyfile /etc/caddy/Caddyfile
    sudo systemctl reload caddy || true
fi

# 5. Build and launch container
if [ "${FRESH_BUILD}" = "true" ]; then
    echo "[info] Performing fresh no-cache build..."
    docker compose build --no-cache
    docker compose up -d --force-recreate --remove-orphans
    echo "[info] Pruning Docker builder cache and unused images..."
    docker builder prune -f
    docker image prune -f
else
    echo "[info] Rebuilding and launching container with cached layers..."
    docker compose up -d --build --remove-orphans
    docker image prune -f
fi

# 6. Validate service health
echo "[info] Validating service health..."
sleep 3
HEALTH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8080/healthz || echo "failed")

if [ "${HEALTH_STATUS}" = "200" ]; then
    echo "[ok] Deployment completed successfully. Service status: HTTP 200 OK."
else
    echo "[warn] Service health check returned: ${HEALTH_STATUS}."
    echo "[hint] Inspect container logs: docker compose logs -n 50"
fi
