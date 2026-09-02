#!/usr/bin/env bash
# DaemonTalk Content Management - Post Lifecycle Module (Publish, Draft, Archive, Delete)

set -euo pipefail

cmd_publish() {
    local slug="${1:-}"
    if [ -z "$slug" ]; then
        echo "[error] Slug argument required. Usage: ./scripts/post.sh publish <slug>"
        exit 1
    fi

    local file
    file=$(find_post_file "$slug" || true)
    if [ -z "$file" ]; then
        echo "[error] Post with slug '$slug' not found."
        exit 1
    fi

    local today
    if grep -q "^status:" "$file"; then
        sed -i -E 's/^status:[[:space:]]*["'"'"']?draft["'"'"']?/status: published/' "$file"
    elif grep -q "^draft:" "$file"; then
        sed -i -E 's/^draft:[[:space:]]*true/status: published/' "$file"
    fi
    sed -i -E "s/^date:[[:space:]]*[0-9]{4}-[0-9]{2}-[0-9]{2}/date: ${today}/" "$file"

    echo "[ok] Article '$slug' published with release date: $today."
}

cmd_draft() {
    local slug="${1:-}"
    if [ -z "$slug" ]; then
        echo "[error] Slug argument required. Usage: ./scripts/post.sh draft <slug>"
        exit 1
    fi

    local file
    file=$(find_post_file "$slug" || true)
    if [ -z "$file" ]; then
        echo "[error] Article with slug '$slug' not found."
        exit 1
    fi

    if grep -q "^status:" "$file"; then
        sed -i -E 's/^status:[[:space:]]*["'"'"']?published["'"'"']?/status: draft/' "$file"
    elif grep -q "^draft:" "$file"; then
        sed -i -E 's/^draft:[[:space:]]*false/status: draft/' "$file"
    fi
    echo "[ok] Article '$slug' reverted to draft."
}

cmd_archive() {
    local slug="${1:-}"
    if [ -z "$slug" ]; then
        echo "[error] Slug argument required. Usage: ./scripts/post.sh archive <slug>"
        exit 1
    fi

    local file
    file=$(find_post_file "$slug" || true)
    if [ -z "$file" ]; then
        echo "[error] Post with slug '$slug' not found."
        exit 1
    fi

    if [[ "$file" == *.md.archive ]]; then
        echo "[info] Post '$slug' is already archived."
        exit 0
    fi

    mv "$file" "${file}.archive"
    echo "[ok] Post '$slug' archived."
}

cmd_restore() {
    local slug="${1:-}"
    if [ -z "$slug" ]; then
        echo "[error] Slug argument required. Usage: ./scripts/post.sh restore <slug>"
        exit 1
    fi

    local file
    file=$(find_post_file "$slug" || true)
    if [ -z "$file" ]; then
        echo "[error] Post with slug '$slug' not found."
        exit 1
    fi

    if [[ "$file" != *.md.archive ]]; then
        echo "[info] Post '$slug' is already active."
        exit 0
    fi

    mv "$file" "${file%.archive}"
    echo "[ok] Post '$slug' restored."
}

cmd_delete() {
    local slug="${1:-}"
    if [ -z "$slug" ]; then
        echo "[error] Slug argument required. Usage: ./scripts/post.sh delete <slug>"
        exit 1
    fi

    local file
    file=$(find_post_file "$slug" || true)
    if [ -z "$file" ]; then
        echo "[error] Post with slug '$slug' not found."
        exit 1
    fi

    read -p "Permanently delete post '$slug' and all associated image assets? (y/N): " -r confirm
    if [[ "$confirm" =~ ^[Yy]$ ]]; then
        rm -f "$file"
        rm -rf "${IMG_BASE_DIR}/${slug}"
        echo "[ok] Post '$slug' and asset directory deleted permanently."
    else
        echo "[info] Deletion cancelled."
    fi
}
