#!/usr/bin/env bash
# DaemonTalk Content Management - Common Shared Library

set -euo pipefail

POSTS_DIR="${POSTS_DIR:-content/posts}"
IMG_BASE_DIR="${IMG_BASE_DIR:-web/static/images/posts}"

find_post_file() {
    local slug="$1"
    for f in "${POSTS_DIR}"/*.md "${POSTS_DIR}"/*.md.archive; do
        if [ -f "$f" ]; then
            if grep -q "^slug: \?\(['\"]\?\)${slug}\1$" "$f" 2>/dev/null; then
                echo "$f"
                return 0
            fi
        fi
    done
    return 1
}

generate_slug() {
    local title="$1"
    local slug
    slug=$(echo "$title" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-z0-9]+/-/g' | sed -E 's/^-+|-+$//g')
    if [ -z "$slug" ]; then
        slug=$(openssl rand -hex 4)
    fi
    echo "$slug"
}
