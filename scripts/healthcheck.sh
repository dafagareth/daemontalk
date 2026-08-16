#!/usr/bin/env bash
# DaemonTalk Production Health Watchdog
# Checks HTTP /healthz and auto-restarts container if unresponsive.
# Usage: ./scripts/healthcheck.sh

set -euo pipefail

APP_DIR="${APP_DIR:-/opt/daemontalk}"
LOG_FILE="/var/log/daemontalk_health.log"
HEALTH_URL="http://127.0.0.1:8080/healthz"
MAX_RETRIES=2
RETRY_COUNT=0

log_msg() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" | sudo tee -a "${LOG_FILE}" > /dev/null
}

check_health() {
    curl -fsSL -m 5 "${HEALTH_URL}" > /dev/null 2>&1
}

if ! check_health; then
    log_msg "[warn] Health check failed for ${HEALTH_URL}. Retrying..."
    while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
        sleep 2
        RETRY_COUNT=$((RETRY_COUNT + 1))
        if check_health; then
            log_msg "[ok] Service recovered on retry ${RETRY_COUNT}."
            exit 0
        fi
    done

    log_msg "[alert] Service unresponsive. Initiating auto-recovery container restart..."
    cd "${APP_DIR}"
    docker compose restart web
    sleep 4

    if check_health; then
        log_msg "[ok] Service auto-recovered and restored to healthy state."
    else
        log_msg "[error] Auto-recovery restart failed. Manual intervention required."
    fi
fi
