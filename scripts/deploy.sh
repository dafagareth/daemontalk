#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
APP_DIR="${APP_DIR:-$SCRIPT_DIR}"
cd "${APP_DIR}"

FRESH_BUILD="false"
for arg in "$@"; do
    if [ "$arg" = "--fresh" ] || [ "$arg" = "-f" ]; then
        FRESH_BUILD="true"
    fi
done

echo "[$(date +'%Y-%m-%d %H:%M:%S')] [deploy] Starting production deployment (fresh_mode: ${FRESH_BUILD})..."

if [ -f "scripts/backup.sh" ]; then
    echo "[info] Creating automatic pre-deploy safety backup..."
    ./scripts/backup.sh || echo "[warn] Pre-deploy backup warning, proceeding..."
fi

echo "[info] Pulling latest commits from git..."
git fetch origin main
git reset --hard origin/main

mkdir -p data backups content/posts web/static/images/posts
chown -R 10001:10001 data content web/static/images/posts 2>/dev/null || true
chmod 750 data || true
chmod -R 775 content web/static/images/posts || true

if [ -f "Caddyfile" ]; then
    echo "[info] Synchronizing Caddy configuration..."
    sudo cp Caddyfile /etc/caddy/Caddyfile
    sudo systemctl reload caddy || true
fi

echo "[info] Pulling latest pre-built image from GHCR..."
docker compose pull web
echo "[info] Restarting container..."
docker compose up -d --force-recreate --remove-orphans
echo "[info] Pruning old images to free up space..."
docker image prune -f

echo "[info] Validating service health..."
sleep 3
HEALTH_STATUS=$(curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8080/healthz || echo "failed")

if [ "${HEALTH_STATUS}" = "200" ]; then
    echo "[ok] Deployment completed successfully. Service status: HTTP 200 OK."
else
    echo "[warn] Service health check returned: ${HEALTH_STATUS}."
    echo "[hint] Inspect container logs: docker compose logs -n 50"
fi
