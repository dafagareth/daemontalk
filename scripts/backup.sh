#!/usr/bin/env bash
# DaemonTalk Production Backup Script
# Safely backs up SQLite databases (WAL mode safe) and markdown content.
# Keeps the last 7 daily archives.
# Usage: ./scripts/backup.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKUP_DIR="${BACKUP_DIR:-$SCRIPT_DIR/backups}"
DATA_DIR="${DATA_DIR:-$SCRIPT_DIR/data}"
CONTENT_DIR="${CONTENT_DIR:-$SCRIPT_DIR/content}"
TIMESTAMP="$(date +'%Y%m%d_%H%M%S')"
ARCHIVE_NAME="daemontalk_backup_${TIMESTAMP}.tar.gz"

mkdir -p "${BACKUP_DIR}"

echo "[$(date +'%Y-%m-%d %H:%M:%S')] [backup] Starting backup procedure..."

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

# Safe SQLite database dump
if [ -d "${DATA_DIR}" ]; then
    mkdir -p "${TMP_DIR}/data"
    for db in "${DATA_DIR}"/*.db; do
        if [ -f "$db" ]; then
            db_name="$(basename "$db")"
            if command -v sqlite3 >/dev/null 2>&1; then
                sqlite3 "$db" ".backup '${TMP_DIR}/data/${db_name}'"
            else
                cp "$db"* "${TMP_DIR}/data/" 2>/dev/null || true
            fi
        fi
    done
fi

# Copy markdown articles
if [ -d "${CONTENT_DIR}" ]; then
    cp -r "${CONTENT_DIR}" "${TMP_DIR}/content"
fi

# Compress into tar.gz
tar -czf "${BACKUP_DIR}/${ARCHIVE_NAME}" -C "${TMP_DIR}" .

echo "[ok] Backup created: ${BACKUP_DIR}/${ARCHIVE_NAME} ($(du -h "${BACKUP_DIR}/${ARCHIVE_NAME}" | cut -f1))"

# Rotate archives older than 7 days
find "${BACKUP_DIR}" -name "daemontalk_backup_*.tar.gz" -type f -mtime +7 -exec rm -f {} +
echo "[ok] Old backup rotation complete."
