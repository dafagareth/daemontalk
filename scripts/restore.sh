#!/usr/bin/env bash
# DaemonTalk Production Restore Script
# Safely restores SQLite databases and content from a backup archive.
# Usage: ./scripts/restore.sh /path/to/archive.tar.gz

set -euo pipefail

if [ $# -lt 1 ]; then
    echo "Usage: $0 <path_to_archive.tar.gz>"
    exit 1
fi

ARCHIVE="$1"

if [ ! -f "${ARCHIVE}" ]; then
    echo "[error] Backup archive '${ARCHIVE}' not found."
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEST_DIR="${DEST_DIR:-$SCRIPT_DIR}"
TIMESTAMP="$(date +'%Y%m%d_%H%M%S')"

echo "[restore] Target archive: ${ARCHIVE}"
echo "[restore] Destination: ${DEST_DIR}"

read -p "Are you sure you want to proceed with restore? (y/N): " -r CONFIRM
if [[ ! "${CONFIRM}" =~ ^[Yy]$ ]]; then
    echo "[info] Restore aborted."
    exit 0
fi

# Stop container to avoid file locks
echo "[info] Stopping Docker container..."
cd "${DEST_DIR}"
docker compose down || true

# Pre-restore safety copy
if [ -d "${DEST_DIR}/data" ]; then
    echo "[info] Creating safety copy at data_prerestore_${TIMESTAMP}..."
    cp -r "${DEST_DIR}/data" "${DEST_DIR}/data_prerestore_${TIMESTAMP}"
fi

# Extract archive
echo "[info] Extracting backup archive..."
TMP_EXTRACT="$(mktemp -d)"
trap 'rm -rf "${TMP_EXTRACT}"' EXIT

tar -xzf "${ARCHIVE}" -C "${TMP_EXTRACT}"

if [ -d "${TMP_EXTRACT}/data" ]; then
    mkdir -p "${DEST_DIR}/data"
    cp -r "${TMP_EXTRACT}/data/"* "${DEST_DIR}/data/" 2>/dev/null || true
fi

if [ -d "${TMP_EXTRACT}/content" ]; then
    cp -r "${TMP_EXTRACT}/content/"* "${DEST_DIR}/content/" 2>/dev/null || true
fi

# Restart container
echo "[info] Restarting Docker container..."
docker compose up -d --build

echo "[ok] Restore completed successfully."
