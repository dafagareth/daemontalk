#!/usr/bin/env bash
# Fast-track script to sync Markdown content and trigger in-memory reload
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${SCRIPT_DIR}"

echo "[info] Pulling latest markdown content..."
git fetch origin main
git checkout origin/main -- content/

# Fix permissions
chown -R 10001:10001 content 2>/dev/null || true
chmod -R 775 content || true

echo "[info] Triggering Webhook to reload memory..."
# We send a local webhook payload to the container
BODY='{"ref":"refs/heads/main"}'
SECRET="${GITHUB_WEBHOOK_SECRET:-}"

if [ -n "$SECRET" ]; then
    SIG=$(echo -n "$BODY" | openssl dgst -sha256 -hmac "$SECRET" | sed 's/^.* //')
    curl -s -X POST -H "X-Hub-Signature-256: sha256=$SIG" -H "X-GitHub-Event: push" -d "$BODY" http://127.0.0.1:8080/api/webhook/github
else
    curl -s -X POST -H "X-GitHub-Event: push" -d "$BODY" http://127.0.0.1:8080/api/webhook/github
fi
echo ""
echo "[ok] Fast content sync complete!"
